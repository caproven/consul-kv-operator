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
	Source ConsulKVSecretSource `json:"source"`
	Secret ConsulKVSecretSecret `json:"secret"`
}

type ConsulKVSecretSource struct {
	// Address of the Consul server
	Address string `json:"address"`

	// Token used for authentication with the Consul server
	// TODO change to secretRef
	// +optional
	Token string `json:"token,omitempty"`
}

type ConsulKVSecretSecret struct {
	// Name of the secret that will be created. Is immutable. Defaults
	// to the ConsulKVSecret name
	// +optional
	Name string `json:"name,omitempty"`

	// Data values that will populate the secret
	Data []ConsulKVSecretDataEntry `json:"data,omitempty"`
}

// ConsulKVSecretDataEntry defines an entry to be populated in the
// secret based on a value in Consul
type ConsulKVSecretDataEntry struct {
	SourceKey string `json:"sourcekey"`
	Key       string `json:"key"`
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
