package examples

import (
	"extractous-go"
	"fmt"
	"log"
)

func main() {
	// Create a new extractor
	extractor := extractous.New()
	defer extractor.Close()

	// Extract text from a file
	content, metadata, err := extractor.ExtractFileToString("document.pdf")
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}

	// Print extracted content
	fmt.Println("Content:")
	fmt.Println(content)

	// Print metadata
	fmt.Println("\nMetadata:")
	for _, key := range metadata.Keys() {
		fmt.Printf("  %s: %s\n", key, metadata.Get(key))
	}
}
