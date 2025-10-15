//! # Extractous FFI
//!
//! High-performance C FFI for the Extractous document extraction library, safe and easy
//! to use from C, Go (via cgo), or Rust.
//!
//! ## Overview
//!
//! This crate provides a **C-compatible Foreign Function Interface (FFI)** for the
//! Extractous library, enabling seamless integration with external programs. Extractous
//! is a fast and efficient solution for extracting content and metadata from various
//! document formats including PDF, Word, Excel, HTML, and more.
//!
//! ## Quick Start
//!
//! ```c
//! // Create extractor
//! CExtractor* extractor = extractous_extractor_new();
//!
//! // Configure extractor (builder pattern)
//! extractor = extractous_extractor_set_xml_output(extractor, true);
//!
//! // Extract content and metadata
//! char* content;
//! CMetadata* metadata;
//! int result = extractous_extractor_extract_file_to_string(
//!     extractor, "document.pdf", &content, &metadata
//! );
//!
//! // Use results...
//! printf("Content: %s\n", content);
//!
//! //Clean up in reverse allocation order
//! extractous_string_free(content);
//! extractous_metadata_free(metadata);
//! extractous_extractor_free(extractor);
//! ```
//!
//! ## Modules
//!
//! | Module     | Purpose                                      |
//! |------------|----------------------------------------------|
//! | extractor  | File, URL, and byte-array extraction         |
//! | stream     | Buffered reading for large documents         |
//! | metadata   | Access and manipulate document metadata      |
//! | config     | PDF, Office, and OCR configuration settings  |
//! | types      | Type definitions, constants, opaque handles  |
//! | errors     | Error codes and human-readable messages      |
//!
//! ## Streams
//!
//! Streams allow reading large documents **without loading all content into memory**.
//!
//! - `extractous_stream_read` → read chunks of data into a buffer
//! - `extractous_stream_read_exact` → read a fixed number of bytes
//! - `extractous_stream_read_all` → read remaining content into a newly allocated buffer
//! - `extractous_buffer_free` → free memory allocated by `read_all`
//!
//! **Tip:** Use streams for large files. For small files, `read_all` is convenient.
//!
//! ## Metadata
//!
//! Metadata is represented by `CMetadata` containing parallel arrays of keys and values.
//! - Keys may have multiple values, comma-separated.
//! - All allocated metadata must be freed with `extractous_metadata_free()`.
//!
//! ## Memory Safety
//!
//! - **Pointer Validity:** All pointers must be valid and properly aligned.
//! - **String Encoding:** Null-terminated UTF-8 strings are required.
//! - **Memory Ownership:** Use the correct `free` function for all allocated objects.
//! - **No Use-After-Free:** Do not use pointers after freeing.
//! - **Thread Safety:** Extractor instances are not thread-safe; use separate instances per thread.
//!
//! # Error Handling
//!
//! This module defines error codes and provides utilities for FFI-safe error handling.
//! All errors are represented as integer codes for C compatibility. Human-readable
//! messages can be retrieved with `extractous_error_message()`.
//!
//! ## Error Codes
//!
//! - `ERR_OK` (0) → Success
//! - Negative values → Error conditions
//!
//! | Code | Meaning                          |
//! |------|----------------------------------|
//! |  0   | Operation successful             |
//! | -1   | Null pointer provided            |
//! | -2   | Invalid UTF-8 string             |
//! | -3   | String conversion/allocation failed |
//! | -4   | Document extraction failed       |
//! | -5   | File system or network I/O error |
//! | -6   | Invalid configuration value      |
//! | -7   | Invalid enumeration value        |
//! | -8   | Unsupported file format          |
//! | -9   | Memory allocation failed         |
//! | -10  | OCR operation failed             |
//!
//! ### Usage Pattern
//!
//! ```c
//! int result = extractous_extractor_extract_file_to_string(
//!     extractor, path, &content, &metadata
//! );
//!
//! if (result != ERR_OK) {
//!     char* error_msg = extractous_error_message(result);
//!     fprintf(stderr, "Extraction failed: %s\n", error_msg);
//!     extractous_string_free(error_msg);
//!     return result;
//! }
//! ```
//!
//! - Returned strings from `extractous_error_message()` must be freed with
//!   `extractous_string_free()`
//! - Do not modify the returned string
//! - Error codes are stable and can be used for programmatic handling
//!
//! ## Version Information
//!
//! FFI Version: 0.1.0  
//! Extractous Core Version: 0.3.0
//!
//! ## Platform Support
//!
//! - Linux (x86_64, aarch64)  
//! - macOS (x86_64, aarch64)  
//! - Windows (x86_64)
//!
//! ## License
//!
//! Apache License 2.0

#![deny(missing_docs)]
#![warn(clippy::all)]
#![allow(clippy::missing_safety_doc)] // Safety docs are in function comments

// Re-export Extractous core library under a consistent alias
pub use extractous as ecore;

// Module declarations
mod config;
mod errors;
mod extractor;
mod metadata;
mod stream;
mod types;

// Public re-exports for C header generation
pub use config::*;
pub use errors::*;
pub use extractor::*;
pub use metadata::*;
pub use stream::*;
pub use types::*;

/// Returns the FFI wrapper version in semver format.
#[no_mangle]
pub extern "C" fn extractous_ffi_version() -> *const libc::c_char {
    concat!(env!("CARGO_PKG_VERSION"), "\0").as_ptr() as *const libc::c_char
}

/// Returns the underlying Extractous core library version.
#[no_mangle]
pub extern "C" fn extractous_core_version() -> *const libc::c_char {
    concat!("0.3.0", "\0").as_ptr() as *const libc::c_char
}
