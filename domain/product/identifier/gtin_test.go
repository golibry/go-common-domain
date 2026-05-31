package identifier

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type GTINTestSuite struct {
	suite.Suite
}

func TestGTINSuite(t *testing.T) {
	suite.Run(t, new(GTINTestSuite))
}

func (s *GTINTestSuite) TestItCanBuildGTINWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"GTIN-8", "96385074", "96385074"},
		{"GTIN-12", "012345678905", "012345678905"},
		{"GTIN-13", "4006 3813 3393 1", "4006381333931"},
		{"GTIN-14", "10012345678902", "10012345678902"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			gtin, err := NewGTIN(tc.input)

			s.NoError(err)
			s.Equal(tc.expected, gtin.Value())
			s.Equal(tc.expected, gtin.String())
		})
	}
}

func (s *GTINTestSuite) TestItFailsToBuildGTINWithInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{"empty", "", ErrEmptyGTIN},
		{"invalid length", "1234567", ErrInvalidGTIN},
		{"invalid chars", "400638133393A", ErrInvalidGTIN},
		{"invalid check digit", "4006381333932", ErrInvalidGTIN},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewGTIN(tc.input)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *GTINTestSuite) TestEquals() {
	gtin1, _ := NewGTIN("4006381333931")
	gtin2, _ := NewGTIN("4006381333931")
	gtin3, _ := NewGTIN("96385074")

	s.True(gtin1.Equals(gtin2))
	s.False(gtin1.Equals(gtin3))
}

func (s *GTINTestSuite) TestJSONSerialization() {
	gtin, _ := NewGTIN("4006381333931")

	jsonData, err := json.Marshal(gtin)
	s.NoError(err)
	s.Equal(`"4006381333931"`, string(jsonData))

	var decoded GTIN
	s.NoError(json.Unmarshal([]byte(`"4006-3813-3393-1"`), &decoded))
	s.Equal(gtin.Value(), decoded.Value())
}

func (s *GTINTestSuite) TestJSONSerializationFailsForNull() {
	var decoded GTIN
	err := json.Unmarshal([]byte(`null`), &decoded)

	s.Error(err)
	s.True(errors.Is(err, domain.ErrNullValue))
}

func (s *GTINTestSuite) TestReconstitute() {
	gtin := ReconstituteGTIN("4006381333931")

	s.Equal("4006381333931", gtin.Value())
}
