package extractous

// CharSet represents character encoding for text extraction.
//
// Character sets determine how bytes are interpreted as text characters. Most
// modern documents use UTF-8, which supports all Unicode characters. Other
// encodings are provided for legacy compatibility.
//
// # Supported Encodings
//
//   - CharSetUTF8: Unicode UTF-8 (default, recommended for all modern uses)
//   - CharSetUSASCII: US-ASCII (7-bit, legacy systems only)
//   - CharSetUTF16BE: UTF-16 Big Endian (some legacy systems)
//
// # When to Use Different Encodings
//
// UTF-8 (default): Use for all modern applications. It's the universal standard
// and supports all languages and symbols.
//
// US-ASCII: Only for legacy systems that cannot handle Unicode. This encoding
// only supports basic English characters (a-z, A-Z, 0-9, and basic punctuation).
// Any characters outside this range will be lost or corrupted.
//
// UTF-16BE: Rare. Only needed for specific legacy systems that require this
// encoding. Modern systems should use UTF-8.
//
// Example:
//
//	// Default UTF-8 (recommended)
//	extractor := extractous.New()
//
//	// Explicitly set UTF-8
//	extractor := extractous.New().
//	    SetEncoding(extractous.CharSetUTF8)
//
//	// Legacy ASCII (not recommended)
//	extractor := extractous.New().
//	    SetEncoding(extractous.CharSetUSASCII)
type CharSet int

const (
	// CharSetUTF8 is UTF-8 encoding (default, recommended).
	//
	// UTF-8 is the universal standard character encoding that supports all
	// Unicode characters, including:
	//   - All human languages (Latin, Cyrillic, Arabic, CJK, etc.)
	//   - Emojis and symbols
	//   - Mathematical notation
	//   - Technical symbols
	//
	// UTF-8 is backward compatible with ASCII, space-efficient, and the
	// default encoding for modern systems.
	//
	// Use this for all new applications unless you have a specific requirement
	// for another encoding.
	CharSetUTF8 CharSet = 0

	// CharSetUSASCII is US-ASCII encoding (7-bit, legacy only).
	//
	// US-ASCII supports only 128 characters:
	//   - Uppercase letters: A-Z
	//   - Lowercase letters: a-z
	//   - Digits: 0-9
	//   - Basic punctuation and symbols
	//
	// Characters outside this range (accented letters, non-Latin scripts, etc.)
	// will be lost or converted to placeholders.
	//
	// Only use this encoding if:
	//   - You need compatibility with very old systems
	//   - Your documents contain only basic English text
	//   - You have explicit requirements for ASCII-only output
	//
	// Warning: This encoding cannot represent most international text.
	CharSetUSASCII CharSet = 1

	// CharSetUTF16BE is UTF-16 Big Endian encoding (legacy).
	//
	// UTF-16 uses 16-bit code units and can represent all Unicode characters.
	// The "Big Endian" variant stores the most significant byte first.
	//
	// This encoding is less space-efficient than UTF-8 for most text and is
	// primarily used for:
	//   - Windows internal APIs (UTF-16LE, not UTF-16BE)
	//   - Java internal string representation
	//   - Some legacy systems
	//
	// Modern systems should use UTF-8 instead. Only use UTF-16BE if you have
	// explicit requirements for this encoding.
	CharSetUTF16BE CharSet = 2
)

// String returns the human-readable name of the character set.
//
// This is useful for logging, debugging, and displaying the current encoding
// configuration to users.
//
// Example:
//
//	charset := extractous.CharSetUTF8
//	fmt.Println(charset.String()) // Output: "UTF-8"
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

// PdfOcrStrategy defines how OCR is applied to PDF documents.
//
// PDF documents can contain different types of text content:
//   - Embedded text: Text that is directly stored in the PDF (selectable text)
//   - Image-based text: Text visible only as pixels in images (not selectable)
//
// The OCR strategy determines how the extractor handles these different types
// of content.
//
// # Strategy Comparison
//
// PdfOcrNoOcr (fastest):
//   - Only extracts embedded text
//   - No OCR processing
//   - Fast and efficient
//   - Good for: PDFs with selectable text, e-books, digital documents
//   - Bad for: Scanned documents, photos of text
//
// PdfOcrAuto (recommended):
//   - Automatically detects pages without embedded text
//   - Performs OCR only on those pages
//   - Balanced performance and accuracy
//   - Good for: Mixed documents, unknown document types
//   - The smart default for most use cases
//
// PdfOcrOcrOnly (specialized):
//   - Only performs OCR, ignores embedded text
//   - Useful when embedded text is corrupt or incorrect
//   - Slow, processes every page with OCR
//   - Good for: PDFs with broken text layers
//   - Bad for: General purpose extraction
//
// PdfOcrOcrAndTextExtraction (comprehensive):
//   - Extracts both embedded text AND performs OCR
//   - Most comprehensive but slowest
//   - May produce duplicate content
//   - Good for: Maximum content extraction, forensic analysis
//   - Bad for: Production systems (very slow)
//
// # Performance Implications
//
// OCR is computationally expensive:
//   - PdfOcrNoOcr: ~100-1000x faster than OCR strategies
//   - PdfOcrAuto: Variable (depends on document content)
//   - PdfOcrOcrOnly: Slowest, processes every page
//   - PdfOcrOcrAndTextExtraction: Slowest, maximum processing
//
// Example:
//
//	// Digital PDF with embedded text (fast)
//	pdfConfig := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrNoOcr)
//
//	// Scanned document (auto OCR when needed)
//	pdfConfig := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrAuto)
//
//	// Force OCR on all pages (slow but comprehensive)
//	pdfConfig := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrOcrAndTextExtraction)
type PdfOcrStrategy int

const (
	// PdfOcrNoOcr extracts only embedded text, no OCR processing (fastest).
	//
	// This strategy only extracts selectable text that is directly embedded in
	// the PDF. It does NOT perform OCR on images or scanned pages.
	//
	// Use when:
	//   - You know the PDFs contain selectable text
	//   - Performance is critical
	//   - Processing digital documents (not scans)
	//
	// This is the fastest strategy, typically 100-1000x faster than OCR-based
	// strategies.
	//
	// Example:
	//
	//	// Fast extraction for digital PDFs
	//	config := extractous.NewPdfConfig().
	//	    SetOcrStrategy(extractous.PdfOcrNoOcr)
	//
	// Note: Scanned documents and images will produce little or no text with
	// this strategy.
	PdfOcrNoOcr PdfOcrStrategy = 0

	// PdfOcrOcrOnly performs OCR on all pages, ignores embedded text.
	//
	// This strategy applies OCR to all pages regardless of whether they contain
	// embedded text. It's useful when the embedded text is corrupt, incorrect,
	// or lower quality than what OCR would produce.
	//
	// Use when:
	//   - The PDF has a broken or incorrect text layer
	//   - You want consistent OCR output across all pages
	//   - You need to extract from pure image PDFs
	//
	// Warning: This is very slow as it performs OCR on every page, even pages
	// that already have good embedded text.
	//
	// Example:
	//
	//	// Force OCR for PDFs with broken text layers
	//	config := extractous.NewPdfConfig().
	//	    SetOcrStrategy(extractous.PdfOcrOcrOnly)
	//
	// Note: Requires Tesseract OCR to be installed on the system.
	PdfOcrOcrOnly PdfOcrStrategy = 1

	// PdfOcrOcrAndTextExtraction extracts both embedded text AND performs OCR.
	//
	// This strategy is the most comprehensive, extracting both:
	//   1. Embedded text from the PDF text layer
	//   2. Text via OCR from images and visual content
	//
	// This can produce duplicate content if the same text appears both as
	// embedded text and in images.
	//
	// Use when:
	//   - Maximum content extraction is required
	//   - You need both text layers for comparison
	//   - Forensic analysis or complete document preservation
	//
	// Warning: This is the slowest strategy, combining all extraction methods.
	// It may also produce duplicate or redundant content.
	//
	// Example:
	//
	//	// Maximum extraction for forensic analysis
	//	config := extractous.NewPdfConfig().
	//	    SetOcrStrategy(extractous.PdfOcrOcrAndTextExtraction)
	//
	// Note: Best for offline processing where completeness matters more than
	// speed or deduplication.
	PdfOcrOcrAndTextExtraction PdfOcrStrategy = 2

	// PdfOcrAuto automatically decides based on page content (recommended).
	//
	// This strategy intelligently detects whether each page has embedded text:
	//   - Pages WITH embedded text: Extract directly (fast)
	//   - Pages WITHOUT embedded text: Apply OCR (slower)
	//
	// This provides the best balance of performance and completeness for
	// documents of unknown type or mixed content.
	//
	// Use when:
	//   - Document type is unknown
	//   - Handling mixed documents (some pages scanned, some digital)
	//   - You want good default behavior
	//
	// This is the recommended strategy for general-purpose PDF extraction.
	//
	// Example:
	//
	//	// Smart extraction that adapts to content
	//	config := extractous.NewPdfConfig().
	//	    SetOcrStrategy(extractous.PdfOcrAuto)
	//
	// Performance: Variable depending on content. Pages with embedded text are
	// processed quickly; only scanned pages incur OCR overhead.
	PdfOcrAuto PdfOcrStrategy = 3
)

// String returns the human-readable name of the OCR strategy.
//
// This is useful for logging, debugging, and displaying the current OCR
// configuration to users.
//
// Example:
//
//	strategy := extractous.PdfOcrAuto
//	fmt.Println(strategy.String()) // Output: "Auto"
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

// Error codes from the FFI layer.
//
// These constants map to error codes returned by the underlying C FFI library.
// They are used internally to construct Go error values. Application code should
// use the exported error variables (ErrNullPointer, ErrIO, etc.) instead of
// checking these raw codes.
//
// Internal use only.
const (
	errOK               = 0  // No error
	errNullPointer      = -1 // Null pointer passed to FFI
	errInvalidUTF8      = -2 // String is not valid UTF-8
	errInvalidString    = -3 // String parameter is invalid
	errExtractionFailed = -4 // Document extraction failed
	errIOError          = -5 // File I/O error
	errInvalidConfig    = -6 // Configuration is invalid
	errInvalidEnum      = -7 // Enum value is invalid
)
