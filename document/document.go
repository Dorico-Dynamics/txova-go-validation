// Package document provides document and file validation for uploads.
package document

import (
	"strings"

	valerrors "github.com/Dorico-Dynamics/txova-go-validation/errors"
)

// Document types.
const (
	DocTypeDriverLicense       = "driver_license"
	DocTypeVehicleRegistration = "vehicle_registration"
	DocTypeInsurance           = "insurance"
	DocTypeIDCard              = "id_card"
	DocTypeProfilePhoto        = "profile_photo"
	DocTypeVehiclePhoto        = "vehicle_photo"
)

// File size limits in bytes.
const (
	MaxDocumentSize     = 5 * 1024 * 1024 // 5 MB
	MaxProfilePhotoSize = 2 * 1024 * 1024 // 2 MB
)

// Image dimension constraints.
const (
	MinImageWidth  = 200
	MinImageHeight = 200
	MaxImageWidth  = 4096
	MaxImageHeight = 4096
)

// Aspect ratio constraints.
const (
	MinAspectRatio = 0.25 // 1:4
	MaxAspectRatio = 4.0  // 4:1
)

// AllowedFormats maps document types to their allowed file extensions.
var AllowedFormats = map[string][]string{
	DocTypeDriverLicense:       {"jpg", "jpeg", "png", "pdf"},
	DocTypeVehicleRegistration: {"jpg", "jpeg", "png", "pdf"},
	DocTypeInsurance:           {"jpg", "jpeg", "png", "pdf"},
	DocTypeIDCard:              {"jpg", "jpeg", "png", "pdf"},
	DocTypeProfilePhoto:        {"jpg", "jpeg", "png"},
	DocTypeVehiclePhoto:        {"jpg", "jpeg", "png"},
}

// MIMETypes maps extensions to expected MIME types.
var MIMETypes = map[string][]string{
	"jpg":  {"image/jpeg"},
	"jpeg": {"image/jpeg"},
	"png":  {"image/png"},
	"pdf":  {"application/pdf"},
}

// AllDocTypes returns a list of all valid document types.
func AllDocTypes() []string {
	return []string{
		DocTypeDriverLicense,
		DocTypeVehicleRegistration,
		DocTypeInsurance,
		DocTypeIDCard,
		DocTypeProfilePhoto,
		DocTypeVehiclePhoto,
	}
}

// ValidateDocType validates that a document type is valid.
func ValidateDocType(docType string) error {
	for _, dt := range AllDocTypes() {
		if dt == docType {
			return nil
		}
	}
	return valerrors.InvalidOptionWithValue("document_type", AllDocTypes(), docType)
}

// ValidateFileSize validates that a file size is within limits for the document type.
func ValidateFileSize(size int64, docType string) error {
	if err := ValidateDocType(docType); err != nil {
		return err
	}

	maxSize := getMaxSize(docType)
	if size > maxSize {
		return valerrors.OutOfRangeWithValue("file_size", 0, maxSize, size)
	}
	if size <= 0 {
		return valerrors.OutOfRangeWithValue("file_size", 1, maxSize, size)
	}

	return nil
}

// getMaxSize returns the maximum file size for a document type.
func getMaxSize(docType string) int64 {
	if docType == DocTypeProfilePhoto {
		return MaxProfilePhotoSize
	}
	return MaxDocumentSize
}

// ValidateMIMEType validates that a MIME type matches the expected type for the extension.
func ValidateMIMEType(mimeType, extension string) error {
	ext := strings.ToLower(strings.TrimPrefix(extension, "."))
	mime := strings.ToLower(strings.TrimSpace(mimeType))

	expectedTypes, ok := MIMETypes[ext]
	if !ok {
		return valerrors.InvalidFormatWithValue("extension", "jpg, jpeg, png, or pdf", extension)
	}

	for _, expected := range expectedTypes {
		if mime == expected {
			return nil
		}
	}

	return valerrors.InvalidFormatWithValue("mime_type", strings.Join(expectedTypes, " or "), mimeType)
}

// ValidateImageDimensions validates that image dimensions are within acceptable range.
func ValidateImageDimensions(width, height int) error {
	if width < MinImageWidth || width > MaxImageWidth {
		return valerrors.OutOfRangeWithValue("width", MinImageWidth, MaxImageWidth, width)
	}
	if height < MinImageHeight || height > MaxImageHeight {
		return valerrors.OutOfRangeWithValue("height", MinImageHeight, MaxImageHeight, height)
	}
	return nil
}

// ValidateAspectRatio validates that an image aspect ratio is within acceptable range.
func ValidateAspectRatio(width, height int) error {
	if height == 0 {
		return valerrors.InvalidFormat("height", "non-zero value")
	}

	ratio := float64(width) / float64(height)
	if ratio < MinAspectRatio || ratio > MaxAspectRatio {
		return valerrors.OutOfRangeWithValue("aspect_ratio", MinAspectRatio, MaxAspectRatio, ratio)
	}
	return nil
}

// ValidateImage validates all image properties at once.
func ValidateImage(width, height int, size int64, docType string) error {
	if err := ValidateFileSize(size, docType); err != nil {
		return err
	}
	if err := ValidateImageDimensions(width, height); err != nil {
		return err
	}
	if err := ValidateAspectRatio(width, height); err != nil {
		return err
	}
	return nil
}

// GetAllowedFormats returns the allowed file formats for a document type.
func GetAllowedFormats(docType string) []string {
	if formats, ok := AllowedFormats[docType]; ok {
		return formats
	}
	return nil
}

// IsAllowedFormat checks if a file extension is allowed for a document type.
func IsAllowedFormat(extension, docType string) bool {
	ext := strings.ToLower(strings.TrimPrefix(extension, "."))
	formats := GetAllowedFormats(docType)
	for _, f := range formats {
		if f == ext {
			return true
		}
	}
	return false
}

// ValidateFormat validates that a file format is allowed for the document type.
func ValidateFormat(extension, docType string) error {
	if err := ValidateDocType(docType); err != nil {
		return err
	}

	if !IsAllowedFormat(extension, docType) {
		formats := GetAllowedFormats(docType)
		return valerrors.InvalidOptionWithValue("format", formats, extension)
	}
	return nil
}

// IsValidDocType returns true if the document type is valid.
func IsValidDocType(docType string) bool {
	return ValidateDocType(docType) == nil
}

// IsImageType returns true if the document type is an image-only type (no PDF).
func IsImageType(docType string) bool {
	return docType == DocTypeProfilePhoto || docType == DocTypeVehiclePhoto
}

// GetMaxFileSize returns the maximum file size for a document type in bytes.
func GetMaxFileSize(docType string) int64 {
	return getMaxSize(docType)
}
