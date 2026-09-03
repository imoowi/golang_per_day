/*
Copyright 2026.

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
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AIServiceSpec defines the desired state of AIService
type AIServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// foo is an example field of AIService. Edit aiservice_types.go to remove/update
	// +optional
	Foo       *string        `json:"foo,omitempty"`
	Model     string         `json:"model,omitempty"`
	Replicas  *int32         `json:"replicas,omitempty"`
	Image     string         `json:"image,omitempty"`
	Resources AIResourceSpec `json:"resources,omitempty"`
}

type AIResourceSpec struct {
	GPU    int32  `json:"gpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// AIServiceStatus defines the observed state of AIService.
type AIServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the AIService resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
	ReadyReplicas int32              `json:"readyReplicas,omitempty"`
	Phase         string             `json:"phase,omitempty"`
	Message       string             `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AIService is the Schema for the aiservices API
type AIService struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of AIService
	// +required
	Spec AIServiceSpec `json:"spec"`

	// status defines the observed state of AIService
	// +optional
	Status AIServiceStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// AIServiceList contains a list of AIService
type AIServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []AIService `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &AIService{}, &AIServiceList{})
		return nil
	})
}
