// Package rating provides rating and review validation for the Txova platform.
package rating

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/Dorico-Dynamics/txova-go-types/rating"
	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

// Review text constraints.
const (
	MinReviewLength = 0
	MaxReviewLength = 500
)

// htmlTagPattern matches HTML tags for stripping.
var htmlTagPattern = regexp.MustCompile(`<[^>]*>`)

// profanityWords contains common profanity terms in Portuguese and English.
// This is a conservative list for flagging, not blocking.
var profanityWords = map[string]bool{
	// English common terms
	"fuck": true, "shit": true, "damn": true, "ass": true, "bitch": true,
	"bastard": true, "crap": true, "piss": true, "dick": true, "cock": true,
	// Portuguese common terms
	"merda": true, "porra": true, "caralho": true, "foda": true, "puta": true,
	"corno": true, "filho da puta": true, "fdp": true, "cabrão": true,
}

// ValidateRating validates that a rating value is within the 1-5 range.
func ValidateRating(value int) error {
	_, err := rating.NewRating(value)
	if err != nil {
		return valerrors.OutOfRangeWithValue("rating", rating.MinRating, rating.MaxRating, value)
	}
	return nil
}

// ValidateReviewText validates the length of review text.
// Text is optional (can be empty) but must not exceed MaxReviewLength characters.
func ValidateReviewText(text string) error {
	length := len([]rune(text)) // Count Unicode characters, not bytes
	if length > MaxReviewLength {
		return valerrors.TooLongWithValue("review", MaxReviewLength, length)
	}
	return nil
}

// SanitizeReviewText sanitizes review text by:
// - Stripping HTML tags
// - Normalizing whitespace (collapsing multiple spaces)
// - Trimming leading/trailing whitespace
func SanitizeReviewText(text string) string {
	// Strip HTML tags
	result := htmlTagPattern.ReplaceAllString(text, "")

	// Normalize whitespace
	result = normalizeWhitespace(result)

	// Trim
	result = strings.TrimSpace(result)

	return result
}

// normalizeWhitespace collapses multiple whitespace characters into a single space.
func normalizeWhitespace(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	inWhitespace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !inWhitespace {
				result.WriteRune(' ')
				inWhitespace = true
			}
		} else {
			result.WriteRune(r)
			inWhitespace = false
		}
	}
	return result.String()
}

// CheckProfanity checks if the text contains potential profanity.
// Returns true if profanity is detected, indicating the text should be flagged for moderation.
// This uses a conservative approach - it only flags, doesn't reject.
func CheckProfanity(text string) bool {
	lower := strings.ToLower(text)

	// Check for exact word matches and partial matches
	for word := range profanityWords {
		if strings.Contains(lower, word) {
			return true
		}
	}

	return false
}

// IsValidRating returns true if the rating is within the 1-5 range.
func IsValidRating(value int) bool {
	return ValidateRating(value) == nil
}

// IsValidReviewText returns true if the review text length is acceptable.
func IsValidReviewText(text string) bool {
	return ValidateReviewText(text) == nil
}

// ValidateAndSanitizeReview validates and sanitizes review text in one operation.
// Returns the sanitized text and any validation error.
func ValidateAndSanitizeReview(text string) (string, error) {
	sanitized := SanitizeReviewText(text)
	if err := ValidateReviewText(sanitized); err != nil {
		return "", err
	}
	return sanitized, nil
}

// ReviewResult contains the result of review validation and processing.
type ReviewResult struct {
	Text            string
	HasProfanity    bool
	RequiresReview  bool
	OriginalLength  int
	SanitizedLength int
}

// ProcessReview validates, sanitizes, and checks a review for profanity.
// Returns a ReviewResult with all processing information.
func ProcessReview(text string) (ReviewResult, error) {
	result := ReviewResult{
		OriginalLength: len([]rune(text)),
	}

	// Sanitize
	sanitized := SanitizeReviewText(text)
	result.Text = sanitized
	result.SanitizedLength = len([]rune(sanitized))

	// Validate
	if err := ValidateReviewText(sanitized); err != nil {
		return result, err
	}

	// Check profanity
	result.HasProfanity = CheckProfanity(sanitized)
	result.RequiresReview = result.HasProfanity

	return result, nil
}
