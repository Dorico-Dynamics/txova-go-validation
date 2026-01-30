package structval

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

// Test structs for validation

type UserRegistration struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,mz_phone"`
	Password string `json:"password" validate:"required,min=8"`
}

type VehicleInfo struct {
	Plate string `json:"plate" validate:"required,mz_plate"`
	Year  int    `json:"year" validate:"required,txova_vehicle_year"`
	Color string `json:"color" validate:"required,oneof=white black silver red blue"`
}

type RideRequest struct {
	PIN    string   `json:"pin" validate:"required,txova_pin"`
	Fare   int64    `json:"fare" validate:"required,txova_money"`
	Rating int      `json:"rating" validate:"omitempty,txova_rating"`
	Pickup Location `json:"pickup" validate:"required,mz_location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type NestedStruct struct {
	User    UserRegistration `json:"user" validate:"required"`
	Vehicle VehicleInfo      `json:"vehicle" validate:"required"`
}

type OptionalFields struct {
	Name   string `json:"name" validate:"omitempty,min=2"`
	Phone  string `json:"phone" validate:"omitempty,mz_phone"`
	Rating int    `json:"rating" validate:"omitempty,txova_rating"`
}

func TestValidate_ValidStruct(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "valid user registration",
			data: UserRegistration{
				Name:     "João Silva",
				Email:    "joao@example.com",
				Phone:    "+258841234567",
				Password: "securepass123",
			},
		},
		{
			name: "valid vehicle info",
			data: VehicleInfo{
				Plate: "AAA-123-MP",
				Year:  2022,
				Color: "white",
			},
		},
		{
			name: "valid ride request",
			data: RideRequest{
				PIN:  "7392", // Non-sequential, non-repeated
				Fare: 10000,
				Pickup: Location{
					Lat: -25.95,
					Lon: 32.58,
				},
			},
		},
		{
			name: "valid ride with rating",
			data: RideRequest{
				PIN:    "4826", // Non-sequential, non-repeated
				Fare:   10000,
				Rating: 5,
				Pickup: Location{
					Lat: -25.95,
					Lon: 32.58,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Validate(tt.data)
			if errs != nil {
				t.Errorf("Validate() returned errors for valid struct: %v", errs)
			}
		})
	}
}

func TestValidate_InvalidStruct(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		expectedFields []string
		expectedCodes  []string
	}{
		{
			name:           "missing required fields",
			data:           UserRegistration{},
			expectedFields: []string{"name", "email", "phone", "password"},
			expectedCodes:  []string{valerrors.CodeRequired, valerrors.CodeRequired, valerrors.CodeRequired, valerrors.CodeRequired},
		},
		{
			name: "invalid phone format",
			data: UserRegistration{
				Name:     "João",
				Email:    "joao@example.com",
				Phone:    "invalid-phone",
				Password: "securepass123",
			},
			expectedFields: []string{"phone"},
			expectedCodes:  []string{valerrors.CodeInvalidFormat},
		},
		{
			name: "invalid email",
			data: UserRegistration{
				Name:     "João",
				Email:    "not-an-email",
				Phone:    "+258841234567",
				Password: "securepass123",
			},
			expectedFields: []string{"email"},
			expectedCodes:  []string{valerrors.CodeInvalidFormat},
		},
		{
			name: "name too short",
			data: UserRegistration{
				Name:     "J",
				Email:    "joao@example.com",
				Phone:    "+258841234567",
				Password: "securepass123",
			},
			expectedFields: []string{"name"},
			expectedCodes:  []string{valerrors.CodeTooShort},
		},
		{
			name: "password too short",
			data: UserRegistration{
				Name:     "João",
				Email:    "joao@example.com",
				Phone:    "+258841234567",
				Password: "short",
			},
			expectedFields: []string{"password"},
			expectedCodes:  []string{valerrors.CodeTooShort},
		},
		{
			name: "invalid plate format",
			data: VehicleInfo{
				Plate: "INVALID",
				Year:  2022,
				Color: "white",
			},
			expectedFields: []string{"plate"},
			expectedCodes:  []string{valerrors.CodeInvalidFormat},
		},
		{
			name: "year out of range",
			data: VehicleInfo{
				Plate: "AAA-123-MP",
				Year:  2005,
				Color: "white",
			},
			expectedFields: []string{"year"},
			expectedCodes:  []string{valerrors.CodeOutOfRange},
		},
		{
			name: "invalid color option",
			data: VehicleInfo{
				Plate: "AAA-123-MP",
				Year:  2022,
				Color: "purple",
			},
			expectedFields: []string{"color"},
			expectedCodes:  []string{valerrors.CodeInvalidOption},
		},
		{
			name: "invalid PIN",
			data: RideRequest{
				PIN:  "1234", // Sequential
				Fare: 10000,
				Pickup: Location{
					Lat: -25.95,
					Lon: 32.58,
				},
			},
			expectedFields: []string{"pin"},
			expectedCodes:  []string{valerrors.CodeInvalidFormat},
		},
		{
			name: "invalid rating",
			data: RideRequest{
				PIN:    "7392",
				Fare:   10000,
				Rating: 10, // Out of 1-5 range
				Pickup: Location{
					Lat: -25.95,
					Lon: 32.58,
				},
			},
			expectedFields: []string{"rating"},
			expectedCodes:  []string{valerrors.CodeOutOfRange},
		},
		{
			name: "location outside Mozambique",
			data: RideRequest{
				PIN:  "7392",
				Fare: 10000,
				Pickup: Location{
					Lat: -34.0, // South Africa - not in Mozambique
					Lon: 18.0,
				},
			},
			expectedFields: []string{"pickup"},
			expectedCodes:  []string{valerrors.CodeOutsideServiceArea},
		},
		{
			name: "negative money value",
			data: RideRequest{
				PIN:  "7392",
				Fare: -100,
				Pickup: Location{
					Lat: -25.95,
					Lon: 32.58,
				},
			},
			expectedFields: []string{"fare"},
			expectedCodes:  []string{valerrors.CodeOutOfRange},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Validate(tt.data)
			if errs == nil {
				t.Fatal("Validate() should return errors for invalid struct")
			}

			for i, expectedField := range tt.expectedFields {
				if !errs.HasField(expectedField) {
					t.Errorf("expected error for field %q", expectedField)
				}
				fieldErrs := errs.GetByField(expectedField)
				if len(fieldErrs) == 0 {
					continue
				}
				if fieldErrs[0].Code != tt.expectedCodes[i] {
					t.Errorf("field %q: expected code %q, got %q", expectedField, tt.expectedCodes[i], fieldErrs[0].Code)
				}
			}
		})
	}
}

func TestValidate_NestedStruct(t *testing.T) {
	t.Run("valid nested struct", func(t *testing.T) {
		data := NestedStruct{
			User: UserRegistration{
				Name:     "João",
				Email:    "joao@example.com",
				Phone:    "+258841234567",
				Password: "securepass123",
			},
			Vehicle: VehicleInfo{
				Plate: "AAA-123-MP",
				Year:  2022,
				Color: "white",
			},
		}

		errs := Validate(data)
		if errs != nil {
			t.Errorf("Validate() returned errors for valid nested struct: %v", errs)
		}
	})

	t.Run("invalid nested user", func(t *testing.T) {
		data := NestedStruct{
			User: UserRegistration{
				Name:     "J", // Too short
				Email:    "invalid",
				Phone:    "bad",
				Password: "short",
			},
			Vehicle: VehicleInfo{
				Plate: "AAA-123-MP",
				Year:  2022,
				Color: "white",
			},
		}

		errs := Validate(data)
		if errs == nil {
			t.Fatal("Validate() should return errors for invalid nested struct")
		}
	})
}

func TestValidate_OptionalFields(t *testing.T) {
	t.Run("all empty is valid", func(t *testing.T) {
		data := OptionalFields{}
		errs := Validate(data)
		if errs != nil {
			t.Errorf("Validate() returned errors for struct with optional empty fields: %v", errs)
		}
	})

	t.Run("valid optional values", func(t *testing.T) {
		data := OptionalFields{
			Name:   "João",
			Phone:  "+258841234567",
			Rating: 5,
		}
		errs := Validate(data)
		if errs != nil {
			t.Errorf("Validate() returned errors for valid optional values: %v", errs)
		}
	})

	t.Run("invalid optional values", func(t *testing.T) {
		data := OptionalFields{
			Name:   "J", // Too short when provided
			Phone:  "bad",
			Rating: 10,
		}
		errs := Validate(data)
		if errs == nil {
			t.Fatal("Validate() should return errors for invalid optional values")
		}

		if !errs.HasField("name") {
			t.Error("expected error for name field")
		}
		if !errs.HasField("phone") {
			t.Error("expected error for phone field")
		}
		if !errs.HasField("rating") {
			t.Error("expected error for rating field")
		}
	})
}

func TestValidateVar(t *testing.T) {
	tests := []struct {
		name    string
		field   interface{}
		tag     string
		wantErr bool
	}{
		{"valid email", "test@example.com", "email", false},
		{"invalid email", "not-an-email", "email", true},
		{"valid phone", "+258841234567", "mz_phone", false},
		{"invalid phone", "invalid", "mz_phone", true},
		{"valid plate", "AAA-123-MP", "mz_plate", false},
		{"invalid plate", "INVALID", "mz_plate", true},
		{"valid pin", "7392", "txova_pin", false},
		{"invalid pin", "1234", "txova_pin", true},
		{"valid rating", 5, "txova_rating", false},
		{"invalid rating", 10, "txova_rating", true},
		{"valid money", int64(100), "txova_money", false},
		{"zero money", int64(0), "txova_money", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateVar(tt.field, tt.tag)
			if tt.wantErr && errs == nil {
				t.Error("ValidateVar() should return error")
			}
			if !tt.wantErr && errs != nil {
				t.Errorf("ValidateVar() unexpected error: %v", errs)
			}
		})
	}
}

func TestRegisterValidation(t *testing.T) {
	t.Run("register custom validation", func(t *testing.T) {
		customFn := func(fl validator.FieldLevel) bool {
			return fl.Field().String() == "custom"
		}

		err := RegisterValidation("test_custom", customFn)
		if err != nil {
			t.Errorf("RegisterValidation() error = %v", err)
		}

		// Test the custom validation
		type TestStruct struct {
			Value string `validate:"test_custom"`
		}

		valid := TestStruct{Value: "custom"}
		errs := Validate(valid)
		if errs != nil {
			t.Errorf("custom validation should pass for valid value: %v", errs)
		}

		invalid := TestStruct{Value: "not-custom"}
		errs = Validate(invalid)
		if errs == nil {
			t.Error("custom validation should fail for invalid value")
		}
	})
}

func TestValidateMzPhone(t *testing.T) {
	type PhoneTest struct {
		Phone string `json:"phone" validate:"required,mz_phone"`
	}

	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid international format", "+258841234567", false},
		{"valid local format", "841234567", false},
		{"valid with spaces", "84 123 4567", false},
		{"valid without plus", "258841234567", false},
		{"invalid prefix", "881234567", true},
		{"too short", "8412345", true},
		{"too long", "8412345678901", true},
		{"non-numeric", "84abcdefg", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := PhoneTest{Phone: tt.phone}
			errs := Validate(data)
			if tt.wantErr && errs == nil {
				t.Error("expected validation error")
			}
			if !tt.wantErr && errs != nil {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestValidateMzPlate(t *testing.T) {
	type PlateTest struct {
		Plate string `json:"plate" validate:"required,mz_plate"`
	}

	tests := []struct {
		name    string
		plate   string
		wantErr bool
	}{
		{"valid standard format", "AAA-123-MP", false},
		{"valid without dashes", "AAA123MP", false},
		{"valid old format", "MP-12-34", false},
		{"valid lowercase", "aaa-123-mp", false},
		{"invalid format", "INVALID", true},
		{"invalid province", "AAA-123-XX", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := PlateTest{Plate: tt.plate}
			errs := Validate(data)
			if tt.wantErr && errs == nil {
				t.Error("expected validation error")
			}
			if !tt.wantErr && errs != nil {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestValidateMzLocation(t *testing.T) {
	type LocationTest struct {
		Location Location `json:"location" validate:"mz_location"`
	}

	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
	}{
		{"Maputo center", -25.95, 32.58, false},
		{"Beira", -19.84, 34.84, false},
		{"northern Mozambique", -12.0, 40.0, false},
		{"outside - South Africa", -34.0, 18.0, true},
		{"outside - equator", 0.0, 0.0, true},
		{"outside - Zimbabwe", -20.0, 29.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := LocationTest{
				Location: Location{Lat: tt.lat, Lon: tt.lon},
			}
			errs := Validate(data)
			if tt.wantErr && errs == nil {
				t.Error("expected validation error")
			}
			if !tt.wantErr && errs != nil {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestValidateTxovaPin(t *testing.T) {
	type PinTest struct {
		PIN string `json:"pin" validate:"required,txova_pin"`
	}

	tests := []struct {
		name    string
		pin     string
		wantErr bool
	}{
		{"valid random pin", "7392", false},
		{"valid pin 4826", "4826", false},
		{"sequential ascending", "1234", true},
		{"sequential descending", "4321", true},
		{"sequential 5678", "5678", true},
		{"all same digits", "1111", true},
		{"too short", "123", true},
		{"too long", "12345", true},
		{"non-numeric", "abcd", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := PinTest{PIN: tt.pin}
			errs := Validate(data)
			if tt.wantErr && errs == nil {
				t.Error("expected validation error")
			}
			if !tt.wantErr && errs != nil {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestValidateTxovaMoney(t *testing.T) {
	type MoneyTest struct {
		Amount int64 `json:"amount" validate:"required,txova_money"`
	}

	tests := []struct {
		name    string
		amount  int64
		wantErr bool
	}{
		{"positive amount", 10000, false},
		{"small positive", 1, false},
		{"large amount", 1000000000, false},
		{"zero", 0, true},
		{"negative", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := MoneyTest{Amount: tt.amount}
			errs := Validate(data)
			if tt.wantErr && errs == nil {
				t.Error("expected validation error")
			}
			if !tt.wantErr && errs != nil {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestValidateTxovaRating(t *testing.T) {
	type RatingTest struct {
		Rating int `json:"rating" validate:"required,txova_rating"`
	}

	tests := []struct {
		name    string
		rating  int
		wantErr bool
	}{
		{"rating 1", 1, false},
		{"rating 3", 3, false},
		{"rating 5", 5, false},
		{"rating 0", 0, true},
		{"rating 6", 6, true},
		{"negative rating", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := RatingTest{Rating: tt.rating}
			errs := Validate(data)
			if tt.wantErr && errs == nil {
				t.Error("expected validation error")
			}
			if !tt.wantErr && errs != nil {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestValidateTxovaVehicleYear(t *testing.T) {
	type YearTest struct {
		Year int `json:"year" validate:"required,txova_vehicle_year"`
	}

	currentYear := time.Now().Year()

	tests := []struct {
		name    string
		year    int
		wantErr bool
	}{
		{"min year 2010", 2010, false},
		{"current year", currentYear, false},
		{"next year", currentYear + 1, false},
		{"year too old", 2009, true},
		{"year too new", currentYear + 2, true},
		{"zero year", 0, true},
		{"negative year", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := YearTest{Year: tt.year}
			errs := Validate(data)
			if tt.wantErr && errs == nil {
				t.Error("expected validation error")
			}
			if !tt.wantErr && errs != nil {
				t.Errorf("unexpected error: %v", errs)
			}
		})
	}
}

func TestFieldNameMapping(t *testing.T) {
	type TestStruct struct {
		UserName string `json:"user_name" validate:"required"`
		Email    string `json:"-" validate:"required"`
		Phone    string `json:"" validate:"required"`
		NoTag    string `validate:"required"`
	}

	data := TestStruct{}
	errs := Validate(data)

	if errs == nil {
		t.Fatal("expected validation errors")
	}

	// Check that JSON tag names are used
	if !errs.HasField("user_name") {
		t.Error("expected error for field 'user_name' (JSON tag)")
	}

	// When json tag is "-", use struct field name
	if !errs.HasField("Email") {
		t.Error("expected error for field 'Email' (struct name when json is '-')")
	}

	// When json tag is empty, use struct field name
	if !errs.HasField("Phone") {
		t.Error("expected error for field 'Phone' (struct name when json is empty)")
	}

	// When no json tag, use struct field name
	if !errs.HasField("NoTag") {
		t.Error("expected error for field 'NoTag' (no json tag)")
	}
}

func TestTranslateError_AllTags(t *testing.T) {
	// Test various validation tags produce correct error codes
	tests := []struct {
		name         string
		data         interface{}
		expectedCode string
	}{
		{
			name: "url validation",
			data: struct {
				URL string `json:"url" validate:"required,url"`
			}{URL: "not-a-url"},
			expectedCode: valerrors.CodeInvalidFormat,
		},
		{
			name: "len validation",
			data: struct {
				Code string `json:"code" validate:"len=4"`
			}{Code: "123"},
			expectedCode: valerrors.CodeInvalidFormat,
		},
		{
			name: "gt validation",
			data: struct {
				Value int `json:"value" validate:"gt=10"`
			}{Value: 5},
			expectedCode: valerrors.CodeOutOfRange,
		},
		{
			name: "lt validation",
			data: struct {
				Value int `json:"value" validate:"lt=10"`
			}{Value: 15},
			expectedCode: valerrors.CodeOutOfRange,
		},
		{
			name: "max string length",
			data: struct {
				Name string `json:"name" validate:"max=5"`
			}{Name: "verylongname"},
			expectedCode: valerrors.CodeTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Validate(tt.data)
			if errs == nil {
				t.Fatal("expected validation error")
			}
			if errs[0].Code != tt.expectedCode {
				t.Errorf("expected code %q, got %q", tt.expectedCode, errs[0].Code)
			}
		})
	}
}

func TestValidateLocationSlice(t *testing.T) {
	type SliceLocationTest struct {
		Coords []float64 `json:"coords" validate:"mz_location"`
	}

	t.Run("valid slice location", func(t *testing.T) {
		data := SliceLocationTest{
			Coords: []float64{-25.95, 32.58}, // Maputo
		}
		errs := Validate(data)
		if errs != nil {
			t.Errorf("unexpected error: %v", errs)
		}
	})

	t.Run("invalid slice location - outside Mozambique", func(t *testing.T) {
		data := SliceLocationTest{
			Coords: []float64{0.0, 0.0},
		}
		errs := Validate(data)
		if errs == nil {
			t.Error("expected validation error")
		}
	})

	t.Run("invalid slice - too short", func(t *testing.T) {
		data := SliceLocationTest{
			Coords: []float64{-25.95},
		}
		errs := Validate(data)
		if errs == nil {
			t.Error("expected validation error")
		}
	})
}

func TestEmptyPhoneAndPlate(t *testing.T) {
	// Test that empty strings for optional mz_phone and mz_plate pass validation
	type OptionalContact struct {
		Phone string `json:"phone" validate:"omitempty,mz_phone"`
		Plate string `json:"plate" validate:"omitempty,mz_plate"`
		PIN   string `json:"pin" validate:"omitempty,txova_pin"`
	}

	data := OptionalContact{}
	errs := Validate(data)
	if errs != nil {
		t.Errorf("empty optional fields should pass validation: %v", errs)
	}
}

func TestMoneyValidationTypes(t *testing.T) {
	t.Run("int type", func(t *testing.T) {
		type IntMoney struct {
			Amount int `json:"amount" validate:"txova_money"`
		}
		errs := Validate(IntMoney{Amount: 100})
		if errs != nil {
			t.Errorf("int money should be valid: %v", errs)
		}
	})

	t.Run("uint type", func(t *testing.T) {
		type UintMoney struct {
			Amount uint `json:"amount" validate:"txova_money"`
		}
		errs := Validate(UintMoney{Amount: 100})
		if errs != nil {
			t.Errorf("uint money should be valid: %v", errs)
		}
	})

	t.Run("float64 type", func(t *testing.T) {
		type FloatMoney struct {
			Amount float64 `json:"amount" validate:"txova_money"`
		}
		errs := Validate(FloatMoney{Amount: 100.50})
		if errs != nil {
			t.Errorf("float money should be valid: %v", errs)
		}
	})

	t.Run("string type fails", func(t *testing.T) {
		type StringMoney struct {
			Amount string `json:"amount" validate:"txova_money"`
		}
		errs := Validate(StringMoney{Amount: "100"})
		if errs == nil {
			t.Error("string money should fail validation")
		}
	})
}

func TestRatingValidationTypes(t *testing.T) {
	t.Run("uint type", func(t *testing.T) {
		type UintRating struct {
			Rating uint `json:"rating" validate:"txova_rating"`
		}
		errs := Validate(UintRating{Rating: 5})
		if errs != nil {
			t.Errorf("uint rating should be valid: %v", errs)
		}
	})

	t.Run("string type fails", func(t *testing.T) {
		type StringRating struct {
			Rating string `json:"rating" validate:"txova_rating"`
		}
		errs := Validate(StringRating{Rating: "5"})
		if errs == nil {
			t.Error("string rating should fail validation")
		}
	})
}

func TestLocationValidationInvalidKind(t *testing.T) {
	type InvalidLocation struct {
		Location string `json:"location" validate:"mz_location"`
	}

	data := InvalidLocation{Location: "-25.95,32.58"}
	errs := Validate(data)
	if errs == nil {
		t.Error("string location should fail mz_location validation")
	}
}
