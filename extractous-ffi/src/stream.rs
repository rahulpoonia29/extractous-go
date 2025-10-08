use crate::ecore::StreamReader as CoreStreamReader;
use crate::errors::*;
use crate::types::*;
use std::io::Read;

/// Read from stream into buffer
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_stream_read(
    handle: *mut CStreamReader,
    buffer: *mut u8,
    buffer_size: libc::size_t,
    bytes_read: *mut libc::size_t,
) -> libc::c_int {
    if handle.is_null() || buffer.is_null() {
        return ERR_NULL_POINTER;
    }

    let reader = &mut *(handle as *mut CoreStreamReader);
    let buf_slice = std::slice::from_raw_parts_mut(buffer, buffer_size);

    match reader.read(buf_slice) {
        Ok(n) => {
            if !bytes_read.is_null() {
                *bytes_read = n;
            }
            ERR_OK
        }
        Err(_) => ERR_IO_ERROR,
    }
}

/// Free stream reader
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_stream_free(handle: *mut CStreamReader) {
    if !handle.is_null() {
        let _ = Box::from_raw(handle as *mut CoreStreamReader);
    }
}
