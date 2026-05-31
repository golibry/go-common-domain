package finance

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrEmptyCurrency   = domain.NewError("currency cannot be empty")
	ErrInvalidCurrency = domain.NewError("currency must be exactly 3 letters")
)

var currencyRegex = regexp.MustCompile(`^[A-Z]{3}$`)

var currencyMinorUnitScales = map[string]int32{
	"AED": 2,
	"AFN": 2,
	"ALL": 2,
	"AMD": 2,
	"ANG": 2,
	"AOA": 2,
	"ARS": 2,
	"AUD": 2,
	"AWG": 2,
	"AZN": 2,
	"BAM": 2,
	"BBD": 2,
	"BDT": 2,
	"BGN": 2,
	"BHD": 3,
	"BIF": 0,
	"BMD": 2,
	"BND": 2,
	"BOB": 2,
	"BRL": 2,
	"BSD": 2,
	"BTN": 2,
	"BWP": 2,
	"BYN": 2,
	"BZD": 2,
	"CAD": 2,
	"CDF": 2,
	"CHF": 2,
	"CLP": 0,
	"CNY": 2,
	"COP": 2,
	"CRC": 2,
	"CUP": 2,
	"CVE": 2,
	"CZK": 2,
	"DJF": 0,
	"DKK": 2,
	"DOP": 2,
	"DZD": 2,
	"EGP": 2,
	"ERN": 2,
	"ETB": 2,
	"EUR": 2,
	"FJD": 2,
	"FKP": 2,
	"GBP": 2,
	"GEL": 2,
	"GHS": 2,
	"GIP": 2,
	"GMD": 2,
	"GNF": 0,
	"GTQ": 2,
	"GYD": 2,
	"HKD": 2,
	"HNL": 2,
	"HTG": 2,
	"HUF": 2,
	"IDR": 2,
	"ILS": 2,
	"INR": 2,
	"IQD": 3,
	"IRR": 2,
	"ISK": 0,
	"JMD": 2,
	"JOD": 3,
	"JPY": 0,
	"KES": 2,
	"KGS": 2,
	"KHR": 2,
	"KMF": 0,
	"KPW": 2,
	"KRW": 0,
	"KWD": 3,
	"KYD": 2,
	"KZT": 2,
	"LAK": 2,
	"LBP": 2,
	"LKR": 2,
	"LRD": 2,
	"LSL": 2,
	"LYD": 3,
	"MAD": 2,
	"MDL": 2,
	"MGA": 2,
	"MKD": 2,
	"MMK": 2,
	"MNT": 2,
	"MOP": 2,
	"MRU": 2,
	"MUR": 2,
	"MVR": 2,
	"MWK": 2,
	"MXN": 2,
	"MYR": 2,
	"MZN": 2,
	"NAD": 2,
	"NGN": 2,
	"NIO": 2,
	"NOK": 2,
	"NPR": 2,
	"NZD": 2,
	"OMR": 3,
	"PAB": 2,
	"PEN": 2,
	"PGK": 2,
	"PHP": 2,
	"PKR": 2,
	"PLN": 2,
	"PYG": 0,
	"QAR": 2,
	"RON": 2,
	"RSD": 2,
	"RUB": 2,
	"RWF": 0,
	"SAR": 2,
	"SBD": 2,
	"SCR": 2,
	"SDG": 2,
	"SEK": 2,
	"SGD": 2,
	"SHP": 2,
	"SLE": 2,
	"SOS": 2,
	"SRD": 2,
	"SSP": 2,
	"STN": 2,
	"SYP": 2,
	"SZL": 2,
	"THB": 2,
	"TJS": 2,
	"TMT": 2,
	"TND": 3,
	"TOP": 2,
	"TRY": 2,
	"TTD": 2,
	"TWD": 2,
	"TZS": 2,
	"UAH": 2,
	"UGX": 0,
	"USD": 2,
	"UYU": 2,
	"UZS": 2,
	"VES": 2,
	"VND": 0,
	"VUV": 0,
	"WST": 2,
	"XAF": 0,
	"XCD": 2,
	"XOF": 0,
	"XPF": 0,
	"YER": 2,
	"ZAR": 2,
	"ZMW": 2,
	"ZWL": 2,
}

type Currency struct {
	value string
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

// MinorUnitScale returns the number of decimal places for this currency's minor unit.
func (c Currency) MinorUnitScale() (int32, bool) {
	scale, ok := currencyMinorUnitScales[c.value]
	return scale, ok
}

// MarshalText returns the normalized currency code as text.
func (c Currency) MarshalText() ([]byte, error) {
	return []byte(c.value), nil
}

// UnmarshalText validates and normalizes text into a Currency.
func (c *Currency) UnmarshalText(text []byte) error {
	currency, err := NewCurrency(string(text))
	if err != nil {
		return err
	}

	*c = currency
	return nil
}

// UnmarshalJSON validates and normalizes a JSON currency string.
func (c *Currency) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return ErrInvalidCurrency
	}

	return c.UnmarshalText([]byte(value))
}

// Scan validates and normalizes a database value into a Currency.
func (c *Currency) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return c.UnmarshalText([]byte(v))
	case []byte:
		return c.UnmarshalText(v)
	default:
		return ErrInvalidCurrency
	}
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
