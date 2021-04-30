package controller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/blakelead/nsinjector/internal/config"
	"github.com/blakelead/nsinjector/internal/utils"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/blakelead/nsinjector/pkg/apis/namespaceresourcesinjector/v1alpha1"
	nriinformers "github.com/blakelead/nsinjector/pkg/generated/informers/externalversions/namespaceresourcesinjector/v1alpha1"
	coreinformers "k8s.io/client-go/informers/core/v1"
)

// Controller structure
type Controller struct {
	nriInformer       nriinformers.NamespaceResourcesInjectorInformer
	namespaceInformer coreinformers.NamespaceInformer
	clients           *config.Clients
	queue             workqueue.RateLimitingInterface
	syncs             []cache.InformerSynced
}

// NamespaceEvent is an Enum
type NamespaceEvent string

// Values associated to NamespaceEvent
const (
	Added   NamespaceEvent = "created"
	Deleted NamespaceEvent = "deleted"
)

// Key represents the objects stored in the queue
type Key struct {
	Value string
	// We need that in order to know which namespace triggered the worker
	Invoker string
	// Wether namespace where created or deleted
	Event NamespaceEvent
}

// NewController creates a new Controller
func NewController(clients *config.Clients, factories *config.InformerFactories) *Controller {

	nriInformer := factories.NRI.Blakelead().V1alpha1().NamespaceResourcesInjectors()
	namespaceInformer := factories.Kube.Core().V1().Namespaces()
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	syncs := []cache.InformerSynced{
		nriInformer.Informer().HasSynced,
		namespaceInformer.Informer().HasSynced,
	}

	c := &Controller{
		nriInformer:       nriInformer,
		namespaceInformer: namespaceInformer,
		clients:           clients,
		queue:             queue,
		syncs:             syncs,
	}

	nriInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAddNRI,
		UpdateFunc: c.onUpdateNRI,
		DeleteFunc: c.onDeleteNRI,
	})

	namespaceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAddNamespace,
		DeleteFunc: c.onDeleteNamespace,
	})

	return c
}

// Run wait for cache sync then starts the workers
func (c *Controller) Run(workers int, stop <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	log.Info("Starting controller")

	if !cache.WaitForCacheSync(stop, c.syncs...) {
		return fmt.Errorf("Failed to wait for caches to sync")
	}

	for i := 0; i < workers; i++ {
		go wait.Until(func() {
			for c.processNextWorkItem() {
			}
		}, time.Second, stop)
	}

	<-stop
	log.Info("Shutting down controller")

	return nil
}

func (c *Controller) onAddNRI(obj interface{}) {
	if key, err := cache.MetaNamespaceKeyFunc(obj); err == nil {
		log.Infof("Added %s to queue because it was created", key)
		c.queue.Add(Key{Value: key})
	}
}

func (c *Controller) onUpdateNRI(old, new interface{}) {
	oldSpec := (old.(*v1alpha1.NamespaceResourcesInjector)).Spec
	newSpec := (new.(*v1alpha1.NamespaceResourcesInjector)).Spec
	if reflect.DeepEqual(oldSpec, newSpec) {
		return
	}
	if key, err := cache.MetaNamespaceKeyFunc(old); err == nil {
		log.Infof("Added %s to queue because it was updated", key)
		c.queue.Add(Key{Value: key})
	}
}

func (c *Controller) onDeleteNRI(obj interface{}) {
	if key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err == nil {
		log.Infof("Added %s to queue because it was deleted", key)
		c.queue.Add(Key{Value: key})
	}
}

// When a namespace is created/deleted the controller should update
// nri that match the namespace

func (c *Controller) onAddNamespace(obj interface{}) {
	ns := obj.(*corev1.Namespace)
	nrilist, err := c.nriInformer.Lister().List(labels.Everything())
	if err != nil {
		log.Error(err)
		return
	}
	for _, nri := range nrilist {
		if nri.CanInject(ns.GetName()) {
			log.Infof("Added %s to queue because namespace %s was %s", nri.GetName(), ns.GetName(), NamespaceEvent(Added))
			c.queue.Add(Key{Value: nri.GetName(), Invoker: ns.GetName(), Event: NamespaceEvent(Added)})
		}
	}
}

func (c *Controller) onDeleteNamespace(obj interface{}) {
	ns := obj.(*corev1.Namespace)
	nrilist, err := c.nriInformer.Lister().List(labels.Everything())
	if err != nil {
		log.Error(err)
		return
	}
	for _, nri := range nrilist {
		if nri.CanInject(ns.GetName()) {
			log.Infof("Added %s to queue because namespace %s was %s", nri.GetName(), ns.GetName(), NamespaceEvent(Added))
			c.queue.Add(Key{Value: nri.GetName(), Invoker: ns.GetName(), Event: NamespaceEvent(Deleted)})
		}
	}
}

func (c *Controller) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)
	err := c.processItem(key.(Key))
	if err == nil {
		c.queue.Forget(key)
		return true
	}
	log.Errorf("Key '%s' processing failed: %v", key.(Key).Value, err)
	c.queue.AddRateLimited(key)
	return true
}

func (c *Controller) processItem(key Key) error {
	// get injector object
	injector, err := c.nriInformer.Lister().Get(key.Value)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Infof("Key '%s' in work queue no longer exists", key.Value)
			return nil
		}
		return err
	}

	// remove deleted namespace from injector status list
	if key.Event == NamespaceEvent(Deleted) {
		err = c.removeFromStatus(injector, key.Invoker)
		if err != nil {
			return err
		}
		return nil
	}

	// get all namespaces that should be injected
	namespaces, err := c.getNamespaces(injector)
	if err != nil {
		return err
	}
	if len(namespaces) == 0 {
		return nil
	}

	for _, namespace := range namespaces {
		// inject resources into namespace
		if err = c.inject(injector, namespace); err != nil {
			return err
		}
		// update injector status with newly injected namespace
		if err := c.addToStatus(injector, namespace); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) getNamespaces(injector *v1alpha1.NamespaceResourcesInjector) ([]*corev1.Namespace, error) {
	namespaces, err := c.namespaceInformer.Lister().List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var filteredNamespaces []*corev1.Namespace
	for _, namespace := range namespaces {
		if injector.CanInject(namespace.GetName()) {
			filteredNamespaces = append(filteredNamespaces, namespace)
		}
	}
	return filteredNamespaces, nil
}

func (c *Controller) inject(injector *v1alpha1.NamespaceResourcesInjector, namespace *corev1.Namespace) error {
	obj := &unstructured.Unstructured{}
	serializer := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(c.clients.Discovery))
	for _, resource := range injector.Spec.Resources {
		_, gvk, err := serializer.Decode([]byte(resource), nil, obj)
		if err != nil {
			return err
		}
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}
		var dynamicResource dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			dynamicResource = c.clients.Dynamic.Resource(mapping.Resource).Namespace(namespace.GetName())
		} else {
			dynamicResource = c.clients.Dynamic.Resource(mapping.Resource)
		}
		obj.SetOwnerReferences(
			[]metav1.OwnerReference{
				*metav1.NewControllerRef(injector, v1alpha1.SchemeGroupVersion.WithKind("NamespaceResourcesInjector")),
			},
		)
		_, err = dynamicResource.Get(context.Background(), obj.GetName(), metav1.GetOptions{})
		if err != nil {
			o, err := dynamicResource.Create(context.Background(), obj, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			log.Infof("Resource %s/%s in namespace %s from injector %s created", o.GetKind(), o.GetName(), o.GetNamespace(), injector.GetName())
		} else {
			o, err := dynamicResource.Update(context.Background(), obj, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
			log.Infof("Resource %s/%s in namespace %s from injector %s updated", o.GetKind(), o.GetName(), o.GetNamespace(), injector.GetName())
		}
	}
	return nil
}

func (c *Controller) addToStatus(injector *v1alpha1.NamespaceResourcesInjector, namespace *corev1.Namespace) error {
	injectorCopy := injector.DeepCopy()
	if !injector.Injected(namespace.GetName()) {
		injectorCopy.Status.InjectedNamespaces = append(injectorCopy.Status.InjectedNamespaces, namespace.GetName())
	}
	_, err := c.clients.NRI.BlakeleadV1alpha1().NamespaceResourcesInjectors().Update(
		context.Background(),
		injectorCopy,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return err
	}
	log.Infof("Added namespace %s to injector %s status", namespace.GetName(), injector.GetName())
	return nil
}

func (c *Controller) removeFromStatus(injector *v1alpha1.NamespaceResourcesInjector, namespace string) error {
	if injector.Injected(namespace) {
		injectorCopy := injector.DeepCopy()
		injectorCopy.Status.InjectedNamespaces = utils.Remove(injectorCopy.Status.InjectedNamespaces, namespace)
		_, err := c.clients.NRI.BlakeleadV1alpha1().NamespaceResourcesInjectors().Update(
			context.Background(),
			injectorCopy,
			metav1.UpdateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Removed namespace %s from injector %s status", namespace, injector.GetName())
	}
	return nil
}
