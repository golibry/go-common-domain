package auth

import (
    "errors"
    "encoding/json"
    "github.com/golibry/go-common-domain/domain"
    "golang.org/x/crypto/bcrypt"
    "strings"
    "unicode"
    "unicode/utf8"
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
	ErrInvalidPasswordChars = domain.NewError("password contains invalid characters")
	ErrPasswordVerifyFailed = domain.NewError("failed to verify password")
)

// Password represents a secure password value object
type Password struct {
	hashedValue string
}

type passwordJSON struct {
    HashedValue string `json:"hashedValue"`
}

// NewPassword creates a new Password instance with validation and secure hashing
func NewPassword(plaintext string) (Password, error) {
	if err := ValidatePassword(plaintext); err != nil {
		return Password{}, err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), BcryptCost)
	if err != nil {
		return Password{}, domain.NewErrorWithWrap(err, "failed to hash password")
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

// NewPasswordFromJSON creates Password from JSON bytes array
func NewPasswordFromJSON(data []byte) (Password, error) {
    var temp passwordJSON

    if err := json.Unmarshal(data, &temp); err != nil {
        return Password{}, domain.NewErrorWithWrap(err, "failed to build password from json")
    }

	if temp.HashedValue == "" {
		return Password{}, domain.NewError(
			"failed to build password from json: missing or empty hashedValue",
		)
	}

	return ReconstitutePassword(temp.HashedValue), nil
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

// String returns a masked representation of the password for logging/debugging
func (p Password) String() string {
	return "[PROTECTED]"
}

// MarshalJSON implements json.Marshaler
func (p Password) MarshalJSON() ([]byte, error) {
    return json.Marshal(
        passwordJSON{
            HashedValue: p.hashedValue,
        },
    )
}

// UnmarshalJSON implements json.Unmarshaler
func (p *Password) UnmarshalJSON(data []byte) error {
    var temp passwordJSON
    if err := json.Unmarshal(data, &temp); err != nil {
        return domain.NewErrorWithWrap(err, "failed to unmarshal password from json")
    }
    if temp.HashedValue == "" {
        return domain.NewError("failed to unmarshal password from json: missing or empty hashedValue")
    }
    p.hashedValue = temp.HashedValue
    return nil
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

	// Check for invalid characters (only printable ASCII and common Unicode)
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
        la, lb, lc, ld := unicode.ToLower(a), unicode.ToLower(b), unicode.ToLower(c), unicode.ToLower(d)
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
