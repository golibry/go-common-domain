package identifier

import (
	"encoding/json"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrEmptyGTIN   = domain.NewError("GTIN cannot be empty")
	ErrInvalidGTIN = domain.NewError("GTIN must be 8, 12, 13, or 14 digits with a valid check digit")
)

type GTIN struct {
	value string
}

// NewGTIN creates a new GTIN with validation and normalization.
func NewGTIN(value string) (GTIN, error) {
	normalized, err := NormalizeGTIN(value)
	if err != nil {
		return GTIN{}, err
	}

	return GTIN{
		value: normalized,
	}, nil
}

// ReconstituteGTIN creates a GTIN from a trusted persisted value.
func ReconstituteGTIN(value string) GTIN {
	return GTIN{
		value: value,
	}
}

// Value returns the GTIN value.
func (g GTIN) Value() string {
	return g.value
}

// Equals compares two GTIN objects for equality.
func (g GTIN) Equals(other GTIN) bool {
	return g.value == other.value
}

// String returns the GTIN value.
func (g GTIN) String() string {
	return g.value
}

// MarshalText returns the GTIN value as text.
func (g GTIN) MarshalText() ([]byte, error) {
	return []byte(g.value), nil
}

// UnmarshalText validates and normalizes text into a GTIN.
func (g *GTIN) UnmarshalText(text []byte) error {
	gtin, err := NewGTIN(string(text))
	if err != nil {
		return err
	}

	*g = gtin
	return nil
}

// UnmarshalJSON validates and normalizes a JSON GTIN string.
func (g *GTIN) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return ErrInvalidGTIN
	}

	return g.UnmarshalText([]byte(value))
}

// NormalizeGTIN removes common separators from a GTIN.
func NormalizeGTIN(value string) (string, error) {
	normalized := normalizeCode(value)
	if err := IsValidGTIN(normalized); err != nil {
		return "", err
	}

	return normalized, nil
}

// IsValidGTIN validates a GTIN-8, GTIN-12, GTIN-13, or GTIN-14.
func IsValidGTIN(value string) error {
	if value == "" {
		return ErrEmptyGTIN
	}

	switch len(value) {
	case 8, 12, 13, 14:
	default:
		return ErrInvalidGTIN
	}

	if !hasGS1CheckDigit(value) {
		return ErrInvalidGTIN
	}

	return nil
}
