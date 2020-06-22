/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/blakelead/nsinjector/pkg/apis/namespaceresourcesinjector/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeNamespaceResourcesInjectors implements NamespaceResourcesInjectorInterface
type FakeNamespaceResourcesInjectors struct {
	Fake *FakeBlakeleadV1alpha1
}

var namespaceresourcesinjectorsResource = schema.GroupVersionResource{Group: "blakelead.com", Version: "v1alpha1", Resource: "namespaceresourcesinjectors"}

var namespaceresourcesinjectorsKind = schema.GroupVersionKind{Group: "blakelead.com", Version: "v1alpha1", Kind: "NamespaceResourcesInjector"}

// Get takes name of the namespaceResourcesInjector, and returns the corresponding namespaceResourcesInjector object, and an error if there is any.
func (c *FakeNamespaceResourcesInjectors) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.NamespaceResourcesInjector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(namespaceresourcesinjectorsResource, name), &v1alpha1.NamespaceResourcesInjector{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceResourcesInjector), err
}

// List takes label and field selectors, and returns the list of NamespaceResourcesInjectors that match those selectors.
func (c *FakeNamespaceResourcesInjectors) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.NamespaceResourcesInjectorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(namespaceresourcesinjectorsResource, namespaceresourcesinjectorsKind, opts), &v1alpha1.NamespaceResourcesInjectorList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.NamespaceResourcesInjectorList{ListMeta: obj.(*v1alpha1.NamespaceResourcesInjectorList).ListMeta}
	for _, item := range obj.(*v1alpha1.NamespaceResourcesInjectorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested namespaceResourcesInjectors.
func (c *FakeNamespaceResourcesInjectors) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(namespaceresourcesinjectorsResource, opts))
}

// Create takes the representation of a namespaceResourcesInjector and creates it.  Returns the server's representation of the namespaceResourcesInjector, and an error, if there is any.
func (c *FakeNamespaceResourcesInjectors) Create(ctx context.Context, namespaceResourcesInjector *v1alpha1.NamespaceResourcesInjector, opts v1.CreateOptions) (result *v1alpha1.NamespaceResourcesInjector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(namespaceresourcesinjectorsResource, namespaceResourcesInjector), &v1alpha1.NamespaceResourcesInjector{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceResourcesInjector), err
}

// Update takes the representation of a namespaceResourcesInjector and updates it. Returns the server's representation of the namespaceResourcesInjector, and an error, if there is any.
func (c *FakeNamespaceResourcesInjectors) Update(ctx context.Context, namespaceResourcesInjector *v1alpha1.NamespaceResourcesInjector, opts v1.UpdateOptions) (result *v1alpha1.NamespaceResourcesInjector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(namespaceresourcesinjectorsResource, namespaceResourcesInjector), &v1alpha1.NamespaceResourcesInjector{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceResourcesInjector), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNamespaceResourcesInjectors) UpdateStatus(ctx context.Context, namespaceResourcesInjector *v1alpha1.NamespaceResourcesInjector, opts v1.UpdateOptions) (*v1alpha1.NamespaceResourcesInjector, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(namespaceresourcesinjectorsResource, "status", namespaceResourcesInjector), &v1alpha1.NamespaceResourcesInjector{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceResourcesInjector), err
}

// Delete takes name of the namespaceResourcesInjector and deletes it. Returns an error if one occurs.
func (c *FakeNamespaceResourcesInjectors) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(namespaceresourcesinjectorsResource, name), &v1alpha1.NamespaceResourcesInjector{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNamespaceResourcesInjectors) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(namespaceresourcesinjectorsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.NamespaceResourcesInjectorList{})
	return err
}

// Patch applies the patch and returns the patched namespaceResourcesInjector.
func (c *FakeNamespaceResourcesInjectors) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NamespaceResourcesInjector, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(namespaceresourcesinjectorsResource, name, pt, data, subresources...), &v1alpha1.NamespaceResourcesInjector{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceResourcesInjector), err
}