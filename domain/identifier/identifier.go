package identifier

import (
	"encoding/json"
	"github.com/golibry/go-common-domain/domain"
	"strconv"
)

var (
	ErrZeroIdentifier     = domain.NewError("identifier cannot be zero")
	ErrNegativeIdentifier = domain.NewError("identifier cannot be negative")
	ErrInvalidIdentifier  = domain.NewError("identifier format is invalid")
)

type Identifier struct {
	value uint64
}

type identifierJSON struct {
	Value uint64 `json:"value"`
}

// NewIdentifier creates a new instance of Identifier with validation
func NewIdentifier(value uint64) (Identifier, error) {
	if err := IsValidIdentifier(value); err != nil {
		return Identifier{}, err
	}

	return Identifier{
		value: value,
	}, nil
}

// NewIdentifierFromInt creates a new instance of Identifier from int64
func NewIdentifierFromInt(value int64) (Identifier, error) {
	if value < 0 {
		return Identifier{}, ErrNegativeIdentifier
	}

	return NewIdentifier(uint64(value))
}

// NewIdentifierFromString creates a new instance of Identifier from string
func NewIdentifierFromString(value string) (Identifier, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return Identifier{}, ErrInvalidIdentifier
	}

	return NewIdentifier(parsed)
}

// ReconstituteIdentifier creates a new Identifier instance without validation
func ReconstituteIdentifier(value uint64) Identifier {
	return Identifier{
		value: value,
	}
}

// NewIdentifierFromJSON creates Identifier from JSON bytes array
func NewIdentifierFromJSON(data []byte) (Identifier, error) {
    var temp identifierJSON

    if err := json.Unmarshal(data, &temp); err != nil {
        return Identifier{}, domain.NewErrorWithWrap(err, "failed to build identifier from json")
    }

	return NewIdentifier(temp.Value)
}

// Value returns the identifier value as uint64
func (i Identifier) Value() uint64 {
	return i.value
}

// Equals compares two Identifier objects for equality
func (i Identifier) Equals(other Identifier) bool {
	return i.value == other.value
}

// String returns a string representation of the identifier
func (i Identifier) String() string {
	return strconv.FormatUint(i.value, 10)
}

// MarshalJSON implements json.Marshaler
func (i Identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		identifierJSON{
			Value: i.value,
		},
	)
}

// IsValidIdentifier validates an identifier (must be positive and non-zero)
func IsValidIdentifier(value uint64) error {
	if value == 0 {
		return ErrZeroIdentifier
	}

	return nil
}
