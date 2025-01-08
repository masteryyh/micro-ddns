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

package detection

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"strings"

	"github.com/masteryyh/micro-ddns/internal/config"
)

type IfaceAddressDetector struct {
	interfaceName string
	stack         config.NetworkStack
	includes      []*net.IPNet
	excludes      []*net.IPNet
	logger        *slog.Logger
}

func NewIfaceAddressDetector(detectionSpec *config.AddressDetectionSpec, stack config.NetworkStack, logger *slog.Logger) *IfaceAddressDetector {
	spec := detectionSpec.Interface

	includes := []*net.IPNet{}
	if detectionSpec.Selector != nil {
		includes = detectionSpec.Selector.GetIncludes()
	}

	excludes := []*net.IPNet{}
	if detectionSpec.Selector != nil {
		excludes = detectionSpec.Selector.GetExcludes()
	}

	return &IfaceAddressDetector{
		interfaceName: spec.Name,
		stack:         stack,
		includes:      includes,
		excludes:      excludes,
		logger:        logger,
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

	var validAddresses []net.IP
	for _, addr := range addrs {
		address := strings.Split(addr.String(), "/")[0]
		var ip net.IP
		if v4 {
			ip = IsValidV4(address)
		} else {
			ip = IsValidV6(address)
		}
		if ip == nil {
			d.logger.Debug("ignoring invalid address", "address", address)
			continue
		}
		validAddresses = append(validAddresses, ip)
	}

	for _, valid := range validAddresses {
		if !AddressExcluded(valid, d.includes, d.excludes) {
			return valid.String(), nil
		}
	}
	return "", errors.New("no valid address found")
}

func (d *IfaceAddressDetector) Detect(_ context.Context) (string, error) {
	if d.stack == config.IPv6 {
		return d.detect(false)
	}
	return d.detect(true)
}
