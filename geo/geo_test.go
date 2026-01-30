package geo

import (
	"math"
	"testing"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
		errCode string
	}{
		// Valid coordinates
		{"valid center", 0, 0, false, ""},
		{"valid Maputo", -25.969, 32.573, false, ""},
		{"valid min lat", -90, 0, false, ""},
		{"valid max lat", 90, 0, false, ""},
		{"valid min lon", 0, -180, false, ""},
		{"valid max lon", 0, 180, false, ""},

		// Invalid coordinates
		{"lat too low", -91, 0, true, valerrors.CodeOutOfRange},
		{"lat too high", 91, 0, true, valerrors.CodeOutOfRange},
		{"lon too low", 0, -181, true, valerrors.CodeOutOfRange},
		{"lon too high", 0, 181, true, valerrors.CodeOutOfRange},
		{"both invalid", -100, 200, true, valerrors.CodeOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCoordinates(tt.lat, tt.lon)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCoordinates(%v, %v) error = %v, wantErr %v", tt.lat, tt.lon, err, tt.wantErr)
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

func TestValidateInMozambique(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
		errCode string
	}{
		// Valid Mozambique locations
		{"Maputo center", -25.969, 32.573, false, ""},
		{"Beira", -19.84, 34.84, false, ""},
		{"Nampula", -15.12, 39.27, false, ""},
		{"Mozambique north", -11.0, 35.0, false, ""},
		{"Mozambique south", -26.0, 32.5, false, ""},

		// Outside Mozambique
		{"South Africa", -29.0, 24.0, true, valerrors.CodeOutsideServiceArea},
		{"Tanzania", -6.0, 35.0, true, valerrors.CodeOutsideServiceArea},
		{"Madagascar", -19.0, 47.0, true, valerrors.CodeOutsideServiceArea},
		{"Zimbabwe", -19.0, 29.0, true, valerrors.CodeOutsideServiceArea},
		{"Atlantic Ocean", -25.0, 15.0, true, valerrors.CodeOutsideServiceArea},

		// Invalid coordinates (checked first)
		{"invalid lat", -100, 35, true, valerrors.CodeOutOfRange},
		{"invalid lon", -20, 200, true, valerrors.CodeOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInMozambique(tt.lat, tt.lon)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInMozambique(%v, %v) error = %v, wantErr %v", tt.lat, tt.lon, err, tt.wantErr)
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

func TestValidateServiceArea(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lon     float64
		area    string
		wantErr bool
		errCode string
	}{
		// Maputo service area
		{"Maputo center", -25.95, 32.5, "maputo", false, ""},
		{"Maputo edge", -26.0, 32.4, "maputo", false, ""},

		// Matola service area
		{"Matola center", -25.95, 32.4, "matola", false, ""},

		// Beira service area
		{"Beira center", -19.8, 34.85, "beira", false, ""},

		// Outside specific service area
		{"Maputo location in Beira area", -25.95, 32.5, "beira", true, valerrors.CodeOutsideServiceArea},
		{"Beira location in Maputo area", -19.8, 34.85, "maputo", true, valerrors.CodeOutsideServiceArea},

		// Invalid area name
		{"invalid area", -25.95, 32.5, "invalid", true, valerrors.CodeInvalidOption},
		{"empty area", -25.95, 32.5, "", true, valerrors.CodeInvalidOption},

		// Invalid coordinates
		{"invalid lat", -100, 32.5, "maputo", true, valerrors.CodeOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServiceArea(tt.lat, tt.lon, tt.area)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServiceArea(%v, %v, %q) error = %v, wantErr %v", tt.lat, tt.lon, tt.area, err, tt.wantErr)
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

func TestValidateAnyServiceArea(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
	}{
		// In service areas
		{"in Maputo", -25.95, 32.5, false},
		{"in Matola", -25.95, 32.4, false},
		{"in Beira", -19.8, 34.85, false},

		// Outside all service areas
		{"Nampula (not active)", -15.12, 39.27, true},
		{"middle of Mozambique", -20.0, 35.0, true},
		{"outside country", -30.0, 30.0, true},

		// Invalid coordinates
		{"invalid coordinates", -100, 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAnyServiceArea(tt.lat, tt.lon)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAnyServiceArea(%v, %v) error = %v, wantErr %v", tt.lat, tt.lon, err, tt.wantErr)
			}
		})
	}
}

func TestGetServiceAreas(t *testing.T) {
	areas := GetServiceAreas()

	if len(areas) != 3 {
		t.Errorf("GetServiceAreas() returned %d areas, want 3", len(areas))
	}

	expected := map[string]bool{"maputo": true, "matola": true, "beira": true}
	for _, area := range areas {
		if !expected[area] {
			t.Errorf("Unexpected area: %s", area)
		}
	}
}

func TestGetServiceArea(t *testing.T) {
	tests := []struct {
		name     string
		area     string
		wantNil  bool
		wantName string
	}{
		{"maputo", "maputo", false, "Maputo"},
		{"matola", "matola", false, "Matola"},
		{"beira", "beira", false, "Beira"},
		{"invalid", "invalid", true, ""},
		{"empty", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := GetServiceArea(tt.area)
			if tt.wantNil {
				if sa != nil {
					t.Errorf("GetServiceArea(%q) = %v, want nil", tt.area, sa)
				}
			} else {
				if sa == nil {
					t.Fatalf("GetServiceArea(%q) = nil, want non-nil", tt.area)
				}
				if sa.Name != tt.wantName {
					t.Errorf("GetServiceArea(%q).Name = %v, want %v", tt.area, sa.Name, tt.wantName)
				}
			}
		})
	}
}

func TestFindServiceArea(t *testing.T) {
	t.Run("in Maputo only", func(t *testing.T) {
		// Point in Maputo but outside Matola (lon > 32.5)
		got := FindServiceArea(-25.85, 32.6)
		if got != "maputo" {
			t.Errorf("FindServiceArea(-25.85, 32.6) = %v, want maputo", got)
		}
	})

	t.Run("in overlapping Maputo/Matola area", func(t *testing.T) {
		// Point in both Maputo and Matola - either is acceptable
		got := FindServiceArea(-25.95, 32.4)
		if got != "maputo" && got != "matola" {
			t.Errorf("FindServiceArea(-25.95, 32.4) = %v, want maputo or matola", got)
		}
	})

	t.Run("in Beira", func(t *testing.T) {
		got := FindServiceArea(-19.8, 34.85)
		if got != "beira" {
			t.Errorf("FindServiceArea(-19.8, 34.85) = %v, want beira", got)
		}
	})

	t.Run("outside all", func(t *testing.T) {
		got := FindServiceArea(-15.0, 39.0)
		if got != "" {
			t.Errorf("FindServiceArea(-15.0, 39.0) = %v, want empty", got)
		}
	})

	t.Run("invalid coords", func(t *testing.T) {
		got := FindServiceArea(-100, 200)
		if got != "" {
			t.Errorf("FindServiceArea(-100, 200) = %v, want empty", got)
		}
	})
}

func TestIsInMozambique(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lon  float64
		want bool
	}{
		{"Maputo", -25.969, 32.573, true},
		{"South Africa", -29.0, 24.0, false},
		{"invalid", -100, 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInMozambique(tt.lat, tt.lon); got != tt.want {
				t.Errorf("IsInMozambique(%v, %v) = %v, want %v", tt.lat, tt.lon, got, tt.want)
			}
		})
	}
}

func TestIsInServiceArea(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lon  float64
		want bool
	}{
		{"in Maputo", -25.95, 32.5, true},
		{"outside service areas", -15.0, 39.0, false},
		{"invalid", -100, 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInServiceArea(tt.lat, tt.lon); got != tt.want {
				t.Errorf("IsInServiceArea(%v, %v) = %v, want %v", tt.lat, tt.lon, got, tt.want)
			}
		})
	}
}

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name       string
		lat1, lon1 float64
		lat2, lon2 float64
		wantMin    float64
		wantMax    float64
		wantErr    bool
	}{
		// Same point
		{"same point", -25.969, 32.573, -25.969, 32.573, 0, 0.001, false},

		// Maputo to Matola (roughly 10-15 km)
		{"Maputo to Matola", -25.969, 32.573, -25.962, 32.459, 10, 15, false},

		// Maputo to Beira (roughly 1000-1200 km)
		{"Maputo to Beira", -25.969, 32.573, -19.84, 34.84, 700, 800, false},

		// Invalid coordinates
		{"invalid lat1", -100, 32.573, -25.969, 32.573, 0, 0, true},
		{"invalid lon2", -25.969, 32.573, -25.969, 200, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got < tt.wantMin || got > tt.wantMax {
					t.Errorf("CalculateDistance() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
				}
			}
		})
	}
}

func TestMozambiqueBounds(t *testing.T) {
	// Verify the bounds constants are reasonable
	if MozambiqueMinLat >= MozambiqueMaxLat {
		t.Error("MozambiqueMinLat should be less than MozambiqueMaxLat")
	}
	if MozambiqueMinLon >= MozambiqueMaxLon {
		t.Error("MozambiqueMinLon should be less than MozambiqueMaxLon")
	}

	// Verify known Mozambique cities are within bounds
	cities := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"Maputo", -25.969, 32.573},
		{"Beira", -19.84, 34.84},
		{"Nampula", -15.12, 39.27},
		{"Pemba", -12.97, 40.52},
		{"Lichinga", -13.31, 35.24},
	}

	for _, city := range cities {
		if city.lat < MozambiqueMinLat || city.lat > MozambiqueMaxLat {
			t.Errorf("%s lat %v is outside Mozambique bounds", city.name, city.lat)
		}
		if city.lon < MozambiqueMinLon || city.lon > MozambiqueMaxLon {
			t.Errorf("%s lon %v is outside Mozambique bounds", city.name, city.lon)
		}
	}
}

func TestServiceAreaBounds(t *testing.T) {
	// Verify all service areas have valid bounds
	for name, sa := range serviceAreas {
		if sa.MinLat >= sa.MaxLat {
			t.Errorf("Service area %s: MinLat (%v) should be less than MaxLat (%v)", name, sa.MinLat, sa.MaxLat)
		}
		if sa.MinLon >= sa.MaxLon {
			t.Errorf("Service area %s: MinLon (%v) should be less than MaxLon (%v)", name, sa.MinLon, sa.MaxLon)
		}

		// Verify service area is within Mozambique
		centerLat := (sa.MinLat + sa.MaxLat) / 2
		centerLon := (sa.MinLon + sa.MaxLon) / 2
		if !IsInMozambique(centerLat, centerLon) {
			t.Errorf("Service area %s center is not within Mozambique", name)
		}
	}
}

func TestCalculateDistance_Precision(t *testing.T) {
	// Test distance calculation precision with known values
	// Distance from equator at prime meridian to 1 degree north should be ~111 km
	dist, err := CalculateDistance(0, 0, 1, 0)
	if err != nil {
		t.Fatalf("CalculateDistance() error = %v", err)
	}

	// Earth circumference is ~40,075 km, so 1 degree should be ~111 km
	expected := 111.0
	tolerance := 2.0 // Allow 2km tolerance

	if math.Abs(dist-expected) > tolerance {
		t.Errorf("CalculateDistance(0,0,1,0) = %v km, want ~%v km (tolerance %v)", dist, expected, tolerance)
	}
}
