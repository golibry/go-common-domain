package geography

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	"github.com/golibry/go-common-domain/domain"
)

const (
	MaxAddressLineLength = 128
	MaxAddressPartLength = 96
)

var (
	ErrMissingAddressLine1   = domain.NewError("address line 1 is required")
	ErrMissingAddressCity    = domain.NewError("address city is required")
	ErrMissingAddressCountry = domain.NewError(
		"address country is required",
	)
	ErrTooLongAddressLine = domain.NewError("address line is too long")
	ErrTooLongAddressPart = domain.NewError("address part is too long")
)

type Address struct {
	line1      string
	line2      string
	city       string
	region     string
	postalCode string
	country    CountryCode
}

type addressJSON struct {
	Line1      string      `json:"line1"`
	Line2      string      `json:"line2,omitempty"`
	City       string      `json:"city"`
	Region     string      `json:"region,omitempty"`
	PostalCode string      `json:"postalCode,omitempty"`
	Country    CountryCode `json:"country"`
}

// NewAddress creates a new Address with validation and normalization.
func NewAddress(line1, line2, city, region, postalCode string, country CountryCode) (Address, error) {
	normalizedLine1 := normalizeAddressPart(line1)
	if normalizedLine1 == "" {
		return Address{}, ErrMissingAddressLine1
	}

	if utf8.RuneCountInString(normalizedLine1) > MaxAddressLineLength {
		return Address{}, ErrTooLongAddressLine
	}

	normalizedLine2 := normalizeAddressPart(line2)
	if utf8.RuneCountInString(normalizedLine2) > MaxAddressLineLength {
		return Address{}, ErrTooLongAddressLine
	}

	normalizedCity := normalizeAddressPart(city)
	if normalizedCity == "" {
		return Address{}, ErrMissingAddressCity
	}

	if utf8.RuneCountInString(normalizedCity) > MaxAddressPartLength {
		return Address{}, ErrTooLongAddressPart
	}

	normalizedRegion := normalizeAddressPart(region)
	if utf8.RuneCountInString(normalizedRegion) > MaxAddressPartLength {
		return Address{}, ErrTooLongAddressPart
	}

	normalizedPostalCode := normalizeAddressPart(postalCode)
	if utf8.RuneCountInString(normalizedPostalCode) > MaxAddressPartLength {
		return Address{}, ErrTooLongAddressPart
	}

	if country.String() == "" {
		return Address{}, ErrMissingAddressCountry
	}

	return Address{
		line1:      normalizedLine1,
		line2:      normalizedLine2,
		city:       normalizedCity,
		region:     normalizedRegion,
		postalCode: normalizedPostalCode,
		country:    country,
	}, nil
}

// ReconstituteAddress creates an Address from trusted persisted fields.
func ReconstituteAddress(line1, line2, city, region, postalCode string, country CountryCode) Address {
	return Address{
		line1:      line1,
		line2:      line2,
		city:       city,
		region:     region,
		postalCode: postalCode,
		country:    country,
	}
}

func (a Address) Line1() string {
	return a.line1
}

func (a Address) Line2() string {
	return a.line2
}

func (a Address) City() string {
	return a.city
}

func (a Address) Region() string {
	return a.region
}

func (a Address) PostalCode() string {
	return a.postalCode
}

func (a Address) Country() CountryCode {
	return a.country
}

// Equals compares two Address objects for equality.
func (a Address) Equals(other Address) bool {
	return a.line1 == other.line1 &&
		a.line2 == other.line2 &&
		a.city == other.city &&
		a.region == other.region &&
		a.postalCode == other.postalCode &&
		a.country.Equals(other.country)
}

// String returns a compact address representation.
func (a Address) String() string {
	parts := []string{a.line1}
	if a.line2 != "" {
		parts = append(parts, a.line2)
	}
	parts = append(parts, a.city)
	if a.region != "" {
		parts = append(parts, a.region)
	}
	if a.postalCode != "" {
		parts = append(parts, a.postalCode)
	}
	parts = append(parts, a.country.String())

	return strings.Join(parts, ", ")
}

// MarshalJSON returns the address as an explicit object.
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(addressJSON{
		Line1:      a.line1,
		Line2:      a.line2,
		City:       a.city,
		Region:     a.region,
		PostalCode: a.postalCode,
		Country:    a.country,
	})
}

// UnmarshalJSON validates and normalizes a JSON address object.
func (a *Address) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var raw addressJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	address, err := NewAddress(raw.Line1, raw.Line2, raw.City, raw.Region, raw.PostalCode, raw.Country)
	if err != nil {
		return err
	}

	*a = address
	return nil
}

func normalizeAddressPart(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
