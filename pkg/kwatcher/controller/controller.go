package controller

import (
	"context"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/kwatcher/config"
	"icm-varnish-k8s-operator/pkg/kwatcher/configmaps"
	"icm-varnish-k8s-operator/pkg/kwatcher/endpoints"
	"icm-varnish-k8s-operator/pkg/kwatcher/events"
	vslabels "icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/logger"
	"strings"

	"github.com/pkg/errors"

	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// PodInfo represents the relevant information of a pod for VCL code
type PodInfo struct {
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
		return errors.Wrap(err, "could not initialize controller")
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
		return errors.Wrap(err, "could not watch configMap")
	}

	backendEpPredicate, err := endpoints.NewPredicate(cfg.EndpointSelectorString, logr)
	if err != nil {
		return errors.Wrap(err, "could not create endpoints predicate")
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
		return errors.Wrap(err, "could not watch endpoints")
	}

	varnishPodsSelector := labels.SelectorFromSet(labels.Set{
		v1alpha1.LabelVarnishOwner:     cfg.VarnishServiceName,
		v1alpha1.LabelVarnishComponent: v1alpha1.VarnishComponentCacheService,
		v1alpha1.LabelVarnishUID:       string(cfg.VarnishServiceUID),
	})
	varnishEpPredicate, err := endpoints.NewPredicate(varnishPodsSelector.String(), logr)
	if err != nil {
		return errors.Wrap(err, "could not create endpoints predicate")
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
		return errors.Wrap(err, "could not watch endpoints")
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

		r.logger.Errorf("%+v", err)
		return reconcile.Result{}, err
	}
	return res, nil
}

func (r *ReconcileVarnish) reconcileWithLogging(request reconcile.Request) (reconcile.Result, error) {
	vs := &v1alpha1.VarnishService{}
	err := r.Get(context.Background(), types.NamespacedName{Namespace: request.Namespace, Name: r.config.VarnishServiceName}, vs)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	r.scheme.Default(vs)

	// ideally, I would pull these from the VarnishService since they could change with an update. For now, I am not doing that since I cannot default the values
	varnishPort := vs.Spec.Service.VarnishPort.Port
	targetPort := int32(vs.Spec.Service.VarnishPort.TargetPort.IntValue())
	entrypointFile := vs.Spec.VCLConfigMap.EntrypointFile

	pod := &v1.Pod{}
	err = r.Get(context.Background(), types.NamespacedName{Namespace: request.Namespace, Name: r.config.PodName}, pod)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	cm, err := r.getConfigMap(r.config.Namespace, r.config.ConfigMapName)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	newFiles, newTemplates := r.filesAndTemplates(cm.Data)

	if err = r.verifyEntrypointExists(newFiles, newTemplates, entrypointFile); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	bks, err := r.getPodInfo(r.config.Namespace, r.config.EndpointSelector, targetPort)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	varnishLabels := labels.SelectorFromSet(vslabels.CombinedComponentLabels(vs, v1alpha1.VarnishComponentCacheService))
	varnishNodes, err := r.getPodInfo(r.config.Namespace, varnishLabels, varnishPort)

	templatizedFiles, err := r.resolveTemplates(newTemplates, targetPort, varnishPort, bks, varnishNodes)
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

	currFiles, err := getCurrentFiles(r.config.VCLDir)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	filesTouched, err := r.reconcileFiles(r.config.VCLDir, currFiles, newFiles)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	if filesTouched {
		if err = r.reconcileVarnish(vs, pod, cm); err != nil {
			return reconcile.Result{}, errors.WithStack(err)
		}
	}

	if err := r.reconcilePod(filesTouched, pod, cm); err != nil {
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

func verifyFilesExist(configMapFiles map[string]string, files ...string) error {
	verify := func(filename string) error {
		if _, found := configMapFiles[filename]; !found {
			return errors.Errorf("%s must exist in configmap, but not found", filename)
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
