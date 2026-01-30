// Package geo provides geographic validation for Mozambique locations and service areas.
package geo

import (
	"github.com/Dorico-Dynamics/txova-go-types/geo"
	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

// Mozambique bounding box coordinates.
const (
	MozambiqueMinLat = -26.9
	MozambiqueMaxLat = -10.3
	MozambiqueMinLon = 30.2
	MozambiqueMaxLon = 41.0
)

// ServiceArea represents a geographic service area.
type ServiceArea struct {
	Name   string
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

// Predefined service areas for Txova operations.
var serviceAreas = map[string]ServiceArea{
	"maputo": {
		Name:   "Maputo",
		MinLat: -26.1,
		MaxLat: -25.8,
		MinLon: 32.3,
		MaxLon: 32.7,
	},
	"matola": {
		Name:   "Matola",
		MinLat: -26.0,
		MaxLat: -25.9,
		MinLon: 32.3,
		MaxLon: 32.5,
	},
	"beira": {
		Name:   "Beira",
		MinLat: -19.9,
		MaxLat: -19.7,
		MinLon: 34.8,
		MaxLon: 34.9,
	},
}

// ValidateCoordinates checks if latitude and longitude are within valid global ranges.
// Latitude must be between -90 and 90, longitude between -180 and 180.
func ValidateCoordinates(lat, lon float64) error {
	if lat < geo.MinLatitude || lat > geo.MaxLatitude {
		return valerrors.OutOfRangeWithValue("latitude", geo.MinLatitude, geo.MaxLatitude, lat)
	}
	if lon < geo.MinLongitude || lon > geo.MaxLongitude {
		return valerrors.OutOfRangeWithValue("longitude", geo.MinLongitude, geo.MaxLongitude, lon)
	}
	return nil
}

// ValidateInMozambique checks if coordinates are within Mozambique's borders.
func ValidateInMozambique(lat, lon float64) error {
	// First validate global ranges
	if err := ValidateCoordinates(lat, lon); err != nil {
		return err
	}

	// Check Mozambique bounds
	if lat < MozambiqueMinLat || lat > MozambiqueMaxLat ||
		lon < MozambiqueMinLon || lon > MozambiqueMaxLon {
		return valerrors.OutsideServiceAreaWithValue("location", lat, lon)
	}

	return nil
}

// ValidateServiceArea checks if coordinates are within a specific service area.
// The area parameter should be one of: "maputo", "matola", "beira".
func ValidateServiceArea(lat, lon float64, area string) error {
	// First validate global ranges
	if err := ValidateCoordinates(lat, lon); err != nil {
		return err
	}

	// Get service area bounds
	sa, exists := serviceAreas[area]
	if !exists {
		return valerrors.InvalidOptionWithValue("area", GetServiceAreas(), area)
	}

	// Check if within service area
	if lat < sa.MinLat || lat > sa.MaxLat || lon < sa.MinLon || lon > sa.MaxLon {
		return valerrors.OutsideServiceAreaWithValue("location", lat, lon)
	}

	return nil
}

// ValidateAnyServiceArea checks if coordinates are within any active service area.
func ValidateAnyServiceArea(lat, lon float64) error {
	// First validate global ranges
	if err := ValidateCoordinates(lat, lon); err != nil {
		return err
	}

	// Check all service areas
	for _, sa := range serviceAreas {
		if lat >= sa.MinLat && lat <= sa.MaxLat && lon >= sa.MinLon && lon <= sa.MaxLon {
			return nil
		}
	}

	return valerrors.OutsideServiceAreaWithValue("location", lat, lon)
}

// GetServiceAreas returns a list of all active service area names.
func GetServiceAreas() []string {
	areas := make([]string, 0, len(serviceAreas))
	for name := range serviceAreas {
		areas = append(areas, name)
	}
	return areas
}

// GetServiceArea returns the service area configuration for a given area name.
// Returns nil if the area doesn't exist.
func GetServiceArea(name string) *ServiceArea {
	sa, exists := serviceAreas[name]
	if !exists {
		return nil
	}
	return &sa
}

// FindServiceArea returns the name of the service area containing the coordinates.
// Returns empty string if not in any service area.
func FindServiceArea(lat, lon float64) string {
	for name, sa := range serviceAreas {
		if lat >= sa.MinLat && lat <= sa.MaxLat && lon >= sa.MinLon && lon <= sa.MaxLon {
			return name
		}
	}
	return ""
}

// IsInMozambique returns true if the coordinates are within Mozambique.
func IsInMozambique(lat, lon float64) bool {
	return ValidateInMozambique(lat, lon) == nil
}

// IsInServiceArea returns true if the coordinates are within any active service area.
func IsInServiceArea(lat, lon float64) bool {
	return ValidateAnyServiceArea(lat, lon) == nil
}

// CalculateDistance returns the distance in kilometers between two points.
// Uses the Haversine formula via the types library.
func CalculateDistance(lat1, lon1, lat2, lon2 float64) (float64, error) {
	loc1, err := geo.NewLocation(lat1, lon1)
	if err != nil {
		return 0, err
	}
	loc2, err := geo.NewLocation(lat2, lon2)
	if err != nil {
		return 0, err
	}
	return geo.DistanceKM(loc1, loc2), nil
}
