package web

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EmailTestSuite struct {
	suite.Suite
}

func TestEmailSuite(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}

func (s *EmailTestSuite) TestItCanBuildNewEmailWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple email",
			input:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "email with uppercase",
			input:    "Test@Example.COM",
			expected: "test@example.com",
		},
		{
			name:     "email with spaces",
			input:    "  test@example.com  ",
			expected: "test@example.com",
		},
		{
			name:     "email with special characters in local part",
			input:    "test.email+tag@example.com",
			expected: "test.email+tag@example.com",
		},
		{
			name:     "email with numbers",
			input:    "user123@domain123.com",
			expected: "user123@domain123.com",
		},
		{
			name:     "email with subdomain",
			input:    "user@mail.example.com",
			expected: "user@mail.example.com",
		},
		{
			name:     "email with all allowed special chars",
			input:    "test.email!#$%&'*+-/=?^_`{|}~@example.com",
			expected: "test.email!#$%&'*+-/=?^_`{|}~@example.com",
		},
		{
			name:     "minimum length email",
			input:    "a@b.c",
			expected: "a@b.c",
		},
		{
			name:     "long but valid email",
			input:    "verylongusernamethatisvalidbutquitelengthy@verylongdomainnamethatisvalidbutquitelengthy.com",
			expected: "verylongusernamethatisvalidbutquitelengthy@verylongdomainnamethatisvalidbutquitelengthy.com",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			email, err := NewEmail(tc.input)
			s.NoError(err)
			s.Equal(tc.expected, email.Value())
		})
	}
}

func (s *EmailTestSuite) TestItFailsToBuildNewEmailFromInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "empty email",
			input:         "",
			expectedError: ErrEmptyEmail,
		},
		{
			name:          "email without @",
			input:         "testexample.com",
			expectedError: ErrMissingAtSymbol,
		},
		{
			name:          "email with multiple @",
			input:         "test@example@com",
			expectedError: ErrMultipleAtSymbols,
		},
		{
			name:          "email with empty local part",
			input:         "@example.com",
			expectedError: ErrEmptyLocalPart,
		},
		{
			name:          "email with empty domain part",
			input:         "test@",
			expectedError: ErrEmptyDomainPart,
		},
		{
			name:          "email with invalid characters",
			input:         "test@exam ple.com",
			expectedError: ErrInvalidDomainPart,
		},
		{
			name:          "email with local part starting with dot",
			input:         ".test@example.com",
			expectedError: ErrInvalidLocalPart,
		},
		{
			name:          "email with local part ending with dot",
			input:         "test.@example.com",
			expectedError: ErrInvalidLocalPart,
		},
		{
			name:          "email with consecutive dots in local part",
			input:         "te..st@example.com",
			expectedError: ErrInvalidLocalPart,
		},
		{
			name:          "email with invalid domain format",
			input:         "test@.example.com",
			expectedError: ErrInvalidDomainPart,
		},
		{
			name:          "email with domain ending with dot",
			input:         "test@example.com.",
			expectedError: ErrInvalidDomainPart,
		},
		{
			name:          "too long email",
			input:         strings.Repeat("a", 250) + "@example.com",
			expectedError: ErrTooLongEmail,
		},
		{
			name:          "too long local part",
			input:         strings.Repeat("a", 65) + "@example.com",
			expectedError: ErrTooLongLocalPart,
		},
		{
			name:          "too short email",
			input:         "a@",
			expectedError: ErrEmptyDomainPart,
		},
		{
			name:          "email with invalid local part characters",
			input:         "test@#@example.com",
			expectedError: ErrMultipleAtSymbols,
		},
		{
			name:          "email with spaces in local part",
			input:         "te st@example.com",
			expectedError: ErrInvalidEmailChars,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			email, err := NewEmail(tc.input)
			s.Error(err)
			s.True(errors.Is(err, tc.expectedError), "Expected error %v, got %v", tc.expectedError, err)
			s.Equal(Email{}, email)
		})
	}
}

func (s *EmailTestSuite) TestEmailNormalization() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase to lowercase",
			input:    "TEST@EXAMPLE.COM",
			expected: "test@example.com",
		},
		{
			name:     "mixed case to lowercase",
			input:    "TeSt@ExAmPlE.CoM",
			expected: "test@example.com",
		},
		{
			name:     "trim leading spaces",
			input:    "   test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "trim trailing spaces",
			input:    "test@example.com   ",
			expected: "test@example.com",
		},
		{
			name:     "trim both leading and trailing spaces",
			input:    "   test@example.com   ",
			expected: "test@example.com",
		},
		{
			name:     "normalize complex email",
			input:    "  TeSt.EmAiL+TaG@ExAmPlE.CoM  ",
			expected: "test.email+tag@example.com",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			normalized, err := NormalizeEmail(tc.input)
			s.NoError(err)
			s.Equal(tc.expected, normalized)
		})
	}
}

func (s *EmailTestSuite) TestEquals() {
	email1, _ := NewEmail("test@example.com")
	email2, _ := NewEmail("test@example.com")
	email3, _ := NewEmail("other@example.com")

	s.True(email1.Equals(email2))
	s.False(email1.Equals(email3))
}

func (s *EmailTestSuite) TestString() {
	email, _ := NewEmail("test@example.com")
	s.Equal("test@example.com", email.String())
}

func (s *EmailTestSuite) TestLocalPart() {
	testCases := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "simple email",
			email:    "test@example.com",
			expected: "test",
		},
		{
			name:     "email with special characters",
			email:    "test.email+tag@example.com",
			expected: "test.email+tag",
		},
		{
			name:     "email with numbers",
			email:    "user123@domain.com",
			expected: "user123",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			email, err := NewEmail(tc.email)
			s.NoError(err)
			s.Equal(tc.expected, email.LocalPart())
		})
	}
}

func (s *EmailTestSuite) TestDomainPart() {
	testCases := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "simple email",
			email:    "test@example.com",
			expected: "example.com",
		},
		{
			name:     "email with subdomain",
			email:    "test@mail.example.com",
			expected: "mail.example.com",
		},
		{
			name:     "email with numbers in domain",
			email:    "test@domain123.com",
			expected: "domain123.com",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			email, err := NewEmail(tc.email)
			s.NoError(err)
			s.Equal(tc.expected, email.DomainPart())
		})
	}
}

func (s *EmailTestSuite) TestJSONSerialization() {
	email, _ := NewEmail("test@example.com")
	jsonData, err := json.Marshal(email)
	s.NoError(err)
	s.JSONEq(`{"value":"test@example.com"}`, string(jsonData))
}

func (s *EmailTestSuite) TestReconstitute() {
	email := ReconstituteEmail("test@example.com")
	s.Equal("test@example.com", email.Value())
}

func (s *EmailTestSuite) TestItCanBuildNewEmailFromValidJSON() {
	jsonData := `{"value":"test@example.com"}`
	email, err := NewEmailFromJSON([]byte(jsonData))
	s.NoError(err)
	s.Equal("test@example.com", email.Value())
}

func (s *EmailTestSuite) TestItFailsToBuildNewEmailFromInvalidJSON() {
	testCases := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON format",
			jsonData: `{"value":"test@example.com"`,
		},
		{
			name:     "invalid email in JSON",
			jsonData: `{"value":"invalid-email"}`,
		},
		{
			name:     "empty JSON",
			jsonData: `{}`,
		},
		{
			name:     "null value",
			jsonData: `{"value":null}`,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			email, err := NewEmailFromJSON([]byte(tc.jsonData))
			s.Error(err)
			s.Equal(Email{}, email)
		})
	}
}

func (s *EmailTestSuite) TestIsValidEmail() {
	testCases := []struct {
		name     string
		email    string
		isValid  bool
		expected error
	}{
		{
			name:     "valid simple email",
			email:    "test@example.com",
			isValid:  true,
			expected: nil,
		},
		{
			name:     "valid email with special characters",
			email:    "test.email+tag@example.com",
			isValid:  true,
			expected: nil,
		},
		{
			name:     "empty email",
			email:    "",
			isValid:  false,
			expected: ErrEmptyEmail,
		},
		{
			name:     "missing @ symbol",
			email:    "testexample.com",
			isValid:  false,
			expected: ErrMissingAtSymbol,
		},
		{
			name:     "multiple @ symbols",
			email:    "test@example@com",
			isValid:  false,
			expected: ErrMultipleAtSymbols,
		},
		{
			name:     "empty local part",
			email:    "@example.com",
			isValid:  false,
			expected: ErrEmptyLocalPart,
		},
		{
			name:     "empty domain part",
			email:    "test@",
			isValid:  false,
			expected: ErrEmptyDomainPart,
		},
		{
			name:     "too long email",
			email:    strings.Repeat("a", 250) + "@example.com",
			isValid:  false,
			expected: ErrTooLongEmail,
		},
		{
			name:     "too long local part",
			email:    strings.Repeat("a", 65) + "@example.com",
			isValid:  false,
			expected: ErrTooLongLocalPart,
		},
		{
			name:     "invalid local part format",
			email:    ".test@example.com",
			isValid:  false,
			expected: ErrInvalidLocalPart,
		},
		{
			name:     "invalid domain part format",
			email:    "test@.example.com",
			isValid:  false,
			expected: ErrInvalidDomainPart,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := IsValidEmail(tc.email)
			if tc.isValid {
				s.NoError(err)
			} else {
				s.Error(err)
				s.True(errors.Is(err, tc.expected), "Expected error %v, got %v", tc.expected, err)
			}
		})
	}
}