package errors

import (
	"encoding/json"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  ValidationError
		want string
	}{
		{
			name: "error without value",
			err:  ValidationError{Field: "email", Code: CodeRequired, Message: "email is required"},
			want: "email: email is required",
		},
		{
			name: "error with value",
			err:  ValidationError{Field: "age", Code: CodeOutOfRange, Message: "age must be between 18 and 120", Value: 150},
			want: "age: age must be between 18 and 120 (value: 150)",
		},
		{
			name: "error with string value",
			err:  ValidationError{Field: "phone", Code: CodeInvalidFormat, Message: "phone has invalid format", Value: "abc"},
			want: "phone: phone has invalid format (value: abc)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	err := New("field", CodeRequired, "field is required")
	if err.Field != "field" {
		t.Errorf("Field = %v, want field", err.Field)
	}
	if err.Code != CodeRequired {
		t.Errorf("Code = %v, want %v", err.Code, CodeRequired)
	}
	if err.Message != "field is required" {
		t.Errorf("Message = %v, want field is required", err.Message)
	}
	if err.Value != nil {
		t.Errorf("Value = %v, want nil", err.Value)
	}
}

func TestNewWithValue(t *testing.T) {
	err := NewWithValue("age", CodeOutOfRange, "age out of range", 150)
	if err.Field != "age" {
		t.Errorf("Field = %v, want age", err.Field)
	}
	if err.Code != CodeOutOfRange {
		t.Errorf("Code = %v, want %v", err.Code, CodeOutOfRange)
	}
	if err.Value != 150 {
		t.Errorf("Value = %v, want 150", err.Value)
	}
}

func TestRequired(t *testing.T) {
	err := Required("username")
	if err.Field != "username" {
		t.Errorf("Field = %v, want username", err.Field)
	}
	if err.Code != CodeRequired {
		t.Errorf("Code = %v, want %v", err.Code, CodeRequired)
	}
	if err.Message != "username is required" {
		t.Errorf("Message = %v, want 'username is required'", err.Message)
	}
}

func TestInvalidFormat(t *testing.T) {
	err := InvalidFormat("email", "valid email address")
	if err.Field != "email" {
		t.Errorf("Field = %v, want email", err.Field)
	}
	if err.Code != CodeInvalidFormat {
		t.Errorf("Code = %v, want %v", err.Code, CodeInvalidFormat)
	}
	if err.Message != "email has invalid format, expected valid email address" {
		t.Errorf("Message = %v", err.Message)
	}
}

func TestInvalidFormatWithValue(t *testing.T) {
	err := InvalidFormatWithValue("phone", "+258XXXXXXXXX", "123")
	if err.Field != "phone" {
		t.Errorf("Field = %v, want phone", err.Field)
	}
	if err.Code != CodeInvalidFormat {
		t.Errorf("Code = %v, want %v", err.Code, CodeInvalidFormat)
	}
	if err.Value != "123" {
		t.Errorf("Value = %v, want 123", err.Value)
	}
}

func TestOutOfRange(t *testing.T) {
	err := OutOfRange("rating", 1, 5)
	if err.Field != "rating" {
		t.Errorf("Field = %v, want rating", err.Field)
	}
	if err.Code != CodeOutOfRange {
		t.Errorf("Code = %v, want %v", err.Code, CodeOutOfRange)
	}
	if err.Message != "rating must be between 1 and 5" {
		t.Errorf("Message = %v", err.Message)
	}
}

func TestOutOfRangeWithValue(t *testing.T) {
	err := OutOfRangeWithValue("rating", 1, 5, 10)
	if err.Value != 10 {
		t.Errorf("Value = %v, want 10", err.Value)
	}
}

func TestTooShort(t *testing.T) {
	err := TooShort("password", 8)
	if err.Field != "password" {
		t.Errorf("Field = %v, want password", err.Field)
	}
	if err.Code != CodeTooShort {
		t.Errorf("Code = %v, want %v", err.Code, CodeTooShort)
	}
	if err.Message != "password must be at least 8 characters" {
		t.Errorf("Message = %v", err.Message)
	}
}

func TestTooShortWithValue(t *testing.T) {
	err := TooShortWithValue("password", 8, 5)
	if err.Value != 5 {
		t.Errorf("Value = %v, want 5", err.Value)
	}
}

func TestTooLong(t *testing.T) {
	err := TooLong("username", 20)
	if err.Field != "username" {
		t.Errorf("Field = %v, want username", err.Field)
	}
	if err.Code != CodeTooLong {
		t.Errorf("Code = %v, want %v", err.Code, CodeTooLong)
	}
	if err.Message != "username must be at most 20 characters" {
		t.Errorf("Message = %v", err.Message)
	}
}

func TestTooLongWithValue(t *testing.T) {
	err := TooLongWithValue("username", 20, 25)
	if err.Value != 25 {
		t.Errorf("Value = %v, want 25", err.Value)
	}
}

func TestInvalidOption(t *testing.T) {
	options := []string{"active", "pending", "suspended"}
	err := InvalidOption("status", options)
	if err.Field != "status" {
		t.Errorf("Field = %v, want status", err.Field)
	}
	if err.Code != CodeInvalidOption {
		t.Errorf("Code = %v, want %v", err.Code, CodeInvalidOption)
	}
	if err.Message != "status must be one of: active, pending, suspended" {
		t.Errorf("Message = %v", err.Message)
	}
}

func TestInvalidOptionWithValue(t *testing.T) {
	options := []string{"standard", "premium"}
	err := InvalidOptionWithValue("service_type", options, "invalid")
	if err.Value != "invalid" {
		t.Errorf("Value = %v, want invalid", err.Value)
	}
}

func TestOutsideServiceArea(t *testing.T) {
	err := OutsideServiceArea("pickup_location")
	if err.Field != "pickup_location" {
		t.Errorf("Field = %v, want pickup_location", err.Field)
	}
	if err.Code != CodeOutsideServiceArea {
		t.Errorf("Code = %v, want %v", err.Code, CodeOutsideServiceArea)
	}
	if err.Message != "pickup_location is outside the service area" {
		t.Errorf("Message = %v", err.Message)
	}
}

func TestOutsideServiceAreaWithValue(t *testing.T) {
	err := OutsideServiceAreaWithValue("location", -25.969, 32.573)
	if err.Value != "-25.969000, 32.573000" {
		t.Errorf("Value = %v, want '-25.969000, 32.573000'", err.Value)
	}
}

func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name   string
		errors ValidationErrors
		want   string
	}{
		{
			name:   "empty errors",
			errors: ValidationErrors{},
			want:   "no validation errors",
		},
		{
			name: "single error",
			errors: ValidationErrors{
				{Field: "email", Code: CodeRequired, Message: "email is required"},
			},
			want: "email: email is required",
		},
		{
			name: "multiple errors",
			errors: ValidationErrors{
				{Field: "email", Code: CodeRequired, Message: "email is required"},
				{Field: "password", Code: CodeTooShort, Message: "password must be at least 8 characters"},
			},
			want: "2 validation errors: email: email is required; password: password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errors.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationErrors_HasErrors(t *testing.T) {
	tests := []struct {
		name   string
		errors ValidationErrors
		want   bool
	}{
		{"empty", ValidationErrors{}, false},
		{"with error", ValidationErrors{{Field: "test"}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errors.HasErrors(); got != tt.want {
				t.Errorf("HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationErrors_HasField(t *testing.T) {
	errors := ValidationErrors{
		{Field: "email", Code: CodeRequired},
		{Field: "password", Code: CodeTooShort},
	}

	tests := []struct {
		name  string
		field string
		want  bool
	}{
		{"existing field", "email", true},
		{"another existing field", "password", true},
		{"non-existing field", "username", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.HasField(tt.field); got != tt.want {
				t.Errorf("HasField(%q) = %v, want %v", tt.field, got, tt.want)
			}
		})
	}
}

func TestValidationErrors_GetByField(t *testing.T) {
	errors := ValidationErrors{
		{Field: "email", Code: CodeRequired},
		{Field: "email", Code: CodeInvalidFormat},
		{Field: "password", Code: CodeTooShort},
	}

	t.Run("multiple errors for field", func(t *testing.T) {
		result := errors.GetByField("email")
		if len(result) != 2 {
			t.Errorf("GetByField('email') returned %d errors, want 2", len(result))
		}
	})

	t.Run("single error for field", func(t *testing.T) {
		result := errors.GetByField("password")
		if len(result) != 1 {
			t.Errorf("GetByField('password') returned %d errors, want 1", len(result))
		}
	})

	t.Run("no errors for field", func(t *testing.T) {
		result := errors.GetByField("username")
		if len(result) != 0 {
			t.Errorf("GetByField('username') returned %d errors, want 0", len(result))
		}
	})
}

func TestValidationErrors_GetByCode(t *testing.T) {
	errors := ValidationErrors{
		{Field: "email", Code: CodeRequired},
		{Field: "username", Code: CodeRequired},
		{Field: "password", Code: CodeTooShort},
	}

	t.Run("multiple errors with code", func(t *testing.T) {
		result := errors.GetByCode(CodeRequired)
		if len(result) != 2 {
			t.Errorf("GetByCode(CodeRequired) returned %d errors, want 2", len(result))
		}
	})

	t.Run("single error with code", func(t *testing.T) {
		result := errors.GetByCode(CodeTooShort)
		if len(result) != 1 {
			t.Errorf("GetByCode(CodeTooShort) returned %d errors, want 1", len(result))
		}
	})

	t.Run("no errors with code", func(t *testing.T) {
		result := errors.GetByCode(CodeOutOfRange)
		if len(result) != 0 {
			t.Errorf("GetByCode(CodeOutOfRange) returned %d errors, want 0", len(result))
		}
	})
}

func TestValidationErrors_First(t *testing.T) {
	t.Run("empty errors", func(t *testing.T) {
		errors := ValidationErrors{}
		if got := errors.First(); got != nil {
			t.Errorf("First() = %v, want nil", got)
		}
	})

	t.Run("with errors", func(t *testing.T) {
		errors := ValidationErrors{
			{Field: "first", Code: CodeRequired},
			{Field: "second", Code: CodeRequired},
		}
		first := errors.First()
		if first == nil {
			t.Fatal("First() = nil, want non-nil")
		}
		if first.Field != "first" {
			t.Errorf("First().Field = %v, want 'first'", first.Field)
		}
	})
}

func TestValidationErrors_Fields(t *testing.T) {
	errors := ValidationErrors{
		{Field: "email", Code: CodeRequired},
		{Field: "email", Code: CodeInvalidFormat},
		{Field: "password", Code: CodeTooShort},
		{Field: "username", Code: CodeRequired},
	}

	fields := errors.Fields()
	if len(fields) != 3 {
		t.Errorf("Fields() returned %d fields, want 3", len(fields))
	}

	expected := map[string]bool{"email": true, "password": true, "username": true}
	for _, f := range fields {
		if !expected[f] {
			t.Errorf("Unexpected field: %s", f)
		}
	}
}

func TestValidationErrors_Add(t *testing.T) {
	var errors ValidationErrors
	errors.Add(Required("email"))
	errors.Add(TooShort("password", 8))

	if len(errors) != 2 {
		t.Errorf("len(errors) = %d, want 2", len(errors))
	}
	if errors[0].Field != "email" {
		t.Errorf("errors[0].Field = %v, want email", errors[0].Field)
	}
	if errors[1].Field != "password" {
		t.Errorf("errors[1].Field = %v, want password", errors[1].Field)
	}
}

func TestValidationErrors_AddAll(t *testing.T) {
	var errors ValidationErrors
	errors.Add(Required("email"))

	moreErrors := ValidationErrors{
		TooShort("password", 8),
		InvalidFormat("phone", "+258XXXXXXXXX"),
	}
	errors.AddAll(moreErrors)

	if len(errors) != 3 {
		t.Errorf("len(errors) = %d, want 3", len(errors))
	}
}

func TestValidationErrors_MarshalJSON(t *testing.T) {
	t.Run("empty errors", func(t *testing.T) {
		var errors ValidationErrors
		data, err := json.Marshal(errors)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(data) != "[]" {
			t.Errorf("MarshalJSON() = %s, want []", string(data))
		}
	})

	t.Run("nil errors", func(t *testing.T) {
		var errors ValidationErrors = nil
		data, err := json.Marshal(errors)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(data) != "[]" {
			t.Errorf("MarshalJSON() = %s, want []", string(data))
		}
	})

	t.Run("with errors", func(t *testing.T) {
		errors := ValidationErrors{
			{Field: "email", Code: CodeRequired, Message: "email is required"},
		}
		data, err := json.Marshal(errors)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}

		var result []ValidationError
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if len(result) != 1 {
			t.Errorf("Unmarshaled %d errors, want 1", len(result))
		}
		if result[0].Field != "email" {
			t.Errorf("result[0].Field = %v, want email", result[0].Field)
		}
	})

	t.Run("with value field", func(t *testing.T) {
		errors := ValidationErrors{
			{Field: "age", Code: CodeOutOfRange, Message: "age out of range", Value: 150},
		}
		data, err := json.Marshal(errors)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}

		var result []map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if result[0]["value"] != float64(150) {
			t.Errorf("value = %v, want 150", result[0]["value"])
		}
	})

	t.Run("without value field omits it", func(t *testing.T) {
		errors := ValidationErrors{
			{Field: "email", Code: CodeRequired, Message: "email is required"},
		}
		data, err := json.Marshal(errors)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}

		var result []map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if _, exists := result[0]["value"]; exists {
			t.Error("value field should be omitted when empty")
		}
	})
}

func TestValidationErrors_ToError(t *testing.T) {
	t.Run("empty returns nil", func(t *testing.T) {
		errors := ValidationErrors{}
		if got := errors.ToError(); got != nil {
			t.Errorf("ToError() = %v, want nil", got)
		}
	})

	t.Run("with errors returns error", func(t *testing.T) {
		errors := ValidationErrors{Required("email")}
		err := errors.ToError()
		if err == nil {
			t.Error("ToError() = nil, want non-nil")
		}
		ve, ok := err.(ValidationErrors)
		if !ok {
			t.Error("ToError() should return ValidationErrors type")
		}
		if len(ve) != 1 {
			t.Errorf("len(ve) = %d, want 1", len(ve))
		}
	})
}

func TestErrorCodes(t *testing.T) {
	// Verify all error codes are defined correctly
	codes := []string{
		CodeRequired,
		CodeInvalidFormat,
		CodeOutOfRange,
		CodeTooShort,
		CodeTooLong,
		CodeInvalidOption,
		CodeOutsideServiceArea,
	}

	expected := []string{
		"REQUIRED",
		"INVALID_FORMAT",
		"OUT_OF_RANGE",
		"TOO_SHORT",
		"TOO_LONG",
		"INVALID_OPTION",
		"OUTSIDE_SERVICE_AREA",
	}

	for i, code := range codes {
		if code != expected[i] {
			t.Errorf("Code %d = %v, want %v", i, code, expected[i])
		}
	}
}
