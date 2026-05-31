package geography

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type AddressTestSuite struct {
	suite.Suite
}

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

func (s *AddressTestSuite) TestItCanBuildAddressWithValidValues() {
	country, _ := NewCountryCode("ro")

	address, err := NewAddress("  Main   Street 1 ", " Apt 2 ", " Bucharest ", " Bucuresti ", " 010101 ", country)

	s.NoError(err)
	s.Equal("Main Street 1", address.Line1())
	s.Equal("Apt 2", address.Line2())
	s.Equal("Bucharest", address.City())
	s.Equal("Bucuresti", address.Region())
	s.Equal("010101", address.PostalCode())
	s.Equal("RO", address.Country().String())
	s.Equal("Main Street 1, Apt 2, Bucharest, Bucuresti, 010101, RO", address.String())
}

func (s *AddressTestSuite) TestItFailsToBuildAddressWithInvalidValues() {
	country, _ := NewCountryCode("RO")
	longValue := strings.Repeat("a", MaxAddressLineLength+1)

	testCases := []struct {
		name          string
		line1         string
		city          string
		country       CountryCode
		expectedError error
	}{
		{"missing line 1", "", "Bucharest", country, ErrMissingAddressLine1},
		{"missing city", "Main Street 1", "", country, ErrMissingAddressCity},
		{"missing country", "Main Street 1", "Bucharest", CountryCode{}, ErrMissingAddressCountry},
		{"too long line", longValue, "Bucharest", country, ErrTooLongAddressLine},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewAddress(tc.line1, "", tc.city, "", "", tc.country)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *AddressTestSuite) TestEquals() {
	country, _ := NewCountryCode("RO")
	address1, _ := NewAddress("Main Street 1", "", "Bucharest", "", "", country)
	address2, _ := NewAddress("Main Street 1", "", "Bucharest", "", "", country)
	address3, _ := NewAddress("Other Street 1", "", "Bucharest", "", "", country)

	s.True(address1.Equals(address2))
	s.False(address1.Equals(address3))
}

func (s *AddressTestSuite) TestJSONSerialization() {
	country, _ := NewCountryCode("RO")
	address, _ := NewAddress("Main Street 1", "", "Bucharest", "", "010101", country)

	jsonData, err := json.Marshal(address)
	s.NoError(err)
	s.JSONEq(`{"line1":"Main Street 1","city":"Bucharest","postalCode":"010101","country":"RO"}`, string(jsonData))

	var decoded Address
	s.NoError(json.Unmarshal([]byte(`{"line1":"Main Street 1","city":"Bucharest","country":"ro"}`), &decoded))
	s.Equal("RO", decoded.Country().String())
	s.Equal("Main Street 1, Bucharest, RO", decoded.String())
}

func (s *AddressTestSuite) TestJSONSerializationFailsForNull() {
	var decoded Address
	err := json.Unmarshal([]byte(`null`), &decoded)

	s.Error(err)
	s.True(errors.Is(err, domain.ErrNullValue))
}

func (s *AddressTestSuite) TestReconstitute() {
	country := ReconstituteCountryCode("RO")
	address := ReconstituteAddress("Main Street 1", "", "Bucharest", "", "010101", country)

	s.Equal("Main Street 1", address.Line1())
	s.Equal("RO", address.Country().String())
}
