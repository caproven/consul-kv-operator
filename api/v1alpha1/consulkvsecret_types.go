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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConsulKVSecretSpec defines the desired state of ConsulKVSecret
type ConsulKVSecretSpec struct {
	Source ConsulKVSecretSource  `json:"source"`
	Values []ConsulKVSecretValue `json:"values"`
	Output ConsulKVSecretOutput  `json:"secret"`
}

// ConsulKVSecretSource describes the Consul server acting as the source
// of values
type ConsulKVSecretSource struct {
	// Host of the Consul server
	Host string `json:"host"`

	// Port of the Consul server
	Port int `json:"port"`

	// Token used for authentication with the Consul server
	// TODO change to secretRef
	// +optional
	Token string `json:"token,omitempty"`
}

// ConsulKVSecretValue defines an entry to be populated in the
// secret based on a value in Consul
type ConsulKVSecretValue struct {
	// SourceKey is the key in the Consul KV store whose value will be pulled
	SourceKey string `json:"sourcekey"`

	// Key is the mapped key in a secret containing the value from Consul
	Key string `json:"key"`
}

// ConsulKVSecretOutput describes the secret generated which holds values
// read from the Consul server
type ConsulKVSecretOutput struct {
	// Name of the secret that will be created. Is immutable. Defaults
	// to the ConsulKVSecret name
	// +optional
	Name string `json:"name,omitempty"`
}

// ConsulKVSecretStatus defines the observed state of ConsulKVSecret
type ConsulKVSecretStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ConsulKVSecret is the Schema for the consulkvsecrets API
type ConsulKVSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConsulKVSecretSpec   `json:"spec,omitempty"`
	Status ConsulKVSecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConsulKVSecretList contains a list of ConsulKVSecret
type ConsulKVSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConsulKVSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConsulKVSecret{}, &ConsulKVSecretList{})
}
