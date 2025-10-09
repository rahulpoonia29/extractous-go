package src

// CharSet represents character encoding for text extraction
type CharSet int

const (
	// CharSetUTF8 is UTF-8 encoding (default)
	CharSetUTF8 CharSet = 0
	// CharSetUSASCII is US-ASCII encoding
	CharSetUSASCII CharSet = 1
	// CharSetUTF16BE is UTF-16 Big Endian encoding
	CharSetUTF16BE CharSet = 2
)

// String returns the string representation of the character set
func (c CharSet) String() string {
	switch c {
	case CharSetUTF8:
		return "UTF-8"
	case CharSetUSASCII:
		return "US-ASCII"
	case CharSetUTF16BE:
		return "UTF-16BE"
	default:
		return "Unknown"
	}
}

// PdfOcrStrategy defines how OCR is applied to PDF documents
type PdfOcrStrategy int

const (
	// PdfOcrNoOcr extracts existing text only, no OCR
	PdfOcrNoOcr PdfOcrStrategy = 0
	// PdfOcrOcrOnly performs OCR only, ignores existing text
	PdfOcrOcrOnly PdfOcrStrategy = 1
	// PdfOcrOcrAndTextExtraction performs both OCR and text extraction
	PdfOcrOcrAndTextExtraction PdfOcrStrategy = 2
	// PdfOcrAuto automatically decides based on document content
	PdfOcrAuto PdfOcrStrategy = 3
)

// String returns the string representation of the OCR strategy
func (s PdfOcrStrategy) String() string {
	switch s {
	case PdfOcrNoOcr:
		return "NoOCR"
	case PdfOcrOcrOnly:
		return "OCROnly"
	case PdfOcrOcrAndTextExtraction:
		return "OCRAndTextExtraction"
	case PdfOcrAuto:
		return "Auto"
	default:
		return "Unknown"
	}
}

// Error codes from FFI
const (
	errOK               = 0
	errNullPointer      = -1
	errInvalidUTF8      = -2
	errInvalidString    = -3
	errExtractionFailed = -4
	errIOError          = -5
	errInvalidConfig    = -6
	errInvalidEnum      = -7
)
