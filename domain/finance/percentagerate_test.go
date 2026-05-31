package finance

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type PercentageRateTestSuite struct {
	suite.Suite
}

func TestPercentageRateSuite(t *testing.T) {
	suite.Run(t, new(PercentageRateTestSuite))
}

func (s *PercentageRateTestSuite) TestItCanBuildNewPercentageRateWithValidValues() {
	testCases := []struct {
		name                string
		input               string
		expectedBasisPoints int64
		expectedPercent     string
		expectedFraction    string
		expectedString      string
	}{
		{
			name:                "zero rate",
			input:               "0",
			expectedBasisPoints: 0,
			expectedPercent:     "0",
			expectedFraction:    "0",
			expectedString:      "0%",
		},
		{
			name:                "whole percent",
			input:               "19",
			expectedBasisPoints: 1900,
			expectedPercent:     "19",
			expectedFraction:    "0.19",
			expectedString:      "19%",
		},
		{
			name:                "decimal percent",
			input:               "19.5",
			expectedBasisPoints: 1950,
			expectedPercent:     "19.5",
			expectedFraction:    "0.195",
			expectedString:      "19.5%",
		},
		{
			name:                "basis point precision",
			input:               "0.01",
			expectedBasisPoints: 1,
			expectedPercent:     "0.01",
			expectedFraction:    "0.0001",
			expectedString:      "0.01%",
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				rate, err := NewPercentageRateFromString(tc.input)
				s.NoError(err)
				s.Equal(tc.expectedBasisPoints, rate.BasisPoints())
				s.Equal(tc.expectedPercent, rate.Percent().String())
				s.Equal(tc.expectedFraction, rate.Fraction().String())
				s.Equal(tc.expectedString, rate.String())
			},
		)
	}
}

func (s *PercentageRateTestSuite) TestItCanBuildNewPercentageRateFromBasisPoints() {
	rate, err := NewPercentageRateFromBasisPoints(1900)
	s.NoError(err)
	s.Equal(int64(1900), rate.BasisPoints())
	s.Equal("19", rate.Percent().String())
}

func (s *PercentageRateTestSuite) TestItFailsToBuildNewPercentageRateFromInvalidValues() {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "invalid format",
			input:         "abc",
			expectedError: ErrInvalidPercentageRate,
		},
		{
			name:          "negative rate",
			input:         "-1",
			expectedError: ErrNegativePercentageRate,
		},
		{
			name:          "too much precision",
			input:         "0.001",
			expectedError: ErrInvalidPercentageRatePrecision,
		},
	}

	for _, tc := range testCases {
		s.Run(
			tc.name, func() {
				_, err := NewPercentageRateFromString(tc.input)
				s.Error(err)
				s.True(errors.Is(err, tc.expectedError))
			},
		)
	}
}

func (s *PercentageRateTestSuite) TestEquals() {
	rate1, _ := NewPercentageRateFromString("19")
	rate2, _ := NewPercentageRateFromBasisPoints(1900)
	rate3, _ := NewPercentageRateFromString("20")

	s.True(rate1.Equals(rate2))
	s.False(rate1.Equals(rate3))
}

func (s *PercentageRateTestSuite) TestJSONSerialization() {
	rate, _ := NewPercentageRateFromString("19.5")

	jsonData, err := json.Marshal(rate)
	s.NoError(err)
	s.Equal(`"19.5"`, string(jsonData))

	var decoded PercentageRate
	s.NoError(json.Unmarshal([]byte(`"19"`), &decoded))
	s.Equal(int64(1900), decoded.BasisPoints())

	s.NoError(json.Unmarshal([]byte(`19.5`), &decoded))
	s.Equal(int64(1950), decoded.BasisPoints())

	s.Error(json.Unmarshal([]byte(`"0.001"`), &decoded))
}

func (s *PercentageRateTestSuite) TestJSONSerializationFailsForNull() {
	var decoded PercentageRate
	err := json.Unmarshal([]byte(`null`), &decoded)

	s.Error(err)
	s.True(errors.Is(err, domain.ErrNullValue))
}

func (s *PercentageRateTestSuite) TestTextSerialization() {
	rate, _ := NewPercentageRateFromString("19.5")

	text, err := rate.MarshalText()
	s.NoError(err)
	s.Equal("19.5", string(text))

	var decoded PercentageRate
	s.NoError(decoded.UnmarshalText([]byte("19")))
	s.Equal(int64(1900), decoded.BasisPoints())

	s.Error(decoded.UnmarshalText([]byte("0.001")))
}

func (s *PercentageRateTestSuite) TestDatabaseValueAndScan() {
	rate, _ := NewPercentageRateFromString("19.5")

	value, err := rate.Value()
	s.NoError(err)
	s.Equal(int64(1950), value)

	var scanned PercentageRate
	s.NoError(scanned.Scan(value))
	s.True(rate.Equals(scanned))

	s.NoError(scanned.Scan(int64(2000)))
	s.Equal(int64(2000), scanned.BasisPoints())

	s.NoError(scanned.Scan("20"))
	s.Equal(int64(2000), scanned.BasisPoints())

	s.NoError(scanned.Scan([]byte("20")))
	s.Equal(int64(2000), scanned.BasisPoints())

	s.Error(scanned.Scan(int64(-1)))
	s.Error(scanned.Scan(123))
}

func (s *PercentageRateTestSuite) TestApplyToMoney() {
	money, _ := NewMoneyFromString("100.00", "EUR")
	rate, _ := NewPercentageRateFromString("19")

	amount, err := rate.ApplyTo(money)
	s.NoError(err)
	s.Equal("19", amount.Amount().String())
	s.Equal(int64(1900), amount.AmountMinorUnits())
	s.Equal("EUR", amount.Currency().String())

	gross, err := rate.AddTo(money)
	s.NoError(err)
	s.Equal("119", gross.Amount().String())

	discounted, err := rate.SubtractFrom(money)
	s.NoError(err)
	s.Equal("81", discounted.Amount().String())
}

func (s *PercentageRateTestSuite) TestApplyToMoneyUsesRoundingMode() {
	money, _ := NewMoneyFromString("10.01", "EUR")
	rate, _ := NewPercentageRateFromString("19")

	halfUp, err := rate.ApplyToWithRounding(money, RoundHalfUp)
	s.NoError(err)
	s.Equal(int64(190), halfUp.AmountMinorUnits())

	roundDown, err := rate.ApplyToWithRounding(money, RoundDown)
	s.NoError(err)
	s.Equal(int64(190), roundDown.AmountMinorUnits())

	roundUp, err := rate.ApplyToWithRounding(money, RoundUp)
	s.NoError(err)
	s.Equal(int64(191), roundUp.AmountMinorUnits())

	_, err = rate.ApplyToWithRounding(money, RoundingMode(999))
	s.Error(err)
	s.True(errors.Is(err, ErrInvalidRoundingMode))
}

func (s *PercentageRateTestSuite) TestReconstitute() {
	rate := ReconstitutePercentageRate(1900)
	s.Equal(int64(1900), rate.BasisPoints())
	s.Equal("19%", rate.String())
}
