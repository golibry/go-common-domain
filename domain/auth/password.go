package auth

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golibry/go-common-domain/domain"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
	BcryptCost        = 12 // Higher cost for better security
)

var (
	ErrPasswordTooShort = domain.NewError(
		"password must be at least %d characters long",
		MinPasswordLength,
	)
	ErrPasswordTooLong = domain.NewError(
		"password cannot exceed %d characters",
		MaxPasswordLength,
	)
	ErrPasswordTooWeak = domain.NewError(
		"password must contain at least one uppercase letter," +
			" one lowercase letter, one number, and one special character",
	)
	ErrPasswordCommon = domain.NewError(
		"password is too common or weak. " +
			"Try to not use common names or repeating characters like \"123456\" or \"123456789\". ",
	)
	ErrInvalidPasswordChars = domain.NewError(
		"password contains invalid characters; only letters, numbers, " +
			"and standard symbols are allowed",
	)
	ErrPasswordVerifyFailed = domain.NewError("failed to verify password")
)

// Password represents a secure password value object
type Password struct {
	hashedValue string
}

// NewPassword creates a new Password instance with validation and secure hashing
func NewPassword(plaintext string) (Password, error) {
	if err := ValidatePassword(plaintext); err != nil {
		return Password{}, err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), BcryptCost)
	if err != nil {
		return Password{}, err
	}

	return Password{
		hashedValue: string(hashedBytes),
	}, nil
}

// ReconstitutePassword creates a Password instance from a pre-hashed value without validation
// This is used when loading passwords from storage
func ReconstitutePassword(hashedValue string) Password {
	return Password{
		hashedValue: hashedValue,
	}
}

// Verify checks if the provided plaintext password matches the stored hash
func (p Password) Verify(plaintext string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plaintext))
	if err == nil {
		return nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrPasswordVerifyFailed
	}
	return domain.NewErrorWithWrap(err, "failed to verify password")
}

// HashedValue returns the hashed password value
func (p Password) HashedValue() string {
	return p.hashedValue
}

// Equals compares two Password objects for equality
func (p Password) Equals(other Password) bool {
	return p.hashedValue == other.hashedValue
}

// String returns a protected string representation of the password
func (p Password) String() string {
	return "[PROTECTED]"
}

// ValidatePassword validates a plaintext password against OWASP security standards
func ValidatePassword(password string) error {
	// Check length constraints
	if utf8.RuneCountInString(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	if utf8.RuneCountInString(password) > MaxPasswordLength {
		return ErrPasswordTooLong
	}

	// Check for invalid characters (only printable characters are allowed)
	for _, r := range password {
		if !unicode.IsPrint(r) {
			return ErrInvalidPasswordChars
		}
	}

	// Check password complexity requirements first
	if err := validatePasswordComplexity(password); err != nil {
		return err
	}

	// Check against common passwords after complexity
	if err := validatePasswordStrength(password); err != nil {
		return err
	}

	return nil
}

// validatePasswordComplexity ensures password meets complexity requirements
func validatePasswordComplexity(password string) error {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsNumber(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrPasswordTooWeak
	}

	return nil
}

// validatePasswordStrength checks against common weak passwords
func validatePasswordStrength(password string) error {
	// Convert to lowercase for comparison
	lowerPassword := strings.ToLower(password)

	// Common weak passwords and patterns
	commonPasswords := []string{
		"password", "123456", "123456789", "12345678", "12345",
		"1234567", "password123", "admin", "qwerty", "abc123",
		"letmein", "monkey", "1234567890", "dragon", "111111",
		"baseball", "iloveyou", "trustno1", "sunshine", "master",
		"welcome", "shadow", "ashley", "football", "jesus",
		"michael", "ninja", "mustang", "password1",
	}

	for _, common := range commonPasswords {
		if lowerPassword == common {
			return ErrPasswordCommon
		}
	}

	// Check for simple patterns
	if isSequentialPattern(password) || isRepeatingPattern(password) {
		return ErrPasswordCommon
	}

	return nil
}

// isSequentialPattern checks for sequential characters like "123456" or "abcdef"
func isSequentialPattern(password string) bool {
	// Build rune slice
	runes := []rune(password)
	if len(runes) < 4 {
		return false
	}

	// helper to check ranges
	isDigit := func(r rune) bool { return r >= '0' && r <= '9' }
	isLetter := func(r rune) bool {
		lr := unicode.ToLower(r)
		return lr >= 'a' && lr <= 'z'
	}

	// Sliding window of 4 runes for ascending/descending sequences
	for i := 0; i <= len(runes)-4; i++ {
		a, b, c, d := runes[i], runes[i+1], runes[i+2], runes[i+3]

		// numeric ascending
		if isDigit(a) && isDigit(b) && isDigit(c) && isDigit(d) {
			if b == a+1 && c == b+1 && d == c+1 {
				return true
			}
			if b == a-1 && c == b-1 && d == c-1 {
				return true
			}
		}

		// alphabetic sequences (case-insensitive)
		la, lb, lc, ld := unicode.ToLower(a), unicode.ToLower(b),
			unicode.ToLower(c), unicode.ToLower(d)
		if isLetter(la) && isLetter(lb) && isLetter(lc) && isLetter(ld) {
			if lb == la+1 && lc == lb+1 && ld == lc+1 {
				return true
			}
			if lb == la-1 && lc == lb-1 && ld == lc-1 {
				return true
			}
		}
	}
	return false
}

// isRepeatingPattern checks for repeating characters like "aaaa" or "1111"
func isRepeatingPattern(password string) bool {
	runes := []rune(password)
	if len(runes) < 4 {
		return false
	}
	count := 1
	for i := 1; i < len(runes); i++ {
		if runes[i] == runes[i-1] {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}
