//! Type definitions for C FFI
//!
//! This module defines opaque handle types and constants that are exposed
//! to the C interface.

use std::os::raw::c_int;

/// Opaque handle to an Extractor instance
///
/// This is an opaque pointer that should only be used through the FFI functions.
/// The actual Extractor is stored on the heap.
#[repr(C)]
pub struct CExtractor {
    _private: [u8; 0],
}

/// Opaque handle to a StreamReader instance
#[repr(C)]
pub struct CStreamReader {
    _private: [u8; 0],
}

/// Opaque handle to a PdfParserConfig instance
#[repr(C)]
pub struct CPdfParserConfig {
    _private: [u8; 0],
}

/// Opaque handle to an OfficeParserConfig instance
#[repr(C)]
pub struct COfficeParserConfig {
    _private: [u8; 0],
}

/// Opaque handle to a TesseractOcrConfig instance
#[repr(C)]
pub struct CTesseractOcrConfig {
    _private: [u8; 0],
}

/// C-compatible metadata structure
///
/// Contains parallel arrays of keys and values, with length stored separately.
/// Both keys and values are null-terminated C strings.
#[repr(C)]
pub struct CMetadata {
    /// Array of key string pointers
    pub keys: *mut *mut libc::c_char,
    /// Array of value string pointers (comma-separated if multiple values)
    pub values: *mut *mut libc::c_char,
    /// Number of key-value pairs
    pub len: libc::size_t,
}

// Character Set Constants (only the 3 supported by extractous)
/// UTF-8 encoding (default)
pub const CHARSET_UTF8: c_int = 0;
/// US-ASCII encoding
pub const CHARSET_US_ASCII: c_int = 1;
/// UTF-16 Big Endian encoding
pub const CHARSET_UTF16BE: c_int = 2;

// PDF OCR Strategy Constants
/// No OCR, extract existing text only
pub const PDF_OCR_NO_OCR: c_int = 0;
/// OCR only, ignore existing text
pub const PDF_OCR_OCR_ONLY: c_int = 1;
/// OCR and extract existing text
pub const PDF_OCR_OCR_AND_TEXT_EXTRACTION: c_int = 2;
/// Automatically decide based on content
pub const PDF_OCR_AUTO: c_int = 3;
