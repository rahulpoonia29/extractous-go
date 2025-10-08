package extractous

// CharSet represents character encodings supported by extractous
type CharSet int

const (
	// CharSetUTF8 is UTF-8 encoding (default)
	CharSetUTF8 CharSet = 0
	// CharSetUSASCII is US-ASCII encoding
	CharSetUSASCII CharSet = 1
	// CharSetUTF16BE is UTF-16 Big Endian encoding
	CharSetUTF16BE CharSet = 2
)

// PdfOcrStrategy represents PDF OCR strategies
type PdfOcrStrategy int

const (
	// PdfOcrNoOcr extracts existing text only, no OCR
	PdfOcrNoOcr PdfOcrStrategy = 0
	// PdfOcrOcrOnly performs OCR only, ignores existing text
	PdfOcrOcrOnly PdfOcrStrategy = 1
	// PdfOcrOcrAndTextExtraction performs both OCR and text extraction
	PdfOcrOcrAndTextExtraction PdfOcrStrategy = 2
	// PdfOcrAuto automatically decides based on content
	PdfOcrAuto PdfOcrStrategy = 3
)
