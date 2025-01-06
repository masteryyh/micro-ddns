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
	"net"
	"strings"
)

func validateAddress(address string) net.IP {
	ip := net.ParseIP(address)
	if ip == nil {
		return nil
	}

	if !ip.IsLoopback() &&
		!ip.IsMulticast() &&
		!ip.IsUnspecified() &&
		!ip.IsInterfaceLocalMulticast() &&
		!ip.IsLinkLocalUnicast() {
		return ip
	}
	return nil
}

func IsValidV4(address string) net.IP {
	if strings.Contains(address, ":") {
		return nil
	}
	return validateAddress(address)
}

func IsValidV6(address string) net.IP {
	if strings.Count(address, ":") < 2 {
		return nil
	}
	return validateAddress(address)
}

func AddressExcluded(address net.IP, includes []*net.IPNet, excludes []*net.IPNet) bool {
	if len(includes) == 0 && len(excludes) == 0 {
		return true
	}

	if len(includes) == 0 {
		for _, exclude := range excludes {
			if exclude.Contains(address) {
				return true
			}
		}
		return false
	}

	if len(excludes) == 0 {
		for _, include := range includes {
			if include.Contains(address) {
				return false
			}
		}
		return true
	}

	var includeCidr *net.IPNet
	for _, include := range includes {
		if include.Contains(address) {
			includeCidr = include
			break
		}
	}

	if includeCidr == nil {
		return false
	}

	var excludeCidr *net.IPNet
	for _, exclude := range excludes {
		if exclude.Contains(address) {
			excludeCidr = exclude
			break
		}
	}
	if excludeCidr == nil {
		return true
	}

	return includeCidr.Contains(excludeCidr.IP)
}
