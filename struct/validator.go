// Package structval provides struct validation using go-playground/validator with custom Txova tags.
package structval

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
	"github.com/Dorico-Dynamics/txova-go-validation/geo"
	"github.com/Dorico-Dynamics/txova-go-validation/phone"
	"github.com/Dorico-Dynamics/txova-go-validation/rating"
	"github.com/Dorico-Dynamics/txova-go-validation/ride"
	"github.com/Dorico-Dynamics/txova-go-validation/vehicle"
)

var (
	once     sync.Once
	validate *validator.Validate
)

// initValidator initializes the singleton validator with custom configuration.
func initValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Use JSON tag names for field names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	// Register custom validation tags.
	// These registrations cannot fail as they are valid tag names with valid functions.
	//nolint:errcheck // Registration errors are not possible with valid inputs
	validate.RegisterValidation("mz_phone", validateMzPhone)
	//nolint:errcheck // Registration errors are not possible with valid inputs
	validate.RegisterValidation("mz_plate", validateMzPlate)
	//nolint:errcheck // Registration errors are not possible with valid inputs
	validate.RegisterValidation("mz_location", validateMzLocation)
	//nolint:errcheck // Registration errors are not possible with valid inputs
	validate.RegisterValidation("txova_pin", validateTxovaPin)
	//nolint:errcheck // Registration errors are not possible with valid inputs
	validate.RegisterValidation("txova_money", validateTxovaMoney)
	//nolint:errcheck // Registration errors are not possible with valid inputs
	validate.RegisterValidation("txova_rating", validateTxovaRating)
}

// getValidator returns the singleton validator instance.
func getValidator() *validator.Validate {
	once.Do(initValidator)
	return validate
}

// Validate validates a struct and returns ValidationErrors.
// Returns nil if validation passes.
func Validate(s interface{}) valerrors.ValidationErrors {
	v := getValidator()

	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return translateErrors(validationErrors)
	}

	// Unexpected error type, wrap it.
	return valerrors.ValidationErrors{
		valerrors.New("_", valerrors.CodeInvalidFormat, err.Error()),
	}
}

// ValidateVar validates a single variable against a tag.
// Returns nil if validation passes.
func ValidateVar(field interface{}, tag string) valerrors.ValidationErrors {
	v := getValidator()

	err := v.Var(field, tag)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return translateErrors(validationErrors)
	}

	return valerrors.ValidationErrors{
		valerrors.New("value", valerrors.CodeInvalidFormat, err.Error()),
	}
}

// RegisterValidation registers a custom validation function.
// Returns an error if the tag is already registered or invalid.
func RegisterValidation(tag string, fn validator.Func) error {
	v := getValidator()
	return v.RegisterValidation(tag, fn)
}

// translateErrors converts go-playground validator errors to our ValidationErrors.
func translateErrors(errs validator.ValidationErrors) valerrors.ValidationErrors {
	if len(errs) == 0 {
		return nil
	}

	result := make(valerrors.ValidationErrors, 0, len(errs))
	for _, err := range errs {
		result = append(result, translateError(err))
	}
	return result
}

// translateError converts a single validator.FieldError to ValidationError.
func translateError(err validator.FieldError) valerrors.ValidationError {
	field := err.Field()
	tag := err.Tag()
	value := err.Value()

	// Handle special cases that need parameter access.
	if ve, ok := translateSpecialTag(err, field, tag, value); ok {
		return ve
	}

	// Handle simple format tags.
	if expected, ok := formatTagExpectations[tag]; ok {
		return valerrors.InvalidFormatWithValue(field, expected, value)
	}

	// Handle range tags.
	if isLowerBoundTag(tag) {
		return valerrors.OutOfRangeWithValue(field, err.Param(), "∞", value)
	}
	if isUpperBoundTag(tag) {
		return valerrors.OutOfRangeWithValue(field, "-∞", err.Param(), value)
	}

	// Default: use tag as expected format.
	return valerrors.InvalidFormatWithValue(field, tag, value)
}

// formatTagExpectations maps validation tags to expected format descriptions.
var formatTagExpectations = map[string]string{
	"email":     "valid email address",
	"url":       "valid URL",
	"mz_phone":  "valid Mozambique phone number",
	"mz_plate":  "valid Mozambique license plate",
	"txova_pin": "4-digit PIN (no sequential or repeated)",
}

// isLowerBoundTag returns true if the tag is a lower bound validation.
func isLowerBoundTag(tag string) bool {
	return tag == "gt" || tag == "gte"
}

// isUpperBoundTag returns true if the tag is an upper bound validation.
func isUpperBoundTag(tag string) bool {
	return tag == "lt" || tag == "lte"
}

// translateSpecialTag handles tags that need special processing.
func translateSpecialTag(err validator.FieldError, field, tag string, value interface{}) (valerrors.ValidationError, bool) {
	switch tag {
	case "required":
		return valerrors.Required(field), true

	case "min":
		return translateMinTag(err, field, value), true

	case "max":
		return translateMaxTag(err, field, value), true

	case "len":
		return valerrors.InvalidFormatWithValue(field, "length "+err.Param(), value), true

	case "oneof":
		options := strings.Split(err.Param(), " ")
		return valerrors.InvalidOptionWithValue(field, options, value), true

	case "mz_location":
		return valerrors.OutsideServiceArea(field), true

	case "txova_money":
		return valerrors.OutOfRangeWithValue(field, 1, "∞", value), true

	case "txova_rating":
		return valerrors.OutOfRangeWithValue(field, 1, 5, value), true

	default:
		return valerrors.ValidationError{}, false
	}
}

// translateMinTag handles the "min" validation tag.
func translateMinTag(err validator.FieldError, field string, value interface{}) valerrors.ValidationError {
	param := err.Param()
	if err.Kind() == reflect.String {
		if s, ok := value.(string); ok {
			return valerrors.TooShortWithValue(field, parseIntParam(param), len(s))
		}
	}
	return valerrors.OutOfRangeWithValue(field, param, "∞", value)
}

// translateMaxTag handles the "max" validation tag.
func translateMaxTag(err validator.FieldError, field string, value interface{}) valerrors.ValidationError {
	param := err.Param()
	if err.Kind() == reflect.String {
		if s, ok := value.(string); ok {
			return valerrors.TooLongWithValue(field, parseIntParam(param), len(s))
		}
	}
	return valerrors.OutOfRangeWithValue(field, "-∞", param, value)
}

// parseIntParam parses a string parameter to int, returning 0 on error.
func parseIntParam(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

// Custom validation functions

// validateMzPhone validates Mozambique phone numbers.
func validateMzPhone(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Empty is handled by required tag
	}
	return phone.Validate(value)
}

// validateMzPlate validates Mozambique license plates.
func validateMzPlate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Empty is handled by required tag
	}
	return vehicle.ValidatePlate(value) == nil
}

// validateMzLocation validates coordinates are within Mozambique.
// Expects a struct with Lat and Lon fields or a slice [lat, lon].
func validateMzLocation(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.Struct:
		return validateLocationStruct(field)
	case reflect.Slice, reflect.Array:
		return validateLocationSlice(field)
	default:
		return false
	}
}

// validateLocationStruct validates a struct with Lat/Latitude and Lon/Longitude fields.
func validateLocationStruct(field reflect.Value) bool {
	var lat, lon float64
	var foundLat, foundLon bool

	// Try common field names for latitude
	for _, name := range []string{"Lat", "Latitude", "lat", "latitude"} {
		if f := field.FieldByName(name); f.IsValid() && f.Kind() == reflect.Float64 {
			lat = f.Float()
			foundLat = true
			break
		}
	}

	// Try common field names for longitude
	for _, name := range []string{"Lon", "Lng", "Longitude", "lon", "lng", "longitude"} {
		if f := field.FieldByName(name); f.IsValid() && f.Kind() == reflect.Float64 {
			lon = f.Float()
			foundLon = true
			break
		}
	}

	if !foundLat || !foundLon {
		return false
	}

	return geo.ValidateInMozambique(lat, lon) == nil
}

// validateLocationSlice validates a [lat, lon] slice.
func validateLocationSlice(field reflect.Value) bool {
	if field.Len() < 2 {
		return false
	}

	lat := field.Index(0)
	lon := field.Index(1)

	if lat.Kind() != reflect.Float64 || lon.Kind() != reflect.Float64 {
		return false
	}

	return geo.ValidateInMozambique(lat.Float(), lon.Float()) == nil
}

// validateTxovaPin validates ride verification PINs.
func validateTxovaPin(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Empty is handled by required tag
	}
	return ride.ValidatePIN(value) == nil
}

// validateTxovaMoney validates positive money amounts.
// Expects an int64 value representing centavos.
func validateTxovaMoney(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() > 0
	case reflect.Float32, reflect.Float64:
		return field.Float() > 0
	default:
		return false
	}
}

// validateTxovaRating validates rating values (1-5).
func validateTxovaRating(fl validator.FieldLevel) bool {
	field := fl.Field()

	var value int
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := field.Int()
		if v < 0 || v > 5 {
			return false
		}
		value = int(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := field.Uint()
		if v > 5 {
			return false
		}
		value = int(v) // #nosec G115 - bounds checked above, max value is 5
	default:
		return false
	}

	return rating.ValidateRating(value) == nil
}
