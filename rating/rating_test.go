package rating

import (
	"strings"
	"testing"

	"github.com/Dorico-Dynamics/txova-go-types/rating"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

func TestValidateRating(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
		errCode string
	}{
		// Valid ratings
		{"rating 1", 1, false, ""},
		{"rating 2", 2, false, ""},
		{"rating 3", 3, false, ""},
		{"rating 4", 4, false, ""},
		{"rating 5", 5, false, ""},

		// Invalid ratings
		{"rating 0", 0, true, valerrors.CodeOutOfRange},
		{"rating -1", -1, true, valerrors.CodeOutOfRange},
		{"rating 6", 6, true, valerrors.CodeOutOfRange},
		{"rating 10", 10, true, valerrors.CodeOutOfRange},
		{"rating 100", 100, true, valerrors.CodeOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRating(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRating(%d) error = %v, wantErr %v", tt.value, err, tt.wantErr)
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

func TestValidateReviewText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		// Valid lengths
		{"empty string", "", false},
		{"short text", "Great driver!", false},
		{"medium text", "The driver was very professional and the car was clean.", false},
		{"max length", strings.Repeat("a", MaxReviewLength), false},

		// Invalid lengths
		{"too long", strings.Repeat("a", MaxReviewLength+1), true},
		{"way too long", strings.Repeat("a", 1000), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReviewText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateReviewText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateReviewText_Unicode(t *testing.T) {
	// Test with Unicode characters to ensure we count runes, not bytes
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{"portuguese chars", "Excelente motorista! Muito obrigado.", false},
		{"emojis count as one", strings.Repeat("😀", MaxReviewLength), false},
		{"emojis too many", strings.Repeat("😀", MaxReviewLength+1), true},
		{"chinese chars", strings.Repeat("中", MaxReviewLength), false},
		{"chinese too many", strings.Repeat("中", MaxReviewLength+1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReviewText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateReviewText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeReviewText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// HTML stripping
		{"strip simple tag", "<b>bold</b>", "bold"},
		{"strip multiple tags", "<p>Hello</p><br/><span>World</span>", "HelloWorld"},
		{"strip script tag", "<script>alert('xss')</script>Hello", "alert('xss')Hello"},
		{"no tags", "plain text", "plain text"},

		// Whitespace normalization
		{"multiple spaces", "hello    world", "hello world"},
		{"tabs", "hello\t\tworld", "hello world"},
		{"newlines", "hello\n\nworld", "hello world"},
		{"mixed whitespace", "hello  \t\n  world", "hello world"},

		// Trimming
		{"leading spaces", "   hello", "hello"},
		{"trailing spaces", "hello   ", "hello"},
		{"both ends", "   hello   ", "hello"},

		// Combined
		{"html and whitespace", "  <b>hello</b>   <i>world</i>  ", "hello world"},
		{"complex", "  <div>Hello,   </div>\n<p>World!</p>  ", "Hello, World!"},

		// Edge cases
		{"empty string", "", ""},
		{"only whitespace", "   ", ""},
		{"only tags", "<br/><hr/>", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeReviewText(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeReviewText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCheckProfanity(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		// Clean text
		{"clean text", "The driver was excellent!", false},
		{"professional review", "Very professional service, highly recommended.", false},
		{"portuguese clean", "Motorista muito bom, obrigado!", false},

		// English profanity
		{"english profanity", "This was shit service", true},
		{"english profanity 2", "What the fuck", true},
		{"english embedded", "This is bullshit", true},

		// Portuguese profanity
		{"portuguese profanity", "Que merda de serviço", true},
		{"portuguese profanity 2", "Vai à porra", true},
		{"portuguese fdp", "Este fdp não sabe conduzir", true},

		// Case insensitive
		{"uppercase", "SHIT", true},
		{"mixed case", "ShIt", true},

		// Edge cases
		{"empty", "", false},
		{"just spaces", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckProfanity(tt.text)
			if got != tt.want {
				t.Errorf("CheckProfanity(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestIsValidRating(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"valid 1", 1, true},
		{"valid 5", 5, true},
		{"invalid 0", 0, false},
		{"invalid 6", 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRating(tt.value); got != tt.want {
				t.Errorf("IsValidRating(%d) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsValidReviewText(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{"valid empty", "", true},
		{"valid short", "Great!", true},
		{"invalid too long", strings.Repeat("a", MaxReviewLength+1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidReviewText(tt.text); got != tt.want {
				t.Errorf("IsValidReviewText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAndSanitizeReview(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"valid with html", "<b>Great</b> driver!", "Great driver!", false},
		{"valid with spaces", "  Hello   World  ", "Hello World", false},
		{"too long after sanitize", strings.Repeat("a", MaxReviewLength+10), "", true},
		{"empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAndSanitizeReview(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndSanitizeReview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ValidateAndSanitizeReview() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessReview(t *testing.T) {
	t.Run("clean review", func(t *testing.T) {
		result, err := ProcessReview("Great driver!")
		if err != nil {
			t.Fatalf("ProcessReview() error = %v", err)
		}
		if result.Text != "Great driver!" {
			t.Errorf("Text = %q, want 'Great driver!'", result.Text)
		}
		if result.HasProfanity {
			t.Error("HasProfanity = true, want false")
		}
		if result.RequiresReview {
			t.Error("RequiresReview = true, want false")
		}
	})

	t.Run("review with html", func(t *testing.T) {
		result, err := ProcessReview("<b>Great</b> service!")
		if err != nil {
			t.Fatalf("ProcessReview() error = %v", err)
		}
		if result.Text != "Great service!" {
			t.Errorf("Text = %q, want 'Great service!'", result.Text)
		}
		if result.OriginalLength != 21 {
			t.Errorf("OriginalLength = %d, want 21", result.OriginalLength)
		}
		if result.SanitizedLength != 14 {
			t.Errorf("SanitizedLength = %d, want 14", result.SanitizedLength)
		}
	})

	t.Run("review with profanity", func(t *testing.T) {
		result, err := ProcessReview("This was shit")
		if err != nil {
			t.Fatalf("ProcessReview() error = %v", err)
		}
		if !result.HasProfanity {
			t.Error("HasProfanity = false, want true")
		}
		if !result.RequiresReview {
			t.Error("RequiresReview = false, want true")
		}
	})

	t.Run("review too long", func(t *testing.T) {
		_, err := ProcessReview(strings.Repeat("a", MaxReviewLength+1))
		if err == nil {
			t.Error("ProcessReview() should return error for too long text")
		}
	})
}

func TestConstants(t *testing.T) {
	// Verify constants match PRD
	if MinReviewLength != 0 {
		t.Errorf("MinReviewLength = %d, want 0", MinReviewLength)
	}
	if MaxReviewLength != 500 {
		t.Errorf("MaxReviewLength = %d, want 500", MaxReviewLength)
	}

	// Verify rating constants from types library
	if rating.MinRating != 1 {
		t.Errorf("MinRating = %d, want 1", rating.MinRating)
	}
	if rating.MaxRating != 5 {
		t.Errorf("MaxRating = %d, want 5", rating.MaxRating)
	}
}
