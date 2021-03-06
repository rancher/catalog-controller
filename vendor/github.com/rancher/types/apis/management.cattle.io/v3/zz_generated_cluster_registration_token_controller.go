package v3

import (
	"context"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	ClusterRegistrationTokenGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ClusterRegistrationToken",
	}
	ClusterRegistrationTokenResource = metav1.APIResource{
		Name:         "clusterregistrationtokens",
		SingularName: "clusterregistrationtoken",
		Namespaced:   true,

		Kind: ClusterRegistrationTokenGroupVersionKind.Kind,
	}
)

type ClusterRegistrationTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterRegistrationToken
}

type ClusterRegistrationTokenHandlerFunc func(key string, obj *ClusterRegistrationToken) error

type ClusterRegistrationTokenLister interface {
	List(namespace string, selector labels.Selector) (ret []*ClusterRegistrationToken, err error)
	Get(namespace, name string) (*ClusterRegistrationToken, error)
}

type ClusterRegistrationTokenController interface {
	Informer() cache.SharedIndexInformer
	Lister() ClusterRegistrationTokenLister
	AddHandler(name string, handler ClusterRegistrationTokenHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ClusterRegistrationTokenInterface interface {
	ObjectClient() *clientbase.ObjectClient
	Create(*ClusterRegistrationToken) (*ClusterRegistrationToken, error)
	GetNamespace(name, namespace string, opts metav1.GetOptions) (*ClusterRegistrationToken, error)
	Get(name string, opts metav1.GetOptions) (*ClusterRegistrationToken, error)
	Update(*ClusterRegistrationToken) (*ClusterRegistrationToken, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespace(name, namespace string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ClusterRegistrationTokenList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ClusterRegistrationTokenController
	AddHandler(name string, sync ClusterRegistrationTokenHandlerFunc)
	AddLifecycle(name string, lifecycle ClusterRegistrationTokenLifecycle)
}

type clusterRegistrationTokenLister struct {
	controller *clusterRegistrationTokenController
}

func (l *clusterRegistrationTokenLister) List(namespace string, selector labels.Selector) (ret []*ClusterRegistrationToken, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*ClusterRegistrationToken))
	})
	return
}

func (l *clusterRegistrationTokenLister) Get(namespace, name string) (*ClusterRegistrationToken, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    ClusterRegistrationTokenGroupVersionKind.Group,
			Resource: "clusterRegistrationToken",
		}, name)
	}
	return obj.(*ClusterRegistrationToken), nil
}

type clusterRegistrationTokenController struct {
	controller.GenericController
}

func (c *clusterRegistrationTokenController) Lister() ClusterRegistrationTokenLister {
	return &clusterRegistrationTokenLister{
		controller: c,
	}
}

func (c *clusterRegistrationTokenController) AddHandler(name string, handler ClusterRegistrationTokenHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*ClusterRegistrationToken))
	})
}

type clusterRegistrationTokenFactory struct {
}

func (c clusterRegistrationTokenFactory) Object() runtime.Object {
	return &ClusterRegistrationToken{}
}

func (c clusterRegistrationTokenFactory) List() runtime.Object {
	return &ClusterRegistrationTokenList{}
}

func (s *clusterRegistrationTokenClient) Controller() ClusterRegistrationTokenController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.clusterRegistrationTokenControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ClusterRegistrationTokenGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &clusterRegistrationTokenController{
		GenericController: genericController,
	}

	s.client.clusterRegistrationTokenControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type clusterRegistrationTokenClient struct {
	client       *Client
	ns           string
	objectClient *clientbase.ObjectClient
	controller   ClusterRegistrationTokenController
}

func (s *clusterRegistrationTokenClient) ObjectClient() *clientbase.ObjectClient {
	return s.objectClient
}

func (s *clusterRegistrationTokenClient) Create(o *ClusterRegistrationToken) (*ClusterRegistrationToken, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*ClusterRegistrationToken), err
}

func (s *clusterRegistrationTokenClient) Get(name string, opts metav1.GetOptions) (*ClusterRegistrationToken, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*ClusterRegistrationToken), err
}

func (s *clusterRegistrationTokenClient) GetNamespace(name, namespace string, opts metav1.GetOptions) (*ClusterRegistrationToken, error) {
	obj, err := s.objectClient.GetNamespace(name, namespace, opts)
	return obj.(*ClusterRegistrationToken), err
}

func (s *clusterRegistrationTokenClient) Update(o *ClusterRegistrationToken) (*ClusterRegistrationToken, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*ClusterRegistrationToken), err
}

func (s *clusterRegistrationTokenClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *clusterRegistrationTokenClient) DeleteNamespace(name, namespace string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespace(name, namespace, options)
}

func (s *clusterRegistrationTokenClient) List(opts metav1.ListOptions) (*ClusterRegistrationTokenList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ClusterRegistrationTokenList), err
}

func (s *clusterRegistrationTokenClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *clusterRegistrationTokenClient) Patch(o *ClusterRegistrationToken, data []byte, subresources ...string) (*ClusterRegistrationToken, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*ClusterRegistrationToken), err
}

func (s *clusterRegistrationTokenClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *clusterRegistrationTokenClient) AddHandler(name string, sync ClusterRegistrationTokenHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *clusterRegistrationTokenClient) AddLifecycle(name string, lifecycle ClusterRegistrationTokenLifecycle) {
	sync := NewClusterRegistrationTokenLifecycleAdapter(name, s, lifecycle)
	s.AddHandler(name, sync)
}
