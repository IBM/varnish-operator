package watch

import (
	"time"

	"github.com/juju/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

// objTypeFor takes a resource string, and returns the struct representing the string
// Not necessarily a complete list.
func objTypeFor(resource string) (runtime.Object, error) {
	switch resource {
	case v1.ResourcePods.String():
		return &v1.Pod{}, nil
	case v1.ResourceServices.String():
		return &v1.Service{}, nil
	case v1.ResourceReplicationControllers.String():
		return &v1.ReplicationController{}, nil
	case v1.ResourceQuotas.String():
		return &v1.ResourceQuota{}, nil
	case v1.ResourceSecrets.String():
		return &v1.Secret{}, nil
	case v1.ResourceConfigMaps.String():
		return &v1.ConfigMap{}, nil
	case v1.ResourcePersistentVolumeClaims.String():
		return &v1.PersistentVolumeClaim{}, nil
	case "endpoints":
		return &v1.Endpoints{}, nil
	default:
		return nil, errors.Errorf("resource %s not in list.", resource)
	}
}

// WatchResource watches the named resource for changes and runs the onChange function with every change.
// For add events, onChange will pass in (nil, newObj), and for delete events, onChange will pass in (oldObj, nil)
func WatchResource(c cache.Getter, resource string, namespace string, selector fields.Selector, onChange func(oldObj, newObj interface{})) (cache.Store, cache.Controller, error) {
	objType, err := objTypeFor(resource)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	listWatch := cache.NewListWatchFromClient(c, resource, namespace, selector)

	handlerFns := cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { onChange(nil, obj) },
		DeleteFunc: func(obj interface{}) { onChange(obj, nil) },
		UpdateFunc: onChange,
	}

	store, controller := cache.NewInformer(listWatch, objType, time.Second*0, handlerFns)
	return store, controller, nil
}
