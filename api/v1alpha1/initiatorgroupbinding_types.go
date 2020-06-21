/*
Copyright 2020 O.Yuanying <yuanying@fraction.jp>

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

// InitiatorGroupBindingSpec defines the desired state of InitiatorGroupBinding
type InitiatorGroupBindingSpec struct {
	// TargetRef is a reference to target resource
	// +kubebuilder:validation:Required
	TargetRef TargetReference `json:"targetRef,omitempty"`

	// InitiatorGroupRef is a reference to initiator group resource
	// +kubebuilder:validation:Required
	InitiatorGroupRef InitiatorGroupReference `json:"initiatorGroupRef,omitempty"`
}

// InitiatorGroupBindingStatus defines the observed state of InitiatorGroupBinding
type InitiatorGroupBindingStatus struct {
}

// +kubebuilder:object:root=true

// InitiatorGroupBinding is the Schema for the initiatorgroupbindings API
// +kubebuilder:resource:shortName=igb,scope=Cluster
// +kubebuilder:subresource:status
type InitiatorGroupBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InitiatorGroupBindingSpec   `json:"spec,omitempty"`
	Status InitiatorGroupBindingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InitiatorGroupBindingList contains a list of InitiatorGroupBinding
type InitiatorGroupBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InitiatorGroupBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InitiatorGroupBinding{}, &InitiatorGroupBindingList{})
}
