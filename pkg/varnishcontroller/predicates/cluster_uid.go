package predicates

import (
	"icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &varnishClusterUIDPredicate{}

type varnishClusterUIDPredicate struct {
	clusterUID types.UID
	logger     *logger.Logger
}

func (ep *varnishClusterUIDPredicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *varnishClusterUIDPredicate) Delete(e event.DeleteEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *varnishClusterUIDPredicate) Update(e event.UpdateEvent) bool {
	newCluster, ok := e.ObjectNew.(*v1alpha1.VarnishCluster)
	if !ok {
		return true
	}
	oldCluster, ok := e.ObjectOld.(*v1alpha1.VarnishCluster)
	if !ok {
		ep.logger.Errorf("Wrong object type. Got %T Expected %T", e.ObjectNew, oldCluster)
		return false
	}

	if e.MetaNew.GetUID() != ep.clusterUID || e.MetaOld.GetUID() != ep.clusterUID {
		return false
	}

	if newCluster.Status.VCL.ConfigMapVersion != oldCluster.Status.VCL.ConfigMapVersion {
		return false
	}

	return true
}

func (ep *varnishClusterUIDPredicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *varnishClusterUIDPredicate) shared(meta metav1.Object, obj runtime.Object) bool {
	if _, ok := obj.(*v1alpha1.VarnishCluster); !ok {
		return true
	}

	if meta.GetUID() != ep.clusterUID {
		return false
	}

	return true
}

// NewVarnishClusterUIDPredicate filters out all VarnishClusters.
func NewVarnishClusterUIDPredicate(clusterUID types.UID, logger *logger.Logger) predicate.Predicate {
	return &varnishClusterUIDPredicate{
		clusterUID: clusterUID,
		logger:     logger,
	}
}
