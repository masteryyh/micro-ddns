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

package ip

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/masteryyh/micro-ddns/internal/config"
)

type IfaceAddressDetector struct {
	interfaceName      string
	localAddressPolicy config.LocalAddressPolicy
	stack              config.NetworkStack
	logger             *slog.Logger
}

func NewIfaceAddressDetector(detectionSpec config.AddressDetectionSpec, stack config.NetworkStack, logger *slog.Logger) *IfaceAddressDetector {
	spec := detectionSpec.Interface
	var policy config.LocalAddressPolicy
	if detectionSpec.LocalAddressPolicy == nil {
		policy = config.LocalAddressPolicyIgnore
	}

	logger.Debug("watching network interface", "interface", spec.Name)
	return &IfaceAddressDetector{
		interfaceName:      spec.Name,
		localAddressPolicy: policy,
		stack:              stack,
		logger:             logger,
	}
}

func (d *IfaceAddressDetector) detect(v4 bool) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var ifaceNeeded net.Interface
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 1 && iface.Name == d.interfaceName {
			ifaceNeeded = iface
			break
		}
	}

	addrs, err := ifaceNeeded.Addrs()
	if err != nil {
		return "", err
	}

	var privateIPs, publicIPs []string
	for _, addr := range addrs {
		address := strings.Split(addr.String(), "/")[0]
		if v4 {
			if !IsValidV4(address) {
				d.logger.Debug("ignoring invalid address", "address", address)
				continue
			}

			if IsPrivate(address) {
				d.logger.Debug("saving private IPv4 address", "address", address)
				privateIPs = append(privateIPs, address)
				continue
			}
			d.logger.Debug("saving public IPv4 address", "address", address)
			publicIPs = append(publicIPs, address)
		} else {
			if !IsValidV6(address) {
				d.logger.Debug("ignoring invalid address", "address", address)
				continue
			}

			if IsPrivate(address) {
				d.logger.Debug("saving private IPv6 address", "address", address)
				privateIPs = append(privateIPs, address)
				continue
			}
			d.logger.Debug("saving public IPv6 address", "address", address)
			publicIPs = append(publicIPs, address)
		}
	}

	var validAddr string
	switch d.localAddressPolicy {
	case config.LocalAddressPolicyAllow:
		if len(publicIPs) == 0 {
			if len(privateIPs) == 0 {
				return "", fmt.Errorf("no valid address found")
			}
			validAddr = privateIPs[0]
		}
		validAddr = publicIPs[0]
	case config.LocalAddressPolicyPrefer:
		if len(privateIPs) == 0 {
			if len(publicIPs) == 0 {
				return "", fmt.Errorf("no valid address found")
			}
			validAddr = publicIPs[0]
		}
		validAddr = privateIPs[0]
	default:
	case config.LocalAddressPolicyIgnore:
		if len(publicIPs) == 0 {
			return "", fmt.Errorf("no valid public address found")
		}
		validAddr = publicIPs[0]
	}

	d.logger.Debug("address selected", "address", validAddr)
	return validAddr, nil
}

func (d *IfaceAddressDetector) Detect(_ context.Context) (string, error) {
	if d.stack == config.IPv6 {
		return d.detect(false)
	}
	return d.detect(true)
}
