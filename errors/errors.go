// Package errors provides structured validation error types for the Txova validation library.
package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Error codes for validation failures.
const (
	// CodeRequired indicates a required field is missing.
	CodeRequired = "REQUIRED"
	// CodeInvalidFormat indicates the value doesn't match expected format.
	CodeInvalidFormat = "INVALID_FORMAT"
	// CodeOutOfRange indicates the value is outside allowed range.
	CodeOutOfRange = "OUT_OF_RANGE"
	// CodeTooShort indicates the value is below minimum length.
	CodeTooShort = "TOO_SHORT"
	// CodeTooLong indicates the value exceeds maximum length.
	CodeTooLong = "TOO_LONG"
	// CodeInvalidOption indicates the value is not in allowed options.
	CodeInvalidOption = "INVALID_OPTION"
	// CodeOutsideServiceArea indicates location is not in a serviceable area.
	CodeOutsideServiceArea = "OUTSIDE_SERVICE_AREA"
)

// ValidationError represents a single validation failure.
type ValidationError struct {
	// Field is the JSON field name that failed validation.
	Field string `json:"field"`
	// Code is the validation error code.
	Code string `json:"code"`
	// Message is a human-readable error message.
	Message string `json:"message"`
	// Value is the invalid value (masked if sensitive).
	Value interface{} `json:"value,omitempty"`
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("%s: %s (value: %v)", e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// New creates a new ValidationError.
func New(field, code, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    code,
		Message: message,
	}
}

// NewWithValue creates a new ValidationError with the invalid value included.
func NewWithValue(field, code, message string, value interface{}) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    code,
		Message: message,
		Value:   value,
	}
}

// Required creates a REQUIRED validation error.
func Required(field string) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeRequired,
		Message: fmt.Sprintf("%s is required", field),
	}
}

// InvalidFormat creates an INVALID_FORMAT validation error.
func InvalidFormat(field, expected string) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeInvalidFormat,
		Message: fmt.Sprintf("%s has invalid format, expected %s", field, expected),
	}
}

// InvalidFormatWithValue creates an INVALID_FORMAT validation error with the invalid value.
func InvalidFormatWithValue(field, expected string, value interface{}) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeInvalidFormat,
		Message: fmt.Sprintf("%s has invalid format, expected %s", field, expected),
		Value:   value,
	}
}

// OutOfRange creates an OUT_OF_RANGE validation error.
func OutOfRange(field string, min, max interface{}) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeOutOfRange,
		Message: fmt.Sprintf("%s must be between %v and %v", field, min, max),
	}
}

// OutOfRangeWithValue creates an OUT_OF_RANGE validation error with the invalid value.
func OutOfRangeWithValue(field string, min, max, value interface{}) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeOutOfRange,
		Message: fmt.Sprintf("%s must be between %v and %v", field, min, max),
		Value:   value,
	}
}

// TooShort creates a TOO_SHORT validation error.
func TooShort(field string, minLength int) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeTooShort,
		Message: fmt.Sprintf("%s must be at least %d characters", field, minLength),
	}
}

// TooShortWithValue creates a TOO_SHORT validation error with the actual length.
func TooShortWithValue(field string, minLength, actualLength int) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeTooShort,
		Message: fmt.Sprintf("%s must be at least %d characters", field, minLength),
		Value:   actualLength,
	}
}

// TooLong creates a TOO_LONG validation error.
func TooLong(field string, maxLength int) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeTooLong,
		Message: fmt.Sprintf("%s must be at most %d characters", field, maxLength),
	}
}

// TooLongWithValue creates a TOO_LONG validation error with the actual length.
func TooLongWithValue(field string, maxLength, actualLength int) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeTooLong,
		Message: fmt.Sprintf("%s must be at most %d characters", field, maxLength),
		Value:   actualLength,
	}
}

// InvalidOption creates an INVALID_OPTION validation error.
func InvalidOption(field string, allowedOptions []string) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeInvalidOption,
		Message: fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowedOptions, ", ")),
	}
}

// InvalidOptionWithValue creates an INVALID_OPTION validation error with the invalid value.
func InvalidOptionWithValue(field string, allowedOptions []string, value interface{}) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeInvalidOption,
		Message: fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowedOptions, ", ")),
		Value:   value,
	}
}

// OutsideServiceArea creates an OUTSIDE_SERVICE_AREA validation error.
func OutsideServiceArea(field string) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeOutsideServiceArea,
		Message: fmt.Sprintf("%s is outside the service area", field),
	}
}

// OutsideServiceAreaWithValue creates an OUTSIDE_SERVICE_AREA error with coordinates.
func OutsideServiceAreaWithValue(field string, lat, lon float64) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    CodeOutsideServiceArea,
		Message: fmt.Sprintf("%s is outside the service area", field),
		Value:   fmt.Sprintf("%.6f, %.6f", lat, lon),
	}
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// Error implements the error interface.
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	if len(ve) == 1 {
		return ve[0].Error()
	}

	var msgs []string
	for _, e := range ve {
		msgs = append(msgs, e.Error())
	}
	return fmt.Sprintf("%d validation errors: %s", len(ve), strings.Join(msgs, "; "))
}

// HasErrors returns true if there are any validation errors.
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// HasField returns true if there is a validation error for the given field.
func (ve ValidationErrors) HasField(field string) bool {
	for _, e := range ve {
		if e.Field == field {
			return true
		}
	}
	return false
}

// GetByField returns all validation errors for the given field.
func (ve ValidationErrors) GetByField(field string) ValidationErrors {
	var result ValidationErrors
	for _, e := range ve {
		if e.Field == field {
			result = append(result, e)
		}
	}
	return result
}

// GetByCode returns all validation errors with the given code.
func (ve ValidationErrors) GetByCode(code string) ValidationErrors {
	var result ValidationErrors
	for _, e := range ve {
		if e.Code == code {
			result = append(result, e)
		}
	}
	return result
}

// First returns the first validation error, or nil if empty.
func (ve ValidationErrors) First() *ValidationError {
	if len(ve) == 0 {
		return nil
	}
	return &ve[0]
}

// Fields returns a list of unique field names that have errors.
func (ve ValidationErrors) Fields() []string {
	seen := make(map[string]bool)
	var fields []string
	for _, e := range ve {
		if !seen[e.Field] {
			seen[e.Field] = true
			fields = append(fields, e.Field)
		}
	}
	return fields
}

// Add appends a validation error to the collection.
func (ve *ValidationErrors) Add(err ValidationError) {
	*ve = append(*ve, err)
}

// AddAll appends multiple validation errors to the collection.
func (ve *ValidationErrors) AddAll(errs ValidationErrors) {
	*ve = append(*ve, errs...)
}

// MarshalJSON implements json.Marshaler for API responses.
func (ve ValidationErrors) MarshalJSON() ([]byte, error) {
	if len(ve) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal([]ValidationError(ve))
}

// ToError returns the ValidationErrors as an error interface, or nil if empty.
func (ve ValidationErrors) ToError() error {
	if len(ve) == 0 {
		return nil
	}
	return ve
}
