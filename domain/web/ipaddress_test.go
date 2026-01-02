package web

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IPAddressTestSuite struct {
	suite.Suite
}

func TestIPAddressSuite(t *testing.T) {
	suite.Run(t, new(IPAddressTestSuite))
}

func (s *IPAddressTestSuite) TestItCanBuildNewIPAddressWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid IPv4 address",
			input:    "192.168.1.1",
			expected: "192.168.1.1",
		},
		{
			name:     "valid IPv4 address with leading zeros",
			input:    "192.168.001.001",
			expected: "192.168.1.1",
		},
		{
			name:     "localhost IPv4",
			input:    "127.0.0.1",
			expected: "127.0.0.1",
		},
		{
			name:     "broadcast IPv4",
			input:    "255.255.255.255",
			expected: "255.255.255.255",
		},
		{
			name:     "zero IPv4",
			input:    "0.0.0.0",
			expected: "0.0.0.0",
		},
		{
			name:     "valid IPv6 address",
			input:    "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: "2001:db8:85a3::8a2e:370:7334",
		},
		{
			name:     "IPv6 with double colon",
			input:    "2001:db8::1",
			expected: "2001:db8::1",
		},
		{
			name:     "IPv6 localhost",
			input:    "::1",
			expected: "::1",
		},
		{
			name:     "IPv6 zero address",
			input:    "::",
			expected: "::",
		},
		{
			name:     "IPv6 full form",
			input:    "fe80:0000:0000:0000:0000:0000:0000:0001",
			expected: "fe80::1",
		},
		{
			name:     "IPv4 with spaces gets trimmed",
			input:    "  192.168.1.1  ",
			expected: "192.168.1.1",
		},
		{
			name:     "IPv6 with spaces gets trimmed",
			input:    "  2001:db8::1  ",
			expected: "2001:db8::1",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				ipAddress, err := NewIPAddress(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, ipAddress.Value())
				s.Equal(tc.expected, ipAddress.String())
			},
		)
	}
}

func (s *IPAddressTestSuite) TestItFailsToBuildNewIPAddressFromInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "empty IP address",
			input:         "",
			expectedError: ErrEmptyIPAddress,
		},
		{
			name:          "only spaces",
			input:         "   ",
			expectedError: ErrEmptyIPAddress,
		},
		{
			name:          "invalid IPv4 - too many octets",
			input:         "192.168.1.1.1",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "invalid IPv4 - octet too large",
			input:         "192.168.1.256",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "invalid IPv4 - negative octet",
			input:         "192.168.1.-1",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "invalid IPv4 - non-numeric",
			input:         "192.168.1.abc",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "invalid IPv4 - missing octets",
			input:         "192.168.1",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "invalid IPv6 - too many groups",
			input:         "2001:0db8:85a3:0000:0000:8a2e:0370:7334:extra",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "invalid IPv6 - invalid characters",
			input:         "2001:0db8:85a3:gggg:0000:8a2e:0370:7334",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "invalid IPv6 - multiple double colons",
			input:         "2001::db8::1",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "completely invalid format",
			input:         "not.an.ip.address",
			expectedError: ErrInvalidIPAddress,
		},
		{
			name:          "random text",
			input:         "hello world",
			expectedError: ErrInvalidIPAddress,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewIPAddress(tc.input)
				s.Error(err)
				s.True(errors.Is(err, tc.expectedError))
			},
		)
	}
}

func (s *IPAddressTestSuite) TestIPAddressNormalization() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "IPv4 with leading zeros",
			input:    "192.168.001.001",
			expected: "192.168.1.1",
		},
		{
			name:     "IPv6 compression",
			input:    "2001:0db8:0000:0000:0000:0000:0000:0001",
			expected: "2001:db8::1",
		},
		{
			name:     "IPv6 case normalization",
			input:    "2001:DB8::1",
			expected: "2001:db8::1",
		},
		{
			name:     "spaces trimmed",
			input:    "  192.168.1.1  ",
			expected: "192.168.1.1",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				normalized, err := NormalizeIPAddress(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, normalized)
			},
		)
	}
}

func (s *IPAddressTestSuite) TestIPv4Detection() {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid IPv4",
			input:    "192.168.1.1",
			expected: true,
		},
		{
			name:     "localhost IPv4",
			input:    "127.0.0.1",
			expected: true,
		},
		{
			name:     "IPv6 address",
			input:    "2001:db8::1",
			expected: false,
		},
		{
			name:     "IPv6 localhost",
			input:    "::1",
			expected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				ipAddress, err := NewIPAddress(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, ipAddress.IsIPv4())
			},
		)
	}
}

func (s *IPAddressTestSuite) TestIPv6Detection() {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid IPv6",
			input:    "2001:db8::1",
			expected: true,
		},
		{
			name:     "IPv6 localhost",
			input:    "::1",
			expected: true,
		},
		{
			name:     "IPv4 address",
			input:    "192.168.1.1",
			expected: false,
		},
		{
			name:     "localhost IPv4",
			input:    "127.0.0.1",
			expected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				ipAddress, err := NewIPAddress(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, ipAddress.IsIPv6())
			},
		)
	}
}

func (s *IPAddressTestSuite) TestEquals() {
	ip1, _ := NewIPAddress("192.168.1.1")
	ip2, _ := NewIPAddress("192.168.1.1")
	ip3, _ := NewIPAddress("192.168.1.2")

	s.True(ip1.Equals(ip2))
	s.False(ip1.Equals(ip3))
}

func (s *IPAddressTestSuite) TestString() {
	ip, _ := NewIPAddress("192.168.1.1")
	s.Equal("192.168.1.1", ip.String())
}

func (s *IPAddressTestSuite) TestJSONSerialization() {
	ip, _ := NewIPAddress("192.168.1.1")
	data, err := json.Marshal(ip)
	s.NoError(err)
	s.JSONEq(`{}`, string(data))
}

func (s *IPAddressTestSuite) TestReconstitute() {
	ip := ReconstituteIPAddress("192.168.1.1")
	s.Equal("192.168.1.1", ip.Value())
}

func (s *IPAddressTestSuite) TestIsValidIPAddress() {
	testCases := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "valid IPv4",
			input:     "192.168.1.1",
			shouldErr: false,
		},
		{
			name:      "valid IPv6",
			input:     "2001:db8::1",
			shouldErr: false,
		},
		{
			name:      "empty string",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "invalid format",
			input:     "not.an.ip",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				err := IsValidIPAddress(tc.input)
				if tc.shouldErr {
					s.Error(err)
				} else {
					s.NoError(err)
				}
			},
		)
	}
}

func (s *IPAddressTestSuite) TestIsValidIPv4Address() {
	testCases := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "valid IPv4",
			input:     "192.168.1.1",
			shouldErr: false,
		},
		{
			name:      "IPv6 address",
			input:     "2001:db8::1",
			shouldErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "invalid format",
			input:     "not.an.ip",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				err := IsValidIPv4Address(tc.input)
				if tc.shouldErr {
					s.Error(err)
				} else {
					s.NoError(err)
				}
			},
		)
	}
}

func (s *IPAddressTestSuite) TestIsValidIPv6Address() {
	testCases := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "valid IPv6",
			input:     "2001:db8::1",
			shouldErr: false,
		},
		{
			name:      "IPv4 address",
			input:     "192.168.1.1",
			shouldErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "invalid format",
			input:     "not.an.ip",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				err := IsValidIPv6Address(tc.input)
				if tc.shouldErr {
					s.Error(err)
				} else {
					s.NoError(err)
				}
			},
		)
	}
}
