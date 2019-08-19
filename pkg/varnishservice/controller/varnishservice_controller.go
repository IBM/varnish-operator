package controller

import (
	"context"
	icmv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"
	"icm-varnish-k8s-operator/pkg/varnishservice/config"
	"icm-varnish-k8s-operator/pkg/varnishservice/pods"
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"go.uber.org/zap"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func NewVarnishReconciler(mgr manager.Manager, cfg *config.Config, logr *logger.Logger) reconcile.Reconciler {
	return &ReconcileVarnishService{
		Client: mgr.GetClient(),
		logger: logr,
		config: cfg,
		scheme: mgr.GetScheme(),
		events: NewEventHandler(mgr.GetRecorder(EventRecorderNameVarnishService)),
	}
}

// Add creates a new VarnishService Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(r reconcile.Reconciler, mgr manager.Manager, logr *logger.Logger) error {
	// Create a new controller
	c, err := controller.New("varnishservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to VarnishService
	logr.Infow("Starting watch loop for VarnishService objects")
	err = c.Watch(&source.Kind{Type: &icmv1alpha1.VarnishService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes in the below resources:
	// Pods
	// ConfigMap
	// StatefulSet
	// Service
	// Role
	// RoleBinding
	// ServiceAccount

	podPredicate, err := pods.NewPredicate(map[string]string{icmv1alpha1.LabelVarnishComponent: icmv1alpha1.VarnishComponentVarnishes})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &v1.Pod{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(
			func(a handler.MapObject) []reconcile.Request {
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Name:      a.Meta.GetLabels()[icmv1alpha1.LabelVarnishOwner],
						Namespace: a.Meta.GetNamespace(),
					}},
				}
			}),
	}, podPredicate)
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbacv1beta1.Role{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbacv1beta1.RoleBinding{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbacv1beta1.ClusterRole{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &rbacv1beta1.ClusterRoleBinding{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.ServiceAccount{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileVarnishService{}

// ReconcileVarnishService reconciles a VarnishService object
type ReconcileVarnishService struct {
	client.Client
	logger *logger.Logger
	config *config.Config
	scheme *runtime.Scheme
	events *EventHandler
}

// Reconcile reads that state of the cluster for a VarnishService object and makes changes based on the state read
// and what is in the VarnishService.Spec
// Automatically generate RBAC rules to allow the Controller to read and write StatefulSets
// +kubebuilder:rbac:groups=icm.ibm.com,resources=varnishservices,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=icm.ibm.com,resources=varnishservices/status,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups="",resources=endpoints,verbs=list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=list;get;watch;update
// +kubebuilder:rbac:groups="",resources=nodes,verbs=watch;list
// +kubebuilder:rbac:groups="admissionregistration.k8s.io",resources="validatingwebhookconfigurations;mutatingwebhookconfigurations",verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=list;watch
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=list;watch;create;update;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings;clusterroles;clusterrolebindings,verbs=list;watch;create;update;delete
func (r *ReconcileVarnishService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
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
			return reconcile.Result{Requeue: true}, nil //retry but do not treat conflicts as errors
		}

		logr.Errorf("%+v", err)
		return reconcile.Result{}, err
	}
	return res, nil
}

func (r *ReconcileVarnishService) reconcileWithContext(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// Fetch the VarnishService instance
	instance := &icmv1alpha1.VarnishService{}
	err := r.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, errors.Wrap(err, "could not read VarnishService")
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
		return reconcile.Result{}, err
	}
	roleName, err := r.reconcileRole(ctx, instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.reconcileRoleBinding(ctx, instance, roleName, serviceAccountName); err != nil {
		return reconcile.Result{}, err
	}
	clusterRoleName, err := r.reconcileClusterRole(ctx, instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.reconcileClusterRoleBinding(ctx, instance, clusterRoleName, serviceAccountName); err != nil {
		return reconcile.Result{}, err
	}
	endpointSelector, err := r.reconcileServiceNoCache(ctx, instance, instanceStatus)
	if err != nil {
		return reconcile.Result{}, err
	}
	headlessServiceName, err := r.reconcileHeadlessService(ctx, instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	varnishSelector, err := r.reconcileStatefulSet(ctx, instance, instanceStatus, serviceAccountName, endpointSelector, headlessServiceName)
	if err != nil {
		return reconcile.Result{}, err
	}
	// TODO: remove extra return var
	if _, err = r.reconcileConfigMap(ctx, varnishSelector, instance, instanceStatus); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.reconcilePodDisruptionBudget(ctx, instance, varnishSelector); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.reconcileService(ctx, instance, instanceStatus, varnishSelector); err != nil {
		return reconcile.Result{}, err
	}

	if !compare.EqualVarnishServiceStatus(&instance.Status, &instanceStatus.Status) {
		logger.FromContext(ctx).Infoc("Updating VarnishService Status", "diff", compare.DiffVarnishServiceStatus(&instance.Status, &instanceStatus.Status))
		if err = r.Status().Update(ctx, instanceStatus); err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "could not update VarnishService Status %s:%s, %s:%s", "name", instance.Name, "namespace", instance.Namespace)
		}
	} else {
		logger.FromContext(ctx).Debugw("No updates for VarnishService status")
	}

	return reconcile.Result{}, nil
}
