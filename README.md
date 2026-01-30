# txova-go-validation

Input validation library providing Mozambique-specific validators, struct validation, and sanitization utilities for Txova services.

## Overview

`txova-go-validation` provides comprehensive input validation tailored for Mozambique, including phone number validation with operator detection, geographic bounds checking, vehicle plate validation, and integration with go-playground/validator.

**Module:** `github.com/Dorico-Dynamics/txova-go-validation`

## Installation

```bash
go get github.com/Dorico-Dynamics/txova-go-validation
```

## Packages

| Package | Import | Description |
|---------|--------|-------------|
| `errors` | `valerrors` | Structured validation error types |
| `phone` | `phone` | Mozambique phone number validation |
| `geo` | `geo` | Geographic coordinate validation |
| `vehicle` | `vehicle` | License plate and vehicle year validation |
| `ride` | `ride` | PIN, distance, and fare validation |
| `rating` | `rating` | Rating and review validation with profanity detection |
| `document` | `document` | Document upload validation |
| `struct` | `structval` | Struct validation with custom tags |
| `sanitize` | `sanitize` | Input sanitization utilities |

## Usage

### Errors Package

Structured validation errors for consistent API responses.

```go
import valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"

// Create individual errors
err := valerrors.Required("email")
err := valerrors.InvalidFormat("phone", "Mozambique phone format (+258XXXXXXXXX)")
err := valerrors.OutOfRange("rating", 1, 5)
err := valerrors.TooLong("review", 500)
err := valerrors.InvalidOption("status", []string{"pending", "active", "completed"})
err := valerrors.OutsideServiceArea("pickup")

// Collect multiple errors
var errs valerrors.ValidationErrors
errs.Add(valerrors.Required("phone"))
errs.Add(valerrors.InvalidFormat("email", "valid email address"))

if errs.HasErrors() {
    // Check specific field
    if errs.HasField("phone") {
        phoneErrs := errs.GetByField("phone")
    }
    
    // Get errors by code
    formatErrs := errs.GetByCode(valerrors.CodeInvalidFormat)
    
    // Get all field names with errors
    fields := errs.Fields() // []string{"phone", "email"}
    
    // JSON serialization for API responses
    jsonBytes, _ := json.Marshal(errs)
}
```

**Error Codes:**

| Code | Description |
|------|-------------|
| `REQUIRED` | Field is required |
| `INVALID_FORMAT` | Format doesn't match expected pattern |
| `OUT_OF_RANGE` | Value outside allowed range |
| `TOO_SHORT` | Below minimum length |
| `TOO_LONG` | Exceeds maximum length |
| `INVALID_OPTION` | Not in allowed options |
| `OUTSIDE_SERVICE_AREA` | Location not serviceable |

### Phone Package

Mozambique phone number validation and normalization.

```go
import "github.com/Dorico-Dynamics/txova-go-validation/phone"

// Validate phone number
valid := phone.Validate("841234567")     // true
valid := phone.Validate("+258841234567") // true
valid := phone.Validate("123456789")     // false (invalid prefix)

// Normalize to international format
normalized, err := phone.Normalize("84 123 4567")
// normalized = "+258841234567"

normalized, err := phone.Normalize("00258841234567")
// normalized = "+258841234567"

// Identify mobile operator
operator := phone.IdentifyOperator("+258841234567") // "Vodacom"
operator := phone.IdentifyOperator("+258831234567") // "Movitel"
operator := phone.IdentifyOperator("+258871234567") // "Tmcel"

// Check specific operator
phone.IsVodacom("+258841234567") // true (prefixes: 82, 84, 85)
phone.IsMovitel("+258831234567") // true (prefixes: 83, 86)
phone.IsTmcel("+258871234567")   // true (prefix: 87)

// Get prefix
prefix := phone.GetPrefix("+258841234567") // "84"
```

**Supported Input Formats:**
- Local: `841234567`
- International: `+258841234567`
- With country code: `258841234567`
- With leading zeros: `00258841234567`
- With separators: `84 123 4567`, `84-123-4567`, `84.123.4567`

**Valid Prefixes:** 82, 83, 84, 85, 86, 87

### Geo Package

Geographic validation for Mozambique locations and service areas.

```go
import "github.com/Dorico-Dynamics/txova-go-validation/geo"

// Validate coordinates are within valid global bounds
err := geo.ValidateCoordinates(-25.969, 32.573)

// Validate coordinates are within Mozambique
err := geo.ValidateInMozambique(-25.969, 32.573)

// Quick checks
geo.IsInMozambique(-25.969, 32.573) // true
geo.IsInServiceArea(-25.969, 32.573) // true (if in any service area)

// Validate within specific service area
err := geo.ValidateServiceArea(-25.969, 32.573, "maputo")

// Validate within any active service area
err := geo.ValidateAnyServiceArea(-25.969, 32.573)

// Find which service area contains coordinates
area := geo.FindServiceArea(-25.969, 32.573) // "maputo" or ""

// Get available service areas
areas := geo.GetServiceAreas() // ["maputo", "matola", "beira"]

// Get service area configuration
config := geo.GetServiceArea("maputo")
// config.MinLat, config.MaxLat, config.MinLon, config.MaxLon

// Calculate distance between two points (Haversine formula)
distanceKM, err := geo.CalculateDistance(lat1, lon1, lat2, lon2)
```

**Mozambique Bounds:**
- Latitude: -26.9 to -10.3
- Longitude: 30.2 to 41.0

**Service Areas:**
- `maputo`: Maputo City
- `matola`: Matola
- `beira`: Beira

### Vehicle Package

Mozambique vehicle validation including license plates and years.

```go
import "github.com/Dorico-Dynamics/txova-go-validation/vehicle"

// Validate license plate
err := vehicle.ValidatePlate("AAA-123-MC")
err := vehicle.ValidatePlate("MC-12-34") // old format

// Normalize plate to standard format
normalized, err := vehicle.NormalizePlate("aaa 123 mp")
// normalized = "AAA-123-MP"

// Quick validation check
vehicle.IsValidPlate("AAA-123-MC") // true

// Check plate format
vehicle.IsStandardFormat("AAA-123-MC") // true (AAA-NNN-LL)
vehicle.IsOldFormat("MC-12-34")        // true (LL-NN-NN)

// Get province information
code := vehicle.GetProvince("AAA-123-MP")     // "MP"
name := vehicle.GetProvinceName("AAA-123-MP") // "Maputo Province"

// Validate vehicle year (2010 to current year + 1)
err := vehicle.ValidateYear(2020)
vehicle.IsValidYear(2020) // true
```

**License Plate Formats:**
- Standard: `AAA-NNN-LL` (e.g., AAA-123-MC)
- Old: `LL-NN-NN` (e.g., MC-12-34)

**Province Codes:**
MC (Maputo City), MP (Maputo Province), GZ (Gaza), IB (Inhambane), SF (Sofala), MN (Manica), TT (Tete), ZB (Zambezia), NP (Nampula), CA (Cabo Delgado), NS (Niassa)

### Ride Package

Ride-specific validation for PIN, distance, fare, and locations.

```go
import "github.com/Dorico-Dynamics/txova-go-validation/ride"

// Validate 4-digit PIN (no sequential or repeated patterns)
err := ride.ValidatePIN("7392")
ride.IsValidPIN("7392") // true
ride.IsValidPIN("1234") // false (sequential)
ride.IsValidPIN("1111") // false (repeated)

// Validate distance (0.5 to 200 km)
err := ride.ValidateDistance(15.5)
ride.IsValidDistance(15.5) // true

// Validate fare (5000 to 5000000 centavos = 50 to 50,000 MZN)
err := ride.ValidateFare(10000) // 100 MZN in centavos
ride.IsValidFare(10000) // true

// Validate fare using Money type from txova-go-types
err := ride.ValidateFareMoney(moneyAmount)

// Validate pickup and dropoff separation (minimum 0.1 km)
err := ride.ValidatePickupDropoff(pickupLat, pickupLon, dropoffLat, dropoffLon)

// Using Location types
err := ride.ValidatePickupDropoffLocations(pickupLocation, dropoffLocation)

// Calculate estimated fare
fare := ride.CalculateEstimatedFare(distanceKM, baseFareCentavos, perKMCentavos)
```

**PIN Rules:**
- Exactly 4 digits
- No sequential patterns (1234, 4321, 5678, 8765, etc.)
- No repeated digits (1111, 2222, 0000, 9999, etc.)

**Fare Limits:**
- Minimum: 5,000 centavos (50 MZN)
- Maximum: 5,000,000 centavos (50,000 MZN)

### Rating Package

Rating and review validation with profanity detection.

```go
import "github.com/Dorico-Dynamics/txova-go-validation/rating"

// Validate rating (1 to 5)
err := rating.ValidateRating(5)
rating.IsValidRating(5) // true

// Validate review text (max 500 characters)
err := rating.ValidateReviewText("Great driver!")
rating.IsValidReviewText("Great driver!") // true

// Sanitize review text (strip HTML, normalize whitespace)
cleaned := rating.SanitizeReviewText("  <b>Great</b>  driver!  ")
// cleaned = "Great driver!"

// Check for profanity
hasProfanity := rating.CheckProfanity("bad word here")

// Validate and sanitize in one step
cleaned, err := rating.ValidateAndSanitizeReview("  <script>x</script>Great!  ")
// cleaned = "Great!"

// Full processing with profanity check
result, err := rating.ProcessReview("  <b>Great driver!</b>  ")
// result.Text = "Great driver!"
// result.HasProfanity = false
// result.RequiresReview = false
// result.OriginalLength = 24
// result.SanitizedLength = 13
```

**Profanity Detection:**
- Detects common profanity in English and Portuguese
- Conservative detection for moderation flagging
- Case-insensitive matching

### Document Package

Document and file upload validation.

```go
import "github.com/Dorico-Dynamics/txova-go-validation/document"

// Validate document type
err := document.ValidateDocType("driver_license")
document.IsValidDocType("driver_license") // true

// Get all valid document types
types := document.AllDocTypes()
// ["driver_license", "vehicle_registration", "insurance", "id_card", "profile_photo", "vehicle_photo"]

// Validate file size (limits vary by document type)
err := document.ValidateFileSize(fileSize, "profile_photo") // max 2MB
err := document.ValidateFileSize(fileSize, "driver_license") // max 5MB

// Validate MIME type and extension
err := document.ValidateMIMEType("image/jpeg", "jpg")

// Validate image dimensions (200-4096 pixels)
err := document.ValidateImageDimensions(1920, 1080)

// Validate aspect ratio (1:4 to 4:1)
err := document.ValidateAspectRatio(1920, 1080)

// Validate all image properties at once
err := document.ValidateImage(width, height, fileSize, "vehicle_photo")

// Get allowed formats for document type
formats := document.GetAllowedFormats("driver_license") // ["jpg", "jpeg", "png", "pdf"]
formats := document.GetAllowedFormats("profile_photo")  // ["jpg", "jpeg", "png"]

// Check if format is allowed
document.IsAllowedFormat("pdf", "driver_license") // true
document.IsAllowedFormat("pdf", "profile_photo")  // false

// Check if document type is an image type
document.IsImageType("profile_photo") // true
document.IsImageType("driver_license") // false

// Get max file size for document type
maxSize := document.GetMaxFileSize("profile_photo") // 2097152 (2MB)
```

**Document Types:**
- `driver_license`: Driver's license (jpg, jpeg, png, pdf) - 5MB max
- `vehicle_registration`: Vehicle registration (jpg, jpeg, png, pdf) - 5MB max
- `insurance`: Insurance document (jpg, jpeg, png, pdf) - 5MB max
- `id_card`: ID card (jpg, jpeg, png, pdf) - 5MB max
- `profile_photo`: Profile photo (jpg, jpeg, png) - 2MB max
- `vehicle_photo`: Vehicle photo (jpg, jpeg, png) - 5MB max

**Image Constraints:**
- Dimensions: 200 to 4096 pixels (width and height)
- Aspect ratio: 1:4 to 4:1

### Sanitize Package

Input sanitization utilities with chainable API.

```go
import "github.com/Dorico-Dynamics/txova-go-validation/sanitize"

// Standalone functions
text := sanitize.TrimWhitespace("  hello  ")           // "hello"
text := sanitize.NormalizeSpaces("hello   world")      // "hello world"
text := sanitize.StripHTML("<b>hello</b>")             // "hello"
text := sanitize.EscapeHTML("<script>")                // "&lt;script&gt;"
text := sanitize.NormalizeName("  joão   silva  ")     // "João Silva"
text := sanitize.NormalizeEmail("  User@EXAMPLE.COM ") // "user@example.com"
text := sanitize.RemoveNonPrintable("hello\x00world")  // "helloworld"
text := sanitize.RemoveControlChars("hello\x00world")  // "helloworld"
text := sanitize.ToUppercase("hello")                  // "HELLO"
text := sanitize.ToLowercase("HELLO")                  // "hello"
text := sanitize.RemoveDigits("abc123")                // "abc"
text := sanitize.KeepDigits("abc123")                  // "123"
text := sanitize.KeepAlphanumeric("abc-123!")          // "abc123"

// Chain multiple functions
result := sanitize.Chain(input, 
    sanitize.StripHTML, 
    sanitize.NormalizeSpaces, 
    sanitize.TrimWhitespace,
)

// Chainable builder pattern
result := sanitize.NewSanitizer().
    StripHTML().
    RemoveNonPrintable().
    NormalizeSpaces().
    TrimWhitespace().
    Apply(input)

// Add custom sanitization function
result := sanitize.NewSanitizer().
    StripHTML().
    Custom(func(s string) string {
        return strings.ReplaceAll(s, "bad", "good")
    }).
    Apply(input)

// Pre-built sanitizers for common use cases
text := sanitize.TextSanitizer().Apply(input)  // StripHTML -> RemoveNonPrintable -> NormalizeSpaces
name := sanitize.NameSanitizer().Apply(input)  // StripHTML -> RemoveNonPrintable -> NormalizeName
email := sanitize.EmailSanitizer().Apply(input) // TrimWhitespace -> NormalizeEmail
phone := sanitize.PhoneSanitizer().Apply(input) // KeepDigits
```

### Struct Package (structval)

Struct validation using go-playground/validator with Txova-specific custom tags.

```go
import structval "github.com/Dorico-Dynamics/txova-go-validation/struct"

// Define struct with validation tags
type CreateRideRequest struct {
    Phone       string    `json:"phone" validate:"required,mz_phone"`
    DriverPhone string    `json:"driver_phone" validate:"omitempty,mz_phone"`
    Pickup      Location  `json:"pickup" validate:"required,mz_location"`
    Dropoff     Location  `json:"dropoff" validate:"required,mz_location"`
    PIN         string    `json:"pin" validate:"required,txova_pin"`
    Fare        int64     `json:"fare" validate:"required,txova_money"`
    Rating      int       `json:"rating" validate:"omitempty,txova_rating"`
    VehicleYear int       `json:"vehicle_year" validate:"required,txova_vehicle_year"`
    Plate       string    `json:"plate" validate:"required,mz_plate"`
}

type Location struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

// Validate struct
request := CreateRideRequest{
    Phone:       "+258841234567",
    Pickup:      Location{Lat: -25.969, Lon: 32.573},
    Dropoff:     Location{Lat: -25.970, Lon: 32.580},
    PIN:         "7392",
    Fare:        10000,
    VehicleYear: 2020,
    Plate:       "AAA-123-MC",
}

errs := structval.Validate(request)
if errs != nil {
    // Handle validation errors
    for _, e := range errs {
        fmt.Printf("Field: %s, Code: %s, Message: %s\n", e.Field, e.Code, e.Message)
    }
}

// Validate single field
errs := structval.ValidateVar("+258841234567", "required,mz_phone")

// Register custom validator
structval.RegisterValidation("my_custom", func(fl validator.FieldLevel) bool {
    return fl.Field().String() != "invalid"
})
```

**Custom Validation Tags:**

| Tag | Description | Valid Examples |
|-----|-------------|----------------|
| `mz_phone` | Mozambique phone number | `+258841234567`, `841234567` |
| `mz_plate` | Mozambique license plate | `AAA-123-MC`, `MC-12-34` |
| `mz_location` | Coordinates within Mozambique | struct with Lat/Lon fields, `[-25.969, 32.573]` |
| `txova_pin` | 4-digit PIN (no sequential/repeated) | `7392`, `4826` |
| `txova_money` | Positive money amount | any positive int64, int, uint, or float |
| `txova_rating` | Rating 1-5 | `1`, `2`, `3`, `4`, `5` |
| `txova_vehicle_year` | Year 2010 to current+1 | `2015`, `2020`, `2025` |

**Standard go-playground/validator Tags:**

All standard tags are supported, including:
- `required` - Field is required
- `omitempty` - Field is optional
- `email` - Valid email format
- `url` - Valid URL format
- `min=N` - Minimum value/length
- `max=N` - Maximum value/length
- `len=N` - Exact length
- `oneof=a b c` - One of specified values

**Location Validation:**

The `mz_location` tag supports multiple formats:
- Struct with `Lat`/`Latitude` and `Lon`/`Longitude` fields
- Slice/array with `[lat, lon]` values

## Dependencies

**Internal:**
- `github.com/Dorico-Dynamics/txova-go-types` - Domain types (Phone, Location, Money, etc.)

**External:**
- `github.com/go-playground/validator/v10` - Struct validation

## Development

### Requirements

- Go 1.25+

### Testing

```bash
go test ./...
```

### Test Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

**Coverage Target:** >90%

### Linting

```bash
golangci-lint run
```

## License

Proprietary - Dorico Dynamics
