//! Stream reader FFI functions for efficient buffered content reading
//!
//! This module provides streaming interfaces for reading extracted content.
//! Streams are useful for:
//! - Processing large documents without loading all content into memory
//! - Implementing custom buffering strategies
//! - Integrating with existing streaming APIs
//!
//! ## Usage Pattern
//!
//! ```c
//! CStreamReader* reader;
//! CMetadata* metadata;
//! int result = extractous_extractor_extract_file(
//!     extractor, "large.pdf", &reader, &metadata
//! );
//!
//! if (result == ERR_OK) {
//!     char buffer[4096];
//!     size_t bytes_read;
//!
//!     while (
//!         extractous_stream_read(
//!             reader,
//!             buffer,
//!             sizeof(buffer),
//!             &bytes_read
//!         ) == ERR_OK
//!         && bytes_read > 0
//!     ) {
//!         // Process buffer data (bytes_read bytes valid)
//!         fwrite(buffer, 1, bytes_read, output_file);
//!     }
//!
//!     extractous_stream_free(reader);
//! }
//!
//! extractous_metadata_free(metadata);
//! ```
//!
//! ## Performance Tips
//!
//! - Use buffer sizes of 4KB-64KB for optimal performance
//! - Reuse the same buffer across multiple read calls
//! - Check for ERR_OK and bytes_read > 0 to detect end of stream
//! - Free the reader when done to release resources
use crate::ecore::StreamReader as CoreStreamReader;
use crate::errors::*;
use crate::types::*;
use std::io::Read;

/// Read data from stream into buffer
///
/// Reads up to `buffer_size` bytes from the stream into the provided buffer.
/// Returns the actual number of bytes read in `bytes_read`.
///
/// ### Arguments
/// * `handle` - Valid StreamReader pointer
/// * `buffer` - Pointer to buffer where data will be written
/// * `buffer_size` - Size of the buffer in bytes
/// * `bytes_read` - Output pointer for number of bytes actually read (can be NULL)
///
/// ### Returns
/// * `ERR_OK` - Read successful (check bytes_read for amount)
/// * `ERR_NULL_POINTER` - handle or buffer is NULL
/// * `ERR_IO_ERROR` - Read operation failed
///
/// ### End of Stream
/// When end of stream is reached:
/// - Function returns `ERR_OK`
/// - `bytes_read` is set to 0
///
/// ### Example
/// ```c
/// char buffer[4096];
/// size_t n;
/// int result = extractous_stream_read(reader, buffer, sizeof(buffer), &n);
///
/// if (result == ERR_OK) {
///     if (n > 0) {
///         // Process n bytes of data in buffer
///     } else {
///         // End of stream reached
///     }
/// } else {
///     // Error occurred
/// }
/// ```
///
/// ### Safety
/// - `handle` must be a valid StreamReader pointer
/// - `buffer` must point to at least `buffer_size` writable bytes
/// - `bytes_read` must be NULL or point to valid size_t
/// - Buffer content is undefined if function returns error
#[no_mangle]
pub unsafe extern "C" fn extractous_stream_read(
    handle: *mut CStreamReader,
    buffer: *mut u8,
    buffer_size: libc::size_t,
    bytes_read: *mut libc::size_t,
) -> libc::c_int {
    if handle.is_null() || buffer.is_null() {
        return ERR_NULL_POINTER;
    }

    if buffer_size == 0 {
        if !bytes_read.is_null() {
            unsafe {
                *bytes_read = 0;
            }
        }
        return ERR_OK;
    }

    unsafe {
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
}

/// Read exactly the requested number of bytes or fail
///
/// Similar to `extractous_stream_read()` but guarantees that exactly
/// `buffer_size` bytes are read unless end of stream or error occurs.
///
/// ### Arguments
/// * `handle` - Valid StreamReader pointer
/// * `buffer` - Pointer to buffer where data will be written
/// * `buffer_size` - Exact number of bytes to read
/// * `bytes_read` - Output pointer for number of bytes actually read (can be NULL)
///
/// ### Returns
/// * `ERR_OK` - Successfully read exactly `buffer_size` bytes
/// * `ERR_OK` with bytes_read < buffer_size - End of stream reached
/// * `ERR_NULL_POINTER` - handle or buffer is NULL
/// * `ERR_IO_ERROR` - Read operation failed
///
/// ### Use Cases
/// - Reading fixed-size headers or chunks
/// - When partial reads are not acceptable
///
/// ### Safety
/// - Same safety requirements as `extractous_stream_read()`
#[no_mangle]
pub unsafe extern "C" fn extractous_stream_read_exact(
    handle: *mut CStreamReader,
    buffer: *mut u8,
    buffer_size: libc::size_t,
    bytes_read: *mut libc::size_t,
) -> libc::c_int {
    if handle.is_null() || buffer.is_null() {
        return ERR_NULL_POINTER;
    }

    if buffer_size == 0 {
        if !bytes_read.is_null() {
            unsafe {
                *bytes_read = 0;
            }
        }
        return ERR_OK;
    }

    unsafe {
        let reader = &mut *(handle as *mut CoreStreamReader);
        let buf_slice = std::slice::from_raw_parts_mut(buffer, buffer_size);

        let mut total_read = 0;
        while total_read < buffer_size {
            match reader.read(&mut buf_slice[total_read..]) {
                Ok(0) => {
                    // End of stream
                    if !bytes_read.is_null() {
                        *bytes_read = total_read;
                    }
                    return ERR_OK;
                }
                Ok(n) => {
                    total_read += n;
                }
                Err(_) => {
                    return ERR_IO_ERROR;
                }
            }
        }

        if !bytes_read.is_null() {
            *bytes_read = total_read;
        }
        ERR_OK
    }
}

/// Read entire remaining stream into a newly allocated buffer
///
/// Reads all remaining data from the stream and returns it in a newly
/// allocated buffer. Useful for smaller streams where you want all data at once.
///
/// **Warning**: This loads all content into memory. For large documents,
/// prefer using `extractous_stream_read()` in a loop.
///
/// ### Arguments
/// * `handle` - Valid StreamReader pointer
/// * `out_buffer` - Output pointer for allocated buffer
/// * `out_size` - Output pointer for buffer size
///
/// ### Returns
/// * `ERR_OK` - Success, buffer allocated and filled
/// * `ERR_NULL_POINTER` - Invalid pointer argument
/// * `ERR_IO_ERROR` - Read operation failed
/// * `ERR_OUT_OF_MEMORY` - Memory allocation failed
///
/// ### Memory Management
/// The caller must free the returned buffer using `extractous_buffer_free()`.
///
/// ### Example
/// ```c
/// uint8_t* data;
/// size_t size;
/// int result = extractous_stream_read_all(reader, &data, &size);
///
/// if (result == ERR_OK) {
///     // Use data (size bytes)
///     process_data(data, size);
///     
///     // Free the buffer
///     extractous_buffer_free(data);
/// }
/// ```
///
/// ### Safety
/// - `handle` must be a valid StreamReader pointer
/// - `out_buffer` and `out_size` must be valid pointers
/// - Returned buffer must be freed with `extractous_buffer_free()`
#[no_mangle]
pub unsafe extern "C" fn extractous_stream_read_all(
    handle: *mut CStreamReader,
    out_buffer: *mut *mut u8,
    out_size: *mut libc::size_t,
) -> libc::c_int {
    if handle.is_null() || out_buffer.is_null() || out_size.is_null() {
        return ERR_NULL_POINTER;
    }

    unsafe {
        let reader = &mut *(handle as *mut CoreStreamReader);
        let mut data = Vec::new();

        match reader.read_to_end(&mut data) {
            Ok(_) => {
                let size = data.len();
                let ptr = data.as_mut_ptr();
                std::mem::forget(data); // Prevent Rust from freeing the Vec

                *out_buffer = ptr;
                *out_size = size;
                ERR_OK
            }
            Err(_) => ERR_IO_ERROR,
        }
    }
}

/// Free a buffer allocated by extractous_stream_read_all
///
/// ### Arguments
/// * `buffer` - Pointer to buffer returned by `extractous_stream_read_all()`
/// * `size` - Size of the buffer in bytes
///
/// ### Safety
/// - `buffer` must be a pointer returned by `extractous_stream_read_all()`
/// - `size` must match the size returned by that function
/// - `buffer` must not be used after this call
/// - Do not call this function twice on the same buffer
///
/// ### Example
/// ```c
/// uint8_t* data;
/// size_t size;
/// extractous_stream_read_all(reader, &data, &size);
/// // ... use data ...
/// extractous_buffer_free(data, size);
/// ```
#[no_mangle]
pub unsafe extern "C" fn extractous_buffer_free(buffer: *mut u8, size: libc::size_t) {
    if !buffer.is_null() && size > 0 {
        unsafe {
            // Reconstruct the Vec to properly deallocate
            let _ = Vec::from_raw_parts(buffer, size, size);
        }
    }
}

/// Free stream reader and release associated resources
///
/// ### Arguments
/// * `handle` - StreamReader pointer to free
///
/// ### Safety
/// - `handle` must be a valid StreamReader pointer
/// - `handle` must not be used after this call
/// - Calling this function twice on the same pointer causes undefined behavior
/// - Safe to call with NULL (no-op)
///
/// ### Example
/// ```c
/// CStreamReader* reader;
/// // ... extract to stream and use reader ...
/// extractous_stream_free(reader);  // Always free when done
/// ```
#[no_mangle]
pub unsafe extern "C" fn extractous_stream_free(handle: *mut CStreamReader) {
    if !handle.is_null() {
        unsafe {
            let _ = Box::from_raw(handle as *mut CoreStreamReader);
        }
    }
}
