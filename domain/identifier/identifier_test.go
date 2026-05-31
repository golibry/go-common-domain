package identifier

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golibry/go-common-domain/domain"
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
				identifier, err := NewIntIdentifier(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, identifier.Value())
				s.EqualValues(tc.expected, identifier.Value())
				s.Equal(tc.expected, identifier.Value())
			},
		)
	}
}

func (s *IdentifierTestSuite) TestItCanBuildNewIdentifierFromIntWithValidValues() {
	testCases := []struct {
		name     string
		input    int64
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
			name:     "max int64",
			input:    9223372036854775807,
			expected: 9223372036854775807,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				identifier, err := NewIntIdentifierFromInt(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, identifier.Value())
			},
		)
	}
}

func (s *IdentifierTestSuite) TestItFailsToBuildNewIdentifierFromInvalidValues() {
	testCases := []struct {
		name          string
		input         int64
		expectedError error
	}{
		{
			name:          "zero identifier",
			input:         0,
			expectedError: ErrZeroIdentifier,
		},
		{
			name:          "negative identifier",
			input:         -1,
			expectedError: ErrInvalidIdentifier,
		},
		{
			name:          "large negative identifier",
			input:         -9223372036854775808,
			expectedError: ErrInvalidIdentifier,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewIntIdentifierFromInt(tc.input)
				s.Error(err)
				s.True(errors.Is(err, tc.expectedError))
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
				identifier, err := NewIntIdentifierFromString(tc.input)
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
	identifier1, _ := NewIntIdentifier(123)
	identifier2, _ := NewIntIdentifier(123)
	identifier3, _ := NewIntIdentifier(456)

	s.True(identifier1.Equals(identifier2))
	s.False(identifier1.Equals(identifier3))
}

func (s *IdentifierTestSuite) TestString() {
	identifier, _ := NewIntIdentifier(12345)
	s.Equal("12345", identifier.String())
}

func (s *IdentifierTestSuite) TestJSONSerialization() {
	identifier, _ := NewIntIdentifier(12345)

	jsonData, err := json.Marshal(identifier)
	s.NoError(err)
	s.Equal(`12345`, string(jsonData))

	var decoded IntIdentifier
	s.NoError(json.Unmarshal([]byte(`456`), &decoded))
	s.Equal(uint64(456), decoded.Value())

	s.NoError(json.Unmarshal([]byte(`"789"`), &decoded))
	s.Equal(uint64(789), decoded.Value())

	s.Error(json.Unmarshal([]byte(`0`), &decoded))
	s.Error(json.Unmarshal([]byte(`-1`), &decoded))
	s.Error(json.Unmarshal([]byte(`"abc"`), &decoded))
}

func (s *IdentifierTestSuite) TestJSONSerializationFailsForNull() {
	var decoded IntIdentifier
	err := json.Unmarshal([]byte(`null`), &decoded)

	s.Error(err)
	s.True(errors.Is(err, domain.ErrNullValue))
}

func (s *IdentifierTestSuite) TestTextSerialization() {
	identifier, _ := NewIntIdentifier(12345)

	text, err := identifier.MarshalText()
	s.NoError(err)
	s.Equal("12345", string(text))

	var decoded IntIdentifier
	s.NoError(decoded.UnmarshalText([]byte("456")))
	s.Equal(uint64(456), decoded.Value())

	s.Error(decoded.UnmarshalText([]byte("0")))
	s.Error(decoded.UnmarshalText([]byte("-1")))
	s.Error(decoded.UnmarshalText([]byte("abc")))
}

func (s *IdentifierTestSuite) TestDatabaseScan() {
	var identifier IntIdentifier

	s.NoError(identifier.Scan(int64(123)))
	s.Equal(uint64(123), identifier.Value())

	s.NoError(identifier.Scan("456"))
	s.Equal(uint64(456), identifier.Value())

	s.NoError(identifier.Scan([]byte("789")))
	s.Equal(uint64(789), identifier.Value())

	s.Error(identifier.Scan(int64(-1)))
	s.Error(identifier.Scan(123))
}

func (s *IdentifierTestSuite) TestReconstitute() {
	identifier := ReconstituteIntIdentifier(12345)
	s.Equal(uint64(12345), identifier.Value())
	s.Equal("12345", identifier.String())
}
