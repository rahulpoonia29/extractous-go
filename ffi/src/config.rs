use crate::ecore::{
    OfficeParserConfig as CoreOfficeConfig, PdfOcrStrategy, PdfParserConfig as CorePdfConfig,
    TesseractOcrConfig as CoreOcrConfig,
};
use crate::types::*;
use std::ffi::CStr;
use std::os::raw::c_char;
use std::ptr;

/// Macro to safely update a config instance behind a raw pointer.
macro_rules! update_config {
    ($handle:expr, $T:ty, |$config_val:ident| $body:block) => {
        if $handle.is_null() {
            return;
        }
        unsafe {
            let config_ptr = $handle as *mut $T;
            let old_config = ptr::read(config_ptr);
            let new_config = {
                let $config_val = old_config;
                $body
            };
            ptr::write(config_ptr, new_config);
        }
    };
}

/// Creates a new PDF parser configuration with default settings.
/// The returned handle must be freed with `extractous_pdf_config_free()`
/// unless passed to an extractor, which will take ownership.
#[must_use]
#[unsafe(no_mangle)]
pub extern "C" fn extractous_pdf_config_new() -> *mut CPdfParserConfig {
    let config = Box::new(CorePdfConfig::new());
    Box::into_raw(config) as *mut CPdfParserConfig
}

/// Frees the memory associated with a PDF parser configuration.
/// Do not call this if the config has been attached to an extractor.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_pdf_config_free(handle: *mut CPdfParserConfig) {
    if !handle.is_null() {
        drop(unsafe { Box::from_raw(handle as *mut CorePdfConfig) });
    }
}

/// Sets the OCR strategy for PDF parsing. Modifies the config in-place.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_pdf_config_set_ocr_strategy(
    handle: *mut CPdfParserConfig,
    strategy: libc::c_int,
) {
    let ocr_strategy = match strategy {
        PDF_OCR_STRATEGY_NO_OCR => PdfOcrStrategy::NO_OCR,
        PDF_OCR_STRATEGY_OCR_ONLY => PdfOcrStrategy::OCR_ONLY,
        PDF_OCR_STRATEGY_OCR_AND_TEXT_EXTRACTION => PdfOcrStrategy::OCR_AND_TEXT_EXTRACTION,
        PDF_OCR_STRATEGY_AUTO => PdfOcrStrategy::AUTO,
        _ => return, // Invalid strategy, do nothing.
    };
    update_config!(handle, CorePdfConfig, |config| {
        config.set_ocr_strategy(ocr_strategy)
    });
}

/// Enables or disables extraction of inline images. Modifies the config in-place.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_inline_images(
    handle: *mut CPdfParserConfig,
    value: bool,
) {
    update_config!(handle, CorePdfConfig, |config| {
        config.set_extract_inline_images(value)
    });
}

/// If enabled, only unique inline images (by digest) will be extracted.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_unique_inline_images_only(
    handle: *mut CPdfParserConfig,
    value: bool,
) {
    update_config!(handle, CorePdfConfig, |config| {
        config.set_extract_unique_inline_images_only(value)
    });
}

/// Enables or disables extraction of text from marked content sections.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_marked_content(
    handle: *mut CPdfParserConfig,
    value: bool,
) {
    update_config!(handle, CorePdfConfig, |config| {
        config.set_extract_marked_content(value)
    });
}

/// Enables or disables extraction of text from annotations.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_annotation_text(
    handle: *mut CPdfParserConfig,
    value: bool,
) {
    update_config!(handle, CorePdfConfig, |config| {
        config.set_extract_annotation_text(value)
    });
}

/// Creates a new Office parser configuration with default settings.
#[must_use]
#[unsafe(no_mangle)]
pub extern "C" fn extractous_office_config_new() -> *mut COfficeParserConfig {
    let config = Box::new(CoreOfficeConfig::new());
    Box::into_raw(config) as *mut COfficeParserConfig
}

/// Frees the memory associated with an Office parser configuration.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_office_config_free(handle: *mut COfficeParserConfig) {
    if !handle.is_null() {
        drop(unsafe { Box::from_raw(handle as *mut CoreOfficeConfig) });
    }
}

/// Enables or disables macro extraction. Modifies the config in-place.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_office_config_set_extract_macros(
    handle: *mut COfficeParserConfig,
    value: bool,
) {
    update_config!(handle, CoreOfficeConfig, |config| {
        config.set_extract_macros(value)
    });
}

/// Enables or disables inclusion of deleted content (track changes).
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_office_config_set_include_deleted_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) {
    update_config!(handle, CoreOfficeConfig, |config| {
        config.set_include_deleted_content(value)
    });
}

/// Enables or disables inclusion of moved-from content (track changes).
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_office_config_set_include_move_from_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) {
    update_config!(handle, CoreOfficeConfig, |config| {
        config.set_include_move_from_content(value)
    });
}

/// Enables or disables inclusion of content from shapes.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_office_config_set_include_shape_based_content(
    handle: *mut COfficeParserConfig,
    value: bool,
) {
    update_config!(handle, CoreOfficeConfig, |config| {
        config.set_include_shape_based_content(value)
    });
}

/// Creates a new Tesseract OCR configuration with default settings.
#[must_use]
#[unsafe(no_mangle)]
pub extern "C" fn extractous_ocr_config_new() -> *mut CTesseractOcrConfig {
    let config = Box::new(CoreOcrConfig::new());
    Box::into_raw(config) as *mut CTesseractOcrConfig
}

/// Frees the memory associated with a Tesseract OCR configuration.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_ocr_config_free(handle: *mut CTesseractOcrConfig) {
    if !handle.is_null() {
        drop(unsafe { Box::from_raw(handle as *mut CoreOcrConfig) });
    }
}

/// Sets the OCR language. Modifies the config in-place.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_ocr_config_set_language(
    handle: *mut CTesseractOcrConfig,
    language: *const c_char,
) {
    if language.is_null() {
        return;
    }
    let lang_str = match unsafe { CStr::from_ptr(language).to_str() } {
        Ok(s) => s,
        Err(_) => return, // Invalid UTF-8, do nothing.
    };
    update_config!(handle, CoreOcrConfig, |config| {
        config.set_language(lang_str)
    });
}

/// Sets the DPI for OCR processing. Modifies the config in-place.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_ocr_config_set_density(
    handle: *mut CTesseractOcrConfig,
    density: i32,
) {
    update_config!(handle, CoreOcrConfig, |config| {
        config.set_density(density)
    });
}

/// Sets the bit depth for OCR processing.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_ocr_config_set_depth(
    handle: *mut CTesseractOcrConfig,
    depth: i32,
) {
    update_config!(handle, CoreOcrConfig, |config| { config.set_depth(depth) });
}

/// Enables or disables image preprocessing for OCR.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_ocr_config_set_enable_image_preprocessing(
    handle: *mut CTesseractOcrConfig,
    value: bool,
) {
    update_config!(handle, CoreOcrConfig, |config| {
        config.set_enable_image_preprocessing(value)
    });
}

/// Sets the timeout for the Tesseract process in seconds.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_ocr_config_set_timeout_seconds(
    handle: *mut CTesseractOcrConfig,
    seconds: i32,
) {
    update_config!(handle, CoreOcrConfig, |config| {
        config.set_timeout_seconds(seconds)
    });
}
