package document

import (
	"testing"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

func TestValidateDocType(t *testing.T) {
	tests := []struct {
		name    string
		docType string
		wantErr bool
	}{
		// Valid types
		{"driver_license", DocTypeDriverLicense, false},
		{"vehicle_registration", DocTypeVehicleRegistration, false},
		{"insurance", DocTypeInsurance, false},
		{"id_card", DocTypeIDCard, false},
		{"profile_photo", DocTypeProfilePhoto, false},
		{"vehicle_photo", DocTypeVehiclePhoto, false},

		// Invalid types
		{"empty", "", true},
		{"invalid", "invalid_type", true},
		{"uppercase", "DRIVER_LICENSE", true},
		{"partial", "driver", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDocType(tt.docType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDocType(%q) error = %v, wantErr %v", tt.docType, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	tests := []struct {
		name    string
		size    int64
		docType string
		wantErr bool
	}{
		// Valid sizes for regular documents (5MB max)
		{"1KB driver license", 1024, DocTypeDriverLicense, false},
		{"1MB insurance", 1024 * 1024, DocTypeInsurance, false},
		{"5MB exactly", MaxDocumentSize, DocTypeIDCard, false},

		// Valid sizes for profile photo (2MB max)
		{"1MB profile photo", 1024 * 1024, DocTypeProfilePhoto, false},
		{"2MB profile photo exactly", MaxProfilePhotoSize, DocTypeProfilePhoto, false},

		// Invalid - too large
		{"5MB+1 driver license", MaxDocumentSize + 1, DocTypeDriverLicense, true},
		{"2MB+1 profile photo", MaxProfilePhotoSize + 1, DocTypeProfilePhoto, true},
		{"10MB document", 10 * 1024 * 1024, DocTypeVehicleRegistration, true},

		// Invalid - zero or negative
		{"zero size", 0, DocTypeDriverLicense, true},
		{"negative size", -1, DocTypeDriverLicense, true},

		// Invalid document type
		{"invalid doc type", 1024, "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileSize(tt.size, tt.docType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileSize(%d, %q) error = %v, wantErr %v", tt.size, tt.docType, err, tt.wantErr)
			}
		})
	}
}

func TestValidateMIMEType(t *testing.T) {
	tests := []struct {
		name      string
		mimeType  string
		extension string
		wantErr   bool
	}{
		// Valid combinations
		{"jpeg with jpg", "image/jpeg", "jpg", false},
		{"jpeg with jpeg", "image/jpeg", "jpeg", false},
		{"png", "image/png", "png", false},
		{"pdf", "application/pdf", "pdf", false},
		{"jpeg with .jpg", "image/jpeg", ".jpg", false},
		{"uppercase ext", "image/png", ".PNG", false},
		{"mime with spaces", " image/jpeg ", "jpg", false},

		// Invalid combinations
		{"png mime with jpg ext", "image/png", "jpg", true},
		{"jpeg mime with pdf ext", "image/jpeg", "pdf", true},
		{"invalid mime", "application/octet-stream", "jpg", true},
		{"invalid extension", "image/jpeg", "gif", true},
		{"empty mime", "", "jpg", true},
		{"empty ext", "image/jpeg", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMIMEType(tt.mimeType, tt.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMIMEType(%q, %q) error = %v, wantErr %v", tt.mimeType, tt.extension, err, tt.wantErr)
			}
		})
	}
}

func TestValidateImageDimensions(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		wantErr bool
		errCode string
	}{
		// Valid dimensions
		{"minimum", MinImageWidth, MinImageHeight, false, ""},
		{"maximum", MaxImageWidth, MaxImageHeight, false, ""},
		{"square 500", 500, 500, false, ""},
		{"landscape", 1920, 1080, false, ""},
		{"portrait", 1080, 1920, false, ""},

		// Invalid - too small
		{"width too small", MinImageWidth - 1, MinImageHeight, true, valerrors.CodeOutOfRange},
		{"height too small", MinImageWidth, MinImageHeight - 1, true, valerrors.CodeOutOfRange},
		{"both too small", 100, 100, true, valerrors.CodeOutOfRange},
		{"zero width", 0, 500, true, valerrors.CodeOutOfRange},
		{"zero height", 500, 0, true, valerrors.CodeOutOfRange},

		// Invalid - too large
		{"width too large", MaxImageWidth + 1, MaxImageHeight, true, valerrors.CodeOutOfRange},
		{"height too large", MaxImageWidth, MaxImageHeight + 1, true, valerrors.CodeOutOfRange},
		{"both too large", 5000, 5000, true, valerrors.CodeOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateImageDimensions(tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateImageDimensions(%d, %d) error = %v, wantErr %v", tt.width, tt.height, err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if ve, ok := err.(valerrors.ValidationError); ok {
					if ve.Code != tt.errCode {
						t.Errorf("error code = %v, want %v", ve.Code, tt.errCode)
					}
				}
			}
		})
	}
}

func TestValidateAspectRatio(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		wantErr bool
	}{
		// Valid ratios
		{"square 1:1", 500, 500, false},
		{"landscape 4:3", 800, 600, false},
		{"portrait 3:4", 600, 800, false},
		{"wide 16:9", 1600, 900, false},
		{"max ratio 4:1", 800, 200, false},
		{"min ratio 1:4", 200, 800, false},

		// Invalid ratios
		{"too wide 5:1", 1000, 200, true},
		{"too tall 1:5", 200, 1000, true},
		{"extreme 10:1", 1000, 100, true},

		// Edge cases
		{"zero height", 500, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAspectRatio(tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAspectRatio(%d, %d) error = %v, wantErr %v", tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}

func TestValidateImage(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		size    int64
		docType string
		wantErr bool
	}{
		// Valid
		{"valid profile photo", 500, 500, 1024 * 1024, DocTypeProfilePhoto, false},
		{"valid vehicle photo", 1920, 1080, 3 * 1024 * 1024, DocTypeVehiclePhoto, false},

		// Invalid size
		{"too large for profile", 500, 500, 3 * 1024 * 1024, DocTypeProfilePhoto, true},

		// Invalid dimensions
		{"too small", 100, 100, 1024, DocTypeProfilePhoto, true},
		{"too large", 5000, 5000, 1024, DocTypeProfilePhoto, true},

		// Invalid ratio
		{"too wide", 1000, 100, 1024, DocTypeProfilePhoto, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateImage(tt.width, tt.height, tt.size, tt.docType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAllowedFormats(t *testing.T) {
	tests := []struct {
		name           string
		docType        string
		wantFormats    []string
		wantIncludePDF bool
	}{
		{"driver_license", DocTypeDriverLicense, []string{"jpg", "jpeg", "png", "pdf"}, true},
		{"profile_photo", DocTypeProfilePhoto, []string{"jpg", "jpeg", "png"}, false},
		{"vehicle_photo", DocTypeVehiclePhoto, []string{"jpg", "jpeg", "png"}, false},
		{"invalid", "invalid", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formats := GetAllowedFormats(tt.docType)
			if tt.wantFormats == nil {
				if formats != nil {
					t.Errorf("GetAllowedFormats(%q) = %v, want nil", tt.docType, formats)
				}
				return
			}

			if len(formats) != len(tt.wantFormats) {
				t.Errorf("GetAllowedFormats(%q) len = %d, want %d", tt.docType, len(formats), len(tt.wantFormats))
			}

			hasPDF := false
			for _, f := range formats {
				if f == "pdf" {
					hasPDF = true
				}
			}
			if hasPDF != tt.wantIncludePDF {
				t.Errorf("GetAllowedFormats(%q) has PDF = %v, want %v", tt.docType, hasPDF, tt.wantIncludePDF)
			}
		})
	}
}

func TestIsAllowedFormat(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		docType   string
		want      bool
	}{
		{"jpg for driver license", "jpg", DocTypeDriverLicense, true},
		{"pdf for driver license", "pdf", DocTypeDriverLicense, true},
		{"png for profile photo", "png", DocTypeProfilePhoto, true},
		{"pdf for profile photo", "pdf", DocTypeProfilePhoto, false},
		{"gif for any", "gif", DocTypeDriverLicense, false},
		{"with dot", ".jpg", DocTypeDriverLicense, true},
		{"uppercase", "JPG", DocTypeDriverLicense, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllowedFormat(tt.extension, tt.docType); got != tt.want {
				t.Errorf("IsAllowedFormat(%q, %q) = %v, want %v", tt.extension, tt.docType, got, tt.want)
			}
		})
	}
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		docType   string
		wantErr   bool
	}{
		{"valid jpg", "jpg", DocTypeDriverLicense, false},
		{"valid pdf", "pdf", DocTypeIDCard, false},
		{"invalid pdf for photo", "pdf", DocTypeProfilePhoto, true},
		{"invalid gif", "gif", DocTypeDriverLicense, true},
		{"invalid doc type", "jpg", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormat(tt.extension, tt.docType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFormat(%q, %q) error = %v, wantErr %v", tt.extension, tt.docType, err, tt.wantErr)
			}
		})
	}
}

func TestIsValidDocType(t *testing.T) {
	if !IsValidDocType(DocTypeDriverLicense) {
		t.Error("IsValidDocType(DocTypeDriverLicense) = false, want true")
	}
	if IsValidDocType("invalid") {
		t.Error("IsValidDocType('invalid') = true, want false")
	}
}

func TestIsImageType(t *testing.T) {
	tests := []struct {
		name    string
		docType string
		want    bool
	}{
		{"profile_photo", DocTypeProfilePhoto, true},
		{"vehicle_photo", DocTypeVehiclePhoto, true},
		{"driver_license", DocTypeDriverLicense, false},
		{"insurance", DocTypeInsurance, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsImageType(tt.docType); got != tt.want {
				t.Errorf("IsImageType(%q) = %v, want %v", tt.docType, got, tt.want)
			}
		})
	}
}

func TestGetMaxFileSize(t *testing.T) {
	tests := []struct {
		name    string
		docType string
		want    int64
	}{
		{"profile_photo", DocTypeProfilePhoto, MaxProfilePhotoSize},
		{"driver_license", DocTypeDriverLicense, MaxDocumentSize},
		{"vehicle_registration", DocTypeVehicleRegistration, MaxDocumentSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMaxFileSize(tt.docType); got != tt.want {
				t.Errorf("GetMaxFileSize(%q) = %d, want %d", tt.docType, got, tt.want)
			}
		})
	}
}

func TestAllDocTypes(t *testing.T) {
	types := AllDocTypes()
	if len(types) != 6 {
		t.Errorf("AllDocTypes() len = %d, want 6", len(types))
	}

	expected := map[string]bool{
		DocTypeDriverLicense:       true,
		DocTypeVehicleRegistration: true,
		DocTypeInsurance:           true,
		DocTypeIDCard:              true,
		DocTypeProfilePhoto:        true,
		DocTypeVehiclePhoto:        true,
	}

	for _, dt := range types {
		if !expected[dt] {
			t.Errorf("Unexpected doc type: %s", dt)
		}
	}
}

func TestConstants(t *testing.T) {
	// Verify constants match PRD
	if MaxDocumentSize != 5*1024*1024 {
		t.Errorf("MaxDocumentSize = %d, want 5MB", MaxDocumentSize)
	}
	if MaxProfilePhotoSize != 2*1024*1024 {
		t.Errorf("MaxProfilePhotoSize = %d, want 2MB", MaxProfilePhotoSize)
	}
	if MinImageWidth != 200 {
		t.Errorf("MinImageWidth = %d, want 200", MinImageWidth)
	}
	if MinImageHeight != 200 {
		t.Errorf("MinImageHeight = %d, want 200", MinImageHeight)
	}
	if MaxImageWidth != 4096 {
		t.Errorf("MaxImageWidth = %d, want 4096", MaxImageWidth)
	}
	if MaxImageHeight != 4096 {
		t.Errorf("MaxImageHeight = %d, want 4096", MaxImageHeight)
	}
}
