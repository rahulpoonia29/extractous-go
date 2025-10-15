//! Type definitions for C FFI interface
//!
//! This module defines all opaque handle types, constants, and data structures
//! that are exposed through the C interface. These definitions are designed to
//! be ABI-compatible with C and can be safely used across the FFI boundary.
//!
//! ## Opaque Handles
//!
//! All complex Rust types are exposed as opaque C pointers. This provides:
//! - Memory safety by preventing direct field access from C
//! - ABI stability across Rust compiler versions
//! - Clear ownership semantics
//!
//! ## Constants
//!
//! All enumeration-like values are exposed as integer constants for C compatibility.
//! These match the internal Rust enums but use a C-friendly representation.

use std::os::raw::c_int;

// ============================================================================
// Opaque Handle Types
// ============================================================================

/// Opaque handle to an Extractor instance
///
/// Represents the main extraction engine. Create with `extractous_extractor_new()`
/// and destroy with `extractous_extractor_free()`.
///
/// ## Thread Safety
/// Not thread-safe. Use separate instances per thread or external synchronization.
///
/// ### Example
/// ```c
/// CExtractor* extractor = extractous_extractor_new();
/// // ... use extractor ...
/// extractous_extractor_free(extractor);
/// ```
#[repr(C)]
pub struct CExtractor {
    _private: [u8; 0],
}

/// Opaque handle to a StreamReader instance
///
/// Represents a buffered stream of extracted content. Read data using
/// `extractous_stream_read()` and free with `extractous_stream_free()`.
///
/// ### Example
/// ```c
/// CStreamReader* reader;
/// CMetadata* metadata;
/// extractous_extractor_extract_file(extractor, "doc.pdf", &reader, &metadata);
///
/// char buffer[4096];
/// size_t bytes_read;
/// while (extractous_stream_read(reader, buffer, sizeof(buffer), &bytes_read) == ERR_OK
///        && bytes_read > 0) {
///     // Process buffer...
/// }
/// extractous_stream_free(reader);
/// ```
#[repr(C)]
pub struct CStreamReader {
    _private: [u8; 0],
}

/// Opaque handle to a PdfParserConfig instance
///
/// Configuration for PDF document parsing. Create with `extractous_pdf_config_new()`,
/// configure with setter functions, and free with `extractous_pdf_config_free()`.
///
/// Note: Setters consume the old handle and return a new one (builder pattern).
#[repr(C)]
pub struct CPdfParserConfig {
    _private: [u8; 0],
}

/// Opaque handle to an OfficeParserConfig instance
///
/// Configuration for Microsoft Office document parsing. Create with
/// `extractous_office_config_new()` and free with `extractous_office_config_free()`.
#[repr(C)]
pub struct COfficeParserConfig {
    _private: [u8; 0],
}

/// Opaque handle to a TesseractOcrConfig instance
///
/// Configuration for Tesseract OCR engine. Create with `extractous_ocr_config_new()`
/// and free with `extractous_ocr_config_free()`.
///
/// Note: Requires Tesseract to be installed on the system.
#[repr(C)]
pub struct CTesseractOcrConfig {
    _private: [u8; 0],
}

// ============================================================================
// Metadata Structure
// ============================================================================

/// C-compatible metadata structure
///
/// Contains document metadata as parallel arrays of keys and values.
/// Multiple values for the same key are comma-separated.
///
/// ### Memory Layout
/// ```text
/// keys[0] -> "author\0"      values[0] -> "John Doe\0"
/// keys[1] -> "title\0"       values[1] -> "My Document\0"
/// keys[2] -> "keywords\0"    values[2] -> "pdf,test,sample\0"
/// ```
///
/// ### Memory Management
/// Must be freed with `extractous_metadata_free()` which will:
/// 1. Free all individual key strings
/// 2. Free all individual value strings
/// 3. Free the key array
/// 4. Free the value array
/// 5. Free the structure itself
///
/// ### Safety
/// - All strings are valid null-terminated UTF-8
/// - Arrays contain exactly `len` elements
/// - Do not modify the structure directly from C
/// - Do not free individual strings; use `extractous_metadata_free()`
#[repr(C)]
pub struct CMetadata {
    /// Array of pointers to key strings (null-terminated UTF-8)
    pub keys: *mut *mut libc::c_char,
    /// Array of pointers to value strings (null-terminated UTF-8, comma-separated if multiple)
    pub values: *mut *mut libc::c_char,
    /// Number of key-value pairs in the arrays
    pub len: libc::size_t,
}

// ============================================================================
// Character Set Constants
// ============================================================================

/// UTF-8 encoding (default, recommended)
///
/// Universal character encoding supporting all languages and emojis.
/// This is the default and recommended encoding for most use cases.
pub const CHARSET_UTF_8: c_int = 0;

/// US-ASCII encoding
///
/// 7-bit ASCII encoding. Use only if you're certain the content
/// contains only basic ASCII characters (0-127).
pub const CHARSET_US_ASCII: c_int = 1;

/// UTF-16 Big Endian encoding
///
/// 16-bit Unicode encoding with big-endian byte order.
/// Less common, use only if specifically required.
pub const CHARSET_UTF_16BE: c_int = 2;

// ============================================================================
// PDF OCR Strategy Constants
// ============================================================================

/// No OCR processing - extract only embedded text
///
/// Fastest option. Extracts only text that is already present in the PDF.
/// Images and scanned pages will not be processed.
///
/// Use when:
/// - PDF contains searchable text
/// - OCR is not needed
/// - Performance is critical
pub const PDF_OCR_NO_OCR: c_int = 0;

/// OCR only - ignore embedded text
///
/// Renders pages as images and performs OCR.
/// Ignores any embedded text in the PDF.
///
/// Use when:
/// - PDF text layer is corrupted or unreliable
/// - You need consistent OCR processing
pub const PDF_OCR_OCR_ONLY: c_int = 1;

/// Combined OCR and text extraction
///
/// Extracts embedded text AND performs OCR on images.
/// Provides most comprehensive extraction but is slower.
///
/// Use when:
/// - PDF has both text and scanned images
/// - Maximum content extraction is needed
pub const PDF_OCR_OCR_AND_TEXT_EXTRACTION: c_int = 2;

/// Automatic OCR strategy selection
///
/// Analyzes the PDF and automatically decides whether to use OCR.
/// Good balance between performance and coverage.
///
/// Use when:
/// - Processing mixed PDFs (some with text, some scanned)
/// - Want automatic optimization
pub const PDF_OCR_AUTO: c_int = 3;

// ============================================================================
// Additional Utility Constants
// ============================================================================

/// Default buffer size for stream reading (4KB)
///
/// Recommended buffer size for efficient stream reading.
/// Balances memory usage and I/O performance.
pub const DEFAULT_BUFFER_SIZE: libc::size_t = 4096;

/// Maximum recommended buffer size (1MB)
///
/// Large buffer for high-performance scenarios.
/// Use when processing very large documents.
pub const MAX_BUFFER_SIZE: libc::size_t = 1024 * 1024;

/// Default string extraction limit (100MB)
///
/// Default maximum length for extracted strings to prevent
/// excessive memory usage on very large documents.
pub const DEFAULT_STRING_MAX_LENGTH: c_int = 100 * 1024 * 1024;
