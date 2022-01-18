package controller

import (
	"context"
	"time"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/config"
	vcreconcile "github.com/ibm/varnish-operator/pkg/varnishcluster/reconcile"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	annotationVarnishClusterName      = "varnish-cluster-name"
	annotationVarnishClusterNamespace = "varnish-cluster-namespace"
)

func SetupVarnishReconciler(ctx context.Context, vcCtrl reconcile.Reconciler, mgr manager.Manager, reconcileChan chan event.GenericEvent) error {
	clusterRoleBindingEventHandler := handler.EnqueueRequestsFromMapFunc(func(a client.Object) []ctrl.Request {
		cr, ok := a.(*rbac.ClusterRoleBinding)
		if !ok {
			return nil
		}

		if cr.Annotations[annotationVarnishClusterNamespace] == "" {
			return nil
		}

		if cr.Annotations[annotationVarnishClusterName] == "" {
			return nil
		}

		return []ctrl.Request{
			{NamespacedName: types.NamespacedName{
				Name:      cr.Annotations[annotationVarnishClusterName],
				Namespace: cr.Annotations[annotationVarnishClusterNamespace],
			}},
		}
	})

	clusterRoleEventHandler := handler.EnqueueRequestsFromMapFunc(func(a client.Object) []ctrl.Request {
		cr, ok := a.(*rbac.ClusterRole)
		if !ok {
			return nil
		}

		if cr.Annotations[annotationVarnishClusterNamespace] == "" {
			return nil
		}

		if cr.Annotations[annotationVarnishClusterName] == "" {
			return nil
		}

		return []ctrl.Request{
			{NamespacedName: types.NamespacedName{
				Name:      cr.Annotations[annotationVarnishClusterName],
				Namespace: cr.Annotations[annotationVarnishClusterNamespace],
			}},
		}
	})

	vcPodsSelector := labels.SelectorFromSet(map[string]string{vcapi.LabelVarnishComponent: vcapi.VarnishComponentVarnish})
	varnishClusterPodsEventHandler := handler.EnqueueRequestsFromMapFunc(func(a client.Object) []ctrl.Request {
		if !vcPodsSelector.Matches(labels.Set(a.GetLabels())) {
			return nil
		}

		return []ctrl.Request{
			{NamespacedName: types.NamespacedName{
				Name:      a.GetLabels()[vcapi.LabelVarnishOwner],
				Namespace: a.GetNamespace(),
			}},
		}
	})

	builder := ctrl.NewControllerManagedBy(mgr)
	builder.Named("varnishcluster")
	builder.For(&vcapi.VarnishCluster{})
	builder.Owns(&v1.ConfigMap{})
	builder.Owns(&appsv1.StatefulSet{})
	builder.Owns(&v1.Service{})
	builder.Owns(&rbac.Role{})
	builder.Owns(&rbac.RoleBinding{})
	builder.Watches(&source.Kind{Type: &rbac.ClusterRole{}}, clusterRoleEventHandler)
	builder.Watches(&source.Kind{Type: &rbac.ClusterRoleBinding{}}, clusterRoleBindingEventHandler)
	builder.Owns(&v1.ServiceAccount{})
	builder.Watches(&source.Channel{Source: reconcileChan}, &handler.EnqueueRequestForObject{})
	builder.Watches(&source.Kind{Type: &v1.Pod{}}, varnishClusterPodsEventHandler)

	serviceMonitorList := &unstructured.UnstructuredList{}
	serviceMonitorList.SetGroupVersionKind(serviceMonitorListGVK)
	err := mgr.GetClient().List(ctx, serviceMonitorList)
	if err != nil {
		if _, ok := errors.Cause(err).(*meta.NoKindMatchError); ok {
			logger.FromContext(ctx).Warn("Can't watch ServiceMonitor. ServiceMonitor Kind is not found. Prometheus operator needs to be installed first.", err)
		} else {
			logger.FromContext(ctx).Error("Can't watch ServiceMonitor: %s", err)
			//the return is intentionally omitted. Better work without that watch than not at all
		}
	} else {
		serviceMonitor := &unstructured.Unstructured{}
		serviceMonitor.SetGroupVersionKind(serviceMonitorGVK)
		builder.Owns(serviceMonitor)
	}

	return builder.Complete(vcCtrl)
}

var _ reconcile.Reconciler = &ReconcileVarnishCluster{}

// ReconcileVarnishCluster reconciles a VarnishCluster object
type ReconcileVarnishCluster struct {
	client.Client
	logger             *logger.Logger
	config             *config.Config
	scheme             *runtime.Scheme
	events             *EventHandler
	reconcileTriggerer *vcreconcile.ReconcileTriggerer
}

func NewVarnishReconciler(mgr manager.Manager, cfg *config.Config, logr *logger.Logger, reconcileChan chan event.GenericEvent) *ReconcileVarnishCluster {
	return &ReconcileVarnishCluster{
		Client:             mgr.GetClient(),
		logger:             logr,
		config:             cfg,
		scheme:             mgr.GetScheme(),
		events:             NewEventHandler(mgr.GetEventRecorderFor(EventRecorderNameVarnishCluster)),
		reconcileTriggerer: vcreconcile.NewReconcileTriggerer(logr, reconcileChan),
	}
}

// Reconcile reads that state of the cluster for a VarnishCluster object and makes changes based on the state read
// and what is in the VarnishCluster.Spec
// Automatically generate RBAC rules to allow the Controller to read and write StatefulSets
// +kubebuilder:rbac:groups=caching.ibm.com,resources=varnishclusters,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=caching.ibm.com,resources=varnishclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=caching.ibm.com,resources=varnishclusters/finalizers,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups="",resources=endpoints;namespaces,verbs=list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=list;get;watch;update
// +kubebuilder:rbac:groups="",resources=nodes,verbs=watch;list
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings;clusterroles;clusterrolebindings,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;delete

func (r *ReconcileVarnishCluster) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	logr := r.logger.With(logger.FieldVarnishCluster, request.Name)
	logr = logr.With(logger.FieldNamespace, request.Namespace)
	ctx = logger.ToContext(ctx, logr)

	logr.Debugw("Reconciling...")
	start := time.Now()
	defer logr.Debugf("Reconciled in %s", time.Since(start).String())
	res, err := r.reconcileWithContext(ctx, request)
	if err != nil {
		if statusErr, ok := errors.Cause(err).(*apierrors.StatusError); ok && statusErr.ErrStatus.Reason == metav1.StatusReasonConflict {
			logr.Info("Conflict occurred. Retrying...", zap.Error(err))
			return ctrl.Result{Requeue: true}, nil //retry but do not treat conflicts as errors
		}

		logr.Errorf("%+v", err)
		return ctrl.Result{}, err
	}
	return res, nil
}

func (r *ReconcileVarnishCluster) reconcileWithContext(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	// Fetch the VarnishCluster instance
	instance := &vcapi.VarnishCluster{}
	err := r.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// There were situations when the VarnishCluser object has been deleted before the finalisation logic is executed below.
			// So make sure the resources supposed to be cleaned up by the finalizer are removed.
			if err := r.deleteCR(ctx, types.NamespacedName{Name: names.ClusterRole(request.Name, request.Namespace)}); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.deleteCRB(ctx, types.NamespacedName{Name: names.ClusterRoleBinding(request.Name, request.Namespace)}); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, errors.Wrap(err, "could not read VarnishCluster")
	}

	r.scheme.Default(instance)

	// For some reason, sometimes Kubernetes returns the object without apiVersion and kind
	// Since the code below relies on that values we set them manually if they are empty
	if instance.APIVersion == "" {
		instance.APIVersion = "caching.ibm.com/v1alpha1"
	}
	if instance.Kind == "" {
		instance.Kind = "VarnishCluster"
	}

	// If VarnishCluster is deleted by the user, but since we have finalizers Kubernetes only marks it as deleted.
	// We need to do our cleanup logic and delete the finalizers in order to let Kubernetes to delete the object from etcd
	// and garbage collect all owned objects.
	// The rest of the code is designed to not be executed when the resource is marked for deletion.
	if !instance.DeletionTimestamp.IsZero() {
		err = r.finalizerCleanUp(ctx, instance)
		return ctrl.Result{}, err
	}

	err = r.reconcileFinalizers(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	instanceStatus := &vcapi.VarnishCluster{}
	instance.ObjectMeta.DeepCopyInto(&instanceStatus.ObjectMeta)
	instance.Status.DeepCopyInto(&instanceStatus.Status)

	err = r.reconcileServiceAccount(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.reconcileRole(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err = r.reconcileRoleBinding(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}
	err = r.reconcileClusterRole(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err = r.reconcileClusterRoleBinding(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}
	endpointSelector, err := r.reconcileServiceNoCache(ctx, instance, instanceStatus)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.reconcileHeadlessService(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err = r.reconcileVarnishSecret(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}
	sts, varnishSelector, err := r.reconcileStatefulSet(ctx, instance, instanceStatus, endpointSelector)
	if err != nil {
		return ctrl.Result{}, err
	}
	// TODO: remove extra return var
	if _, err = r.reconcileConfigMap(ctx, varnishSelector, instance, instanceStatus); err != nil {
		return ctrl.Result{}, err
	}

	if err = r.reconcilePodDisruptionBudget(ctx, instance, varnishSelector); err != nil {
		return ctrl.Result{}, err
	}
	if err = r.reconcileService(ctx, instance, instanceStatus, varnishSelector); err != nil {
		return ctrl.Result{}, err
	}

	if err = r.reconcileDelayedRollingUpdate(ctx, instance, instanceStatus, sts); err != nil {
		return ctrl.Result{}, err
	}

	if err = r.reconcileServiceMonitor(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	if err = r.reconcileGrafanaDashboard(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	if !compare.EqualVarnishClusterStatus(&instance.Status, &instanceStatus.Status) {
		logger.FromContext(ctx).Infoc("Updating VarnishCluster Status", "diff", compare.DiffVarnishClusterStatus(&instance.Status, &instanceStatus.Status))
		if err = r.Status().Update(ctx, instanceStatus); err != nil {
			return ctrl.Result{}, errors.Wrapf(err, "could not update VarnishCluster Status %s:%s, %s:%s", "name", instance.Name, "namespace", instance.Namespace)
		}
	} else {
		logger.FromContext(ctx).Debugw("No updates for VarnishCluster status")
	}

	return ctrl.Result{}, nil
}
