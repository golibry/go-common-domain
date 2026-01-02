package web

import (
	"net/url"
	"strings"

	"github.com/golibry/go-common-domain/domain"
)

const MaxURLLength = 2048

var (
	ErrEmptyURL   = domain.NewError("URL cannot be empty")
	ErrInvalidURL = domain.NewError("URL format is invalid")
	ErrTooLongURL = domain.NewError("URL is too long")
)

type URL struct {
	value string
}

// NewURL creates a new instance of URL with validation and normalization
func NewURL(value string) (URL, error) {
	normalized, err := NormalizeURL(value)
	if err != nil {
		return URL{}, err
	}

	return URL{
		value: normalized,
	}, nil
}

// ReconstituteURL creates a new URL instance without validation or normalization.
//
// ReconstituteURL should only be used with values that were previously validated
// and persisted by this package. Passing an arbitrary or invalid value may cause
// component accessors like Scheme(), Host(), and Path() to return empty results
// due to failed parsing.
func ReconstituteURL(value string) URL {
	return URL{
		value: value,
	}
}

// Value returns the URL value
func (u URL) Value() string {
	return u.value
}

// String returns a string representation of the phone number
func (p URL) String() string {
	return p.value
}

// Parsed returns the URL representation of the raw url string (value)
func (u URL) Parsed() url.URL {
	parsed, _ := url.Parse(u.value)
	if parsed == nil {
		return url.URL{}
	}
	return *parsed
}

// Equals compares two URL objects for equality
func (u URL) Equals(other URL) bool {
	return u.value == other.value
}

// NormalizeURL normalizes a URL by trimming spaces and ensuring a proper format
func NormalizeURL(urlStr string) (string, error) {
	// Trim spaces from the beginning and end
	urlStr = strings.TrimSpace(urlStr)

	parsed, err := IsValidURL(urlStr)
	if err != nil {
		return "", err
	}

	return parsed.String(), nil
}

// IsValidURL validates a URL
func IsValidURL(urlStr string) (*url.URL, error) {
	if urlStr == "" {
		return nil, ErrEmptyURL
	}

	if len(urlStr) > MaxURLLength {
		return nil, ErrTooLongURL
	}

	// Parse the URL to check if it's valid
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return nil, ErrInvalidURL
	}

	// Check if scheme and host are present
	if parsed.Scheme == "" {
		return nil, ErrInvalidURL
	}

	// Enforce allowed schemes
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
		// ok
	default:
		return nil, ErrInvalidURL
	}

	// For absolute URLs, the host should be present
	if parsed.IsAbs() && parsed.Host == "" {
		return nil, ErrInvalidURL
	}

	return parsed, nil
}
