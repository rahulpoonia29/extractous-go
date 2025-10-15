package extractous

/*
#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

// Metadata represents document metadata as key-value pairs.
//
// Metadata contains information about a document such as author, title, creation
// date, modification date, and other document properties. Each metadata field can
// have multiple values, so values are stored as string slices.
//
// # Common Metadata Fields
//
// Different document formats provide different metadata fields. Common fields
// include:
//
//   - "title" - Document title
//   - "author" - Document author(s)
//   - "creator" - Application that created the document
//   - "producer" - Application that produced the PDF (for PDFs)
//   - "subject" - Document subject/description
//   - "keywords" - Document keywords
//   - "created" - Creation date/time
//   - "modified" - Last modification date/time
//   - "Content-Type" - MIME type of the document
//   - "dc:title" - Dublin Core title (some formats)
//   - "dc:creator" - Dublin Core creator (some formats)
//
// # Multi-valued Fields
//
// Some fields can have multiple values, particularly "author" and "keywords":
//
//	metadata := Metadata{
//	    "author": []string{"Alice", "Bob"},
//	    "keywords": []string{"report", "quarterly", "finance"},
//	}
//
// # Case Sensitivity
//
// Metadata keys are case-sensitive. Some formats use lowercase keys ("author"),
// others use mixed case ("Author" or "dc:creator"). Always check the actual keys
// returned from extraction.
//
// # Usage Examples
//
// Basic access:
//
//	content, metadata, err := extractor.ExtractFileToString("document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get single value
//	title := metadata.Get("title")
//	author := metadata.Get("author")
//
//	// Get all values
//	allAuthors := metadata.GetAll("author")
//	for _, author := range allAuthors {
//	    fmt.Println("Author:", author)
//	}
//
//	// Check existence
//	if metadata.Has("keywords") {
//	    keywords := metadata.GetAll("keywords")
//	    fmt.Println("Keywords:", keywords)
//	}
//
// Iterate all metadata:
//
//	for _, key := range metadata.Keys() {
//	    values := metadata.GetAll(key)
//	    fmt.Printf("%s: %v\n", key, values)
//	}
//
// # Empty Metadata
//
// If a document has no metadata, an empty Metadata map is returned (not nil).
// Always safe to call methods on Metadata even when empty.
type Metadata map[string][]string

// metadataWrapper wraps C metadata for proper cleanup.
//
// This is an internal type used to manage the lifecycle of C metadata pointers.
// It ensures that C resources are freed when the Go garbage collector determines
// they are no longer needed.
//
// Internal use only.
type metadataWrapper struct {
	ptr *C.struct_CMetadata
}

// newMetadata converts C metadata to Go and sets up cleanup.
//
// This function:
//  1. Converts C metadata structure to Go map
//  2. Sets up a finalizer to free C resources
//  3. Splits comma-separated values into slices
//  4. Trims whitespace from values
//
// Returns an empty Metadata map if the C pointer is nil.
//
// Internal use only.
func newMetadata(cMeta *C.struct_CMetadata) Metadata {
	if cMeta == nil {
		return make(Metadata)
	}

	// Create wrapper for proper cleanup
	wrapper := &metadataWrapper{ptr: cMeta}
	runtime.SetFinalizer(wrapper, (*metadataWrapper).free)

	// Convert to Go map
	result := make(Metadata, int(cMeta.len))

	if cMeta.len == 0 {
		return result
	}

	// Convert C arrays to Go slices
	keys := unsafe.Slice(cMeta.keys, cMeta.len)
	values := unsafe.Slice(cMeta.values, cMeta.len)

	for i := 0; i < int(cMeta.len); i++ {
		key := C.GoString(keys[i])
		value := C.GoString(values[i])

		// Values are comma-separated in C, split them into individual values
		valueSlice := strings.Split(value, ",")
		// Trim whitespace from each value
		for j := range valueSlice {
			valueSlice[j] = strings.TrimSpace(valueSlice[j])
		}
		result[key] = valueSlice
	}

	return result
}

// free releases C metadata resources.
//
// This is called automatically by the garbage collector via the finalizer.
// Application code should not call this directly.
//
// Internal use only.
func (m *metadataWrapper) free() {
	if m.ptr != nil {
		C.extractous_metadata_free(m.ptr)
		m.ptr = nil
	}
}

// Get returns the first value for a metadata key.
//
// If the key exists and has one or more values, the first value is returned.
// If the key doesn't exist or has no values, an empty string is returned.
//
// This is the most convenient method for metadata fields that typically have
// a single value (like "title", "author", "created").
//
// Parameters:
//   - key: Metadata field name (case-sensitive)
//
// Returns:
//   - First value as a string, or "" if not found
//
// Example:
//
//	title := metadata.Get("title")
//	if title == "" {
//	    fmt.Println("No title")
//	} else {
//	    fmt.Println("Title:", title)
//	}
//
//	// For potentially multi-valued fields, Get returns first value
//	firstAuthor := metadata.Get("author")
func (m Metadata) Get(key string) string {
	if vals, ok := m[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// GetAll returns all values for a metadata key.
//
// Some metadata fields can have multiple values (particularly "author" and
// "keywords"). This method returns all values as a slice.
//
// If the key doesn't exist, nil is returned (not an empty slice).
//
// Parameters:
//   - key: Metadata field name (case-sensitive)
//
// Returns:
//   - Slice of all values, or nil if key not found
//
// Example:
//
//	// Get all authors
//	authors := metadata.GetAll("author")
//	if authors != nil {
//	    for _, author := range authors {
//	        fmt.Println("Author:", author)
//	    }
//	}
//
//	// Get all keywords
//	keywords := metadata.GetAll("keywords")
//	if keywords != nil {
//	    fmt.Println("Keywords:", strings.Join(keywords, ", "))
//	}
//
//	// Check for nil vs empty
//	if metadata.GetAll("nonexistent") == nil {
//	    fmt.Println("Key not found")
//	}
func (m Metadata) GetAll(key string) []string {
	return m[key]
}

// Has checks if a metadata key exists.
//
// Returns true if the key exists in the metadata, even if it has empty values.
// Returns false if the key doesn't exist.
//
// This is useful for distinguishing between a missing key and a key with an
// empty value.
//
// Parameters:
//   - key: Metadata field name (case-sensitive)
//
// Returns:
//   - true if key exists, false otherwise
//
// Example:
//
//	if metadata.Has("author") {
//	    author := metadata.Get("author")
//	    if author == "" {
//	        fmt.Println("Author field exists but is empty")
//	    } else {
//	        fmt.Println("Author:", author)
//	    }
//	} else {
//	    fmt.Println("No author field")
//	}
//
//	// Conditional processing
//	if metadata.Has("keywords") {
//	    processKeywords(metadata.GetAll("keywords"))
//	}
func (m Metadata) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// Keys returns all metadata keys.
//
// Returns a slice of all keys present in the metadata. The order is not
// guaranteed and may vary between calls.
//
// This is useful for iterating over all metadata fields without knowing the
// specific keys in advance.
//
// Returns:
//   - Slice of all metadata keys (may be empty but never nil)
//
// Example:
//
//	// Print all metadata
//	for _, key := range metadata.Keys() {
//	    values := metadata.GetAll(key)
//	    fmt.Printf("%s: %v\n", key, values)
//	}
//
//	// Count metadata fields
//	fmt.Printf("Found %d metadata fields\n", len(metadata.Keys()))
//
//	// Filter specific fields
//	for _, key := range metadata.Keys() {
//	    if strings.HasPrefix(key, "dc:") {
//	        // Process Dublin Core fields
//	        fmt.Printf("%s = %s\n", key, metadata.Get(key))
//	    }
//	}
//
// Note: The returned slice is a new allocation and can be modified without
// affecting the Metadata.
func (m Metadata) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
