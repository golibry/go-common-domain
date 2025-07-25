package finance

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CurrencyTestSuite struct {
	suite.Suite
}

func TestCurrencySuite(t *testing.T) {
	suite.Run(t, new(CurrencyTestSuite))
}

func (s *CurrencyTestSuite) TestItCanBuildNewCurrencyWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "USD currency",
			input:    "USD",
			expected: "USD",
		},
		{
			name:     "lowercase currency",
			input:    "usd",
			expected: "USD",
		},
		{
			name:     "mixed case currency",
			input:    "UsD",
			expected: "USD",
		},
		{
			name:     "currency with spaces",
			input:    " EUR ",
			expected: "EUR",
		},
		{
			name:     "GBP currency",
			input:    "GBP",
			expected: "GBP",
		},
		{
			name:     "JPY currency",
			input:    "jpy",
			expected: "JPY",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			currency, err := NewCurrency(tc.input)
			s.NoError(err)
			s.Equal(tc.expected, currency.Value())
			s.Equal(tc.expected, currency.String())
		})
	}
}

func (s *CurrencyTestSuite) TestItFailsToBuildNewCurrencyFromInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "empty currency",
			input:         "",
			expectedError: ErrEmptyCurrency,
		},
		{
			name:          "currency with only spaces",
			input:         "   ",
			expectedError: ErrEmptyCurrency,
		},
		{
			name:          "currency too short",
			input:         "US",
			expectedError: ErrInvalidCurrency,
		},
		{
			name:          "currency too long",
			input:         "USDD",
			expectedError: ErrInvalidCurrency,
		},
		{
			name:          "currency with numbers",
			input:         "US1",
			expectedError: ErrInvalidCurrency,
		},
		{
			name:          "currency with special characters",
			input:         "US-",
			expectedError: ErrInvalidCurrency,
		},
		{
			name:          "currency with spaces in middle",
			input:         "U S D",
			expectedError: ErrInvalidCurrency,
		},
		{
			name:          "single character",
			input:         "U",
			expectedError: ErrInvalidCurrency,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewCurrency(tc.input)
			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *CurrencyTestSuite) TestCurrencyNormalization() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "converts to uppercase",
			input:    "usd",
			expected: "USD",
		},
		{
			name:     "trims whitespace",
			input:    "  eur  ",
			expected: "EUR",
		},
		{
			name:     "handles mixed case",
			input:    "gBp",
			expected: "GBP",
		},
		{
			name:     "already uppercase",
			input:    "JPY",
			expected: "JPY",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			normalized, err := NormalizeCurrency(tc.input)
			s.NoError(err)
			s.Equal(tc.expected, normalized)
		})
	}
}

func (s *CurrencyTestSuite) TestEquals() {
	currency1, _ := NewCurrency("USD")
	currency2, _ := NewCurrency("usd")
	currency3, _ := NewCurrency("EUR")

	s.True(currency1.Equals(currency2))
	s.False(currency1.Equals(currency3))
}

func (s *CurrencyTestSuite) TestString() {
	currency, _ := NewCurrency("usd")
	s.Equal("USD", currency.String())
}

func (s *CurrencyTestSuite) TestJSONSerialization() {
	currency, _ := NewCurrency("USD")
	
	jsonData, err := json.Marshal(currency)
	s.NoError(err)
	s.JSONEq(`{"value":"USD"}`, string(jsonData))
}

func (s *CurrencyTestSuite) TestReconstitute() {
	currency := ReconstituteCurrency("USD")
	s.Equal("USD", currency.Value())
	s.Equal("USD", currency.String())
}

func (s *CurrencyTestSuite) TestItCanBuildNewCurrencyFromValidJSON() {
	jsonData := `{"value":"USD"}`
	
	currency, err := NewCurrencyFromJSON([]byte(jsonData))
	s.NoError(err)
	s.Equal("USD", currency.Value())
}

func (s *CurrencyTestSuite) TestItFailsToBuildNewCurrencyFromInvalidJSON() {
	testCases := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON format",
			jsonData: `{"value":"USD"`,
		},
		{
			name:     "invalid currency in JSON",
			jsonData: `{"value":"USDD"}`,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewCurrencyFromJSON([]byte(tc.jsonData))
			s.Error(err)
		})
	}
}