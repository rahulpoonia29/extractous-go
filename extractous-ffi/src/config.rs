use crate::ecore::{PdfOcrStrategy, PdfParserConfig as CorePdfConfig};
// use crate::errors::*;
use crate::types::*;
use std::ptr;

/// Create new PDF parser config
#[unsafe(no_mangle)]
pub extern "C" fn extractous_pdf_config_new() -> *mut CPdfParserConfig {
    let config = Box::new(CorePdfConfig::new());
    Box::into_raw(config) as *mut CPdfParserConfig
}

/// Set OCR strategy for PDF config
#[unsafe(no_mangle)]
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
        _ => return ptr::null_mut(),
    };

    let old_config = Box::from_raw(handle as *mut CorePdfConfig) ;
    let new_config = old_config.set_ocr_strategy(ocr_strategy);

    Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
}

/// Set extract inline images flag
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_pdf_config_set_extract_inline_images(
    handle: *mut CPdfParserConfig,
    value: bool,
) -> *mut CPdfParserConfig {
    if handle.is_null() {
        return ptr::null_mut();
    }

    let old_config = unsafe { Box::from_raw(handle as *mut CorePdfConfig) };
    let new_config = old_config.set_extract_inline_images(value);

    Box::into_raw(Box::new(new_config)) as *mut CPdfParserConfig
}

/// Free PDF config
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_pdf_config_free(handle: *mut CPdfParserConfig) {
    if !handle.is_null() {
        let _ = unsafe { Box::from_raw(handle as *mut CorePdfConfig) };
    }
}

// Similar implementations for OfficeParserConfig and TesseractOcrConfig...
