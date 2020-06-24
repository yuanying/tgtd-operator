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

// InitiatorGroupSpec defines the desired state of InitiatorGroup
type InitiatorGroupSpec struct {
	// Addresses are used for initiator address
	Addresses []string `json:"addresses,omitempty"`

	// NodeSelector is a selector to select initiator nodes
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// InitiatorNameStrategy is a strategy how to decide initiator name from node object
	InitiatorNameStrategy InitiatorGroupInitiatorNameStrategy `json:"initiatorNameStrategy,omitempty"`
}

type InitiatorGroupInitiatorNameStrategy struct {
	// Type is a strategy type. Can be "NodeName" or "AnnotationKey". Default is NodeName
	// +kubebuilder:default:NodeName
	// +kubebuilder:validation:Enum=NodeName;Annotation
	Type InitiatorGroupInitiatorNameStrategyType `json:"type,omitempty"`

	// InitiatorNamePrefix is used for generating initiator name from node name
	InitiatorNamePrefix *string `json:"initiatorNamePrefix,omitempty"`

	// AnnotationKey is used for retrieving initiator name from annotation
	AnnotationKey *string `json:"annotationKey,omitempty"`
}

type InitiatorGroupInitiatorNameStrategyType string

const (
	NodeNameInitiatorNameStrategy = "NodeName"

	AnnotationInitiatorNameStrategy = "Annotation"
)

// InitiatorGroupStatus defines the observed state of InitiatorGroup
type InitiatorGroupStatus struct {
	// Addresses are addresses of initiators
	Addresses []string `json:"addresses,omitempty"`

	// Initiators are names of initiators
	Initiators []string `json:"initiators,omitempty"`
}

type InitiatorGroupReference struct {
	// Name is unique to reference a initiator group resource.
	Name string `json:"name,omitempty"`
}

// +kubebuilder:object:root=true

// InitiatorGroup is the Schema for the initiatorgroups API
// +kubebuilder:resource:shortName=ig,scope=Cluster
// +kubebuilder:subresource:status
type InitiatorGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InitiatorGroupSpec   `json:"spec,omitempty"`
	Status InitiatorGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InitiatorGroupList contains a list of InitiatorGroup
type InitiatorGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InitiatorGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InitiatorGroup{}, &InitiatorGroupList{})
}
