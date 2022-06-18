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

package v1

import (
	"context"
	v1 "ekiOperator/apis/iaaseki/v1"
	scheme "ekiOperator/generated/iaaseki/clientset/versioned/scheme"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EkiMonitorsGetter has a method to return a EkiMonitorInterface.
// A group's client should implement this interface.
type EkiMonitorsGetter interface {
	EkiMonitors(namespace string) EkiMonitorInterface
}

// EkiMonitorInterface has methods to work with EkiMonitor resources.
type EkiMonitorInterface interface {
	Create(ctx context.Context, ekiMonitor *v1.EkiMonitor, opts metav1.CreateOptions) (*v1.EkiMonitor, error)
	Update(ctx context.Context, ekiMonitor *v1.EkiMonitor, opts metav1.UpdateOptions) (*v1.EkiMonitor, error)
	UpdateStatus(ctx context.Context, ekiMonitor *v1.EkiMonitor, opts metav1.UpdateOptions) (*v1.EkiMonitor, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.EkiMonitor, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.EkiMonitorList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.EkiMonitor, err error)
	EkiMonitorExpansion
}

// ekiMonitors implements EkiMonitorInterface
type ekiMonitors struct {
	client rest.Interface
	ns     string
}

// newEkiMonitors returns a EkiMonitors
func newEkiMonitors(c *IaasekiV1Client, namespace string) *ekiMonitors {
	return &ekiMonitors{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the ekiMonitor, and returns the corresponding ekiMonitor object, and an error if there is any.
func (c *ekiMonitors) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.EkiMonitor, err error) {
	result = &v1.EkiMonitor{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("ekimonitors").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EkiMonitors that match those selectors.
func (c *ekiMonitors) List(ctx context.Context, opts metav1.ListOptions) (result *v1.EkiMonitorList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.EkiMonitorList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("ekimonitors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested ekiMonitors.
func (c *ekiMonitors) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("ekimonitors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a ekiMonitor and creates it.  Returns the server's representation of the ekiMonitor, and an error, if there is any.
func (c *ekiMonitors) Create(ctx context.Context, ekiMonitor *v1.EkiMonitor, opts metav1.CreateOptions) (result *v1.EkiMonitor, err error) {
	result = &v1.EkiMonitor{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("ekimonitors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ekiMonitor).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a ekiMonitor and updates it. Returns the server's representation of the ekiMonitor, and an error, if there is any.
func (c *ekiMonitors) Update(ctx context.Context, ekiMonitor *v1.EkiMonitor, opts metav1.UpdateOptions) (result *v1.EkiMonitor, err error) {
	result = &v1.EkiMonitor{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("ekimonitors").
		Name(ekiMonitor.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ekiMonitor).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *ekiMonitors) UpdateStatus(ctx context.Context, ekiMonitor *v1.EkiMonitor, opts metav1.UpdateOptions) (result *v1.EkiMonitor, err error) {
	result = &v1.EkiMonitor{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("ekimonitors").
		Name(ekiMonitor.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ekiMonitor).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the ekiMonitor and deletes it. Returns an error if one occurs.
func (c *ekiMonitors) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("ekimonitors").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *ekiMonitors) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("ekimonitors").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched ekiMonitor.
func (c *ekiMonitors) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.EkiMonitor, err error) {
	result = &v1.EkiMonitor{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("ekimonitors").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
