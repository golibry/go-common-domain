package finance

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/golibry/go-common-domain/domain"
	"github.com/shopspring/decimal"
)

const basisPointsPerPercent int64 = 100
const basisPointsPerWhole int64 = 10000

var (
	ErrNegativePercentageRate          = domain.NewError("percentage rate cannot be negative")
	ErrInvalidPercentageRate           = domain.NewError("percentage rate format is invalid")
	ErrInvalidPercentageRatePrecision  = domain.NewError("percentage rate has more precision than basis points allow")
	ErrPercentageRateBasisPointsTooBig = domain.NewError("percentage rate basis points value is too large")
	ErrInvalidRoundingMode             = domain.NewError("rounding mode is invalid")
	ErrInvalidPercentageRateScanValue  = domain.NewError("percentage rate scan value must be integer basis points")
)

type RoundingMode int

const (
	RoundHalfUp RoundingMode = iota
	RoundDown
	RoundUp
)

type PercentageRate struct {
	basisPoints int64
}

// NewPercentageRateFromBasisPoints creates a rate from integer basis points.
func NewPercentageRateFromBasisPoints(basisPoints int64) (PercentageRate, error) {
	if basisPoints < 0 {
		return PercentageRate{}, ErrNegativePercentageRate
	}

	return PercentageRate{
		basisPoints: basisPoints,
	}, nil
}

// NewPercentageRateFromString creates a rate from a percentage string, such as "19" or "19.5".
func NewPercentageRateFromString(value string) (PercentageRate, error) {
	percentage, err := decimal.NewFromString(value)
	if err != nil {
		return PercentageRate{}, ErrInvalidPercentageRate
	}

	if percentage.IsNegative() {
		return PercentageRate{}, ErrNegativePercentageRate
	}

	basisPoints := percentage.Mul(decimal.NewFromInt(basisPointsPerPercent))
	if !basisPoints.IsInteger() {
		return PercentageRate{}, ErrInvalidPercentageRatePrecision
	}

	raw := basisPoints.BigInt()
	if !raw.IsInt64() {
		return PercentageRate{}, ErrPercentageRateBasisPointsTooBig
	}

	return NewPercentageRateFromBasisPoints(raw.Int64())
}

// ReconstitutePercentageRate creates a PercentageRate without validation.
func ReconstitutePercentageRate(basisPoints int64) PercentageRate {
	return PercentageRate{
		basisPoints: basisPoints,
	}
}

// BasisPoints returns the rate as integer basis points.
func (r PercentageRate) BasisPoints() int64 {
	return r.basisPoints
}

// Percent returns the rate as a percentage value.
func (r PercentageRate) Percent() decimal.Decimal {
	return decimal.New(r.basisPoints, -2)
}

// Fraction returns the rate as a multiplier fraction.
func (r PercentageRate) Fraction() decimal.Decimal {
	return decimal.New(r.basisPoints, -4)
}

// Equals compares two PercentageRate objects for equality.
func (r PercentageRate) Equals(other PercentageRate) bool {
	return r.basisPoints == other.basisPoints
}

// String returns the rate as a percentage string.
func (r PercentageRate) String() string {
	return fmt.Sprintf("%s%%", r.Percent().String())
}

// MarshalText returns the percentage value as text.
func (r PercentageRate) MarshalText() ([]byte, error) {
	return []byte(r.Percent().String()), nil
}

// UnmarshalText validates and normalizes text into a PercentageRate.
func (r *PercentageRate) UnmarshalText(text []byte) error {
	rate, err := NewPercentageRateFromString(string(text))
	if err != nil {
		return err
	}

	*r = rate
	return nil
}

// MarshalJSON returns the percentage value as a decimal-safe string.
func (r PercentageRate) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Percent().String())
}

// UnmarshalJSON validates and normalizes a JSON percentage string or number.
func (r *PercentageRate) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		var number json.Number
		if numberErr := json.Unmarshal(data, &number); numberErr != nil {
			return ErrInvalidPercentageRate
		}
		value = number.String()
	}

	return r.UnmarshalText([]byte(value))
}

// Value returns the canonical basis points for database storage.
func (r PercentageRate) Value() (driver.Value, error) {
	return r.basisPoints, nil
}

// Scan validates and normalizes a database value into a PercentageRate.
func (r *PercentageRate) Scan(value any) error {
	switch v := value.(type) {
	case int64:
		rate, err := NewPercentageRateFromBasisPoints(v)
		if err != nil {
			return err
		}

		*r = rate
		return nil
	case string:
		rate, err := NewPercentageRateFromString(v)
		if err != nil {
			return err
		}

		*r = rate
		return nil
	case []byte:
		rate, err := NewPercentageRateFromString(string(v))
		if err != nil {
			return err
		}

		*r = rate
		return nil
	default:
		return ErrInvalidPercentageRateScanValue
	}
}

// ApplyTo returns the money amount represented by this rate using half-up rounding.
func (r PercentageRate) ApplyTo(money Money) (Money, error) {
	return r.ApplyToWithRounding(money, RoundHalfUp)
}

// ApplyToWithRounding returns the money amount represented by this rate.
func (r PercentageRate) ApplyToWithRounding(money Money, mode RoundingMode) (Money, error) {
	amountMinor, err := multiplyMinorUnitsByBasisPoints(money.AmountMinorUnits(), r.basisPoints, mode)
	if err != nil {
		return Money{}, err
	}

	return NewMoneyFromMinorUnits(amountMinor, money.Currency(), money.Scale())
}

// AddTo returns money plus this rate amount using half-up rounding.
func (r PercentageRate) AddTo(money Money) (Money, error) {
	return r.AddToWithRounding(money, RoundHalfUp)
}

// AddToWithRounding returns money plus this rate amount.
func (r PercentageRate) AddToWithRounding(money Money, mode RoundingMode) (Money, error) {
	rateAmount, err := r.ApplyToWithRounding(money, mode)
	if err != nil {
		return Money{}, err
	}

	return money.Add(rateAmount)
}

// SubtractFrom returns money minus this rate amount using half-up rounding.
func (r PercentageRate) SubtractFrom(money Money) (Money, error) {
	return r.SubtractFromWithRounding(money, RoundHalfUp)
}

// SubtractFromWithRounding returns money minus this rate amount.
func (r PercentageRate) SubtractFromWithRounding(money Money, mode RoundingMode) (Money, error) {
	rateAmount, err := r.ApplyToWithRounding(money, mode)
	if err != nil {
		return Money{}, err
	}

	return money.Subtract(rateAmount)
}

func multiplyMinorUnitsByBasisPoints(amountMinor, basisPoints int64, mode RoundingMode) (int64, error) {
	if amountMinor < 0 {
		return 0, ErrNegativeAmount
	}

	if basisPoints < 0 {
		return 0, ErrNegativePercentageRate
	}

	numerator := big.NewInt(amountMinor)
	numerator.Mul(numerator, big.NewInt(basisPoints))

	rounded, err := roundBigIntQuotient(numerator, big.NewInt(basisPointsPerWhole), mode)
	if err != nil {
		return 0, err
	}

	return checkedMoneyMinorAmount(rounded)
}

func roundBigIntQuotient(numerator, denominator *big.Int, mode RoundingMode) (*big.Int, error) {
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.QuoRem(numerator, denominator, remainder)

	switch mode {
	case RoundDown:
		return quotient, nil
	case RoundUp:
		if remainder.Sign() > 0 {
			quotient.Add(quotient, big.NewInt(1))
		}
		return quotient, nil
	case RoundHalfUp:
		doubleRemainder := new(big.Int).Mul(remainder, big.NewInt(2))
		if doubleRemainder.Cmp(denominator) >= 0 {
			quotient.Add(quotient, big.NewInt(1))
		}
		return quotient, nil
	default:
		return nil, ErrInvalidRoundingMode
	}
}
