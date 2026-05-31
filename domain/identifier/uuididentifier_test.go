package identifier

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type UUIDIdentifierTestSuite struct {
	suite.Suite
}

func TestUUIDIdentifierSuite(t *testing.T) {
	suite.Run(t, new(UUIDIdentifierTestSuite))
}

func (s *UUIDIdentifierTestSuite) TestItCanBuildUUIDIdentifierWithValidValue() {
	identifier, err := NewUUIDIdentifier(" 550E8400-E29B-41D4-A716-446655440000 ")

	s.NoError(err)
	s.Equal("550e8400-e29b-41d4-a716-446655440000", identifier.Value())
	s.Equal("550e8400-e29b-41d4-a716-446655440000", identifier.String())
}

func (s *UUIDIdentifierTestSuite) TestItFailsToBuildUUIDIdentifierWithInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{"empty", "", ErrEmptyUUIDIdentifier},
		{"invalid format", "550e8400e29b41d4a716446655440000", ErrInvalidUUIDIdentifier},
		{"invalid characters", "550e8400-e29b-41d4-a716-44665544000z", ErrInvalidUUIDIdentifier},
		{"zero uuid", zeroUUIDIdentifier, ErrZeroUUIDIdentifier},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewUUIDIdentifier(tc.input)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *UUIDIdentifierTestSuite) TestEquals() {
	identifier1, _ := NewUUIDIdentifier("550e8400-e29b-41d4-a716-446655440000")
	identifier2, _ := NewUUIDIdentifier("550e8400-e29b-41d4-a716-446655440000")
	identifier3, _ := NewUUIDIdentifier("7c9e6679-7425-40de-944b-e07fc1f90ae7")

	s.True(identifier1.Equals(identifier2))
	s.False(identifier1.Equals(identifier3))
}

func (s *UUIDIdentifierTestSuite) TestJSONSerialization() {
	identifier, _ := NewUUIDIdentifier("550e8400-e29b-41d4-a716-446655440000")

	jsonData, err := json.Marshal(identifier)
	s.NoError(err)
	s.Equal(`"550e8400-e29b-41d4-a716-446655440000"`, string(jsonData))

	var decoded UUIDIdentifier
	s.NoError(json.Unmarshal([]byte(`"550E8400-E29B-41D4-A716-446655440000"`), &decoded))
	s.Equal("550e8400-e29b-41d4-a716-446655440000", decoded.Value())
}

func (s *UUIDIdentifierTestSuite) TestJSONSerializationFailsForNull() {
	var decoded UUIDIdentifier
	err := json.Unmarshal([]byte(`null`), &decoded)

	s.Error(err)
	s.True(errors.Is(err, domain.ErrNullValue))
}

func (s *UUIDIdentifierTestSuite) TestTextSerialization() {
	identifier, _ := NewUUIDIdentifier("550e8400-e29b-41d4-a716-446655440000")

	text, err := identifier.MarshalText()
	s.NoError(err)
	s.Equal("550e8400-e29b-41d4-a716-446655440000", string(text))

	var decoded UUIDIdentifier
	s.NoError(decoded.UnmarshalText([]byte("550E8400-E29B-41D4-A716-446655440000")))
	s.Equal(identifier.Value(), decoded.Value())
}

func (s *UUIDIdentifierTestSuite) TestReconstitute() {
	identifier := ReconstituteUUIDIdentifier("550e8400-e29b-41d4-a716-446655440000")

	s.Equal("550e8400-e29b-41d4-a716-446655440000", identifier.Value())
}
