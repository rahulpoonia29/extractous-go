//! Error code definitions and error handling utilities
//!
//! This module provides a error handling system for the FFI interface.
//! All errors are represented as integer codes for C compatibility, with
//! human-readable messages available through helper functions.
//!
//! ## Error Philosophy
//!
//! - **Positive values**: Reserved for future use (warnings, info)
//! - **Zero (0)**: Success (`ERR_OK`)
//! - **Negative values**: Error conditions
//!
//! ## Usage Pattern
//!
//! ```c
//! int result = extractous_extractor_extract_file_to_string(
//!     extractor, path, &content, &metadata
//! );
//!
//! if (result != ER R_OK) {
//!     char* error_msg = extractous_error_message(result);
//!     fprintf(stderr, "Extraction failed: %s\n", error_msg);
//!     extractous_string_free(error_msg);
//!     return result;
//! }
//! ```

use std::cell::RefCell;
use std::error::Error as StdError;
use std::ffi::CString;
use std::os::raw::c_int;
// use std::ffi::{CStr, CString};

// ============================================================================
// Error Code Constants
// ============================================================================

/// Success - operation completed without errors
///
/// This is the only non-error return value. All operations that complete
/// successfully will return this code.
pub const ERR_OK: c_int = 0;

/// Null pointer provided as argument
///
/// Returned when a required pointer argument is NULL.
/// Check all pointer arguments before calling FFI functions.
///
/// Common causes:
/// - Forgot to allocate output parameter
/// - Accidentally passed NULL instead of valid pointer
/// - Double-free caused pointer to become invalid
pub const ERR_NULL_POINTER: c_int = -1;

/// Invalid UTF-8 string encoding
///
/// Returned when a C string argument contains invalid UTF-8 sequences.
/// All string arguments must be valid UTF-8.
///
/// Common causes:
/// - String contains binary data
/// - Wrong encoding used (e.g., Latin-1 instead of UTF-8)
/// - Corrupted string data
pub const ERR_INVALID_UTF8: c_int = -2;

/// String conversion or allocation failed
///
/// Returned when internal string operations fail, typically due to:
/// - Null bytes in unexpected positions
/// - Memory allocation failure
/// - String contains invalid characters for the operation
pub const ERR_INVALID_STRING: c_int = -3;

/// Document extraction failed
///
/// General extraction error when the specific cause is unknown or internal.
/// The document may be:
/// - Corrupted or malformed
/// - Using an unsupported format variant
/// - Encrypted without proper credentials
/// - Too complex for the parser
pub const ERR_EXTRACTION_FAILED: c_int = -4;

/// File system or network I/O error
///
/// Returned when file or network operations fail.
///
/// Common causes:
/// - File not found
/// - Permission denied
/// - Network timeout
/// - Disk full
/// - Path too long
pub const ERR_IO_ERROR: c_int = -5;

/// Invalid configuration value
///
/// Returned when configuration parameters are invalid.
///
/// Common causes:
/// - Out of range values
/// - Incompatible configuration combinations
/// - Invalid enum constants
pub const ERR_INVALID_CONFIG: c_int = -6;

/// Invalid enumeration value
///
/// Returned when an enum constant (like charset or OCR strategy) is invalid.
/// Only use the documented constants.
pub const ERR_INVALID_ENUM: c_int = -7;

/// Unsupported file format
///
/// The file format is not supported by extractous or the parser
/// for this format is not available.
pub const ERR_UNSUPPORTED_FORMAT: c_int = -8;

/// Memory allocation failed
///
/// Extremely rare - indicates the system is out of memory.
pub const ERR_OUT_OF_MEMORY: c_int = -9;

/// OCR operation failed
///
/// OCR processing failed, possibly because:
/// - Tesseract is not installed
/// - Invalid language data
/// - Image format not supported
pub const ERR_OCR_FAILED: c_int = -10;

// ============================================================================
// Error Message Functions
// ============================================================================

/// Get human-readable error message for error code
///
/// Returns a newly allocated string containing the error description.
/// The caller must free the returned string with `extractous_string_free()`.
///
/// ### Arguments
/// * `code` - Error code returned by an extractous function
///
/// ### Returns
/// Pointer to null-terminated UTF-8 string, or NULL if allocation fails.
///
/// ### Example
/// ```c
/// int err = extractous_extractor_extract_file(...);
/// if (err != ERR_OK) {
///     char* msg = extractous_error_message(err);
///     printf("Error: %s\n", msg);
///     extractous_string_free(msg);
/// }
/// ```
///
/// ### Safety
/// - Return value must be freed with `extractous_string_free()`
/// - Do not modify the returned string
#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_message(code: c_int) -> *mut libc::c_char {
    let msg = match code {
        ERR_OK => "Operation completed successfully",
        ERR_NULL_POINTER => "Null pointer provided as argument",
        ERR_INVALID_UTF8 => "Invalid UTF-8 string encoding",
        ERR_INVALID_STRING => "String conversion or allocation failed",
        ERR_EXTRACTION_FAILED => "Document extraction failed",
        ERR_IO_ERROR => "File system or network I/O error",
        ERR_INVALID_CONFIG => "Invalid configuration value",
        ERR_INVALID_ENUM => "Invalid enumeration value",
        ERR_UNSUPPORTED_FORMAT => "Unsupported file format",
        ERR_OUT_OF_MEMORY => "Memory allocation failed",
        ERR_OCR_FAILED => "OCR operation failed",
        _ => "Unknown error code",
    };

    match CString::new(msg) {
        Ok(s) => s.into_raw(),
        Err(_) => std::ptr::null_mut(),
    }
}

/// Get error category description
///
/// Returns a high-level categorization of the error.
///
/// ### Arguments
/// * `code` - Error code
///
/// ### Returns
/// Static string pointer (do not free)
///
/// ### Safety
/// Return value points to static memory and must not be freed.
#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_category(code: c_int) -> *const libc::c_char {
    let category = match code {
        ERR_OK => "success\0",
        ERR_NULL_POINTER | ERR_INVALID_UTF8 | ERR_INVALID_STRING | ERR_INVALID_CONFIG
        | ERR_INVALID_ENUM => "invalid_argument\0",
        ERR_IO_ERROR => "io_error\0",
        ERR_EXTRACTION_FAILED | ERR_UNSUPPORTED_FORMAT | ERR_OCR_FAILED => "extraction_error\0",
        ERR_OUT_OF_MEMORY => "resource_error\0",
        _ => "unknown\0",
    };

    category.as_ptr() as *const libc::c_char
}

// ============================================================================
// Internal Error Conversion Utilities
// ============================================================================

/// Convert extractous core errors to FFI error codes
///
/// This is an internal utility function that maps Rust errors from the
/// extractous core library to C-compatible error codes.
///
/// ### Arguments
/// * `err` - Reference to an extractous error
///
/// ### Returns
/// Appropriate error code constant
pub(crate) fn extractous_error_to_code(err: &crate::ecore::Error) -> c_int {
    use std::io::ErrorKind;

    // Get error message for better classification
    let err_msg = err.to_string().to_lowercase();

    // Check for specific error patterns in the message
    if err_msg.contains("unsupported") || err_msg.contains("unknown format") {
        return ERR_UNSUPPORTED_FORMAT;
    }

    if err_msg.contains("ocr") || err_msg.contains("tesseract") {
        return ERR_OCR_FAILED;
    }

    if err_msg.contains("memory") || err_msg.contains("allocation") {
        return ERR_OUT_OF_MEMORY;
    }

    // Walk the error's source chain to find the underlying cause
    let mut source = err.source();
    while let Some(cause) = source {
        // Check if the cause is a standard I/O error
        if let Some(io_err) = cause.downcast_ref::<std::io::Error>() {
            return match io_err.kind() {
                ErrorKind::NotFound => ERR_IO_ERROR,
                ErrorKind::PermissionDenied => ERR_IO_ERROR,
                ErrorKind::ConnectionRefused => ERR_IO_ERROR,
                ErrorKind::ConnectionReset => ERR_IO_ERROR,
                ErrorKind::TimedOut => ERR_IO_ERROR,
                ErrorKind::WriteZero => ERR_IO_ERROR,
                ErrorKind::Interrupted => ERR_IO_ERROR,
                ErrorKind::UnexpectedEof => ERR_IO_ERROR,
                ErrorKind::InvalidData => ERR_EXTRACTION_FAILED,
                _ => ERR_IO_ERROR,
            };
        }
        source = cause.source();
    }

    // If no specific error was identified, return general extraction failure
    ERR_EXTRACTION_FAILED
}

thread_local! {
    /// Stores the last error that occurred on this thread
    /// We store the actual Error object, not the formatted string,
    /// to defer the expensive formatting until it's requested
    static LAST_ERROR: RefCell<Option<Box<dyn std::error::Error + Send>>> = RefCell::new(None);
}

/// Store the last error in thread-local storage
///
/// This is called internally whenever an FFI function returns an error code.
/// The error is stored as-is without any formatting, making this very cheap.
pub(crate) fn set_last_error(err: impl std::error::Error + Send + 'static) {
    LAST_ERROR.with(|cell| {
        *cell.borrow_mut() = Some(Box::new(err));
    });
}

/// Retrieve detailed debug information for the last error on this thread
///
/// This function formats the stored error with full debug representation
/// including error chain and backtrace (if RUST_BACKTRACE=1).
///
/// **Important**: This clears the stored error after retrieval.
/// Subsequent calls return NULL unless a new error occurs.
///
/// # Returns
/// - Pointer to null-terminated UTF-8 string with debug info
/// - NULL if no error has occurred on this thread
///
/// # Safety
/// - Returned string must be freed with `extractous_string_free()`
/// - This function is thread-safe (uses thread-local storage)
///
/// # Performance
/// This function does expensive string formatting. Only call it when
/// you actually need debug information (e.g., logging, debugging).
///
/// # Example
/// ```
/// int code = extractous_extractor_extract_file_to_string(...);
/// if (code != ERR_OK) {
///     // Get user-facing message (fast)
///     char* msg = extractous_error_message(code);
///     printf("Error: %s\n", msg);
///     extractous_string_free(msg);
///     
///     // Optionally get debug info (slower, for developers)
///     char* debug = extractous_error_get_last_debug();
///     if (debug) {
///         fprintf(stderr, "Debug details:\n%s\n", debug);
///         extractous_string_free(debug);
///     }
/// }
/// ```
#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_get_last_debug() -> *mut libc::c_char {
    LAST_ERROR.with(|cell| {
        if let Some(err) = cell.borrow_mut().take() {
            // Format the error with full debug representation
            // This includes:
            // - Main error message
            // - Complete source chain (all nested errors)
            // - Backtrace if RUST_BACKTRACE=1 or RUST_BACKTRACE=full
            let mut debug_output = format!("Error: {}\n", err);

            // Walk the error chain
            let mut source = err.source();
            let mut level = 0;
            while let Some(cause) = source {
                debug_output.push_str(&format!("\nCaused by:\n  {}: {}", level, cause));
                source = cause.source();
                level += 1;
            }

            // Add backtrace if available (requires RUST_BACKTRACE=1)
            // The Debug formatter will include it
            debug_output.push_str(&format!("\n\nDebug representation:\n{:?}", err));

            // Convert to C string
            match CString::new(debug_output) {
                Ok(s) => s.into_raw(),
                Err(_) => std::ptr::null_mut(),
            }
        } else {
            std::ptr::null_mut()
        }
    })
}

/// Check if there is a stored error for the current thread
///
/// This is useful to check if debug info is available before
/// calling the more expensive `extractous_error_get_last_debug()`.
///
/// # Returns
/// - 1 if an error is stored
/// - 0 if no error is stored
#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_has_debug() -> libc::c_int {
    LAST_ERROR.with(|cell| if cell.borrow().is_some() { 1 } else { 0 })
}

/// Clear any stored error for the current thread without retrieving it
///
/// This is useful to reset error state without the overhead of formatting.
#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_clear_last() {
    LAST_ERROR.with(|cell| {
        *cell.borrow_mut() = None;
    });
}
