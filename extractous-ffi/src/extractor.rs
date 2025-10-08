//! Extractor FFI functions
//!
//! This module provides the main extraction interface.

use crate::ecore::{CharSet, Extractor as CoreExtractor};
use crate::errors::*;
use crate::metadata::metadata_to_c;
use crate::types::*;
use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use std::ptr;

// ============================================================================
// Extractor Lifecycle
// ============================================================================

/// Create a new Extractor with default configuration
///
/// # Returns
///
/// Pointer to new Extractor, or NULL on failure.
/// Must be freed with `extractous_extractor_free`.
#[no_mangle]
pub extern "C" fn extractous_extractor_new() -> *mut CExtractor {
    let extractor = Box::new(CoreExtractor::new());
    Box::into_raw(extractor) as *mut CExtractor
}

/// Free an Extractor instance
///
/// # Safety
///
/// - `handle` must be a valid pointer returned by `extractous_extractor_new`
/// - `handle` must not be used after this call
/// - Calling this twice on the same pointer causes undefined behavior
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_free(handle: *mut CExtractor) {
    if !handle.is_null() {
        drop(Box::from_raw(handle as *mut CoreExtractor));
    }
}

// ============================================================================
// Extractor Configuration (Builder Pattern)
// ============================================================================

/// Set maximum length for extracted string content
///
/// # Safety
///
/// - `handle` must be a valid Extractor pointer
/// - Returns a NEW handle; old handle is consumed and must not be used
///
/// # Returns
///
/// New Extractor handle with updated config, or NULL on error.
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_set_extract_string_max_length(
    handle: *mut CExtractor,
    max_length: libc::c_int,
) -> *mut CExtractor {
    if handle.is_null() {
        return ptr::null_mut();
    }

    let old_extractor = Box::from_raw(handle as *mut CoreExtractor);
    let new_extractor = old_extractor.set_extract_string_max_length(max_length);

    Box::into_raw(Box::new(new_extractor)) as *mut CExtractor
}

/// Set character encoding for extraction
///
/// # Safety
///
/// - `handle` must be a valid Extractor pointer
/// - `encoding` must be a valid CHARSET_* constant
/// - Returns a NEW handle; old handle is consumed
///
/// # Returns
///
/// New Extractor handle, or NULL if encoding is invalid.
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_set_encoding(
    handle: *mut CExtractor,
    encoding: libc::c_int,
) -> *mut CExtractor {
    if handle.is_null() {
        return ptr::null_mut();
    }

    let charset = match encoding {
        CHARSET_UTF8 => CharSet::UTF_8,
        CHARSET_US_ASCII => CharSet::US_ASCII,
        CHARSET_UTF16BE => CharSet::UTF_16BE,
        _ => return ptr::null_mut(),
    };

    let old_extractor = Box::from_raw(handle as *mut CoreExtractor);
    let new_extractor = old_extractor.set_encoding(charset);

    Box::into_raw(Box::new(new_extractor)) as *mut CExtractor
}

/// Set PDF parser configuration
///
/// # Safety
///
/// - `handle` must be a valid Extractor pointer
/// - `config` must be a valid PdfParserConfig pointer
/// - Returns a NEW handle; old handle is consumed
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_set_pdf_config(
    handle: *mut CExtractor,
    config: *mut CPdfParserConfig,
) -> *mut CExtractor {
    if handle.is_null() || config.is_null() {
        return ptr::null_mut();
    }

    let pdf_config = &*(config as *mut crate::ecore::PdfParserConfig);
    let old_extractor = Box::from_raw(handle as *mut CoreExtractor);
    let new_extractor = old_extractor.set_pdf_config(pdf_config.clone());

    Box::into_raw(Box::new(new_extractor)) as *mut CExtractor
}

/// Set Office parser configuration
///
/// # Safety
///
/// Same safety requirements as `extractous_extractor_set_pdf_config`.
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_set_office_config(
    handle: *mut CExtractor,
    config: *mut COfficeParserConfig,
) -> *mut CExtractor {
    if handle.is_null() || config.is_null() {
        return ptr::null_mut();
    }

    let office_config = &*(config as *mut crate::ecore::OfficeParserConfig);
    let old_extractor = Box::from_raw(handle as *mut CoreExtractor);
    let new_extractor = old_extractor.set_office_config(office_config.clone());

    Box::into_raw(Box::new(new_extractor)) as *mut CExtractor
}

/// Set OCR configuration
///
/// # Safety
///
/// Same safety requirements as `extractous_extractor_set_pdf_config`.
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_set_ocr_config(
    handle: *mut CExtractor,
    config: *mut CTesseractOcrConfig,
) -> *mut CExtractor {
    if handle.is_null() || config.is_null() {
        return ptr::null_mut();
    }

    let ocr_config = &*(config as *mut crate::ecore::TesseractOcrConfig);
    let old_extractor = Box::from_raw(handle as *mut CoreExtractor);
    let new_extractor = old_extractor.set_ocr_config(ocr_config.clone());

    Box::into_raw(Box::new(new_extractor)) as *mut CExtractor
}

// ============================================================================
// Extraction Functions
// ============================================================================

/// Extract file content to string
///
/// # Safety
///
/// - `handle` must be a valid Extractor pointer
/// - `path` must be a valid null-terminated UTF-8 string
/// - `out_content` and `out_metadata` must be valid pointers
/// - Caller must free returned content with `extractous_string_free`
/// - Caller must free returned metadata with `extractous_metadata_free`
///
/// # Returns
///
/// ERR_OK on success, error code on failure.
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_extract_file_to_string(
    handle: *mut CExtractor,
    path: *const c_char,
    out_content: *mut *mut c_char,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    // Null pointer checks
    if handle.is_null() || path.is_null() || out_content.is_null() || out_metadata.is_null() {
        return ERR_NULL_POINTER;
    }

    // Convert C string to Rust str
    let path_str = match CStr::from_ptr(path).to_str() {
        Ok(s) => s,
        Err(_) => return ERR_INVALID_UTF8,
    };

    // Get reference to extractor
    let extractor = &*(handle as *mut CoreExtractor);

    // Perform extraction
    match extractor.extract_file_to_string(path_str) {
        Ok((content, metadata)) => {
            // Convert content to C string
            *out_content = match CString::new(content) {
                Ok(s) => s.into_raw(),
                Err(_) => return ERR_INVALID_STRING,
            };

            // Convert metadata to C structure
            *out_metadata = metadata_to_c(metadata);

            ERR_OK
        }
        Err(_) => ERR_EXTRACTION_FAILED,
    }
}

/// Extract file content to stream
///
/// # Safety
///
/// - `handle` must be a valid Extractor pointer
/// - `path` must be a valid null-terminated UTF-8 string
/// - `out_reader` and `out_metadata` must be valid pointers
/// - Caller must free returned reader with `extractous_stream_free`
/// - Caller must free returned metadata with `extractous_metadata_free`
///
/// # Returns
///
/// ERR_OK on success, error code on failure.
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_extract_file(
    handle: *mut CExtractor,
    path: *const c_char,
    out_reader: *mut *mut CStreamReader,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    if handle.is_null() || path.is_null() || out_reader.is_null() || out_metadata.is_null() {
        return ERR_NULL_POINTER;
    }

    let path_str = match CStr::from_ptr(path).to_str() {
        Ok(s) => s,
        Err(_) => return ERR_INVALID_UTF8,
    };

    let extractor = &*(handle as *mut CoreExtractor);

    match extractor.extract_file(path_str) {
        Ok((reader, metadata)) => {
            *out_reader = Box::into_raw(Box::new(reader)) as *mut CStreamReader;
            *out_metadata = metadata_to_c(metadata);
            ERR_OK
        }
        Err(_) => ERR_EXTRACTION_FAILED,
    }
}

/// Extract from byte array to string
///
/// # Safety
///
/// - `handle` must be a valid Extractor pointer
/// - `data` must point to at least `data_len` valid bytes
/// - `out_content` and `out_metadata` must be valid pointers
///
/// # Returns
///
/// ERR_OK on success, error code on failure.
#[no_mangle]
pub unsafe extern "C" fn extractous_extractor_extract_bytes_to_string(
    handle: *mut CExtractor,
    data: *const u8,
    data_len: libc::size_t,
    out_content: *mut *mut c_char,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    if handle.is_null() || data.is_null() || out_content.is_null() || out_metadata.is_null() {
        return ERR_NULL_POINTER;
    }

    let bytes = std::slice::from_raw_parts(data, data_len);
    let extractor = &*(handle as *mut CoreExtractor);

    match extractor.extract_bytes_to_string(bytes) {
        Ok((content, metadata)) => {
            *out_content = match CString::new(content) {
                Ok(s) => s.into_raw(),
                Err(_) => return ERR_INVALID_STRING,
            };
            *out_metadata = metadata_to_c(metadata);
            ERR_OK
        }
        Err(_) => ERR_EXTRACTION_FAILED,
    }
}

// ============================================================================
// String Memory Management
// ============================================================================

/// Free a string allocated by Rust
///
/// # Safety
///
/// - `s` must be a pointer returned by an extractous function
/// - `s` must not be used after this call
/// - Calling this twice on the same pointer causes undefined behavior
#[no_mangle]
pub unsafe extern "C" fn extractous_string_free(s: *mut c_char) {
    if !s.is_null() {
        drop(CString::from_raw(s));
    }
}
