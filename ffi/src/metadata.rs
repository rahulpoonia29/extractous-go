//! Metadata extraction and manipulation utilities
//!
//! This module provides utilities for working with document metadata extracted
//! during content extraction. Metadata includes information such as:
//! - Author, title, subject
//! - Creation and modification dates
//! - Document format and version
//! - Page count, word count
//! - Custom properties
//!
//! ## Structure
//!
//! Metadata is represented as a `CMetadata` structure containing parallel arrays
//! of keys and values. Each key may have multiple values, which are comma-separated.
//!
//! ## Memory Management
//!
//! Metadata structures are allocated by extraction functions and must be freed
//! by the caller using `extractous_metadata_free()`.

use crate::types::CMetadata;
use std::collections::HashMap;
use std::ffi::{CStr, CString};

/// Convert a Rust HashMap to a C-compatible metadata structure.
///
/// This is an internal utility function. The returned `CMetadata` structure
/// must be freed using `extractous_metadata_free()`.
///
/// ### Safety
/// Allocates memory for keys and values that must be freed by the caller.
pub(crate) unsafe fn metadata_to_c(metadata: HashMap<String, Vec<String>>) -> *mut CMetadata {
    let len = metadata.len();

    if len == 0 {
        // Return empty metadata structure
        return Box::into_raw(Box::new(CMetadata {
            keys: std::ptr::null_mut(),
            values: std::ptr::null_mut(),
            len: 0,
        }));
    }

    let mut keys: Vec<*mut libc::c_char> = Vec::with_capacity(len);
    let mut values: Vec<*mut libc::c_char> = Vec::with_capacity(len);

    for (key, value_vec) in metadata {
        // Convert key to C string
        keys.push(CString::new(key).unwrap().into_raw());

        // Join multiple values with comma separator
        let joined = value_vec.join(",");
        values.push(CString::new(joined).unwrap().into_raw());
    }

    let keys_ptr = keys.as_mut_ptr();
    let values_ptr = values.as_mut_ptr();

    // Prevent Rust from dropping the vectors (we need the pointers)
    std::mem::forget(keys);
    std::mem::forget(values);

    Box::into_raw(Box::new(CMetadata {
        keys: keys_ptr,
        values: values_ptr,
        len,
    }))
}

/// Free a metadata structure and all associated memory.
///
/// Frees:
/// 1. All individual key strings
/// 2. All individual value strings
/// 3. The key array
/// 4. The value array
/// 5. The `CMetadata` structure itself
///
/// ### Arguments
/// * `metadata` - Pointer to a `CMetadata` structure to free
///
/// ### Safety
/// - `metadata` must be a pointer returned by an extraction function
/// - `metadata` must not be used after this call
/// - Safe to call with NULL (no-op)
///
/// ### Example
/// ```c
/// CMetadata* metadata;
/// extractous_extractor_extract_file_to_string(extractor, path, &content, &metadata);
/// // ... use metadata ...
/// extractous_metadata_free(metadata);
/// ```
#[no_mangle]
pub unsafe extern "C" fn extractous_metadata_free(metadata: *mut CMetadata) {
    if metadata.is_null() {
        return;
    }

    let m = Box::from_raw(metadata);

    // Free all key and value strings
    for i in 0..m.len {
        if !m.keys.is_null() {
            let _ = CString::from_raw(*m.keys.add(i));
        }
        if !m.values.is_null() {
            let _ = CString::from_raw(*m.values.add(i));
        }
    }

    // Free the arrays themselves if allocated
    if !m.keys.is_null() && m.len > 0 {
        let _ = Vec::from_raw_parts(m.keys, m.len, m.len);
    }
    if !m.values.is_null() && m.len > 0 {
        let _ = Vec::from_raw_parts(m.values, m.len, m.len);
    }
}
