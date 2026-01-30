// Package ride provides ride-related validation for distance, fare, PIN, and locations.
package ride

import (
	"github.com/Dorico-Dynamics/txova-go-types/geo"
	"github.com/Dorico-Dynamics/txova-go-types/money"
	"github.com/Dorico-Dynamics/txova-go-types/ride"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

// Distance constraints in kilometers.
const (
	MinDistanceKM = 0.5
	MaxDistanceKM = 200.0
)

// Fare constraints in centavos (MZN * 100).
const (
	MinFareCentavos = 5000    // 50 MZN
	MaxFareCentavos = 5000000 // 50,000 MZN
)

// Minimum separation between pickup and dropoff in kilometers.
const MinPickupDropoffSeparationKM = 0.1

// ValidatePIN validates a 4-digit ride verification PIN.
// Uses the types library which enforces no sequential (1234, 4321) or repeated (1111) patterns.
func ValidatePIN(input string) error {
	_, err := ride.ParsePIN(input)
	if err != nil {
		return valerrors.InvalidFormatWithValue("pin", "4-digit PIN (no sequential or repeated)", input)
	}
	return nil
}

// ValidateDistance validates that a ride distance is within acceptable range.
func ValidateDistance(km float64) error {
	if km < MinDistanceKM || km > MaxDistanceKM {
		return valerrors.OutOfRangeWithValue("distance", MinDistanceKM, MaxDistanceKM, km)
	}
	return nil
}

// ValidateFare validates that a fare amount (in centavos) is within acceptable range.
func ValidateFare(centavos int64) error {
	if centavos < MinFareCentavos || centavos > MaxFareCentavos {
		return valerrors.OutOfRangeWithValue("fare", MinFareCentavos, MaxFareCentavos, centavos)
	}
	return nil
}

// ValidateFareMoney validates a Money amount is within acceptable fare range.
func ValidateFareMoney(m money.Money) error {
	return ValidateFare(m.Centavos())
}

// ValidatePickupDropoff validates that pickup and dropoff locations are sufficiently separated.
// Returns an error if the locations are too close together.
func ValidatePickupDropoff(pickupLat, pickupLon, dropoffLat, dropoffLon float64) error {
	pickup, err := geo.NewLocation(pickupLat, pickupLon)
	if err != nil {
		return valerrors.InvalidFormatWithValue("pickup", "valid coordinates", err.Error())
	}

	dropoff, err := geo.NewLocation(dropoffLat, dropoffLon)
	if err != nil {
		return valerrors.InvalidFormatWithValue("dropoff", "valid coordinates", err.Error())
	}

	distance := geo.DistanceKM(pickup, dropoff)
	if distance < MinPickupDropoffSeparationKM {
		return valerrors.New("pickup_dropoff", valerrors.CodeOutOfRange,
			"pickup and dropoff must be at least 100 meters apart")
	}

	return nil
}

// ValidatePickupDropoffLocations validates pickup and dropoff using Location types.
func ValidatePickupDropoffLocations(pickup, dropoff geo.Location) error {
	if pickup.IsZero() {
		return valerrors.Required("pickup")
	}
	if dropoff.IsZero() {
		return valerrors.Required("dropoff")
	}

	distance := geo.DistanceKM(pickup, dropoff)
	if distance < MinPickupDropoffSeparationKM {
		return valerrors.New("pickup_dropoff", valerrors.CodeOutOfRange,
			"pickup and dropoff must be at least 100 meters apart")
	}

	return nil
}

// IsValidPIN returns true if the PIN is valid.
func IsValidPIN(input string) bool {
	return ValidatePIN(input) == nil
}

// IsValidDistance returns true if the distance is within acceptable range.
func IsValidDistance(km float64) bool {
	return ValidateDistance(km) == nil
}

// IsValidFare returns true if the fare (in centavos) is within acceptable range.
func IsValidFare(centavos int64) bool {
	return ValidateFare(centavos) == nil
}

// CalculateEstimatedFare calculates an estimated fare based on distance.
// This is a simplified calculation for validation purposes.
// Returns fare in centavos.
func CalculateEstimatedFare(distanceKM float64, baseFareCentavos, perKMCentavos int64) int64 {
	return baseFareCentavos + int64(distanceKM*float64(perKMCentavos))
}
