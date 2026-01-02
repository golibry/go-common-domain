package web

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DomainNameTestSuite struct {
	suite.Suite
}

func TestDomainNameSuite(t *testing.T) {
	suite.Run(t, new(DomainNameTestSuite))
}

func (s *DomainNameTestSuite) TestItCanBuildNewDomainNameWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple domain",
			input:    "example.com",
			expected: "example.com",
		},
		{
			name:     "subdomain",
			input:    "www.example.com",
			expected: "www.example.com",
		},
		{
			name:     "multiple subdomains",
			input:    "api.v1.example.com",
			expected: "api.v1.example.com",
		},
		{
			name:     "domain with numbers",
			input:    "test123.example.org",
			expected: "test123.example.org",
		},
		{
			name:     "domain with hyphens",
			input:    "my-site.example-domain.net",
			expected: "my-site.example-domain.net",
		},
		{
			name:     "single character domain",
			input:    "a.b",
			expected: "a.b",
		},
		{
			name:     "uppercase gets normalized",
			input:    "EXAMPLE.COM",
			expected: "example.com",
		},
		{
			name:     "mixed case gets normalized",
			input:    "MyDomain.Example.COM",
			expected: "mydomain.example.com",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				domainName, err := NewDomainName(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, domainName.Value())
				s.Equal(tc.expected, domainName.String())
			},
		)
	}
}

func (s *DomainNameTestSuite) TestItFailsToBuildNewDomainNameFromInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "empty domain",
			input:         "",
			expectedError: ErrEmptyDomainName,
		},
		{
			name:          "only spaces",
			input:         "   ",
			expectedError: ErrEmptyDomainName,
		},
		{
			name:          "starts with dot",
			input:         ".example.com",
			expectedError: ErrStartsOrEndsWithDot,
		},
		{
			name:          "ends with dot",
			input:         "example.com.",
			expectedError: ErrStartsOrEndsWithDot,
		},
		{
			name:          "consecutive dots",
			input:         "example..com",
			expectedError: ErrConsecutiveDots,
		},
		{
			name:          "starts with hyphen",
			input:         "-example.com",
			expectedError: ErrStartsOrEndsWithHyphen,
		},
		{
			name:          "ends with hyphen",
			input:         "example-.com",
			expectedError: ErrStartsOrEndsWithHyphen,
		},
		{
			name:          "label starts with hyphen",
			input:         "example.-test.com",
			expectedError: ErrStartsOrEndsWithHyphen,
		},
		{
			name:          "label ends with hyphen",
			input:         "example.test-.com",
			expectedError: ErrStartsOrEndsWithHyphen,
		},
		{
			name:          "invalid characters",
			input:         "example@.com",
			expectedError: ErrInvalidDomainNameChars,
		},
		{
			name:          "underscore not allowed",
			input:         "example_test.com",
			expectedError: ErrInvalidDomainNameChars,
		},
		{
			name:          "space in domain",
			input:         "example .com",
			expectedError: ErrInvalidDomainNameChars,
		},
		{
			name:          "too long domain",
			input:         strings.Repeat("a", 250) + ".com",
			expectedError: ErrTooLongDomainName,
		},
		{
			name:          "too long label",
			input:         strings.Repeat("a", 64) + ".com",
			expectedError: ErrTooLongDomainLabel,
		},
		{
			name:          "only dot",
			input:         ".",
			expectedError: ErrStartsOrEndsWithDot,
		},
		{
			name:          "only hyphen",
			input:         "-",
			expectedError: ErrStartsOrEndsWithHyphen,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewDomainName(tc.input)
				s.Error(err)
				s.True(
					errors.Is(err, tc.expectedError),
					"Expected error %v, got %v",
					tc.expectedError,
					err,
				)
			},
		)
	}
}

func (s *DomainNameTestSuite) TestDomainNameNormalization() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase to lowercase",
			input:    "EXAMPLE.COM",
			expected: "example.com",
		},
		{
			name:     "mixed case to lowercase",
			input:    "MyDomain.Example.COM",
			expected: "mydomain.example.com",
		},
		{
			name:     "trim spaces",
			input:    "  example.com  ",
			expected: "example.com",
		},
		{
			name:     "trim spaces and normalize case",
			input:    "  EXAMPLE.COM  ",
			expected: "example.com",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				normalized, err := NormalizeDomainName(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, normalized)
			},
		)
	}
}

func (s *DomainNameTestSuite) TestEquals() {
	domain1, _ := NewDomainName("example.com")
	domain2, _ := NewDomainName("example.com")
	domain3, _ := NewDomainName("different.com")

	s.True(domain1.Equals(domain2))
	s.False(domain1.Equals(domain3))
}

func (s *DomainNameTestSuite) TestString() {
	domain, _ := NewDomainName("example.com")
	s.Equal("example.com", domain.String())
}

func (s *DomainNameTestSuite) TestJSONSerialization() {
	domain, _ := NewDomainName("example.com")

	jsonData, err := json.Marshal(domain)
	s.NoError(err)
	s.JSONEq(`{}`, string(jsonData))
}

func (s *DomainNameTestSuite) TestReconstitute() {
	// Test that reconstitute bypasses validation
	domain := ReconstituteDomainName("invalid..domain")
	s.Equal("invalid..domain", domain.Value())
	s.Equal("invalid..domain", domain.String())
}

func (s *DomainNameTestSuite) TestIsValidDomainName() {
	testCases := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name:    "valid simple domain",
			input:   "example.com",
			isValid: true,
		},
		{
			name:    "valid subdomain",
			input:   "www.example.com",
			isValid: true,
		},
		{
			name:    "valid with numbers",
			input:   "test123.example.org",
			isValid: true,
		},
		{
			name:    "valid with hyphens",
			input:   "my-site.example-domain.net",
			isValid: true,
		},
		{
			name:    "invalid empty",
			input:   "",
			isValid: false,
		},
		{
			name:    "invalid consecutive dots",
			input:   "example..com",
			isValid: false,
		},
		{
			name:    "invalid starts with dot",
			input:   ".example.com",
			isValid: false,
		},
		{
			name:    "invalid ends with dot",
			input:   "example.com.",
			isValid: false,
		},
		{
			name:    "invalid starts with hyphen",
			input:   "-example.com",
			isValid: false,
		},
		{
			name:    "invalid ends with hyphen",
			input:   "example-.com",
			isValid: false,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				err := IsValidDomainName(tc.input)
				if tc.isValid {
					s.NoError(err)
				} else {
					s.Error(err)
				}
			},
		)
	}
}
