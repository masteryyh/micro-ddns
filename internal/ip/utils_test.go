package ip

import (
	"net"
	"testing"
)

func parseCidr(cidr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return ipnet
}

func TestAddressExcluded(t *testing.T) {
	type args struct {
		address  net.IP
		includes []*net.IPNet
		excludes []*net.IPNet
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Included test #0",
			args: args{
				address: net.ParseIP("192.168.31.254"),
				includes: []*net.IPNet{
					parseCidr("192.168.31.0/24"),
				},
			},
			want: false,
		},
		{
			name: "Included test #1",
			args: args{
				address: net.ParseIP("113.88.200.10"),
				includes: []*net.IPNet{
					parseCidr("113.88.0.0/16"),
				},
			},
			want: false,
		},
		{
			name: "Included test #2",
			args: args{
				address: net.ParseIP("192.168.100.10"),
				includes: []*net.IPNet{
					parseCidr("172.18.0.0/16"),
					parseCidr("192.168.1.0/24"),
					parseCidr("192.168.100.0/24"),
				},
			},
			want: false,
		},
		{
			name: "Included test #3",
			args: args{
				address: net.ParseIP("172.16.100.5"),
				includes: []*net.IPNet{
					parseCidr("172.16.100.0/24"),
				},
				excludes: []*net.IPNet{
					parseCidr("172.0.0.0/8"),
				},
			},
			want: false,
		},
		{
			name: "Excluded test #0",
			args: args{
				address: net.ParseIP("10.100.105.20"),
				excludes: []*net.IPNet{
					parseCidr("10.0.0.0/8"),
				},
			},
			want: true,
		},
		{
			name: "Excluded test #1",
			args: args{
				address: net.ParseIP("1.0.0.1"),
				excludes: []*net.IPNet{
					parseCidr("1.0.0.0/24"),
				},
			},
			want: true,
		},
		{
			name: "Excluded test #2",
			args: args{
				address: net.ParseIP("172.16.20.105"),
				excludes: []*net.IPNet{
					parseCidr("192.168.1.0/24"),
					parseCidr("172.0.0.0/8"),
					parseCidr("172.16.0.0/16"),
				},
			},
			want: true,
		},
		{
			name: "Excluded test #3",
			args: args{
				address: net.ParseIP("10.100.50.24"),
				includes: []*net.IPNet{
					parseCidr("10.100.0.0/16"),
				},
				excludes: []*net.IPNet{
					parseCidr("10.100.50.0/24"),
				},
			},
			want: true,
		},
		{
			name: "Excluded test #4",
			args: args{
				address: net.ParseIP("192.168.52.28"),
				includes: []*net.IPNet{
					parseCidr("192.168.0.0/16"),
				},
				excludes: []*net.IPNet{
					parseCidr("192.168.0.0/16"),
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddressExcluded(tt.args.address, tt.args.includes, tt.args.excludes); got != tt.want {
				t.Errorf("AddressExcluded() = %v, want %v", got, tt.want)
			}
		})
	}
}
