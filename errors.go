package src

/*
#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
)

// Sentinel errors that can be checked with errors.Is
var (
	// ErrNullPointer indicates a null pointer was provided
	ErrNullPointer = errors.New("null pointer provided")
	// ErrInvalidUTF8 indicates invalid UTF-8 encoding
	ErrInvalidUTF8 = errors.New("invalid UTF-8 string")
	// ErrInvalidString indicates string conversion failed
	ErrInvalidString = errors.New("string conversion failed")
	// ErrExtraction indicates document extraction failed
	ErrExtraction = errors.New("extraction failed")
	// ErrIO indicates an IO operation failed
	ErrIO = errors.New("IO error")
	// ErrInvalidConfig indicates invalid configuration
	ErrInvalidConfig = errors.New("invalid configuration")
	// ErrInvalidEnum indicates invalid enum value
	ErrInvalidEnum = errors.New("invalid enum value")
)

// ExtractError wraps detailed extraction error information
type ExtractError struct {
	Code    int
	Message string
	Err     error // Wrapped sentinel error
}

// Error implements the error interface
func (e *ExtractError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("extractous error (code %d): %s", e.Code, e.Message)
	}
	return fmt.Sprintf("extractous error (code %d)", e.Code)
}

// Unwrap allows errors.Is and errors.As to work
func (e *ExtractError) Unwrap() error {
	return e.Err
}

// newError creates an ExtractError from a C error code
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

	// Get detailed error message from FFI
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
