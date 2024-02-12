/*
Copyright 2020 The Knative Authors

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

	v1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSidecars implements SidecarInterface
type FakeSidecars struct {
	Fake *FakeNetworkingV1beta1
	ns   string
}

var sidecarsResource = v1beta1.SchemeGroupVersion.WithResource("sidecars")

var sidecarsKind = v1beta1.SchemeGroupVersion.WithKind("Sidecar")

// Get takes name of the sidecar, and returns the corresponding sidecar object, and an error if there is any.
func (c *FakeSidecars) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Sidecar, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(sidecarsResource, c.ns, name), &v1beta1.Sidecar{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Sidecar), err
}

// List takes label and field selectors, and returns the list of Sidecars that match those selectors.
func (c *FakeSidecars) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.SidecarList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(sidecarsResource, sidecarsKind, c.ns, opts), &v1beta1.SidecarList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.SidecarList{ListMeta: obj.(*v1beta1.SidecarList).ListMeta}
	for _, item := range obj.(*v1beta1.SidecarList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sidecars.
func (c *FakeSidecars) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(sidecarsResource, c.ns, opts))

}

// Create takes the representation of a sidecar and creates it.  Returns the server's representation of the sidecar, and an error, if there is any.
func (c *FakeSidecars) Create(ctx context.Context, sidecar *v1beta1.Sidecar, opts v1.CreateOptions) (result *v1beta1.Sidecar, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(sidecarsResource, c.ns, sidecar), &v1beta1.Sidecar{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Sidecar), err
}

// Update takes the representation of a sidecar and updates it. Returns the server's representation of the sidecar, and an error, if there is any.
func (c *FakeSidecars) Update(ctx context.Context, sidecar *v1beta1.Sidecar, opts v1.UpdateOptions) (result *v1beta1.Sidecar, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(sidecarsResource, c.ns, sidecar), &v1beta1.Sidecar{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Sidecar), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSidecars) UpdateStatus(ctx context.Context, sidecar *v1beta1.Sidecar, opts v1.UpdateOptions) (*v1beta1.Sidecar, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(sidecarsResource, "status", c.ns, sidecar), &v1beta1.Sidecar{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Sidecar), err
}

// Delete takes name of the sidecar and deletes it. Returns an error if one occurs.
func (c *FakeSidecars) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(sidecarsResource, c.ns, name, opts), &v1beta1.Sidecar{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSidecars) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(sidecarsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.SidecarList{})
	return err
}

// Patch applies the patch and returns the patched sidecar.
func (c *FakeSidecars) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Sidecar, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(sidecarsResource, c.ns, name, pt, data, subresources...), &v1beta1.Sidecar{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Sidecar), err
}