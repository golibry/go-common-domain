package identifier

import (
	"encoding/json"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrEmptyCNK   = domain.NewError("CNK cannot be empty")
	ErrInvalidCNK = domain.NewError("CNK must be 7 digits with a valid Luhn check digit")
)

type CNK struct {
	value string
}

// NewCNK creates a new CNK with validation and normalization.
func NewCNK(value string) (CNK, error) {
	normalized, err := NormalizeCNK(value)
	if err != nil {
		return CNK{}, err
	}

	return CNK{
		value: normalized,
	}, nil
}

// ReconstituteCNK creates a CNK from a trusted persisted value.
func ReconstituteCNK(value string) CNK {
	return CNK{
		value: value,
	}
}

// Value returns the CNK value.
func (c CNK) Value() string {
	return c.value
}

// Equals compares two CNK objects for equality.
func (c CNK) Equals(other CNK) bool {
	return c.value == other.value
}

// String returns the CNK value.
func (c CNK) String() string {
	return c.value
}

// MarshalText returns the CNK value as text.
func (c CNK) MarshalText() ([]byte, error) {
	return []byte(c.value), nil
}

// UnmarshalText validates and normalizes text into a CNK.
func (c *CNK) UnmarshalText(text []byte) error {
	cnk, err := NewCNK(string(text))
	if err != nil {
		return err
	}

	*c = cnk
	return nil
}

// UnmarshalJSON validates and normalizes a JSON CNK string.
func (c *CNK) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return ErrInvalidCNK
	}

	return c.UnmarshalText([]byte(value))
}

// NormalizeCNK removes common separators from a CNK.
func NormalizeCNK(value string) (string, error) {
	normalized := normalizeCode(value)
	if err := IsValidCNK(normalized); err != nil {
		return "", err
	}

	return normalized, nil
}

// IsValidCNK validates a 7-digit CNK with a Luhn check digit.
func IsValidCNK(value string) error {
	if value == "" {
		return ErrEmptyCNK
	}

	if len(value) != 7 || !hasLuhnCheckDigit(value) {
		return ErrInvalidCNK
	}

	return nil
}
