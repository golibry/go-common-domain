package auth

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PasswordTestSuite struct {
	suite.Suite
}

func TestPasswordSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}

func (s *PasswordTestSuite) TestItCanCreateNewPasswordWithValidInput() {
	testCases := []struct {
		name     string
		password string
	}{
		{"minimum valid password", "Abc123!@"},
		{"recommended length password", "MySecure123!"},
		{"password with various special chars", "Test123!@#$%"},
		{"longer secure password", "MyVerySecurePassword123!@#"},
		{"password with unicode", "Tëst123!@"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			password, err := NewPassword(tc.password)
			s.NoError(err)
			s.NotEmpty(password.HashedValue())
			s.Equal("[PROTECTED]", password.String())
		})
	}
}

func (s *PasswordTestSuite) TestItFailsToCreatePasswordThatIsTooShort() {
	testCases := []struct {
		name     string
		password string
	}{
		{"empty password", ""},
		{"one character", "A"},
		{"seven characters", "Abc123!"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewPassword(tc.password)
			s.Error(err)
			s.True(errors.Is(err, ErrPasswordTooShort))
		})
	}
}

func (s *PasswordTestSuite) TestItFailsToCreatePasswordThatIsTooLong() {
	longPassword := strings.Repeat("A", MaxPasswordLength+1) + "bc123!@"
	_, err := NewPassword(longPassword)
	s.Error(err)
	s.True(errors.Is(err, ErrPasswordTooLong))
}

func (s *PasswordTestSuite) TestItFailsToCreatePasswordThatIsTooWeak() {
	testCases := []struct {
		name     string
		password string
		reason   string
	}{
		{"no uppercase", "abc123!@", "missing uppercase"},
		{"no lowercase", "ABC123!@", "missing lowercase"},
		{"no numbers", "Abcdef!@", "missing numbers"},
		{"no special chars", "Abc12345", "missing special characters"},
		{"only letters", "AbcDefGh", "missing numbers and special chars"},
		{"only numbers", "12345678", "missing letters and special chars"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewPassword(tc.password)
			s.Error(err)
			s.True(errors.Is(err, ErrPasswordTooWeak))
		})
	}
}

func (s *PasswordTestSuite) TestItFailsToCreateCommonPasswords() {
	testCases := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{"password", "password", ErrPasswordTooWeak}, // 8 chars, fails complexity (no numbers/special chars)
		{"Password123", "Password123", ErrPasswordTooWeak}, // 11 chars, fails complexity (no special chars)
		{"123456", "123456", ErrPasswordTooShort}, // 6 chars, fails length first
		{"qwerty", "qwerty", ErrPasswordTooShort}, // 6 chars, fails length first
		{"admin", "admin", ErrPasswordTooShort},   // 5 chars, fails length first
		{"letmein", "letmein", ErrPasswordTooShort}, // 7 chars, fails length first
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewPassword(tc.password)
			s.Error(err)
			s.True(errors.Is(err, tc.expectedErr))
		})
	}
}

func (s *PasswordTestSuite) TestItFailsToCreatePasswordsWithSequentialPatterns() {
	testCases := []struct {
		name     string
		password string
	}{
		{"numeric sequence", "Test1234!"},
		{"alphabetic sequence", "Abcdef1!"},
		{"reverse sequence", "Test4321!"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewPassword(tc.password)
			s.Error(err)
			s.True(errors.Is(err, ErrPasswordCommon))
		})
	}
}

func (s *PasswordTestSuite) TestItFailsToCreatePasswordsWithRepeatingPatterns() {
	testCases := []struct {
		name     string
		password string
	}{
		{"repeating numbers", "Test1111!"},
		{"repeating letters", "Aaaaa123!"},
		{"repeating special chars", "Test123!!!!"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewPassword(tc.password)
			s.Error(err)
			s.True(errors.Is(err, ErrPasswordCommon))
		})
	}
}

func (s *PasswordTestSuite) TestItFailsToCreatePasswordWithInvalidCharacters() {
	// Password with non-printable characters
	invalidPassword := "Test123!\x00\x01"
	_, err := NewPassword(invalidPassword)
	s.Error(err)
	s.True(errors.Is(err, ErrInvalidPasswordChars))
}

func (s *PasswordTestSuite) TestPasswordVerification() {
	plaintext := "MySecure123!@"
	password, err := NewPassword(plaintext)
	s.NoError(err)

	// Test correct password verification
	err = password.Verify(plaintext)
	s.NoError(err)

	// Test incorrect password verification
	err = password.Verify("WrongPassword123!")
	s.Error(err)
	s.True(errors.Is(err, ErrPasswordVerifyFailed))

	// Test empty password verification
	err = password.Verify("")
	s.Error(err)
	s.True(errors.Is(err, ErrPasswordVerifyFailed))
}

func (s *PasswordTestSuite) TestPasswordEquals() {
	plaintext := "MySecure123!@"
	password1, err := NewPassword(plaintext)
	s.NoError(err)

	password2, err := NewPassword(plaintext)
	s.NoError(err)

	// Different passwords with same plaintext should have different hashes (due to salt)
	s.False(password1.Equals(password2))

	// Same password object should equal itself
	s.True(password1.Equals(password1))

	// Reconstituted password with same hash should be equal
	reconstituted := ReconstitutePassword(password1.HashedValue())
	s.True(password1.Equals(reconstituted))
}

func (s *PasswordTestSuite) TestPasswordString() {
	password, err := NewPassword("MySecure123!@")
	s.NoError(err)
	s.Equal("[PROTECTED]", password.String())
}

func (s *PasswordTestSuite) TestJSONSerialization() {
	plaintext := "MySecure123!@"
	password, err := NewPassword(plaintext)
	s.NoError(err)

	// Test marshaling
	jsonData, err := json.Marshal(password)
	s.NoError(err)
	s.Contains(string(jsonData), "hashedValue")
	s.NotContains(string(jsonData), plaintext)

	// Test unmarshaling
	var unmarshaled Password
	err = json.Unmarshal(jsonData, &unmarshaled)
	s.NoError(err)
	s.True(password.Equals(unmarshaled))

	// Verify the unmarshaled password works
	err = unmarshaled.Verify(plaintext)
	s.NoError(err)
}

func (s *PasswordTestSuite) TestNewPasswordFromJSON() {
	plaintext := "MySecure123!@"
	originalPassword, err := NewPassword(plaintext)
	s.NoError(err)

	// Create JSON data
	jsonData, err := json.Marshal(originalPassword)
	s.NoError(err)

	// Test NewPasswordFromJSON
	password, err := NewPasswordFromJSON(jsonData)
	s.NoError(err)
	s.True(originalPassword.Equals(password))

	// Verify the password works
	err = password.Verify(plaintext)
	s.NoError(err)
}

func (s *PasswordTestSuite) TestItFailsToBuildPasswordFromInvalidJSON() {
	invalidJSON := []byte(`{"invalid": "json"}`)
	_, err := NewPasswordFromJSON(invalidJSON)
	s.Error(err)
	s.Contains(err.Error(), "failed to build password from json")
}

func (s *PasswordTestSuite) TestReconstitutePassword() {
	plaintext := "MySecure123!@"
	originalPassword, err := NewPassword(plaintext)
	s.NoError(err)

	// Test reconstitution
	reconstituted := ReconstitutePassword(originalPassword.HashedValue())
	s.True(originalPassword.Equals(reconstituted))
	s.Equal(originalPassword.HashedValue(), reconstituted.HashedValue())

	// Verify the reconstituted password works
	err = reconstituted.Verify(plaintext)
	s.NoError(err)
}

func (s *PasswordTestSuite) TestPasswordValidation() {
	// Test valid passwords
	validPasswords := []string{
		"MySecure123!@",
		"Strong123!@",
		"Test@123Pass",
		"Str0ng!P@ssw0rd",
	}

	for _, pwd := range validPasswords {
		s.Run("valid: "+pwd, func() {
			err := ValidatePassword(pwd)
			s.NoError(err)
		})
	}

	// Test invalid passwords
	invalidPasswords := []struct {
		password string
		expected error
	}{
		{"short", ErrPasswordTooShort},
		{"NoNumbers!@", ErrPasswordTooWeak},
		{"nonumbers123", ErrPasswordTooWeak},
		{"NOLOWERCASE123!", ErrPasswordTooWeak},
		{"NoSpecialChars123", ErrPasswordTooWeak},
		{"password", ErrPasswordTooWeak}, // fails complexity first
		{"Test1234!", ErrPasswordCommon},
		{"Test1111!", ErrPasswordCommon},
	}

	for _, tc := range invalidPasswords {
		s.Run("invalid: "+tc.password, func() {
			err := ValidatePassword(tc.password)
			s.Error(err)
			s.True(errors.Is(err, tc.expected))
		})
	}
}

func (s *PasswordTestSuite) TestPasswordHashingConsistency() {
	plaintext := "MySecure123!@"
	
	// Create multiple passwords with same plaintext
	password1, err := NewPassword(plaintext)
	s.NoError(err)
	
	password2, err := NewPassword(plaintext)
	s.NoError(err)

	// Hashes should be different (bcrypt uses random salt)
	s.NotEqual(password1.HashedValue(), password2.HashedValue())

	// But both should verify correctly
	s.NoError(password1.Verify(plaintext))
	s.NoError(password2.Verify(plaintext))
}

func (s *PasswordTestSuite) TestPasswordSecurityProperties() {
	plaintext := "MySecure123!@"
	password, err := NewPassword(plaintext)
	s.NoError(err)

	// Hash should not contain plaintext
	s.NotContains(password.HashedValue(), plaintext)
	s.NotContains(password.HashedValue(), "MySecure")
	s.NotContains(password.HashedValue(), "123")

	// Hash should be bcrypt format (starts with $2a$, $2b$, or $2y$)
	hash := password.HashedValue()
	s.True(strings.HasPrefix(hash, "$2a$") || 
		   strings.HasPrefix(hash, "$2b$") || 
		   strings.HasPrefix(hash, "$2y$"))

	// Hash should contain cost factor
	s.Contains(hash, "$12$") // BcryptCost = 12
}