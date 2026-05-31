package temporal

import (
	"encoding/json"
	"time"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrZeroTimeRangeStart = domain.NewError("time range start cannot be zero")
	ErrZeroTimeRangeEnd   = domain.NewError("time range end cannot be zero")
	ErrInvalidTimeRange   = domain.NewError("time range start must be before or equal to end")
	ErrMissingStartTime   = domain.NewError("time range start is required")
	ErrMissingEndTime     = domain.NewError("time range end is required")
	ErrInvalidTimeFormat  = domain.NewError("time must use RFC3339 format")
)

type TimeRange struct {
	startTime time.Time
	endTime   time.Time
}

type timeRangeJSON struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// NewTimeRange creates a new TimeRange with validation.
func NewTimeRange(startTime, endTime time.Time) (TimeRange, error) {
	if startTime.IsZero() {
		return TimeRange{}, ErrZeroTimeRangeStart
	}

	if endTime.IsZero() {
		return TimeRange{}, ErrZeroTimeRangeEnd
	}

	if startTime.After(endTime) {
		return TimeRange{}, ErrInvalidTimeRange
	}

	return TimeRange{
		startTime: startTime,
		endTime:   endTime,
	}, nil
}

// NewTimeRangeFromString creates a new TimeRange from RFC3339 values.
func NewTimeRangeFromString(startTime, endTime string) (TimeRange, error) {
	start, err := parseTime(startTime)
	if err != nil {
		return TimeRange{}, err
	}

	end, err := parseTime(endTime)
	if err != nil {
		return TimeRange{}, err
	}

	return NewTimeRange(start, end)
}

// ReconstituteTimeRange creates a TimeRange from trusted persisted timestamps.
func ReconstituteTimeRange(startTime, endTime time.Time) TimeRange {
	return TimeRange{
		startTime: startTime,
		endTime:   endTime,
	}
}

// StartTime returns the range start time.
func (r TimeRange) StartTime() time.Time {
	return r.startTime
}

// EndTime returns the range end time.
func (r TimeRange) EndTime() time.Time {
	return r.endTime
}

// Contains reports whether timestamp is inside the range, inclusive.
func (r TimeRange) Contains(timestamp time.Time) bool {
	return !timestamp.Before(r.startTime) && !timestamp.After(r.endTime)
}

// Duration returns the range duration.
func (r TimeRange) Duration() time.Duration {
	return r.endTime.Sub(r.startTime)
}

// Equals compares two TimeRange objects for equality.
func (r TimeRange) Equals(other TimeRange) bool {
	return r.startTime.Equal(other.startTime) && r.endTime.Equal(other.endTime)
}

// String returns the time range as RFC3339/RFC3339.
func (r TimeRange) String() string {
	return r.startTime.Format(time.RFC3339) + "/" + r.endTime.Format(time.RFC3339)
}

// MarshalJSON returns the time range as an explicit object.
func (r TimeRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(timeRangeJSON{
		StartTime: r.startTime.Format(time.RFC3339),
		EndTime:   r.endTime.Format(time.RFC3339),
	})
}

// UnmarshalJSON validates and normalizes a JSON time range object.
func (r *TimeRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var raw timeRangeJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.StartTime == "" {
		return ErrMissingStartTime
	}

	if raw.EndTime == "" {
		return ErrMissingEndTime
	}

	timeRange, err := NewTimeRangeFromString(raw.StartTime, raw.EndTime)
	if err != nil {
		return err
	}

	*r = timeRange
	return nil
}

func parseTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, ErrInvalidTimeFormat
	}

	return parsed, nil
}
