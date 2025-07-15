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
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestProviderSpecValidation tests DNS provider specification validation
func TestProviderSpecValidation(t *testing.T) {
	// Test Cloudflare with API token
	cloudflareAPITokenConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-api-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "cloudflare-api-token.yaml")

	err := os.WriteFile(configFile, []byte(cloudflareAPITokenConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid Cloudflare API token config: %v", err)
	}

	config := GetConfig()
	if config.Provider[0].GetType() != DNSProviderCloudflare {
		t.Errorf("Expected Cloudflare provider type, got %s", config.Provider[0].GetType())
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestCloudflareGlobalAPIKey tests Cloudflare global API key configuration
func TestCloudflareGlobalAPIKey(t *testing.T) {
	cloudflareGlobalAPIConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare:
      globalApiKey: "test-global-key"
      email: "test@example.com"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "cloudflare-global-api.yaml")

	err := os.WriteFile(configFile, []byte(cloudflareGlobalAPIConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid Cloudflare global API config: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestAliCloudProvider tests AliCloud provider configuration
func TestAliCloudProvider(t *testing.T) {
	alicloudConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: alicloud
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: alicloud
    alicloud:
      accessKeyId: "test-access-key"
      accessKeySecret: "test-secret-key"
      line: "default"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "alicloud.yaml")

	err := os.WriteFile(configFile, []byte(alicloudConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid AliCloud config: %v", err)
	}

	config := GetConfig()
	if config.Provider[0].GetType() != DNSProviderAliCloud {
		t.Errorf("Expected AliCloud provider type, got %s", config.Provider[0].GetType())
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestDNSPodProvider tests DNSPod provider configuration
func TestDNSPodProvider(t *testing.T) {
	dnspodConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: dnspod
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: dnspod
    dnspod:
      secretId: "test-secret-id"
      secretKey: "test-secret-key"
      lineId: "0"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "dnspod.yaml")

	err := os.WriteFile(configFile, []byte(dnspodConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid DNSPod config: %v", err)
	}

	config := GetConfig()
	if config.Provider[0].GetType() != DNSProviderDNSPod {
		t.Errorf("Expected DNSPod provider type, got %s", config.Provider[0].GetType())
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestRFC2136Provider tests RFC2136 provider configuration
func TestRFC2136Provider(t *testing.T) {
	rfc2136Config := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: rfc2136
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: rfc2136
    rfc2136:
      address: "192.168.1.1"
      port: 53
      useTcp: false
      tsig:
        keyName: "test-key"
        key: "dGVzdC1rZXk="
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "rfc2136.yaml")

	err := os.WriteFile(configFile, []byte(rfc2136Config), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid RFC2136 config: %v", err)
	}

	config := GetConfig()
	if config.Provider[0].GetType() != DNSProviderRFC2136 {
		t.Errorf("Expected RFC2136 provider type, got %s", config.Provider[0].GetType())
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestHuaweiCloudProvider tests Huawei Cloud provider configuration
func TestHuaweiCloudProvider(t *testing.T) {
	huaweiConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: huawei
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: huawei
    huawei:
      accessKey: "test-access-key"
      secretAccessKey: "test-secret-key"
      region: "cn-north-1"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "huawei.yaml")

	err := os.WriteFile(configFile, []byte(huaweiConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid Huawei Cloud config: %v", err)
	}

	config := GetConfig()
	if config.Provider[0].GetType() != DNSProviderHuaweiCloud {
		t.Errorf("Expected HuaweiCloud provider type, got %s", config.Provider[0].GetType())
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestJDCloudProvider tests JD Cloud provider configuration
func TestJDCloudProvider(t *testing.T) {
	jdConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: jd
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: jd
    jd:
      accessKey: "test-access-key"
      secretKey: "test-secret-key"
      viewId: -1
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "jd.yaml")

	err := os.WriteFile(configFile, []byte(jdConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid JD Cloud config: %v", err)
	}

	config := GetConfig()
	if config.Provider[0].GetType() != DNSProviderJDCloud {
		t.Errorf("Expected JDCloud provider type, got %s", config.Provider[0].GetType())
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInvalidProviderConfigurations tests various invalid provider configurations
func TestInvalidProviderConfigurations(t *testing.T) {
	testCases := []struct {
		name        string
		config      string
		expectError string
	}{
		{
			name: "Missing Cloudflare credentials",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare: {}
`,
			expectError: "must choose between api token or global api key with email",
		},
		{
			name: "Missing email for global API key",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare:
      globalApiKey: "test-key"
`,
			expectError: "must choose between api token or global api key with email",
		},
		{
			name: "Missing AliCloud access key",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: alicloud
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: alicloud
    alicloud:
      accessKeySecret: "test-secret"
`,
			expectError: "AccessKeyID cannot be empty",
		},
		{
			name: "Invalid RFC2136 port",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: rfc2136
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: rfc2136
    rfc2136:
      address: "192.168.1.1"
      port: 70000
`,
			expectError: "port 70000 is invalid",
		},
		{
			name: "Multiple providers in one spec",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: multi
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: multi
    cloudflare:
      apiToken: "test-token"
    alicloud:
      accessKeyId: "test-id"
      accessKeySecret: "test-secret"
`,
			expectError: "only 1 provider can be used within 1 spec",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configFile := filepath.Join(tempDir, "invalid-provider.yaml")

			err := os.WriteFile(configFile, []byte(tc.config), 0644)
			if err != nil {
				t.Fatalf("Failed to create test config file: %v", err)
			}

			err = InitializeConfig(configFile)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tc.name)
			} else if !strings.Contains(err.Error(), tc.expectError) {
				t.Errorf("Expected error containing '%s', got: %v", tc.expectError, err)
			}

			// Reset global config for next test
			globalConfig = nil
			configOnce = sync.Once{}
		})
	}
}
