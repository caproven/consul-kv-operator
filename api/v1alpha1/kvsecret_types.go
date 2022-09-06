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

// KVSecretSpec defines the desired state of KVSecret
type KVSecretSpec struct {
	Source SourceSpec   `json:"source"`
	Values []KeyMapping `json:"values"`
	Output OutputSpec   `json:"secret"`
}

// SourceSpec describes the Consul server acting as the source
// of values
type SourceSpec struct {
	// Host of the Consul server
	Host string `json:"host"`

	// Port of the Consul server
	Port int `json:"port"`

	// Token used for authentication with the Consul server
	// TODO change to secretRef
	// +optional
	Token string `json:"token,omitempty"`
}

// KeyMapping defines an entry to be populated in the
// secret based on a value in Consul
type KeyMapping struct {
	// SourceKey is the key in the Consul KV store whose value will be pulled
	SourceKey string `json:"sourcekey"`

	// Key is the mapped key in a secret containing the value from Consul
	Key string `json:"key"`
}

// OutputSpec describes the secret generated which holds values
// read from the Consul server
type OutputSpec struct {
	// Name of the secret that will be created. Is immutable. Defaults
	// to the KVSecret name
	// +optional
	Name string `json:"name,omitempty"`
}

// KVSecretStatus defines the observed state of KVSecret
type KVSecretStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KVSecret is the Schema for the kvsecrets API
type KVSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KVSecretSpec   `json:"spec,omitempty"`
	Status KVSecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KVSecretList contains a list of KVSecret
type KVSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KVSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KVSecret{}, &KVSecretList{})
}
