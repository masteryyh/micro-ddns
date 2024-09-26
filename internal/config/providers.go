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

import "fmt"

// AliCloudSpec is the information of AliCloud API credential and extra settings
type AliCloudSpec struct {
	// AccessKeyID is the AccessKey of the account
	AccessKeyID string `json:"accessKeyId" yaml:"accessKeyId"`

	// AccessKeySecret is the AccessKeySecret of the account
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"`

	// Line is the resolve line of the record
	Line *string `json:"line,omitempty" yaml:"line,omitempty"`
}

func (spec *AliCloudSpec) Validate() error {
	if spec.AccessKeyID == "" {
		return fmt.Errorf("AccessKeyID cannot be empty")
	}

	if spec.AccessKeySecret == "" {
		return fmt.Errorf("AccessKeySecret cannot be empty")
	}

	return nil
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

func (spec *CloudflareSpec) Validate() error {
	if spec.APIToken == nil || *spec.APIToken == "" {
		if spec.GlobalAPIKey == nil || *spec.GlobalAPIKey == "" {
			return fmt.Errorf("must choose between api token or global api key with email")
		}

		if spec.Email == nil || *spec.Email == "" {
			return fmt.Errorf("must choose between api token or global api key with email")
		}
	}

	return nil
}

// DNSPodSpec is the information of Tencent DNSPod API credential and extra settings
type DNSPodSpec struct {
	// SecretID is the SecretID in your credential
	SecretID string `json:"secretId" yaml:"secretId"`

	// SecretKey is the SecretKey in your credential
	SecretKey string `json:"secretKey" yaml:"secretKey"`

	// LineID is the ID of line, leave empty for default line (0)
	LineID *string `json:"lineId,omitempty" yaml:"lineId,omitempty"`
}

func (spec *DNSPodSpec) Validate() error {
	if spec.SecretID == "" {
		return fmt.Errorf("SecretID cannot be empty")
	}

	if spec.SecretKey == "" {
		return fmt.Errorf("SecretKey cannot be empty")
	}

	return nil
}

// HuaweiCloudSpec is the information of Huawei Cloud credential and settings
type HuaweiCloudSpec struct {
	// AccessKey is the access key (AK) of the account
	AccessKey string `json:"accessKey" yaml:"accessKey"`

	// SecretAccessKey is the secret key (SK) of the account
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey"`

	// Region is the region of resources, this decides the endpoint of the API
	Region string `json:"region" yaml:"region"`
}

func (spec *HuaweiCloudSpec) Validate() error {
	if spec.AccessKey == "" {
		return fmt.Errorf("AcessKey cannot be empty")
	}

	if spec.SecretAccessKey == "" {
		return fmt.Errorf("SecretAccessKey cannot be empty")
	}

	if spec.Region == "" {
		return fmt.Errorf("region cannot be empty")
	}

	return nil
}

// JDCloudSpec is the information of JDCloud credential and settings
type JDCloudSpec struct {
	// AccessKey is the access key of the account
	AccessKey string `json:"accessKey" yaml:"accessKey"`

	// SecretKey is the secret key of the account
	SecretKey string `json:"secretKey" yaml:"secretKey"`

	// ViewID is the resolve line ID of the DNS record, leave it empty for default value -1
	ViewID *int `json:"viewId,omitempty" yaml:"viewId,omitempty"`
}

func (spec *JDCloudSpec) Validate() error {
	if spec.AccessKey == "" {
		return fmt.Errorf("AccessKey cannot be empty")
	}

	if spec.SecretKey == "" {
		return fmt.Errorf("SecretKey cannot be empty")
	}

	return nil
}

// RFC2136Spec is the information about an RFC 2136 compliant DNS server
type RFC2136Spec struct {
	// Address of the DNS server
	Address string `json:"address" yaml:"address"`

	// Port of the DNS server, leave empty for default value (53)
	Port *int `json:"port,omitempty" yaml:"port,omitempty"`

	// UseTCP specifies if TCP should be used instead of UDP to contact DNS server
	UseTCP *bool `json:"useTcp,omitempty" yaml:"useTcp,omitempty"`

	TSIG *TSIGSpec `json:"tsig,omitempty" yaml:"tsig,omitempty"`

	GSSTSIG *GSSTSIGSpec `json:"gssTsig,omitempty" yaml:"gssTsig,omitempty"`
}

func (spec *RFC2136Spec) Validate() error {
	if spec.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if spec.Port != nil && (*spec.Port < 1 || *spec.Port > 65535) {
		return fmt.Errorf("port %d is invalid", *spec.Port)
	}

	if spec.TSIG != nil {
		return spec.TSIG.Validate()
	} else if spec.GSSTSIG != nil {
		return spec.GSSTSIG.Validate()
	}

	return nil
}

// TSIGSpec is the information about TSIG authentication
type TSIGSpec struct {
	// KeyName is the name of TSIG key
	KeyName string `json:"keyName,omitempty" yaml:"keyName,omitempty"`

	// Key is the key for TSIG
	Key string `json:"key,omitempty" yaml:"key,omitempty"`
}

func (spec *TSIGSpec) Validate() error {
	if spec.KeyName == "" {
		return fmt.Errorf("key name cannot be empty")
	}

	if spec.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	return nil
}

type GSSTSIGSpec struct {
	// Domain is the domain of the directory service
	Domain string `json:"domain" yaml:"domain"`

	// Username is the name of the user to authenticate
	Username string `json:"username" yaml:"username"`

	// Password is the password of the user
	Password string `json:"password" yaml:"password"`
}

func (spec *GSSTSIGSpec) Validate() error {
	if spec.Domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	if spec.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if spec.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}
