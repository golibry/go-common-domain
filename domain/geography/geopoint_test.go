package geography

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golibry/go-common-domain/domain"
	"github.com/stretchr/testify/suite"
)

type GeoPointTestSuite struct {
	suite.Suite
}

func TestGeoPointSuite(t *testing.T) {
	suite.Run(t, new(GeoPointTestSuite))
}

func (s *GeoPointTestSuite) TestItCanBuildGeoPointWithValidValues() {
	point, err := NewGeoPoint(44.4268, 26.1025)

	s.NoError(err)
	s.Equal(44.4268, point.Latitude())
	s.Equal(26.1025, point.Longitude())
	s.Equal("44.4268,26.1025", point.String())
}

func (s *GeoPointTestSuite) TestItFailsToBuildGeoPointWithInvalidValues() {
	testCases := []struct {
		name          string
		latitude      float64
		longitude     float64
		expectedError error
	}{
		{"latitude too low", -90.1, 0, ErrInvalidLatitude},
		{"latitude too high", 90.1, 0, ErrInvalidLatitude},
		{"longitude too low", 0, -180.1, ErrInvalidLongitude},
		{"longitude too high", 0, 180.1, ErrInvalidLongitude},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := NewGeoPoint(tc.latitude, tc.longitude)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *GeoPointTestSuite) TestEquals() {
	point1, _ := NewGeoPoint(44.4268, 26.1025)
	point2, _ := NewGeoPoint(44.4268, 26.1025)
	point3, _ := NewGeoPoint(45.0, 26.1025)

	s.True(point1.Equals(point2))
	s.False(point1.Equals(point3))
}

func (s *GeoPointTestSuite) TestJSONSerialization() {
	point, _ := NewGeoPoint(44.4268, 26.1025)

	jsonData, err := json.Marshal(point)
	s.NoError(err)
	s.JSONEq(`{"latitude":44.4268,"longitude":26.1025}`, string(jsonData))

	var decoded GeoPoint
	s.NoError(json.Unmarshal([]byte(`{"latitude":0,"longitude":0}`), &decoded))
	s.Equal(0.0, decoded.Latitude())
	s.Equal(0.0, decoded.Longitude())
}

func (s *GeoPointTestSuite) TestJSONSerializationFailsForInvalidValues() {
	testCases := []struct {
		name          string
		jsonData      string
		expectedError error
	}{
		{"null", `null`, domain.ErrNullValue},
		{"missing latitude", `{"longitude":26.1025}`, ErrMissingLatitude},
		{"missing longitude", `{"latitude":44.4268}`, ErrMissingLongitude},
		{"invalid latitude", `{"latitude":90.1,"longitude":26.1025}`, ErrInvalidLatitude},
		{"invalid longitude", `{"latitude":44.4268,"longitude":180.1}`, ErrInvalidLongitude},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var decoded GeoPoint
			err := json.Unmarshal([]byte(tc.jsonData), &decoded)

			s.Error(err)
			s.True(errors.Is(err, tc.expectedError))
		})
	}
}

func (s *GeoPointTestSuite) TestReconstitute() {
	point := ReconstituteGeoPoint(44.4268, 26.1025)

	s.Equal(44.4268, point.Latitude())
	s.Equal(26.1025, point.Longitude())
}
