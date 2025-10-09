package extractous

import "extractous-go/src"

// Re-export main types for convenience
type (
	Extractor    = src.Extractor
	StreamReader = src.StreamReader
	Metadata     = src.Metadata
	PdfConfig    = src.PdfConfig
	OfficeConfig = src.OfficeConfig
	OcrConfig    = src.OcrConfig
	CharSet      = src.CharSet
)

// Re-export enums
const (
	CharSetUTF8    = src.CharSetUTF8
	CharSetUSASCII = src.CharSetUSASCII
	CharSetUTF16BE = src.CharSetUTF16BE
)

// Re-export PDF OCR strategies
type PdfOcrStrategy = src.PdfOcrStrategy

const (
	PdfOcrNoOcr                = src.PdfOcrNoOcr
	PdfOcrOcrOnly              = src.PdfOcrOcrOnly
	PdfOcrOcrAndTextExtraction = src.PdfOcrOcrAndTextExtraction
	PdfOcrAuto                 = src.PdfOcrAuto
)

// Re-export constructor functions
var (
	New             = src.New
	NewPdfConfig    = src.NewPdfConfig
	NewOfficeConfig = src.NewOfficeConfig
	NewOcrConfig    = src.NewOcrConfig
)

// Re-export errors
var (
	ErrNullPointer   = src.ErrNullPointer
	ErrInvalidUTF8   = src.ErrInvalidUTF8
	ErrInvalidString = src.ErrInvalidString
	ErrExtraction    = src.ErrExtraction
	ErrIO            = src.ErrIO
	ErrInvalidConfig = src.ErrInvalidConfig
	ErrInvalidEnum   = src.ErrInvalidEnum
)
