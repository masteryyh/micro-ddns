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

func validateAddress(address string) bool {
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}

	return !ip.IsLoopback() &&
		!ip.IsMulticast() &&
		!ip.IsUnspecified() &&
		!ip.IsInterfaceLocalMulticast() &&
		!ip.IsLinkLocalUnicast()
}

func IsPrivate(address string) bool {
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}
	return ip.IsPrivate()
}

func IsValidV4(address string) bool {
	if strings.Contains(address, ":") {
		return false
	}
	return validateAddress(address)
}

func IsValidV6(address string) bool {
	if strings.Count(address, ":") < 2 {
		return false
	}
	return validateAddress(address)
}
