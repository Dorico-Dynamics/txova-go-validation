package phone

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// Valid formats
		{"local format", "841234567", true},
		{"international with plus", "+258841234567", true},
		{"international without plus", "258841234567", true},
		{"with 00 prefix", "00258841234567", true},
		{"with spaces", "84 123 4567", true},
		{"with dashes", "84-123-4567", true},
		{"with dots", "84.123.4567", true},
		{"with mixed separators", "84 123-4567", true},

		// All valid prefixes
		{"prefix 82 (Vodacom)", "821234567", true},
		{"prefix 83 (Movitel)", "831234567", true},
		{"prefix 84 (Vodacom)", "841234567", true},
		{"prefix 85 (Vodacom)", "851234567", true},
		{"prefix 86 (Movitel)", "861234567", true},
		{"prefix 87 (Tmcel)", "871234567", true},

		// Invalid formats
		{"empty string", "", false},
		{"too short", "8412345", false},
		{"too long local", "8412345678", false},
		{"invalid prefix 80", "801234567", false},
		{"invalid prefix 81", "811234567", false},
		{"invalid prefix 88", "881234567", false},
		{"invalid prefix 89", "891234567", false},
		{"letters only", "abcdefghi", false},
		{"mixed letters numbers", "84abc4567", false},
		{"wrong country code", "+254841234567", false},
		{"wrong country code no plus", "254841234567", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.input); got != tt.want {
				t.Errorf("Validate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// Valid formats normalize to +258XXXXXXXXX
		{"local format", "841234567", "+258841234567", false},
		{"international with plus", "+258841234567", "+258841234567", false},
		{"international without plus", "258841234567", "+258841234567", false},
		{"with 00 prefix", "00258841234567", "+258841234567", false},
		{"with spaces", "84 123 4567", "+258841234567", false},
		{"with dashes", "84-123-4567", "+258841234567", false},
		{"with dots", "84.123.4567", "+258841234567", false},
		{"with parentheses", "(84) 123 4567", "+258841234567", false},
		{"with mixed separators", "84 123-4567", "+258841234567", false},
		{"leading spaces", "  841234567", "+258841234567", false},
		{"trailing spaces", "841234567  ", "+258841234567", false},
		{"plus with spaces", "+ 258 84 123 4567", "+258841234567", false},

		// All valid prefixes
		{"prefix 82", "821234567", "+258821234567", false},
		{"prefix 83", "831234567", "+258831234567", false},
		{"prefix 84", "841234567", "+258841234567", false},
		{"prefix 85", "851234567", "+258851234567", false},
		{"prefix 86", "861234567", "+258861234567", false},
		{"prefix 87", "871234567", "+258871234567", false},

		// Invalid formats
		{"empty string", "", "", true},
		{"only spaces", "   ", "", true},
		{"only plus", "+", "", true},
		{"too short", "8412345", "", true},
		{"too long local", "8412345678", "", true},
		{"invalid prefix", "801234567", "", true},
		{"wrong country code", "+254841234567", "", true},
		{"letters only", "abcdefghi", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Normalize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIdentifyOperator(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Vodacom (82, 84, 85)
		{"prefix 82 Vodacom", "821234567", "Vodacom"},
		{"prefix 84 Vodacom", "841234567", "Vodacom"},
		{"prefix 85 Vodacom", "851234567", "Vodacom"},
		{"prefix 82 international", "+258821234567", "Vodacom"},

		// Movitel (83, 86)
		{"prefix 83 Movitel", "831234567", "Movitel"},
		{"prefix 86 Movitel", "861234567", "Movitel"},
		{"prefix 83 international", "+258831234567", "Movitel"},

		// Tmcel (87)
		{"prefix 87 Tmcel", "871234567", "Tmcel"},
		{"prefix 87 international", "+258871234567", "Tmcel"},

		// Invalid
		{"invalid number", "invalid", ""},
		{"empty string", "", ""},
		{"invalid prefix", "801234567", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IdentifyOperator(tt.input); got != tt.want {
				t.Errorf("IdentifyOperator(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"local 82", "821234567", "82"},
		{"local 83", "831234567", "83"},
		{"local 84", "841234567", "84"},
		{"local 85", "851234567", "85"},
		{"local 86", "861234567", "86"},
		{"local 87", "871234567", "87"},
		{"international", "+258841234567", "84"},
		{"with spaces", "84 123 4567", "84"},
		{"invalid", "invalid", ""},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPrefix(tt.input); got != tt.want {
				t.Errorf("GetPrefix(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsVodacom(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"prefix 82", "821234567", true},
		{"prefix 84", "841234567", true},
		{"prefix 85", "851234567", true},
		{"prefix 83 (Movitel)", "831234567", false},
		{"prefix 86 (Movitel)", "861234567", false},
		{"prefix 87 (Tmcel)", "871234567", false},
		{"invalid", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVodacom(tt.input); got != tt.want {
				t.Errorf("IsVodacom(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsMovitel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"prefix 83", "831234567", true},
		{"prefix 86", "861234567", true},
		{"prefix 82 (Vodacom)", "821234567", false},
		{"prefix 84 (Vodacom)", "841234567", false},
		{"prefix 85 (Vodacom)", "851234567", false},
		{"prefix 87 (Tmcel)", "871234567", false},
		{"invalid", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMovitel(tt.input); got != tt.want {
				t.Errorf("IsMovitel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsTmcel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"prefix 87", "871234567", true},
		{"prefix 82 (Vodacom)", "821234567", false},
		{"prefix 83 (Movitel)", "831234567", false},
		{"invalid", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTmcel(tt.input); got != tt.want {
				t.Errorf("IsTmcel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalize_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// Edge cases for 00 prefix handling
		{"00 prefix correct", "00258841234567", "+258841234567", false},
		{"00 prefix wrong country", "00254841234567", "", true},

		// Edge cases for plus handling
		{"plus only", "+", "", true},
		{"plus with wrong country", "+254841234567", "", true},
		{"plus with spaces before digits", "+  258841234567", "+258841234567", false},

		// Boundary cases
		{"exactly 9 digits valid", "841234567", "+258841234567", false},
		{"exactly 12 digits with 258", "258841234567", "+258841234567", false},
		{"exactly 14 digits with 00258", "00258841234567", "+258841234567", false},

		// Numbers that look valid but aren't
		{"11 digits", "84123456789", "", true},
		{"13 digits", "2588412345678", "", true},
		{"10 digits", "8412345678", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Normalize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
