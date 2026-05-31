package temporal

import (
	"encoding/json"
	"time"

	"github.com/golibry/go-common-domain/domain"
)

const dateLayout = "2006-01-02"

var (
	ErrZeroDateRangeStart = domain.NewError("date range start date cannot be zero")
	ErrZeroDateRangeEnd   = domain.NewError("date range end date cannot be zero")
	ErrInvalidDateRange   = domain.NewError("date range start date must be before or equal to end date")
	ErrMissingStartDate   = domain.NewError("date range start date is required")
	ErrMissingEndDate     = domain.NewError("date range end date is required")
	ErrInvalidDateFormat  = domain.NewError("date must use YYYY-MM-DD format")
)

type DateRange struct {
	startDate time.Time
	endDate   time.Time
}

type dateRangeJSON struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// NewDateRange creates a new DateRange with validation and date-only normalization.
func NewDateRange(startDate, endDate time.Time) (DateRange, error) {
	startDate = dateOnly(startDate)
	endDate = dateOnly(endDate)

	if startDate.IsZero() {
		return DateRange{}, ErrZeroDateRangeStart
	}

	if endDate.IsZero() {
		return DateRange{}, ErrZeroDateRangeEnd
	}

	if startDate.After(endDate) {
		return DateRange{}, ErrInvalidDateRange
	}

	return DateRange{
		startDate: startDate,
		endDate:   endDate,
	}, nil
}

// NewDateRangeFromString creates a new DateRange from YYYY-MM-DD values.
func NewDateRangeFromString(startDate, endDate string) (DateRange, error) {
	start, err := parseDate(startDate)
	if err != nil {
		return DateRange{}, err
	}

	end, err := parseDate(endDate)
	if err != nil {
		return DateRange{}, err
	}

	return NewDateRange(start, end)
}

// ReconstituteDateRange creates a DateRange from trusted persisted dates.
func ReconstituteDateRange(startDate, endDate time.Time) DateRange {
	return DateRange{
		startDate: startDate,
		endDate:   endDate,
	}
}

// StartDate returns the range start date.
func (r DateRange) StartDate() time.Time {
	return r.startDate
}

// EndDate returns the range end date.
func (r DateRange) EndDate() time.Time {
	return r.endDate
}

// Contains reports whether date is inside the range, inclusive.
func (r DateRange) Contains(date time.Time) bool {
	date = dateOnly(date)
	return !date.Before(r.startDate) && !date.After(r.endDate)
}

// Days returns the inclusive number of calendar days in the range.
func (r DateRange) Days() int {
	return int(r.endDate.Sub(r.startDate).Hours()/24) + 1
}

// Equals compares two DateRange objects for equality.
func (r DateRange) Equals(other DateRange) bool {
	return r.startDate.Equal(other.startDate) && r.endDate.Equal(other.endDate)
}

// String returns the date range as YYYY-MM-DD/YYYY-MM-DD.
func (r DateRange) String() string {
	return r.startDate.Format(dateLayout) + "/" + r.endDate.Format(dateLayout)
}

// MarshalJSON returns the date range as an explicit object.
func (r DateRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(dateRangeJSON{
		StartDate: r.startDate.Format(dateLayout),
		EndDate:   r.endDate.Format(dateLayout),
	})
}

// UnmarshalJSON validates and normalizes a JSON date range object.
func (r *DateRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var raw dateRangeJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.StartDate == "" {
		return ErrMissingStartDate
	}

	if raw.EndDate == "" {
		return ErrMissingEndDate
	}

	dateRange, err := NewDateRangeFromString(raw.StartDate, raw.EndDate)
	if err != nil {
		return err
	}

	*r = dateRange
	return nil
}

func parseDate(value string) (time.Time, error) {
	parsed, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, ErrInvalidDateFormat
	}

	return parsed, nil
}

func dateOnly(value time.Time) time.Time {
	if value.IsZero() {
		return time.Time{}
	}

	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}
