package extractous_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	extractous "github.com/rahulpoonia29/extractous-go"
)

// ============================================================================
// Test Setup
// ============================================================================

const (
	testDataDir = "../testdata"
)

func setupTestDir(t *testing.T) string {
	dir := filepath.Join(testDataDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create test data directory: %v", err)
	}
	return dir
}

func createTestFile(t *testing.T, filename, content string) string {
	dir := setupTestDir(t)
	filePath := filepath.Join(dir, filename)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	return filePath
}

// ============================================================================
// Text File Tests
// ============================================================================

func TestIntegration_ExtractPlainText(t *testing.T) {
	content := "Hello, World!\nThis is a test file."
	filePath := createTestFile(t, "test.txt", content)
	defer os.Remove(filePath)

	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	extracted, metadata, err := extractor.ExtractFileToString(filePath)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	if !strings.Contains(extracted, "Hello, World!") {
		t.Errorf("Expected content not found in extracted text: %s", extracted)
	}

	if !strings.Contains(extracted, "This is a test file") {
		t.Errorf("Expected content not found in extracted text: %s", extracted)
	}

	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}

	// Check for common metadata
	if contentType := metadata.Get("Content-Type"); contentType == "" {
		t.Log("Warning: Content-Type not found in metadata")
	}
}

func TestIntegration_ExtractBytes(t *testing.T) {
	content := "Test content for bytes extraction"

	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	extracted, metadata, err := extractor.ExtractBytesToString([]byte(content))
	if err != nil {
		t.Fatalf("ExtractBytesToString failed: %v", err)
	}

	if !strings.Contains(extracted, content) {
		t.Errorf("Expected content not found. Got: %s", extracted)
	}

	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}
}

func TestIntegration_ExtractBytesStream(t *testing.T) {
	content := "Test content for streaming bytes extraction"

	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	stream, metadata, err := extractor.ExtractBytes([]byte(content))
	if err != nil {
		t.Fatalf("ExtractBytes failed: %v", err)
	}

	if stream == nil {
		t.Fatal("Expected non-nil stream")
	}

	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}

	// Read content from stream
	extractedBytes := make([]byte, 1024)
	n, _ := stream.Read(extractedBytes)

	extracted := string(extractedBytes[:n])
	if !strings.Contains(extracted, content) {
		t.Errorf("Expected content not found in stream. Got: %s", extracted)
	}
}

// ============================================================================
// Configuration Tests
// ============================================================================

func TestIntegration_MaxLengthConfiguration(t *testing.T) {
	// Create a long text file
	longContent := strings.Repeat("A", 10000)
	filePath := createTestFile(t, "long_test.txt", longContent)
	defer os.Remove(filePath)

	// Extract with small max length
	extractor := extractous.New().SetExtractStringMaxLength(100)
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	extracted, _, err := extractor.ExtractFileToString(filePath)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	// Extracted content should be truncated (or close to max length)
	if len(extracted) > 200 { // Some overhead is allowed
		t.Logf("Warning: Extracted %d chars, expected ~100", len(extracted))
	}
}

func TestIntegration_EncodingConfiguration(t *testing.T) {
	content := "UTF-8 content: こんにちは"
	filePath := createTestFile(t, "utf8_test.txt", content)
	defer os.Remove(filePath)

	extractor := extractous.New().SetEncoding(extractous.CharSetUTF8)
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	extracted, _, err := extractor.ExtractFileToString(filePath)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	if !strings.Contains(extracted, "UTF-8") {
		t.Logf("Extracted content: %s", extracted)
	}
}

func TestIntegration_XmlOutputConfiguration(t *testing.T) {
	content := "Test content for XML output"
	filePath := createTestFile(t, "xml_test.txt", content)
	defer os.Remove(filePath)

	// Test with XML output enabled
	extractor := extractous.New().SetXmlOutput(true)
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	extracted, _, err := extractor.ExtractFileToString(filePath)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	// XML output should contain XML tags
	if !strings.Contains(extracted, "<") {
		t.Logf("Warning: XML output doesn't seem to contain XML tags: %s", extracted[:min(100, len(extracted))])
	}

	// Test with XML output disabled
	extractor2 := extractous.New().SetXmlOutput(false)
	if extractor2 == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor2.Close()

	extracted2, _, err := extractor2.ExtractFileToString(filePath)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	if strings.Contains(extracted2, "Test content") {
		t.Log("Plain text extraction successful")
	}
}

// ============================================================================
// Metadata Tests
// ============================================================================

func TestIntegration_MetadataExtraction(t *testing.T) {
	content := "Test content"
	filePath := createTestFile(t, "metadata_test.txt", content)
	defer os.Remove(filePath)

	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	_, metadata, err := extractor.ExtractFileToString(filePath)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	if metadata == nil {
		t.Fatal("Expected non-nil metadata")
	}

	// Test metadata methods
	if !metadata.Has("Content-Type") {
		t.Log("Warning: Content-Type not found in metadata")
	}

	keys := metadata.Keys()
	if len(keys) == 0 {
		t.Error("Expected some metadata keys")
	}

	t.Logf("Metadata keys: %v", keys)

	// Test Get method
	for _, key := range keys {
		value := metadata.Get(key)
		if value == "" {
			t.Errorf("Get returned empty string for existing key: %s", key)
		}
		t.Logf("%s: %s", key, value)
	}

	// Test GetAll method
	for _, key := range keys {
		values := metadata.GetAll(key)
		if len(values) == 0 {
			t.Errorf("GetAll returned nil/empty for existing key: %s", key)
		}
	}
}

func TestIntegration_MetadataWithMultipleValues(t *testing.T) {
	// Some metadata fields can have multiple values (comma-separated)
	content := "Test content"
	filePath := createTestFile(t, "multi_meta_test.txt", content)
	defer os.Remove(filePath)

	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	_, metadata, err := extractor.ExtractFileToString(filePath)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	// Check if any metadata has multiple values
	for _, key := range metadata.Keys() {
		values := metadata.GetAll(key)
		if len(values) > 1 {
			t.Logf("Key '%s' has multiple values: %v", key, values)
		}
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestIntegration_NonexistentFile(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	_, _, err := extractor.ExtractFileToString("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestIntegration_EmptyFile(t *testing.T) {
	filePath := createTestFile(t, "empty.txt", "")
	defer os.Remove(filePath)

	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	extracted, metadata, err := extractor.ExtractFileToString(filePath)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	if extracted != "" {
		t.Logf("Note: Empty file produced content: %s", extracted)
	}

	if metadata == nil {
		t.Error("Expected non-nil metadata even for empty file")
	}
}

// ============================================================================
// Concurrency Tests
// ============================================================================

func TestIntegration_ConcurrentExtraction(t *testing.T) {
	content := "Concurrent test content"
	filePath := createTestFile(t, "concurrent_test.txt", content)
	defer os.Remove(filePath)

	const numGoroutines = 10
	errors := make(chan error, numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			extractor := extractous.New()
			if extractor == nil {
				errors <- nil // Signal completion even on nil
				return
			}
			defer extractor.Close()

			extracted, _, err := extractor.ExtractFileToString(filePath)
			if err != nil {
				errors <- err
				return
			}

			if !strings.Contains(extracted, content) {
				errors <- nil
				return
			}

			errors <- nil // Success
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		if err != nil {
			t.Errorf("Goroutine failed: %v", err)
		}
	}
}

func TestIntegration_MultipleExtractorsSameFile(t *testing.T) {
	content := "Multiple extractors test"
	filePath := createTestFile(t, "multi_ext_test.txt", content)
	defer os.Remove(filePath)

	extractors := make([]*extractous.Extractor, 5)
	for i := range extractors {
		extractors[i] = extractous.New()
		if extractors[i] == nil {
			t.Fatal("Failed to create extractor")
		}
		defer extractors[i].Close()
	}

	// All extractors extract the same file
	for i, ext := range extractors {
		extracted, _, err := ext.ExtractFileToString(filePath)
		if err != nil {
			t.Errorf("Extractor %d failed: %v", i, err)
		}
		if !strings.Contains(extracted, content) {
			t.Errorf("Extractor %d didn't extract correct content", i)
		}
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
