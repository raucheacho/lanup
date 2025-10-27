package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		// Valid private IPs - 10.0.0.0/8
		{"10.0.0.0", "10.0.0.0", true},
		{"10.0.0.1", "10.0.0.1", true},
		{"10.255.255.255", "10.255.255.255", true},
		{"10.1.2.3", "10.1.2.3", true},

		// Valid private IPs - 172.16.0.0/12
		{"172.16.0.0", "172.16.0.0", true},
		{"172.16.0.1", "172.16.0.1", true},
		{"172.31.255.255", "172.31.255.255", true},
		{"172.20.10.5", "172.20.10.5", true},

		// Valid private IPs - 192.168.0.0/16
		{"192.168.0.0", "192.168.0.0", true},
		{"192.168.0.1", "192.168.0.1", true},
		{"192.168.255.255", "192.168.255.255", true},
		{"192.168.1.100", "192.168.1.100", true},

		// Invalid - public IPs
		{"8.8.8.8", "8.8.8.8", false},
		{"1.1.1.1", "1.1.1.1", false},
		{"172.15.0.1", "172.15.0.1", false},
		{"172.32.0.1", "172.32.0.1", false},
		{"192.167.0.1", "192.167.0.1", false},
		{"192.169.0.1", "192.169.0.1", false},

		// Invalid - localhost
		{"127.0.0.1", "127.0.0.1", false},

		// Invalid - malformed
		{"invalid", "invalid", false},
		{"", "", false},
		{"256.1.1.1", "256.1.1.1", false},
		{"192.168.1", "192.168.1", false},

		// Edge cases
		{"9.255.255.255", "9.255.255.255", false},
		{"11.0.0.0", "11.0.0.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrivateIP(tt.ip)
			assert.Equal(t, tt.expected, result, "IsPrivateIP(%s) = %v, want %v", tt.ip, result, tt.expected)
		})
	}
}

func TestPrioritizeInterfaces(t *testing.T) {
	tests := []struct {
		name       string
		interfaces []NetworkInfo
		expected   *NetworkInfo
	}{
		{
			name:       "empty list",
			interfaces: []NetworkInfo{},
			expected:   nil,
		},
		{
			name: "single wifi interface",
			interfaces: []NetworkInfo{
				{IP: "192.168.1.100", Interface: "wlan0", Type: "wifi"},
			},
			expected: &NetworkInfo{IP: "192.168.1.100", Interface: "wlan0", Type: "wifi"},
		},
		{
			name: "single ethernet interface",
			interfaces: []NetworkInfo{
				{IP: "192.168.1.100", Interface: "eth0", Type: "ethernet"},
			},
			expected: &NetworkInfo{IP: "192.168.1.100", Interface: "eth0", Type: "ethernet"},
		},
		{
			name: "prefer wifi over virtual",
			interfaces: []NetworkInfo{
				{IP: "172.17.0.1", Interface: "docker0", Type: "virtual"},
				{IP: "192.168.1.100", Interface: "wlan0", Type: "wifi"},
			},
			expected: &NetworkInfo{IP: "192.168.1.100", Interface: "wlan0", Type: "wifi"},
		},
		{
			name: "prefer ethernet over virtual",
			interfaces: []NetworkInfo{
				{IP: "172.17.0.1", Interface: "docker0", Type: "virtual"},
				{IP: "192.168.1.100", Interface: "eth0", Type: "ethernet"},
			},
			expected: &NetworkInfo{IP: "192.168.1.100", Interface: "eth0", Type: "ethernet"},
		},
		{
			name: "first wifi or ethernet from physical",
			interfaces: []NetworkInfo{
				{IP: "192.168.1.50", Interface: "eth0", Type: "ethernet"},
				{IP: "192.168.1.100", Interface: "wlan0", Type: "wifi"},
			},
			expected: &NetworkInfo{IP: "192.168.1.50", Interface: "eth0", Type: "ethernet"},
		},
		{
			name: "multiple physical interfaces - first wifi or ethernet",
			interfaces: []NetworkInfo{
				{IP: "192.168.1.50", Interface: "eth0", Type: "ethernet"},
				{IP: "192.168.1.100", Interface: "wlan0", Type: "wifi"},
				{IP: "10.0.0.1", Interface: "en1", Type: "ethernet"},
			},
			expected: &NetworkInfo{IP: "192.168.1.50", Interface: "eth0", Type: "ethernet"},
		},
		{
			name: "only virtual interfaces",
			interfaces: []NetworkInfo{
				{IP: "172.17.0.1", Interface: "docker0", Type: "virtual"},
				{IP: "172.18.0.1", Interface: "br-123", Type: "virtual"},
			},
			expected: &NetworkInfo{IP: "172.17.0.1", Interface: "docker0", Type: "virtual"},
		},
		{
			name: "complex mix - first physical wifi or ethernet",
			interfaces: []NetworkInfo{
				{IP: "172.17.0.1", Interface: "docker0", Type: "virtual"},
				{IP: "192.168.1.50", Interface: "eth0", Type: "ethernet"},
				{IP: "172.18.0.1", Interface: "veth0", Type: "virtual"},
				{IP: "192.168.1.100", Interface: "wlan0", Type: "wifi"},
			},
			expected: &NetworkInfo{IP: "192.168.1.50", Interface: "eth0", Type: "ethernet"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PrioritizeInterfaces(tt.interfaces)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.IP, result.IP)
				assert.Equal(t, tt.expected.Interface, result.Interface)
				assert.Equal(t, tt.expected.Type, result.Type)
			}
		})
	}
}

func TestClassifyInterface(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		expected      string
	}{
		// Virtual interfaces
		{"docker0", "docker0", "virtual"},
		{"docker1", "docker1", "virtual"},
		{"veth0", "veth0", "virtual"},
		{"veth123abc", "veth123abc", "virtual"},
		{"br-123456", "br-123456", "virtual"},
		{"virbr0", "virbr0", "virtual"},
		{"vmnet0", "vmnet0", "virtual"},
		{"vboxnet0", "vboxnet0", "virtual"},

		// WiFi interfaces
		{"wlan0", "wlan0", "wifi"},
		{"wlan1", "wlan1", "wifi"},
		{"wl0", "wl0", "wifi"},
		{"wifi0", "wifi0", "wifi"},

		// Ethernet interfaces
		{"eth0", "eth0", "ethernet"},
		{"eth1", "eth1", "ethernet"},
		{"en0", "en0", "ethernet"},
		{"en1", "en1", "ethernet"},
		{"em0", "em0", "ethernet"},
		{"eno1", "eno1", "ethernet"},
		{"enp0s3", "enp0s3", "ethernet"},
		{"ens33", "ens33", "ethernet"},

		// Unknown defaults to ethernet
		{"unknown0", "unknown0", "ethernet"},
		{"myinterface", "myinterface", "ethernet"},

		// Case insensitive
		{"WLAN0", "WLAN0", "wifi"},
		{"ETH0", "ETH0", "ethernet"},
		{"DOCKER0", "DOCKER0", "virtual"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyInterface(tt.interfaceName)
			assert.Equal(t, tt.expected, result)
		})
	}
}
