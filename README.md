# txova-go-validation

Input validation library providing Mozambique-specific validators, struct validation, and sanitization utilities for Txova services.

## Overview

`txova-go-validation` provides comprehensive input validation tailored for Mozambique, including phone number validation with operator detection, geographic bounds checking, vehicle plate validation, and integration with go-playground/validator.

**Module:** `github.com/txova/txova-go-validation`

## Features

- **Phone Validation** - Mozambique phone numbers with operator detection
- **Geographic Validation** - Coordinate and service area validation
- **Vehicle Validation** - License plate format and vehicle year validation
- **Ride Validation** - PIN, distance, and fare validation
- **Struct Validation** - Custom validator tags for go-playground/validator
- **Sanitization** - Input sanitization utilities

## Packages

| Package | Description |
|---------|-------------|
| `phone` | Mozambique phone number validation |
| `geo` | Geographic coordinate validation |
| `vehicle` | License plate and vehicle validation |
| `ride` | Ride-specific validations |
| `rating` | Rating and review validation |
| `document` | Document upload validation |
| `struct` | Struct validation with custom tags |
| `sanitize` | Input sanitization utilities |

## Installation

```bash
go get github.com/txova/txova-go-validation
```

## Usage

### Phone Validation

```go
import "github.com/txova/txova-go-validation/phone"

// Validate and normalize
normalized, err := phone.Normalize("841234567")
// Returns: "+258841234567"

// Identify operator
operator := phone.GetOperator("+258841234567")
// Returns: "Vodacom"

// Valid prefixes: 82, 83, 84, 85, 86, 87
```

### Geographic Validation

```go
import "github.com/txova/txova-go-validation/geo"

// Check if within Mozambique
valid := geo.InMozambique(-25.9692, 32.5732)

// Check if within service area
inService := geo.InServiceArea(-25.9692, 32.5732)

// Mozambique bounds: -26.9 to -10.3 lat, 30.2 to 41.0 lon
```

### Vehicle Validation

```go
import "github.com/txova/txova-go-validation/vehicle"

// Validate license plate
valid := vehicle.ValidatePlate("AAA-123-MZ")

// Normalize plate
normalized := vehicle.NormalizePlate("aaa 123 mz")
// Returns: "AAA-123-MZ"

// Province codes: MZ, MC, GA, IN, SO, MA, TE, ZA, NA, CD, NI
```

### Struct Validation

```go
import "github.com/txova/txova-go-validation/struct"

type CreateRideRequest struct {
    Phone    string  `json:"phone" validate:"required,mz_phone"`
    Pickup   Location `json:"pickup" validate:"required,mz_location"`
    PIN      string  `json:"pin" validate:"required,txova_pin"`
}

validator := struct.New()
errors := validator.Validate(request)
```

### Custom Validation Tags

| Tag | Description |
|-----|-------------|
| `mz_phone` | Mozambique phone number |
| `mz_plate` | Mozambique license plate |
| `mz_location` | Within Mozambique bounds |
| `txova_pin` | Valid ride PIN (4 digits, no sequential) |
| `txova_money` | Valid money amount |

### Sanitization

```go
import "github.com/txova/txova-go-validation/sanitize"

// Clean user input
name := sanitize.NormalizeName("  joão   silva  ")
// Returns: "João Silva"

email := sanitize.NormalizeEmail("  User@Example.COM  ")
// Returns: "user@example.com"

text := sanitize.StripHTML("<script>alert('xss')</script>Hello")
// Returns: "Hello"
```

## Error Codes

| Code | Description |
|------|-------------|
| `REQUIRED` | Field is required |
| `INVALID_FORMAT` | Format doesn't match pattern |
| `OUT_OF_RANGE` | Value outside allowed range |
| `TOO_SHORT` | Below minimum length |
| `TOO_LONG` | Exceeds maximum length |
| `INVALID_OPTION` | Not in allowed options |
| `OUTSIDE_SERVICE_AREA` | Location not serviceable |

## Dependencies

**Internal:**
- `txova-go-types`

**External:**
- `github.com/go-playground/validator/v10` - Struct validation

## Development

### Requirements

- Go 1.25+

### Testing

```bash
go test ./...
```

### Test Coverage Target

> 90%

## License

Proprietary - Dorico Dynamics
