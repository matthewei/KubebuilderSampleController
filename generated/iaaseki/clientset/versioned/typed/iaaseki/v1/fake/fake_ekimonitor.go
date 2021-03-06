/*
Copyright 2022.

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
	iaasekiv1 "ekiOperator/apis/iaaseki/v1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEkiMonitors implements EkiMonitorInterface
type FakeEkiMonitors struct {
	Fake *FakeIaasekiV1
	ns   string
}

var ekimonitorsResource = schema.GroupVersionResource{Group: "iaaseki.cmss", Version: "v1", Resource: "ekimonitors"}

var ekimonitorsKind = schema.GroupVersionKind{Group: "iaaseki.cmss", Version: "v1", Kind: "EkiMonitor"}

// Get takes name of the ekiMonitor, and returns the corresponding ekiMonitor object, and an error if there is any.
func (c *FakeEkiMonitors) Get(ctx context.Context, name string, options v1.GetOptions) (result *iaasekiv1.EkiMonitor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(ekimonitorsResource, c.ns, name), &iaasekiv1.EkiMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*iaasekiv1.EkiMonitor), err
}

// List takes label and field selectors, and returns the list of EkiMonitors that match those selectors.
func (c *FakeEkiMonitors) List(ctx context.Context, opts v1.ListOptions) (result *iaasekiv1.EkiMonitorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(ekimonitorsResource, ekimonitorsKind, c.ns, opts), &iaasekiv1.EkiMonitorList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &iaasekiv1.EkiMonitorList{ListMeta: obj.(*iaasekiv1.EkiMonitorList).ListMeta}
	for _, item := range obj.(*iaasekiv1.EkiMonitorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested ekiMonitors.
func (c *FakeEkiMonitors) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(ekimonitorsResource, c.ns, opts))

}

// Create takes the representation of a ekiMonitor and creates it.  Returns the server's representation of the ekiMonitor, and an error, if there is any.
func (c *FakeEkiMonitors) Create(ctx context.Context, ekiMonitor *iaasekiv1.EkiMonitor, opts v1.CreateOptions) (result *iaasekiv1.EkiMonitor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(ekimonitorsResource, c.ns, ekiMonitor), &iaasekiv1.EkiMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*iaasekiv1.EkiMonitor), err
}

// Update takes the representation of a ekiMonitor and updates it. Returns the server's representation of the ekiMonitor, and an error, if there is any.
func (c *FakeEkiMonitors) Update(ctx context.Context, ekiMonitor *iaasekiv1.EkiMonitor, opts v1.UpdateOptions) (result *iaasekiv1.EkiMonitor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(ekimonitorsResource, c.ns, ekiMonitor), &iaasekiv1.EkiMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*iaasekiv1.EkiMonitor), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeEkiMonitors) UpdateStatus(ctx context.Context, ekiMonitor *iaasekiv1.EkiMonitor, opts v1.UpdateOptions) (*iaasekiv1.EkiMonitor, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(ekimonitorsResource, "status", c.ns, ekiMonitor), &iaasekiv1.EkiMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*iaasekiv1.EkiMonitor), err
}

// Delete takes name of the ekiMonitor and deletes it. Returns an error if one occurs.
func (c *FakeEkiMonitors) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(ekimonitorsResource, c.ns, name, opts), &iaasekiv1.EkiMonitor{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEkiMonitors) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(ekimonitorsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &iaasekiv1.EkiMonitorList{})
	return err
}

// Patch applies the patch and returns the patched ekiMonitor.
func (c *FakeEkiMonitors) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *iaasekiv1.EkiMonitor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(ekimonitorsResource, c.ns, name, pt, data, subresources...), &iaasekiv1.EkiMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*iaasekiv1.EkiMonitor), err
}
