package extractous

/*
#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

// PdfConfig configures PDF document parsing behavior.
//
// PdfConfig provides fine-grained control over PDF extraction including OCR
// strategy, inline image extraction, marked content, and annotation text.
// Configuration objects use the builder pattern where each setter consumes
// the receiver and returns a new instance.
//
// # Default Configuration
//
// The default configuration includes:
//   - OCR strategy: PdfOcrNo (no OCR)
//   - Extract inline images: false
//   - Extract unique inline images only: true
//   - Extract marked content: false
//   - Extract annotation text: true
//
// # Usage Pattern
//
// Always use the returned value from setter methods:
//
//	// WRONG - config is consumed and invalid
//	config := extractous.NewPdfConfig()
//	config.SetOcrStrategy(extractous.PdfOcrAuto) // config is now invalid!
//
//	// CORRECT - use returned value
//	config := extractous.NewPdfConfig()
//	config = config.SetOcrStrategy(extractous.PdfOcrAuto)
//
//	// BEST - chain method calls
//	config := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrAuto).
//	    SetExtractAnnotationText(true)
//
// # Memory Management
//
// PdfConfig uses finalizers for automatic cleanup. When passed to
// Extractor.SetPdfConfig, the config is consumed and ownership transfers
// to the extractor. Do not use the config after passing it to an extractor.
//
// Example:
//
//		config := extractous.NewPdfConfig().
//		    SetOcrStrategy(extractous.PdfOcrAuto).
//		    SetExtractInlineImages(true)
//
//		extractor := extractous.New().
//		    SetPdfConfig(config) // config is consumed here.
//	        // Do not use 'config' after this point
type PdfConfig struct {
	ptr *C.struct_CPdfParserConfig
}

// NewPdfConfig creates a new PDF configuration with default settings.
//
// Returns a PdfConfig with default values suitable for most use cases.
// Returns nil if configuration creation failed (extremely rare).
//
// Example:
//
//	config := extractous.NewPdfConfig()
//	// Use with default settings or customize
func NewPdfConfig() *PdfConfig {
	ptr := C.extractous_pdf_config_new()
	if ptr == nil {
		return nil
		
	}

	cfg := &PdfConfig{ptr: ptr}
	runtime.SetFinalizer(cfg, (*PdfConfig).free)
	return cfg
}

// SetOcrStrategy sets the OCR (Optical Character Recognition) strategy for PDFs.
// This determines how OCR is applied to PDF documents:
//   - PdfOcrNo: No OCR, only extract embedded text (fastest, default)
//   - PdfOcrAuto: OCR only pages without embedded text (balanced)
//   - PdfOcrAlways: OCR all pages regardless of embedded text (slowest, most complete)
//
// OCR requires Tesseract to be installed on the system. If Tesseract is not
// available, OCR strategies will fail with an error.
//
// Parameters:
//   - strategy: One of PdfOcrNo, PdfOcrAuto, or PdfOcrAlways
//
// Example:
//
//	// Extract scanned PDFs automatically
//	config := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrAuto)
//
//	// Always OCR for maximum text extraction
//	config := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrAlways)
//
// This method consumes the receiver and returns a new PdfConfig.
// Returns nil if the strategy is invalid.
func (c *PdfConfig) SetOcrStrategy(strategy PdfOcrStrategy) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_ocr_strategy(c.ptr, C.int(strategy))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetExtractInlineImages enables or disables extraction of inline embedded images.
//
// When enabled, inline images embedded in the PDF content stream are extracted
// and included in the output. This can significantly increase extraction time
// and output size for image-heavy PDFs.
//
// Default: false (inline images are not extracted)
//
// Parameters:
//   - value: true to extract inline images, false to skip them
//
// Example:
//
//	// Extract all content including inline images
//	config := extractous.NewPdfConfig().
//	    SetExtractInlineImages(true)
//
// Note: This is separate from inline image OCR, which is controlled by
// SetOcrStrategy.
//
// This method consumes the receiver and returns a new PdfConfig.
func (c *PdfConfig) SetExtractInlineImages(value bool) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_extract_inline_images(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetExtractUniqueInlineImagesOnly extracts each unique inline image only once.
//
// When enabled (default), duplicate inline images are deduplicated and only
// extracted once. This is useful for PDFs that reuse the same images across
// multiple pages (e.g., logos, watermarks).
//
// When disabled, every occurrence of an inline image is extracted separately,
// which can be useful for preserving exact document structure but increases
// output size.
//
// Default: true (extract each unique image once)
//
// Parameters:
//   - value: true for deduplication, false to extract all occurrences
//
// Example:
//
//	// Extract every image occurrence (no deduplication)
//	config := extractous.NewPdfConfig().
//	    SetExtractInlineImages(true).
//	    SetExtractUniqueInlineImagesOnly(false)
//
// Note: This setting has no effect if SetExtractInlineImages is false.
//
// This method consumes the receiver and returns a new PdfConfig.
func (c *PdfConfig) SetExtractUniqueInlineImagesOnly(value bool) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_extract_unique_inline_images_only(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetExtractMarkedContent enables extraction of PDF marked content structure.
//
// Marked content in PDFs provides structural information like paragraph
// boundaries, headings, tables, etc. When enabled, this structure is preserved
// in the extracted output (particularly useful with XML output).
//
// Default: false (marked content is not extracted)
//
// Parameters:
//   - value: true to preserve marked content structure, false to ignore it
//
// Example:
//
//	// Preserve document structure
//	config := extractous.NewPdfConfig().
//	    SetExtractMarkedContent(true)
//
//	extractor := extractous.New().
//	    SetXmlOutput(true). // XML output shows structure
//	    SetPdfConfig(config)
//
// This is most useful when combined with XML output mode to preserve the
// document's logical structure.
//
// This method consumes the receiver and returns a new PdfConfig.
func (c *PdfConfig) SetExtractMarkedContent(value bool) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_extract_marked_content(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetExtractAnnotationText enables extraction of text from PDF annotations.
//
// PDF annotations include comments, highlights, sticky notes, and other markup.
// When enabled, the text content of these annotations is included in the
// extracted output.
//
// Default: true (annotation text is extracted)
//
// Parameters:
//   - value: true to extract annotation text, false to skip it
//
// Example:
//
//	// Skip annotations, only extract document content
//	config := extractous.NewPdfConfig().
//	    SetExtractAnnotationText(false)
//
//	// Include annotations (default behavior)
//	config := extractous.NewPdfConfig().
//	    SetExtractAnnotationText(true)
//
// Useful when you want only the document's primary content without user
// comments and markup.
//
// This method consumes the receiver and returns a new PdfConfig.
func (c *PdfConfig) SetExtractAnnotationText(value bool) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_extract_annotation_text(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

func (c *PdfConfig) free() {
	if c.ptr != nil {
		C.extractous_pdf_config_free(c.ptr)
		c.ptr = nil
	}
}

// OfficeConfig configures Microsoft Office document parsing behavior.
//
// OfficeConfig controls extraction from Microsoft Office formats (DOCX, XLSX,
// PPTX) and OpenOffice formats (ODT, ODS, ODP). It provides control over macro
// extraction, revision tracking content, and shape-based content.
//
// # Default Configuration
//
// The default configuration includes:
//   - Extract macros: true
//   - Include deleted content: false
//   - Include move from content: false
//   - Include shape-based content: true
//
// # Security Consideration
//
// Macros can contain malicious code. If extracting untrusted Office documents,
// consider disabling macro extraction with SetExtractMacros(false).
//
// # Usage Pattern
//
// Configuration objects use the builder pattern:
//
//	config := extractous.NewOfficeConfig().
//	    SetExtractMacros(false).           // More secure
//	    SetIncludeShapeBasedContent(true). // Extract text boxes
//	    SetIncludeDeletedContent(false)    // Skip tracked deletions
//
// Example:
//
//	// Secure configuration for untrusted documents
//	config := extractous.NewOfficeConfig().
//	    SetExtractMacros(false).
//	    SetIncludeDeletedContent(false)
//
//	extractor := extractous.New().
//	    SetOfficeConfig(config)
type OfficeConfig struct {
	ptr *C.struct_COfficeParserConfig
}

// NewOfficeConfig creates a new Office configuration with default settings.
//
// Returns an OfficeConfig with default values suitable for most use cases.
// Returns nil if configuration creation failed (extremely rare).
//
// Example:
//
//	config := extractous.NewOfficeConfig()
func NewOfficeConfig() *OfficeConfig {
	ptr := C.extractous_office_config_new()
	if ptr == nil {
		return nil
	}

	cfg := &OfficeConfig{ptr: ptr}
	runtime.SetFinalizer(cfg, (*OfficeConfig).free)
	return cfg
}

// SetExtractMacros enables or disables extraction of VBA macros.
//
// VBA (Visual Basic for Applications) macros are embedded scripts in Office
// documents. When enabled, macro code is extracted as text.
//
// Default: true (macros are extracted)
//
// Security Warning: Macros can contain malicious code. If processing untrusted
// documents, consider disabling macro extraction to avoid exposing potentially
// harmful code.
//
// Parameters:
//   - value: true to extract macros, false to skip them
//
// Example:
//
//	// Secure configuration - don't extract macros
//	config := extractous.NewOfficeConfig().
//	    SetExtractMacros(false)
//
//	// Full extraction including macros
//	config := extractous.NewOfficeConfig().
//	    SetExtractMacros(true)
//
// This method consumes the receiver and returns a new OfficeConfig.
func (c *OfficeConfig) SetExtractMacros(value bool) *OfficeConfig {
	newPtr := C.extractous_office_config_set_extract_macros(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetIncludeDeletedContent includes content marked as deleted in track changes.
//
// Office documents with track changes enabled can contain "deleted" content that
// still exists in the file but is marked for removal. When enabled, this deleted
// content is included in extraction.
//
// Default: false (deleted content is not extracted)
//
// Parameters:
//   - value: true to include deleted content, false to skip it
//
// Example:
//
//	// Include all content including deletions
//	config := extractous.NewOfficeConfig().
//	    SetIncludeDeletedContent(true)
//
//	// Only extract current document state (default)
//	config := extractous.NewOfficeConfig().
//	    SetIncludeDeletedContent(false)
//
// Useful for forensic analysis or when you need to see the full editing history.
//
// This method consumes the receiver and returns a new OfficeConfig.
func (c *OfficeConfig) SetIncludeDeletedContent(value bool) *OfficeConfig {
	newPtr := C.extractous_office_config_set_include_deleted_content(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetIncludeMoveFromContent includes "move from" content in track changes.
//
// When content is moved in a document with track changes, it appears both at the
// source ("move from") and destination ("move to") locations. When enabled, the
// "move from" content is included.
//
// Default: false (move from content is not extracted)
//
// Parameters:
//   - value: true to include move from content, false to skip it
//
// Example:
//
//	// Include moved content sources
//	config := extractous.NewOfficeConfig().
//	    SetIncludeMoveFromContent(true)
//
// This can result in duplicate content if both the source and destination of
// moved text are extracted.
//
// This method consumes the receiver and returns a new OfficeConfig.
func (c *OfficeConfig) SetIncludeMoveFromContent(value bool) *OfficeConfig {
	newPtr := C.extractous_office_config_set_include_move_from_content(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetIncludeShapeBasedContent includes text from shapes and text boxes.
//
// Office documents can contain text in shapes, text boxes, SmartArt, and other
// graphical elements. When enabled, this text is extracted along with the main
// document content.
//
// Default: true (shape-based content is extracted)
//
// Parameters:
//   - value: true to extract shape text, false to skip it
//
// Example:
//
//	// Extract all text including shapes and text boxes
//	config := extractous.NewOfficeConfig().
//	    SetIncludeShapeBasedContent(true)
//
//	// Only extract main document text (skip shapes)
//	config := extractous.NewOfficeConfig().
//	    SetIncludeShapeBasedContent(false)
//
// Disabling this can be useful when you only want the primary document flow
// without sidebars, callouts, and other graphical text elements.
//
// This method consumes the receiver and returns a new OfficeConfig.
func (c *OfficeConfig) SetIncludeShapeBasedContent(value bool) *OfficeConfig {
	newPtr := C.extractous_office_config_set_include_shape_based_content(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

func (c *OfficeConfig) free() {
	if c.ptr != nil {
		C.extractous_office_config_free(c.ptr)
		c.ptr = nil
	}
}

// OcrConfig configures Tesseract OCR (Optical Character Recognition) settings.
//
// OcrConfig controls how Tesseract OCR extracts text from images and scanned
// documents. Tesseract must be installed on the system for OCR to work.
//
// # System Requirements
//
// Tesseract OCR must be installed:
//   - Ubuntu/Debian: apt-get install tesseract-ocr
//   - macOS: brew install tesseract
//   - Windows: Download from https://github.com/UB-Mannheim/tesseract/wiki
//
// Language data files must also be installed for the languages you want to use.
//
// # Default Configuration
//
// The default configuration includes:
//   - Language: "eng" (English)
//   - DPI density: 300
//   - Color depth: 32
//   - Image preprocessing: false
//   - Timeout: 300 seconds (5 minutes)
//
// # Performance
//
// OCR is computationally expensive. Consider:
//   - Setting timeouts to prevent hanging on complex images
//   - Using lower DPI for faster processing (at the cost of accuracy)
//   - Enabling image preprocessing for better results on low-quality scans
//
// # Usage Pattern
//
// Configuration objects use the builder pattern:
//
//	config := extractous.NewOcrConfig().
//	    SetLanguage("eng").              // English language
//	    SetDensity(300).                 // 300 DPI
//	    SetEnableImagePreprocessing(true). // Better accuracy
//	    SetTimeoutSeconds(120)           // 2 minute timeout
//
// Example:
//
//	// High-quality OCR with preprocessing
//	ocrConfig := extractous.NewOcrConfig().
//	    SetLanguage("eng").
//	    SetDensity(400).
//	    SetEnableImagePreprocessing(true).
//	    SetTimeoutSeconds(300)
//
//	pdfConfig := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrAuto)
//
//	extractor := extractous.New().
//	    SetOcrConfig(ocrConfig).
//	    SetPdfConfig(pdfConfig)
type OcrConfig struct {
	ptr *C.struct_CTesseractOcrConfig
}

// NewOcrConfig creates a new OCR configuration with default settings.
//
// Returns an OcrConfig with default values suitable for English text at 300 DPI.
// Returns nil if configuration creation failed (extremely rare).
//
// Example:
//
//	config := extractous.NewOcrConfig()
func NewOcrConfig() *OcrConfig {
	ptr := C.extractous_ocr_config_new()
	if ptr == nil {
		return nil
	}

	cfg := &OcrConfig{ptr: ptr}
	runtime.SetFinalizer(cfg, (*OcrConfig).free)
	return cfg
}

// SetLanguage sets the OCR language(s) for text recognition.
//
// Tesseract supports many languages identified by 3-letter codes. Multiple
// languages can be specified with "+" separators for documents with mixed
// languages.
//
// Default: "eng" (English)
//
// Common language codes:
//   - "eng" - English
//   - "fra" - French
//   - "deu" - German
//   - "spa" - Spanish
//   - "ita" - Italian
//   - "por" - Portuguese
//   - "rus" - Russian
//   - "chi_sim" - Chinese Simplified
//   - "chi_tra" - Chinese Traditional
//   - "jpn" - Japanese
//   - "kor" - Korean
//   - "ara" - Arabic
//
// Parameters:
//   - lang: Language code(s), e.g., "eng", "fra", or "eng+fra" for multiple
//
// Example:
//
//	// English only
//	config := extractous.NewOcrConfig().
//	    SetLanguage("eng")
//
//	// Multiple languages
//	config := extractous.NewOcrConfig().
//	    SetLanguage("eng+fra+deu")
//
//	// Chinese simplified
//	config := extractous.NewOcrConfig().
//	    SetLanguage("chi_sim")
//
// Note: The corresponding language data must be installed. On Ubuntu/Debian,
// install with: apt-get install tesseract-ocr-<lang>
//
// This method consumes the receiver and returns a new OcrConfig.
func (c *OcrConfig) SetLanguage(lang string) *OcrConfig {
	cLang := cString(lang)
	defer freeString(cLang)

	newPtr := C.extractous_ocr_config_set_language(c.ptr, cLang)
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetDensity sets the DPI (dots per inch) density for image processing.
//
// Higher DPI improves OCR accuracy for small text but increases processing time
// and memory usage. Lower DPI is faster but may miss small text or details.
//
// Default: 300 DPI (recommended for most documents)
//
// Recommended values:
//   - 150-200 DPI: Fast processing, lower quality scans
//   - 300 DPI: Standard (default), good balance
//   - 400-600 DPI: High quality, small text, slower
//
// Parameters:
//   - dpi: DPI value, typically between 150 and 600
//
// Example:
//
//	// Fast processing
//	config := extractous.NewOcrConfig().
//	    SetDensity(200)
//
//	// High quality
//	config := extractous.NewOcrConfig().
//	    SetDensity(400)
//
// This method consumes the receiver and returns a new OcrConfig.
func (c *OcrConfig) SetDensity(dpi int) *OcrConfig {
	newPtr := C.extractous_ocr_config_set_density(c.ptr, C.int32_t(dpi))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetDepth sets the color depth (bits per pixel) for image processing.
//
// Higher color depth preserves more color information but uses more memory.
// For pure text recognition, lower depths are usually sufficient.
//
// Default: 32 bits per pixel
//
// Common values:
//   - 1: Black and white (1-bit)
//   - 8: Grayscale (8-bit)
//   - 24: True color (24-bit RGB)
//   - 32: True color with alpha (32-bit RGBA)
//
// Parameters:
//   - depth: Bit depth, typically 1, 8, 24, or 32
//
// Example:
//
//	// Grayscale for text-only documents
//	config := extractous.NewOcrConfig().
//	    SetDepth(8)
//
//	// Full color
//	config := extractous.NewOcrConfig().
//	    SetDepth(32)
//
// For most text OCR, grayscale (8-bit) is sufficient and uses less memory.
//
// This method consumes the receiver and returns a new OcrConfig.
func (c *OcrConfig) SetDepth(depth int) *OcrConfig {
	newPtr := C.extractous_ocr_config_set_depth(c.ptr, C.int32_t(depth))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetEnableImagePreprocessing enables or disables image preprocessing.
//
// When enabled, Tesseract applies preprocessing steps like noise removal,
// deskewing, and contrast enhancement to improve OCR accuracy. This is
// particularly helpful for poor-quality scans, photos of documents, or images
// with noise and artifacts.
//
// Default: false (no preprocessing)
//
// Preprocessing impacts:
//   - Pro: Better accuracy on low-quality images
//   - Pro: Handles skewed or rotated text better
//   - Con: Slower processing
//   - Con: May over-process high-quality scans
//
// Parameters:
//   - value: true to enable preprocessing, false to disable
//
// Example:
//
//	// For high-quality scans (no preprocessing needed)
//	config := extractous.NewOcrConfig().
//	    SetEnableImagePreprocessing(false)
//
//	// For low-quality scans or photos
//	config := extractous.NewOcrConfig().
//	    SetEnableImagePreprocessing(true)
//
// Enable this when working with:
//   - Photos of documents
//   - Faxes and low-quality scans
//   - Skewed or rotated images
//   - Images with backgrounds or noise
//
// This method consumes the receiver and returns a new OcrConfig.
func (c *OcrConfig) SetEnableImagePreprocessing(value bool) *OcrConfig {
	newPtr := C.extractous_ocr_config_set_enable_image_preprocessing(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetTimeoutSeconds sets the maximum processing time for OCR in seconds.
//
// OCR processing can be very slow for complex images. Setting a timeout prevents
// the process from hanging indefinitely on problematic images.
//
// Default: 300 seconds (5 minutes)
//
// Recommended values:
//   - 30-60 seconds: Quick processing, may timeout on complex pages
//   - 120-300 seconds: Standard (default)
//   - 600+ seconds: Very complex documents with many pages
//
// Parameters:
//   - seconds: Timeout in seconds, 0 for no timeout (not recommended)
//
// Example:
//
//	// Quick timeout for simple documents
//	config := extractous.NewOcrConfig().
//	    SetTimeoutSeconds(60)
//
//	// Extended timeout for complex documents
//	config := extractous.NewOcrConfig().
//	    SetTimeoutSeconds(600)
//
// When a timeout occurs, extraction will fail with an ErrOcrFailed error.
// Consider batch processing with retries for production systems.
//
// This method consumes the receiver and returns a new OcrConfig.
func (c *OcrConfig) SetTimeoutSeconds(seconds int) *OcrConfig {
	newPtr := C.extractous_ocr_config_set_timeout_seconds(c.ptr, C.int32_t(seconds))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

func (c *OcrConfig) free() {
	if c.ptr != nil {
		C.extractous_ocr_config_free(c.ptr)
		c.ptr = nil
	}
}
