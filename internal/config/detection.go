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
	"net"

	"github.com/masteryyh/micro-ddns/pkg/utils"
)

type AddressDetectionType string

const (
	AddressDetectionIface      AddressDetectionType = "Interface"
	AddressDetectionThirdParty AddressDetectionType = "ThirdParty"
	AddressDetectionSSH        AddressDetectionType = "SSH"
)

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

// SSHHostSpec is the specification of remote SSH host
type SSHHostSpec struct {
	// Address of remote host
	Address string `json:"address" yaml:"address"`

	// Port of remote host, default is 22
	Port *uint16 `json:"port,omitempty" yaml:"port,omitempty"`
}

func (spec *SSHHostSpec) Validate() error {
	if spec.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if spec.Port == nil {
		spec.Port = utils.Uint16Ptr(22)
	}

	if *spec.Port == 0 {
		return fmt.Errorf("port cannot be 0")
	}

	return nil
}

// SSHCredentialSpec is the credential used to connect to remote SSH host
type SSHCredentialSpec struct {
	// User that connect to the host
	User string `json:"user" yaml:"user"`

	// Password of remote host, unsafe since it's trasmitted in plaintext
	Password *string `json:"password,omitempty" yaml:"password,omitempty"`

	// PrivateKey for SSH authentication, recommended
	PrivateKey *string `json:"privateKey,omitempty" yaml:"privateKey,omitempty"`

	// Passphrase for private key if it's encrypted
	Passphrase *string `json:"passphrase,omitempty" yaml:"passphrase,omitempty"`
}

func (spec *SSHCredentialSpec) Validate() error {
	if spec.User == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if utils.IsEmpty(spec.Password) && utils.IsEmpty(spec.PrivateKey) {
		return fmt.Errorf("password or private key must be specified")
	}

	if !utils.IsEmpty(spec.Password) && !utils.IsEmpty(spec.PrivateKey) {
		return fmt.Errorf("only one of password and private key can be specified")
	}

	return nil
}

// SSHDetectionSpec defines how should we connect to remote machine and fetch IP address using SSH
type SSHDetectionSpec struct {
	Host *SSHHostSpec `json:"host" yaml:"host"`

	Credential *SSHCredentialSpec `json:"credential" yaml:"credential"`

	// Command that fetch IP address
	// for IPv4 it's `ip addr show dev <iface> | grep -oE 'inet ([^/]+)' | awk '{print $2}'`
	// for IPv6 it's `ip -6 addr show dev <iface> | grep -oE 'inet6 ([^/]+)' | awk '{print $2}'`
	Command *string `json:"command,omitempty" yaml:"command,omitempty"`

	// Interface that fetch addresses from
	Interface *string `json:"interface,omitempty" yaml:"interface,omitempty"`
}

func (spec *SSHDetectionSpec) Validate() error {
	if spec.Host == nil {
		return fmt.Errorf("host cannot be empty")
	}

	if err := spec.Host.Validate(); err != nil {
		return err
	}

	if spec.Credential == nil {
		return fmt.Errorf("credential cannot be empty")
	}

	if err := spec.Credential.Validate(); err != nil {
		return err
	}

	if utils.IsEmpty(spec.Command) && utils.IsEmpty(spec.Interface) {
		return fmt.Errorf("interface must be specified if command is not specified")
	}
	return nil
}

// IPAddressSelectorSpec defines how should we select an IP address from a set of IP addresses
// If multiple address are found, the first one will be used
type IPAddressSelectorSpec struct {
	// IncludeCIDRs is the list of CIDRs that should be included
	// If address is in one of these CIDRs, it will be selected
	IncludeCIDRs []string `json:"includeCIDRs,omitempty" yaml:"includeCIDRs,omitempty"`

	// ExcludeCIDRs is the list of CIDRs that should be excluded
	ExcludeCIDRs []string `json:"excludeCIDRs,omitempty" yaml:"excludeCIDRs,omitempty"`

	includes []*net.IPNet

	excludes []*net.IPNet
}

func (spec *IPAddressSelectorSpec) GetIncludes() []*net.IPNet {
	if len(spec.includes) == 0 {
		return []*net.IPNet{}
	}
	return spec.includes
}

func (spec *IPAddressSelectorSpec) GetExcludes() []*net.IPNet {
	if len(spec.excludes) == 0 {
		return []*net.IPNet{}
	}
	return spec.excludes
}

func (spec *IPAddressSelectorSpec) Validate() error {
	if len(spec.IncludeCIDRs) > 0 {
		for _, cidr := range spec.IncludeCIDRs {
			ip, n, err := net.ParseCIDR(cidr)
			if err != nil {
				return fmt.Errorf("invalid CIDR %s", cidr)
			}
			if !n.IP.Equal(ip) {
				return fmt.Errorf("invalid CIDR %s", cidr)
			}
			spec.includes = append(spec.includes, n)
		}
	}

	if len(spec.ExcludeCIDRs) > 0 {
		for _, cidr := range spec.ExcludeCIDRs {
			ip, n, err := net.ParseCIDR(cidr)
			if err != nil {
				return fmt.Errorf("invalid CIDR %s", cidr)
			}
			if !n.IP.Equal(ip) {
				return fmt.Errorf("invalid CIDR %s", cidr)
			}
			spec.excludes = append(spec.excludes, n)
		}
	}
	return nil
}

// AddressDetectionSpec defines how should we detect current IP address
type AddressDetectionSpec struct {
	// Name of this address detection specification
	Name string `json:"name" yaml:"name"`

	detectionType AddressDetectionType

	Selector *IPAddressSelectorSpec `json:"selector,omitempty" yaml:"selector,omitempty"`

	Interface *NetworkInterfaceDetectionSpec `json:"interface,omitempty" yaml:"interface,omitempty"`

	API *ThirdPartyServiceSpec `json:"api,omitempty" yaml:"api,omitempty"`

	SSH *SSHDetectionSpec `json:"ssh,omitempty" yaml:"ssh,omitempty"`
}

func (spec *AddressDetectionSpec) Validate() error {
	if spec.Selector != nil {
		if err := spec.Selector.Validate(); err != nil {
			return err
		}
	}

	if spec.Interface != nil {
		spec.detectionType = AddressDetectionIface
		return spec.Interface.Validate()
	} else if spec.API != nil {
		spec.detectionType = AddressDetectionThirdParty
		return spec.API.Validate()
	} else if spec.SSH != nil {
		spec.detectionType = AddressDetectionSSH
		return spec.SSH.Validate()
	}
	return fmt.Errorf("must specify a detection method")
}

func (spec *AddressDetectionSpec) GetDetectionType() AddressDetectionType {
	return spec.detectionType
}
