package identifier

import (
	"strconv"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrZeroIdentifier    = domain.NewError("identifier cannot be zero")
	ErrInvalidIdentifier = domain.NewError("identifier format is invalid")
)

type IntIdentifier struct {
	value uint64
}

// NewIntIdentifier creates a new instance of IntIdentifier with validation
func NewIntIdentifier(value uint64) (IntIdentifier, error) {
	if err := IsValidIntIdentifier(value); err != nil {
		return IntIdentifier{}, err
	}

	return IntIdentifier{
		value: value,
	}, nil
}

// NewIntIdentifierFromInt creates a new instance of IntIdentifier from int64
func NewIntIdentifierFromInt(value int64) (IntIdentifier, error) {
	return NewIntIdentifier(uint64(value))
}

// NewIntIdentifierFromString creates a new instance of IntIdentifier from string
func NewIntIdentifierFromString(value string) (IntIdentifier, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return IntIdentifier{}, ErrInvalidIdentifier
	}

	return NewIntIdentifier(parsed)
}

// ReconstituteIntIdentifier creates a new IntIdentifier instance without validation
func ReconstituteIntIdentifier(value uint64) IntIdentifier {
	return IntIdentifier{
		value: value,
	}
}

// Value returns the identifier value as uint64
func (i IntIdentifier) Value() uint64 {
	return i.value
}

// Equals compares two IntIdentifier objects for equality
func (i IntIdentifier) Equals(other IntIdentifier) bool {
	return i.value == other.value
}

// String returns a string representation of the identifier
func (i IntIdentifier) String() string {
	return strconv.FormatUint(i.value, 10)
}

// IsValidIntIdentifier validates an identifier (must be positive and non-zero)
func IsValidIntIdentifier(value uint64) error {
	if value == 0 {
		return ErrZeroIdentifier
	}

	return nil
}
