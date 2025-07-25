package web

import (
	"encoding/json"
	"github.com/golibry/go-common-domain/domain"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	MaxDomainNameLength = 253
	MaxLabelLength      = 63
	MinDomainNameLength = 1
)

var (
	ErrEmptyDomainName        = domain.NewError("domain name cannot be empty")
	ErrInvalidDomainNameChars = domain.NewError("domain name contains invalid characters")
	ErrTooLongDomainName      = domain.NewError("domain name is too long")
	ErrTooLongDomainLabel     = domain.NewError("domain name label is too long")
	ErrInvalidDomainFormat    = domain.NewError("domain name has invalid format")
	ErrConsecutiveDots        = domain.NewError("domain name cannot have consecutive dots")
	ErrStartsOrEndsWithDot    = domain.NewError("domain name cannot start or end with a dot")
	ErrStartsOrEndsWithHyphen = domain.NewError("domain name label cannot start or end with hyphen")
)

// domainNameRegex validates basic domain name format
var domainNameRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)

type DomainName struct {
	value string
}

type domainNameJSON struct {
	Value string `json:"value"`
}

// NewDomainName creates a new instance of DomainName with validation and normalization
func NewDomainName(value string) (DomainName, error) {
	normalized, err := NormalizeDomainName(value)
	if err != nil {
		return DomainName{}, err
	}

	return DomainName{
		value: normalized,
	}, nil
}

// ReconstituteDomainName creates a new DomainName instance without validation or normalization
func ReconstituteDomainName(value string) DomainName {
	return DomainName{
		value: value,
	}
}

// NewDomainNameFromJSON creates DomainName from JSON bytes array
func NewDomainNameFromJSON(data []byte) (DomainName, error) {
	var temp domainNameJSON

	if err := json.Unmarshal(data, &temp); err != nil {
		return DomainName{}, domain.NewError("failed to build domain name from json: %s", err)
	}

	newDomainName, err := NewDomainName(temp.Value)
	if err != nil {
		return DomainName{}, err
	}

	return newDomainName, nil
}

// Value returns the domain name value
func (d DomainName) Value() string {
	return d.value
}

// Equals compares two DomainName objects for equality
func (d DomainName) Equals(other DomainName) bool {
	return d.value == other.value
}

// String returns a string representation of the domain name
func (d DomainName) String() string {
	return d.value
}

// MarshalJSON implements json.Marshaler
func (d DomainName) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		domainNameJSON{
			Value: d.value,
		},
	)
}

// NormalizeDomainName normalizes a domain name by converting to lowercase and trimming spaces
func NormalizeDomainName(domainName string) (string, error) {
	// Trim spaces from the beginning and end
	domainName = strings.TrimSpace(domainName)

	// Convert to lowercase
	domainName = strings.ToLower(domainName)

	if err := IsValidDomainName(domainName); err != nil {
		return domainName, err
	}

	return domainName, nil
}

// IsValidDomainName validates a domain name according to RFC standards
func IsValidDomainName(domainName string) error {
	if domainName == "" {
		return ErrEmptyDomainName
	}

	if utf8.RuneCountInString(domainName) > MaxDomainNameLength {
		return ErrTooLongDomainName
	}

	if utf8.RuneCountInString(domainName) < MinDomainNameLength {
		return ErrEmptyDomainName
	}

	// Check for consecutive dots
	if strings.Contains(domainName, "..") {
		return ErrConsecutiveDots
	}

	// Check if starts or ends with dot
	if strings.HasPrefix(domainName, ".") || strings.HasSuffix(domainName, ".") {
		return ErrStartsOrEndsWithDot
	}

	// Split into labels and validate each
	labels := strings.Split(domainName, ".")
	for _, label := range labels {
		if err := isValidDomainLabel(label); err != nil {
			return err
		}
	}

	// Use regex for final validation
	if !domainNameRegex.MatchString(domainName) {
		return ErrInvalidDomainFormat
	}

	return nil
}

// isValidDomainLabel validates a single domain label
func isValidDomainLabel(label string) error {
	if label == "" {
		return ErrInvalidDomainFormat
	}

	if utf8.RuneCountInString(label) > MaxLabelLength {
		return ErrTooLongDomainLabel
	}

	// Check if starts or ends with hyphen
	if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
		return ErrStartsOrEndsWithHyphen
	}

	// Check for valid characters (letters, numbers, hyphens only)
	for _, r := range label {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return ErrInvalidDomainNameChars
		}
	}

	return nil
}
