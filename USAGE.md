# txova-go-validation Usage Guide

Complete usage guide for the Txova validation library, providing Mozambique-specific validators, struct validation, and sanitization utilities.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Packages](#packages)
  - [errors - Structured Validation Errors](#errors-package)
  - [phone - Phone Number Validation](#phone-package)
  - [geo - Geographic Validation](#geo-package)
  - [vehicle - Vehicle Validation](#vehicle-package)
  - [ride - Ride Validation](#ride-package)
  - [rating - Rating & Review Validation](#rating-package)
  - [document - Document Upload Validation](#document-package)
  - [sanitize - Input Sanitization](#sanitize-package)
  - [struct - Struct Validation](#struct-package)
- [Integration Patterns](#integration-patterns)
- [Error Handling](#error-handling)
- [API Response Examples](#api-response-examples)

---

## Installation

```bash
go get github.com/Dorico-Dynamics/txova-go-validation
```

**Requirements:**
- Go 1.25+
- Dependencies are managed automatically via go modules

---

## Quick Start

```go
package main

import (
    "encoding/json"
    "fmt"

    valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
    "github.com/Dorico-Dynamics/txova-go-validation/phone"
    "github.com/Dorico-Dynamics/txova-go-validation/sanitize"
    structval "github.com/Dorico-Dynamics/txova-go-validation/struct"
)

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Phone string `json:"phone" validate:"required,mz_phone"`
    Email string `json:"email" validate:"required,email"`
}

func main() {
    // Sanitize inputs first
    name := sanitize.NameSanitizer().Apply("  joão   silva  ")
    email := sanitize.EmailSanitizer().Apply("  USER@EXAMPLE.COM  ")
    phoneInput := "+258 84 123 4567"

    // Normalize phone
    normalizedPhone, err := phone.Normalize(phoneInput)
    if err != nil {
        fmt.Println("Invalid phone:", err)
        return
    }

    // Create request
    req := CreateUserRequest{
        Name:  name,   // "João Silva"
        Phone: normalizedPhone, // "+258841234567"
        Email: email,  // "user@example.com"
    }

    // Validate struct
    errs := structval.Validate(req)
    if errs != nil {
        // Convert to JSON for API response
        jsonBytes, _ := json.MarshalIndent(errs, "", "  ")
        fmt.Println(string(jsonBytes))
        return
    }

    fmt.Println("Validation passed!")
}
```

---

## Packages

### errors Package

Structured validation errors for consistent API responses.

**Import:**
```go
import valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
```

#### Error Codes

| Code | Constant | Description |
|------|----------|-------------|
| `REQUIRED` | `CodeRequired` | Field is required but missing |
| `INVALID_FORMAT` | `CodeInvalidFormat` | Value doesn't match expected format |
| `OUT_OF_RANGE` | `CodeOutOfRange` | Numeric value outside allowed range |
| `TOO_SHORT` | `CodeTooShort` | String below minimum length |
| `TOO_LONG` | `CodeTooLong` | String exceeds maximum length |
| `INVALID_OPTION` | `CodeInvalidOption` | Value not in allowed options |
| `OUTSIDE_SERVICE_AREA` | `CodeOutsideServiceArea` | Location not serviceable |

#### Creating Errors

```go
// Basic error constructors
err := valerrors.Required("email")
err := valerrors.InvalidFormat("phone", "Mozambique phone format (+258XXXXXXXXX)")
err := valerrors.OutOfRange("rating", 1, 5)
err := valerrors.TooShort("password", 8)
err := valerrors.TooLong("bio", 500)
err := valerrors.InvalidOption("status", []string{"pending", "active", "completed"})
err := valerrors.OutsideServiceArea("pickup")

// With value included (useful for debugging)
err := valerrors.InvalidFormatWithValue("phone", "Mozambique format", "123456")
err := valerrors.OutOfRangeWithValue("rating", 1, 5, 10)
err := valerrors.TooShortWithValue("password", 8, 4)
```

#### Collecting Multiple Errors

```go
var errs valerrors.ValidationErrors

// Add errors
errs.Add(valerrors.Required("email"))
errs.Add(valerrors.InvalidFormat("phone", "Mozambique format"))
errs.Add(valerrors.TooShort("password", 8))

// Check if any errors exist
if errs.HasErrors() {
    // Check specific field
    if errs.HasField("email") {
        emailErrs := errs.GetByField("email")
        // Handle email-specific errors
    }

    // Get errors by code
    formatErrs := errs.GetByCode(valerrors.CodeInvalidFormat)

    // Get first error
    first := errs.First()

    // Get all unique field names
    fields := errs.Fields() // []string{"email", "phone", "password"}
}
```

#### JSON Serialization

```go
errs := valerrors.ValidationErrors{
    valerrors.Required("email"),
    valerrors.InvalidFormat("phone", "Mozambique format"),
}

jsonBytes, _ := json.Marshal(errs)
// Output:
// [
//   {"field":"email","code":"REQUIRED","message":"email is required"},
//   {"field":"phone","code":"INVALID_FORMAT","message":"phone must be in Mozambique format format"}
// ]
```

---

### phone Package

Mozambique phone number validation, normalization, and operator detection.

**Import:**
```go
import "github.com/Dorico-Dynamics/txova-go-validation/phone"
```

#### Constants

```go
const MozambiqueCountryCode = "258"
// Valid prefixes: 82, 83, 84, 85, 86, 87
```

#### Validation

```go
// Quick validation check
valid := phone.Validate("+258841234567") // true
valid := phone.Validate("841234567")      // true
valid := phone.Validate("123456789")      // false (invalid prefix)
valid := phone.Validate("84123")          // false (too short)
```

#### Normalization

Converts any valid format to standardized `+258XXXXXXXXX` format.

```go
// All these normalize to "+258841234567"
normalized, err := phone.Normalize("841234567")
normalized, err := phone.Normalize("+258841234567")
normalized, err := phone.Normalize("258841234567")
normalized, err := phone.Normalize("00258841234567")
normalized, err := phone.Normalize("84 123 4567")
normalized, err := phone.Normalize("84-123-4567")
normalized, err := phone.Normalize("84.123.4567")

if err != nil {
    // Invalid phone number
}
```

#### Operator Detection

| Operator | Prefixes |
|----------|----------|
| Vodacom | 82, 84, 85 |
| Movitel | 83, 86 |
| Tmcel | 87 |

```go
// Identify operator
operator := phone.IdentifyOperator("+258841234567") // "Vodacom"
operator := phone.IdentifyOperator("+258831234567") // "Movitel"
operator := phone.IdentifyOperator("+258871234567") // "Tmcel"

// Check specific operator
phone.IsVodacom("+258841234567") // true
phone.IsMovitel("+258831234567") // true
phone.IsTmcel("+258871234567")   // true

// Get prefix
prefix := phone.GetPrefix("+258841234567") // "84"
```

---

### geo Package

Geographic validation for Mozambique locations and service areas.

**Import:**
```go
import "github.com/Dorico-Dynamics/txova-go-validation/geo"
```

#### Mozambique Bounds

```go
const (
    MozambiqueMinLat = -26.9
    MozambiqueMaxLat = -10.3
    MozambiqueMinLon = 30.2
    MozambiqueMaxLon = 41.0
)
```

#### Service Areas

| Area | City |
|------|------|
| `maputo` | Maputo City |
| `matola` | Matola |
| `beira` | Beira |

#### Coordinate Validation

```go
// Validate global coordinate bounds
err := geo.ValidateCoordinates(-25.969, 32.573)

// Validate within Mozambique
err := geo.ValidateInMozambique(-25.969, 32.573)

// Quick checks
geo.IsInMozambique(-25.969, 32.573) // true
geo.IsInMozambique(40.0, -74.0)     // false (New York)
```

#### Service Area Validation

```go
// Check if in any service area
geo.IsInServiceArea(-25.95, 32.5) // true (Maputo area)

// Validate specific service area
err := geo.ValidateServiceArea(-25.95, 32.5, "maputo")

// Validate any service area
err := geo.ValidateAnyServiceArea(-25.95, 32.5)

// Find which service area
area := geo.FindServiceArea(-25.95, 32.5) // "maputo"
area := geo.FindServiceArea(-19.84, 34.84) // "beira"
area := geo.FindServiceArea(-20.0, 35.0)   // "" (not in any)

// Get all service areas
areas := geo.GetServiceAreas() // ["maputo", "matola", "beira"]

// Get service area config
config := geo.GetServiceArea("maputo")
// config.MinLat, config.MaxLat, config.MinLon, config.MaxLon
```

#### Distance Calculation

```go
// Calculate distance between two points (Haversine formula)
distanceKM, err := geo.CalculateDistance(
    -25.969, 32.573,  // Maputo
    -19.84, 34.84,    // Beira
)
// distanceKM ≈ 730 km
```

---

### vehicle Package

Mozambique vehicle validation including license plates and years.

**Import:**
```go
import "github.com/Dorico-Dynamics/txova-go-validation/vehicle"
```

#### License Plate Formats

| Format | Pattern | Example |
|--------|---------|---------|
| Standard | `AAA-NNN-LL` | `AAA-123-MC` |
| Old | `LL-NN-NN` | `MC-12-34` |

#### Province Codes

| Code | Province |
|------|----------|
| MC | Maputo City |
| MP | Maputo Province |
| GZ | Gaza |
| IB | Inhambane |
| SF | Sofala |
| MN | Manica |
| TT | Tete |
| ZB | Zambezia |
| NP | Nampula |
| CA | Cabo Delgado |
| NS | Niassa |

#### Plate Validation

```go
// Validate plate
err := vehicle.ValidatePlate("AAA-123-MC")  // nil (valid)
err := vehicle.ValidatePlate("MC-12-34")    // nil (valid old format)
err := vehicle.ValidatePlate("XXX-123-XX")  // error (invalid province)

// Quick check
vehicle.IsValidPlate("AAA-123-MC") // true

// Format detection
vehicle.IsStandardFormat("AAA-123-MC") // true
vehicle.IsOldFormat("MC-12-34")        // true
```

#### Plate Normalization

```go
// Normalize to standard format with dashes
normalized, err := vehicle.NormalizePlate("aaa123mp")
// normalized = "AAA-123-MP"

normalized, err := vehicle.NormalizePlate("AAA 123 MC")
// normalized = "AAA-123-MC"
```

#### Province Information

```go
// Get province code
code := vehicle.GetProvince("AAA-123-MP") // "MP"

// Get province name
name := vehicle.GetProvinceName("AAA-123-MP") // "Maputo Province"
name := vehicle.GetProvinceName("AAA-123-MC") // "Maputo City"
```

#### Vehicle Year Validation

```go
const MinVehicleYear = 2010
// Maximum: current year + 1

// Validate year
err := vehicle.ValidateYear(2020)  // nil (valid)
err := vehicle.ValidateYear(2005)  // error (before 2010)
err := vehicle.ValidateYear(2030)  // error (too far in future)

// Quick check
vehicle.IsValidYear(2020) // true
```

---

### ride Package

Ride-specific validation for PIN, distance, fare, and locations.

**Import:**
```go
import "github.com/Dorico-Dynamics/txova-go-validation/ride"
```

#### Constants

```go
const (
    MinDistanceKM                = 0.5
    MaxDistanceKM                = 200.0
    MinFareCentavos              = 5000     // 50 MZN
    MaxFareCentavos              = 5000000  // 50,000 MZN
    MinPickupDropoffSeparationKM = 0.1      // 100 meters
)
```

#### PIN Validation

4-digit PIN with security rules:
- No sequential patterns: `1234`, `4321`, `5678`, `8765`, etc.
- No repeated digits: `1111`, `2222`, `0000`, `9999`, etc.

```go
// Valid PINs
ride.IsValidPIN("7392") // true
ride.IsValidPIN("4826") // true

// Invalid PINs
ride.IsValidPIN("1234") // false (sequential)
ride.IsValidPIN("4321") // false (reverse sequential)
ride.IsValidPIN("1111") // false (repeated)
ride.IsValidPIN("123")  // false (too short)

// With error details
err := ride.ValidatePIN("1234")
// Error: PIN cannot be sequential
```

#### Distance Validation

```go
// Valid: 0.5 to 200 km
err := ride.ValidateDistance(15.5) // nil
err := ride.ValidateDistance(0.3)  // error (too short)
err := ride.ValidateDistance(250)  // error (too long)

ride.IsValidDistance(15.5) // true
```

#### Fare Validation

All amounts in centavos (MZN * 100).

```go
// Valid: 5,000 to 5,000,000 centavos (50 to 50,000 MZN)
err := ride.ValidateFare(10000)    // nil (100 MZN)
err := ride.ValidateFare(3000)     // error (below minimum)
err := ride.ValidateFare(6000000)  // error (above maximum)

ride.IsValidFare(10000) // true

// Using Money type from txova-go-types
err := ride.ValidateFareMoney(moneyAmount)
```

#### Pickup/Dropoff Validation

Ensures minimum 100m separation between pickup and dropoff.

```go
// Using coordinates
err := ride.ValidatePickupDropoff(
    -25.969, 32.573,  // Pickup
    -25.970, 32.580,  // Dropoff
)

// Using Location types from txova-go-types
err := ride.ValidatePickupDropoffLocations(pickupLocation, dropoffLocation)
```

#### Fare Estimation

```go
// Calculate: base + (distance * perKM)
fare := ride.CalculateEstimatedFare(
    10.5,   // Distance in km
    5000,   // Base fare: 50 MZN
    1000,   // Per km: 10 MZN
)
// fare = 5000 + (10.5 * 1000) = 15500 centavos (155 MZN)
```

---

### rating Package

Rating and review validation with profanity detection.

**Import:**
```go
import "github.com/Dorico-Dynamics/txova-go-validation/rating"
```

#### Constants

```go
const (
    MinReviewLength = 0
    MaxReviewLength = 500  // Unicode characters
)
```

#### Rating Validation

```go
// Valid ratings: 1 to 5
err := rating.ValidateRating(4)  // nil
err := rating.ValidateRating(0)  // error
err := rating.ValidateRating(6)  // error

rating.IsValidRating(4) // true
```

#### Review Text Validation

```go
// Validates length (0-500 characters)
err := rating.ValidateReviewText("Great service!")  // nil
err := rating.ValidateReviewText(veryLongText)      // error if >500 chars

rating.IsValidReviewText("Great service!") // true
```

#### Review Sanitization

```go
// Strips HTML, normalizes whitespace, trims
sanitized := rating.SanitizeReviewText("  <b>Great</b>   service!  ")
// sanitized = "Great service!"
```

#### Profanity Detection

Detects common profanity in English and Portuguese. Flags for moderation rather than rejecting.

```go
// Check for profanity
hasProfanity := rating.CheckProfanity("This was great service")  // false
hasProfanity := rating.CheckProfanity("This was shit service")   // true
```

#### Combined Processing

```go
// Validate and sanitize in one step
sanitized, err := rating.ValidateAndSanitizeReview("  <b>Great!</b>  ")
// sanitized = "Great!"

// Full processing with profanity check
result, err := rating.ProcessReview("  <b>Great driver!</b>  ")
if err == nil {
    fmt.Println(result.Text)            // "Great driver!"
    fmt.Println(result.HasProfanity)    // false
    fmt.Println(result.RequiresReview)  // false
    fmt.Println(result.OriginalLength)  // 24
    fmt.Println(result.SanitizedLength) // 13
}
```

---

### document Package

Document and file upload validation.

**Import:**
```go
import "github.com/Dorico-Dynamics/txova-go-validation/document"
```

#### Document Types

| Type | Constant | Max Size | Allowed Formats |
|------|----------|----------|-----------------|
| Driver License | `DocTypeDriverLicense` | 5 MB | jpg, jpeg, png, pdf |
| Vehicle Registration | `DocTypeVehicleRegistration` | 5 MB | jpg, jpeg, png, pdf |
| Insurance | `DocTypeInsurance` | 5 MB | jpg, jpeg, png, pdf |
| ID Card | `DocTypeIDCard` | 5 MB | jpg, jpeg, png, pdf |
| Profile Photo | `DocTypeProfilePhoto` | 2 MB | jpg, jpeg, png |
| Vehicle Photo | `DocTypeVehiclePhoto` | 5 MB | jpg, jpeg, png |

#### Image Constraints

```go
const (
    MinImageWidth  = 200
    MinImageHeight = 200
    MaxImageWidth  = 4096
    MaxImageHeight = 4096
    MinAspectRatio = 0.25  // 1:4
    MaxAspectRatio = 4.0   // 4:1
)
```

#### Document Type Validation

```go
// Validate type
err := document.ValidateDocType("driver_license") // nil
err := document.ValidateDocType("invalid_type")   // error

document.IsValidDocType("driver_license") // true

// Get all types
types := document.AllDocTypes()
// ["driver_license", "vehicle_registration", "insurance", "id_card", "profile_photo", "vehicle_photo"]

// Check if image type (vs document type)
document.IsImageType("profile_photo")  // true
document.IsImageType("driver_license") // false
```

#### File Size Validation

```go
// Validate size for document type
err := document.ValidateFileSize(1024*1024, "profile_photo")   // nil (1MB < 2MB limit)
err := document.ValidateFileSize(3*1024*1024, "profile_photo") // error (3MB > 2MB limit)

// Get max size for type
maxSize := document.GetMaxFileSize("profile_photo")  // 2097152 (2MB)
maxSize := document.GetMaxFileSize("driver_license") // 5242880 (5MB)
```

#### Format Validation

```go
// Get allowed formats
formats := document.GetAllowedFormats("driver_license")
// ["jpg", "jpeg", "png", "pdf"]

formats := document.GetAllowedFormats("profile_photo")
// ["jpg", "jpeg", "png"]

// Check if format allowed
document.IsAllowedFormat("pdf", "driver_license") // true
document.IsAllowedFormat("pdf", "profile_photo")  // false

// Validate format
err := document.ValidateFormat("jpg", "profile_photo")  // nil
err := document.ValidateFormat("pdf", "profile_photo")  // error
```

#### MIME Type Validation

```go
// Validate MIME type matches extension
err := document.ValidateMIMEType("image/jpeg", "jpg")  // nil
err := document.ValidateMIMEType("image/png", "jpg")   // error (mismatch)
```

#### Image Validation

```go
// Validate dimensions (200-4096 pixels)
err := document.ValidateImageDimensions(1920, 1080) // nil
err := document.ValidateImageDimensions(100, 100)   // error (too small)

// Validate aspect ratio (1:4 to 4:1)
err := document.ValidateAspectRatio(1920, 1080)  // nil (16:9)
err := document.ValidateAspectRatio(100, 1000)   // error (1:10)

// Combined image validation
err := document.ValidateImage(
    1920,           // width
    1080,           // height
    1024*1024,      // size in bytes
    "profile_photo", // document type
)
```

---

### sanitize Package

Input sanitization utilities with chainable API.

**Import:**
```go
import "github.com/Dorico-Dynamics/txova-go-validation/sanitize"
```

#### Standalone Functions

**Whitespace:**
```go
sanitize.TrimWhitespace("  hello  ")      // "hello"
sanitize.NormalizeSpaces("hello   world") // "hello world"
```

**HTML:**
```go
sanitize.StripHTML("<b>hello</b>")           // "hello"
sanitize.EscapeHTML("<script>alert()</script>") // "&lt;script&gt;alert()&lt;/script&gt;"
```

**Name/Email:**
```go
sanitize.NormalizeName("  JOHN   doe  ")       // "John Doe"
sanitize.NormalizeEmail("  USER@EXAMPLE.COM ") // "user@example.com"
```

**Character Filtering:**
```go
sanitize.RemoveNonPrintable("hello\x00world")  // "helloworld"
sanitize.RemoveControlChars("hello\x00world")  // "helloworld"
sanitize.ToUppercase("hello")                  // "HELLO"
sanitize.ToLowercase("HELLO")                  // "hello"
sanitize.RemoveDigits("abc123")                // "abc"
sanitize.KeepDigits("abc123")                  // "123"
sanitize.KeepAlphanumeric("abc-123!")          // "abc123"
```

#### Function Chaining

```go
result := sanitize.Chain(
    "  <B>HELLO</B>  ",
    sanitize.StripHTML,
    sanitize.TrimWhitespace,
    sanitize.ToLowercase,
)
// result = "hello"
```

#### Builder Pattern

```go
result := sanitize.NewSanitizer().
    StripHTML().
    RemoveNonPrintable().
    NormalizeSpaces().
    TrimWhitespace().
    Apply("  <b>Hello</b>   World  ")
// result = "Hello World"

// With custom function
result := sanitize.NewSanitizer().
    StripHTML().
    Custom(func(s string) string {
        return strings.ReplaceAll(s, "bad", "good")
    }).
    Apply("<b>This is bad</b>")
// result = "This is good"
```

#### Pre-built Sanitizers

```go
// Text: StripHTML -> RemoveNonPrintable -> NormalizeSpaces
text := sanitize.TextSanitizer().Apply("<p>Hello   World</p>")
// "Hello World"

// Names: StripHTML -> RemoveNonPrintable -> NormalizeName
name := sanitize.NameSanitizer().Apply("  <b>JOHN</b>   doe  ")
// "John Doe"

// Emails: TrimWhitespace -> NormalizeEmail (lowercase)
email := sanitize.EmailSanitizer().Apply("  USER@EXAMPLE.COM  ")
// "user@example.com"

// Phone: KeepDigits only
phone := sanitize.PhoneSanitizer().Apply("+258 84-123-4567")
// "258841234567"
```

---

### struct Package

Struct validation using go-playground/validator with custom Txova tags.

**Import:**
```go
import structval "github.com/Dorico-Dynamics/txova-go-validation/struct"
```

#### Custom Validation Tags

| Tag | Description | Valid Examples |
|-----|-------------|----------------|
| `mz_phone` | Mozambique phone number | `+258841234567`, `841234567` |
| `mz_plate` | Mozambique license plate | `AAA-123-MC`, `MC-12-34` |
| `mz_location` | Location within Mozambique | struct with Lat/Lon, `[-25.969, 32.573]` |
| `txova_pin` | 4-digit PIN (no sequential/repeated) | `7392`, `4826` |
| `txova_money` | Positive money amount | any positive number |
| `txova_rating` | Rating 1-5 | `1`, `2`, `3`, `4`, `5` |
| `txova_vehicle_year` | Year 2010 to current+1 | `2015`, `2020`, `2025` |

#### Standard Tags (go-playground/validator)

```go
validate:"required"           // Field is required
validate:"omitempty"          // Field is optional
validate:"email"              // Valid email format
validate:"url"                // Valid URL format
validate:"min=N"              // Minimum length/value
validate:"max=N"              // Maximum length/value
validate:"len=N"              // Exact length
validate:"oneof=a b c"        // One of specified values
validate:"gt=N"               // Greater than
validate:"gte=N"              // Greater than or equal
validate:"lt=N"               // Less than
validate:"lte=N"              // Less than or equal
```

#### Struct Validation

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Phone    string `json:"phone" validate:"required,mz_phone"`
    Password string `json:"password" validate:"required,min=8"`
}

user := CreateUserRequest{
    Name:     "João Silva",
    Email:    "joao@example.com",
    Phone:    "+258841234567",
    Password: "securepass123",
}

errs := structval.Validate(user)
if errs != nil {
    // Handle validation errors
    for _, err := range errs {
        fmt.Printf("Field: %s, Code: %s, Message: %s\n",
            err.Field, err.Code, err.Message)
    }
}
```

#### Location Validation

```go
// Using struct with Lat/Lon fields
type Location struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

type RideRequest struct {
    Pickup  Location `json:"pickup" validate:"required,mz_location"`
    Dropoff Location `json:"dropoff" validate:"required,mz_location"`
}

// Or using slice/array
type RideRequestAlt struct {
    Pickup  []float64 `json:"pickup" validate:"required,mz_location"`
    Dropoff []float64 `json:"dropoff" validate:"required,mz_location"`
}
```

#### Complete Ride Example

```go
type CreateRideRequest struct {
    PassengerPhone string   `json:"passenger_phone" validate:"required,mz_phone"`
    DriverPhone    string   `json:"driver_phone" validate:"omitempty,mz_phone"`
    Pickup         Location `json:"pickup" validate:"required,mz_location"`
    Dropoff        Location `json:"dropoff" validate:"required,mz_location"`
    PIN            string   `json:"pin" validate:"required,txova_pin"`
    Fare           int64    `json:"fare" validate:"required,txova_money"`
    VehicleYear    int      `json:"vehicle_year" validate:"required,txova_vehicle_year"`
    Plate          string   `json:"plate" validate:"required,mz_plate"`
    Rating         int      `json:"rating" validate:"omitempty,txova_rating"`
}

ride := CreateRideRequest{
    PassengerPhone: "+258841234567",
    Pickup:         Location{Lat: -25.969, Lon: 32.573},
    Dropoff:        Location{Lat: -25.970, Lon: 32.580},
    PIN:            "7392",
    Fare:           15000,
    VehicleYear:    2020,
    Plate:          "AAA-123-MC",
}

errs := structval.Validate(ride)
```

#### Single Field Validation

```go
// Validate single value against tag
errs := structval.ValidateVar("+258841234567", "required,mz_phone")
errs := structval.ValidateVar(4, "required,txova_rating")
```

#### Custom Validators

```go
// Register custom validation function
err := structval.RegisterValidation("my_custom", func(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    return value != "forbidden"
})

type MyStruct struct {
    Field string `validate:"my_custom"`
}
```

---

## Integration Patterns

### Request Validation Pipeline

```go
func validateCreateRideRequest(raw RawRequest) (*CreateRideRequest, error) {
    // 1. Sanitize inputs
    phone := sanitize.PhoneSanitizer().Apply(raw.Phone)

    // 2. Normalize phone
    normalizedPhone, err := phone.Normalize(phone)
    if err != nil {
        return nil, fmt.Errorf("invalid phone: %w", err)
    }

    // 3. Create validated struct
    req := &CreateRideRequest{
        Phone:   normalizedPhone,
        Pickup:  raw.Pickup,
        Dropoff: raw.Dropoff,
        PIN:     raw.PIN,
        Fare:    raw.Fare,
    }

    // 4. Validate struct
    errs := structval.Validate(req)
    if errs != nil {
        return nil, errs.ToError()
    }

    // 5. Additional business validations
    err = ride.ValidatePickupDropoff(
        req.Pickup.Lat, req.Pickup.Lon,
        req.Dropoff.Lat, req.Dropoff.Lon,
    )
    if err != nil {
        return nil, err
    }

    return req, nil
}
```

### Review Processing

```go
func processReview(reviewText string, ratingValue int) (*ProcessedReview, error) {
    // Validate rating
    if err := rating.ValidateRating(ratingValue); err != nil {
        return nil, err
    }

    // Process review text
    result, err := rating.ProcessReview(reviewText)
    if err != nil {
        return nil, err
    }

    return &ProcessedReview{
        Rating:         ratingValue,
        Text:           result.Text,
        RequiresReview: result.HasProfanity || result.RequiresReview,
    }, nil
}
```

### Document Upload Validation

```go
func validateDocumentUpload(
    docType string,
    fileSize int64,
    mimeType string,
    extension string,
    width, height int,
) error {
    // Validate document type
    if err := document.ValidateDocType(docType); err != nil {
        return err
    }

    // Validate file size
    if err := document.ValidateFileSize(fileSize, docType); err != nil {
        return err
    }

    // Validate format
    if err := document.ValidateFormat(extension, docType); err != nil {
        return err
    }

    // Validate MIME type
    if err := document.ValidateMIMEType(mimeType, extension); err != nil {
        return err
    }

    // If image, validate dimensions
    if document.IsImageType(docType) {
        if err := document.ValidateImageDimensions(width, height); err != nil {
            return err
        }
        if err := document.ValidateAspectRatio(width, height); err != nil {
            return err
        }
    }

    return nil
}
```

---

## Error Handling

### Field-Level Handling

```go
errs := structval.Validate(request)
if errs != nil {
    // Handle specific fields
    if errs.HasField("phone") {
        phoneErrs := errs.GetByField("phone")
        for _, e := range phoneErrs {
            log.Printf("Phone error: %s", e.Message)
        }
    }

    // Handle by error type
    formatErrs := errs.GetByCode(valerrors.CodeInvalidFormat)
    requiredErrs := errs.GetByCode(valerrors.CodeRequired)
}
```

### Converting to Standard Error

```go
errs := structval.Validate(request)
if errs != nil {
    // Convert to standard error interface
    return errs.ToError()
}
```

---

## API Response Examples

### Validation Error Response

```go
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)

    errs := structval.Validate(req)
    if errs != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": false,
            "errors":  errs,
        })
        return
    }

    // Process valid request...
}
```

**Response:**
```json
{
  "success": false,
  "errors": [
    {
      "field": "email",
      "code": "REQUIRED",
      "message": "email is required"
    },
    {
      "field": "phone",
      "code": "INVALID_FORMAT",
      "message": "phone must be a valid Mozambique phone number"
    }
  ]
}
```

### Grouped by Field

```go
func groupErrorsByField(errs valerrors.ValidationErrors) map[string][]string {
    result := make(map[string][]string)
    for _, field := range errs.Fields() {
        fieldErrs := errs.GetByField(field)
        for _, e := range fieldErrs {
            result[field] = append(result[field], e.Message)
        }
    }
    return result
}
```

**Response:**
```json
{
  "success": false,
  "errors": {
    "email": ["email is required"],
    "phone": ["phone must be a valid Mozambique phone number"],
    "password": ["password must be at least 8 characters"]
  }
}
```

---

## Best Practices

1. **Sanitize before validating** - Always sanitize user input before validation
2. **Use struct tags for declarative validation** - Cleaner than imperative checks
3. **Combine struct validation with business rules** - Struct tags for format, functions for business logic
4. **Return all errors at once** - Better UX than one error at a time
5. **Use JSON tags for field names** - Consistent naming between API and validation errors
6. **Normalize early** - Normalize phones, plates, etc. at the entry point
7. **Use pre-built sanitizers** - Consistent sanitization across the application
