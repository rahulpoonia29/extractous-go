package main

import (
	"fmt"
	"io"
	"log"

	"github.com/rahulpoonia29/extractous-go"
)

func main() {
	// Create a new extractor
	extractor := extractous.New()
	if extractor == nil {
		log.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

	// Extract to a stream for large files
	reader, metadata, err := extractor.ExtractFile("large_document.pdf")
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}
	defer reader.Close()

	// Read and print in chunks
	buffer := make([]byte, 4096)
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatalf("Read failed: %v", err)
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buffer[:n]))
	}

	// Print metadata
	fmt.Println("\nMetadata:")
	for key, values := range metadata {
		fmt.Printf("%s: %v\n", key, values)
	}
}
