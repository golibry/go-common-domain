package temporal

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type DateRangeTestSuite struct {
	suite.Suite
}

func TestDateRangeSuite(t *testing.T) {
	suite.Run(t, new(DateRangeTestSuite))
}

func (s *DateRangeTestSuite) TestItCanBuildDateRangeWithValidValues() {
	start := time.Date(2026, 5, 1, 10, 30, 0, 0, time.UTC)
	end := time.Date(2026, 5, 3, 22, 15, 0, 0, time.UTC)

	dateRange, err := NewDateRange(start, end)

	s.NoError(err)
	s.Equal("2026-05-01", dateRange.StartDate().Format(dateLayout))
	s.Equal("2026-05-03", dateRange.EndDate().Format(dateLayout))
	s.Equal(3, dateRange.Days())
	s.True(dateRange.Contains(time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)))
	s.False(dateRange.Contains(time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC)))
	s.Equal("2026-05-01/2026-05-03", dateRange.String())
}

func (s *DateRangeTestSuite) TestItCanBuildDateRangeFromStrings() {
	dateRange, err := NewDateRangeFromString("2026-05-01", "2026-05-01")

	s.NoError(err)
	s.Equal(1, dateRange.Days())
}

func (s *DateRangeTestSuite) TestItFailsToBuildDateRangeWithInvalidValues() {
	validDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name          string
		start         time.Time
		end           time.Time
		expectedError error
	}{
		{"zero start", time.Time{}, validDate, ErrZeroDateRangeStart},
		{"zero end", validDate, time.Time{}, ErrZeroDateRangeEnd},
		{"start after end", validDate.AddDate(0, 0, 1), validDate, ErrInvalidDateRange},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewDateRange(tc.start, tc.end)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *DateRangeTestSuite) TestJSONSerialization() {
	dateRange, _ := NewDateRangeFromString("2026-05-01", "2026-05-03")

	jsonData, err := json.Marshal(dateRange)
	s.NoError(err)
	s.JSONEq(`{"startDate":"2026-05-01","endDate":"2026-05-03"}`, string(jsonData))

	var decoded DateRange
	s.NoError(json.Unmarshal([]byte(`{"startDate":"2026-05-01","endDate":"2026-05-03"}`), &decoded))
	s.True(dateRange.Equals(decoded))
}

func (s *DateRangeTestSuite) TestJSONSerializationFailsForInvalidValues() {
	testCases := []struct {
		name          string
		jsonData      string
		expectedError error
	}{
		{"null", `null`, domain.ErrNullValue},
		{"missing start", `{"endDate":"2026-05-03"}`, ErrMissingStartDate},
		{"missing end", `{"startDate":"2026-05-01"}`, ErrMissingEndDate},
		{"invalid date", `{"startDate":"2026/05/01","endDate":"2026-05-03"}`, ErrInvalidDateFormat},
		{"invalid range", `{"startDate":"2026-05-03","endDate":"2026-05-01"}`, ErrInvalidDateRange},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var decoded DateRange
			err := json.Unmarshal([]byte(tc.jsonData), &decoded)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *DateRangeTestSuite) TestReconstitute() {
	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC)

	dateRange := ReconstituteDateRange(start, end)

	s.Equal(start, dateRange.StartDate())
	s.Equal(end, dateRange.EndDate())
}
