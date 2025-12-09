package web

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/golibry/go-common-domain/domain"
)

const (
	MaxEmailLength      = 254 // RFC 5321 limit
	MaxLocalPartLength  = 64  // RFC 5321 limit
	MaxDomainPartLength = 253 // RFC 5321 limit
	MinEmailLength      = 3   // a@b minimum
)

var (
	ErrEmptyEmail         = domain.NewError("email address cannot be empty")
	ErrInvalidEmailFormat = domain.NewError("email address has invalid format")
	ErrTooLongEmail       = domain.NewError("email address is too long")
	ErrTooLongLocalPart   = domain.NewError("email local part is too long")
	ErrTooLongDomainPart  = domain.NewError("email domain part is too long")
	ErrInvalidEmailChars  = domain.NewError("email address contains invalid characters")
	ErrMissingAtSymbol    = domain.NewError("email address must contain exactly one @ symbol")
	ErrMultipleAtSymbols  = domain.NewError("email address cannot contain multiple @ symbols")
	ErrEmptyLocalPart     = domain.NewError("email local part cannot be empty")
	ErrEmptyDomainPart    = domain.NewError("email domain part cannot be empty")
	ErrInvalidLocalPart   = domain.NewError("email local part has invalid format")
	ErrInvalidDomainPart  = domain.NewError("email domain part has invalid format")
)

// emailRegex validates basic email format according to RFC 5322 (simplified)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

type Email struct {
	value string
}

type emailJSON struct {
	Value string `json:"value"`
}

// NewEmail creates a new instance of Email with validation and normalization
func NewEmail(value string) (Email, error) {
	normalized, err := NormalizeEmail(value)
	if err != nil {
		return Email{}, err
	}

	return Email{
		value: normalized,
	}, nil
}

// ReconstituteEmail creates a new Email instance without validation or normalization
func ReconstituteEmail(value string) Email {
	return Email{
		value: value,
	}
}

// NewEmailFromJSON creates Email from JSON bytes array
func NewEmailFromJSON(data []byte) (Email, error) {
	var temp emailJSON

	if err := json.Unmarshal(data, &temp); err != nil {
		return Email{}, domain.NewErrorWithWrap(err, "failed to build email from json")
	}

	newEmail, err := NewEmail(temp.Value)
	if err != nil {
		return Email{}, err
	}

	return newEmail, nil
}

// Value returns the email address value
func (e Email) Value() string {
	return e.value
}

// LocalPart returns the local part of the email address (before @)
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// DomainPart returns the domain part of the email address (after @)
func (e Email) DomainPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// Equals compares two Email objects for equality
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// String returns a string representation of the email address
func (e Email) String() string {
	return e.value
}

// MarshalJSON implements json.Marshaler
func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		emailJSON{
			Value: e.value,
		},
	)
}

// NormalizeEmail normalizes an email address by converting to lowercase and trimming spaces
func NormalizeEmail(email string) (string, error) {
	// Trim spaces from the beginning and end
	email = strings.TrimSpace(email)

	// Convert to lowercase
	email = strings.ToLower(email)

	if err := IsValidEmail(email); err != nil {
		return "", err
	}

	return email, nil
}

// IsValidEmail validates an email address according to RFC standards
func IsValidEmail(email string) error {
	if email == "" {
		return ErrEmptyEmail
	}

	if utf8.RuneCountInString(email) > MaxEmailLength {
		return ErrTooLongEmail
	}

	// Check for exactly one @ symbol
	atCount := strings.Count(email, "@")
	if atCount == 0 {
		return ErrMissingAtSymbol
	}
	if atCount > 1 {
		return ErrMultipleAtSymbols
	}

	// Split into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ErrInvalidEmailFormat
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Validate local part
	if err := isValidLocalPart(localPart); err != nil {
		return err
	}

	// Validate domain part
	if err := isValidEmailDomainPart(domainPart); err != nil {
		return err
	}

	// Check minimum length after validating parts (for more specific error messages)
	if utf8.RuneCountInString(email) < MinEmailLength {
		return ErrInvalidEmailFormat
	}

	// Use regex for final validation
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmailFormat
	}

	return nil
}

// isValidLocalPart validates the local part of an email address (before @)
func isValidLocalPart(localPart string) error {
	if localPart == "" {
		return ErrEmptyLocalPart
	}

	if utf8.RuneCountInString(localPart) > MaxLocalPartLength {
		return ErrTooLongLocalPart
	}

	// Check for invalid starting/ending characters
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return ErrInvalidLocalPart
	}

	// Check for consecutive dots
	if strings.Contains(localPart, "..") {
		return ErrInvalidLocalPart
	}

	// Check for valid characters in local part
	// RFC 5322 allows: a-z A-Z 0-9 . ! # $ % & ' * + - / = ? ^ _ ` { | } ~
	for _, r := range localPart {
		if !isValidLocalPartChar(r) {
			return ErrInvalidEmailChars
		}
	}

	return nil
}

// isValidEmailDomainPart validates the domain part of an email address (after @)
func isValidEmailDomainPart(domainPart string) error {
	if domainPart == "" {
		return ErrEmptyDomainPart
	}

	if utf8.RuneCountInString(domainPart) > MaxDomainPartLength {
		return ErrTooLongDomainPart
	}

	// Use the existing domain validation logic
	if err := IsValidDomainName(domainPart); err != nil {
		return ErrInvalidDomainPart
	}

	return nil
}

// isValidLocalPartChar checks if a character is valid in the local part of an email
func isValidLocalPartChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '.' || r == '!' || r == '#' || r == '$' || r == '%' ||
		r == '&' || r == '\'' || r == '*' || r == '+' || r == '-' ||
		r == '/' || r == '=' || r == '?' || r == '^' || r == '_' ||
		r == '`' || r == '{' || r == '|' || r == '}' || r == '~'
}
