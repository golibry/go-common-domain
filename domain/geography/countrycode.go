package geography

import (
	"encoding/json"
	"github.com/golibry/go-common-domain/domain"
	"regexp"
	"strings"
)

var (
	ErrEmptyCountryCode   = domain.NewError("country code cannot be empty")
	ErrInvalidCountryCode = domain.NewError("country code must be exactly 2 letters")
)

var countryCodeRegex = regexp.MustCompile(`^[A-Z]{2}$`)

type CountryCode struct {
	value string
}

type countryCodeJSON struct {
	Value string `json:"value"`
}

// NewCountryCode creates a new instance of CountryCode with validation and normalization
func NewCountryCode(value string) (CountryCode, error) {
	normalized, err := NormalizeCountryCode(value)
	if err != nil {
		return CountryCode{}, err
	}

	return CountryCode{
		value: normalized,
	}, nil
}

// ReconstituteCountryCode creates a new CountryCode instance without validation or normalization
func ReconstituteCountryCode(value string) CountryCode {
	return CountryCode{
		value: value,
	}
}

// NewCountryCodeFromJSON creates CountryCode from JSON bytes array
func NewCountryCodeFromJSON(data []byte) (CountryCode, error) {
	var temp countryCodeJSON

	if err := json.Unmarshal(data, &temp); err != nil {
		return CountryCode{}, domain.NewError("failed to build country code from json: %s", err)
	}

	newCountryCode, err := NewCountryCode(temp.Value)
	if err != nil {
		return CountryCode{}, err
	}

	return newCountryCode, nil
}

// Value returns the country code value
func (c CountryCode) Value() string {
	return c.value
}

// Equals compares two CountryCode objects for equality
func (c CountryCode) Equals(other CountryCode) bool {
	return c.value == other.value
}

// String returns a string representation of the country code
func (c CountryCode) String() string {
	return c.value
}

// MarshalJSON implements json.Marshaler
func (c CountryCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		countryCodeJSON{
			Value: c.value,
		},
	)
}

// NormalizeCountryCode normalizes a country code by trimming spaces and converting to uppercase
func NormalizeCountryCode(countryCode string) (string, error) {
	// Trim spaces and convert to uppercase
	normalized := strings.ToUpper(strings.TrimSpace(countryCode))
	
	if err := IsValidCountryCode(normalized); err != nil {
		return "", err
	}
	
	return normalized, nil
}

// IsValidCountryCode validates a country code (must be exactly 2 uppercase letters)
func IsValidCountryCode(countryCode string) error {
	if countryCode == "" {
		return ErrEmptyCountryCode
	}
	
	if !countryCodeRegex.MatchString(countryCode) {
		return ErrInvalidCountryCode
	}
	
	return nil
}