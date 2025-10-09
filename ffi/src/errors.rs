//! Error code definitions and error handling utilities

use std::error::Error as StdError;
use std::ffi::CString;
use std::os::raw::c_int;

// ============================================================================
// Error Codes
// ============================================================================
pub const ERR_OK: c_int = 0;
pub const ERR_NULL_POINTER: c_int = -1;
pub const ERR_INVALID_UTF8: c_int = -2;
pub const ERR_INVALID_STRING: c_int = -3;
pub const ERR_EXTRACTION_FAILED: c_int = -4;
pub const ERR_IO_ERROR: c_int = -5;
pub const ERR_INVALID_CONFIG: c_int = -6;
pub const ERR_INVALID_ENUM: c_int = -7;

// ============================================================================
// Error Message Function
// ============================================================================
#[no_mangle]
pub extern "C" fn extractous_error_message(code: c_int) -> *mut libc::c_char {
    let msg = match code {
        ERR_OK => "No error",
        ERR_NULL_POINTER => "Null pointer provided",
        ERR_INVALID_UTF8 => "Invalid UTF-8 string",
        ERR_INVALID_STRING => "String conversion failed",
        ERR_EXTRACTION_FAILED => "Extraction failed",
        ERR_IO_ERROR => "IO error",
        ERR_INVALID_CONFIG => "Invalid configuration",
        ERR_INVALID_ENUM => "Invalid enum value",
        _ => "Unknown error",
    };

    match CString::new(msg) {
        Ok(s) => s.into_raw(),
        Err(_) => std::ptr::null_mut(),
    }
}

// ============================================================================
// Error Conversion Function
// ============================================================================

/// Helper to convert extractous errors to error codes
pub(crate) fn extractous_error_to_code(err: &crate::ecore::Error) -> c_int {
    use std::io::ErrorKind;

    // Walk the error's source chain to find the underlying cause.
    let mut source = err.source();
    while let Some(cause) = source {
        // Check if the cause is a standard I/O error.
        if let Some(io_err) = cause.downcast_ref::<std::io::Error>() {
            return match io_err.kind() {
                ErrorKind::NotFound => ERR_IO_ERROR,
                ErrorKind::PermissionDenied => ERR_IO_ERROR,
                // You can add more specific IO errors here if needed
                _ => ERR_IO_ERROR,
            };
        }
        source = cause.source();
    }

    // If no specific I/O error was found, return a general extraction failure.
    ERR_EXTRACTION_FAILED
}
