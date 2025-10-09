//! Configuration structures for extraction

use crate::ecore::{
    OfficeParserConfig as CoreOfficeConfig, PdfOcrStrategy, PdfParserConfig as CorePdfConfig,
    TesseractOcrConfig as CoreOcrConfig,
};
use crate::types::*;
use std::ffi::CStr;
use std::ptr;

// ============================================================================
// PDF Parser Config
// ============================================================================

/// Create new PDF parser config with default settings.
#[no_mangle]
pub extern "C" fn extractous_pdf_config_new() -> *mut CPdfParserConfig {
    let config = Box::new(CorePdfConfig::new());
    Box::into_raw(config) as *mut CPdfParserConfig
}

/// Sets the OCR strategy for PDF parsing.
///
/// # Safety
/// - `handle` must be a valid PdfParserConfig pointer.
/// - `strategy` must be a valid PDF_OCR_* constant.
/// - Returns a NEW handle; old handle is consumed and must not be used.
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

    let old_config = Box::from_raw(handle as *mut CorePdfConfig);
    let new_config = old_config.set_ocr_strategy(ocr_strategy);
    Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
}

/// If true, extract the literal inline embedded OBXImages. Use with caution.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_inline_images(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CorePdfConfig);
    let new_config = old_config.set_extract_inline_images(value);
    Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
}

/// If true, extract each unique inline image only once.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_unique_inline_images_only(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CorePdfConfig);
    let new_config = old_config.set_extract_unique_inline_images_only(value);
    Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
}

/// If true, try to extract text and its marked structure.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_marked_content(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CorePdfConfig);
    let new_config = old_config.set_extract_marked_content(value);
    Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
}

/// If true, try to extract the text of annotations.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_annotation_text(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CorePdfConfig);
    let new_config = old_config.set_extract_annotation_text(value);
    Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
}

/// Free PDF config.
///
/// # Safety
/// - `handle` must be a valid PdfParserConfig pointer.
/// - `handle` must not be used after this call.
#[no_mangle]
pub unsafe extern "C" fn extractous_pdf_config_free(handle: *mut CPdfParserConfig) {
    if !handle.is_null() {
        drop(Box::from_raw(handle as *mut CorePdfConfig));
    }
}

// ============================================================================
// Office Parser Config
// ============================================================================

/// Create new Office parser config.
#[no_mangle]
pub extern "C" fn extractous_office_config_new() -> *mut COfficeParserConfig {
    let config = Box::new(CoreOfficeConfig::new());
    Box::into_raw(config) as *mut COfficeParserConfig
}

/// Sets whether MSOffice parsers should extract macros.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_set_extract_macros(
    handle: *mut COfficeParserConfig,
    value: bool,
) -> *mut COfficeParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CoreOfficeConfig);
    let new_config = old_config.set_extract_macros(value);
    Box::into_raw(Box::new(new_config)) as *mut COfficeParserConfig
}

/// Whether to include deleted content from DOCX files.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_set_include_deleted_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) -> *mut COfficeParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CoreOfficeConfig);
    let new_config = old_config.set_include_deleted_content(value);
    Box::into_raw(Box::new(new_config)) as *mut COfficeParserConfig
}

/// Whether to include content from "moveFrom" sections in DOCX.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_set_include_move_from_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) -> *mut COfficeParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CoreOfficeConfig);
    let new_config = old_config.set_include_move_from_content(value);
    Box::into_raw(Box::new(new_config)) as *mut COfficeParserConfig
}

/// Whether to include text from drawing shapes.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_set_include_shape_based_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) -> *mut COfficeParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CoreOfficeConfig);
    let new_config = old_config.set_include_shape_based_content(value);
    Box::into_raw(Box::new(new_config)) as *mut COfficeParserConfig
}

/// Free Office config.
///
/// # Safety
/// - `handle` must be valid and not used after this call.
#[no_mangle]
pub unsafe extern "C" fn extractous_office_config_free(handle: *mut COfficeParserConfig) {
    if !handle.is_null() {
        drop(Box::from_raw(handle as *mut CoreOfficeConfig));
    }
}

// ============================================================================
// Tesseract OCR Config
// ============================================================================

/// Create new Tesseract OCR config.
#[no_mangle]
pub extern "C" fn extractous_ocr_config_new() -> *mut CTesseractOcrConfig {
    let config = Box::new(CoreOcrConfig::new());
    Box::into_raw(config) as *mut CTesseractOcrConfig
}

/// Sets the OCR language.
///
/// # Safety
/// - `handle` must be a valid TesseractOcrConfig pointer.
/// - `language` must be a valid null-terminated UTF-8 string.
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_language(
    handle: *mut CTesseractOcrConfig,
    language: *const libc::c_char,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() || language.is_null() {
        return ptr::null_mut();
    }

    let lang_str = match CStr::from_ptr(language).to_str() {
        Ok(s) => s,
        Err(_) => return ptr::null_mut(),
    };

    let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
    let new_config = old_config.set_language(lang_str);
    Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
}

/// Sets the DPI (dots per inch) for OCR.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_density(
    handle: *mut CTesseractOcrConfig,
    density: i32,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
    let new_config = old_config.set_density(density);
    Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
}

/// Sets the color depth for OCR.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_depth(
    handle: *mut CTesseractOcrConfig,
    depth: i32,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
    let new_config = old_config.set_depth(depth);
    Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
}

/// Sets whether to enable image preprocessing for OCR.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_enable_image_preprocessing(
    handle: *mut CTesseractOcrConfig,
    value: bool,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
    let new_config = old_config.set_enable_image_preprocessing(value);
    Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
}

/// Sets the timeout in seconds for the OCR process.
///
/// # Safety
/// - Returns a NEW handle; old handle is consumed.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_set_timeout_seconds(
    handle: *mut CTesseractOcrConfig,
    seconds: i32,
) -> *mut CTesseractOcrConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }
    let old_config = Box::from_raw(handle as *mut CoreOcrConfig);
    let new_config = old_config.set_timeout_seconds(seconds);
    Box::into_raw(Box::new(new_config)) as *mut CTesseractOcrConfig
}

/// Free OCR config.
///
/// # Safety
/// - `handle` must be valid and not used after this call.
#[no_mangle]
pub unsafe extern "C" fn extractous_ocr_config_free(handle: *mut CTesseractOcrConfig) {
    if !handle.is_null() {
        drop(Box::from_raw(handle as *mut CoreOcrConfig));
    }
}
