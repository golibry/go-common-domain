package web

import (
	"net"
	"strings"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrEmptyIPAddress     = domain.NewError("IP address cannot be empty")
	ErrInvalidIPAddress   = domain.NewError("IP address has invalid format")
	ErrInvalidIPv4Address = domain.NewError("IPv4 address has invalid format")
	ErrInvalidIPv6Address = domain.NewError("IPv6 address has invalid format")
)

type IPAddress struct {
	value string
}

// NewIPAddress creates a new instance of IPAddress with validation and normalization
func NewIPAddress(value string) (IPAddress, error) {
	normalized, err := NormalizeIPAddress(value)
	if err != nil {
		return IPAddress{}, err
	}

	return IPAddress{
		value: normalized,
	}, nil
}

// ReconstituteIPAddress creates a new IPAddress instance without validation or normalization
func ReconstituteIPAddress(value string) IPAddress {
	return IPAddress{
		value: value,
	}
}

// Value returns the IP address value
func (ip IPAddress) Value() string {
	return ip.value
}

// IsIPv4 returns true if the IP address is IPv4
func (ip IPAddress) IsIPv4() bool {
	parsedIP := net.ParseIP(ip.value)
	return parsedIP != nil && parsedIP.To4() != nil
}

// IsIPv6 returns true if the IP address is IPv6
func (ip IPAddress) IsIPv6() bool {
	parsedIP := net.ParseIP(ip.value)
	return parsedIP != nil && parsedIP.To4() == nil
}

// Equals compares two IPAddress objects for equality
func (ip IPAddress) Equals(other IPAddress) bool {
	return ip.value == other.value
}

// String returns a string representation of the IP address
func (ip IPAddress) String() string {
	return ip.value
}

// NormalizeIPAddress normalizes an IP address by trimming spaces and standardizing format
func NormalizeIPAddress(ipAddress string) (string, error) {
	// Trim spaces from the beginning and end
	ipAddress = strings.TrimSpace(ipAddress)

	// Preprocess IPv4 addresses to remove leading zeros
	preprocessed := preprocessIPv4(ipAddress)

	if err := IsValidIPAddress(preprocessed); err != nil {
		return "", err
	}

	// Parse and format to ensure consistent representation
	parsedIP := net.ParseIP(preprocessed)
	if parsedIP == nil {
		return "", ErrInvalidIPAddress
	}

	// For IPv4, ensure standard dotted decimal notation
	if parsedIP.To4() != nil {
		return parsedIP.To4().String(), nil
	}

	// For IPv6, use standard representation
	return parsedIP.String(), nil
}

// IsValidIPAddress validates an IP address (both IPv4 and IPv6)
func IsValidIPAddress(ipAddress string) error {
	if ipAddress == "" {
		return ErrEmptyIPAddress
	}

	parsedIP := net.ParseIP(ipAddress)
	if parsedIP == nil {
		return ErrInvalidIPAddress
	}

	return nil
}

// IsValidIPv4Address validates specifically an IPv4 address
func IsValidIPv4Address(ipAddress string) error {
	if ipAddress == "" {
		return ErrEmptyIPAddress
	}

	parsedIP := net.ParseIP(ipAddress)
	if parsedIP == nil || parsedIP.To4() == nil {
		return ErrInvalidIPv4Address
	}

	return nil
}

// IsValidIPv6Address validates specifically an IPv6 address
func IsValidIPv6Address(ipAddress string) error {
	if ipAddress == "" {
		return ErrEmptyIPAddress
	}

	parsedIP := net.ParseIP(ipAddress)
	if parsedIP == nil || parsedIP.To4() != nil {
		return ErrInvalidIPv6Address
	}

	return nil
}

// preprocessIPv4 removes leading zeros from IPv4 addresses to avoid octal interpretation
func preprocessIPv4(ipAddress string) string {
	// Check if it looks like an IPv4 address (contains dots but not colons)
	if !strings.Contains(ipAddress, ":") && strings.Contains(ipAddress, ".") {
		parts := strings.Split(ipAddress, ".")
		if len(parts) == 4 {
			for i, part := range parts {
				// Remove leading zeros but keep at least one digit
				if len(part) > 1 && part[0] == '0' {
					// Remove leading zeros
					trimmed := strings.TrimLeft(part, "0")
					if trimmed == "" {
						trimmed = "0"
					}
					parts[i] = trimmed
				}
			}
			return strings.Join(parts, ".")
		}
	}
	return ipAddress
}
