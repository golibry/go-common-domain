package geography

import (
	"encoding/json"
	"fmt"

	"github.com/golibry/go-common-domain/domain"
)

var (
	ErrMissingLatitude  = domain.NewError("latitude is required")
	ErrMissingLongitude = domain.NewError("longitude is required")
	ErrInvalidLatitude  = domain.NewError("latitude must be between -90 and 90")
	ErrInvalidLongitude = domain.NewError("longitude must be between -180 and 180")
)

type GeoPoint struct {
	latitude  float64
	longitude float64
}

type geoPointJSON struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// NewGeoPoint creates a new GeoPoint with validation.
func NewGeoPoint(latitude, longitude float64) (GeoPoint, error) {
	if latitude < -90 || latitude > 90 {
		return GeoPoint{}, ErrInvalidLatitude
	}

	if longitude < -180 || longitude > 180 {
		return GeoPoint{}, ErrInvalidLongitude
	}

	return GeoPoint{
		latitude:  latitude,
		longitude: longitude,
	}, nil
}

// ReconstituteGeoPoint creates a GeoPoint from trusted persisted coordinates.
func ReconstituteGeoPoint(latitude, longitude float64) GeoPoint {
	return GeoPoint{
		latitude:  latitude,
		longitude: longitude,
	}
}

// Latitude returns the latitude coordinate.
func (p GeoPoint) Latitude() float64 {
	return p.latitude
}

// Longitude returns the longitude coordinate.
func (p GeoPoint) Longitude() float64 {
	return p.longitude
}

// Equals compares two GeoPoint objects for equality.
func (p GeoPoint) Equals(other GeoPoint) bool {
	return p.latitude == other.latitude && p.longitude == other.longitude
}

// String returns a compact coordinate representation.
func (p GeoPoint) String() string {
	return fmt.Sprintf("%g,%g", p.latitude, p.longitude)
}

// MarshalJSON returns the geo point as an explicit object.
func (p GeoPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(geoPointJSON{
		Latitude:  p.latitude,
		Longitude: p.longitude,
	})
}

// UnmarshalJSON validates and normalizes a JSON geo point object.
func (p *GeoPoint) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return domain.ErrNullValue
	}

	var raw struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Latitude == nil {
		return ErrMissingLatitude
	}

	if raw.Longitude == nil {
		return ErrMissingLongitude
	}

	geoPoint, err := NewGeoPoint(*raw.Latitude, *raw.Longitude)
	if err != nil {
		return err
	}

	*p = geoPoint
	return nil
}
