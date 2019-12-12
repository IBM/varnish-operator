package controller

import (
	"context"
	"icm-varnish-k8s-operator/api/v1alpha1"
	vclabels "icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishcontroller/config"
	"icm-varnish-k8s-operator/pkg/varnishcontroller/events"
	"icm-varnish-k8s-operator/pkg/varnishcontroller/predicates"
	"icm-varnish-k8s-operator/pkg/varnishcontroller/varnishadm"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/intstr"

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
}

// SetupVarnishReconciler creates a new VarnishCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func SetupVarnishReconciler(mgr manager.Manager, cfg *config.Config, varnish varnishadm.VarnishAdministrator, logr *logger.Logger) error {
	r := &ReconcileVarnish{
		config:       cfg,
		logger:       logr,
		Client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		varnish:      varnish,
		eventHandler: events.NewEventHandler(mgr.GetEventRecorderFor(events.EventRecorderName), cfg.PodName),
	}

	podMapFunc := &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(
			func(a handler.MapObject) []reconcile.Request {
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Namespace: cfg.Namespace,
						Name:      cfg.PodName,
					}},
				}
			})}

	builder := ctrl.NewControllerManagedBy(mgr)
	builder.Named("varnish-controller")

	// Normally the `For` function receives the type the operator owns. For the operator it's the VarnishCluster.
	// For varnish-controller we don't have such resource but still need to provide something for the `For` function, it will fail otherwise.
	// For that reason we provide the v1alpha1.VarnishCluster resource but filter it all out. The operator will receive events from Kubernetes but won't react.
	builder.For(&v1alpha1.VarnishCluster{})
	builder.WithEventFilter(predicates.NewIgnoreVarnishClusterPredicate())

	builder.Watches(&source.Kind{Type: &v1.ConfigMap{}}, podMapFunc)
	builder.WithEventFilter(predicates.NewConfigMapNamePredicate(cfg.ConfigMapName, logr))

	builder.Watches(&source.Kind{Type: &v1.Endpoints{}}, podMapFunc)

	backendsLabels, err := labels.ConvertSelectorToLabelsMap(cfg.EndpointSelectorString)
	if err != nil {
		return err
	}

	varnishPodsSelector := labels.SelectorFromSet(labels.Set{
		v1alpha1.LabelVarnishOwner:     cfg.VarnishClusterName,
		v1alpha1.LabelVarnishComponent: v1alpha1.VarnishComponentCacheService,
		v1alpha1.LabelVarnishUID:       string(cfg.VarnishClusterUID),
	})

	endpointsSelectors := []labels.Selector{labels.SelectorFromSet(backendsLabels), varnishPodsSelector}
	builder.WithEventFilter(predicates.NewEndpointsSelectors(endpointsSelectors, logr))

	return builder.Complete(r)
}

var _ reconcile.Reconciler = &ReconcileVarnish{}

type ReconcileVarnish struct {
	client.Client
	config       *config.Config
	logger       *logger.Logger
	scheme       *runtime.Scheme
	eventHandler *events.EventHandler
	varnish      varnishadm.VarnishAdministrator
}

func (r *ReconcileVarnish) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

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
	vc := &v1alpha1.VarnishCluster{}
	err := r.Get(ctx, types.NamespacedName{Namespace: request.Namespace, Name: r.config.VarnishClusterName}, vc)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	r.scheme.Default(vc)

	varnishPort := int32(v1alpha1.VarnishPort)
	targetPort := *vc.Spec.Backend.Port
	entrypointFileName := *vc.Spec.VCL.EntrypointFileName

	pod := &v1.Pod{}
	err = r.Get(ctx, types.NamespacedName{Namespace: request.Namespace, Name: r.config.PodName}, pod)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	cm, err := r.getConfigMap(ctx, r.config.Namespace, r.config.ConfigMapName)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	newFiles, newTemplates := r.filesAndTemplates(cm.Data)

	if err = r.verifyEntrypointExists(newFiles, newTemplates, entrypointFileName); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	bks, backendPortNumber, err := r.getPodInfo(ctx, r.config.Namespace, r.config.EndpointSelector, targetPort)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	varnishLabels := labels.SelectorFromSet(vclabels.CombinedComponentLabels(vc, v1alpha1.VarnishComponentCacheService))
	varnishNodes, _, err := r.getPodInfo(ctx, r.config.Namespace, varnishLabels, intstr.FromString(v1alpha1.VarnishPortName))
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	templatizedFiles, err := r.resolveTemplates(newTemplates, backendPortNumber, varnishPort, bks, varnishNodes)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	for fileName, contents := range templatizedFiles {
		if _, found := newFiles[fileName]; found {
			// TODO: probably want to create event for this
			return reconcile.Result{}, errors.Errorf("ConfigMap has %s and %s.tmpl entries. Cannot include file and template with same name", fileName, fileName)
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

	if filesTouched {
		if err = r.reconcileVarnish(ctx, vc, pod, cm); err != nil {
			return reconcile.Result{}, errors.WithStack(err)
		}
	}

	if err := r.reconcilePod(ctx, filesTouched, pod, cm); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	return reconcile.Result{}, nil
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
