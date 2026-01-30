// Package vehicle provides Mozambique vehicle validation including license plates and year.
package vehicle

import (
	"errors"
	"time"

	"github.com/Dorico-Dynamics/txova-go-types/vehicle"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

// Vehicle year constraints.
const (
	MinVehicleYear = 2010
)

// ValidatePlate validates a Mozambique license plate format.
// Accepts both standard (AAA-NNN-LL) and old (LL-NN-NN) formats.
func ValidatePlate(input string) error {
	_, err := vehicle.ParseLicensePlate(input)
	if err != nil {
		if errors.Is(err, vehicle.ErrInvalidProvinceCode) {
			return valerrors.InvalidFormat("plate", "valid Mozambique province code")
		}
		return valerrors.InvalidFormatWithValue("plate", "AAA-NNN-LL or LL-NN-NN", input)
	}
	return nil
}

// NormalizePlate normalizes a license plate to standard format with dashes.
// Returns the normalized plate string or an error if invalid.
func NormalizePlate(input string) (string, error) {
	plate, err := vehicle.ParseLicensePlate(input)
	if err != nil {
		if errors.Is(err, vehicle.ErrInvalidProvinceCode) {
			return "", valerrors.InvalidFormat("plate", "valid Mozambique province code")
		}
		return "", valerrors.InvalidFormatWithValue("plate", "AAA-NNN-LL or LL-NN-NN", input)
	}
	return plate.String(), nil
}

// ValidateYear validates a vehicle year is within acceptable range.
// Year must be between MinVehicleYear (2010) and current year + 1.
func ValidateYear(year int) error {
	maxYear := time.Now().Year() + 1
	if year < MinVehicleYear || year > maxYear {
		return valerrors.OutOfRangeWithValue("year", MinVehicleYear, maxYear, year)
	}
	return nil
}

// GetProvince extracts the province code from a license plate.
// Returns the province code string or empty if invalid.
func GetProvince(input string) string {
	plate, err := vehicle.ParseLicensePlate(input)
	if err != nil {
		return ""
	}
	return plate.Province().String()
}

// GetProvinceName returns the full province name for a license plate.
// Returns empty string if the plate is invalid.
func GetProvinceName(input string) string {
	plate, err := vehicle.ParseLicensePlate(input)
	if err != nil {
		return ""
	}
	return plate.Province().ProvinceName()
}

// IsStandardFormat returns true if the plate is in standard format (AAA-NNN-LL).
func IsStandardFormat(input string) bool {
	plate, err := vehicle.ParseLicensePlate(input)
	if err != nil {
		return false
	}
	return plate.IsStandardFormat()
}

// IsOldFormat returns true if the plate is in old format (LL-NN-NN).
func IsOldFormat(input string) bool {
	plate, err := vehicle.ParseLicensePlate(input)
	if err != nil {
		return false
	}
	return plate.IsOldFormat()
}

// IsValidPlate returns true if the input is a valid Mozambique license plate.
func IsValidPlate(input string) bool {
	return ValidatePlate(input) == nil
}

// IsValidYear returns true if the year is within acceptable range.
func IsValidYear(year int) bool {
	return ValidateYear(year) == nil
}
