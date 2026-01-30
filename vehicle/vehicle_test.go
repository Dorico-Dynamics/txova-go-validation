package vehicle

import (
	"testing"
	"time"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

func TestValidatePlate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errCode string
	}{
		// Valid standard format (AAA-NNN-LL)
		{"standard format", "AAA-123-MC", false, ""},
		{"standard lowercase", "aaa-123-mc", false, ""},
		{"standard no dashes", "AAA123MC", false, ""},
		{"standard with spaces", "AAA 123 MC", false, ""},

		// Valid old format (LL-NN-NN)
		{"old format", "MC-12-34", false, ""},
		{"old format lowercase", "mc-12-34", false, ""},
		{"old format no dashes", "MC1234", false, ""},

		// All valid province codes
		{"province MC", "AAA-123-MC", false, ""},
		{"province MP", "AAA-123-MP", false, ""},
		{"province GZ", "AAA-123-GZ", false, ""},
		{"province IB", "AAA-123-IB", false, ""},
		{"province SF", "AAA-123-SF", false, ""},
		{"province MN", "AAA-123-MN", false, ""},
		{"province TT", "AAA-123-TT", false, ""},
		{"province ZB", "AAA-123-ZB", false, ""},
		{"province NP", "AAA-123-NP", false, ""},
		{"province CA", "AAA-123-CA", false, ""},
		{"province NS", "AAA-123-NS", false, ""},

		// Invalid formats
		{"empty string", "", true, valerrors.CodeInvalidFormat},
		{"invalid province", "AAA-123-XX", true, valerrors.CodeInvalidFormat},
		{"too short", "AA-12", true, valerrors.CodeInvalidFormat},
		{"random string", "invalid", true, valerrors.CodeInvalidFormat},
		{"numbers only", "12345678", true, valerrors.CodeInvalidFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePlate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePlate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
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

func TestNormalizePlate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// Standard format normalization
		{"standard format", "AAA-123-MC", "AAA-123-MC", false},
		{"standard lowercase", "aaa-123-mc", "AAA-123-MC", false},
		{"standard no dashes", "AAA123MC", "AAA-123-MC", false},
		{"standard with spaces", "AAA 123 MC", "AAA-123-MC", false},
		{"standard mixed case", "Aaa-123-Mc", "AAA-123-MC", false},

		// Old format normalization
		{"old format", "MC-12-34", "MC-12-34", false},
		{"old format lowercase", "mc-12-34", "MC-12-34", false},
		{"old format no dashes", "MC1234", "MC-12-34", false},
		{"old format with spaces", "MC 12 34", "MC-12-34", false},

		// Invalid
		{"empty string", "", "", true},
		{"invalid province", "AAA-123-XX", "", true},
		{"invalid format", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizePlate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizePlate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizePlate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateYear(t *testing.T) {
	currentYear := time.Now().Year()
	maxYear := currentYear + 1

	tests := []struct {
		name    string
		year    int
		wantErr bool
	}{
		// Valid years
		{"minimum year", MinVehicleYear, false},
		{"current year", currentYear, false},
		{"next year", maxYear, false},
		{"mid range", 2020, false},

		// Invalid years
		{"too old", 2009, true},
		{"year 2000", 2000, true},
		{"year 1990", 1990, true},
		{"too new", maxYear + 1, true},
		{"far future", 2050, true},
		{"zero", 0, true},
		{"negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateYear(tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateYear(%d) error = %v, wantErr %v", tt.year, err, tt.wantErr)
			}
		})
	}
}

func TestGetProvince(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"standard MC", "AAA-123-MC", "MC"},
		{"standard MP", "AAA-123-MP", "MP"},
		{"old format MC", "MC-12-34", "MC"},
		{"old format GZ", "GZ-99-01", "GZ"},
		{"invalid", "invalid", ""},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProvince(tt.input); got != tt.want {
				t.Errorf("GetProvince(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetProvinceName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"MC", "AAA-123-MC", "Maputo City"},
		{"MP", "AAA-123-MP", "Maputo Province"},
		{"GZ", "AAA-123-GZ", "Gaza"},
		{"IB", "AAA-123-IB", "Inhambane"},
		{"SF", "AAA-123-SF", "Sofala"},
		{"invalid", "invalid", ""},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProvinceName(tt.input); got != tt.want {
				t.Errorf("GetProvinceName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsStandardFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"standard format", "AAA-123-MC", true},
		{"old format", "MC-12-34", false},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStandardFormat(tt.input); got != tt.want {
				t.Errorf("IsStandardFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsOldFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"old format", "MC-12-34", true},
		{"standard format", "AAA-123-MC", false},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOldFormat(tt.input); got != tt.want {
				t.Errorf("IsOldFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidPlate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid standard", "AAA-123-MC", true},
		{"valid old", "MC-12-34", true},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPlate(tt.input); got != tt.want {
				t.Errorf("IsValidPlate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidYear(t *testing.T) {
	currentYear := time.Now().Year()

	tests := []struct {
		name string
		year int
		want bool
	}{
		{"valid current", currentYear, true},
		{"valid min", MinVehicleYear, true},
		{"invalid old", 2009, false},
		{"invalid future", 2050, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidYear(tt.year); got != tt.want {
				t.Errorf("IsValidYear(%d) = %v, want %v", tt.year, got, tt.want)
			}
		})
	}
}
