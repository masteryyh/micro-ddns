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
	"regexp"
)

const (
	ipv4RegexString = `^((25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)$`
	ipv6RegexString = `^^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:([0-9a-fA-F]{1,4}|:)|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]).){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]).){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
)

var (
	ipv4Regex = regexp.MustCompile(ipv4RegexString)
	ipv6Regex = regexp.MustCompile(ipv6RegexString)

	ipv4InvalidBlocksString = []string{
		"0.0.0.0/8",
		"127.0.0.0/8",
		"255.255.255.255/32",
	}
	ipv4PrivateBlocksString = []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"198.18.0.0/15",
	}

	ipv6InvalidBlocksString = []string{
		"::/128",
		"::1/128",
		"fe80::/10",
		"ff00::/8",
	}
	ipv6PrivateBlocksString = []string{
		"fc00::/7",
	}

	ipv4InvalidBlocks []*net.IPNet
	ipv4PrivateBlocks []*net.IPNet
	ipv6InvalidBlocks []*net.IPNet
	ipv6PrivateBlocks []*net.IPNet
)

func init() {
	for _, block := range ipv4PrivateBlocksString {
		if _, cidr, err := net.ParseCIDR(block); err == nil {
			ipv4PrivateBlocks = append(ipv4PrivateBlocks, cidr)
		}
	}

	for _, block := range ipv4InvalidBlocksString {
		if _, cidr, err := net.ParseCIDR(block); err == nil {
			ipv4InvalidBlocks = append(ipv4InvalidBlocks, cidr)
		}
	}

	for _, block := range ipv6InvalidBlocksString {
		if _, cidr, err := net.ParseCIDR(block); err == nil {
			ipv6InvalidBlocks = append(ipv6InvalidBlocks, cidr)
		}
	}

	for _, block := range ipv6PrivateBlocksString {
		if _, cidr, err := net.ParseCIDR(block); err == nil {
			ipv6PrivateBlocks = append(ipv6PrivateBlocks, cidr)
		}
	}
}

func validateAddressV4(address string) bool {
	if !ipv4Regex.MatchString(address) {
		return false
	}
	return !isLoopbackV4(address)
}

func isLoopbackV4(address string) bool {
	ip := net.ParseIP(address)
	for _, block := range ipv4InvalidBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func isPrivateV4(address string) bool {
	ip := net.ParseIP(address)
	for _, block := range ipv4PrivateBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func isInvalidV6(address string) bool {
	ip := net.ParseIP(address)
	for _, block := range ipv6InvalidBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func isULA(address string) bool {
	ip := net.ParseIP(address)
	for _, block := range ipv6PrivateBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func validateAddressV6(address string) bool {
	if ip := net.ParseIP(address); ip == nil {
		return false
	}
	return !isInvalidV6(address)
}
