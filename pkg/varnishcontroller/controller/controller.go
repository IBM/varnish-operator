package controller

import (
	"context"
	"github.com/ibm/varnish-operator/api/v1alpha1"
	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/config"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/events"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/metrics"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/predicates"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/varnishadm"
	"strings"
	"time"

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
	r := &ReconcileVarnish{
		config:       cfg,
		logger:       logr,
		Client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		varnish:      varnish,
		eventHandler: events.NewEventHandler(mgr.GetEventRecorderFor(events.EventRecorderName), cfg.PodName),
		metrics:      metrics,
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

	builder.For(&v1alpha1.VarnishCluster{})
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
	builder.WithEventFilter(predicates.NewVarnishControllerPredicate(cfg.VarnishClusterUID, endpointsSelectors, nil))

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
	metrics      *metrics.VarnishControllerMetrics
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

	varnishPort := int32(v1alpha1.VarnishPort)
	entrypointFileName := *vc.Spec.VCL.EntrypointFileName

	pod := &v1.Pod{}
	err = r.Get(ctx, types.NamespacedName{Namespace: request.Namespace, Name: r.config.PodName}, pod)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
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

	if err := r.reconcilePod(ctx, filesTouched, pod, cm, localWeight, remoteWeight); err != nil {
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
