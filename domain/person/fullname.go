package person

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golibry/go-common-domain/domain"
	"strings"
	"unicode"
	"unicode/utf8"
)

const MaxNamePartLength = 128

var (
	ErrEmptyNamePart        = domain.NewError("name part cannot be empty")
	ErrInvalidNamePartChars = domain.NewError("name part contains invalid characters")
	ErrTooLongNamePart      = domain.NewError("name part is too long")
)

type FullName struct {
	firstName  string
	middleName string
	lastName   string
}

type fullNameJSON struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
}

// NewFullName creates a new instance of FullName
func NewFullName(firstName, middleName, lastName string) (FullName, error) {
	normalizedFirst, err := NormalizeNamePart(firstName)
	if err != nil {
		return FullName{}, fmt.Errorf("%w (first name)", err)
	}

	normalizedMiddle, err := NormalizeNamePart(middleName)
	if err != nil && !errors.Is(err, ErrEmptyNamePart) {
		return FullName{}, fmt.Errorf("%w (middle name)", err)
	}

	normalizedLast, err := NormalizeNamePart(lastName)
	if err != nil {
		return FullName{}, fmt.Errorf("%w (last name)", err)
	}

	return FullName{
		firstName:  normalizedFirst,
		middleName: normalizedMiddle,
		lastName:   normalizedLast,
	}, nil
}

// ReconstituteFullName creates a new FullName instance without validation or normalization
func ReconstituteFullName(firstName, middleName, lastName string) FullName {
	return FullName{
		firstName:  firstName,
		middleName: middleName,
		lastName:   lastName,
	}
}

// NewFullNameFromJSON creates FullName from JSON bytes array
func NewFullNameFromJSON(data []byte) (FullName, error) {
    var temp fullNameJSON

    if err := json.Unmarshal(data, &temp); err != nil {
        return FullName{}, domain.NewErrorWithWrap(err, "failed to build full name from json")
    }

	newFullName, err := NewFullName(temp.FirstName, temp.MiddleName, temp.LastName)
	if err != nil {
		return FullName{}, err
	}

	return newFullName, nil
}

// FirstName returns the first name
func (f FullName) FirstName() string {
	return f.firstName
}

// MiddleName returns the middle name
func (f FullName) MiddleName() string {
	return f.middleName
}

// LastName returns the last name
func (f FullName) LastName() string {
	return f.lastName
}

// Equals compares two FullName objects for equality
func (f FullName) Equals(other FullName) bool {
	return f.firstName == other.firstName &&
		f.middleName == other.middleName &&
		f.lastName == other.lastName
}

// String returns a string representation of the full name
func (f FullName) String() string {
	if f.middleName == "" {
		return fmt.Sprintf("%s %s", f.firstName, f.lastName)
	}
	return fmt.Sprintf("%s %s %s", f.firstName, f.middleName, f.lastName)
}

// MarshalJSON implements json.Marshaler
func (f FullName) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		fullNameJSON{
			FirstName:  f.firstName,
			MiddleName: f.middleName,
			LastName:   f.lastName,
		},
	)
}

func NormalizeNamePart(namePart string) (string, error) {
	// Trim spaces from the beginning and end
	namePart = strings.TrimSpace(namePart)
	var result strings.Builder
	var prevRune rune

	for i, r := range namePart {
		// Skip this character if it's a special character (space, hyphen, apostrophe, period)
		// and the previous character was also a special character
		if i > 0 && (r == ' ' || r == '-' || r == '\'' || r == '.') &&
			(prevRune == ' ' || prevRune == '-' || prevRune == '\'' || prevRune == '.') {
			continue
		}

		result.WriteRune(r)
		prevRune = r
	}

 resultStr := result.String()
 if err := IsValidNamePart(resultStr); err != nil {
     return "", err
 }

	return resultStr, nil
}

func IsValidNamePart(namePart string) error {
	if namePart == "" {
		return ErrEmptyNamePart
	}

	if utf8.RuneCountInString(namePart) > MaxNamePartLength {
		return ErrTooLongNamePart
	}

	// Check if the namePart starts or ends with invalid characters
	firstRune, _ := utf8.DecodeRuneInString(namePart)
	lastRune, _ := utf8.DecodeLastRuneInString(namePart)

	if firstRune == '-' || firstRune == '\'' || firstRune == '.' {
		return ErrInvalidNamePartChars
	}
	if lastRune == '-' || lastRune == '\'' || lastRune == '.' {
		return ErrInvalidNamePartChars
	}

	for _, r := range namePart {
		// Check if the character is valid.
		// Valid characters: Unicode letters, spaces, hyphens, apostrophes, periods
		if !unicode.IsLetter(r) && r != ' ' && r != '-' && r != '\'' && r != '.' {
			return ErrInvalidNamePartChars
		}
	}

	return nil
}
