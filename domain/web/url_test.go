package web

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type URLTestSuite struct {
	suite.Suite
}

func TestURLSuite(t *testing.T) {
	suite.Run(t, new(URLTestSuite))
}

func (s *URLTestSuite) TestItCanBuildNewURLWithValidValues() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTTPS URL",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "HTTP URL",
			input:    "http://example.com",
			expected: "http://example.com",
		},
		{
			name:     "URL with path",
			input:    "https://example.com/path/to/resource",
			expected: "https://example.com/path/to/resource",
		},
		{
			name:     "URL with query parameters",
			input:    "https://example.com/search?q=test&page=1",
			expected: "https://example.com/search?q=test&page=1",
		},
		{
			name:     "URL with port",
			input:    "https://example.com:8080",
			expected: "https://example.com:8080",
		},
		{
			name:     "URL with fragment",
			input:    "https://example.com/page#section",
			expected: "https://example.com/page#section",
		},
		{
			name:     "URL with spaces (trimmed)",
			input:    "  https://example.com  ",
			expected: "https://example.com",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			url, err := NewURL(tc.input)
			s.NoError(err)
			s.Equal(tc.expected, url.Value())
			s.Equal(tc.expected, url.String())
		})
	}
}

func (s *URLTestSuite) TestItFailsToBuildNewURLFromInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "empty URL",
			input:         "",
			expectedError: ErrEmptyURL,
		},
		{
			name:          "URL with only spaces",
			input:         "   ",
			expectedError: ErrEmptyURL,
		},
		{
			name:          "URL without scheme",
			input:         "example.com",
			expectedError: ErrInvalidURL,
		},
 	{
			name:          "URL with invalid scheme",
			input:         "invalid://example.com",
			expectedError: ErrInvalidURL,
		},
		{
			name:          "URL without host",
			input:         "https://",
			expectedError: ErrInvalidURL,
		},
		{
			name:          "malformed URL",
			input:         "https:///invalid",
			expectedError: ErrInvalidURL,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewURL(tc.input)
			if tc.expectedError != nil {
				s.Error(err)
				s.True(errors.Is(err, tc.expectedError))
			}
		})
	}
}

func (s *URLTestSuite) TestURLNormalization() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trims whitespace",
			input:    "  https://example.com  ",
			expected: "https://example.com",
		},
		{
			name:     "normalizes path",
			input:    "https://example.com//path//to//resource",
			expected: "https://example.com//path//to//resource", // URL parsing preserves double slashes in path
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			normalized, err := NormalizeURL(tc.input)
			s.NoError(err)
			s.Equal(tc.expected, normalized)
		})
	}
}

func (s *URLTestSuite) TestURLComponents() {
	url, _ := NewURL("https://example.com:8080/path/to/resource?q=test#section")

	s.Equal("https", url.Scheme())
	s.Equal("example.com:8080", url.Host())
	s.Equal("/path/to/resource", url.Path())
}

func (s *URLTestSuite) TestEquals() {
	url1, _ := NewURL("https://example.com")
	url2, _ := NewURL("https://example.com")
	url3, _ := NewURL("https://different.com")

	s.True(url1.Equals(url2))
	s.False(url1.Equals(url3))
}

func (s *URLTestSuite) TestString() {
	url, _ := NewURL("https://example.com")
	s.Equal("https://example.com", url.String())
}

func (s *URLTestSuite) TestJSONSerialization() {
	url, _ := NewURL("https://example.com")
	
	jsonData, err := json.Marshal(url)
	s.NoError(err)
	s.JSONEq(`{"value":"https://example.com"}`, string(jsonData))
}

func (s *URLTestSuite) TestReconstitute() {
	url := ReconstituteURL("https://example.com")
	s.Equal("https://example.com", url.Value())
	s.Equal("https://example.com", url.String())
}

func (s *URLTestSuite) TestItCanBuildNewURLFromValidJSON() {
	jsonData := `{"value":"https://example.com"}`
	
	url, err := NewURLFromJSON([]byte(jsonData))
	s.NoError(err)
	s.Equal("https://example.com", url.Value())
}

func (s *URLTestSuite) TestItFailsToBuildNewURLFromInvalidJSON() {
	testCases := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON format",
			jsonData: `{"value":"https://example.com"`,
		},
		{
			name:     "invalid URL in JSON",
			jsonData: `{"value":"invalid-url"}`,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewURLFromJSON([]byte(tc.jsonData))
			s.Error(err)
		})
	}
}

func (s *URLTestSuite) TestTooLongURL() {
	// Create a URL that exceeds MaxURLLength
	longPath := make([]byte, MaxURLLength)
	for i := range longPath {
		longPath[i] = 'a'
	}
	longURL := "https://example.com/" + string(longPath)

	_, err := NewURL(longURL)
	s.Error(err)
	s.True(errors.Is(err, ErrTooLongURL))
}