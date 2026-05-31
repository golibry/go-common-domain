package identifier

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type EANTestSuite struct {
	suite.Suite
}

func TestEANSuite(t *testing.T) {
	suite.Run(t, new(EANTestSuite))
}

func (s *EANTestSuite) TestItCanBuildEANWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"EAN-8", "96385074", "96385074"},
		{"EAN-13", "4006 3813 3393 1", "4006381333931"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ean, err := NewEAN(tc.input)

			s.NoError(err)
			s.Equal(tc.expected, ean.Value())
			s.Equal(tc.expected, ean.String())
		})
	}
}

func (s *EANTestSuite) TestItFailsToBuildEANWithInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{"empty", "", ErrEmptyEAN},
		{"invalid length", "012345678905", ErrInvalidEAN},
		{"invalid chars", "400638133393A", ErrInvalidEAN},
		{"invalid check digit", "4006381333932", ErrInvalidEAN},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewEAN(tc.input)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *EANTestSuite) TestEquals() {
	ean1, _ := NewEAN("4006381333931")
	ean2, _ := NewEAN("4006381333931")
	ean3, _ := NewEAN("96385074")

	s.True(ean1.Equals(ean2))
	s.False(ean1.Equals(ean3))
}

func (s *EANTestSuite) TestJSONSerialization() {
	ean, _ := NewEAN("4006381333931")

	jsonData, err := json.Marshal(ean)
	s.NoError(err)
	s.Equal(`"4006381333931"`, string(jsonData))

	var decoded EAN
	s.NoError(json.Unmarshal([]byte(`"4006-3813-3393-1"`), &decoded))
	s.Equal(ean.Value(), decoded.Value())
}

func (s *EANTestSuite) TestJSONSerializationFailsForNull() {
	var decoded EAN
	err := json.Unmarshal([]byte(`null`), &decoded)

	s.Error(err)
	s.True(errors.Is(err, domain.ErrNullValue))
}

func (s *EANTestSuite) TestReconstitute() {
	ean := ReconstituteEAN("4006381333931")

	s.Equal("4006381333931", ean.Value())
}
