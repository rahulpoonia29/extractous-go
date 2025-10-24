use crate::ecore::Error;
use std::cell::RefCell;
use std::error::Error as StdError;
use std::ffi::CString;
use std::os::raw::{c_char, c_int};
use std::ptr;

pub const ERR_OK: c_int = 0;
pub const ERR_NULL_POINTER: c_int = -1;
pub const ERR_INVALID_UTF8: c_int = -2;
pub const ERR_INVALID_STRING: c_int = -3;
pub const ERR_EXTRACTION_FAILED: c_int = -4;
pub const ERR_IO_ERROR: c_int = -5;
pub const ERR_INVALID_CONFIG: c_int = -6;
pub const ERR_INVALID_ENUM: c_int = -7;
pub const ERR_UNSUPPORTED_FORMAT: c_int = -8;
pub const ERR_OUT_OF_MEMORY: c_int = -9;
pub const ERR_OCR_FAILED: c_int = -10;

pub(crate) fn extractous_error_to_code(err: &Error) -> c_int {
    match err {
        Error::IoError(_) => ERR_IO_ERROR,
        Error::Utf8Error(_) => ERR_INVALID_UTF8,

        // For unknown errors, inspect the message content
        Error::ParseError(msg) | Error::Unknown(msg) => {
            let lower_msg = msg.to_lowercase();
            if lower_msg.contains("ocr") {
                ERR_OCR_FAILED
            } else if lower_msg.contains("unsupported") {
                ERR_UNSUPPORTED_FORMAT
            } else if lower_msg.contains("config") {
                ERR_INVALID_CONFIG
            } else {
                // Default to general extraction failure
                ERR_EXTRACTION_FAILED
            }
        }

        Error::JniError(jni_err) => {
            let error_string = jni_err.to_string();
            let lower_error_string = error_string.to_lowercase();

            if lower_error_string.contains("javaexception") {
                // This string appears when the error is due to a Java-side exception,
                // which is the case your `jnicallmethodlocal` handles. This is a strong
                // indicator of a failure within Tika's processing.
                ERR_EXTRACTION_FAILED
            } else if lower_error_string.contains("nomemory") {
                ERR_OUT_OF_MEMORY
            } else {
                ERR_EXTRACTION_FAILED
            }
        }

        Error::JniEnvCall(_) => ERR_EXTRACTION_FAILED,
    }
}

#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_message(code: c_int) -> *mut c_char {
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
        Err(_) => ptr::null_mut(),
    }
}

thread_local! {
    /// Stores the last detailed error that occurred on the current thread
    static LAST_ERROR: RefCell<Option<Box<dyn StdError + Send>>> = RefCell::new(None);
}

pub(crate) fn set_last_error(err: impl StdError + Send + 'static) {
    LAST_ERROR.with(|cell| {
        *cell.borrow_mut() = Some(Box::new(err));
    });
}

/// Retrieves a detailed debug report for the last error on this thread
/// full error chain and a backtrace if RUST_BACKTRACE=1
#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_get_last_debug() -> *mut c_char {
    LAST_ERROR.with(|cell| {
        if let Some(err) = cell.borrow_mut().take() {
            let mut debug_output = format!("Error: {}", err);
            let mut source = err.source();
            if source.is_some() {
                debug_output.push_str("\n\nCaused by:");
            }
            let mut level = 0;
            while let Some(cause) = source {
                debug_output.push_str(&format!("\n    {}: {}", level, cause));
                source = cause.source();
                level += 1;
            }
            debug_output.push_str(&format!("\n\nDebug Representation:\n{:?}", err));
            match CString::new(debug_output) {
                Ok(s) => s.into_raw(),
                Err(_) => ptr::null_mut(),
            }
        } else {
            ptr::null_mut()
        }
    })
}

/// Checks if debug information is available for the current thread
#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_has_debug() -> c_int {
    LAST_ERROR.with(|cell| if cell.borrow().is_some() { 1 } else { 0 })
}

#[unsafe(no_mangle)]
pub extern "C" fn extractous_error_clear_last() {
    LAST_ERROR.with(|cell| {
        *cell.borrow_mut() = None;
    });
}
