package controller

import (
	"context"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/kwatcher/config"
	"icm-varnish-k8s-operator/pkg/kwatcher/configmaps"
	"icm-varnish-k8s-operator/pkg/kwatcher/endpoints"
	"icm-varnish-k8s-operator/pkg/kwatcher/events"
	"icm-varnish-k8s-operator/pkg/logger"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.uber.org/zap"

	"github.com/juju/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// VarnishNode represents a Varnish pod
type VarnishNode struct {
	IP      string
	Port    int32
	PodName string
}

// Backend represents pods that are cached by Varnish
type Backend struct {
	IP         string
	NodeLabels map[string]string
	PodName    string
}

// Add creates a new VarnishService Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, cfg *config.Config, logr *logger.Logger) error {
	r := &ReconcileVarnish{
		config:       cfg,
		logger:       logr,
		Client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		eventHandler: events.NewEventHandler(mgr.GetRecorder(events.EventRecorderName), cfg.PodName),
	}

	// Create a new controller
	c, err := controller.New("varnishservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return errors.Annotate(err, "could not initialize controller")
	}

	err = c.Watch(&source.Kind{Type: &v1.ConfigMap{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(
			func(a handler.MapObject) []reconcile.Request {
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Namespace: cfg.Namespace,
						Name:      cfg.PodName,
					}},
				}
			}),
	}, configmaps.Predicate(cfg.ConfigMapName))
	if err != nil {
		return errors.Annotate(err, "could not watch configMap")
	}

	backendEpPredicate, err := endpoints.NewPredicate(cfg.EndpointSelectorString, logr)
	if err != nil {
		return errors.Annotate(err, "could not create endpoints predicate")
	}
	err = c.Watch(&source.Kind{Type: &v1.Endpoints{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(
			func(a handler.MapObject) []reconcile.Request {
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Namespace: cfg.Namespace,
						Name:      cfg.PodName,
					}},
				}
			}),
	}, backendEpPredicate)
	if err != nil {
		return errors.Annotate(err, "could not watch endpoints")
	}

	varnishPodsSelector := labels.SelectorFromSet(labels.Set{
		v1alpha1.LabelVarnishOwner:     cfg.VarnishServiceName,
		v1alpha1.LabelVarnishComponent: v1alpha1.VarnishComponentCachedService,
		v1alpha1.LabelVarnishUID:       string(cfg.VarnishServiceUID),
	})
	varnishEpPredicate, err := endpoints.NewPredicate(varnishPodsSelector.String(), logr)
	if err != nil {
		return errors.Annotate(err, "could not create endpoints predicate")
	}
	err = c.Watch(&source.Kind{Type: &v1.Endpoints{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(
			func(a handler.MapObject) []reconcile.Request {
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Namespace: cfg.Namespace,
						Name:      cfg.PodName,
					}},
				}
			}),
	}, varnishEpPredicate)
	if err != nil {
		return errors.Annotate(err, "could not watch endpoints")
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileVarnish{}

type ReconcileVarnish struct {
	client.Client
	config       *config.Config
	logger       *logger.Logger
	scheme       *runtime.Scheme
	eventHandler *events.EventHandler
}

func (r *ReconcileVarnish) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	r.logger.Debugw("Reconciling...")
	res, err := r.reconcileWithLogging(request)
	if err != nil {
		if statusErr, ok := errors.Cause(err).(*apierrors.StatusError); ok && statusErr.ErrStatus.Reason == metav1.StatusReasonConflict {
			r.logger.Info("Conflict occurred. Retrying...", zap.Error(err))
			return reconcile.Result{Requeue: true}, nil //retry but do not treat conflicts as errors
		}

		r.logger.Error(zap.Error(err))
		return reconcile.Result{}, err
	}
	return res, nil
}

func (r *ReconcileVarnish) reconcileWithLogging(request reconcile.Request) (reconcile.Result, error) {
	vs := &v1alpha1.VarnishService{}
	err := r.Get(context.Background(), types.NamespacedName{Namespace: request.Namespace, Name: r.config.VarnishServiceName}, vs)
	if err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	pod := &v1.Pod{}
	err = r.Get(context.Background(), types.NamespacedName{Namespace: request.Namespace, Name: r.config.PodName}, pod)
	if err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	cm, err := r.getConfigMap(r.config.Namespace, r.config.ConfigMapName)
	if err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	newFiles := make(map[string][]byte, len(cm.Data))
	for k, v := range cm.Data {
		newFiles[k] = []byte(v)
	}

	if err = verifyFilesExist(newFiles, r.config.DefaultFile, r.config.BackendsTmplFile); err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	bks, err := r.getBackends(r.config.Namespace, r.config.EndpointSelector, r.config.TargetPort)
	if err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	var varnishNodes []VarnishNode
	if varnishNodes, err = r.getVarnishNodes(vs); err != nil {
		return reconcile.Result{}, err
	}

	backendsFile, err := r.resolveTemplate(newFiles[r.config.BackendsTmplFile], r.config.TargetPort, bks, varnishNodes)
	if err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	delete(newFiles, r.config.BackendsTmplFile)
	newFiles[r.config.BackendsFile] = backendsFile

	currFiles, err := getCurrentFiles(r.config.VCLDir)
	if err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	filesTouched, err := r.reconcileFiles(r.config.VCLDir, currFiles, newFiles)
	if err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	if filesTouched {
		if err = r.reconcileVarnish(vs, pod, cm); err != nil {
			return reconcile.Result{}, errors.Trace(err)
		}
	}

	if err := r.reconcilePod(filesTouched, pod, cm); err != nil {
		return reconcile.Result{}, errors.Trace(err)
	}

	return reconcile.Result{}, nil
}

func verifyFilesExist(configMapFiles map[string][]byte, files ...string) error {
	verify := func(filename string) error {
		if _, found := configMapFiles[filename]; !found {
			return errors.NotFoundf("%s must exist in configmap, but not found", filename)
		}
		return nil
	}

	for _, file := range files {
		if err := verify(file); err != nil {
			return err
		}
	}
	return nil
}
