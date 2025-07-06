package person

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FullNameTestSuite struct {
	suite.Suite
}

func TestFullNameSuite(t *testing.T) {
	suite.Run(t, new(FullNameTestSuite))
}

func (s *FullNameTestSuite) TestItCanBuildNewFullNameWithValidParts() {
	testCases := []struct {
		name       string
		firstName  string
		middleName string
		lastName   string
	}{
		{"Standard name", "John", "William", "Doe"},
		{"Empty middle name", "John", "", "Doe"},
		{"Compound first name", "Jean-Claude", "Van", "Damme"},
		{"Compound last name", "Mary", "Jane", "O'Connor"},
		{"Hyphenated last name", "Sarah", "Jessica", "Parker-Davis"},
		{"Name with periods", "J.R", "", "Tolkien"},
		{"Name with normalization", "  John  ", "  William  ", "  Doe  "},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				fullName, err := NewFullName(tc.firstName, tc.middleName, tc.lastName)

				s.NoError(err)
				s.NotNil(fullName)

				// For names with spaces, we expect them to be trimmed
				expectedFirstName := strings.TrimSpace(tc.firstName)
				expectedMiddleName := strings.TrimSpace(tc.middleName)
				expectedLastName := strings.TrimSpace(tc.lastName)

				s.Equal(expectedFirstName, fullName.FirstName())
				s.Equal(expectedMiddleName, fullName.MiddleName())
				s.Equal(expectedLastName, fullName.LastName())
			},
		)
	}
}

func (s *FullNameTestSuite) TestItFailsToBuildNewFullNameFromInvalidFirstName() {
	testCases := []struct {
		name     string
		input    string
		expected error
	}{
		{"Empty first name", "", ErrEmptyNamePart},
		{"Invalid characters", "John123", ErrInvalidNamePartChars},
		{"Too long name", strings.Repeat("A", MaxNamePartLength+1), ErrTooLongNamePart},
		{"Starts with hyphen", "-John", ErrInvalidNamePartChars},
		{"Ends with apostrophe", "John'", ErrInvalidNamePartChars},
		{"Period at the end", "J.R.R.", ErrInvalidNamePartChars},
		{"Period at the start", ".J.R.R", ErrInvalidNamePartChars},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewFullName(tc.input, "William", "Doe")
				s.Error(err)
				s.True(
					errors.Is(err, tc.expected),
					"Expected error containing %v, got %v",
					tc.expected,
					err,
				)
			},
		)
	}
}

func (s *FullNameTestSuite) TestItFailsToBuildNewFullNameFromInvalidLastName() {
	testCases := []struct {
		name     string
		input    string
		expected error
	}{
		{"Empty last name", "", ErrEmptyNamePart},
		{"Invalid characters", "Doe123", ErrInvalidNamePartChars},
		{"Too long name", strings.Repeat("A", MaxNamePartLength+1), ErrTooLongNamePart},
		{"Starts with period", ".Doe", ErrInvalidNamePartChars},
		{"Ends with hyphen", "Doe-", ErrInvalidNamePartChars},
		{"Period at the end", "J.R.R.", ErrInvalidNamePartChars},
		{"Period at the start", ".J.R.R", ErrInvalidNamePartChars},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewFullName("John", "William", tc.input)
				s.Error(err)
				s.True(
					errors.Is(err, tc.expected),
					"Expected error containing %v, got %v",
					tc.expected,
					err,
				)
			},
		)
	}
}

// TestNameNormalization tests the normalization of name parts
func (s *FullNameTestSuite) TestNameNormalization() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Trim spaces", "  John  ", "John"},
		{"Collapse multiple spaces", "John  Doe", "John Doe"},
		{"Collapse multiple hyphens", "Smith--Jones", "Smith-Jones"},
		{"Collapse multiple apostrophes", "O''Brien", "O'Brien"},
		{"Collapse multiple periods", "J..R..R", "J.R.R"},
		{"Mixed special characters", "Smith-  -Jones", "Smith-Jones"},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				normalized, err := NormalizeNamePart(tc.input)
				s.NoError(err)
				s.Equal(tc.expected, normalized)
			},
		)
	}
}

// TestEquals tests the Equals method
func (s *FullNameTestSuite) TestEquals() {
	name1, _ := NewFullName("John", "William", "Doe")

	// Same values
	name2, _ := NewFullName("John", "William", "Doe")
	s.True(name1.Equals(*name2))

	// Different values
	name3, _ := NewFullName("Jane", "William", "Doe")
	s.False(name1.Equals(*name3))
}

// TestString tests the String method
func (s *FullNameTestSuite) TestString() {
	name1, _ := NewFullName("John", "William", "Doe")
	s.Equal("John William Doe", name1.String())

	name2, _ := NewFullName("John", "", "Doe")
	s.Equal("John Doe", name2.String())
}

// TestJSONSerialization tests the JSON marshaling and unmarshalling
func (s *FullNameTestSuite) TestJSONSerialization() {
	name, _ := NewFullName("John", "William", "Doe")

	jsonData, _ := json.Marshal(name)
	var unmarshalledName FullName
	_ = json.Unmarshal(jsonData, &unmarshalledName)

	s.True(name.Equals(unmarshalledName))
}

// TestReconstitute tests the Reconstitute function
func (s *FullNameTestSuite) TestReconstitute() {
	firstname := "John"
	lastName := "William"
	middleName := "Doe"
	fullName := Reconstitute(firstname, middleName, lastName)

	s.NotNil(fullName)
	s.Equal(firstname, fullName.FirstName())
	s.Equal(middleName, fullName.MiddleName())
	s.Equal(lastName, fullName.LastName())
}
