package net

import (
	"fmt"
	"net"
	"strings"
)

// NetworkInfo contains information about a network interface
type NetworkInfo struct {
	IP        string
	Interface string
	Type      string // wifi, ethernet, virtual
}

// DetectLocalIP detects the local IP address on the LAN
// It returns the most appropriate private IP address found
func DetectLocalIP() (*NetworkInfo, error) {
	interfaces, err := GetAllInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	if len(interfaces) == 0 {
		return nil, fmt.Errorf("no active network interfaces found")
	}

	selected := PrioritizeInterfaces(interfaces)
	if selected == nil {
		return nil, fmt.Errorf("no suitable private IP address found")
	}

	return selected, nil
}

// GetAllInterfaces returns all network interfaces with valid private IPs
func GetAllInterfaces() ([]NetworkInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []NetworkInfo

	for _, iface := range ifaces {
		// Skip interfaces that are down
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Only consider IPv4 addresses
			if ip == nil || ip.To4() == nil {
				continue
			}

			ipStr := ip.String()

			// Only include private IPs
			if !IsPrivateIP(ipStr) {
				continue
			}

			netInfo := NetworkInfo{
				IP:        ipStr,
				Interface: iface.Name,
				Type:      classifyInterface(iface.Name),
			}

			result = append(result, netInfo)
		}
	}

	return result, nil
}

// IsPrivateIP validates that an IP belongs to RFC 1918 private ranges
// Valid ranges: 192.168.x.x, 10.x.x.x, 172.16-31.x.x
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Convert to 4-byte representation
	ip = ip.To4()
	if ip == nil {
		return false
	}

	// 10.0.0.0/8
	if ip[0] == 10 {
		return true
	}

	// 172.16.0.0/12 (172.16.0.0 - 172.31.255.255)
	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return true
	}

	// 192.168.0.0/16
	if ip[0] == 192 && ip[1] == 168 {
		return true
	}

	return false
}

// PrioritizeInterfaces selects the best interface from a list
// Priority: physical interfaces (wifi, ethernet) over virtual interfaces
func PrioritizeInterfaces(interfaces []NetworkInfo) *NetworkInfo {
	if len(interfaces) == 0 {
		return nil
	}

	var physical []NetworkInfo
	var virtual []NetworkInfo

	for _, iface := range interfaces {
		if iface.Type == "virtual" {
			virtual = append(virtual, iface)
		} else {
			physical = append(physical, iface)
		}
	}

	// Prefer physical interfaces
	if len(physical) > 0 {
		// Among physical, prefer wifi and ethernet
		for _, iface := range physical {
			if iface.Type == "wifi" || iface.Type == "ethernet" {
				return &iface
			}
		}
		// Return first physical if no wifi/ethernet found
		return &physical[0]
	}

	// Fall back to virtual if no physical found
	if len(virtual) > 0 {
		return &virtual[0]
	}

	return nil
}

// classifyInterface determines the type of network interface based on its name
func classifyInterface(name string) string {
	nameLower := strings.ToLower(name)

	// Virtual interfaces
	if strings.HasPrefix(nameLower, "docker") ||
		strings.HasPrefix(nameLower, "veth") ||
		strings.HasPrefix(nameLower, "br-") ||
		strings.HasPrefix(nameLower, "virbr") ||
		strings.HasPrefix(nameLower, "vmnet") ||
		strings.HasPrefix(nameLower, "vbox") {
		return "virtual"
	}

	// WiFi interfaces
	if strings.HasPrefix(nameLower, "wlan") ||
		strings.HasPrefix(nameLower, "wl") ||
		strings.HasPrefix(nameLower, "wifi") ||
		strings.Contains(nameLower, "wi-fi") {
		return "wifi"
	}

	// Ethernet interfaces
	if strings.HasPrefix(nameLower, "eth") ||
		strings.HasPrefix(nameLower, "en") ||
		strings.HasPrefix(nameLower, "em") ||
		strings.HasPrefix(nameLower, "eno") ||
		strings.HasPrefix(nameLower, "enp") ||
		strings.HasPrefix(nameLower, "ens") {
		return "ethernet"
	}

	// Default to ethernet for unknown physical interfaces
	return "ethernet"
}
