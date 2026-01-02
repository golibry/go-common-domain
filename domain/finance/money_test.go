package finance

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type MoneyTestSuite struct {
	suite.Suite
}

func TestMoneySuite(t *testing.T) {
	suite.Run(t, new(MoneyTestSuite))
}

func (s *MoneyTestSuite) TestItCanBuildNewMoneyWithValidValues() {
	testCases := []struct {
		name             string
		amount           string
		currency         string
		expectedAmount   string
		expectedCurrency string
	}{
		{
			name:             "USD money",
			amount:           "100.50",
			currency:         "USD",
			expectedAmount:   "100.5",
			expectedCurrency: "USD",
		},
		{
			name:             "zero amount",
			amount:           "0",
			currency:         "EUR",
			expectedAmount:   "0",
			expectedCurrency: "EUR",
		},
		{
			name:             "decimal amount",
			amount:           "99.99",
			currency:         "GBP",
			expectedAmount:   "99.99",
			expectedCurrency: "GBP",
		},
		{
			name:             "large amount",
			amount:           "1000000.00",
			currency:         "JPY",
			expectedAmount:   "1000000",
			expectedCurrency: "JPY",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				money, err := NewMoneyFromString(tc.amount, tc.currency)
				s.NoError(err)
				s.Equal(tc.expectedAmount, money.Amount().String())
				s.Equal(tc.expectedCurrency, money.Currency().String())
			},
		)
	}
}

func (s *MoneyTestSuite) TestItFailsToBuildNewMoneyFromInvalidValues() {
	testCases := []struct {
		name          string
		amount        string
		currency      string
		expectedError error
	}{
		{
			name:          "negative amount",
			amount:        "-100.50",
			currency:      "USD",
			expectedError: ErrNegativeAmount,
		},
		{
			name:          "invalid currency",
			amount:        "100.50",
			currency:      "INVALID",
			expectedError: ErrInvalidCurrency,
		},
		{
			name:          "invalid amount format",
			amount:        "abc",
			currency:      "USD",
			expectedError: nil, // Will be a different error about format
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewMoneyFromString(tc.amount, tc.currency)
				s.Error(err)
				if tc.expectedError != nil {
					s.True(errors.Is(err, tc.expectedError))
				}
			},
		)
	}
}

func (s *MoneyTestSuite) TestMoneyArithmetic() {
	usd, _ := NewCurrency("USD")
	eur, _ := NewCurrency("EUR")

	money1, _ := NewMoney(decimal.NewFromFloat(100.50), usd)
	money2, _ := NewMoney(decimal.NewFromFloat(50.25), usd)
	money3, _ := NewMoney(decimal.NewFromFloat(25.00), eur)

	s.Run(
		"addition with same currency", func() {
			result, err := money1.Add(money2)
			s.NoError(err)
			s.Equal("150.75", result.Amount().String())
			s.True(result.Currency().Equals(usd))
		},
	)

	s.Run(
		"addition with different currency fails", func() {
			_, err := money1.Add(money3)
			s.Error(err)
		},
	)

	s.Run(
		"subtraction with same currency", func() {
			result, err := money1.Subtract(money2)
			s.NoError(err)
			s.Equal("50.25", result.Amount().String())
			s.True(result.Currency().Equals(usd))
		},
	)

	s.Run(
		"subtraction with different currency fails", func() {
			_, err := money1.Subtract(money3)
			s.Error(err)
		},
	)

	s.Run(
		"subtraction resulting in negative fails", func() {
			_, err := money2.Subtract(money1)
			s.Error(err)
			s.True(errors.Is(err, ErrNegativeAmount))
		},
	)

	s.Run(
		"multiplication", func() {
			result, err := money1.Multiply(decimal.NewFromFloat(2))
			s.NoError(err)
			s.Equal("201", result.Amount().String())
			s.True(result.Currency().Equals(usd))
		},
	)

	s.Run(
		"multiplication by negative factor fails", func() {
			_, err := money1.Multiply(decimal.NewFromFloat(-1))
			s.Error(err)
			s.True(errors.Is(err, ErrNegativeAmount))
		},
	)

	s.Run(
		"division", func() {
			result, err := money1.Divide(decimal.NewFromFloat(2))
			s.NoError(err)
			s.Equal("50.25", result.Amount().String())
			s.True(result.Currency().Equals(usd))
		},
	)

	s.Run(
		"division by zero fails", func() {
			_, err := money1.Divide(decimal.Zero)
			s.Error(err)
		},
	)

	s.Run(
		"division by negative factor fails", func() {
			_, err := money1.Divide(decimal.NewFromFloat(-1))
			s.Error(err)
			s.True(errors.Is(err, ErrNegativeAmount))
		},
	)
}

func (s *MoneyTestSuite) TestEquals() {
	usd, _ := NewCurrency("USD")
	eur, _ := NewCurrency("EUR")

	money1, _ := NewMoney(decimal.NewFromFloat(100.50), usd)
	money2, _ := NewMoney(decimal.NewFromFloat(100.50), usd)
	money3, _ := NewMoney(decimal.NewFromFloat(100.50), eur)
	money4, _ := NewMoney(decimal.NewFromFloat(200.00), usd)

	s.True(money1.Equals(money2))
	s.False(money1.Equals(money3)) // Different currency
	s.False(money1.Equals(money4)) // Different amount
}

func (s *MoneyTestSuite) TestString() {
	usd, _ := NewCurrency("USD")
	money, _ := NewMoney(decimal.NewFromFloat(100.50), usd)
	s.Equal("100.5 USD", money.String())
}

func (s *MoneyTestSuite) TestReconstitute() {
	usd, _ := NewCurrency("USD")
	amount := decimal.NewFromFloat(100.50)
	money := ReconstituteMoney(amount, usd)

	s.Equal("100.5", money.Amount().String())
	s.Equal("USD", money.Currency().String())
}
