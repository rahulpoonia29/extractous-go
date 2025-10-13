// Example demonstrating basic document extraction
package main

import (
	"fmt"
	"log"

	"github.com/rahulpoonia29/extractous-go"
)

func main() {
	// Create a new extractor instance
	extractor := src.New()
	defer extractor.Close()

	// Extract content from a file
	content, metadata, err := extractor.ExtractFileToString("sample.pdf")
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}

	// Print results
	fmt.Println("=== Extracted Content ===")
	fmt.Println(content)
	fmt.Println()

	fmt.Println("=== Metadata ===")
	for key, value := range metadata {
		fmt.Printf("%s: %s\n", key, value)
	}
}
