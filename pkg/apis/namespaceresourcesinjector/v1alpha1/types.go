package v1alpha1

import (
	"regexp"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceResourcesInjector describes a NamespaceResourcesInjector resource
type NamespaceResourcesInjector struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec   NamespaceResourcesInjectorSpec   `json:"spec"`
	Status NamespaceResourcesInjectorStatus `json:"status"`
}

func (nri *NamespaceResourcesInjector) CanInject(namespace string) bool {
	for _, item := range nri.Spec.Namespaces {
		if match, _ := regexp.MatchString(item, namespace); match {
			return true
		}
	}
	return false
}

func (nri *NamespaceResourcesInjector) Injected(namespace string) bool {
	for _, item := range nri.Status.InjectedNamespaces {
		if item == namespace {
			return true
		}
	}
	return false
}

// NamespaceResourcesInjectorSpec is the spec for a NamespaceResourcesInjector resource
type NamespaceResourcesInjectorSpec struct {
	Namespaces []string `json:"namespaces"`
	Resources  []string `json:"resources"`
}

// NamespaceResourcesInjectorStatus is the status for a NamespaceResourcesInjector resource
type NamespaceResourcesInjectorStatus struct {
	InjectedNamespaces []string `json:"injectedNamespaces"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceResourcesInjectorList is a list of NamespaceResourcesInjector resources
type NamespaceResourcesInjectorList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []NamespaceResourcesInjector `json:"items"`
}
