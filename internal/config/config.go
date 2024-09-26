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

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/masteryyh/micro-ddns/pkg/utils"
	"gopkg.in/yaml.v3"
)

var (
	config Config

	domainRegex    = regexp.MustCompile(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`)
	subdomainRegex = regexp.MustCompile(`^([a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)|([a-zA-Z0-9]*@[a-zA-Z0-9]*)$`)
)

type NetworkStack string

const (
	IPv4 NetworkStack = "IPv4"
	IPv6 NetworkStack = "IPv6"
)

type AddressDetectionType string

const (
	AddressDetectionIface      AddressDetectionType = "Interface"
	AddressDetectionThirdParty AddressDetectionType = "ThirdParty"
)

type LocalAddressPolicy string

const (
	LocalAddressPolicyIgnore LocalAddressPolicy = "Ignore"
	LocalAddressPolicyAllow  LocalAddressPolicy = "Allow"
	LocalAddressPolicyPrefer LocalAddressPolicy = "Prefer"
)

type DNSProvider string

const (
	DNSProviderCloudflare  DNSProvider = "Cloudflare"
	DNSProviderAliCloud    DNSProvider = "AliCloud"
	DNSProviderDNSPod      DNSProvider = "DNSPod"
	DNSProviderHuaweiCloud DNSProvider = "HuaweiCloud"
	DNSProviderJDCloud     DNSProvider = "JDCloud"
	DNSProviderRFC2136     DNSProvider = "RFC2136"
)

// DNSProviderSpec is the specification of DNS provider, currently only Cloudflare
// is supported
type DNSProviderSpec struct {
	// Name is the name of DNS provider
	Name DNSProvider `json:"name" yaml:"name"`

	Cloudflare *CloudflareSpec `json:"cloudflare,omitempty" yaml:"cloudflare,omitempty"`

	AliCloud *AliCloudSpec `json:"alicloud,omitempty" yaml:"alicloud,omitempty"`

	DNSPod *DNSPodSpec `json:"dnspod,omitempty" yaml:"dnspod,omitempty"`

	Huawei *HuaweiCloudSpec `json:"huawei,omitempty" yaml:"huawei,omitempty"`

	JD *JDCloudSpec `json:"jd,omitempty" yaml:"jd,omitempty"`

	RFC2136 *RFC2136Spec `json:"rfc2136,omitempty" yaml:"rfc2136,omitempty"`
}

func (spec *DNSProviderSpec) Validate() error {
	if spec.Name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}

	count := 0
	if spec.Cloudflare != nil {
		count++
	}
	if spec.AliCloud != nil {
		count++
	}
	if spec.DNSPod != nil {
		count++
	}
	if spec.Huawei != nil {
		count++
	}
	if spec.JD != nil {
		count++
	}
	if spec.RFC2136 != nil {
		count++
	}

	if count == 0 {
		return fmt.Errorf("no provider specified")
	}

	if count > 1 {
		return fmt.Errorf("only 1 provider can be used within 1 spec")
	}

	if spec.Cloudflare != nil {
		return spec.Validate()
	} else if spec.AliCloud != nil {
		return spec.AliCloud.Validate()
	} else if spec.DNSPod != nil {
		return spec.DNSPod.Validate()
	} else if spec.Huawei != nil {
		return spec.Huawei.Validate()
	} else if spec.JD != nil {
		return spec.JD.Validate()
	} else if spec.RFC2136 != nil {
		return spec.RFC2136.Validate()
	}

	return nil
}

// NetworkInterfaceDetectionSpec defines how should we get IP address from an interface
// By default the first address detected will be used
type NetworkInterfaceDetectionSpec struct {
	// Name is the name of interface
	Name string `json:"name" yaml:"name"`
}

func (spec *NetworkInterfaceDetectionSpec) Validate() error {
	if spec.Name == "" {
		return fmt.Errorf("interface name cannot be empty")
	}
	return nil
}

// ThirdPartyServiceSpec defines how should we access third party API to get our IP address
type ThirdPartyServiceSpec struct {
	// URL is the URL of third-party API
	URL string `json:"url" yaml:"url"`

	// JsonPath is the path to the address if data returned by API is JSON-formatted
	JsonPath *string `json:"jsonPath,omitempty" yaml:"jsonPath,omitempty"`

	// Params will be added to the URL
	Params *map[string]string `json:"params,omitempty" yaml:"params,omitempty"`

	// Headers will be added to the request header if not empty
	Headers *map[string]string `json:"customHeaders,omitempty" yaml:"customHeaders,omitempty"`

	// Username is the username for HTTP basic authentication if required
	Username *string `json:"username,omitempty" yaml:"username,omitempty"`

	// Password is the password for HTTP basic authentication if required
	Password *string `json:"password,omitempty" yaml:"password,omitempty"`
}

func (spec *ThirdPartyServiceSpec) Validate() error {
	if spec.URL == "" {
		return fmt.Errorf("url cannot be empty")
	}

	if spec.JsonPath != nil && *spec.JsonPath == "" {
		spec.JsonPath = nil
	}

	return nil
}

// AddressDetectionSpec defines how should we detect current IP address
type AddressDetectionSpec struct {
	// Type is the type of address detection method
	// Currently we can acquire address by network interface or 3rd-party API
	Type AddressDetectionType `json:"type" yaml:"type"`

	// LocalAddressPolicy defines how should we process addresses
	// LocalAddressPolicyIgnore means the operation would fail when no public address presents on the interface
	// LocalAddressPolicyAllow means local addresses will be used for DNS record, but only if no public address presents on the interface
	// LocalAddressPolicyPrefer means local addresses will be used for DNS record even public address presents on the interface
	LocalAddressPolicy *LocalAddressPolicy `json:"localAddressPolicy,omitempty" yaml:"localAddressPolicy,omitempty"`

	Interface *NetworkInterfaceDetectionSpec `json:"interface,omitempty" yaml:"interface,omitempty"`

	API *ThirdPartyServiceSpec `json:"api,omitempty" yaml:"api,omitempty"`
}

func (spec *AddressDetectionSpec) Validate() error {
	t := string(spec.Type)
	if t == "" {
		return fmt.Errorf("detection type is needed")
	}

	if t != "Interface" && t != "ThirdParty" {
		return fmt.Errorf("unknown detection type %s", t)
	}

	if spec.LocalAddressPolicy == nil {
		spec.LocalAddressPolicy = (*LocalAddressPolicy)(utils.StringPtr("Ignore"))
	}

	p := *spec.LocalAddressPolicy
	if p != "Ignore" && p != "Prefer" && p != "Allow" {
		return fmt.Errorf("unknown localAddressPolicy %s", p)
	}

	if spec.Interface == nil && spec.API == nil {
		return fmt.Errorf("must define a address detection method")
	}

	if spec.Interface != nil {
		return spec.Interface.Validate()
	} else {
		return spec.API.Validate()
	}
}

// DDNSSpec is the specification of DDNS service
type DDNSSpec struct {
	// Name is the name of the specification
	Name string `json:"name" yaml:"name"`

	// Domain is the domain of user
	Domain string `json:"domain" yaml:"domain"`

	// Subdomain is the subdomain to update, use "@" if no subdomain is used
	Subdomain string `json:"subdomain" yaml:"subdomain"`

	// Stack determines if IPv4 or IPv6 is used
	Stack NetworkStack `json:"stack" yaml:"stack"`

	// Cron is the cron expression about how should we schedule this task
	Cron string `json:"cron" yaml:"cron"`

	Detection AddressDetectionSpec `json:"detection" yaml:"detection"`

	Provider DNSProviderSpec `json:"provider" yaml:"provider"`
}

func (spec *DDNSSpec) Validate() error {
	if spec.Name == "" {
		return fmt.Errorf("name is needed for a DDNS spec")
	}

	if spec.Domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	if !domainRegex.MatchString(spec.Domain) {
		return fmt.Errorf("%s is not a valid domain", spec.Domain)
	}

	if spec.Subdomain == "" {
		return fmt.Errorf("subdomain cannot be empty, use \"@\" if you want to use zone apex")
	}

	if !subdomainRegex.MatchString(spec.Subdomain) {
		return fmt.Errorf("%s is not a valid subdomain", spec.Subdomain)
	}

	stack := string(spec.Stack)
	if stack == "" {
		return fmt.Errorf("stack cannot be empty, must be one of IPv4 or IPv6")
	}

	if stack != "IPv4" && stack != "IPv6" {
		return fmt.Errorf("%s is not a valid stack, must be one of IPv4 or IPv6", stack)
	}

	if spec.Cron == "" {
		return fmt.Errorf("crontab cannot be empty")
	}

	return spec.Detection.Validate()
}

// Config is the configuration of this application
type Config struct {
	DDNS []*DDNSSpec `json:"ddns,omitempty" yaml:"ddns,omitempty"`
}

func (c *Config) Validate() error {
	for _, v := range c.DDNS {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func ReadConfigOrGet(path string) (*Config, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("config path points to a directory")
	}

	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("config path points to an unknown file type")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fileType := parts[len(parts)-1]
	switch fileType {
	case "yaml", "yml":
		if err := yaml.Unmarshal(content, &config); err != nil {
			return nil, err
		}
	case "json":
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("config path points to an unknown file type")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
