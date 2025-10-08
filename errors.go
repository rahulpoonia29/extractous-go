package extractous

import "errors"
import "C"

var (
	ErrNullPointer      = errors.New("null pointer provided")
	ErrInvalidUTF8      = errors.New("invalid UTF-8 string")
	ErrInvalidString    = errors.New("string conversion failed")
	ErrExtractionFailed = errors.New("extraction failed")
	ErrIOError          = errors.New("I/O error")
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrInvalidEnum      = errors.New("invalid enum value")
)

func errorFromCode(code C.int) error {
	switch code {
	case C.ERR_OK:
		return nil
	case C.ERR_NULL_POINTER:
		return ErrNullPointer
	case C.ERR_INVALID_UTF8:
		return ErrInvalidUTF8
	case C.ERR_INVALID_STRING:
		return ErrInvalidString
	case C.ERR_EXTRACTION_FAILED:
		return ErrExtractionFailed
	case C.ERR_IO_ERROR:
		return ErrIOError
	case C.ERR_INVALID_CONFIG:
		return ErrInvalidConfig
	case C.ERR_INVALID_ENUM:
		return ErrInvalidEnum
	default:
		return errors.New("unknown error")
	}
}
