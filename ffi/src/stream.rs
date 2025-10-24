use crate::ecore::StreamReader as CoreStreamReader;
use crate::errors::*;
use crate::types::*;
use std::io::Read;

/// Reads data from a stream into a user-provided buffer.
///
/// Returns the actual number of bytes read via the `bytes_read` output parameter.
/// Reaching the end of the stream is indicated by `ERR_OK` and `*bytes_read == 0`.
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
    if !bytes_read.is_null() {
        unsafe { *bytes_read = 0 };
    }
    if buffer_size == 0 {
        return ERR_OK;
    }

    let reader = unsafe { &mut *(handle as *mut CoreStreamReader) };
    let buf_slice = unsafe { std::slice::from_raw_parts_mut(buffer, buffer_size) };

    match reader.read(buf_slice) {
        Ok(n) => {
            if !bytes_read.is_null() {
                unsafe { *bytes_read = n };
            }
            ERR_OK
        }
        Err(_) => ERR_IO_ERROR,
    }
}

/// Reads exactly `buffer_size` bytes from the stream.
///
/// Function will continue reading until the buffer is full, or the end of
/// the stream is reached, or an error occurs.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_stream_read_exact(
    handle: *mut CStreamReader,
    buffer: *mut u8,
    buffer_size: libc::size_t,
    bytes_read: *mut libc::size_t,
) -> libc::c_int {
    if handle.is_null() || buffer.is_null() || bytes_read.is_null() {
        return ERR_NULL_POINTER;
    }
    if buffer_size == 0 {
        return ERR_OK;
    }

    unsafe { *bytes_read = 0 };

    let reader = unsafe { &mut *(handle as *mut CoreStreamReader) };
    // slice representing the user-provided buffer
    let total_buf_slice = unsafe { std::slice::from_raw_parts_mut(buffer, buffer_size) };

    let mut total_bytes_read = 0;
    while total_bytes_read < buffer_size {
        // In each loop, we try to read into the remaining part of the buffer
        let remaining_buf = &mut total_buf_slice[total_bytes_read..];

        match reader.read(remaining_buf) {
            Ok(0) => {
                // `read` returned 0, which signifies the end of the stream
                // We break the loop and will return the total bytes we've read
                break;
            }
            Ok(n) => {
                total_bytes_read += n;
            }
            Err(e) if e.kind() == std::io::ErrorKind::Interrupted => {
                // The read was interrupted by a signal. This is recoverable so we just continue
                continue;
            }
            Err(_) => {
                // A non-recoverable I/O error occurred.
                return ERR_IO_ERROR;
            }
        }
    }

    unsafe { *bytes_read = total_bytes_read };
    ERR_OK
}

/// Reads the remaining stream into a newly allocated buffer.
#[must_use]
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_stream_read_all(
    handle: *mut CStreamReader,
    out_buffer: *mut *mut u8,
    out_size: *mut libc::size_t,
) -> libc::c_int {
    if handle.is_null() || out_buffer.is_null() || out_size.is_null() {
        return ERR_NULL_POINTER;
    }

    let reader = unsafe { &mut *(handle as *mut CoreStreamReader) };
    let mut data_vec = Vec::new();

    match reader.read_to_end(&mut data_vec) {
        Ok(_) => {
            data_vec.shrink_to_fit();

            let size = data_vec.len();
            let ptr = data_vec.as_mut_ptr();
            std::mem::forget(data_vec);

            unsafe { *out_buffer = ptr };
            unsafe { *out_size = size };
            ERR_OK
        }
        Err(_) => ERR_IO_ERROR,
    }
}

/// Frees a buffer allocated by `extractous_stream_read_all`.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_buffer_free(buffer: *mut u8, size: libc::size_t) {
    if buffer.is_null() || size == 0 {
        return;
    }
    let _ = unsafe { Vec::from_raw_parts(buffer, size, size) };
}

/// Frees a stream reader and releases its resources.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn extractous_stream_free(handle: *mut CStreamReader) {
    if !handle.is_null() {
        // Reconstruct the Box and let Rust's drop handler deallocate it.
        let _ = unsafe { Box::from_raw(handle as *mut CoreStreamReader) };
    }
}
