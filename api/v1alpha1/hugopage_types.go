/*
Copyright 2023.

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

// HugoPageSpec defines the desired state of HugoPage
type HugoPageSpec struct {
	// Repository specifies the target Repository to pull from for building the hugo site
	// +kubebuilder:validation:Required
	// +kubebuilder:example:https://github.com/cedi/cedi.github.io.git
	Repository string `json:"repository"`

	// Branch specifies the branch from which to build the site. (default: main)
	// +kubebuilder:default:main
	Branch string `json:"branch,omitempty"`

	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// BuildType configures how the Hugo-Site is rebuild.
	//
	// Poll takes the configured polling interval to rebuild the page
	// Webhook requires a CI/CD Pipeline to call the Webhook URL of this page to re-build the site
	// +kubebuilder:validation:Enum=cron
	// +kubebuilder:default:cron
	BuildType string `json:"type,omitempty"`

	// CronInterval is the polling interval in which the hugo-site is refreshed as a cron syntax string
	// +kubebuilder:default:"*/5 * * * *"
	CronInterval string `json:"interval,omitempty"`
}

// HugoPageStatus defines the observed state of HugoPage
type HugoPageStatus struct {
	// LastBuild is a date-time when the Hugo Page was last built
	// +kubebuilder:validation:Format:date-time
	LastBuild string `json:"lastbuild"`

	// Commit contains the commit-id of the current build
	Commit string `json:"commit"`

	// Status contains the status of the last build action
	// +kubebuilder:validation:Enum=Failed;Success;Cancelled
	Status string `json:"status"`
}

// HugoPage is the Schema for the HugoPages API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="URL",type=string,JSONPath=`.spec.url`
// +kubebuilder:printcolumn:name="LastBuild",type=string,JSONPath=`.status.LastBuild`
// +kubebuilder:printcolumn:name="Commit",type=string,JSONPath=`.status.Commit`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.Status`
// +k8s:openapi-gen=true
type HugoPage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HugoPageSpec   `json:"spec,omitempty"`
	Status HugoPageStatus `json:"status,omitempty"`
}

// HugoPageList contains a list of HugoPage
// +kubebuilder:object:root=true
type HugoPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HugoPage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HugoPage{}, &HugoPageList{})
}
