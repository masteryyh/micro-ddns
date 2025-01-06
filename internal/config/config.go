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

type DNSProvider string

const (
	DNSProviderCloudflare  DNSProvider = "Cloudflare"
	DNSProviderAliCloud    DNSProvider = "AliCloud"
	DNSProviderDNSPod      DNSProvider = "DNSPod"
	DNSProviderHuaweiCloud DNSProvider = "HuaweiCloud"
	DNSProviderJDCloud     DNSProvider = "JDCloud"
	DNSProviderRFC2136     DNSProvider = "RFC2136"
)

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
