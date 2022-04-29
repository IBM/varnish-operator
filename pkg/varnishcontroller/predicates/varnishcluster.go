package predicates

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/logger"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &LabelMatcherPredicate{}

type varnishClusterPredicate struct {
	uuid   types.UID
	logger *logger.Logger
}

func NewVarnishClusterPredicate(uuid types.UID, logr *logger.Logger) predicate.Predicate {
	if logr == nil {
		logr = logger.NewNopLogger()
	}
	return &varnishClusterPredicate{
		logger: logr,
		uuid:   uuid,
	}
}

func (p *varnishClusterPredicate) Create(e event.CreateEvent) bool {
	vc, ok := e.Object.(*v1alpha1.VarnishCluster)
	if !ok {
		return true
	}

	if vc.GetUID() != p.uuid {
		return false
	}

	return true
}

func (p *varnishClusterPredicate) Delete(e event.DeleteEvent) bool {
	_, ok := e.Object.(*v1alpha1.VarnishCluster)
	if !ok {
		return false
	}

	return true
}

func (p *varnishClusterPredicate) Update(e event.UpdateEvent) bool {
	if reflect.TypeOf(e.ObjectNew) != reflect.TypeOf(e.ObjectOld) {
		return false
	}

	newCluster := e.ObjectNew.(*v1alpha1.VarnishCluster)
	oldCluster := e.ObjectOld.(*v1alpha1.VarnishCluster)
	if e.ObjectNew.GetUID() != p.uuid || e.ObjectOld.GetUID() != p.uuid {
		return false
	}

	if newCluster.Status.VCL.ConfigMapVersion != oldCluster.Status.VCL.ConfigMapVersion {
		return true
	}

	if !cmp.Equal(newCluster.Spec.Backend, oldCluster.Spec.Backend) {
		return true
	}

	return true
}

func (p *varnishClusterPredicate) Generic(e event.GenericEvent) bool {
	vc, ok := e.Object.(*v1alpha1.VarnishCluster)
	if !ok {
		return true
	}

	if vc.GetUID() != p.uuid {
		return false
	}

	return true
}
