# txova-go-validation

## Overview
Input validation library providing Mozambique-specific validators, struct validation, and sanitization utilities for all user input across Txova services.

**Module:** `github.com/txova/txova-go-validation`

---

## Packages

### `phone` - Phone Number Validation

**Mozambique Phone Rules:**
| Rule | Description |
|------|-------------|
| Country code | +258 (required in storage) |
| Length | 9 digits after country code |
| Valid prefixes | 82, 83, 84, 85, 86, 87 |
| Format | +258XXXXXXXXX |

**Accepted Input Formats:**
| Input | Normalized |
|-------|------------|
| 841234567 | +258841234567 |
| +258841234567 | +258841234567 |
| 258841234567 | +258841234567 |
| 00258841234567 | +258841234567 |
| 84 123 4567 | +258841234567 |

**Operators by Prefix:**
| Prefix | Operator |
|--------|----------|
| 82 | Vodacom |
| 84 | Vodacom |
| 85 | Vodacom |
| 83 | Movitel |
| 86 | Movitel |
| 87 | Tmcel |

**Requirements:**
- Validate: returns bool
- Normalize: returns normalized string or error
- Identify operator: returns operator name
- Strip all non-digit characters except leading +

---

### `geo` - Geographic Validation

**Coordinate Validation:**
| Rule | Value |
|------|-------|
| Latitude range | -90 to +90 |
| Longitude range | -180 to +180 |
| Precision | 6 decimal places |

**Mozambique Bounds:**
| Boundary | Value |
|----------|-------|
| Min latitude | -26.9 |
| Max latitude | -10.3 |
| Min longitude | 30.2 |
| Max longitude | 41.0 |

**Service Area Validation:**
| City | Bounds |
|------|--------|
| Maputo | -26.1 to -25.8 lat, 32.3 to 32.7 lon |
| Matola | -26.0 to -25.9 lat, 32.3 to 32.5 lon |
| Beira | -19.9 to -19.7 lat, 34.8 to 34.9 lon |

**Requirements:**
- Validate coordinates are within valid ranges
- Validate location is within Mozambique
- Validate location is within active service area
- Return specific error for out-of-bounds (not generic validation error)

---

### `vehicle` - Vehicle Validation

**License Plate Format (Mozambique):**
| Format | Pattern | Example |
|--------|---------|---------|
| Standard | AAA-NNN-LL | AAA-123-MZ |
| Old format | AA-NN-NN | MC-12-34 |

Where: A=letter, N=number, L=province code

**Province Codes:**
| Code | Province |
|------|----------|
| MZ | Maputo City |
| MC | Maputo Province |
| GA | Gaza |
| IN | Inhambane |
| SO | Sofala |
| MA | Manica |
| TE | Tete |
| ZA | Zambezia |
| NA | Nampula |
| CD | Cabo Delgado |
| NI | Niassa |

**Vehicle Year Validation:**
| Rule | Value |
|------|-------|
| Minimum year | 2010 |
| Maximum year | current + 1 |

**Requirements:**
- Normalize plate to uppercase, standard format
- Validate plate format matches pattern
- Validate vehicle year is within range
- Support both old and new plate formats

---

### `ride` - Ride Validation

**PIN Validation:**
| Rule | Value |
|------|-------|
| Length | 4 digits |
| Format | Numeric only |
| No sequential | Reject 1234, 4321 |
| No repeated | Reject 1111, 2222 |

**Distance Validation:**
| Rule | Value |
|------|-------|
| Minimum | 0.5 km |
| Maximum | 200 km |

**Fare Validation:**
| Rule | Value |
|------|-------|
| Minimum | 50 MZN |
| Maximum | 50,000 MZN |

**Requirements:**
- Validate PIN meets security rules
- Validate pickup != dropoff (minimum distance)
- Validate estimated fare within range
- Validate service type is available in area

---

### `rating` - Rating Validation

| Rule | Value |
|------|-------|
| Minimum | 1 |
| Maximum | 5 |
| Type | Integer only |

**Review Text Validation:**
| Rule | Value |
|------|-------|
| Min length | 0 (optional) |
| Max length | 500 characters |
| Profanity filter | Required |

**Requirements:**
- Validate rating is 1-5 integer
- Sanitize review text (strip HTML)
- Flag potential profanity for moderation

---

### `document` - Document Validation

**Supported Document Types:**
| Type | Max Size | Formats |
|------|----------|---------|
| Driver's license | 5 MB | JPG, PNG, PDF |
| Vehicle registration | 5 MB | JPG, PNG, PDF |
| Insurance | 5 MB | JPG, PNG, PDF |
| ID card | 5 MB | JPG, PNG, PDF |
| Profile photo | 2 MB | JPG, PNG |
| Vehicle photo | 5 MB | JPG, PNG |

**Image Validation:**
| Rule | Value |
|------|-------|
| Min dimensions | 200x200 px |
| Max dimensions | 4096x4096 px |
| Aspect ratio | 1:4 to 4:1 |

**Requirements:**
- Validate file size before upload
- Validate MIME type matches extension
- Validate image dimensions
- Reject corrupted files

---

### `struct` - Struct Validation

**Integration with go-playground/validator:**

**Common Validation Tags:**
| Tag | Description |
|-----|-------------|
| required | Field must be present |
| email | Valid email format |
| min | Minimum value/length |
| max | Maximum value/length |
| oneof | Value from allowed list |
| uuid | Valid UUID format |

**Custom Validation Tags:**
| Tag | Description |
|-----|-------------|
| mz_phone | Mozambique phone number |
| mz_plate | Mozambique license plate |
| mz_location | Within Mozambique bounds |
| txova_pin | Valid ride PIN |
| txova_money | Valid money amount |

**Requirements:**
- Register custom validators on init
- Return structured validation errors
- Map field names to JSON names in errors
- Support nested struct validation

---

### `sanitize` - Input Sanitization

| Function | Description |
|----------|-------------|
| TrimWhitespace | Remove leading/trailing spaces |
| NormalizeSpaces | Collapse multiple spaces |
| StripHTML | Remove HTML tags |
| EscapeHTML | Escape HTML entities |
| NormalizeName | Capitalize first letter of each word |
| NormalizeEmail | Lowercase and trim |

**Requirements:**
- Apply sanitization before validation
- Provide chainable sanitizer
- Never modify original input (return new value)

---

## Validation Error Format

| Field | Description |
|-------|-------------|
| field | JSON field name |
| code | Validation error code |
| message | Human-readable message |
| value | Invalid value (masked if sensitive) |

**Error Codes:**
| Code | Description |
|------|-------------|
| REQUIRED | Field is required |
| INVALID_FORMAT | Format doesn't match pattern |
| OUT_OF_RANGE | Value outside allowed range |
| TOO_SHORT | Below minimum length |
| TOO_LONG | Exceeds maximum length |
| INVALID_OPTION | Not in allowed options |
| OUTSIDE_SERVICE_AREA | Location not serviceable |

---

## Dependencies

**Internal:**
- `txova-go-types`

**External:**
- `github.com/go-playground/validator/v10` — Struct validation

---

## Success Metrics
| Metric | Target |
|--------|--------|
| Test coverage | > 90% |
| False positive rate | < 0.1% |
| Validation latency | < 1ms |
| Custom validator coverage | 100% |
