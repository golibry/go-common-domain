package finance

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/golibry/go-common-domain/domain"
	"github.com/shopspring/decimal"
)

const (
	DefaultMoneyScale int32 = 2
	MaxMoneyScale     int32 = 18
)

var (
	ErrNegativeAmount              = domain.NewError("money amount cannot be negative")
	ErrInvalidMoneyScale           = domain.NewError("money scale must be between 0 and 18")
	ErrInvalidMoneyAmountPrecision = domain.NewError("money amount has more precision than scale allows")
	ErrMoneyAmountTooLarge         = domain.NewError("money amount is too large")
	ErrMissingMoneyAmount          = domain.NewError("money amount is required")
	ErrMissingMoneyCurrency        = domain.NewError("money currency is required")
)

type Money struct {
	amountMinor int64
	currency    Currency
	scale       int32
}

type moneyJSON struct {
	Amount   string   `json:"amount"`
	Currency Currency `json:"currency"`
	Scale    int32    `json:"scale"`
}

// NewMoney creates a new instance of Money with validation
func NewMoney(amount decimal.Decimal, currency Currency) (Money, error) {
	return NewMoneyWithScale(amount, currency, DefaultMoneyScale)
}

// NewMoneyWithScale creates a new instance of Money with validation and explicit minor-unit scale.
func NewMoneyWithScale(amount decimal.Decimal, currency Currency, scale int32) (Money, error) {
	if err := IsValidMoneyAmount(amount); err != nil {
		return Money{}, err
	}

	amountMinor, err := decimalToMinorUnits(amount, scale)
	if err != nil {
		return Money{}, err
	}

	return NewMoneyFromMinorUnits(amountMinor, currency, scale)
}

// NewMoneyFromMinorUnits creates a new instance of Money from integer minor units.
func NewMoneyFromMinorUnits(amountMinor int64, currency Currency, scale int32) (Money, error) {
	if err := IsValidMoneyScale(scale); err != nil {
		return Money{}, err
	}

	if amountMinor < 0 {
		return Money{}, ErrNegativeAmount
	}

	return Money{
		amountMinor: amountMinor,
		currency:    currency,
		scale:       scale,
	}, nil
}

// NewMoneyFromString creates a new instance of Money from string amount and currency
func NewMoneyFromString(amountStr, currencyStr string) (Money, error) {
	return NewMoneyFromStringWithScale(amountStr, currencyStr, DefaultMoneyScale)
}

// NewMoneyFromStringWithScale creates a new instance of Money from string amount, currency, and explicit scale.
func NewMoneyFromStringWithScale(amountStr, currencyStr string, scale int32) (Money, error) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return Money{}, domain.NewErrorWithWrap(err, "invalid amount format")
	}

	currency, err := NewCurrency(currencyStr)
	if err != nil {
		return Money{}, err
	}

	return NewMoneyWithScale(amount, currency, scale)
}

// ReconstituteMoney creates a new Money instance without validation.
// Amounts are rounded half-up to the default scale before being stored as minor units.
func ReconstituteMoney(amount decimal.Decimal, currency Currency) Money {
	amount = amount.Round(DefaultMoneyScale)
	amountMinor, _ := decimalToMinorUnits(amount, DefaultMoneyScale)

	return Money{
		amountMinor: amountMinor,
		currency:    currency,
		scale:       DefaultMoneyScale,
	}
}

// ReconstituteMoneyFromMinorUnits creates a new Money instance from raw fields without validation.
func ReconstituteMoneyFromMinorUnits(amountMinor int64, currency Currency, scale int32) Money {
	return Money{
		amountMinor: amountMinor,
		currency:    currency,
		scale:       scale,
	}
}

// Amount returns the money amount
func (m Money) Amount() decimal.Decimal {
	return decimal.New(m.amountMinor, -m.scale)
}

// AmountMinorUnits returns the canonical integer minor-unit amount.
func (m Money) AmountMinorUnits() int64 {
	return m.amountMinor
}

// Currency returns the money currency
func (m Money) Currency() Currency {
	return m.currency
}

// Scale returns the number of decimal places represented by one major unit.
func (m Money) Scale() int32 {
	return m.scale
}

// Equals compares two Money objects for equality
func (m Money) Equals(other Money) bool {
	return m.Amount().Equal(other.Amount()) && m.currency.Equals(other.currency)
}

// String returns a string representation of the money
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.Amount().String(), m.currency.String())
}

// MarshalJSON returns money as an object with a decimal-safe string amount.
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(moneyJSON{
		Amount:   m.Amount().StringFixed(m.scale),
		Currency: m.currency,
		Scale:    m.scale,
	})
}

// UnmarshalJSON validates and normalizes a JSON money object.
func (m *Money) UnmarshalJSON(data []byte) error {
	var raw struct {
		Amount   json.RawMessage `json:"amount"`
		Currency Currency        `json:"currency"`
		Scale    *int32          `json:"scale"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw.Amount) == 0 {
		return ErrMissingMoneyAmount
	}

	if raw.Currency.String() == "" {
		return ErrMissingMoneyCurrency
	}

	amount, err := decodeMoneyAmount(raw.Amount)
	if err != nil {
		return err
	}

	scale := DefaultMoneyScale
	if raw.Scale != nil {
		scale = *raw.Scale
	}

	money, err := NewMoneyFromStringWithScale(amount, raw.Currency.String(), scale)
	if err != nil {
		return err
	}

	*m = money
	return nil
}

// Add adds another Money object to this one (must have the same currency)
func (m Money) Add(other Money) (Money, error) {
	if !m.currency.Equals(other.currency) {
		return Money{}, domain.NewError(
			"cannot add money with different currencies: %s and %s",
			m.currency.String(),
			other.currency.String(),
		)
	}

	scale := maxScale(m.scale, other.scale)
	newAmount := m.Amount().Add(other.Amount())
	return NewMoneyWithScale(newAmount, m.currency, scale)
}

// Subtract subtracts another Money object from this one (must have same currency)
func (m Money) Subtract(other Money) (Money, error) {
	if !m.currency.Equals(other.currency) {
		return Money{}, domain.NewError(
			"cannot subtract money with different currencies: %s and %s",
			m.currency.String(),
			other.currency.String(),
		)
	}

	scale := maxScale(m.scale, other.scale)
	newAmount := m.Amount().Sub(other.Amount())
	if newAmount.IsNegative() {
		return Money{}, ErrNegativeAmount
	}

	return NewMoneyWithScale(newAmount, m.currency, scale)
}

// Multiply multiplies the money amount by a factor
func (m Money) Multiply(factor decimal.Decimal) (Money, error) {
	return m.MultiplyWithRounding(factor, RoundHalfUp)
}

// MultiplyWithRounding multiplies the money amount by a factor using the provided rounding mode.
func (m Money) MultiplyWithRounding(factor decimal.Decimal, mode RoundingMode) (Money, error) {
	newAmount := m.Amount().Mul(factor)
	if newAmount.IsNegative() {
		return Money{}, ErrNegativeAmount
	}

	return newMoneyWithRoundedAmount(newAmount, m.currency, m.scale, mode)
}

// Divide divides the money amount by a divisor
func (m Money) Divide(divisor decimal.Decimal) (Money, error) {
	return m.DivideWithRounding(divisor, RoundHalfUp)
}

// DivideWithRounding divides the money amount by a divisor using the provided rounding mode.
func (m Money) DivideWithRounding(divisor decimal.Decimal, mode RoundingMode) (Money, error) {
	if divisor.IsZero() {
		return Money{}, domain.NewError("cannot divide by zero")
	}

	newAmount := m.Amount().Div(divisor)
	if newAmount.IsNegative() {
		return Money{}, ErrNegativeAmount
	}

	return newMoneyWithRoundedAmount(newAmount, m.currency, m.scale, mode)
}

// IsValidMoneyAmount validates a money amount (must not be negative)
func IsValidMoneyAmount(amount decimal.Decimal) error {
	if amount.IsNegative() {
		return ErrNegativeAmount
	}
	return nil
}

// IsValidMoneyScale validates a money scale.
func IsValidMoneyScale(scale int32) error {
	if scale < 0 || scale > MaxMoneyScale {
		return ErrInvalidMoneyScale
	}
	return nil
}

func decodeMoneyAmount(data json.RawMessage) (string, error) {
	var amount string
	if err := json.Unmarshal(data, &amount); err == nil {
		return amount, nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err != nil {
		return "", domain.NewErrorWithWrap(err, "invalid amount format")
	}

	return number.String(), nil
}

func decimalToMinorUnits(amount decimal.Decimal, scale int32) (int64, error) {
	if err := IsValidMoneyScale(scale); err != nil {
		return 0, err
	}

	scaled := amount.Shift(scale)
	if !scaled.IsInteger() {
		return 0, ErrInvalidMoneyAmountPrecision
	}

	amountMinor := scaled.BigInt()
	if !amountMinor.IsInt64() {
		return 0, ErrMoneyAmountTooLarge
	}

	return amountMinor.Int64(), nil
}

func newMoneyWithRoundedAmount(amount decimal.Decimal, currency Currency, scale int32, mode RoundingMode) (Money, error) {
	rounded, err := roundMoneyAmount(amount, scale, mode)
	if err != nil {
		return Money{}, err
	}

	return NewMoneyWithScale(rounded, currency, scale)
}

func roundMoneyAmount(amount decimal.Decimal, scale int32, mode RoundingMode) (decimal.Decimal, error) {
	if err := IsValidMoneyScale(scale); err != nil {
		return decimal.Zero, err
	}

	switch mode {
	case RoundHalfUp:
		return amount.Round(scale), nil
	case RoundDown:
		return amount.RoundDown(scale), nil
	case RoundUp:
		return amount.RoundUp(scale), nil
	default:
		return decimal.Zero, ErrInvalidRoundingMode
	}
}

func maxScale(first, second int32) int32 {
	if first > second {
		return first
	}

	return second
}

func checkedMoneyMinorAmount(amount *big.Int) (int64, error) {
	if !amount.IsInt64() {
		return 0, ErrMoneyAmountTooLarge
	}

	if amount.Sign() < 0 {
		return 0, ErrNegativeAmount
	}

	return amount.Int64(), nil
}
