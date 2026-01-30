# txova-go-validation Execution Plan

**Version:** 1.0  
**Module:** `github.com/Dorico-Dynamics/txova-go-validation`  
**Target Test Coverage:** >90%  
**External Dependencies:** `go-playground/validator/v10`  
**Internal Dependencies:** `txova-go-types`

---

## Dependency: txova-go-types Integration

This library depends on `txova-go-types` for domain invariant types. The validation library handles **business rules** while the types library enforces **domain invariants**.

| txova-go-types | Used In | Purpose |
|----------------|---------|---------|
| `contact.PhoneNumber` | `phone` package | Parse raw input, leverage existing normalization |
| `contact.Operator` | `phone` package | Operator identification (derived from prefix) |
| `vehicle.LicensePlate` | `vehicle` package | Parse raw input, leverage format validation |
| `vehicle.ProvinceCode` | `vehicle` package | Province extraction from plates |
| `ride.PIN` | `ride` package | PIN format validation (sequential/repeated checks) |
| `rating.Rating` | `rating` package | Rating range validation (1-5) |
| `geo.Location` | `geo` package | Coordinate types for bounds checking |
| `geo.BoundingBox` | `geo` package | Service area containment checks |
| `money.Money` | `ride` package | Fare range validation |
| `enums.*` | `struct` package | Enum validation in custom validators |

**Key Principle**: Use types library for parsing/construction (domain invariants), use validation library for business rules (contextual validation).

---

## Phase 1: Project Setup & Core Infrastructure (Week 1) - COMPLETE

### 1.1 Project Initialization
- [x] Initialize Go module with `go mod init github.com/Dorico-Dynamics/txova-go-validation`
- [x] Add dependency: `go get github.com/Dorico-Dynamics/txova-go-types@v1.1.1`
- [x] Add dependency: `go get github.com/go-playground/validator/v10@v10.30.1`
- [x] Create directory structure for all packages
- [x] Set up `.gitignore` for Go projects
- [ ] Configure linting with strict rules

### 1.2 Package: `errors` - Validation Error Types - COMPLETE (100% coverage)
- [x] Define `ValidationError` struct (field, code, message, value)
- [x] Define error code constants (REQUIRED, INVALID_FORMAT, OUT_OF_RANGE, TOO_SHORT, TOO_LONG, INVALID_OPTION, OUTSIDE_SERVICE_AREA)
- [x] Implement `Error()` method for error interface
- [x] Implement `ValidationErrors` slice type with helper methods
- [x] Implement JSON marshaling for API responses
- [x] Implement `HasField(field string)` helper method
- [x] Implement `GetByField(field string)` helper method
- [x] Write comprehensive tests for error types

**Deliverables:**
- [x] `errors/` package with structured validation errors
- [x] Test suite covering all error scenarios (100% coverage)

---

## Phase 2: Phone & Geographic Validation (Week 2) - COMPLETE

### 2.1 Package: `phone` - Phone Number Validation - COMPLETE (90.2% coverage)
- [x] Implement `Validate(input string) bool` - returns true if valid Mozambique number
- [x] Implement `Normalize(input string) (string, error)` - normalizes to +258XXXXXXXXX
- [x] Implement `IdentifyOperator(input string) string` - returns operator name
- [x] Handle all input formats: local (841234567), international (+258841234567), with country code (258841234567), with 00 prefix (00258841234567), with spaces (84 123 4567)
- [x] Strip non-digit characters except leading +
- [x] Validate prefix is 82, 83, 84, 85, 86, or 87
- [x] Use `types/contact.PhoneNumber` for final parsing and `types/contact.Operator` for operator identification
- [x] Write tests covering all input format variations

**Deliverables:**
- [x] `phone/` package with validation and normalization
- [x] Test suite with format normalization matrix

### 2.2 Package: `geo` - Geographic Validation - COMPLETE (100% coverage)
- [x] Implement `ValidateCoordinates(lat, lon float64) error` - checks global valid ranges
- [x] Implement `ValidateInMozambique(lat, lon float64) error` - checks Mozambique bounds
- [x] Implement `ValidateServiceArea(lat, lon float64, area string) error` - checks city service areas
- [x] Define service area bounds for Maputo, Matola, Beira
- [x] Use `types/geo.Location` for coordinate types
- [x] Use `types/geo.DistanceKM` for distance calculations
- [x] Return specific error codes (OUTSIDE_SERVICE_AREA vs OUT_OF_RANGE)
- [x] Implement `GetServiceAreas() []string` - returns list of active service areas
- [x] Write tests with real Mozambique coordinates

**Deliverables:**
- [x] `geo/` package with coordinate and service area validation
- [x] Test suite with boundary edge cases

---

## Phase 3: Vehicle & Ride Validation (Week 3)

### 3.1 Package: `vehicle` - Vehicle Validation
- [ ] Implement `ValidatePlate(input string) error` - validates Mozambique plate format
- [ ] Implement `NormalizePlate(input string) (string, error)` - normalizes to standard format
- [ ] Implement `ValidateYear(year int) error` - validates vehicle year (2010 to current+1)
- [ ] Support standard format (AAA-NNN-LL) and old format (AA-NN-NN)
- [ ] Use `types/vehicle.LicensePlate` for parsing and format validation
- [ ] Use `types/vehicle.ProvinceCode` for province validation
- [ ] Normalize: uppercase, add dashes if missing
- [ ] Write tests for both plate formats and year edge cases

**Deliverables:**
- [ ] `vehicle/` package with plate and year validation
- [ ] Test suite covering format variations

### 3.2 Package: `ride` - Ride Validation
- [ ] Implement `ValidatePIN(input string) error` - validates 4-digit PIN format
- [ ] Use `types/ride.PIN` for PIN parsing and invariant validation (no sequential, no repeated)
- [ ] Implement `ValidateDistance(km float64) error` - validates distance (0.5 to 200 km)
- [ ] Implement `ValidateFare(amount int64) error` - validates fare (50 to 50,000 MZN in centavos)
- [ ] Implement `ValidatePickupDropoff(pickup, dropoff types/geo.Location) error` - ensures minimum separation
- [ ] Use `types/money.Money` for fare validation
- [ ] Write tests for all validation rules

**Deliverables:**
- [ ] `ride/` package with ride request validation
- [ ] Test suite covering distance, fare, and PIN rules

---

## Phase 4: Rating & Document Validation (Week 4)

### 4.1 Package: `rating` - Rating Validation
- [ ] Implement `ValidateRating(value int) error` - validates 1-5 range
- [ ] Use `types/rating.Rating` for range validation
- [ ] Implement `ValidateReviewText(text string) error` - validates length (0-500 chars)
- [ ] Implement `SanitizeReviewText(text string) string` - strips HTML, normalizes whitespace
- [ ] Implement `CheckProfanity(text string) bool` - flags potential profanity for moderation
- [ ] Load profanity word list (Portuguese and English common terms)
- [ ] Write tests for rating validation and text sanitization

**Deliverables:**
- [ ] `rating/` package with rating and review validation
- [ ] Test suite covering edge cases and profanity detection

### 4.2 Package: `document` - Document Validation
- [ ] Define document type constants (driver_license, vehicle_registration, insurance, id_card, profile_photo, vehicle_photo)
- [ ] Implement `ValidateFileSize(size int64, docType string) error` - validates max size by type
- [ ] Implement `ValidateMIMEType(mimeType, extension string) error` - validates MIME matches extension
- [ ] Implement `ValidateImageDimensions(width, height int) error` - validates 200x200 to 4096x4096
- [ ] Implement `ValidateAspectRatio(width, height int) error` - validates 1:4 to 4:1 ratio
- [ ] Implement `GetAllowedFormats(docType string) []string` - returns allowed formats
- [ ] Write tests for all document validation rules

**Deliverables:**
- [ ] `document/` package with file validation
- [ ] Test suite covering size, format, and dimension rules

---

## Phase 5: Struct Validation & Custom Tags (Week 5)

### 5.1 Package: `struct` - Struct Validation Integration
- [ ] Initialize singleton validator instance with custom configuration
- [ ] Implement `Validate(s interface{}) ValidationErrors` - validates struct and returns errors
- [ ] Implement `RegisterValidation(tag string, fn validator.Func) error` - registers custom validator
- [ ] Map struct field names to JSON tag names in error responses
- [ ] Support nested struct validation
- [ ] Write tests for struct validation with various tag combinations

### 5.2 Custom Validation Tags
- [ ] Register `mz_phone` - uses `phone.Validate()`
- [ ] Register `mz_plate` - uses `vehicle.ValidatePlate()` with `types/vehicle.LicensePlate`
- [ ] Register `mz_location` - uses `geo.ValidateInMozambique()` with `types/geo.Location`
- [ ] Register `txova_pin` - uses `ride.ValidatePIN()` with `types/ride.PIN`
- [ ] Register `txova_money` - validates positive money amount using `types/money.Money`
- [ ] Register `txova_rating` - uses `rating.ValidateRating()` with `types/rating.Rating`
- [ ] Write tests for all custom validation tags

**Deliverables:**
- [ ] `struct/` package with go-playground/validator integration
- [ ] All custom validation tags registered and tested
- [ ] Field-to-JSON name mapping implemented

---

## Phase 6: Sanitization & Integration (Week 6)

### 6.1 Package: `sanitize` - Input Sanitization
- [ ] Implement `TrimWhitespace(s string) string` - removes leading/trailing whitespace
- [ ] Implement `NormalizeSpaces(s string) string` - collapses multiple spaces to single
- [ ] Implement `StripHTML(s string) string` - removes all HTML tags
- [ ] Implement `EscapeHTML(s string) string` - escapes HTML entities
- [ ] Implement `NormalizeName(s string) string` - capitalizes first letter of each word
- [ ] Implement `NormalizeEmail(s string) string` - lowercase and trim
- [ ] Implement `Chain(input string, fns ...func(string) string) string` - chainable sanitization
- [ ] All functions return new values, never modify input
- [ ] Write tests for all sanitization functions

**Deliverables:**
- [ ] `sanitize/` package with all sanitization utilities
- [ ] Test suite covering edge cases and chaining

### 6.2 Cross-Package Integration
- [ ] Verify all packages work together without circular dependencies
- [ ] Ensure validation error format is consistent across all packages
- [ ] Verify types library integration works correctly
- [ ] Test combined sanitization + validation workflows

### 6.3 Quality Assurance
- [ ] Run full test suite and verify >90% coverage
- [ ] Run linter and fix all issues
- [ ] Run `go vet` and address all warnings
- [ ] Run `gosec` security scan and fix issues
- [ ] Benchmark validation latency (target: <1ms per validation)

### 6.4 Documentation
- [ ] Add package-level documentation (doc.go) for each package
- [ ] Document all exported types and functions with godoc comments
- [ ] Update README.md with usage examples
- [ ] Create CHANGELOG.md with v1.0.0 release notes

**Deliverables:**
- [ ] Complete, tested library
- [ ] v1.0.0 release tagged and published
- [ ] >90% test coverage verified

---

## Success Criteria

| Criteria | Target |
|----------|--------|
| Test Coverage | >90% |
| False Positive Rate | <0.1% |
| Validation Latency | <1ms |
| Custom Validator Coverage | 100% |
| Linting Errors | 0 |
| `go vet` Warnings | 0 |
| `gosec` Issues | 0 |

---

## Package Dependency Order

```
errors (no internal deps)
    ↓
sanitize (no internal deps)
    ↓
phone (depends: errors, txova-go-types/contact)
    ↓
geo (depends: errors, txova-go-types/geo)
    ↓
vehicle (depends: errors, txova-go-types/vehicle)
    ↓
ride (depends: errors, geo, txova-go-types/ride, txova-go-types/money, txova-go-types/geo)
    ↓
rating (depends: errors, sanitize, txova-go-types/rating)
    ↓
document (depends: errors)
    ↓
struct (depends: all above, go-playground/validator)
```

---

## txova-go-types Usage Summary

| Validation Package | Types Used | Integration Point |
|--------------------|------------|-------------------|
| `phone` | `contact.PhoneNumber`, `contact.Operator` | `Normalize()` uses `ParsePhoneNumber()`, `IdentifyOperator()` uses `Operator()` method |
| `geo` | `geo.Location`, `geo.BoundingBox` | `ValidateServiceArea()` uses `BoundingBox.Contains()` |
| `vehicle` | `vehicle.LicensePlate`, `vehicle.ProvinceCode` | `NormalizePlate()` uses `ParseLicensePlate()` |
| `ride` | `ride.PIN`, `money.Money`, `geo.Location` | `ValidatePIN()` uses `ParsePIN()`, `ValidateFare()` uses `Money` comparison |
| `rating` | `rating.Rating` | `ValidateRating()` uses `NewRating()` constructor |
| `struct` | All above via custom tags | Custom validators wrap respective package functions |

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Types library version mismatch | Pin specific version in go.mod, test against that version |
| Validator library breaking changes | Pin go-playground/validator to v10, monitor releases |
| Profanity filter false positives | Use conservative word list, flag for moderation rather than reject |
| Service area bounds inaccuracy | Validate bounds with GPS data from local team |
| Slow validation performance | Benchmark critical paths, use regexp caching |
| Circular dependencies | Strict package layering, errors package at bottom |
