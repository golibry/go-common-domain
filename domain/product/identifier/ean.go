package identifier

import (
	"encoding/json"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrEmptyEAN   = domain.NewError("EAN cannot be empty")
	ErrInvalidEAN = domain.NewError("EAN must be 8 or 13 digits with a valid check digit")
)

type EAN struct {
	value string
}

// NewEAN creates a new EAN with validation and normalization.
func NewEAN(value string) (EAN, error) {
	normalized, err := NormalizeEAN(value)
	if err != nil {
		return EAN{}, err
	}

	return EAN{
		value: normalized,
	}, nil
}

// ReconstituteEAN creates an EAN from a trusted persisted value.
func ReconstituteEAN(value string) EAN {
	return EAN{
		value: value,
	}
}

// Value returns the EAN value.
func (e EAN) Value() string {
	return e.value
}

// Equals compares two EAN objects for equality.
func (e EAN) Equals(other EAN) bool {
	return e.value == other.value
}

// String returns the EAN value.
func (e EAN) String() string {
	return e.value
}

// MarshalText returns the EAN value as text.
func (e EAN) MarshalText() ([]byte, error) {
	return []byte(e.value), nil
}

// UnmarshalText validates and normalizes text into an EAN.
func (e *EAN) UnmarshalText(text []byte) error {
	ean, err := NewEAN(string(text))
	if err != nil {
		return err
	}

	*e = ean
	return nil
}

// UnmarshalJSON validates and normalizes a JSON EAN string.
func (e *EAN) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return ErrInvalidEAN
	}

	return e.UnmarshalText([]byte(value))
}

// NormalizeEAN removes common separators from an EAN.
func NormalizeEAN(value string) (string, error) {
	normalized := normalizeCode(value)
	if err := IsValidEAN(normalized); err != nil {
		return "", err
	}

	return normalized, nil
}

// IsValidEAN validates an EAN-8 or EAN-13.
func IsValidEAN(value string) error {
	if value == "" {
		return ErrEmptyEAN
	}

	switch len(value) {
	case 8, 13:
	default:
		return ErrInvalidEAN
	}

	if !hasGS1CheckDigit(value) {
		return ErrInvalidEAN
	}

	return nil
}
