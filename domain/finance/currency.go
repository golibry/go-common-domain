package finance

import (
	"encoding/json"
	"github.com/golibry/go-common-domain/domain"
	"regexp"
	"strings"
)

var (
	ErrEmptyCurrency   = domain.NewError("currency cannot be empty")
	ErrInvalidCurrency = domain.NewError("currency must be exactly 3 letters")
)

var currencyRegex = regexp.MustCompile(`^[A-Z]{3}$`)

type Currency struct {
	value string
}

type currencyJSON struct {
	Value string `json:"value"`
}

// NewCurrency creates a new instance of Currency with validation and normalization
func NewCurrency(value string) (Currency, error) {
	normalized, err := NormalizeCurrency(value)
	if err != nil {
		return Currency{}, err
	}

	return Currency{
		value: normalized,
	}, nil
}

// ReconstituteCurrency creates a new Currency instance without validation or normalization
func ReconstituteCurrency(value string) Currency {
	return Currency{
		value: value,
	}
}

// NewCurrencyFromJSON creates Currency from JSON bytes array
func NewCurrencyFromJSON(data []byte) (Currency, error) {
    var temp currencyJSON

    if err := json.Unmarshal(data, &temp); err != nil {
        return Currency{}, domain.NewErrorWithWrap(err, "failed to build currency from json")
    }

	newCurrency, err := NewCurrency(temp.Value)
	if err != nil {
		return Currency{}, err
	}

	return newCurrency, nil
}

// Value returns the currency value
func (c Currency) Value() string {
	return c.value
}

// Equals compares two Currency objects for equality
func (c Currency) Equals(other Currency) bool {
	return c.value == other.value
}

// String returns a string representation of the currency
func (c Currency) String() string {
	return c.value
}

// MarshalJSON implements json.Marshaler
func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		currencyJSON{
			Value: c.value,
		},
	)
}

// NormalizeCurrency normalizes a currency by trimming spaces and converting to uppercase
func NormalizeCurrency(currency string) (string, error) {
	// Trim spaces and convert to uppercase
	normalized := strings.ToUpper(strings.TrimSpace(currency))
	
	if err := IsValidCurrency(normalized); err != nil {
		return "", err
	}
	
	return normalized, nil
}

// IsValidCurrency validates a currency (must be exactly 3 uppercase letters)
func IsValidCurrency(currency string) error {
	if currency == "" {
		return ErrEmptyCurrency
	}
	
	if !currencyRegex.MatchString(currency) {
		return ErrInvalidCurrency
	}
	
	return nil
}