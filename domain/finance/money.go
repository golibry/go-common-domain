package finance

import (
	"encoding/json"
	"fmt"
	"github.com/golibry/go-common-domain/domain"
	"github.com/shopspring/decimal"
)

var (
	ErrNegativeAmount = domain.NewError("money amount cannot be negative")
)

type Money struct {
	amount   decimal.Decimal
	currency Currency
}

type moneyJSON struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// NewMoney creates a new instance of Money with validation
func NewMoney(amount decimal.Decimal, currency Currency) (Money, error) {
	if err := IsValidMoneyAmount(amount); err != nil {
		return Money{}, err
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// NewMoneyFromString creates a new instance of Money from string amount and currency
func NewMoneyFromString(amountStr, currencyStr string) (Money, error) {
	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return Money{}, domain.NewError("invalid amount format: %s", err)
	}

	currency, err := NewCurrency(currencyStr)
	if err != nil {
		return Money{}, err
	}

	return NewMoney(amount, currency)
}

// ReconstituteMoney creates a new Money instance without validation
func ReconstituteMoney(amount decimal.Decimal, currency Currency) Money {
	return Money{
		amount:   amount,
		currency: currency,
	}
}

// NewMoneyFromJSON creates Money from JSON bytes array
func NewMoneyFromJSON(data []byte) (Money, error) {
	var temp moneyJSON

	if err := json.Unmarshal(data, &temp); err != nil {
		return Money{}, domain.NewError("failed to build money from json: %s", err)
	}

	return NewMoneyFromString(temp.Amount, temp.Currency)
}

// Amount returns the money amount
func (m Money) Amount() decimal.Decimal {
	return m.amount
}

// Currency returns the money currency
func (m Money) Currency() Currency {
	return m.currency
}

// Equals compares two Money objects for equality
func (m Money) Equals(other Money) bool {
	return m.amount.Equal(other.amount) && m.currency.Equals(other.currency)
}

// String returns a string representation of the money
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.amount.String(), m.currency.String())
}

// Add adds another Money object to this one (must have same currency)
func (m Money) Add(other Money) (Money, error) {
	if !m.currency.Equals(other.currency) {
		return Money{}, domain.NewError("cannot add money with different currencies: %s and %s", m.currency.String(), other.currency.String())
	}

	newAmount := m.amount.Add(other.amount)
	return Money{
		amount:   newAmount,
		currency: m.currency,
	}, nil
}

// Subtract subtracts another Money object from this one (must have same currency)
func (m Money) Subtract(other Money) (Money, error) {
	if !m.currency.Equals(other.currency) {
		return Money{}, domain.NewError("cannot subtract money with different currencies: %s and %s", m.currency.String(), other.currency.String())
	}

	newAmount := m.amount.Sub(other.amount)
	if newAmount.IsNegative() {
		return Money{}, ErrNegativeAmount
	}

	return Money{
		amount:   newAmount,
		currency: m.currency,
	}, nil
}

// Multiply multiplies the money amount by a factor
func (m Money) Multiply(factor decimal.Decimal) (Money, error) {
	newAmount := m.amount.Mul(factor)
	if newAmount.IsNegative() {
		return Money{}, ErrNegativeAmount
	}

	return Money{
		amount:   newAmount,
		currency: m.currency,
	}, nil
}

// MarshalJSON implements json.Marshaler
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		moneyJSON{
			Amount:   m.amount.String(),
			Currency: m.currency.String(),
		},
	)
}

// IsValidMoneyAmount validates a money amount (must not be negative)
func IsValidMoneyAmount(amount decimal.Decimal) error {
	if amount.IsNegative() {
		return ErrNegativeAmount
	}
	return nil
}