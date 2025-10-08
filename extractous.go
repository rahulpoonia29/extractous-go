// Package extractous provides Go bindings for the extractous document extraction library.
//
// Extractous is a fast and efficient library for extracting text and metadata
// from various document formats including PDF, Word, Excel, PowerPoint, and many more.
//
// Example usage:
//
//	extractor, err := extractous.New()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer extractor.Close()
//
//	content, metadata, err := extractor.ExtractFileToString("document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Println("Content:", content)
//	fmt.Println("Metadata:", metadata)
package extractous

const (
	// Version of the extractous-go bindings
	Version = "0.1.0"
)
