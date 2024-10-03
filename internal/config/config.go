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
	"sync"

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
	// Name of the provider specification
	Name string `json:"name" yaml:"name"`

	providerType DNSProvider

	Cloudflare *CloudflareSpec `json:"cloudflare,omitempty" yaml:"cloudflare,omitempty"`

	AliCloud *AliCloudSpec `json:"alicloud,omitempty" yaml:"alicloud,omitempty"`

	DNSPod *DNSPodSpec `json:"dnspod,omitempty" yaml:"dnspod,omitempty"`

	Huawei *HuaweiCloudSpec `json:"huawei,omitempty" yaml:"huawei,omitempty"`

	JD *JDCloudSpec `json:"jd,omitempty" yaml:"jd,omitempty"`

	RFC2136 *RFC2136Spec `json:"rfc2136,omitempty" yaml:"rfc2136,omitempty"`
}

func (spec *DNSProviderSpec) Validate() error {
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
		spec.providerType = DNSProviderCloudflare
		return spec.Cloudflare.Validate()
	} else if spec.AliCloud != nil {
		spec.providerType = DNSProviderAliCloud
		return spec.AliCloud.Validate()
	} else if spec.DNSPod != nil {
		spec.providerType = DNSProviderDNSPod
		return spec.DNSPod.Validate()
	} else if spec.Huawei != nil {
		spec.providerType = DNSProviderHuaweiCloud
		return spec.Huawei.Validate()
	} else if spec.JD != nil {
		spec.providerType = DNSProviderJDCloud
		return spec.JD.Validate()
	} else if spec.RFC2136 != nil {
		spec.providerType = DNSProviderRFC2136
		return spec.RFC2136.Validate()
	}

	return nil
}

func (spec *DNSProviderSpec) GetType() DNSProvider {
	return spec.providerType
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
	// Name of this address detection specification
	Name string `json:"name" yaml:"name"`

	detectionType AddressDetectionType

	// LocalAddressPolicy defines how should we process addresses
	// LocalAddressPolicyIgnore means the operation would fail when no public address presents on the interface
	// LocalAddressPolicyAllow means local addresses will be used for DNS record, but only if no public address presents on the interface
	// LocalAddressPolicyPrefer means local addresses will be used for DNS record even public address presents on the interface
	LocalAddressPolicy *LocalAddressPolicy `json:"localAddressPolicy,omitempty" yaml:"localAddressPolicy,omitempty"`

	Interface *NetworkInterfaceDetectionSpec `json:"interface,omitempty" yaml:"interface,omitempty"`

	API *ThirdPartyServiceSpec `json:"api,omitempty" yaml:"api,omitempty"`
}

func (spec *AddressDetectionSpec) Validate() error {
	if spec.LocalAddressPolicy == nil {
		spec.LocalAddressPolicy = (*LocalAddressPolicy)(utils.StringPtr("Ignore"))
	}

	p := *spec.LocalAddressPolicy
	if p != "Ignore" && p != "Prefer" && p != "Allow" {
		return fmt.Errorf("unknown localAddressPolicy %s", p)
	}

	if spec.Interface != nil {
		spec.detectionType = AddressDetectionIface
		return spec.Interface.Validate()
	} else if spec.API != nil {
		spec.detectionType = AddressDetectionThirdParty
		return spec.API.Validate()
	}
	return fmt.Errorf("must specify a detection method")
}

func (spec *AddressDetectionSpec) GetDetectionType() AddressDetectionType {
	return spec.detectionType
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

	// ProviderRef is the name of the DNS provider specification defined by user
	ProviderRef string `json:"providerRef" yaml:"providerRef"`

	// DetectionRef is the name of the address detection specification defined by user
	DetectionRef string `json:"detectionRef" yaml:"detectionRef"`

	detectionSpec *AddressDetectionSpec

	providerSpec *DNSProviderSpec
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

	if spec.ProviderRef == "" {
		return fmt.Errorf("providerref cannot be empty")
	}

	if spec.DetectionRef == "" {
		return fmt.Errorf("detectionref cannot be empty")
	}

	return nil
}

func (spec *DDNSSpec) GetDetectionSpec() *AddressDetectionSpec {
	return spec.detectionSpec
}

func (spec *DDNSSpec) GetProviderSpec() *DNSProviderSpec {
	return spec.providerSpec
}

// Config is the configuration of this application
type Config struct {
	DDNS []*DDNSSpec `json:"ddns" yaml:"ddns"`

	Detection []*AddressDetectionSpec `json:"detection" yaml:"detection"`

	Provider []*DNSProviderSpec `json:"provider" yaml:"provider"`
}

func (c *Config) Validate() error {
	if len(c.DDNS) == 0 {
		return fmt.Errorf("must have at least 1 ddns spec")
	}

	if len(c.Detection) == 0 {
		return fmt.Errorf("must have at least 1 detection spec")
	}

	if len(c.Provider) == 0 {
		return fmt.Errorf("must have at least 1 provider")
	}

	var validateWg sync.WaitGroup
	validateWg.Add(3)

	var ddnsErr error
	ddns := make(map[string]*DDNSSpec)
	go func(wg *sync.WaitGroup) {
		for _, spec := range c.DDNS {
			if _, exists := ddns[spec.Name]; exists {
				ddnsErr = fmt.Errorf("ddns spec %s already exists", spec.Name)
				break
			}
			if ddnsErr = spec.Validate(); ddnsErr != nil {
				break
			}
			ddns[spec.Name] = spec
		}
		wg.Done()
	}(&validateWg)

	var detectionErr error
	detects := make(map[string]*AddressDetectionSpec)
	go func(wg *sync.WaitGroup) {
		for _, spec := range c.Detection {
			if _, exists := detects[spec.Name]; exists {
				detectionErr = fmt.Errorf("detection spec %s already exists", spec.Name)
				break
			}
			if detectionErr = spec.Validate(); detectionErr != nil {
				break
			}
			detects[spec.Name] = spec
		}
		wg.Done()
	}(&validateWg)

	var providerErr error
	providers := make(map[string]*DNSProviderSpec)
	go func(wg *sync.WaitGroup) {
		for _, spec := range c.Provider {
			if _, exists := providers[spec.Name]; exists {
				providerErr = fmt.Errorf("provider spec %s already exists", spec.Name)
				break
			}
			if providerErr = spec.Validate(); providerErr != nil {
				break
			}
			providers[spec.Name] = spec
		}
		wg.Done()
	}(&validateWg)

	validateWg.Wait()

	if ddnsErr != nil {
		return ddnsErr
	}
	if detectionErr != nil {
		return detectionErr
	}
	if providerErr != nil {
		return providerErr
	}

	for k := range ddns {
		ddnsSpec := ddns[k]

		detectionName := ddnsSpec.DetectionRef
		if _, exists := detects[detectionName]; !exists {
			return fmt.Errorf("ddns spec %s referenced unknown detection spec %s", ddnsSpec.Name, detectionName)
		}
		ddnsSpec.detectionSpec = detects[detectionName]

		providerName := ddnsSpec.ProviderRef
		if _, exists := providers[providerName]; !exists {
			return fmt.Errorf("ddns spec %s referenced unknown provider spec %s", ddnsSpec.Name, providerName)
		}
		ddnsSpec.providerSpec = providers[providerName]
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
