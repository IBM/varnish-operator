package varnishservice

import (
	"context"
	"errors"
	icmv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/compare"
	"icm-varnish-k8s-operator/pkg/logger"

	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new VarnishService Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileVarnishService{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		events: NewEventHandler(mgr.GetRecorder(EventRecorderNameVarnishService)),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("varnishservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to VarnishService
	logger.Infow("Starting watch loop for VarnishService objects")
	err = c.Watch(&source.Kind{Type: &icmv1alpha1.VarnishService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes in the below resources:
	// ConfigMap
	// Deployment
	// Service
	// Role
	// RoleBinding
	// ServiceAccount

	err = c.Watch(&source.Kind{Type: &v1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &icmv1alpha1.VarnishService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
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
	scheme *runtime.Scheme
	events *EventHandler
}

// Reconcile reads that state of the cluster for a VarnishService object and makes changes based on the state read
// and what is in the VarnishService.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=icm.ibm.com,resources=varnishservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=icm.ibm.com,resources=varnishservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=endpoints,verbs=list;watch
func (r *ReconcileVarnishService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the VarnishService instance
	instance := &icmv1alpha1.VarnishService{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, logger.RErrorw(err, "could not read VarnishService")
	}

	r.scheme.Default(instance)
	logr := logger.With("name", instance.Name, "namespace", instance.Namespace)

	instanceStatus := &icmv1alpha1.VarnishService{}
	instance.ObjectMeta.DeepCopyInto(&instanceStatus.ObjectMeta)
	instance.Status.DeepCopyInto(&instanceStatus.Status)

	configmapSelector, err := r.reconcileConfigMap(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	serviceAccountName, err := r.reconcileServiceAccount(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	roleName, err := r.reconcileRole(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.reconcileRoleBinding(instance, roleName, serviceAccountName); err != nil {
		return reconcile.Result{}, err
	}
	applicationPort, err := getApplicationPort(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	endpointSelector, err := r.reconcileNoCachedService(instance, instanceStatus, applicationPort)
	if err != nil {
		return reconcile.Result{}, err
	}
	varnishSelector, err := r.reconcileDeployment(instance, instanceStatus, serviceAccountName, applicationPort, endpointSelector, configmapSelector)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.reconcileCachedService(instance, instanceStatus, applicationPort, varnishSelector); err != nil {
		return reconcile.Result{}, err
	}

	if !compare.EqualVarnishServiceStatus(&instance.Status, &instanceStatus.Status) {
		logr.Debugw("Updating VarnishService Status", "diff", compare.DiffVarnishServiceStatus(&instance.Status, &instanceStatus.Status))
		if err = r.Status().Update(context.TODO(), instanceStatus); err != nil {
			return reconcile.Result{}, logger.RErrorw(err, "could not update VarnishService Status", "name", instance.Name, "namespace", instance.Namespace)
		}
	} else {
		logr.Debugw("No updates for VarnishService status")
	}

	return reconcile.Result{}, nil
}

func getApplicationPort(instance *icmv1alpha1.VarnishService) (*v1.ServicePort, error) {
	logr := logger.With("name", instance.Name, "namespace", instance.Namespace)

	if len(instance.Spec.Service.Ports) != 1 {
		err := errors.New("must specify exactly one port in service spec")
		logr.RErrorw(err, "")
		return nil, err
	}
	port := instance.Spec.Service.Ports[0]
	if port.TargetPort == (intstr.IntOrString{}) {
		port.TargetPort = intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: port.Port,
		}
	}
	return &port, nil
}
