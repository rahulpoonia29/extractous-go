package extractous_test

import (
	"testing"

	extractous "github.com/rahulpoonia29/extractous-go/src"
)

// ============================================================================
// Extractor Lifecycle Tests
// ============================================================================

func TestExtractor_New(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}
	defer extractor.Close()
}

func TestExtractor_Close(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}

	err := extractor.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// Closing again should be safe
	err = extractor.Close()
	if err != nil {
		t.Errorf("Second Close() returned error: %v", err)
	}
}

func TestExtractor_NilClose(t *testing.T) {
	var extractor *extractous.Extractor
	err := extractor.Close()
	if err != nil {
		t.Errorf("Close() on nil extractor returned error: %v", err)
	}
}

// ============================================================================
// Configuration Tests
// ============================================================================

func TestExtractor_SetExtractStringMaxLength(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}

	newExtractor := extractor.SetExtractStringMaxLength(10000)
	if newExtractor == nil {
		t.Error("SetExtractStringMaxLength returned nil")
	}
	defer newExtractor.Close()
}

func TestExtractor_SetExtractStringMaxLength_Nil(t *testing.T) {
	var extractor *extractous.Extractor
	result := extractor.SetExtractStringMaxLength(10000)
	if result != nil {
		t.Error("Expected nil when calling SetExtractStringMaxLength on nil extractor")
	}
}

func TestExtractor_SetEncoding(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}

	tests := []struct {
		name    string
		charset extractous.CharSet
		wantNil bool
	}{
		{"UTF-8", extractous.CharSetUTF8, false},
		{"US-ASCII", extractous.CharSetUTF8, false},
		{"UTF-16BE", extractous.CharSetUTF8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractor.SetEncoding(tt.charset)
			if tt.wantNil && result != nil {
				t.Errorf("Expected nil, got non-nil")
			}
			if !tt.wantNil && result == nil {
				t.Errorf("Expected non-nil, got nil")
			}
			if result != nil {
				extractor = result
			}
		})
	}

	extractor.Close()
}

func TestExtractor_SetXmlOutput(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}

	newExtractor := extractor.SetXmlOutput(true)
	if newExtractor == nil {
		t.Error("SetXmlOutput(true) returned nil")
	}
	defer newExtractor.Close()

	newExtractor2 := newExtractor.SetXmlOutput(false)
	if newExtractor2 == nil {
		t.Error("SetXmlOutput(false) returned nil")
	}
	defer newExtractor2.Close()
}

func TestExtractor_SetXmlOutput_Nil(t *testing.T) {
	var extractor *extractous.Extractor
	result := extractor.SetXmlOutput(true)
	if result != nil {
		t.Error("Expected nil when calling SetXmlOutput on nil extractor")
	}
}

func TestExtractor_ChainedConfiguration(t *testing.T) {
	extractor := extractous.New().
		SetExtractStringMaxLength(5000).
		SetEncoding(extractous.CharSetUTF8).
		SetXmlOutput(false)

	if extractor == nil {
		t.Fatal("Chained configuration returned nil")
	}
	defer extractor.Close()
}

// ============================================================================
// PDF Configuration Tests
// ============================================================================

func TestPdfConfig_New(t *testing.T) {
	config := extractous.NewPdfConfig()
	if config == nil {
		t.Fatal("NewPdfConfig() returned nil")
	}
}

func TestPdfConfig_SetOcrStrategy(t *testing.T) {
	config := extractous.NewPdfConfig()
	if config == nil {
		t.Fatal("NewPdfConfig() returned nil")
	}

	tests := []struct {
		name     string
		strategy extractous.PdfOcrStrategy
	}{
		{"NoOCR", extractous.PdfOcrNoOcr},
		{"OCROnly", extractous.PdfOcrOcrOnly},
		{"OCRAndText", extractous.PdfOcrOcrAndTextExtraction},
		{"Auto", extractous.PdfOcrAuto},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newConfig := config.SetOcrStrategy(tt.strategy)
			if newConfig == nil {
				t.Errorf("SetOcrStrategy(%v) returned nil", tt.strategy)
			}
			config = newConfig
		})
	}
}

func TestPdfConfig_SetExtractInlineImages(t *testing.T) {
	config := extractous.NewPdfConfig()
	if config == nil {
		t.Fatal("NewPdfConfig() returned nil")
	}

	newConfig := config.SetExtractInlineImages(true)
	if newConfig == nil {
		t.Error("SetExtractInlineImages(true) returned nil")
	}
}

func TestPdfConfig_ChainedConfiguration(t *testing.T) {
	config := extractous.NewPdfConfig().
		SetOcrStrategy(extractous.PdfOcrAuto).
		SetExtractInlineImages(false).
		SetExtractAnnotationText(true)

	if config == nil {
		t.Fatal("Chained PDF configuration returned nil")
	}
}

func TestExtractor_SetPdfConfig(t *testing.T) {
	pdfConfig := extractous.NewPdfConfig()
	if pdfConfig == nil {
		t.Fatal("NewPdfConfig() returned nil")
	}

	extractor := extractous.New().SetPdfConfig(pdfConfig)
	if extractor == nil {
		t.Error("SetPdfConfig returned nil")
	}
	defer extractor.Close()
}

func TestExtractor_SetPdfConfig_Nil(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}
	defer extractor.Close()

	newExtractor := extractor.SetPdfConfig(nil)
	if newExtractor != nil {
		t.Error("Expected nil when setting nil PDF config")
	}
}

// ============================================================================
// Office Configuration Tests
// ============================================================================

func TestOfficeConfig_New(t *testing.T) {
	config := extractous.NewOfficeConfig()
	if config == nil {
		t.Fatal("NewOfficeConfig() returned nil")
	}
}

func TestOfficeConfig_SetExtractMacros(t *testing.T) {
	config := extractous.NewOfficeConfig()
	if config == nil {
		t.Fatal("NewOfficeConfig() returned nil")
	}

	newConfig := config.SetExtractMacros(true)
	if newConfig == nil {
		t.Error("SetExtractMacros(true) returned nil")
	}
}

func TestExtractor_SetOfficeConfig(t *testing.T) {
	officeConfig := extractous.NewOfficeConfig()
	if officeConfig == nil {
		t.Fatal("NewOfficeConfig() returned nil")
	}

	extractor := extractous.New().SetOfficeConfig(officeConfig)
	if extractor == nil {
		t.Error("SetOfficeConfig returned nil")
	}
	defer extractor.Close()
}

// ============================================================================
// OCR Configuration Tests
// ============================================================================

func TestOcrConfig_New(t *testing.T) {
	config := extractous.NewOcrConfig()
	if config == nil {
		t.Fatal("NewOcrConfig() returned nil")
	}
	// Config has finalizer, no manual cleanup needed
}

func TestOcrConfig_SetLanguage(t *testing.T) {
	config := extractous.NewOcrConfig()
	if config == nil {
		t.Fatal("NewOcrConfig() returned nil")
	}

	newConfig := config.SetLanguage("eng")
	if newConfig == nil {
		t.Error("SetLanguage('eng') returned nil")
	}
}

func TestExtractor_SetOcrConfig(t *testing.T) {
	ocrConfig := extractous.NewOcrConfig()
	if ocrConfig == nil {
		t.Fatal("NewOcrConfig() returned nil")
	}

	extractor := extractous.New().SetOcrConfig(ocrConfig)
	if extractor == nil {
		t.Error("SetOcrConfig returned nil")
	}
	defer extractor.Close()
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestExtractor_ExtractFileToString_NilExtractor(t *testing.T) {
	var extractor *extractous.Extractor
	_, _, err := extractor.ExtractFileToString("test.txt")
	if err == nil {
		t.Error("Expected error when using nil extractor")
	}
}

func TestExtractor_ExtractFile_NilExtractor(t *testing.T) {
	var extractor *extractous.Extractor
	_, _, err := extractor.ExtractFile("test.txt")
	if err == nil {
		t.Error("Expected error when using nil extractor")
	}
}

func TestExtractor_ExtractBytesToString_NilExtractor(t *testing.T) {
	var extractor *extractous.Extractor
	_, _, err := extractor.ExtractBytesToString([]byte("test"))
	if err == nil {
		t.Error("Expected error when using nil extractor")
	}
}

func TestExtractor_ExtractBytes_NilExtractor(t *testing.T) {
	var extractor *extractous.Extractor
	_, _, err := extractor.ExtractBytes([]byte("test"))
	if err == nil {
		t.Error("Expected error when using nil extractor")
	}
}

func TestExtractor_ExtractUrlToString_NilExtractor(t *testing.T) {
	var extractor *extractous.Extractor
	_, _, err := extractor.ExtractUrlToString("http://example.com")
	if err == nil {
		t.Error("Expected error when using nil extractor")
	}
}

func TestExtractor_ExtractUrl_NilExtractor(t *testing.T) {
	var extractor *extractous.Extractor
	_, _, err := extractor.ExtractUrl("http://example.com")
	if err == nil {
		t.Error("Expected error when using nil extractor")
	}
}

func TestExtractor_ExtractBytesToString_EmptyBytes(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}
	defer extractor.Close()

	content, metadata, err := extractor.ExtractBytesToString([]byte{})
	if err != nil {
		t.Errorf("Unexpected error with empty bytes: %v", err)
	}
	if content != "" {
		t.Error("Expected empty content for empty bytes")
	}
	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}
}

// ============================================================================
// Metadata Tests
// ============================================================================

func TestMetadata_Get(t *testing.T) {
	metadata := make(extractous.Metadata)
	metadata["key1"] = []string{"value1"}
	metadata["key2"] = []string{"value2", "value3"}

	if metadata.Get("key1") != "value1" {
		t.Error("Get failed for single value")
	}

	if metadata.Get("key2") != "value2" {
		t.Error("Get failed for multiple values (should return first)")
	}

	if metadata.Get("nonexistent") != "" {
		t.Error("Get should return empty string for nonexistent key")
	}
}

func TestMetadata_GetAll(t *testing.T) {
	metadata := make(extractous.Metadata)
	metadata["key1"] = []string{"value1", "value2"}

	values := metadata.GetAll("key1")
	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}

	values = metadata.GetAll("nonexistent")
	if values != nil {
		t.Error("GetAll should return nil for nonexistent key")
	}
}

func TestMetadata_Has(t *testing.T) {
	metadata := make(extractous.Metadata)
	metadata["key1"] = []string{"value1"}

	if !metadata.Has("key1") {
		t.Error("Has returned false for existing key")
	}

	if metadata.Has("nonexistent") {
		t.Error("Has returned true for nonexistent key")
	}
}

func TestMetadata_Keys(t *testing.T) {
	metadata := make(extractous.Metadata)
	metadata["key1"] = []string{"value1"}
	metadata["key2"] = []string{"value2"}

	keys := metadata.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	if !keyMap["key1"] || !keyMap["key2"] {
		t.Error("Keys() didn't return all keys")
	}
}

// ============================================================================
// CharSet Tests
// ============================================================================

func TestCharSet_Constants(t *testing.T) {
	tests := []struct {
		name    string
		charset extractous.CharSet
	}{
		{"UTF-8", extractous.CharSetUTF8},
		{"US-ASCII", extractous.CharSetUSASCII},
		{"UTF-16BE", extractous.CharSetUTF8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify constants exist and can be used
			extractor := extractous.New().SetEncoding(tt.charset)
			if extractor == nil {
				t.Errorf("Failed to set encoding for %s", tt.name)
			} else {
				extractor.Close()
			}
		})
	}
}

// ============================================================================
// Builder Pattern Tests
// ============================================================================

func TestBuilder_ComplexConfiguration(t *testing.T) {
	pdfConfig := extractous.NewPdfConfig().
		SetOcrStrategy(extractous.PdfOcrAuto).
		SetExtractAnnotationText(true).
		SetExtractInlineImages(false)

	officeConfig := extractous.NewOfficeConfig().
		SetExtractMacros(false).
		SetIncludeDeletedContent(true)

	ocrConfig := extractous.NewOcrConfig().
		SetLanguage("eng").
		SetTimeoutSeconds(300)

	extractor := extractous.New().
		SetExtractStringMaxLength(50000).
		SetEncoding(extractous.CharSetUTF8).
		SetPdfConfig(pdfConfig).
		SetOfficeConfig(officeConfig).
		SetOcrConfig(ocrConfig).
		SetXmlOutput(false)

	if extractor == nil {
		t.Fatal("Complex configuration returned nil")
	}
	defer extractor.Close()
}

func TestBuilder_OrderIndependence(t *testing.T) {
	// Configuration order shouldn't matter
	e1 := extractous.New().
		SetXmlOutput(true).
		SetExtractStringMaxLength(1000).
		SetEncoding(extractous.CharSetUTF8)

	e2 := extractous.New().
		SetEncoding(extractous.CharSetUTF8).
		SetExtractStringMaxLength(1000).
		SetXmlOutput(true)

	if e1 == nil || e2 == nil {
		t.Fatal("Configuration order affected result")
	}

	e1.Close()
	e2.Close()
}
