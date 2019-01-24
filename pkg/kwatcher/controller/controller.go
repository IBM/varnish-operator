package controller

import (
	"context"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/kwatcher/config"
	"icm-varnish-k8s-operator/pkg/kwatcher/configmaps"
	"icm-varnish-k8s-operator/pkg/kwatcher/endpoints"
	"icm-varnish-k8s-operator/pkg/kwatcher/events"
	"icm-varnish-k8s-operator/pkg/logger"

	"go.uber.org/zap"

	"github.com/juju/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

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

	epPredicate, err := endpoints.Predicate(cfg.EndpointSelectorString)
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
	}, epPredicate)
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
	res, err := r.reconcileWithLogging(request)
	if err != nil {
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

	backendsFile, err := resolveTemplate(newFiles[r.config.BackendsTmplFile], r.config.TargetPort, bks)
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
