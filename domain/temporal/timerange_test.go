package temporal

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type TimeRangeTestSuite struct {
	suite.Suite
}

func TestTimeRangeSuite(t *testing.T) {
	suite.Run(t, new(TimeRangeTestSuite))
}

func (s *TimeRangeTestSuite) TestItCanBuildTimeRangeWithValidValues() {
	start := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 1, 12, 30, 0, 0, time.UTC)

	timeRange, err := NewTimeRange(start, end)

	s.NoError(err)
	s.Equal(start, timeRange.StartTime())
	s.Equal(end, timeRange.EndTime())
	s.Equal(150*time.Minute, timeRange.Duration())
	s.True(timeRange.Contains(time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC)))
	s.False(timeRange.Contains(time.Date(2026, 5, 1, 13, 0, 0, 0, time.UTC)))
	s.Equal("2026-05-01T10:00:00Z/2026-05-01T12:30:00Z", timeRange.String())
}

func (s *TimeRangeTestSuite) TestItCanBuildTimeRangeFromStrings() {
	timeRange, err := NewTimeRangeFromString("2026-05-01T10:00:00Z", "2026-05-01T10:00:00Z")

	s.NoError(err)
	s.Equal(time.Duration(0), timeRange.Duration())
}

func (s *TimeRangeTestSuite) TestItFailsToBuildTimeRangeWithInvalidValues() {
	validTime := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)

	testCases := []struct {
		name          string
		start         time.Time
		end           time.Time
		expectedError error
	}{
		{"zero start", time.Time{}, validTime, ErrZeroTimeRangeStart},
		{"zero end", validTime, time.Time{}, ErrZeroTimeRangeEnd},
		{"start after end", validTime.Add(time.Minute), validTime, ErrInvalidTimeRange},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewTimeRange(tc.start, tc.end)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *TimeRangeTestSuite) TestJSONSerialization() {
	timeRange, _ := NewTimeRangeFromString("2026-05-01T10:00:00Z", "2026-05-01T12:30:00Z")

	jsonData, err := json.Marshal(timeRange)
	s.NoError(err)
	s.JSONEq(`{"startTime":"2026-05-01T10:00:00Z","endTime":"2026-05-01T12:30:00Z"}`, string(jsonData))

	var decoded TimeRange
	s.NoError(json.Unmarshal([]byte(`{"startTime":"2026-05-01T10:00:00Z","endTime":"2026-05-01T12:30:00Z"}`), &decoded))
	s.True(timeRange.Equals(decoded))
}

func (s *TimeRangeTestSuite) TestJSONSerializationFailsForInvalidValues() {
	testCases := []struct {
		name          string
		jsonData      string
		expectedError error
	}{
		{"null", `null`, domain.ErrNullValue},
		{"missing start", `{"endTime":"2026-05-01T12:30:00Z"}`, ErrMissingStartTime},
		{"missing end", `{"startTime":"2026-05-01T10:00:00Z"}`, ErrMissingEndTime},
		{"invalid time", `{"startTime":"2026-05-01 10:00","endTime":"2026-05-01T12:30:00Z"}`, ErrInvalidTimeFormat},
		{"invalid range", `{"startTime":"2026-05-01T12:30:00Z","endTime":"2026-05-01T10:00:00Z"}`, ErrInvalidTimeRange},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var decoded TimeRange
			err := json.Unmarshal([]byte(tc.jsonData), &decoded)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *TimeRangeTestSuite) TestReconstitute() {
	start := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 1, 12, 30, 0, 0, time.UTC)

	timeRange := ReconstituteTimeRange(start, end)

	s.Equal(start, timeRange.StartTime())
	s.Equal(end, timeRange.EndTime())
}
