package geography

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CountryCodeTestSuite struct {
	suite.Suite
}

func TestCountryCodeSuite(t *testing.T) {
	suite.Run(t, new(CountryCodeTestSuite))
}

func (s *CountryCodeTestSuite) TestItCanBuildNewCountryCodeWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "US country code",
			input:    "US",
			expected: "US",
		},
		{
			name:     "lowercase country code",
			input:    "us",
			expected: "US",
		},
		{
			name:     "mixed case country code",
			input:    "Us",
			expected: "US",
		},
		{
			name:     "country code with spaces",
			input:    " US ",
			expected: "US",
		},
		{
			name:     "UK country code",
			input:    "GB",
			expected: "GB",
		},
		{
			name:     "Canada country code",
			input:    "ca",
			expected: "CA",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				countryCode, err := NewCountryCode(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, countryCode.Value())
				s.Equal(tc.expected, countryCode.String())
			},
		)
	}
}

func (s *CountryCodeTestSuite) TestItFailsToBuildNewCountryCodeFromInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "empty country code",
			input:         "",
			expectedError: ErrEmptyCountryCode,
		},
		{
			name:          "country code with only spaces",
			input:         "   ",
			expectedError: ErrEmptyCountryCode,
		},
		{
			name:          "country code too short",
			input:         "U",
			expectedError: ErrInvalidCountryCode,
		},
		{
			name:          "country code too long",
			input:         "USA",
			expectedError: ErrInvalidCountryCode,
		},
		{
			name:          "country code with numbers",
			input:         "U1",
			expectedError: ErrInvalidCountryCode,
		},
		{
			name:          "country code with special characters",
			input:         "U-",
			expectedError: ErrInvalidCountryCode,
		},
		{
			name:          "country code with spaces in middle",
			input:         "U S",
			expectedError: ErrInvalidCountryCode,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewCountryCode(tc.input)
				s.Error(err)
				s.True(errors.Is(err, tc.expectedError))
			},
		)
	}
}

func (s *CountryCodeTestSuite) TestCountryCodeNormalization() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "converts to uppercase",
			input:    "us",
			expected: "US",
		},
		{
			name:     "trims whitespace",
			input:    "  gb  ",
			expected: "GB",
		},
		{
			name:     "handles mixed case",
			input:    "cA",
			expected: "CA",
		},
		{
			name:     "already uppercase",
			input:    "FR",
			expected: "FR",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				normalized, err := NormalizeCountryCode(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, normalized)
			},
		)
	}
}

func (s *CountryCodeTestSuite) TestEquals() {
	countryCode1, _ := NewCountryCode("US")
	countryCode2, _ := NewCountryCode("us")
	countryCode3, _ := NewCountryCode("GB")

	s.True(countryCode1.Equals(countryCode2))
	s.False(countryCode1.Equals(countryCode3))
}

func (s *CountryCodeTestSuite) TestString() {
	countryCode, _ := NewCountryCode("us")
	s.Equal("US", countryCode.String())
}

func (s *CountryCodeTestSuite) TestJSONSerialization() {
	countryCode, _ := NewCountryCode("US")

	jsonData, err := json.Marshal(countryCode)
	s.NoError(err)
	s.JSONEq(`{}`, string(jsonData))
}

func (s *CountryCodeTestSuite) TestReconstitute() {
	countryCode := ReconstituteCountryCode("US")
	s.Equal("US", countryCode.Value())
	s.Equal("US", countryCode.String())
}
