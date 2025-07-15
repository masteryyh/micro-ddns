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

// TestValidConfig tests a valid configuration
func TestValidConfig(t *testing.T) {
	// Create a temporary valid YAML config file
	validYAMLConfig := `ddns:
  - name: test-instance
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
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "valid-config.yaml")

	err := os.WriteFile(configFile, []byte(validYAMLConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test InitializeConfig
	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid config: %v", err)
	}

	// Test GetConfig
	config := GetConfig()
	if config == nil {
		t.Fatal("GetConfig returned nil")
	}

	// Validate the parsed config
	if len(config.DDNS) != 1 {
		t.Errorf("Expected 1 DDNS instance, got %d", len(config.DDNS))
	}

	if config.DDNS[0].Name != "test-instance" {
		t.Errorf("Expected DDNS name 'test-instance', got '%s'", config.DDNS[0].Name)
	}

	if config.DDNS[0].Domain != "example.com" {
		t.Errorf("Expected domain 'example.com', got '%s'", config.DDNS[0].Domain)
	}

	if config.DDNS[0].Stack != IPv4 {
		t.Errorf("Expected stack IPv4, got %s", config.DDNS[0].Stack)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestValidJSONConfig tests a valid JSON configuration
func TestValidJSONConfig(t *testing.T) {
	validJSONConfig := `{
  "ddns": [
    {
      "name": "json-test",
      "domain": "example.org",
      "subdomain": "api",
      "stack": "IPv6",
      "cron": "0 */6 * * *",
      "providerRef": "alicloud",
      "detectionRef": "interface"
    }
  ],
  "detection": [
    {
      "name": "interface",
      "interface": {
        "name": "eth0"
      }
    }
  ],
  "provider": [
    {
      "name": "alicloud",
      "alicloud": {
        "accessKeyId": "test-id",
        "accessKeySecret": "test-secret"
      }
    }
  ]
}`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "valid-config.json")

	err := os.WriteFile(configFile, []byte(validJSONConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid JSON config: %v", err)
	}

	config := GetConfig()
	if config.DDNS[0].Stack != IPv6 {
		t.Errorf("Expected stack IPv6, got %s", config.DDNS[0].Stack)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInvalidFileExtension tests unsupported file extensions
func TestInvalidFileExtension(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.txt")

	err := os.WriteFile(configFile, []byte("invalid content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for unsupported file extension, got nil")
	}

	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Expected file reading error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestNonExistentFile tests reading a non-existent config file
func TestNonExistentFile(t *testing.T) {
	configFile := "/path/to/nonexistent/config.yaml"

	err := InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Expected file reading error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInvalidYAMLSyntax tests invalid YAML syntax
func TestInvalidYAMLSyntax(t *testing.T) {
	invalidYAML := `ddns:
  - name: test
    domain: example.com
    invalid: [ unclosed bracket
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid-syntax.yaml")

	err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for invalid YAML syntax, got nil")
	}

	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Expected file reading error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInvalidJSONSyntax tests invalid JSON syntax
func TestInvalidJSONSyntax(t *testing.T) {
	invalidJSON := `{
  "ddns": [
    {
      "name": "test",
      "domain": "example.com",
      "invalid": { unclosed brace
    }
  ]
}`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid-syntax.json")

	err := os.WriteFile(configFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for invalid JSON syntax, got nil")
	}

	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Expected file reading error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestEmptyDDNSSection tests config with no DDNS instances
func TestEmptyDDNSSection(t *testing.T) {
	emptyDDNSConfig := `ddns: []

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "empty-ddns.yaml")

	err := os.WriteFile(configFile, []byte(emptyDDNSConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for empty DDNS section, got nil")
	}

	if !strings.Contains(err.Error(), "must have at least 1 ddns spec") {
		t.Errorf("Expected DDNS validation error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestMissingRequiredFields tests config with missing required fields
func TestMissingRequiredFields(t *testing.T) {
	invalidConfig := `ddns:
  - name: ""
    domain: ""
    subdomain: ""
    stack: ""
    cron: ""
    providerRef: ""
    detectionRef: ""

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "missing-fields.yaml")

	err := os.WriteFile(configFile, []byte(invalidConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for missing required fields, got nil")
	}

	// Should contain validation error about empty name
	if !strings.Contains(err.Error(), "name is needed for a DDNS spec") {
		t.Errorf("Expected name validation error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInvalidDomain tests invalid domain format
func TestInvalidDomain(t *testing.T) {
	invalidDomainConfig := `ddns:
  - name: test
    domain: "invalid-domain"
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
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid-domain.yaml")

	err := os.WriteFile(configFile, []byte(invalidDomainConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for invalid domain, got nil")
	}

	if !strings.Contains(err.Error(), "is not a valid domain") {
		t.Errorf("Expected domain validation error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInvalidStack tests invalid stack value
func TestInvalidStack(t *testing.T) {
	invalidStackConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv5
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
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid-stack.yaml")

	err := os.WriteFile(configFile, []byte(invalidStackConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for invalid stack, got nil")
	}

	if !strings.Contains(err.Error(), "is not a valid stack") {
		t.Errorf("Expected stack validation error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInvalidReference tests invalid provider/detection references
func TestInvalidReference(t *testing.T) {
	invalidRefConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: nonexistent-provider
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid-ref.yaml")

	err := os.WriteFile(configFile, []byte(invalidRefConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for invalid provider reference, got nil")
	}

	if !strings.Contains(err.Error(), "referenced unknown provider spec") {
		t.Errorf("Expected provider reference error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestDuplicateNames tests duplicate names in configurations
func TestDuplicateNames(t *testing.T) {
	duplicateNamesConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: api
  - name: test
    domain: example.org
    subdomain: api
    stack: IPv6
    cron: "0 */6 * * *"
    providerRef: cloudflare
    detectionRef: api

detection:
  - name: api
    api:
      url: https://api.ipify.org/

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "duplicate-names.yaml")

	err := os.WriteFile(configFile, []byte(duplicateNamesConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for duplicate DDNS names, got nil")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected duplicate name error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestComplexValidConfig tests a complex but valid configuration
func TestComplexValidConfig(t *testing.T) {
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

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "complex-config.yaml")

	err := os.WriteFile(configFile, []byte(complexConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with complex valid config: %v", err)
		// Reset global config for next test
		globalConfig = nil
		configOnce = sync.Once{}
		return
	}

	config := GetConfig()

	// Validate multiple DDNS instances
	if len(config.DDNS) != 3 {
		t.Errorf("Expected 3 DDNS instances, got %d", len(config.DDNS))
	}

	// Validate multiple detection methods
	if len(config.Detection) != 3 {
		t.Errorf("Expected 3 detection methods, got %d", len(config.Detection))
	}

	// Validate multiple providers
	if len(config.Provider) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(config.Provider))
	}

	// Validate zone apex subdomain
	found := false
	for _, ddns := range config.DDNS {
		if ddns.Name == "internal-dns" && ddns.Subdomain == "@" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find zone apex configuration")
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestEnvironmentVariableConfig tests configuration with environment variables
func TestEnvironmentVariableConfig(t *testing.T) {
	// Set environment variables
	os.Setenv("MICRO_DDNS_DDNS_0_NAME", "env-test")
	os.Setenv("MICRO_DDNS_PROVIDER_0_CLOUDFLARE_APITOKEN", "env-token")
	defer func() {
		os.Unsetenv("MICRO_DDNS_DDNS_0_NAME")
		os.Unsetenv("MICRO_DDNS_PROVIDER_0_CLOUDFLARE_APITOKEN")
	}()

	envConfig := `ddns:
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
      apiToken: "default-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "env-config.yaml")

	err := os.WriteFile(configFile, []byte(envConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with environment config: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestGetConfigPanic tests GetConfig panic when not initialized
func TestGetConfigPanic(t *testing.T) {
	// Ensure global config is nil
	globalConfig = nil
	configOnce = sync.Once{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling GetConfig without initialization")
		}
	}()

	GetConfig()
}

// TestEdgeCaseSubdomains tests edge case subdomain validations
func TestEdgeCaseSubdomains(t *testing.T) {
	testCases := []struct {
		name      string
		subdomain string
		expectErr bool
	}{
		{"zone apex", "@", false},
		{"simple subdomain", "www", false},
		{"multi-level subdomain", "api.v1", false},
		{"numeric subdomain", "123", false},
		{"empty subdomain", "", true},
		{"invalid characters", "sub-domain", false},
		{"special characters", "sub_domain", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configContent := `ddns:
  - name: test
    domain: example.com
    subdomain: "` + tc.subdomain + `"
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
      apiToken: "test-token"
`

			tempDir := t.TempDir()
			configFile := filepath.Join(tempDir, "subdomain-test.yaml")

			err := os.WriteFile(configFile, []byte(configContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test config file: %v", err)
			}

			err = InitializeConfig(configFile)

			if tc.expectErr && err == nil {
				t.Errorf("Expected error for subdomain '%s', got nil", tc.subdomain)
			} else if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error for subdomain '%s': %v", tc.subdomain, err)
			}

			// Reset global config for next test
			globalConfig = nil
			configOnce = sync.Once{}
		})
	}
}

// TestProviderValidations tests different provider configurations
func TestProviderValidations(t *testing.T) {
	// Test missing Cloudflare credentials
	missingCreds := `ddns:
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
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "missing-creds.yaml")

	err := os.WriteFile(configFile, []byte(missingCreds), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for missing Cloudflare credentials, got nil")
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestDetectionValidations tests different detection configurations
func TestDetectionValidations(t *testing.T) {
	// Test missing detection method
	missingMethod := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: api

detection:
  - name: api

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "missing-method.yaml")

	err := os.WriteFile(configFile, []byte(missingMethod), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err == nil {
		t.Error("Expected error for missing detection method, got nil")
	}

	if !strings.Contains(err.Error(), "must specify a detection method") {
		t.Errorf("Expected detection method error, got: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}
