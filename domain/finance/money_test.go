package finance

import (
	"encoding/json"
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

func (s *MoneyTestSuite) TestItCanBuildNewMoneyFromMinorUnits() {
	eur, _ := NewCurrency("EUR")

	money, err := NewMoneyFromMinorUnits(1099, eur, 2)
	s.NoError(err)
	s.Equal("10.99", money.Amount().String())
	s.Equal(int64(1099), money.AmountMinorUnits())
	s.Equal(int32(2), money.Scale())
	s.Equal("EUR", money.Currency().String())
}

func (s *MoneyTestSuite) TestItCanBuildNewMoneyWithExplicitScale() {
	jpy, _ := NewCurrency("JPY")

	money, err := NewMoneyWithScale(decimal.NewFromInt(1000), jpy, 0)
	s.NoError(err)
	s.Equal("1000", money.Amount().String())
	s.Equal(int64(1000), money.AmountMinorUnits())
	s.Equal(int32(0), money.Scale())
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

func (s *MoneyTestSuite) TestItFailsToBuildNewMoneyWithInvalidScaleOrPrecision() {
	usd, _ := NewCurrency("USD")

	_, err := NewMoneyWithScale(decimal.NewFromFloat(10.999), usd, 2)
	s.Error(err)
	s.True(errors.Is(err, ErrInvalidMoneyAmountPrecision))

	_, err = NewMoneyFromMinorUnits(100, usd, -1)
	s.Error(err)
	s.True(errors.Is(err, ErrInvalidMoneyScale))

	_, err = NewMoneyFromMinorUnits(-100, usd, 2)
	s.Error(err)
	s.True(errors.Is(err, ErrNegativeAmount))
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

func (s *MoneyTestSuite) TestMoneyMultiplicationUsesRoundingMode() {
	money, _ := NewMoneyFromString("10.03", "USD")
	factor := decimal.NewFromFloat(0.5)

	halfUp, err := money.Multiply(factor)
	s.NoError(err)
	s.Equal("5.02", halfUp.Amount().String())
	s.Equal(int64(502), halfUp.AmountMinorUnits())

	roundDown, err := money.MultiplyWithRounding(factor, RoundDown)
	s.NoError(err)
	s.Equal("5.01", roundDown.Amount().String())
	s.Equal(int64(501), roundDown.AmountMinorUnits())

	roundUp, err := money.MultiplyWithRounding(factor, RoundUp)
	s.NoError(err)
	s.Equal("5.02", roundUp.Amount().String())
	s.Equal(int64(502), roundUp.AmountMinorUnits())

	_, err = money.MultiplyWithRounding(factor, RoundingMode(999))
	s.Error(err)
	s.True(errors.Is(err, ErrInvalidRoundingMode))
}

func (s *MoneyTestSuite) TestMoneyDivisionUsesRoundingMode() {
	money, _ := NewMoneyFromString("10.05", "USD")
	divisor := decimal.NewFromInt(2)

	halfUp, err := money.Divide(divisor)
	s.NoError(err)
	s.Equal("5.03", halfUp.Amount().String())
	s.Equal(int64(503), halfUp.AmountMinorUnits())

	roundDown, err := money.DivideWithRounding(divisor, RoundDown)
	s.NoError(err)
	s.Equal("5.02", roundDown.Amount().String())
	s.Equal(int64(502), roundDown.AmountMinorUnits())

	roundUp, err := money.DivideWithRounding(divisor, RoundUp)
	s.NoError(err)
	s.Equal("5.03", roundUp.Amount().String())
	s.Equal(int64(503), roundUp.AmountMinorUnits())

	_, err = money.DivideWithRounding(divisor, RoundingMode(999))
	s.Error(err)
	s.True(errors.Is(err, ErrInvalidRoundingMode))
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

func (s *MoneyTestSuite) TestJSONSerialization() {
	money, _ := NewMoneyFromString("100.50", "usd")

	jsonData, err := json.Marshal(money)
	s.NoError(err)
	s.JSONEq(`{"amount":"100.50","currency":"USD","scale":2}`, string(jsonData))

	var decoded Money
	s.NoError(json.Unmarshal([]byte(`{"amount":"99.99","currency":"eur"}`), &decoded))
	s.Equal("99.99", decoded.Amount().String())
	s.Equal("EUR", decoded.Currency().String())
	s.Equal(int64(9999), decoded.AmountMinorUnits())
	s.Equal(int32(2), decoded.Scale())

	s.NoError(json.Unmarshal([]byte(`{"amount":42.25,"currency":"gbp"}`), &decoded))
	s.Equal("42.25", decoded.Amount().String())
	s.Equal("GBP", decoded.Currency().String())

	s.NoError(json.Unmarshal([]byte(`{"amount":"42.125","currency":"gbp","scale":3}`), &decoded))
	s.Equal("42.125", decoded.Amount().String())
	s.Equal(int64(42125), decoded.AmountMinorUnits())
	s.Equal(int32(3), decoded.Scale())
}

func (s *MoneyTestSuite) TestJSONSerializationFailsForInvalidValues() {
	testCases := []struct {
		name          string
		jsonData      string
		expectedError error
	}{
		{
			name:          "negative amount",
			jsonData:      `{"amount":"-1","currency":"USD"}`,
			expectedError: ErrNegativeAmount,
		},
		{
			name:          "invalid amount",
			jsonData:      `{"amount":"abc","currency":"USD"}`,
			expectedError: nil,
		},
		{
			name:          "invalid currency",
			jsonData:      `{"amount":"1","currency":"USDD"}`,
			expectedError: ErrInvalidCurrency,
		},
		{
			name:          "missing amount",
			jsonData:      `{"currency":"USD"}`,
			expectedError: nil,
		},
		{
			name:          "missing currency",
			jsonData:      `{"amount":"1"}`,
			expectedError: ErrEmptyCurrency,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				var decoded Money
				err := json.Unmarshal([]byte(tc.jsonData), &decoded)
				s.Error(err)
				if tc.expectedError != nil {
					s.True(errors.Is(err, tc.expectedError))
				}
			},
		)
	}
}

func (s *MoneyTestSuite) TestReconstitute() {
	usd, _ := NewCurrency("USD")
	amount := decimal.NewFromFloat(100.50)
	money := ReconstituteMoney(amount, usd)

	s.Equal("100.5", money.Amount().String())
	s.Equal(int64(10050), money.AmountMinorUnits())
	s.Equal("USD", money.Currency().String())
	s.Equal(int32(DefaultMoneyScale), money.Scale())
}

func (s *MoneyTestSuite) TestReconstituteRoundsToDefaultScale() {
	usd, _ := NewCurrency("USD")
	amount := decimal.RequireFromString("10.999")

	money := ReconstituteMoney(amount, usd)

	s.Equal("11", money.Amount().String())
	s.Equal(int64(1100), money.AmountMinorUnits())
	s.Equal(int32(DefaultMoneyScale), money.Scale())
}

func (s *MoneyTestSuite) TestReconstituteFromMinorUnits() {
	kwd, _ := NewCurrency("KWD")

	money := ReconstituteMoneyFromMinorUnits(10999, kwd, 3)

	s.Equal("10.999", money.Amount().String())
	s.Equal(int64(10999), money.AmountMinorUnits())
	s.Equal("KWD", money.Currency().String())
	s.Equal(int32(3), money.Scale())
}
