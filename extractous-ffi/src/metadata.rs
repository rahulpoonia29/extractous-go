use crate::types::CMetadata;
use std::collections::HashMap;
use std::ffi::CString;

/// Convert Rust HashMap to C metadata structure
pub unsafe fn metadata_to_c(metadata: HashMap<String, Vec<String>>) -> *mut CMetadata {
    let len = metadata.len();
    let mut keys: Vec<*mut libc::c_char> = Vec::with_capacity(len);
    let mut values: Vec<*mut libc::c_char> = Vec::with_capacity(len);

    for (key, value_vec) in metadata {
        keys.push(CString::new(key).unwrap().into_raw());
        let joined = value_vec.join(",");
        values.push(CString::new(joined).unwrap().into_raw());
    }

    let keys_ptr = keys.as_mut_ptr();
    let values_ptr = values.as_mut_ptr();
    std::mem::forget(keys);
    std::mem::forget(values);

    Box::into_raw(Box::new(CMetadata {
        keys: keys_ptr,
        values: values_ptr,
        len,
    }))
}

/// Free metadata structure
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_metadata_free(meta: *mut CMetadata) {
    if meta.is_null() {
        return;
    }

    let m = Box::from_raw(meta);

    for i in 0..m.len {
        let _ = CString::from_raw(*m.keys.add(i));
        let _ = CString::from_raw(*m.values.add(i));
    }

    let _ = Vec::from_raw_parts(m.keys, m.len, m.len);
    let _ = Vec::from_raw_parts(m.values, m.len, m.len);
}
