package src

/*
#cgo CFLAGS: -I${SRCDIR}/../include

#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import (
	"io"
	"runtime"
	"unsafe"
)

// StreamReader implements io.Reader for streaming document content
type StreamReader struct {
	ptr    *C.struct_CStreamReader
	closed bool
}

// newStreamReader creates a StreamReader from a C pointer
func newStreamReader(ptr *C.struct_CStreamReader) *StreamReader {
	if ptr == nil {
		return nil
	}

	reader := &StreamReader{ptr: ptr}
	runtime.SetFinalizer(reader, (*StreamReader).Close)
	return reader
}

// Read implements io.Reader
func (r *StreamReader) Read(p []byte) (n int, err error) {
	if r.closed || r.ptr == nil {
		return 0, io.EOF
	}

	if len(p) == 0 {
		return 0, nil
	}

	var bytesRead C.size_t
	code := C.extractous_stream_read(
		r.ptr,
		(*C.uint8_t)(unsafe.Pointer(&p[0])),
		C.size_t(len(p)),
		&bytesRead,
	)

	if code != errOK {
		return 0, newError(code)
	}

	if bytesRead == 0 {
		return 0, io.EOF
	}

	return int(bytesRead), nil
}

// Close closes the stream and releases resources
func (r *StreamReader) Close() error {
	if r.closed || r.ptr == nil {
		return nil
	}

	C.extractous_stream_free(r.ptr)
	r.ptr = nil
	r.closed = true
	return nil
}
