package identifier

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrEmptyUUIDIdentifier   = domain.NewError("UUID identifier cannot be empty")
	ErrInvalidUUIDIdentifier = domain.NewError("UUID identifier format is invalid")
	ErrZeroUUIDIdentifier    = domain.NewError("UUID identifier cannot be zero")
)

var uuidIdentifierRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

const zeroUUIDIdentifier = "00000000-0000-0000-0000-000000000000"

type UUIDIdentifier struct {
	value string
}

// NewUUIDIdentifier creates a new UUIDIdentifier with validation and normalization.
func NewUUIDIdentifier(value string) (UUIDIdentifier, error) {
	normalized, err := NormalizeUUIDIdentifier(value)
	if err != nil {
		return UUIDIdentifier{}, err
	}

	return UUIDIdentifier{
		value: normalized,
	}, nil
}

// ReconstituteUUIDIdentifier creates a UUIDIdentifier from a trusted persisted value.
func ReconstituteUUIDIdentifier(value string) UUIDIdentifier {
	return UUIDIdentifier{
		value: value,
	}
}

// Value returns the UUID value.
func (i UUIDIdentifier) Value() string {
	return i.value
}

// Equals compares two UUIDIdentifier objects for equality.
func (i UUIDIdentifier) Equals(other UUIDIdentifier) bool {
	return i.value == other.value
}

// String returns the UUID value.
func (i UUIDIdentifier) String() string {
	return i.value
}

// MarshalText returns the UUID value as text.
func (i UUIDIdentifier) MarshalText() ([]byte, error) {
	return []byte(i.value), nil
}

// UnmarshalText validates and normalizes text into a UUIDIdentifier.
func (i *UUIDIdentifier) UnmarshalText(text []byte) error {
	identifier, err := NewUUIDIdentifier(string(text))
	if err != nil {
		return err
	}

	*i = identifier
	return nil
}

// UnmarshalJSON validates and normalizes a JSON UUID string.
func (i *UUIDIdentifier) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return ErrInvalidUUIDIdentifier
	}

	return i.UnmarshalText([]byte(value))
}

// NormalizeUUIDIdentifier trims spaces and lowercases a UUID.
func NormalizeUUIDIdentifier(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if err := IsValidUUIDIdentifier(normalized); err != nil {
		return "", err
	}

	return normalized, nil
}

// IsValidUUIDIdentifier validates a canonical UUID string.
func IsValidUUIDIdentifier(value string) error {
	if value == "" {
		return ErrEmptyUUIDIdentifier
	}

	if !uuidIdentifierRegex.MatchString(value) {
		return ErrInvalidUUIDIdentifier
	}

	if value == zeroUUIDIdentifier {
		return ErrZeroUUIDIdentifier
	}

	return nil
}
