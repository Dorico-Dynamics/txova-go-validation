// Package sanitize provides input sanitization utilities for the Txova platform.
// All functions return new values and never modify the input.
package sanitize

import (
	"regexp"
	"strings"
	"unicode"
)

// htmlTagPattern matches HTML tags for stripping.
var htmlTagPattern = regexp.MustCompile(`<[^>]*>`)

// multiSpacePattern matches multiple consecutive whitespace characters.
var multiSpacePattern = regexp.MustCompile(`\s+`)

// TrimWhitespace removes leading and trailing whitespace from a string.
func TrimWhitespace(s string) string {
	return strings.TrimSpace(s)
}

// NormalizeSpaces collapses multiple consecutive whitespace characters into a single space.
// Also trims leading and trailing whitespace.
func NormalizeSpaces(s string) string {
	result := multiSpacePattern.ReplaceAllString(s, " ")
	return strings.TrimSpace(result)
}

// StripHTML removes all HTML tags from a string.
// Does not decode HTML entities.
func StripHTML(s string) string {
	return htmlTagPattern.ReplaceAllString(s, "")
}

// EscapeHTML escapes HTML special characters to their entity equivalents.
// Escapes: & < > " '.
func EscapeHTML(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		switch r {
		case '&':
			result.WriteString("&amp;")
		case '<':
			result.WriteString("&lt;")
		case '>':
			result.WriteString("&gt;")
		case '"':
			result.WriteString("&quot;")
		case '\'':
			result.WriteString("&#39;")
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

// NormalizeName normalizes a name by trimming whitespace,
// collapsing multiple spaces, and capitalizing the first letter of each word.
func NormalizeName(s string) string {
	s = NormalizeSpaces(s)
	if s == "" {
		return ""
	}

	words := strings.Fields(s)
	for i, word := range words {
		if word != "" {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToLower(runes[j])
			}
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// NormalizeEmail normalizes an email address by trimming whitespace
// and converting to lowercase.
func NormalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// RemoveNonPrintable removes non-printable characters from a string.
// Keeps printable ASCII and common Unicode characters.
func RemoveNonPrintable(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if unicode.IsPrint(r) || r == '\n' || r == '\r' || r == '\t' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// RemoveControlChars removes control characters except newline, carriage return, and tab.
func RemoveControlChars(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if !unicode.IsControl(r) || r == '\n' || r == '\r' || r == '\t' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToUppercase converts a string to uppercase.
func ToUppercase(s string) string {
	return strings.ToUpper(s)
}

// ToLowercase converts a string to lowercase.
func ToLowercase(s string) string {
	return strings.ToLower(s)
}

// RemoveDigits removes all digit characters from a string.
func RemoveDigits(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if !unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// KeepDigits keeps only digit characters in a string.
func KeepDigits(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// KeepAlphanumeric keeps only alphanumeric characters in a string.
func KeepAlphanumeric(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Func is a function type for sanitization operations.
type Func func(string) string

// Chain applies multiple sanitization functions in sequence.
// Functions are applied left to right.
func Chain(input string, fns ...Func) string {
	result := input
	for _, fn := range fns {
		result = fn(result)
	}
	return result
}

// Sanitizer provides a chainable API for building sanitization pipelines.
type Sanitizer struct {
	fns []Func
}

// NewSanitizer creates a new Sanitizer instance.
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		fns: make([]Func, 0),
	}
}

// TrimWhitespace adds whitespace trimming to the pipeline.
func (s *Sanitizer) TrimWhitespace() *Sanitizer {
	s.fns = append(s.fns, TrimWhitespace)
	return s
}

// NormalizeSpaces adds space normalization to the pipeline.
func (s *Sanitizer) NormalizeSpaces() *Sanitizer {
	s.fns = append(s.fns, NormalizeSpaces)
	return s
}

// StripHTML adds HTML stripping to the pipeline.
func (s *Sanitizer) StripHTML() *Sanitizer {
	s.fns = append(s.fns, StripHTML)
	return s
}

// EscapeHTML adds HTML escaping to the pipeline.
func (s *Sanitizer) EscapeHTML() *Sanitizer {
	s.fns = append(s.fns, EscapeHTML)
	return s
}

// NormalizeName adds name normalization to the pipeline.
func (s *Sanitizer) NormalizeName() *Sanitizer {
	s.fns = append(s.fns, NormalizeName)
	return s
}

// NormalizeEmail adds email normalization to the pipeline.
func (s *Sanitizer) NormalizeEmail() *Sanitizer {
	s.fns = append(s.fns, NormalizeEmail)
	return s
}

// ToUppercase adds uppercase conversion to the pipeline.
func (s *Sanitizer) ToUppercase() *Sanitizer {
	s.fns = append(s.fns, ToUppercase)
	return s
}

// ToLowercase adds lowercase conversion to the pipeline.
func (s *Sanitizer) ToLowercase() *Sanitizer {
	s.fns = append(s.fns, ToLowercase)
	return s
}

// RemoveNonPrintable adds non-printable character removal to the pipeline.
func (s *Sanitizer) RemoveNonPrintable() *Sanitizer {
	s.fns = append(s.fns, RemoveNonPrintable)
	return s
}

// RemoveControlChars adds control character removal to the pipeline.
func (s *Sanitizer) RemoveControlChars() *Sanitizer {
	s.fns = append(s.fns, RemoveControlChars)
	return s
}

// KeepDigits adds digit-only filtering to the pipeline.
func (s *Sanitizer) KeepDigits() *Sanitizer {
	s.fns = append(s.fns, KeepDigits)
	return s
}

// KeepAlphanumeric adds alphanumeric-only filtering to the pipeline.
func (s *Sanitizer) KeepAlphanumeric() *Sanitizer {
	s.fns = append(s.fns, KeepAlphanumeric)
	return s
}

// Custom adds a custom sanitization function to the pipeline.
func (s *Sanitizer) Custom(fn Func) *Sanitizer {
	s.fns = append(s.fns, fn)
	return s
}

// Apply applies all sanitization functions to the input.
func (s *Sanitizer) Apply(input string) string {
	return Chain(input, s.fns...)
}

// Common pre-built sanitizers

// TextSanitizer returns a sanitizer for general text input.
// Strips HTML, normalizes spaces, removes non-printable characters.
func TextSanitizer() *Sanitizer {
	return NewSanitizer().
		StripHTML().
		RemoveNonPrintable().
		NormalizeSpaces()
}

// NameSanitizer returns a sanitizer for name fields.
// Strips HTML, normalizes spaces, and capitalizes words.
func NameSanitizer() *Sanitizer {
	return NewSanitizer().
		StripHTML().
		RemoveNonPrintable().
		NormalizeName()
}

// EmailSanitizer returns a sanitizer for email addresses.
// Trims whitespace and converts to lowercase.
func EmailSanitizer() *Sanitizer {
	return NewSanitizer().
		TrimWhitespace().
		NormalizeEmail()
}

// PhoneSanitizer returns a sanitizer for phone numbers.
// Keeps only digits.
func PhoneSanitizer() *Sanitizer {
	return NewSanitizer().
		KeepDigits()
}
