package identifier

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type CNKTestSuite struct {
	suite.Suite
}

func TestCNKSuite(t *testing.T) {
	suite.Run(t, new(CNKTestSuite))
}

func (s *CNKTestSuite) TestItCanBuildCNKWithValidValue() {
	cnk, err := NewCNK("123-4566")

	s.NoError(err)
	s.Equal("1234566", cnk.Value())
	s.Equal("1234566", cnk.String())
}

func (s *CNKTestSuite) TestItFailsToBuildCNKWithInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{"empty", "", ErrEmptyCNK},
		{"invalid length", "123456", ErrInvalidCNK},
		{"invalid chars", "123456A", ErrInvalidCNK},
		{"invalid check digit", "1234567", ErrInvalidCNK},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewCNK(tc.input)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *CNKTestSuite) TestEquals() {
	cnk1, _ := NewCNK("1234566")
	cnk2, _ := NewCNK("1234566")
	cnk3, _ := NewCNK("7654320")

	s.True(cnk1.Equals(cnk2))
	s.False(cnk1.Equals(cnk3))
}

func (s *CNKTestSuite) TestJSONSerialization() {
	cnk, _ := NewCNK("1234566")

	jsonData, err := json.Marshal(cnk)
	s.NoError(err)
	s.Equal(`"1234566"`, string(jsonData))

	var decoded CNK
	s.NoError(json.Unmarshal([]byte(`"123-4566"`), &decoded))
	s.Equal(cnk.Value(), decoded.Value())
}

func (s *CNKTestSuite) TestJSONSerializationFailsForNull() {
	var decoded CNK
	err := json.Unmarshal([]byte(`null`), &decoded)

	s.Error(err)
	s.True(errors.Is(err, domain.ErrNullValue))
}

func (s *CNKTestSuite) TestReconstitute() {
	cnk := ReconstituteCNK("1234566")

	s.Equal("1234566", cnk.Value())
}
