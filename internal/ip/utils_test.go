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

import "testing"

func TestValidIsLoopbackV4(t *testing.T) {
	validIPs := []string{
		"127.0.0.1",
		"127.0.8.1",
		"0.0.0.0",
		"0.1.1.1",
		"255.255.255.255",
	}

	for _, ip := range validIPs {
		if !isLoopbackV4(ip) {
			t.Errorf("should be loopback but it says not: %s", ip)
		}
	}
}

func TestInvalidIsLoopbackV4(t *testing.T) {
	invalidIPs := []string{
		"192.168.31.1",
		"114.114.114.114",
		"255.255.255.254",
		"255.255.255.0",
		"128.0.0.1",
		"1.0.0.1",
		"1.1.1.1",
		"8.8.8.8",
	}

	for _, ip := range invalidIPs {
		if isLoopbackV4(ip) {
			t.Errorf("should not be loopback but it says is: %s", ip)
		}
	}
}

func TestValidValidateAddressV4(t *testing.T) {
	validIPs := []string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.0",
		"8.8.8.8",
		"1.1.1.1",
		"192.168.31.254",
		"114.114.114.114",
		"169.254.0.1",
		"223.255.255.254",
		"100.64.0.1",
	}

	for _, ip := range validIPs {
		if !validateAddressV4(ip) {
			t.Errorf("should be valid address but it says not: %s", ip)
		}
	}
}

func TestInvalidValidateAddressV4(t *testing.T) {
	invalidIPs := []string{
		"192.168.1. 1",
		" 192.168.1.1",
		"192.168.1.1\n",
		"192.168.1.1#",
		"192.168.1.1/24",
		"192.168.1.1.0",
		".192.168.1.1",
		"192.168.1.1.",
		"192.168.1.01",
		"300.168.1.1",
	}

	for _, ip := range invalidIPs {
		if validateAddressV4(ip) {
			t.Errorf("should be invalid address but it says valid: %s", ip)
		}
	}
}

func TestValidIsPrivateV4(t *testing.T) {
	validIPs := []string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.0",
		"192.168.32.25",
		"10.244.0.1",
		"172.18.0.224",
		"172.16.7.2",
		"192.168.100.254",
		"192.168.100.121",
		"169.254.31.20",
		"198.18.0.1",
	}

	for _, ip := range validIPs {
		if !isPrivateV4(ip) {
			t.Errorf("should be private ipv4 but it says not: %s", ip)
		}
	}
}

func TestInvalidIsPrivateV4(t *testing.T) {
	invalidIPs := []string{
		"8.8.8.8",
		"140.82.114.4",
		"64.233.187.99",
		"192.0.2.1",
		"198.51.100.1",
		"203.0.113.1",
		"172.15.255.255",
		"172.32.0.0",
		"192.167.255.255",
		"192.169.0.0",
		"11.0.0.1",
		"41.77.60.1",
		"74.125.200.106",
		"143.95.32.84",
		"157.240.1.35",
		"173.194.39.78",
		"184.150.153.75",
		"199.19.85.1",
		"207.148.248.143",
		"216.58.212.238",
	}

	for _, ip := range invalidIPs {
		if isPrivateV4(ip) {
			t.Errorf("should not be private ipv4 but it says is: %s", ip)
		}
	}
}

func TestValidIsInvalidV6(t *testing.T) {
	data := []string{
		"::",
		"::1",
		"fe80::",
		"fe80::1",
		"fe80::1234",
		"fe80::abcd:efff:cccc:dddd",
		"febf:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"ff00::",
		"ff02::1",
		"ff02::2",
		"ff05::1:3",
		"ff01::101",
		"ff00::abcd:1234",
		"ff10::1234:5678",
		"ff0f:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"fe80::aede:48ff:fe00:1122",
		"ff00::c633:6400",
		"ff02::5",
		"ff08::1",
	}

	for _, ip := range data {
		if !isInvalidV6(ip) {
			t.Errorf("should not be invalid ipv6 but it says is: %s", ip)
		}
	}
}

func TestInvalidIsInvalidV6(t *testing.T) {
	data := []string{
		"2001:db8::1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2607:f0d0:1002:51::4",
		"2001:db8:abcd:0012::0",
		"2001:db8:1234:0000:0000:0000:0000:0001",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		"2607:f0d0:1002:0051:0000:0000:0000:0004",
		"2001:0db8:abcd:001f::1",
		"2001:db8:1234:ffff:ffff:ffff:ffff:ffff",
		"2001:db8:0000:abcd:0000:0000:2345:6789",
		"3ffe:1900:4545:3:200:f8ff:fe21:67cf",
		"2001:0db8:0123:4567:89ab:cdef:1234:5678",
		"2002:c0a8:6301::c0a8:6301",
		"2001:0db8:0000:0000:0000:ff00:0042:8329",
		"2001:db8:abcd:0012:0000:0000:0000:0000",
		"2001:db8:1234:5678:9abc:def0:1234:5678",
		"2607:f0d0:1002:0051:0000:0000:0000:0001",
		"2001:0db8:3333:4444:5555:6666:7777:8888",
		"2001:db8:0000:00ff:0000:0000:0000:fffe",
		"2001:db8:abcd:ef00:ffff:0000:0000:ffff",
	}

	for _, ip := range data {
		if isInvalidV6(ip) {
			t.Errorf("should not be invalid ipv6 but it says is: %s", ip)
		}
	}
}

func TestValidIsULA(t *testing.T) {
	data := []string{
		"fc00:1234:5678:9abc:def0:1234:5678:9abc",
		"fc01:abcd:ef01:2345:6789:abcd:ef01:2345",
		"fc12:3456:789a:bcde:f012:3456:789a:bcde",
		"fc23:4567:89ab:cdef:0123:4567:89ab:cdef",
		"fc34:5678:9abc:def0:1234:5678:9abc:def0",
		"fc45:6789:abcd:ef01:2345:6789:abcd:ef01",
		"fc56:789a:bcde:f012:3456:789a:bcde:f012",
		"fc67:89ab:cdef:0123:4567:89ab:cdef:0123",
		"fc78:9abc:def0:1234:5678:9abc:def0:1234",
		"fc89:abcd:ef01:2345:6789:abcd:ef01:2345",
		"fc9a:bcde:f012:3456:789a:bcde:f012:3456",
		"fcab:cdef:0123:4567:89ab:cdef:0123:4567",
		"fcbc:def0:1234:5678:9abc:def0:1234:5678",
		"fccd:ef01:2345:6789:abcd:ef01:2345:6789",
		"fcde:f012:3456:789a:bcde:f012:3456:789a",
		"fcef:0123:4567:89ab:cdef:0123:4567:89ab",
		"fcf0:1234:5678:9abc:def0:1234:5678:9abc",
		"fcf1:2345:6789:abcd:ef01:2345:6789:abcd",
		"fcf2:3456:789a:bcde:f012:3456:789a:bcde",
		"fcf3:4567:89ab:cdef:0123:4567:89ab:cdef",
	}

	for _, ip := range data {
		if !isULA(ip) {
			t.Errorf("should be ULA but it says not: %s", ip)
		}
	}
}

// TODO: fix unit tests
func TestValidValidateAddressV6(t *testing.T) {
	data := []string{
		"2001:db8::1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2001:db8:1234:5678:9abc:def0:1234:5678",
		"2607:f0d0:1002:51::4",
		"2001:0db8:abcd:0012::0",
		"2001:db8:1234::",
		"2001:db8:abcd:ef00:ffff:0000:0000:ffff",
		"3ffe:1900:4545:3:200:f8ff:fe21:67cf",
		"2001:0db8::1234:5678:9abc:def0:1234:5678",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		"2607:f0d0:1002:0051:0000:0000:0000:0004",
		"2001:0db8:abcd:001f::1",
		"2001:db8:1234:ffff:ffff:ffff:ffff:ffff",
		"2001:db8:0000:abcd:0000:0000:2345:6789",
		"2001:0db8:0123:4567:89ab:cdef:1234:5678",
		"2002:c0a8:6301::c0a8:6301",
		"2001:0db8:0000:0000:0000:ff00:0042:8329",
		"2001:db8:abcd:0012:0000:0000:0000:0000",
		"2001:db8:1234:5678:9abc:def0:1234:5678",
		"2607:f0d0:1002:0051:0000:0000:0000:0001",
	}

	for _, ip := range data {
		if !validateAddressV6(ip) {
			t.Errorf("should be valid IPv6 but it says not: %s", ip)
		}
	}
}

// TODO: fix unit tests
func TestInvalidValidateAddressV6(t *testing.T) {
	data := []string{
		"2001::85a3::0370:7334",
		"2001:db8:1234:5678:9abc:def0:1234:5678:1000",
		"2001:db8:gabc:0000",
		"2001:db8::12345",
		"::1234::5678",
		"2001:db8::xyz",
		"fe80::1234::abcd",
		"2001:db8:1234:5678:9abc:def0:1234:5678:1234",
		"12345::6789",
		"fe80:::1234",
		"%eth0:2001:db8::1234",
		"2001:db8::g123",
		"2001:db8::12345678",
		"::1::",
		"fe80::1234::",
		"2001:db8::88888",
		"fe80::abcd:ef01:2345:6789:abcd",
		"2001:db8::11111",
		"2001:db8:1234:5678::abcd::ef01",
		"2001:db8:0000:00ff::1234:5678:9abc",
	}

	for _, ip := range data {
		if validateAddressV6(ip) {
			t.Errorf("should be invalid IPv6 but it says is: %s", ip)
		}
	}
}
