package ride

import (
	"testing"

	"github.com/Dorico-Dynamics/txova-go-types/geo"
	"github.com/Dorico-Dynamics/txova-go-types/money"
	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

func TestValidatePIN(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid PINs
		{"valid 4 digits", "7392", false},
		{"valid another", "4826", false},
		{"valid with zeros", "0392", false},

		// Invalid PINs - sequential
		{"sequential ascending", "1234", true},
		{"sequential descending", "4321", true},
		{"sequential 5678", "5678", true},
		{"sequential 8765", "8765", true},

		// Invalid PINs - repeated
		{"all same 1111", "1111", true},
		{"all same 2222", "2222", true},
		{"all same 0000", "0000", true},
		{"all same 9999", "9999", true},

		// Invalid PINs - format
		{"too short", "123", true},
		{"too long", "12345", true},
		{"empty", "", true},
		{"letters", "abcd", true},
		{"mixed", "12ab", true},
		{"with spaces", "12 34", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePIN(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePIN(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDistance(t *testing.T) {
	tests := []struct {
		name    string
		km      float64
		wantErr bool
		errCode string
	}{
		// Valid distances
		{"minimum", MinDistanceKM, false, ""},
		{"maximum", MaxDistanceKM, false, ""},
		{"mid range", 50.0, false, ""},
		{"just above min", 0.51, false, ""},
		{"just below max", 199.9, false, ""},

		// Invalid distances
		{"too short", 0.1, true, valerrors.CodeOutOfRange},
		{"zero", 0, true, valerrors.CodeOutOfRange},
		{"negative", -5, true, valerrors.CodeOutOfRange},
		{"too long", 250, true, valerrors.CodeOutOfRange},
		{"just below min", 0.49, true, valerrors.CodeOutOfRange},
		{"just above max", 200.1, true, valerrors.CodeOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDistance(tt.km)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDistance(%v) error = %v, wantErr %v", tt.km, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if ve, ok := err.(valerrors.ValidationError); ok {
					if ve.Code != tt.errCode {
						t.Errorf("error code = %v, want %v", ve.Code, tt.errCode)
					}
				}
			}
		})
	}
}

func TestValidateFare(t *testing.T) {
	tests := []struct {
		name     string
		centavos int64
		wantErr  bool
	}{
		// Valid fares
		{"minimum", MinFareCentavos, false},
		{"maximum", MaxFareCentavos, false},
		{"mid range", 100000, false}, // 1000 MZN
		{"just above min", 5001, false},
		{"just below max", 4999999, false},

		// Invalid fares
		{"too low", 4999, true},
		{"zero", 0, true},
		{"negative", -1000, true},
		{"too high", 5000001, true},
		{"way too high", 10000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFare(tt.centavos)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFare(%v) error = %v, wantErr %v", tt.centavos, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFareMoney(t *testing.T) {
	tests := []struct {
		name    string
		mzn     float64
		wantErr bool
	}{
		{"valid 100 MZN", 100, false},
		{"valid 1000 MZN", 1000, false},
		{"valid min 50 MZN", 50, false},
		{"too low 49 MZN", 49, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := money.FromMZN(tt.mzn)
			err := ValidateFareMoney(m)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFareMoney(%v MZN) error = %v, wantErr %v", tt.mzn, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePickupDropoff(t *testing.T) {
	// Maputo center
	maputoLat, maputoLon := -25.969, 32.573

	tests := []struct {
		name       string
		pickupLat  float64
		pickupLon  float64
		dropoffLat float64
		dropoffLon float64
		wantErr    bool
	}{
		// Valid - sufficient separation
		{"far apart", maputoLat, maputoLon, maputoLat + 0.01, maputoLon + 0.01, false},
		{"1km apart", maputoLat, maputoLon, maputoLat + 0.009, maputoLon, false},

		// Invalid - too close
		{"same point", maputoLat, maputoLon, maputoLat, maputoLon, true},
		{"very close", maputoLat, maputoLon, maputoLat + 0.0001, maputoLon + 0.0001, true},

		// Invalid coordinates
		{"invalid pickup lat", -100, maputoLon, maputoLat, maputoLon, true},
		{"invalid dropoff lon", maputoLat, maputoLon, maputoLat, 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePickupDropoff(tt.pickupLat, tt.pickupLon, tt.dropoffLat, tt.dropoffLon)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePickupDropoff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePickupDropoffLocations(t *testing.T) {
	maputo := geo.MustNewLocation(-25.969, 32.573)
	nearby := geo.MustNewLocation(-25.970, 32.574)
	farAway := geo.MustNewLocation(-25.980, 32.590)

	tests := []struct {
		name    string
		pickup  geo.Location
		dropoff geo.Location
		wantErr bool
		errCode string
	}{
		{"valid far apart", maputo, farAway, false, ""},
		{"valid nearby but OK", maputo, nearby, false, ""},
		{"same location", maputo, maputo, true, valerrors.CodeOutOfRange},
		{"zero pickup", geo.Location{}, farAway, true, valerrors.CodeRequired},
		{"zero dropoff", maputo, geo.Location{}, true, valerrors.CodeRequired},
		{"both zero", geo.Location{}, geo.Location{}, true, valerrors.CodeRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePickupDropoffLocations(tt.pickup, tt.dropoff)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePickupDropoffLocations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if ve, ok := err.(valerrors.ValidationError); ok {
					if ve.Code != tt.errCode {
						t.Errorf("error code = %v, want %v", ve.Code, tt.errCode)
					}
				}
			}
		})
	}
}

func TestIsValidPIN(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid", "7392", true},
		{"sequential", "1234", false},
		{"repeated", "1111", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPIN(tt.input); got != tt.want {
				t.Errorf("IsValidPIN(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidDistance(t *testing.T) {
	tests := []struct {
		name string
		km   float64
		want bool
	}{
		{"valid", 10, true},
		{"min", MinDistanceKM, true},
		{"max", MaxDistanceKM, true},
		{"too short", 0.1, false},
		{"too long", 300, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDistance(tt.km); got != tt.want {
				t.Errorf("IsValidDistance(%v) = %v, want %v", tt.km, got, tt.want)
			}
		})
	}
}

func TestIsValidFare(t *testing.T) {
	tests := []struct {
		name     string
		centavos int64
		want     bool
	}{
		{"valid", 10000, true},
		{"min", MinFareCentavos, true},
		{"max", MaxFareCentavos, true},
		{"too low", 1000, false},
		{"too high", 10000000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidFare(tt.centavos); got != tt.want {
				t.Errorf("IsValidFare(%v) = %v, want %v", tt.centavos, got, tt.want)
			}
		})
	}
}

func TestCalculateEstimatedFare(t *testing.T) {
	tests := []struct {
		name     string
		distance float64
		baseFare int64
		perKM    int64
		want     int64
	}{
		{"0km", 0, 5000, 1000, 5000},
		{"10km", 10, 5000, 1000, 15000},
		{"5km", 5, 3000, 500, 5500},
		{"1km", 1, 5000, 2000, 7000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateEstimatedFare(tt.distance, tt.baseFare, tt.perKM)
			if got != tt.want {
				t.Errorf("CalculateEstimatedFare(%v, %v, %v) = %v, want %v",
					tt.distance, tt.baseFare, tt.perKM, got, tt.want)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Verify constants are reasonable
	if MinDistanceKM <= 0 {
		t.Error("MinDistanceKM should be positive")
	}
	if MaxDistanceKM <= MinDistanceKM {
		t.Error("MaxDistanceKM should be greater than MinDistanceKM")
	}
	if MinFareCentavos <= 0 {
		t.Error("MinFareCentavos should be positive")
	}
	if MaxFareCentavos <= MinFareCentavos {
		t.Error("MaxFareCentavos should be greater than MinFareCentavos")
	}
	if MinPickupDropoffSeparationKM <= 0 {
		t.Error("MinPickupDropoffSeparationKM should be positive")
	}

	// Verify constants match PRD
	if MinDistanceKM != 0.5 {
		t.Errorf("MinDistanceKM = %v, want 0.5", MinDistanceKM)
	}
	if MaxDistanceKM != 200 {
		t.Errorf("MaxDistanceKM = %v, want 200", MaxDistanceKM)
	}
	if MinFareCentavos != 5000 {
		t.Errorf("MinFareCentavos = %v, want 5000 (50 MZN)", MinFareCentavos)
	}
	if MaxFareCentavos != 5000000 {
		t.Errorf("MaxFareCentavos = %v, want 5000000 (50,000 MZN)", MaxFareCentavos)
	}
}
