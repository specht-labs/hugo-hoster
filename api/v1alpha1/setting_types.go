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

// SettingSpec defines the desired state of Setting
type SettingSpec struct {
	// TLS configure if you want to enable TLS
	// +kubebuilder:default:={enable: false}
	TLS TLSSpec `json:"tls,omitempty"`

	// IngressClassName makes it possible to override the ingress-class
	// +kubebuilder:default:=nginx
	IngressClassName string `json:"ingressClassName,omitempty"`

	// S3Config contains the configuration of the S3 bucket to upload the pages to
	// +kubebuilder:validation:Required
	S3Config S3Config `json:"s3_config"`

	// ProxyURL is the URL from which the static files are served. Most of the time this is the same as `spec.S3Config.Endpoint``. If this is empty, `spec.S3Config.Endpoint` is used
	// +kubebuilder:example:=https://f003.backblazeb2.com/file
	ProxyURL string `json:"serving_url,omitempty"`

	// NginxProxyReplica is the number of replicas for each page
	// +kubebuilder:default:=1
	// +kubebuilder:validation:Optional
	NginxProxyReplica int32 `json:"nginxProxyReplica,omitempty"`
}

type S3Config struct {
	// S3Endpoint is the Endpoint URL for the S3 Bucket to upload pages to
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=https://s3.eu-central-003.backblazeb2.com
	Endpoint string `json:"endpoint"`

	// BucketName is the S3 Bucket to which all sites are uploaded
	// +kubebuilder:validation:Required
	BucketName string `json:"bucketname"`

	// SecretName is the name of the Kubernetes Secret that contains the S3 AccessKeyId and the AccessKey
	// +kubebuilder:validation:Required
	SecretName string `json:"secretName"`

	// AccessKeyIDRef is the name of the key in the S3SecretName that contains the AccessKeyId
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=AccessKeyId
	AccessKeyIDRef string `json:"accessKeyIdKeyName"`

	// AccessKeyRef is the name of the key in the S3SecretName that contains the AccessKey
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=AccessKey
	AccessKeyRef string `json:"accessKeyKeyName"`
}

// TLSSpec holds the TLS configuration used
type TLSSpec struct {
	// +kubebuilder:default:=false
	Enable      bool              `json:"enable,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type SettingStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="TlsEnabled",type=string,JSONPath=`.spec.tls.enable`
// +kubebuilder:printcolumn:name="S3Endpoint",type=string,JSONPath=`.spec.s3_config.endpoint`
// Setting is the Schema for the settings API
type Setting struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SettingSpec   `json:"spec,omitempty"`
	Status SettingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SettingList contains a list of Setting
type SettingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Setting `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Setting{}, &SettingList{})
}
