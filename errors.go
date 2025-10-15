package extractous

/*
#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
)

// Sentinel errors that can be checked with errors.Is.
//
// These errors represent common failure modes in document extraction. They can
// be used with errors.Is() and errors.As() for error handling and classification.
//
// # Error Handling Pattern
//
//	content, _, err := extractor.ExtractFileToString("document.pdf")
//	if err != nil {
//	    if errors.Is(err, extractous.ErrIO) {
//	        // File not found or not readable
//	        log.Printf("File error: %v", err)
//	    } else if errors.Is(err, extractous.ErrExtraction) {
//	        // Document format issue or corruption
//	        log.Printf("Extraction error: %v", err)
//	    } else {
//	        // Other error
//	        log.Printf("Unknown error: %v", err)
//	    }
//	    return
//	}
//
// # Error Unwrapping
//
// All errors returned by this package can be unwrapped to get the sentinel error:
//
//	var extractErr *extractous.ExtractError
//	if errors.As(err, &extractErr) {
//	    fmt.Printf("Error code: %d\n", extractErr.Code)
//	    fmt.Printf("Message: %s\n", extractErr.Message)
//	    fmt.Printf("Type: %v\n", errors.Unwrap(extractErr))
//	}
var (
	// ErrNullPointer indicates a null pointer was encountered internally.
	//
	// This error typically indicates a programming error or corrupted internal
	// state. It should not occur during normal operation.
	//
	// Common causes:
	//   - Using an extractor after calling Close()
	//   - Passing a nil configuration to a setter
	//   - Internal FFI layer issues
	//
	// Example:
	//
	//	extractor := extractous.New()
	//	extractor.Close()
	//	_, _, err := extractor.ExtractFileToString("doc.pdf")
	//	if errors.Is(err, extractous.ErrNullPointer) {
	//	    // Extractor was already closed
	//	}
	ErrNullPointer = errors.New("null pointer provided")

	// ErrInvalidUTF8 indicates a string contains invalid UTF-8 sequences.
	//
	// This can occur when:
	//   - Extracting documents with corrupt character encoding
	//   - Document claims to be UTF-8 but contains invalid sequences
	//   - Character set mismatch between document and configuration
	//
	// Example handling:
	//
	//	if errors.Is(err, extractous.ErrInvalidUTF8) {
	//	    // Try with a different encoding
	//	    extractor = extractor.SetEncoding(extractous.CharSetUSASCII)
	//	    content, _, err = extractor.ExtractFileToString(path)
	//	}
	ErrInvalidUTF8 = errors.New("invalid UTF-8 string")

	// ErrInvalidString indicates string conversion or processing failed.
	//
	// This error occurs when:
	//   - String parameters cannot be converted properly
	//   - Extracted text contains characters that cannot be represented
	//   - Internal string buffer operations fail
	//
	// This is less common than ErrInvalidUTF8 and typically indicates a more
	// fundamental issue with the document or extraction process.
	ErrInvalidString = errors.New("string conversion failed")

	// ErrExtraction indicates document extraction failed.
	//
	// This is the most common error and can have many causes:
	//   - Unsupported document format
	//   - Corrupted or malformed document
	//   - Document is encrypted or password-protected
	//   - Document uses unsupported features
	//   - OCR processing failed
	//
	// The ExtractError.Message field usually contains specific details about
	// what went wrong.
	//
	// Example handling:
	//
	//	if errors.Is(err, extractous.ErrExtraction) {
	//	    var extractErr *extractous.ExtractError
	//	    if errors.As(err, &extractErr) {
	//	        log.Printf("Extraction failed: %s", extractErr.Message)
	//	        // Try alternative extraction method or skip document
	//	    }
	//	}
	ErrExtraction = errors.New("extraction failed")

	// ErrIO indicates an I/O operation failed.
	//
	// Common causes:
	//   - File not found
	//   - Permission denied
	//   - Disk read/write error
	//   - Network error (for URL extraction)
	//   - Out of disk space
	//
	// Example handling:
	//
	//	if errors.Is(err, extractous.ErrIO) {
	//	    if os.IsNotExist(err) {
	//	        log.Println("File not found")
	//	    } else if os.IsPermission(err) {
	//	        log.Println("Permission denied")
	//	    } else {
	//	        log.Printf("I/O error: %v", err)
	//	    }
	//	}
	ErrIO = errors.New("IO error")

	// ErrInvalidConfig indicates the provided configuration is invalid.
	//
	// This error occurs when:
	//   - Configuration parameter is out of valid range
	//   - Conflicting configuration options
	//   - Required configuration is missing
	//
	// Example:
	//
	//	config := extractous.NewOcrConfig().SetDensity(-100) // Invalid DPI
	//	extractor := extractous.New().SetOcrConfig(config)
	//	// May return ErrInvalidConfig when used
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrInvalidEnum indicates an invalid enum value was provided.
	//
	// This typically indicates a programming error where an enum constant
	// is used incorrectly or has an unexpected value.
	//
	// Example:
	//
	//	strategy := extractous.PdfOcrStrategy(999) // Invalid value
	//	config := extractous.NewPdfConfig().SetOcrStrategy(strategy)
	//	// May return ErrInvalidEnum
	ErrInvalidEnum = errors.New("invalid enum value")
)

// ExtractError wraps detailed extraction error information.
//
// ExtractError provides structured error information including an error code,
// a human-readable message, and a sentinel error for classification. It
// implements the error interface and supports error unwrapping for use with
// errors.Is() and errors.As().
//
// # Fields
//
//   - Code: Numeric error code from the FFI layer (negative values)
//   - Message: Detailed error message from the underlying extraction library
//   - Err: Wrapped sentinel error (one of ErrNullPointer, ErrIO, etc.)
//
// # Usage
//
// Use errors.Is() to check error types:
//
//	if errors.Is(err, extractous.ErrExtraction) {
//	    // Handle extraction error
//	}
//
// Use errors.As() to access error details:
//
//	var extractErr *extractous.ExtractError
//	if errors.As(err, &extractErr) {
//	    fmt.Printf("Error code: %d\n", extractErr.Code)
//	    fmt.Printf("Message: %s\n", extractErr.Message)
//	}
//
// # Example
//
//	content, _, err := extractor.ExtractFileToString("corrupt.pdf")
//	if err != nil {
//	    var extractErr *extractous.ExtractError
//	    if errors.As(err, &extractErr) {
//	        switch {
//	        case extractErr.Code == -4:
//	            log.Printf("Extraction failed: %s", extractErr.Message)
//	        case extractErr.Code == -5:
//	            log.Printf("I/O error: %s", extractErr.Message)
//	        default:
//	            log.Printf("Error %d: %s", extractErr.Code, extractErr.Message)
//	        }
//	    }
//	}
// errors.go

// ExtractError wraps detailed extraction error information.
type ExtractError struct {
    Code    int    // Numeric error code from FFI layer
    Message string // User-facing error message
    Err     error  // Wrapped sentinel error for errors.Is()
}

// newError creates an ExtractError from a C error code.
// Debug details are NOT fetched here to avoid memory overhead.
// Users must explicitly call Debug() to get detailed information.
func newError(code C.int) error {
    if code == errOK {
        return nil
    }

    var sentinelErr error
    switch int(code) {
    case errNullPointer:
        sentinelErr = ErrNullPointer
    case errInvalidUTF8:
        sentinelErr = ErrInvalidUTF8
    case errInvalidString:
        sentinelErr = ErrInvalidString
    case errExtractionFailed:
        sentinelErr = ErrExtraction
    case errIOError:
        sentinelErr = ErrIO
    case errInvalidConfig:
        sentinelErr = ErrInvalidConfig
    case errInvalidEnum:
        sentinelErr = ErrInvalidEnum
    default:
        sentinelErr = fmt.Errorf("unknown error code: %d", code)
    }

    // Get user-facing error message (fast, small string)
    cMsg := C.extractous_error_message(code)
    var msg string
    if cMsg != nil {
        msg = goString(cMsg)
        C.extractous_string_free(cMsg)
    }

    return &ExtractError{
        Code:    int(code),
        Message: msg,
        Err:     sentinelErr,
    }
}

// Error implements the error interface.
// Returns user-friendly error message.
func (e *ExtractError) Error() string {
    if e.Message != "" {
        return fmt.Sprintf("extractous error (code %d): %s", e.Code, e.Message)
    }
    return fmt.Sprintf("extractous error (code %d)", e.Code)
}

// Unwrap returns the underlying sentinel error for errors.Is() support
func (e *ExtractError) Unwrap() error {
    return e.Err
}

// Debug retrieves detailed debug information for the last error
// that occurred on the current thread.
//
// This function is EXPENSIVE - it formats the full error chain with
// backtrace (if RUST_BACKTRACE=1). Only call it when you actually
// need detailed debugging information.
//
// **Important**: This clears the stored error. Subsequent calls to
// Debug() will return empty string unless a new error occurs.
//
// Example:
//
//  _, _, err := extractor.ExtractFileToString("corrupt.pdf")
//  if err != nil {
//      var extractErr *extractous.ExtractError
//      if errors.As(err, &extractErr) {
//          // Show user-facing error
//          fmt.Printf("Error: %s\n", extractErr.Error())
//
//          // Optionally get debug details (for developers only)
//          if debug := extractErr.Debug(); debug != "" {
//              log.Printf("DEBUG:\n%s", debug)
//          }
//      }
//  }
func (e *ExtractError) Debug() string {
    cDebug := C.extractous_error_get_last_debug()
    if cDebug == nil {
        return ""
    }
    defer C.extractous_string_free(cDebug)
    return goString(cDebug)
}

// HasDebug checks if debug information is available for the last error
// on the current thread without retrieving it.
//
// This is useful to avoid the overhead of Debug() when no error is stored.
func (e *ExtractError) HasDebug() bool {
    return C.extractous_error_has_debug() != 0
}
