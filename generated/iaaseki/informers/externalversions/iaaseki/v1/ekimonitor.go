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
// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	iaasekiv1 "ekiOperator/apis/iaaseki/v1"
	versioned "ekiOperator/generated/iaaseki/clientset/versioned"
	internalinterfaces "ekiOperator/generated/iaaseki/informers/externalversions/internalinterfaces"
	v1 "ekiOperator/generated/iaaseki/listers/iaaseki/v1"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// EkiMonitorInformer provides access to a shared informer and lister for
// EkiMonitors.
type EkiMonitorInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.EkiMonitorLister
}

type ekiMonitorInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewEkiMonitorInformer constructs a new informer for EkiMonitor type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewEkiMonitorInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredEkiMonitorInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredEkiMonitorInformer constructs a new informer for EkiMonitor type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredEkiMonitorInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.IaasekiV1().EkiMonitors(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.IaasekiV1().EkiMonitors(namespace).Watch(context.TODO(), options)
			},
		},
		&iaasekiv1.EkiMonitor{},
		resyncPeriod,
		indexers,
	)
}

func (f *ekiMonitorInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredEkiMonitorInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *ekiMonitorInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&iaasekiv1.EkiMonitor{}, f.defaultInformer)
}

func (f *ekiMonitorInformer) Lister() v1.EkiMonitorLister {
	return v1.NewEkiMonitorLister(f.Informer().GetIndexer())
}
