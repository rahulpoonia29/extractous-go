//! Configuration structures and utilities for document extraction
//!
//! This module provides configuration interfaces for different parser types:
//! - **PDF Parser**: Controls PDF extraction behavior including OCR strategies
//! - **Office Parser**: Configures Microsoft Office document parsing options
//! - **Tesseract OCR**: Manages optical character recognition settings
//!
//! ## Configuration Pattern
//!
//! All configuration objects follow a builder pattern where setter functions
//! consume the old configuration and return a new one:
//!
//! ```c
//! CPdfParserConfig* config = extractous_pdf_config_new();
//! config = extractous_pdf_config_set_ocr_strategy(config, PDF_OCR_AUTO);
//! config = extractous_pdf_config_set_extract_inline_images(config, true);
//!
//! CExtractor* extractor = extractous_extractor_new();
//! extractor = extractous_extractor_set_pdf_config(extractor, config);
//!
//! // Config is now owned by extractor, don't free it separately
//! ```
//!
//! ## Safety Notes
//!
//! - Setter functions consume the input handle; do not use the old handle after calling setters
//! - Free standalone configs with appropriate free functions
//! - Configs attached to an extractor will be freed when the extractor is freed

use crate::ecore::{
    OfficeParserConfig as CoreOfficeConfig, PdfOcrStrategy, PdfParserConfig as CorePdfConfig,
    TesseractOcrConfig as CoreOcrConfig,
};
use crate::types::*;
use std::ffi::CStr;
use std::ptr;

// ============================================================================
// PDF Parser Configuration
// ============================================================================

/// Create a new PDF parser configuration with default settings
///
/// ### Default configuration:
/// - OCR strategy: NO_OCR (fastest, text extraction only)
/// - Extract inline images: false
/// - Extract unique inline images only: true
/// - Extract marked content: false
/// - Extract annotation text: false
///
/// Returns
/// Pointer to new PdfParserConfig. Must be freed with `extractous_pdf_config_free()`
/// unless attached to an extractor.
///
/// ```c
/// CPdfParserConfig* config = extractous_pdf_config_new();
/// if (config == NULL) {
///     // Handle allocation error
/// }
/// ```
#[no_mangle]
pub extern "C" fn extractous_pdf_config_new() -> *mut CPdfParserConfig {
    let config = Box::new(CorePdfConfig::new());
    Box::into_raw(config) as *mut CPdfParserConfig
}

/// Set the OCR strategy for PDF parsing
///
/// Determines how OCR is applied to PDF documents.
///
/// ##### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `strategy` - PDF_OCR_NO_OCR, PDF_OCR_OCR_ONLY, PDF_OCR_OCR_AND_TEXT_EXTRACTION, PDF_OCR_AUTO
///
/// ##### Returns
/// New PdfParserConfig handle with updated strategy, or NULL if invalid.
/// The input handle is consumed and must not be used.
///
/// ### Strategy Guide
/// - `PDF_OCR_NO_OCR`: Fastest, text-based PDFs only
/// - `PDF_OCR_OCR_ONLY`: Scanned documents, ignore existing text
/// - `PDF_OCR_OCR_AND_TEXT_EXTRACTION`: Mixed content, thorough extraction
/// - `PDF_OCR_AUTO`: Let the library decide (recommended)
///
/// ##### Safety
/// - Input handle is consumed; do not use after this call
/// - Returns NULL if handle is NULL or strategy is invalid
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_ocr_strategy(
    handle: *mut CPdfParserConfig,
    strategy: libc::c_int,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }

    let ocr_strategy = match strategy {
        PDF_OCR_NO_OCR => PdfOcrStrategy::NO_OCR,
        PDF_OCR_OCR_ONLY => PdfOcrStrategy::OCR_ONLY,
        PDF_OCR_OCR_AND_TEXT_EXTRACTION => PdfOcrStrategy::OCR_AND_TEXT_EXTRACTION,
        PDF_OCR_AUTO => PdfOcrStrategy::AUTO,
        _ => return ptr::null_mut(), // Invalid strategy
    };

    unsafe {
        let old_config = Box::from_raw(handle as *mut CorePdfConfig);
        let new_config = old_config.set_ocr_strategy(ocr_strategy);
        Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
    }
}

/// Enable or disable extraction of inline images from PDF
///
/// When enabled, extracts embedded image data from the PDF.
/// Can significantly increase memory usage and processing time.
///
/// ##### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `value` - true to extract inline images, false to skip
///
/// ##### Returns
/// New PdfParserConfig handle or NULL if handle is invalid.
///
/// ### Performance Impact
/// - Disabled (default): Fast, minimal memory
/// - Enabled: Slower, higher memory usage
///
/// ##### Safety
/// Input handle is consumed; do not use after this call.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_inline_images(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CorePdfConfig);
        let new_config = old_config.set_extract_inline_images(value);
        Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
    }
}

/// Extract each unique inline image only once
///
/// When enabled with inline image extraction, deduplicates repeated images.
///
/// ### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `value` - true for deduplication (recommended), false to extract all
///
/// ### Returns
/// New PdfParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_unique_inline_images_only(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CorePdfConfig);
        let new_config = old_config.set_extract_unique_inline_images_only(value);
        Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
    }
}

/// Extract text with marked content structure
///
/// Attempts to preserve document structure markers from the PDF.
///
/// ### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `value` - true to extract marked content, false otherwise
///
/// ### Returns
/// New PdfParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_marked_content(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CorePdfConfig);
        let new_config = old_config.set_extract_marked_content(value);
        Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
    }
}

/// Extract text from PDF annotations
///
/// Includes comments, highlights, and other annotation content.
///
/// ### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `value` - true to extract annotations, false to skip
///
/// ### Returns
/// New PdfParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_annotation_text(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CorePdfConfig);
        let new_config = old_config.set_extract_annotation_text(value);
        Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
    }
}

/// Free PDF parser configuration
///
/// ### Safety
/// - `handle` must be a valid PdfParserConfig pointer
/// - `handle` must not be used after this call
/// - Do not call this if config was attached to an extractor (it will be freed automatically)
///
/// ### Example
/// ```c
/// CPdfParserConfig* config = extractous_pdf_config_new();
/// // Use config...
/// extractous_pdf_config_free(config);  // Only if not attached to extractor
/// ```
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_free(handle: *mut CPdfParserConfig) {
    if !handle.is_null() {
        unsafe {
            drop(Box::from_raw(handle as *mut CorePdfConfig));
        }
    }
}

// ============================================================================
// Office Parser Configuration
// ============================================================================

/// Create a new Office parser configuration with default settings
///
/// Default configuration:
/// - Extract macros: false
/// - Include deleted content: false
/// - Include move-from content: false
/// - Include shape-based content: true
///
/// ### Returns
/// Pointer to new OfficeParserConfig. Must be freed with `extractous_office_config_free()`
/// unless attached to an extractor.
#[no_mangle]
pub extern "C" fn extractous_office_config_new() -> *mut COfficeParserConfig {
    let config = Box::new(CoreOfficeConfig::new());
    Box::into_raw(config) as *mut COfficeParserConfig
}

/// Enable or disable macro extraction from Office documents
///
/// **Security Warning**: Macros can contain malicious code. Only enable this
/// if you trust the document source and need macro content.
///
/// ### Arguments
/// * `handle` - Valid OfficeParserConfig pointer
/// * `value` - true to extract macros (security risk), false to skip (safer)
///
/// ### Returns
/// New OfficeParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_set_extract_macros(
    handle: *mut COfficeParserConfig,
    value: bool,
) -> *mut COfficeParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CoreOfficeConfig);
        let new_config = old_config.set_extract_macros(value);
        Box::into_raw(Box::new(new_config)) as *mut COfficeParserConfig
    }
}

/// Include deleted content from DOCX track changes
///
/// When enabled, extracts text that was deleted but is still present in
/// the document's revision history.
///
/// ### Arguments
/// * `handle` - Valid OfficeParserConfig pointer
/// * `value` - true to include deleted text, false to skip
///
/// ### Returns
/// New OfficeParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_set_include_deleted_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) -> *mut COfficeParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CoreOfficeConfig);
        let new_config = old_config.set_include_deleted_content(value);
        Box::into_raw(Box::new(new_config)) as *mut COfficeParserConfig
    }
}

/// Include "move-from" content in DOCX documents
///
/// Extracts text that was moved from one location to another during editing.
///
/// ### Arguments
/// * `handle` - Valid OfficeParserConfig pointer
/// * `value` - true to include moved text, false to skip
///
/// ### Returns
/// New OfficeParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_set_include_move_from_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) -> *mut COfficeParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CoreOfficeConfig);
        let new_config = old_config.set_include_move_from_content(value);
        Box::into_raw(Box::new(new_config)) as *mut COfficeParserConfig
    }
}

/// Include text from drawing shapes and text boxes
///
/// When enabled, extracts text from shapes, text boxes, and other drawing objects.
///
/// ### Arguments
/// * `handle` - Valid OfficeParserConfig pointer
/// * `value` - true to include shape text (recommended), false to skip
///
/// ### Returns
/// New OfficeParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_set_include_shape_based_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) -> *mut COfficeParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CoreOfficeConfig);
        let new_config = old_config.set_include_shape_based_content(value);
        Box::into_raw(Box::new(new_config)) as *mut COfficeParserConfig
    }
}

/// Free Office parser configuration
///
/// ### Safety
/// - `handle` must be valid and not used after this call
/// - Do not call this if config was attached to an extractor
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_free(handle: *mut COfficeParserConfig) {
    if !handle.is_null() {
        unsafe {
            drop(Box::from_raw(handle as *mut CoreOfficeConfig));
        }
    }
}

// ============================================================================
// Tesseract OCR Configuration
// ============================================================================

/// Create a new Tesseract OCR configuration with default settings
///
/// Default configuration:
/// - Language: "eng" (English)
/// - Density: 300 DPI
/// - Depth: 32 bits
/// - Image preprocessing: true
/// - Timeout: 300 seconds
///
/// ### Prerequisites
/// Tesseract must be installed on the system with appropriate language data files.
///
/// ### Returns
/// Pointer to new TesseractOcrConfig. Must be freed with `extractous_ocr_config_free()`
/// unless attached to an extractor.
#[no_mangle]
pub extern "C" fn extractous_ocr_config_new() -> *mut CTesseractOcrConfig {
    let config = Box::new(CoreOcrConfig::new());
    Box::into_raw(config) as *mut CTesseractOcrConfig
}

/// Set the OCR language
///
/// Specifies which language(s) Tesseract should use for recognition.
/// Multiple languages can be specified with '+' separator (e.g., "eng+fra").
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `language` - Null-terminated UTF-8 language code (e.g., "eng", "deu", "eng+fra")
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle or language is invalid.
///
/// ### Common Language Codes
/// - "eng" - English
/// - "deu" - German
/// - "fra" - French
/// - "spa" - Spanish
///
/// ### Requirements
/// The specified language data must be installed on the system.
/// On Debian/Ubuntu: `apt install tesseract-ocr-[lang]`
///
/// ### Safety
/// Input handle is consumed. Language string must be valid UTF-8.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_language(
    handle: *mut CTesseractOcrConfig,
    language: *const libc::c_char,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() || language.is_null() {
        return ptr::null_mut();
    }

    unsafe {
        let lang_str = match CStr::from_ptr(language).to_str() {
            Ok(s) => s,
            Err(_) => return ptr::null_mut(),
        };

        let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
        let new_config = old_config.set_language(lang_str);
        Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
    }
}

/// Set the DPI (dots per inch) for OCR processing
///
/// Higher DPI values can improve accuracy but increase processing time.
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `density` - DPI value (recommended: 150-600, default: 300)
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle is invalid.
///
/// ### Recommendations
/// - 150 DPI: Fast, lower quality
/// - 300 DPI: Balanced (default)
/// - 400-600 DPI: High quality, slower
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_density(
    handle: *mut CTesseractOcrConfig,
    density: i32,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
        let new_config = old_config.set_density(density);
        Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
    }
}

/// Set the color depth for OCR processing
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `depth` - Bit depth (typically 8, 24, or 32)
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_depth(
    handle: *mut CTesseractOcrConfig,
    depth: i32,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
        let new_config = old_config.set_depth(depth);
        Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
    }
}

/// Enable or disable image preprocessing for OCR
///
/// Preprocessing can improve OCR accuracy by normalizing image quality,
/// adjusting contrast, removing noise, etc.
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `value` - true to enable preprocessing (recommended), false to disable
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_enable_image_preprocessing(
    handle: *mut CTesseractOcrConfig,
    value: bool,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
        let new_config = old_config.set_enable_image_preprocessing(value);
        Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
    }
}

/// Set timeout for OCR processing
///
/// Prevents OCR from running indefinitely on problematic images.
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `seconds` - Timeout in seconds (0 = no timeout, default: 300)
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle is invalid.
///
/// # Recommendations
/// - 60-120 seconds: Fast processing, may timeout on complex images
/// - 300 seconds: Default, handles most documents
/// - 600+ seconds: Very complex documents
///
/// ### Safety
/// Input handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_timeout_seconds(
    handle: *mut CTesseractOcrConfig,
    seconds: i32,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    unsafe {
        let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
        let new_config = old_config.set_timeout_seconds(seconds);
        Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
    }
}

/// Free Tesseract OCR configuration
///
/// ### Safety
/// - `handle` must be valid and not used after this call
/// - Do not call this if config was attached to an extractor
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_free(handle: *mut CTesseractOcrConfig) {
    if !handle.is_null() {
        unsafe {
            drop(Box::from_raw(handle as *mut CoreOcrConfig));
        }
    }
}
