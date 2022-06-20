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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type DeploymentSpec struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Replicas int32  `json:"replicas"`
}

type ServiceSpec struct {
	Name string `json:"name"`
}

type IngressSpec struct {
	Name string `json:"name"`
}

// EkiMonitorSpec defines the desired state of EkiMonitor
type EkiMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Deployment DeploymentSpec `json:"deployment"`
	Service    ServiceSpec    `json:"service"`
	Ingress    IngressSpec    `json:"ingress"`
}

// EkiMonitorStatus defines the observed state of EkiMonitor
type EkiMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EkiMonitor is the Schema for the ekimonitors API
type EkiMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EkiMonitorSpec   `json:"spec,omitempty"`
	Status EkiMonitorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EkiMonitorList contains a list of EkiMonitor
type EkiMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EkiMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EkiMonitor{}, &EkiMonitorList{})
}
