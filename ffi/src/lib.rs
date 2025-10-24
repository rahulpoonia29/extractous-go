//! This crate provides a **C-compatible Foreign Function Interface (FFI)** for the
//! Extractous library. Extractous is a fast and efficient solution for extracting
//! content and metadata from various document formats including PDF, Word, Excel, and more.
//!
//! This FFI layer is meticulously designed for safety and performance, featuring:
//! - Opaque pointers to prevent unsafe access to internal data structures.
//! - A robust, thread-safe error handling mechanism with on-demand debug info.
//! - A clear memory ownership model with explicit `_new` and `_free` functions.
//!
//! ## Quick Start
//!
//! ```
//! // 1. Create an extractor instance.
//! CExtractor* extractor = extractous_extractor_new();
//!
//! // 2. Configure the extractor. Setters modify the object in-place.
//! //    DO NOT re-assign the pointer.
//! extractous_extractor_set_xml_output(extractor, true);
//!
//! // 3. Extract content and metadata from a file.
//! char* content = NULL;
//! CMetadata* metadata = NULL;
//! int result = extractous_extractor_extract_file_to_string(
//!     extractor, "document.pdf", &content, &metadata
//! );
//!
//! // 4. Check for errors and handle them.
//! if (result != ERR_OK) {
//!     // Handle the error (see Error Handling section).
//!     fprintf(stderr, "Extraction failed with code: %d\n", result);
//! } else {
//!     // 5. Use the results.
//!     printf("Content: %s\n", content);
//! }
//!
//! // 6. Clean up all allocated resources in reverse order.
//! extractous_string_free(content);
//! extractous_metadata_free(metadata);
//! extractous_extractor_free(extractor);
//! ```
//!
//! ## Thread Safety
//!
//! - **Extractor Instances**: `CExtractor` and its associated config/stream objects are
//!   **NOT thread-safe**. Do not share a handle across threads. The recommended pattern is
//!   to create one `CExtractor` instance per thread that needs it.
//! - **Error Handling**: The error reporting system **IS thread-safe**. Each thread stores
//!   its own last error information independently, preventing race conditions. You can safely
//!   call error-handling functions from any thread.
//!
//! # Advanced Error Handling
//!
//! This library uses a powerful two-tier error system for maximum performance and diagnostics.
//!
//! ### Tier 1: Fast Path (Error Codes)
//!
//! All FFI functions return an integer error code. `ERR_OK` (0) signifies success. This allows
//! for a very fast check without any overhead.
//!
//! ### Tier 2: Slow Path (On-Demand Detailed Info)
//!
//! When an error occurs, you can request more information on demand.
//!
//! **1. Get the Error Category:**
//! Use `extractous_error_category()` to get a stable, machine-readable string
//! representing the *type* of error. This is perfect for building idiomatic Go error wrappers.
//! The returned pointer is static and **must not be freed**.
//!
//! **2. Get a Simple Message:**
//! Use `extractous_error_message()` to get a simple, human-readable description.
//! The returned string **must be freed** with `extractous_string_free()`.
//!
//! **3. Get a Full Debug Report:**
//! If `extractous_error_has_debug()` returns `1`, you can call `extractous_error_get_last_debug()`
//! to get a detailed report, including the full error chain and a backtrace (if enabled with `RUST_BACKTRACE=1`).
//! The returned string **must be freed**.
//!
//! ### Go Usage Pattern
//!
//! ```
//! // (Inside a function that calls the FFI)
//! resultCode := C.some_extractous_function(...)
//! if resultCode != C.ERR_OK {
//!     // Get stable category for idiomatic error wrapping.
//!     category := C.GoString(C.extractous_error_category(resultCode))
//!
//!     // Get the simple message for the error string.
//!     msgCStr := C.extractous_error_message(resultCode)
//!     defer C.extractous_string_free(msgCStr)
//!     message := C.GoString(msgCStr)
//!
//!     var baseError error
//!     switch category {
//!     case "io_error": baseError = ErrIO
//!     default: baseError = ErrUnknown
//!     }
//!
//!     // Optionally log the full debug info for developers.
//!     if C.extractous_error_has_debug() != 0 {
//!         debugCStr := C.extractous_error_get_last_debug()
//!         defer C.extractous_string_free(debugCStr)
//!         log.Printf("Full debug details: %s", C.GoString(debugCStr))
//!     }
//!
//!     return fmt.Errorf("%w: %s", baseError, message)
//! }
//! ```
#![warn(clippy::all)]
#![allow(clippy::missing_safety_doc)]

// Re-export the core library under a consistent, private alias.
pub use extractous as ecore;

// Module declarations.
mod config;
mod errors;
mod extractor;
mod metadata;
mod stream;
mod types;

// Publicly re-export all FFI-safe functions and types for C header generation.
pub use config::*;
pub use errors::*;
pub use extractor::*;
pub use metadata::*;
pub use stream::*;
pub use types::*;

/// Returns the FFI wrapper version as a null-terminated UTF-8 string.
/// The returned pointer is to a static string and must not be freed.
#[unsafe(no_mangle)]
pub extern "C" fn extractous_ffi_version() -> *const libc::c_char {
    // Use a static byte array with a null terminator for guaranteed memory safety.
    static VERSION: &[u8] = concat!(env!("CARGO_PKG_VERSION"), "\0").as_bytes();
    VERSION.as_ptr() as *const libc::c_char
}

/// Returns the underlying Extractous core library version.
/// The returned pointer is to a static string and must not be freed.
#[unsafe(no_mangle)]
pub extern "C" fn extractous_core_version() -> *const libc::c_char {
    static VERSION: &[u8] = b"0.3.0\0";
    VERSION.as_ptr() as *const libc::c_char
}
