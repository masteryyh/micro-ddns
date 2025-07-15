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

// TestAPIDetectionValidation tests API-based address detection validation
func TestAPIDetectionValidation(t *testing.T) {
	// Test valid API detection
	validAPIConfig := `ddns:
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
      jsonPath: ".ip"
      customHeaders:
        User-Agent: "micro-ddns/test"
        Accept: "application/json"
      params:
        format: "json"
      username: "testuser"
      password: "testpass"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "api-detection.yaml")

	err := os.WriteFile(configFile, []byte(validAPIConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid API detection config: %v", err)
	}

	config := GetConfig()
	if config.Detection[0].GetDetectionType() != AddressDetectionThirdParty {
		t.Errorf("Expected ThirdParty detection type, got %s", config.Detection[0].GetDetectionType())
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInterfaceDetectionValidation tests interface-based address detection validation
func TestInterfaceDetectionValidation(t *testing.T) {
	interfaceConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: interface

detection:
  - name: interface
    interface:
      name: eth0
    selector:
      includeCIDRs:
        - "192.168.0.0/16"
        - "10.0.0.0/8"
      excludeCIDRs:
        - "169.254.0.0/16"
        - "127.0.0.0/8"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "interface-detection.yaml")

	err := os.WriteFile(configFile, []byte(interfaceConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid interface detection config: %v", err)
	}

	config := GetConfig()
	if config.Detection[0].GetDetectionType() != AddressDetectionIface {
		t.Errorf("Expected Interface detection type, got %s", config.Detection[0].GetDetectionType())
	}

	// Verify selector configuration
	if config.Detection[0].Selector == nil {
		t.Error("Expected selector configuration, got nil")
	} else {
		includes := config.Detection[0].Selector.GetIncludes()
		excludes := config.Detection[0].Selector.GetExcludes()

		if len(includes) != 2 {
			t.Errorf("Expected 2 include CIDRs, got %d", len(includes))
		}

		if len(excludes) != 2 {
			t.Errorf("Expected 2 exclude CIDRs, got %d", len(excludes))
		}
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestSSHDetectionValidation tests SSH-based address detection validation
func TestSSHDetectionValidation(t *testing.T) {
	// Test SSH with password
	sshPasswordConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      host:
        address: "192.168.1.100"
        port: 22
      credential:
        user: "testuser"
        password: "testpass"
      interface: "eth0"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "ssh-password.yaml")

	err := os.WriteFile(configFile, []byte(sshPasswordConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid SSH password config: %v", err)
	}

	config := GetConfig()
	if config.Detection[0].GetDetectionType() != AddressDetectionSSH {
		t.Errorf("Expected SSH detection type, got %s", config.Detection[0].GetDetectionType())
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestSSHWithPrivateKey tests SSH detection with private key
func TestSSHWithPrivateKey(t *testing.T) {
	sshKeyConfig := `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      host:
        address: "server.example.com"
        port: 2022
      credential:
        user: "admin"
        privateKey: "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC..."
        passphrase: "optional-passphrase"
      command: "custom-ip-command"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "ssh-key.yaml")

	err := os.WriteFile(configFile, []byte(sshKeyConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with valid SSH key config: %v", err)
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestInvalidDetectionConfigurations tests various invalid detection configurations
func TestInvalidDetectionConfigurations(t *testing.T) {
	testCases := []struct {
		name        string
		config      string
		expectError string
	}{
		{
			name: "Missing detection method",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: empty

detection:
  - name: empty

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "must specify a detection method",
		},
		{
			name: "Empty API URL",
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
      url: ""

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "url cannot be empty",
		},
		{
			name: "Empty interface name",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: interface

detection:
  - name: interface
    interface:
      name: ""

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "interface name cannot be empty",
		},
		{
			name: "SSH missing host",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      credential:
        user: "testuser"
        password: "testpass"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "host cannot be empty",
		},
		{
			name: "SSH missing credential",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      host:
        address: "192.168.1.100"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "credential cannot be empty",
		},
		{
			name: "SSH empty address",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      host:
        address: ""
      credential:
        user: "testuser"
        password: "testpass"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "address cannot be empty",
		},
		{
			name: "SSH empty user",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      host:
        address: "192.168.1.100"
      credential:
        user: ""
        password: "testpass"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "user cannot be empty",
		},
		{
			name: "SSH no authentication method",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      host:
        address: "192.168.1.100"
      credential:
        user: "testuser"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "password or private key must be specified",
		},
		{
			name: "SSH both password and private key",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      host:
        address: "192.168.1.100"
      credential:
        user: "testuser"
        password: "testpass"
        privateKey: "test-key"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "only one of password and private key can be specified",
		},
		{
			name: "SSH missing interface for default command",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: ssh

detection:
  - name: ssh
    ssh:
      host:
        address: "192.168.1.100"
      credential:
        user: "testuser"
        password: "testpass"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "interface must be specified if command is not specified",
		},
		{
			name: "Invalid include CIDR",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: interface

detection:
  - name: interface
    interface:
      name: eth0
    selector:
      includeCIDRs:
        - "192.168.1.0/24"
        - "invalid-cidr"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "invalid CIDR invalid-cidr",
		},
		{
			name: "Invalid exclude CIDR",
			config: `ddns:
  - name: test
    domain: example.com
    subdomain: www
    stack: IPv4
    cron: "*/30 * * * *"
    providerRef: cloudflare
    detectionRef: interface

detection:
  - name: interface
    interface:
      name: eth0
    selector:
      excludeCIDRs:
        - "127.0.0.0/8"
        - "192.168.1.1/24"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`,
			expectError: "invalid CIDR 192.168.1.1/24",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configFile := filepath.Join(tempDir, "invalid-detection.yaml")

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

// TestAPIDetectionWithJSONPath tests API detection with JSON path extraction
func TestAPIDetectionWithJSONPath(t *testing.T) {
	jsonPathConfig := `ddns:
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
      url: https://httpbin.org/ip
      jsonPath: ".origin"

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "json-path.yaml")

	err := os.WriteFile(configFile, []byte(jsonPathConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with JSON path config: %v", err)
	}

	config := GetConfig()
	if config.Detection[0].API.JsonPath == nil || *config.Detection[0].API.JsonPath != ".origin" {
		t.Error("Expected JSON path '.origin' to be set")
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}

// TestEmptyJSONPath tests that empty JSON path is handled correctly
func TestEmptyJSONPath(t *testing.T) {
	emptyJSONPathConfig := `ddns:
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
      jsonPath: ""

provider:
  - name: cloudflare
    cloudflare:
      apiToken: "test-token"
`

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "empty-json-path.yaml")

	err := os.WriteFile(configFile, []byte(emptyJSONPathConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = InitializeConfig(configFile)
	if err != nil {
		t.Errorf("InitializeConfig failed with empty JSON path config: %v", err)
	}

	config := GetConfig()
	// Empty JSON path should be converted to nil
	if config.Detection[0].API.JsonPath != nil {
		t.Error("Expected empty JSON path to be nil")
	}

	// Reset global config for next test
	globalConfig = nil
	configOnce = sync.Once{}
}
