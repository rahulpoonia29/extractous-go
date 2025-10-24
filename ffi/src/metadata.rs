use crate::types::CMetadata;
use std::collections::HashMap;
use std::ffi::CString;
use std::os::raw::c_char;
use std::ptr;

/// Convert a Rust HashMap to a C-compatible metadata structure.
pub(crate) unsafe fn metadata_to_c(metadata: HashMap<String, Vec<String>>) -> *mut CMetadata {
    if metadata.is_empty() {
        return Box::into_raw(Box::new(CMetadata {
            keys: ptr::null_mut(),
            values: ptr::null_mut(),
            len: 0,
        }));
    }

    let capacity = metadata.len();
    let mut keys: Vec<*mut c_char> = Vec::with_capacity(capacity);
    let mut values: Vec<*mut c_char> = Vec::with_capacity(capacity);

    for (key, value_vec) in metadata {
        // CString::new will return an error if the string contains `\0`.
        let c_key = match CString::new(key) {
            Ok(s) => s.into_raw(),
            Err(_) => continue, // Skip metadata with invalid keys.
        };

        let joined_values = value_vec.join(", ");
        let c_value = match CString::new(joined_values) {
            Ok(s) => s.into_raw(),
            Err(_) => {
                // Clean up the already-allocated key if the value is invalid.
                let _ = unsafe { CString::from_raw(c_key) };
                continue;
            }
        };

        keys.push(c_key);
        values.push(c_value);
    }

    // Final length is derived from the vectors after they are populated.
    // Guarantees that the length matches the number of allocated pointers.
    let final_len = keys.len();
    assert_eq!(final_len, values.len());

    if final_len == 0 {
        return Box::into_raw(Box::new(CMetadata {
            keys: ptr::null_mut(),
            values: ptr::null_mut(),
            len: 0,
        }));
    }

    keys.shrink_to_fit();
    values.shrink_to_fit();

    let keys_ptr = keys.as_mut_ptr();
    let values_ptr = values.as_mut_ptr();
    std::mem::forget(keys);
    std::mem::forget(values);

    Box::into_raw(Box::new(CMetadata {
        keys: keys_ptr,
        values: values_ptr,
        len: final_len,
    }))
}

/// Frees a metadata structure and all associated memory.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_metadata_free(metadata: *mut CMetadata) {
    if metadata.is_null() {
        return;
    }

    // Take ownership of CMetadata struct.
    let m = unsafe { Box::from_raw(metadata) };

    let keys_vec = unsafe { Vec::from_raw_parts(m.keys, m.len, m.len) };
    let values_vec = unsafe { Vec::from_raw_parts(m.values, m.len, m.len) };

    // Drop to free the memory for each CString.
    for key_ptr in keys_vec {
        let _ = unsafe { CString::from_raw(key_ptr) };
    }

    for value_ptr in values_vec {
        let _ = unsafe { CString::from_raw(value_ptr) };
    }
}
