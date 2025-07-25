package contact

import (
	"encoding/json"
	"github.com/golibry/go-common-domain/domain"
	"regexp"
	"strings"
	"unicode"
)

const MaxPhoneNumberLength = 20

var (
	ErrEmptyPhoneNumber        = domain.NewError("phone number cannot be empty")
	ErrInvalidPhoneNumberChars = domain.NewError("phone number contains invalid characters")
	ErrTooLongPhoneNumber      = domain.NewError("phone number is too long")
	ErrTooShortPhoneNumber     = domain.NewError("phone number is too short")
)

var phoneNumberRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

type PhoneNumber struct {
	value string
}

type phoneNumberJSON struct {
	Value string `json:"value"`
}

// NewPhoneNumber creates a new instance of PhoneNumber with validation and normalization
func NewPhoneNumber(value string) (PhoneNumber, error) {
	normalized, err := NormalizePhoneNumber(value)
	if err != nil {
		return PhoneNumber{}, err
	}

	return PhoneNumber{
		value: normalized,
	}, nil
}

// ReconstitutePhoneNumber creates a new PhoneNumber instance without validation or normalization
func ReconstitutePhoneNumber(value string) PhoneNumber {
	return PhoneNumber{
		value: value,
	}
}

// NewPhoneNumberFromJSON creates PhoneNumber from JSON bytes array
func NewPhoneNumberFromJSON(data []byte) (PhoneNumber, error) {
	var temp phoneNumberJSON

	if err := json.Unmarshal(data, &temp); err != nil {
		return PhoneNumber{}, domain.NewError("failed to build phone number from json: %s", err)
	}

	newPhoneNumber, err := NewPhoneNumber(temp.Value)
	if err != nil {
		return PhoneNumber{}, err
	}

	return newPhoneNumber, nil
}

// Value returns the phone number value
func (p PhoneNumber) Value() string {
	return p.value
}

// Equals compares two PhoneNumber objects for equality
func (p PhoneNumber) Equals(other PhoneNumber) bool {
	return p.value == other.value
}

// String returns a string representation of the phone number
func (p PhoneNumber) String() string {
	return p.value
}

// MarshalJSON implements json.Marshaler
func (p PhoneNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		phoneNumberJSON{
			Value: p.value,
		},
	)
}

// NormalizePhoneNumber normalizes a phone number by removing spaces, dashes, parentheses, and dots
func NormalizePhoneNumber(phoneNumber string) (string, error) {
	// Trim spaces from the beginning and end
	phoneNumber = strings.TrimSpace(phoneNumber)

	// First check for invalid characters before normalization
	for _, r := range phoneNumber {
		// Allow digits, plus sign, spaces, dashes, parentheses, and dots
		if !unicode.IsDigit(r) && r != '+' && r != ' ' && r != '-' && r != '(' && r != ')' && r != '.' {
			return "", ErrInvalidPhoneNumberChars
		}
	}

	var result strings.Builder

	for _, r := range phoneNumber {
		// Keep only digits and plus sign
		if unicode.IsDigit(r) || r == '+' {
			result.WriteRune(r)
		}
	}

	normalized := result.String()

	if err := IsValidPhoneNumber(normalized); err != nil {
		return "", err
	}

	return normalized, nil
}

// IsValidPhoneNumber validates a phone number
func IsValidPhoneNumber(phoneNumber string) error {
	if phoneNumber == "" {
		return ErrEmptyPhoneNumber
	}

	if len(phoneNumber) > MaxPhoneNumberLength {
		return ErrTooLongPhoneNumber
	}

	if len(phoneNumber) < 3 {
		return ErrTooShortPhoneNumber
	}

	// Check for invalid characters (should only contain digits and optionally start with +)
	for i, r := range phoneNumber {
		if i == 0 && r == '+' {
			continue
		}
		if !unicode.IsDigit(r) {
			return ErrInvalidPhoneNumberChars
		}
	}

	// Use regex for final validation
	if !phoneNumberRegex.MatchString(phoneNumber) {
		return ErrInvalidPhoneNumberChars
	}

	return nil
}
