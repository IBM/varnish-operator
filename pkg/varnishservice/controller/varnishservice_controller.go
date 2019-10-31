package controller

import (
	"context"
	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"
	"icm-varnish-k8s-operator/pkg/varnishservice/config"
	vsreconcile "icm-varnish-k8s-operator/pkg/varnishservice/reconcile"
	"time"

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

func SetupVarnishReconciler(mgr manager.Manager, cfg *config.Config, logr *logger.Logger, reconcileChan chan event.GenericEvent) error {
	vsCtrl := &ReconcileVarnishService{
		Client:             mgr.GetClient(),
		logger:             logr,
		config:             cfg,
		scheme:             mgr.GetScheme(),
		events:             NewEventHandler(mgr.GetEventRecorderFor(EventRecorderNameVarnishService)),
		reconcileTriggerer: vsreconcile.NewReconcileTriggerer(logr, reconcileChan),
	}

	builder := ctrl.NewControllerManagedBy(mgr)
	builder.Named("varnishservice")
	builder.For(&icmv1alpha1.VarnishService{})
	builder.Owns(&v1.ConfigMap{})
	builder.Owns(&appsv1.StatefulSet{})
	builder.Owns(&v1.Service{})
	builder.Owns(&rbac.Role{})
	builder.Owns(&rbac.RoleBinding{})
	builder.Owns(&rbac.ClusterRole{})
	builder.Owns(&rbac.ClusterRoleBinding{})
	builder.Owns(&v1.ServiceAccount{})
	builder.Watches(&source.Channel{Source: reconcileChan}, &handler.EnqueueRequestForObject{})

	vsPodsSelector := labels.SelectorFromSet(map[string]string{icmv1alpha1.LabelVarnishComponent: icmv1alpha1.VarnishComponentVarnish})
	builder.Watches(&source.Kind{Type: &v1.Pod{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(
			func(a handler.MapObject) []ctrl.Request {
				if !vsPodsSelector.Matches(labels.Set(a.Meta.GetLabels())) {
					return nil
				}

				return []ctrl.Request{
					{NamespacedName: types.NamespacedName{
						Name:      a.Meta.GetLabels()[icmv1alpha1.LabelVarnishOwner],
						Namespace: a.Meta.GetNamespace(),
					}},
				}
			}),
	})

	return builder.Complete(vsCtrl)
}

var _ reconcile.Reconciler = &ReconcileVarnishService{}

// ReconcileVarnishService reconciles a VarnishService object
type ReconcileVarnishService struct {
	client.Client
	logger             *logger.Logger
	config             *config.Config
	scheme             *runtime.Scheme
	events             *EventHandler
	reconcileTriggerer *vsreconcile.ReconcileTriggerer
}

// Reconcile reads that state of the cluster for a VarnishService object and makes changes based on the state read
// and what is in the VarnishService.Spec
// Automatically generate RBAC rules to allow the Controller to read and write StatefulSets
// +kubebuilder:rbac:groups=icm.ibm.com,resources=varnishservices,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=icm.ibm.com,resources=varnishservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups="",resources=endpoints,verbs=list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=list;get;watch;update
// +kubebuilder:rbac:groups="",resources=nodes,verbs=watch;list
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings;clusterroles;clusterrolebindings,verbs=list;watch;create;update;delete

func (r *ReconcileVarnishService) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	logr := r.logger.With(logger.FieldVarnishService, request.Name)
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

func (r *ReconcileVarnishService) reconcileWithContext(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	// Fetch the VarnishService instance
	instance := &icmv1alpha1.VarnishService{}
	err := r.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, errors.Wrap(err, "could not read VarnishService")
	}

	r.scheme.Default(instance)

	// For some reason, sometimes Kubernetes returns the object without apiVersion and kind
	// Since the code below relies on that values we set them manually if they are empty
	if instance.APIVersion == "" {
		instance.APIVersion = "icm.ibm.com/v1alpha1"
	}
	if instance.Kind == "" {
		instance.Kind = "VarnishService"
	}

	instanceStatus := &icmv1alpha1.VarnishService{}
	instance.ObjectMeta.DeepCopyInto(&instanceStatus.ObjectMeta)
	instance.Status.DeepCopyInto(&instanceStatus.Status)

	serviceAccountName, err := r.reconcileServiceAccount(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	roleName, err := r.reconcileRole(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err = r.reconcileRoleBinding(ctx, instance, roleName, serviceAccountName); err != nil {
		return ctrl.Result{}, err
	}
	clusterRoleName, err := r.reconcileClusterRole(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err = r.reconcileClusterRoleBinding(ctx, instance, clusterRoleName, serviceAccountName); err != nil {
		return ctrl.Result{}, err
	}
	endpointSelector, err := r.reconcileServiceNoCache(ctx, instance, instanceStatus)
	if err != nil {
		return ctrl.Result{}, err
	}
	headlessServiceName, err := r.reconcileHeadlessService(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	sts, varnishSelector, err := r.reconcileStatefulSet(ctx, instance, instanceStatus, serviceAccountName, endpointSelector, headlessServiceName)
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

	if !compare.EqualVarnishServiceStatus(&instance.Status, &instanceStatus.Status) {
		logger.FromContext(ctx).Infoc("Updating VarnishService Status", "diff", compare.DiffVarnishServiceStatus(&instance.Status, &instanceStatus.Status))
		if err = r.Status().Update(ctx, instanceStatus); err != nil {
			return ctrl.Result{}, errors.Wrapf(err, "could not update VarnishService Status %s:%s, %s:%s", "name", instance.Name, "namespace", instance.Namespace)
		}
	} else {
		logger.FromContext(ctx).Debugw("No updates for VarnishService status")
	}

	return ctrl.Result{}, nil
}
