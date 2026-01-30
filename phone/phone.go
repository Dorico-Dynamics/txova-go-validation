// Package phone provides Mozambique phone number validation and normalization.
package phone

import (
	"regexp"
	"strings"

	"github.com/Dorico-Dynamics/txova-go-types/contact"
)

// MozambiqueCountryCode is the country calling code for Mozambique.
const MozambiqueCountryCode = "258"

// validPrefixes are the valid Mozambique mobile prefixes.
var validPrefixes = map[string]bool{
	"82": true,
	"83": true,
	"84": true,
	"85": true,
	"86": true,
	"87": true,
}

// digitsOnly matches all non-digit characters.
var digitsOnly = regexp.MustCompile(`\D`)

// Validate checks if the input is a valid Mozambique phone number.
// Returns true if the number can be parsed and normalized to a valid format.
func Validate(input string) bool {
	_, err := Normalize(input)
	return err == nil
}

// Normalize converts various phone number formats to the standard +258XXXXXXXXX format.
// Accepts formats:
//   - Local: 841234567
//   - International: +258841234567
//   - With country code: 258841234567
//   - With 00 prefix: 00258841234567
//   - With spaces/dashes: 84 123 4567, 84-123-4567
//
// Returns the normalized phone number string or an error if invalid.
func Normalize(input string) (string, error) {
	if input == "" {
		return "", contact.ErrInvalidPhoneNumber
	}

	// Remove all non-digit characters except leading +
	hasPlus := strings.HasPrefix(strings.TrimSpace(input), "+")
	digits := digitsOnly.ReplaceAllString(input, "")

	if digits == "" {
		return "", contact.ErrInvalidPhoneNumber
	}

	// Normalize to 9 digits (local number without country code)
	var localNumber string

	switch {
	case len(digits) == 9:
		// Local format: 841234567
		localNumber = digits
	case len(digits) == 12 && strings.HasPrefix(digits, MozambiqueCountryCode):
		// International format without +: 258841234567
		localNumber = digits[3:]
	case len(digits) == 12 && hasPlus:
		// International format with +: +258841234567 (+ already stripped)
		if strings.HasPrefix(digits, MozambiqueCountryCode) {
			localNumber = digits[3:]
		} else {
			return "", contact.ErrInvalidPhoneNumber
		}
	case len(digits) == 14 && strings.HasPrefix(digits, "00"+MozambiqueCountryCode):
		// With 00 prefix: 00258841234567
		localNumber = digits[5:]
	default:
		return "", contact.ErrInvalidPhoneNumber
	}

	// Validate prefix
	if len(localNumber) != 9 {
		return "", contact.ErrInvalidPhoneNumber
	}

	prefix := localNumber[:2]
	if !validPrefixes[prefix] {
		return "", contact.ErrInvalidMobilePrefix
	}

	return "+" + MozambiqueCountryCode + localNumber, nil
}

// IdentifyOperator returns the mobile network operator name for the given phone number.
// Returns an empty string if the number is invalid or operator cannot be determined.
func IdentifyOperator(input string) string {
	normalized, err := Normalize(input)
	if err != nil {
		return ""
	}

	// Parse using types library to get operator
	phone, err := contact.ParsePhoneNumber(normalized)
	if err != nil {
		return ""
	}

	return phone.Operator().String()
}

// GetPrefix extracts the mobile prefix from a phone number.
// Returns the 2-digit prefix (82-87) or empty string if invalid.
func GetPrefix(input string) string {
	normalized, err := Normalize(input)
	if err != nil {
		return ""
	}

	// Normalized format is +258XXXXXXXXX, prefix is at position 4-5
	if len(normalized) >= 6 {
		return normalized[4:6]
	}
	return ""
}

// IsVodacom returns true if the phone number belongs to Vodacom (prefixes 82, 84, 85).
func IsVodacom(input string) bool {
	prefix := GetPrefix(input)
	return prefix == "82" || prefix == "84" || prefix == "85"
}

// IsMovitel returns true if the phone number belongs to Movitel (prefixes 83, 86).
func IsMovitel(input string) bool {
	prefix := GetPrefix(input)
	return prefix == "83" || prefix == "86"
}

// IsTmcel returns true if the phone number belongs to Tmcel (prefix 87).
func IsTmcel(input string) bool {
	return GetPrefix(input) == "87"
}
