//! Error code definitions and error handling utilities

use std::ffi::CString;
use std::os::raw::c_int;

// ============================================================================
// Error Codes
// ============================================================================

/// Success, no error occurred
pub const ERR_OK: c_int = 0;

/// Null pointer was provided as an argument
pub const ERR_NULL_POINTER: c_int = -1;

/// Invalid UTF-8 string encoding
pub const ERR_INVALID_UTF8: c_int = -2;

/// String conversion failed
pub const ERR_INVALID_STRING: c_int = -3;

/// Document extraction failed
pub const ERR_EXTRACTION_FAILED: c_int = -4;

/// I/O operation failed
pub const ERR_IO_ERROR: c_int = -5;

/// Invalid configuration provided
pub const ERR_INVALID_CONFIG: c_int = -6;

/// Invalid enum value provided
pub const ERR_INVALID_ENUM: c_int = -7;

// ============================================================================
// Error Message Function
// ============================================================================

/// Get human-readable error message for error code
///
/// # Safety
///
/// The returned string must be freed with `extractous_string_free`.
///
/// # Returns
///
/// Pointer to null-terminated C string, or NULL on error.
#[no_mangle]
pub extern "C" fn extractous_error_message(code: c_int) -> *mut libc::c_char {
    let msg = match code {
        ERR_OK => "No error",
        ERR_NULL_POINTER => "Null pointer provided",
        ERR_INVALID_UTF8 => "Invalid UTF-8 string",
        ERR_INVALID_STRING => "String conversion failed",
        ERR_EXTRACTION_FAILED => "Extraction failed",
        ERR_IO_ERROR => "I/O error",
        ERR_INVALID_CONFIG => "Invalid configuration",
        ERR_INVALID_ENUM => "Invalid enum value",
        _ => "Unknown error",
    };

    match CString::new(msg) {
        Ok(s) => s.into_raw(),
        Err(_) => std::ptr::null_mut(),
    }
}
