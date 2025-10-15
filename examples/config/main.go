package main

import (
	"fmt"
	"log"

	"github.com/rahulpoonia29/extractous-go"
)

func main() {
	// Create PDF config with OCR
	pdfConfig := extractous.NewPdfConfig().
		SetOcrStrategy(extractous.PdfOcrAuto).
		SetExtractAnnotationText(true)

	// Create OCR config
	ocrConfig := extractous.NewOcrConfig().
		SetLanguage("eng").
		SetTimeoutSeconds(60)

	// Create extractor with configurations
	extractor := extractous.New().
		SetExtractStringMaxLength(50000).
		SetEncoding(extractous.CharSetUTF8).
		SetPdfConfig(pdfConfig).
		SetOcrConfig(ocrConfig)
	if extractor == nil {
		log.Fatal("Failed to create configured extractor")
	}
	defer extractor.Close()

	// Extract from a URL
	content, metadata, err := extractor.ExtractURLToString("https://example.com/sample.pdf")
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}

	fmt.Println("Extracted Content:")
	fmt.Println(content)

	fmt.Println("\nMetadata:")
	for key, values := range metadata {
		fmt.Printf("%s: %v\n", key, values)
	}
}
