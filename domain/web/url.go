package web

import (
	"encoding/json"
	"github.com/golibry/go-common-domain/domain"
	"net/url"
	"strings"
)

const MaxURLLength = 2048

var (
	ErrEmptyURL    = domain.NewError("URL cannot be empty")
	ErrInvalidURL  = domain.NewError("URL format is invalid")
	ErrTooLongURL  = domain.NewError("URL is too long")
)

type URL struct {
	value string
}

type urlJSON struct {
	Value string `json:"value"`
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

// ReconstituteURL creates a new URL instance without validation or normalization
func ReconstituteURL(value string) URL {
	return URL{
		value: value,
	}
}

// NewURLFromJSON creates URL from JSON bytes array
func NewURLFromJSON(data []byte) (URL, error) {
	var temp urlJSON

	if err := json.Unmarshal(data, &temp); err != nil {
		return URL{}, domain.NewError("failed to build URL from json: %s", err)
	}

	newURL, err := NewURL(temp.Value)
	if err != nil {
		return URL{}, err
	}

	return newURL, nil
}

// Value returns the URL value
func (u URL) Value() string {
	return u.value
}

// Scheme returns the URL scheme (e.g., "https", "http")
func (u URL) Scheme() string {
	parsed, _ := url.Parse(u.value)
	return parsed.Scheme
}

// Host returns the URL host
func (u URL) Host() string {
	parsed, _ := url.Parse(u.value)
	return parsed.Host
}

// Path returns the URL path
func (u URL) Path() string {
	parsed, _ := url.Parse(u.value)
	return parsed.Path
}

// Equals compares two URL objects for equality
func (u URL) Equals(other URL) bool {
	return u.value == other.value
}

// String returns a string representation of the URL
func (u URL) String() string {
	return u.value
}

// MarshalJSON implements json.Marshaler
func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		urlJSON{
			Value: u.value,
		},
	)
}

// NormalizeURL normalizes a URL by trimming spaces and ensuring proper format
func NormalizeURL(urlStr string) (string, error) {
	// Trim spaces from the beginning and end
	urlStr = strings.TrimSpace(urlStr)
	
	if err := IsValidURL(urlStr); err != nil {
		return "", err
	}
	
	// Parse and reconstruct to normalize the URL
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return "", ErrInvalidURL
	}
	
	return parsed.String(), nil
}

// IsValidURL validates a URL
func IsValidURL(urlStr string) error {
	if urlStr == "" {
		return ErrEmptyURL
	}
	
	if len(urlStr) > MaxURLLength {
		return ErrTooLongURL
	}
	
	// Parse the URL to check if it's valid
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return ErrInvalidURL
	}
	
	// Check if scheme and host are present
	if parsed.Scheme == "" {
		return ErrInvalidURL
	}
	
	// For absolute URLs, host should be present
	if parsed.IsAbs() && parsed.Host == "" {
		return ErrInvalidURL
	}
	
	return nil
}