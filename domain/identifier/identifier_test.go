package identifier

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IdentifierTestSuite struct {
	suite.Suite
}

func TestIdentifierSuite(t *testing.T) {
	suite.Run(t, new(IdentifierTestSuite))
}

func (s *IdentifierTestSuite) TestItCanBuildNewIdentifierWithValidValues() {
	testCases := []struct {
		name     string
		input    uint64
		expected uint64
	}{
		{
			name:     "small positive number",
			input:    1,
			expected: 1,
		},
		{
			name:     "medium positive number",
			input:    12345,
			expected: 12345,
		},
		{
			name:     "large positive number",
			input:    18446744073709551615, // max uint64
			expected: 18446744073709551615,
		},
		{
			name:     "typical ID",
			input:    999999,
			expected: 999999,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				identifier, err := NewIdentifier(tc.input)
				s.NoError(err)
    s.Equal(tc.expected, identifier.Value())
    s.EqualValues(tc.expected, identifier.Value())
    s.Equal(tc.expected, identifier.Value())
			},
		)
	}
}

func (s *IdentifierTestSuite) TestItFailsToBuildNewIdentifierFromInvalidValues() {
	testCases := []struct {
		name          string
		input         uint64
		expectedError error
	}{
		{
			name:          "zero identifier",
			input:         0,
			expectedError: ErrZeroIdentifier,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewIdentifier(tc.input)
				s.Error(err)
				s.True(errors.Is(err, tc.expectedError))
			},
		)
	}
}

func (s *IdentifierTestSuite) TestItCanBuildNewIdentifierFromInt() {
	testCases := []struct {
		name          string
		input         int64
		expected      uint64
		expectedError error
	}{
		{
			name:     "positive int",
			input:    123,
			expected: 123,
		},
		{
			name:          "negative int",
			input:         -123,
			expectedError: ErrNegativeIdentifier,
		},
		{
			name:          "zero int",
			input:         0,
			expectedError: ErrZeroIdentifier,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				identifier, err := NewIdentifierFromInt(tc.input)
				if tc.expectedError != nil {
					s.Error(err)
					s.True(errors.Is(err, tc.expectedError))
				} else {
					s.NoError(err)
					s.Equal(tc.expected, identifier.Value())
				}
			},
		)
	}
}

func (s *IdentifierTestSuite) TestItCanBuildNewIdentifierFromString() {
	testCases := []struct {
		name          string
		input         string
		expected      uint64
		expectedError error
	}{
		{
			name:     "valid string number",
			input:    "123",
			expected: 123,
		},
		{
			name:     "large string number",
			input:    "999999999999",
			expected: 999999999999,
		},
		{
			name:          "zero string",
			input:         "0",
			expectedError: ErrZeroIdentifier,
		},
		{
			name:          "invalid string",
			input:         "abc",
			expectedError: ErrInvalidIdentifier,
		},
		{
			name:          "negative string",
			input:         "-123",
			expectedError: ErrInvalidIdentifier,
		},
		{
			name:          "empty string",
			input:         "",
			expectedError: ErrInvalidIdentifier,
		},
		{
			name:          "string with spaces",
			input:         " 123 ",
			expectedError: ErrInvalidIdentifier,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				identifier, err := NewIdentifierFromString(tc.input)
				if tc.expectedError != nil {
					s.Error(err)
					s.True(errors.Is(err, tc.expectedError))
				} else {
					s.NoError(err)
					s.Equal(tc.expected, identifier.Value())
				}
			},
		)
	}
}

func (s *IdentifierTestSuite) TestEquals() {
	identifier1, _ := NewIdentifier(123)
	identifier2, _ := NewIdentifier(123)
	identifier3, _ := NewIdentifier(456)

	s.True(identifier1.Equals(identifier2))
	s.False(identifier1.Equals(identifier3))
}

func (s *IdentifierTestSuite) TestString() {
	identifier, _ := NewIdentifier(12345)
	s.Equal("12345", identifier.String())
}

func (s *IdentifierTestSuite) TestJSONSerialization() {
	identifier, _ := NewIdentifier(12345)

	jsonData, err := json.Marshal(identifier)
	s.NoError(err)
	s.JSONEq(`{"value":12345}`, string(jsonData))
}

func (s *IdentifierTestSuite) TestReconstitute() {
	identifier := ReconstituteIdentifier(12345)
	s.Equal(uint64(12345), identifier.Value())
	s.Equal("12345", identifier.String())
}

func (s *IdentifierTestSuite) TestItCanBuildNewIdentifierFromValidJSON() {
	jsonData := `{"value":12345}`

	identifier, err := NewIdentifierFromJSON([]byte(jsonData))
	s.NoError(err)
	s.Equal(uint64(12345), identifier.Value())
}

func (s *IdentifierTestSuite) TestItFailsToBuildNewIdentifierFromInvalidJSON() {
	testCases := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON format",
			jsonData: `{"value":12345`,
		},
		{
			name:     "zero identifier in JSON",
			jsonData: `{"value":0}`,
		},
		{
			name:     "string value in JSON",
			jsonData: `{"value":"123"}`,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewIdentifierFromJSON([]byte(tc.jsonData))
				s.Error(err)
			},
		)
	}
}
