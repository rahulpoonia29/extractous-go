use crate::ecore::{CharSet, Extractor as CoreExtractor};
use crate::errors::*;
use crate::metadata::metadata_to_c;
use crate::types::*;
use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use std::ptr;

/// Creates a new `Extractor` with a default configuration.
/// The returned handle must be freed with `extractous_extractor_free`.
#[unsafe(no_mangle)]
#[must_use]
pub extern "C" fn extractous_extractor_new() -> *mut CExtractor {
    let extractor = Box::new(CoreExtractor::new());
    Box::into_raw(extractor) as *mut CExtractor
}

/// Frees the memory associated with an `Extractor` handle.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_free(handle: *mut CExtractor) {
    if !handle.is_null() {
        unsafe {
            drop(Box::from_raw(handle as *mut CoreExtractor));
        }
    }
}

/// A macro to safely update an Extractor instance behind a raw pointer.
macro_rules! update_extractor {
    ($handle:expr, |$extractor_val:ident| $body:block) => {
        if $handle.is_null() {
            return;
        }
        unsafe {
            let extractor_ptr = $handle as *mut CoreExtractor;
            let old_extractor = ptr::read(extractor_ptr);
            let new_extractor = {
                let $extractor_val = old_extractor;
                $body
            };
            ptr::write(extractor_ptr, new_extractor);
        }
    };
}

/// Sets the maximum length for extracted string content.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_set_extract_string_max_length_mut(
    handle: *mut CExtractor,
    max_length: libc::c_int,
) {
    update_extractor!(handle, |extractor| {
        extractor.set_extract_string_max_length(max_length as i32)
    });
}

/// Sets the character encoding for the extracted text.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_set_encoding_mut(
    handle: *mut CExtractor,
    encoding: libc::c_int,
) {
    update_extractor!(handle, |extractor| {
        let charset = match encoding {
            CHARSET_UTF_8 => CharSet::UTF_8,
            CHARSET_US_ASCII => CharSet::US_ASCII,
            CHARSET_UTF_16BE => CharSet::UTF_16BE,
            _ => return,
        };
        extractor.set_encoding(charset)
    });
}

/// Sets the configuration for the PDF parser.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_set_pdf_config_mut(
    handle: *mut CExtractor,
    config: *const CPdfParserConfig,
) {
    if config.is_null() {
        return;
    }
    update_extractor!(handle, |extractor| {
        let pdf_config = &*(config as *const crate::ecore::PdfParserConfig);
        extractor.set_pdf_config(pdf_config.clone())
    });
}

/// Sets the configuration for the Office document parser.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_set_office_config_mut(
    handle: *mut CExtractor,
    config: *const COfficeParserConfig,
) {
    if config.is_null() {
        return;
    }
    update_extractor!(handle, |extractor| {
        let office_config = &*(config as *const crate::ecore::OfficeParserConfig);
        extractor.set_office_config(office_config.clone())
    });
}

/// Sets the configuration for Tesseract OCR.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_set_ocr_config_mut(
    handle: *mut CExtractor,
    config: *const CTesseractOcrConfig,
) {
    if config.is_null() {
        return;
    }
    update_extractor!(handle, |extractor| {
        let ocr_config = &*(config as *const crate::ecore::TesseractOcrConfig);
        extractor.set_ocr_config(ocr_config.clone())
    });
}

/// Sets whether to output structured XML instead of plain text.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_set_xml_output_mut(
    handle: *mut CExtractor,
    xml_output: bool,
) {
    update_extractor!(handle, |extractor| { extractor.set_xml_output(xml_output) });
}

// Macro to handle the common extraction logic and error wrapping.
macro_rules! perform_extraction {
    (
        $handle:expr,
        $out_ptr1:expr,
        $out_ptr2:expr,
        $extractor_call:expr,
        $success_handler:expr
    ) => {{
        if $handle.is_null() || $out_ptr1.is_null() || $out_ptr2.is_null() {
            return ERR_NULL_POINTER;
        }

        // Safely get a shared reference to the extractor.
        let extractor = unsafe { &*($handle as *const CoreExtractor) };

        match $extractor_call(extractor) {
            Ok((res1, res2)) => {
                $success_handler($out_ptr1, $out_ptr2, res1, res2);
                ERR_OK
            }
            Err(e) => {
                let code = extractous_error_to_code(&e);
                set_last_error(e);
                code
            }
        }
    }};
}

/// Extracts content and metadata from a local file path into a string.
///
/// Output strings must be freed with `extractous_string_free`.
/// Output metadata must be freed with `extractous_metadata_free`.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_extract_file_to_string(
    handle: *mut CExtractor,
    path: *const c_char,
    out_content: *mut *mut c_char,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    if path.is_null() {
        return ERR_NULL_POINTER;
    }
    let path_str = match unsafe { CStr::from_ptr(path).to_str() } {
        Ok(s) => s,
        Err(_) => return ERR_INVALID_UTF8,
    };

    perform_extraction!(
        handle,
        out_content,
        out_metadata,
        |extractor: &CoreExtractor| extractor.extract_file_to_string(path_str),
        |out_c: *mut *mut c_char, out_m: *mut *mut CMetadata, content, metadata| {
            unsafe {
                *out_c = CString::new(content).map_or(ptr::null_mut(), |s| s.into_raw());
                *out_m = metadata_to_c(metadata);
            }
        }
    )
}

/// Extracts content and metadata from a local file path into a stream.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_extract_file(
    handle: *mut CExtractor,
    path: *const c_char,
    out_reader: *mut *mut CStreamReader,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    if path.is_null() {
        return ERR_NULL_POINTER;
    }
    let path_str = match unsafe { CStr::from_ptr(path).to_str() } {
        Ok(s) => s,
        Err(_) => return ERR_INVALID_UTF8,
    };

    perform_extraction!(
        handle,
        out_reader,
        out_metadata,
        |extractor: &CoreExtractor| extractor.extract_file(path_str),
        |out_r: *mut *mut CStreamReader, out_m: *mut *mut CMetadata, reader, metadata| {
            unsafe {
                *out_r = Box::into_raw(Box::new(reader)) as *mut CStreamReader;
                *out_m = metadata_to_c(metadata);
            }
        }
    )
}

/// Extracts content and metadata from a byte slice into a string.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_extract_bytes_to_string(
    handle: *mut CExtractor,
    data: *const u8,
    data_len: libc::size_t,
    out_content: *mut *mut c_char,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    if data.is_null() {
        return ERR_NULL_POINTER;
    }
    let bytes = unsafe { std::slice::from_raw_parts(data, data_len) };

    perform_extraction!(
        handle,
        out_content,
        out_metadata,
        |extractor: &CoreExtractor| extractor.extract_bytes_to_string(bytes),
        |out_c: *mut *mut c_char, out_m: *mut *mut CMetadata, content, metadata| {
            unsafe {
                *out_c = CString::new(content).map_or(ptr::null_mut(), |s| s.into_raw());
                *out_m = metadata_to_c(metadata);
            }
        }
    )
}

/// Extracts content and metadata from a byte slice into a stream.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_extract_bytes(
    handle: *mut CExtractor,
    data: *const u8,
    data_len: libc::size_t,
    out_reader: *mut *mut CStreamReader,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    if data.is_null() {
        return ERR_NULL_POINTER;
    }
    let bytes = unsafe { std::slice::from_raw_parts(data, data_len) };

    perform_extraction!(
        handle,
        out_reader,
        out_metadata,
        |extractor: &CoreExtractor| extractor.extract_bytes(bytes),
        |out_r: *mut *mut CStreamReader, out_m: *mut *mut CMetadata, reader, metadata| {
            unsafe {
                *out_r = Box::into_raw(Box::new(reader)) as *mut CStreamReader;
                *out_m = metadata_to_c(metadata);
            }
        }
    )
}

/// Extracts content and metadata from a URL into a string.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_extract_url_to_string(
    handle: *mut CExtractor,
    url: *const c_char,
    out_content: *mut *mut c_char,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    if url.is_null() {
        return ERR_NULL_POINTER;
    }
    let url_str = match unsafe { CStr::from_ptr(url).to_str() } {
        Ok(s) => s,
        Err(_) => return ERR_INVALID_UTF8,
    };

    perform_extraction!(
        handle,
        out_content,
        out_metadata,
        |extractor: &CoreExtractor| extractor.extract_url_to_string(url_str),
        |out_c: *mut *mut c_char, out_m: *mut *mut CMetadata, content, metadata| {
            unsafe {
                *out_c = CString::new(content).map_or(ptr::null_mut(), |s| s.into_raw());
                *out_m = metadata_to_c(metadata);
            }
        }
    )
}

/// Extracts content and metadata from a URL into a stream.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_extractor_extract_url(
    handle: *mut CExtractor,
    url: *const c_char,
    out_reader: *mut *mut CStreamReader,
    out_metadata: *mut *mut CMetadata,
) -> libc::c_int {
    if url.is_null() {
        return ERR_NULL_POINTER;
    }
    let url_str = match unsafe { CStr::from_ptr(url).to_str() } {
        Ok(s) => s,
        Err(_) => return ERR_INVALID_UTF8,
    };

    perform_extraction!(
        handle,
        out_reader,
        out_metadata,
        |extractor: &CoreExtractor| extractor.extract_url(url_str),
        |out_r: *mut *mut CStreamReader, out_m: *mut *mut CMetadata, reader, metadata| {
            unsafe {
                *out_r = Box::into_raw(Box::new(reader)) as *mut CStreamReader;
                *out_m = metadata_to_c(metadata);
            }
        }
    )
}

/// Frees a C-style string that was allocated by this library.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_string_free(s: *mut c_char) {
    if !s.is_null() {
        drop(unsafe { CString::from_raw(s) });
    }
}
