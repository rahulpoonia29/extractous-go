package main

import (
	"bytes"
	"errors"
	"extractous-go"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test file paths
const (
	testDataDir = "../testdata"
	testPDF     = "sample.pdf"
	testDOCX    = "sample.docx"
	testTXT     = "sample.txt"
)

func TestNew(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}
	defer extractor.Close()
}

func TestExtractFileToString(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	testFile := filepath.Join(testDataDir, testTXT)
	content, metadata, err := extractor.ExtractFileToString(testFile)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	if content == "" {
		t.Error("Expected non-empty content")
	}

	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}

	t.Logf("Content length: %d bytes", len(content))
	t.Logf("Metadata keys: %d", len(metadata))
}

func TestExtractFileToStringPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping PDF test in short mode")
	}

	extractor := extractous.New()
	defer extractor.Close()

	testFile := filepath.Join(testDataDir, testPDF)
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found", testFile)
	}

	content, metadata, err := extractor.ExtractFileToString(testFile)
	if err != nil {
		t.Fatalf("ExtractFileToString failed: %v", err)
	}

	if content == "" {
		t.Error("Expected non-empty content from PDF")
	}

	// Check common PDF metadata fields
	if contentType := metadata.Get("Content-Type"); contentType == "" {
		t.Error("Expected Content-Type in metadata")
	}

	t.Logf("PDF Content length: %d bytes", len(content))
	t.Logf("Content-Type: %s", metadata.Get("Content-Type"))
}

func TestExtractFile(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	testFile := filepath.Join(testDataDir, testTXT)
	reader, metadata, err := extractor.ExtractFile(testFile)
	if err != nil {
		t.Fatalf("ExtractFile failed: %v", err)
	}
	defer reader.Close()

	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	if len(content) == 0 {
		t.Error("Expected non-empty content")
	}

	t.Logf("Stream content length: %d bytes", len(content))
}

func TestExtractBytesToString(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	testContent := []byte("This is a test document.\nWith multiple lines.\n")

	content, metadata, err := extractor.ExtractBytesToString(testContent)
	if err != nil {
		t.Fatalf("ExtractBytesToString failed: %v", err)
	}

	if !strings.Contains(content, "test document") {
		t.Errorf("Expected content to contain 'test document', got: %s", content)
	}

	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}
}

func TestExtractBytes(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	testContent := []byte("Sample text content")

	reader, metadata, err := extractor.ExtractBytes(testContent)
	if err != nil {
		t.Fatalf("ExtractBytes failed: %v", err)
	}
	defer reader.Close()

	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	if len(content) == 0 {
		t.Error("Expected non-empty content")
	}
}

func TestExtractEmptyBytes(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	content, metadata, err := extractor.ExtractBytesToString([]byte{})
	if err != nil {
		t.Fatalf("ExtractBytesToString with empty bytes failed: %v", err)
	}

	if content != "" {
		t.Errorf("Expected empty content, got: %s", content)
	}

	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}
}

func TestExtractNonExistentFile(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	_, _, err := extractor.ExtractFileToString("nonexistent_file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	if !errors.Is(err, extractous.ErrIO) {
		t.Errorf("Expected ErrIO, got: %v", err)
	}
}

func TestPdfConfig(t *testing.T) {
	config := extractous.NewPdfConfig()
	if config == nil {
		t.Fatal("NewPdfConfig returned nil")
	}

	// Test builder pattern
	config = config.
		SetOcrStrategy(extractous.PdfOcrAuto).
		SetExtractInlineImages(true).
		SetExtractAnnotationText(true)

	if config == nil {
		t.Error("Config builder returned nil")
	}

	// Use config with extractor
	extractor := extractous.New().SetPdfConfig(config)
	if extractor == nil {
		t.Error("SetPdfConfig returned nil")
	}
	defer extractor.Close()
}

func TestOfficeConfig(t *testing.T) {
	config := extractous.NewOfficeConfig()
	if config == nil {
		t.Fatal("NewOfficeConfig returned nil")
	}

	config = config.
		SetExtractMacros(true).
		SetIncludeDeletedContent(false).
		SetIncludeShapeBasedContent(true)

	if config == nil {
		t.Error("Config builder returned nil")
	}

	extractor := extractous.New().SetOfficeConfig(config)
	if extractor == nil {
		t.Error("SetOfficeConfig returned nil")
	}
	defer extractor.Close()
}

func TestOcrConfig(t *testing.T) {
	config := extractous.NewOcrConfig()
	if config == nil {
		t.Fatal("NewOcrConfig returned nil")
	}

	config = config.
		SetLanguage("eng").
		SetDensity(300).
		SetDepth(8).
		SetEnableImagePreprocessing(true).
		SetTimeoutSeconds(30)

	if config == nil {
		t.Error("Config builder returned nil")
	}

	extractor := extractous.New().SetOcrConfig(config)
	if extractor == nil {
		t.Error("SetOcrConfig returned nil")
	}
	defer extractor.Close()
}

func TestSetEncoding(t *testing.T) {
	extractor := extractous.New()
	if extractor == nil {
		t.Fatal("New() returned nil")
	}
	defer extractor.Close()

	extractor = extractor.SetEncoding(extractous.CharSetUTF8)
	if extractor == nil {
		t.Error("SetEncoding returned nil")
	}
}

func TestSetExtractStringMaxLength(t *testing.T) {
	extractor := extractous.New().SetExtractStringMaxLength(1000)
	if extractor == nil {
		t.Fatal("SetExtractStringMaxLength returned nil")
	}
	defer extractor.Close()

	// Extract short content
	content, _, err := extractor.ExtractBytesToString(bytes.Repeat([]byte("a"), 2000))
	if err != nil {
		t.Fatalf("Extraction failed: %v", err)
	}

	// Content should be truncated
	if len(content) > 1100 { // Some margin for overhead
		t.Errorf("Expected content <= ~1000 chars, got %d", len(content))
	}
}

func TestMetadata(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	testFile := filepath.Join(testDataDir, testTXT)
	_, metadata, err := extractor.ExtractFileToString(testFile)
	if err != nil {
		t.Fatalf("Extraction failed: %v", err)
	}

	// Test metadata methods
	keys := metadata.Keys()
	if len(keys) == 0 {
		t.Error("Expected at least one metadata key")
	}

	// Test Get method
	contentType := metadata.Get("Content-Type")
	if contentType == "" {
		t.Error("Expected Content-Type metadata")
	}

	// Test Has method
	if !metadata.Has("Content-Type") {
		t.Error("Expected Has to return true for Content-Type")
	}

	// Test non-existent key
	if metadata.Has("NonExistentKey") {
		t.Error("Expected Has to return false for non-existent key")
	}
}

func TestStreamReaderMultipleReads(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	testFile := filepath.Join(testDataDir, testTXT)
	reader, _, err := extractor.ExtractFile(testFile)
	if err != nil {
		t.Fatalf("ExtractFile failed: %v", err)
	}
	defer reader.Close()

	// Read in chunks
	var total int
	buf := make([]byte, 64)
	for {
		n, err := reader.Read(buf)
		total += n
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
	}

	if total == 0 {
		t.Error("Expected to read some bytes")
	}

	t.Logf("Read %d bytes total", total)
}

func TestStreamReaderClose(t *testing.T) {
	extractor := extractous.New()
	defer extractor.Close()

	testFile := filepath.Join(testDataDir, testTXT)
	reader, _, err := extractor.ExtractFile(testFile)
	if err != nil {
		t.Fatalf("ExtractFile failed: %v", err)
	}

	err = reader.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Reading after close should return EOF
	buf := make([]byte, 64)
	n, err := reader.Read(buf)
	if err != io.EOF {
		t.Errorf("Expected EOF after close, got: %v (read %d bytes)", err, n)
	}
}

func TestConcurrentExtraction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	extractor := extractous.New()
	defer extractor.Close()

	testFile := filepath.Join(testDataDir, testTXT)

	// Run multiple extractions concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			_, _, err := extractor.ExtractFileToString(testFile)
			if err != nil {
				t.Errorf("Concurrent extraction %d failed: %v", id, err)
			}
			done <- true
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmark tests
func BenchmarkExtractFileToString(b *testing.B) {
	extractor := extractous.New()
	defer extractor.Close()

	testFile := filepath.Join(testDataDir, testTXT)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := extractor.ExtractFileToString(testFile)
		if err != nil {
			b.Fatalf("Extraction failed: %v", err)
		}
	}
}

func BenchmarkExtractBytesToString(b *testing.B) {
	extractor := extractous.New()
	defer extractor.Close()

	testContent := []byte("This is a test document with some content.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := extractor.ExtractBytesToString(testContent)
		if err != nil {
			b.Fatalf("Extraction failed: %v", err)
		}
	}
}
