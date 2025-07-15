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
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// BenchmarkConfigParsing benchmarks the configuration parsing performance
func BenchmarkConfigParsing(b *testing.B) {
	// Create a complex test configuration
	complexConfig := `ddns:
  - name: homelab-v4
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "TZ=UTC */15 * * * *"
    providerRef: cloudflare
    detectionRef: api-v4
  - name: homelab-v6
    domain: example.com
    subdomain: www
    stack: IPv6
    cron: "TZ=America/New_York */30 * * * *"
    providerRef: cloudflare
    detectionRef: api-v6
  - name: internal-dns
    domain: internal.com
    subdomain: "@"
    stack: IPv4
    cron: "0 */2 * * *"
    providerRef: rfc2136
    detectionRef: interface

detection:
  - name: api-v4
    api:
      url: https://api.ipify.org/
      customHeaders:
        User-Agent: "micro-ddns/test"
  - name: api-v6
    api:
      url: https://api6.ipify.org/
  - name: interface
    interface:
      name: eth0
    selector:
      includeCIDRs:
        - "192.168.0.0/16"
        - "10.0.0.0/8"
      excludeCIDRs:
        - "169.254.0.0/16"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-cloudflare-token"
  - name: rfc2136
    rfc2136:
      address: "192.168.1.1"
      port: 53
      useTcp: false
      tsig:
        keyName: "test-key"
        key: "dGVzdC1rZXk="
`

	tempDir := b.TempDir()
	configFile := filepath.Join(tempDir, "benchmark-config.yaml")

	err := os.WriteFile(configFile, []byte(complexConfig), 0644)
	if err != nil {
		b.Fatalf("Failed to create test config file: %v", err)
	}

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Reset global config for each iteration
		globalConfig = nil
		configOnce = sync.Once{}

		err := InitializeConfig(configFile)
		if err != nil {
			b.Fatalf("InitializeConfig failed: %v", err)
		}

		// Access the config to ensure it's fully loaded
		config := GetConfig()
		if len(config.DDNS) != 3 {
			b.Fatalf("Expected 3 DDNS instances, got %d", len(config.DDNS))
		}
	}
}

// BenchmarkConfigValidation benchmarks the configuration validation performance
func BenchmarkConfigValidation(b *testing.B) {
	// Create a configuration struct for validation
	config := &Config{
		DDNS: []*DDNSSpec{
			{
				Name:         "test",
				Domain:       "example.com",
				Subdomain:    "www",
				Stack:        IPv4,
				Cron:         "*/30 * * * *",
				ProviderRef:  "cloudflare",
				DetectionRef: "api",
			},
		},
		Detection: []*AddressDetectionSpec{
			{
				Name: "api",
				API: &ThirdPartyServiceSpec{
					URL: "https://api.ipify.org/",
				},
			},
		},
		Provider: []*DNSProviderSpec{
			{
				Name: "cloudflare",
				Cloudflare: &CloudflareSpec{
					APIToken: stringPtr("test-token"),
				},
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := config.Validate()
		if err != nil {
			b.Fatalf("Config validation failed: %v", err)
		}
	}
}

// BenchmarkMultipleConfigs benchmarks parsing multiple configuration instances
func BenchmarkMultipleConfigs(b *testing.B) {
	// Configuration with many DDNS instances
	manyInstancesConfig := `ddns:`

	// Generate 50 DDNS instances
	for i := 0; i < 50; i++ {
		manyInstancesConfig += `
  - name: instance-` + fmt.Sprintf("%d", i) + `
    domain: example` + fmt.Sprintf("%d", i%10) + `.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: api`
	}

	manyInstancesConfig += `

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := b.TempDir()
	configFile := filepath.Join(tempDir, "many-instances.yaml")

	err := os.WriteFile(configFile, []byte(manyInstancesConfig), 0644)
	if err != nil {
		b.Fatalf("Failed to create test config file: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Reset global config for each iteration
		globalConfig = nil
		configOnce = sync.Once{}

		err := InitializeConfig(configFile)
		if err != nil {
			b.Fatalf("InitializeConfig failed: %v", err)
		}
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
