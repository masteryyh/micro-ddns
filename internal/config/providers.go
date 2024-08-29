/*
Copyright © 2024 masteryyh <yyh991013@163.com>

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

package config

// AliCloudSpec is the information of AliCloud API credential and extra settings
type AliCloudSpec struct {
	// AccessKeyID is the AccessKey of the account
	AccessKeyID string `json:"accessKeyId" yaml:"accessKeyId"`

	// AccessKeySecret is the AccessKeySecret of the account
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`

	// RegionID is the region of the domain
	RegionID string `json:"regionId" yaml:"regionId"`

	// Line is the resolve line of the record
	Line *string `json:"line,omitempty" yaml:"line,omitempty"`
}

// CloudflareSpec is the information of Cloudflare API credential
type CloudflareSpec struct {
	// APIToken is the fine-grained token generated by user
	APIToken *string `json:"apiToken,omitempty" yaml:"apiToken,omitempty"`

	// GlobalAPIKey is the global API key with full permission to access Cloudflare API
	GlobalAPIKey *string `json:"globalApiKey,omitempty" yaml:"globalApiKey,omitempty"`

	// Email is the email of user account, required when using Global API key
	Email *string `json:"email,omitempty" yaml:"email,omitempty"`
}

// DNSPodSpec is the information of Tencent DNSPod API credential and extra settings
type DNSPodSpec struct {
	// SecretID is the SecretID in your credential
	SecretID string `json:"secretId" yaml:"secretKey"`

	// SecretKey is the SecretKey in your credential
	SecretKey string `json:"secretKey" yaml:"secretKey"`

	// Region is the region of your resource
	Region string `json:"region" yaml:"region"`

	// LineID is the ID of line, leave empty for default line (0)
	LineID *string `json:"lineId,omitempty" yaml:"lineId,omitempty"`
}