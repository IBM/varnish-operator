package controller

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	ctrlBuilder "sigs.k8s.io/controller-runtime/pkg/builder"

	"github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/config"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/events"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/metrics"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/predicates"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/varnishadm"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// PodInfo represents the relevant information of a pod for VCL code
type PodInfo struct {
	IP         string
	NodeLabels map[string]string
	PodName    string
	Weight     float64
}

// SetupVarnishReconciler creates a new VarnishCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func SetupVarnishReconciler(mgr manager.Manager, cfg *config.Config, varnish varnishadm.VarnishAdministrator, metrics *metrics.VarnishControllerMetrics, logr *logger.Logger) error {
	// stub, backends selector will be set and updated on reconcile
	backendsSelector := labels.SelectorFromSet(labels.Set{})
	backendNamespacePredicate := predicates.NewNamespacesMatcherPredicate([]string{cfg.Namespace}, logr)
	backendLabelsPredicate := predicates.NewLabelMatcherPredicate(backendsSelector, logr)

	r := &ReconcileVarnish{
		config:                     cfg,
		logger:                     logr,
		Client:                     mgr.GetClient(),
		scheme:                     mgr.GetScheme(),
		varnish:                    varnish,
		eventHandler:               events.NewEventHandler(mgr.GetEventRecorderFor(events.EventRecorderName), cfg.PodName),
		metrics:                    metrics,
		backendsSelectorPredicate:  backendLabelsPredicate,
		backendsNamespacePredicate: backendNamespacePredicate,
	}

	podMapFunc := handler.EnqueueRequestsFromMapFunc(
		func(a client.Object) []reconcile.Request {
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Namespace: cfg.Namespace,
					Name:      cfg.PodName,
				}},
			}
		})

	builder := ctrl.NewControllerManagedBy(mgr)
	builder.Named("varnish-controller")

	builder.For(&v1alpha1.VarnishCluster{}, ctrlBuilder.WithPredicates(predicates.NewVarnishClusterPredicate(r.config.VarnishClusterUID, logr)))

	builder.Watches(
		&source.Kind{Type: &v1.Pod{}},
		podMapFunc,
		ctrlBuilder.WithPredicates(
			backendNamespacePredicate,
			backendLabelsPredicate,
		),
	)

	varnishPodsSelector := labels.SelectorFromSet(labels.Set{
		v1alpha1.LabelVarnishOwner:     cfg.VarnishClusterName,
		v1alpha1.LabelVarnishComponent: v1alpha1.VarnishComponentVarnish,
		v1alpha1.LabelVarnishUID:       string(cfg.VarnishClusterUID),
	})
	builder.Watches(
		&source.Kind{Type: &v1.Pod{}},
		podMapFunc,
		ctrlBuilder.WithPredicates(
			predicates.NewNamespacesMatcherPredicate([]string{cfg.Namespace}, logr),
			predicates.NewLabelMatcherPredicate(varnishPodsSelector, logr),
		),
	)
	//builder.WithEventFilter(predicates.NewDebugPredicate(logr))

	return builder.Complete(r)
}

var _ reconcile.Reconciler = &ReconcileVarnish{}

type ReconcileVarnish struct {
	client.Client
	config                     *config.Config
	logger                     *logger.Logger
	scheme                     *runtime.Scheme
	eventHandler               *events.EventHandler
	varnish                    varnishadm.VarnishAdministrator
	metrics                    *metrics.VarnishControllerMetrics
	backendsNamespacePredicate *predicates.NamespacesMatcherPredicate
	backendsSelectorPredicate  *predicates.LabelMatcherPredicate
}

func (r *ReconcileVarnish) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logr := r.logger.With(logger.FieldVarnishCluster, r.config.VarnishClusterName)
	logr = logr.With(logger.FieldPodName, r.config.PodName)
	logr = logr.With(logger.FieldNamespace, r.config.Namespace)

	logr.Debugw("Reconciling...")
	start := time.Now()
	defer logr.Debugf("Reconciled in %s", time.Since(start).String())
	ctx = logger.ToContext(ctx, logr)
	res, err := r.reconcileWithContext(ctx, request)
	if err != nil {
		if statusErr, ok := errors.Cause(err).(*apierrors.StatusError); ok && statusErr.ErrStatus.Reason == metav1.StatusReasonConflict {
			logr.Info("Conflict occurred. Retrying...", zap.Error(err))
			return reconcile.Result{Requeue: true}, nil //retry but do not treat conflicts as errors
		}

		logr.Errorf("%+v", err)
		return reconcile.Result{}, err
	}
	return res, nil
}

func (r *ReconcileVarnish) reconcileWithContext(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logr := logger.FromContext(ctx)
	vc := &v1alpha1.VarnishCluster{}
	err := r.Get(ctx, types.NamespacedName{Namespace: request.Namespace, Name: r.config.VarnishClusterName}, vc)
	if err != nil {
		// Happens when VarnishCluster is already deleted but the pods are still alive.
		// They are going to be deleted soon by the StatefulSet controller.
		if apierrors.IsNotFound(err) {
			logr.Infow("VarnishCluster is not found. Aborting, as the pod is going to be deleted soon.")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, errors.WithStack(err)
	}

	r.scheme.Default(vc)

	if len(vc.Spec.Backend.Namespaces) > 0 {
		r.backendsNamespacePredicate.Namespaces = vc.Spec.Backend.Namespaces
	} else {
		r.backendsNamespacePredicate.Namespaces = []string{r.config.Namespace}
	}
	r.backendsSelectorPredicate.Selector = labels.SelectorFromSet(vc.Spec.Backend.Selector)

	varnishPort := int32(v1alpha1.VarnishPort)
	entrypointFileName := *vc.Spec.VCL.EntrypointFileName

	pod := &v1.Pod{}
	err = r.Get(ctx, types.NamespacedName{Namespace: request.Namespace, Name: r.config.PodName}, pod)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	var haproxyConfigMap *v1.ConfigMap
	if vc.Spec.HaproxySidecar.Enabled {
		haproxyConfigMap, err = r.reconcileHaproxyConfig(ctx, vc)
		if err != nil {
			return reconcile.Result{}, errors.WithStack(err)
		}
	}

	cm, err := r.getConfigMap(ctx, r.config.Namespace, *vc.Spec.VCL.ConfigMapName)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	newFiles, newTemplates := r.filesAndTemplates(cm.Data)

	if err = r.verifyEntrypointExists(newFiles, newTemplates, entrypointFileName); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	bks, backendPortNumber, localWeight, remoteWeight, err := r.getBackendEndpoints(ctx, vc)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	varnishNodes, err := r.getVarnishEndpoints(ctx, vc)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	templatizedFiles, err := r.resolveTemplates(newTemplates, backendPortNumber, varnishPort, bks, varnishNodes)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	for fileName, contents := range templatizedFiles {
		if _, found := newFiles[fileName]; found {
			errMsg := fmt.Sprintf("VCL ConfigMap %s has %s and %s.tmpl entries. Cannot include file and template with same name",
				*vc.Spec.VCL.ConfigMapName, fileName, fileName)
			r.eventHandler.Warning(vc, events.EventReasonInvalidVCLConfigMap, errMsg)
			return reconcile.Result{}, errors.Errorf(errMsg)
		}
		newFiles[fileName] = contents
	}

	currFiles, err := getCurrentFiles(config.VCLConfigDir)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	filesTouched, err := r.reconcileFiles(ctx, config.VCLConfigDir, currFiles, newFiles)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	configName, err := r.varnish.GetActiveConfigurationName()
	if err != nil {
		return reconcile.Result{}, err
	}

	// reload if files changed, or we didn't load the VCL yet (happens when only the container restarted and not the whole pod)
	if filesTouched || configName == "boot" {
		if err = r.reconcileVarnish(ctx, vc, pod, cm); err != nil {
			return reconcile.Result{}, errors.WithStack(err)
		}
	}

	if err := r.reconcilePod(ctx, filesTouched, pod, cm, haproxyConfigMap, localWeight, remoteWeight); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileVarnish) templatizeHaproxyConfig(instance *v1alpha1.VarnishCluster, tmpl string) (string, error) {
	t, err := template.New("haproxy-config").Parse(tmpl)
	if err != nil {
		return "", errors.WithStack(err)
	}
	var b bytes.Buffer
	if err := t.Execute(&b, instance.Spec.HaproxySidecar); err != nil {
		return "", errors.WithStack(err)
	}
	return b.String(), nil
}

func (r *ReconcileVarnish) updateHaproxyConfig(haproxyConfigFileName string, cfgData string) (bool, error) {
	if b, err := os.ReadFile(haproxyConfigFileName); err != nil {
		if strings.Compare(string(b), cfgData) != 0 {
			if err := os.WriteFile(haproxyConfigFileName, []byte(cfgData), 0644); err != nil {
				return false, err
			}
			return true, nil
		}
	} else {
		if err := os.WriteFile(haproxyConfigFileName, []byte(cfgData), 0644); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (r *ReconcileVarnish) hupHaproxy() error {
	cmdString := "/haproxy-scripts/haproxy-hup.sh"
	cmd := exec.Command(cmdString)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileVarnish) reconcileHaproxyConfig(ctx context.Context, vc *v1alpha1.VarnishCluster) (*v1.ConfigMap, error) {
	logr := logger.FromContext(ctx)
	haproxyConfigMap, err := r.getConfigMap(ctx, r.config.Namespace, names.HaproxyConfigMap(vc.Name))
	if err != nil {
		return nil, err
	}

	cfgData, err := r.templatizeHaproxyConfig(vc, haproxyConfigMap.Data[v1alpha1.HaproxyConfigFileName])
	if err != nil {
		return nil, err
	}

	haproxyConfigFileName := v1alpha1.HaproxyConfigDir + "/" + v1alpha1.HaproxyConfigFileName
	haproxyConfigUpdated, err := r.updateHaproxyConfig(haproxyConfigFileName, cfgData)
	if err != nil {
		return nil, err
	}

	if haproxyConfigUpdated {
		logr.Info("hup'ing haproxy")
		if err := r.hupHaproxy(); err != nil {
			return nil, err
		}
	}
	return haproxyConfigMap, nil
}

func (r *ReconcileVarnish) filesAndTemplates(data map[string]string) (files, templates map[string]string) {
	files = make(map[string]string, len(data))
	templates = make(map[string]string)
	for fileName, contents := range data {
		if strings.HasSuffix(fileName, ".tmpl") {
			templates[fileName] = contents
		} else {
			files[fileName] = contents
		}
	}
	return
}

func (r *ReconcileVarnish) verifyEntrypointExists(files, templates map[string]string, entrypoint string) error {
	_, fileFound := files[entrypoint]
	_, templateFound := templates[entrypoint+".tmpl"]
	if !fileFound && !templateFound {
		return errors.Errorf("%s must exist in configmap, but not found", entrypoint)
	}
	return nil
}
