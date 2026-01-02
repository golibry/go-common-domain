package person

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golibry/go-common-domain/domain"
)

const MaxNamePartLength = 128

var (
	ErrEmptyNamePart        = domain.NewError("name part cannot be empty")
	ErrInvalidNamePartChars = domain.NewError("name part contains invalid characters; allowed: letters (Unicode), spaces, hyphens (-), apostrophes ('), and periods (.). Name parts cannot start or end with a hyphen, apostrophe, or period.")
	ErrTooLongNamePart      = domain.NewError("name part is too long")
)

type FullName struct {
	firstName  string
	middleName string
	lastName   string
}

// NewFullName creates a new instance of FullName.
// It normalizes parts (trimming and collapsing repeated separators) and validates them.
// Allowed characters are Unicode letters, spaces, hyphens (-), apostrophes ('), and periods (.).
// Name parts cannot start or end with a hyphen, apostrophe, or period.
// The middle name can be empty or a single-letter initial followed by a period (e.g., "F.").
func NewFullName(firstName, middleName, lastName string) (FullName, error) {
	normalizedFirst, _ := NormalizeNamePart(firstName)
	if err := IsValidNamePart(normalizedFirst); err != nil {
		return FullName{}, fmt.Errorf("%w (first name)", err)
	}

	normalizedMiddle, _ := NormalizeNamePart(middleName)
	if normalizedMiddle != "" {
		if err := IsValidNamePart(normalizedMiddle); err != nil {
			if !isInitialWithPeriod(normalizedMiddle) {
				return FullName{}, fmt.Errorf("%w (middle name)", err)
			}
		}
	}

	normalizedLast, _ := NormalizeNamePart(lastName)
	if err := IsValidNamePart(normalizedLast); err != nil {
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
	// Note: Validation is intentionally separated from normalization.
	// Callers should validate the normalized value via IsValidNamePart or custom rules.
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

// isInitialWithPeriod reports whether the provided string is a single
// Unicode letter followed by a period, e.g., "F.". This is allowed
// for the middle name only.
func isInitialWithPeriod(s string) bool {
	if utf8.RuneCountInString(s) != 2 {
		return false
	}
	r1, size := utf8.DecodeRuneInString(s)
	r2, _ := utf8.DecodeRuneInString(s[size:])
	return unicode.IsLetter(r1) && r2 == '.'
}
