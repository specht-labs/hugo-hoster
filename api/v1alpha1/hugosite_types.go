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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HugoSitePreviewSpec controls if and how previews are built for this page
type HugoSitePreviewSpec struct {
	// Enabled controls if previews are build
	// +kubebuilder:default:false
	Enabled bool `json:"enabled"`

	// Branch (regex) controls from which branches the preview is built
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:[*]
	Branch []string `json:"branches"`

	// ExcludeBranch controls which branch to exclude from preview builds, such as `main` or `site`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:[main, site]
	ExcludeBranch []string `json:"exclude"`
}

// HugoSiteSpec defines the desired state of HugoSite
type HugoSiteSpec struct {
	// Repository specifies the target Repository to pull from for building the hugo site
	// +kubebuilder:validation:Required
	// +kubebuilder:example:https://github.com/cedi/cedi.github.io.git
	Repository string `json:"repository"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	URL string `json:"url"`

	// MainBranch specifies the branch from which to build the site. (default: main)
	// +kubebuilder:default:main
	MainBranch string `json:"mainbranch,omitempty"`

	// Preview is the Preview object to enable or disable preview builds
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:nil
	Preview *HugoSitePreviewSpec `json:"preview,omitempty"`

	// BuildType configures how the Hugo-Site is rebuild.
	//
	// Poll takes the configured polling interval to rebuild the page
	// Webhook requires a CI/CD Pipeline to call the Webhook URL of this page to re-build the site
	// +kubebuilder:validation:Enum=Poll;Webhook
	// +kubebuilder:default:Poll
	BuildType string `json:"type"`

	// PollInterval is the polling interval in which the hugo-site is refreshed as a go time interval string
	// +kubebuilder:validation:Pattern:^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$
	// +kubebuilder:default:5m
	PollInterval time.Duration `json:"interval"`
}

// HugoSiteStatus defines the observed state of HugoSite
type HugoSiteStatus struct {
	// LastBuild is a date-time when the Hugo Page was last built
	// +kubebuilder:validation:Format:date-time
	LastBuild string `json:"lastbuild"`

	// Commit contains the commit-id of the current build
	Commit string `json:"commit"`

	// Status contains the status of the last build action
	// +kubebuilder:validation:Enum=Failed;Success;Cancelled
	Status string `json:"status"`
}

// HugoSite is the Schema for the hugosites API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="URL",type=string,JSONPath=`.spec.url`
// +kubebuilder:printcolumn:name="LastBuild",type=string,JSONPath=`.status.LastBuild`
// +kubebuilder:printcolumn:name="Commit",type=string,JSONPath=`.status.Commit`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.Status`
// +k8s:openapi-gen=true
type HugoSite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HugoSiteSpec   `json:"spec,omitempty"`
	Status HugoSiteStatus `json:"status,omitempty"`
}

// HugoSiteList contains a list of HugoSite
// +kubebuilder:object:root=true
type HugoSiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HugoSite `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HugoSite{}, &HugoSiteList{})
}
