//! # Extractous FFI
//!
//! C FFI bindings for extractous Go integration.
//!
//! This crate provides C-compatible FFI functions that wrap the extractous
//! Rust library for use in Go via cgo.
//!
//! ## Safety
//!
//! All public FFI functions are marked as `unsafe` or `extern "C"` where appropriate.
//! Callers must ensure:
//! - Pointers are valid and properly aligned
//! - String pointers point to valid null-terminated C strings
//! - Objects are freed using the correct free function
//! - No use-after-free by calling free functions multiple times

// Re-export extractous as ecore for internal use
pub use extractous as ecore;

// Module declarations
mod config;
mod errors;
mod extractor;
mod metadata;
mod stream;
mod types;

// Re-export for C header generation
pub use errors::*;
pub use types::*;
