package controller

import (
	"icm-varnish-k8s-operator/operator/controller/pkg/handlers"
	"time"

	"github.com/juju/errors"

	icmapiv1alpha1 "icm-varnish-k8s-operator/operator/controller/pkg/apis/icm/v1alpha1"
	vsclientset "icm-varnish-k8s-operator/operator/controller/pkg/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Controller describes everything the Varnish Service operator acts on
type Controller struct {
	clientset    vsclientset.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler handlers.Handler
	maxRetries   int
}

// New creates a new instance of Controller
func New(clientset vsclientset.Interface, eventHandler handlers.Handler, maxRetries int) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return clientset.IcmV1alpha1().VarnishServices(metav1.NamespaceAll).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return clientset.IcmV1alpha1().VarnishServices(metav1.NamespaceAll).Watch(options)
			},
		},
		&icmapiv1alpha1.VarnishService{},
		0, // skip resync
		cache.Indexers{},
	)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if key, err := cache.MetaNamespaceKeyFunc(obj); err == nil {
				queue.Add(CacheKey{Key: key, Act: Add})
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if key, err := cache.MetaNamespaceKeyFunc(newObj); err == nil {
				queue.Add(CacheKey{Key: key, Act: Update})
			}
		},
		DeleteFunc: func(obj interface{}) {
			if key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err == nil {
				queue.Add(CacheKey{Key: key, Act: Delete})
			}
		},
	})

	return &Controller{
		clientset:    clientset,
		queue:        queue,
		informer:     informer,
		eventHandler: &handlers.Default{},
		maxRetries:   maxRetries,
	}
}

// NewVarnishServiceController creates a new instance of Controller for a VarnishService
func NewVarnishServiceController(clientset vsclientset.Interface, maxRetries int) *Controller {
	return New(clientset, &handlers.VarnishServiceHandler{}, maxRetries)
}

// Run starts the Controller listening for its resource
func (c *Controller) Run(stopCh <-chan struct{}) {
	// catch panics
	defer utilruntime.HandleCrash()
	// shut down the queue when done
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	// wait for cache to sync before starting worker
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(errors.New("Timed out waiting for caches to sync"))
		return
	}

	wait.Until(c.runWorker, time.Second, stopCh)
}

// runWorker processes items in a loop until it returns false
func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

// processNextItem pulls a value off the queue and processes.
// Returns false when shutting down
func (c *Controller) processNextItem() bool {
	// pull next item off queue, which should be a key to lookup something in cache
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// indicate to the queue that work is complete
	defer c.queue.Done(key)

	err := c.processItem(key.(CacheKey))
	if err == nil {
		c.queue.Forget(key)
	} else if c.queue.NumRequeues(key) < c.maxRetries {
		c.queue.AddRateLimited(key)
	} else {
		c.queue.Forget(key)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(key CacheKey) error {
	obj, exists, err := c.informer.GetIndexer().GetByKey(key.Key)
	if err != nil {
		return errors.Errorf("Error fetching object with key %s from store: %v", key, err)
	}

	switch key.Act {
	case Add:
		return c.eventHandler.ObjectAdded(obj)
	case Update:
		return c.eventHandler.ObjectUpdated(obj)
	case Delete:
		return c.eventHandler.ObjectDeleted(obj)
	}
	if !exists {
		return c.eventHandler.ObjectDeleted(obj)
	}
	return c.eventHandler.ObjectAdded(obj)
}
