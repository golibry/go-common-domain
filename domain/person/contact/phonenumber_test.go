package contact

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PhoneNumberTestSuite struct {
	suite.Suite
}

func TestPhoneNumberSuite(t *testing.T) {
	suite.Run(t, new(PhoneNumberTestSuite))
}

func (s *PhoneNumberTestSuite) TestItCanBuildNewPhoneNumberWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "US phone number with country code",
			input:    "+1234567890",
			expected: "+1234567890",
		},
		{
			name:     "Phone number with spaces and dashes",
			input:    "+1 (234) 567-890",
			expected: "+1234567890",
		},
		{
			name:     "Phone number with dots",
			input:    "+1.234.567.890",
			expected: "+1234567890",
		},
		{
			name:     "International phone number",
			input:    "+44123456789",
			expected: "+44123456789",
		},
		{
			name:     "Phone number without country code",
			input:    "1234567890",
			expected: "1234567890",
		},
		{
			name:     "Short phone number",
			input:    "123",
			expected: "123",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				phoneNumber, err := NewPhoneNumber(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, phoneNumber.Value())
				s.Equal(tc.expected, phoneNumber.String())
			},
		)
	}
}

func (s *PhoneNumberTestSuite) TestItFailsToBuildNewPhoneNumberFromInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "empty phone number",
			input:         "",
			expectedError: ErrEmptyPhoneNumber,
		},
		{
			name:          "phone number with only spaces",
			input:         "   ",
			expectedError: ErrEmptyPhoneNumber,
		},
		{
			name:          "phone number with letters",
			input:         "+1234abc567",
			expectedError: ErrInvalidPhoneNumberChars,
		},
		{
			name:          "phone number starting with zero",
			input:         "+0123456789",
			expectedError: ErrInvalidPhoneNumberChars,
		},
		{
			name:          "phone number too short",
			input:         "12",
			expectedError: ErrTooShortPhoneNumber,
		},
		{
			name:          "phone number too long",
			input:         "+123456789012345678901",
			expectedError: ErrTooLongPhoneNumber,
		},
		{
			name:          "phone number with multiple plus signs",
			input:         "++1234567890",
			expectedError: ErrInvalidPhoneNumberChars,
		},
		{
			name:          "phone number with plus in middle",
			input:         "123+4567890",
			expectedError: ErrInvalidPhoneNumberChars,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewPhoneNumber(tc.input)
				s.Error(err)
				s.True(errors.Is(err, tc.expectedError))
			},
		)
	}
}

func (s *PhoneNumberTestSuite) TestPhoneNumberNormalization() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes spaces",
			input:    "+1 234 567 890",
			expected: "+1234567890",
		},
		{
			name:     "removes dashes",
			input:    "+1-234-567-890",
			expected: "+1234567890",
		},
		{
			name:     "removes parentheses",
			input:    "+1(234)567890",
			expected: "+1234567890",
		},
		{
			name:     "removes dots",
			input:    "+1.234.567.890",
			expected: "+1234567890",
		},
		{
			name:     "removes mixed formatting",
			input:    "+1 (234) 567-890",
			expected: "+1234567890",
		},
		{
			name:     "trims whitespace",
			input:    "  +1234567890  ",
			expected: "+1234567890",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				normalized, err := NormalizePhoneNumber(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, normalized)
			},
		)
	}
}

func (s *PhoneNumberTestSuite) TestEquals() {
	phoneNumber1, _ := NewPhoneNumber("+1234567890")
	phoneNumber2, _ := NewPhoneNumber("+1234567890")
	phoneNumber3, _ := NewPhoneNumber("+9876543210")

	s.True(phoneNumber1.Equals(phoneNumber2))
	s.False(phoneNumber1.Equals(phoneNumber3))
}

func (s *PhoneNumberTestSuite) TestString() {
	phoneNumber, _ := NewPhoneNumber("+1234567890")
	s.Equal("+1234567890", phoneNumber.String())
}

func (s *PhoneNumberTestSuite) TestJSONSerialization() {
	phoneNumber, _ := NewPhoneNumber("+1234567890")

	jsonData, err := json.Marshal(phoneNumber)
	s.NoError(err)
	s.JSONEq(`{}`, string(jsonData))
}

func (s *PhoneNumberTestSuite) TestReconstitute() {
	phoneNumber := ReconstitutePhoneNumber("+1234567890")
	s.Equal("+1234567890", phoneNumber.Value())
	s.Equal("+1234567890", phoneNumber.String())
}
