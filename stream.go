package extractous

import (
	"io"
	"runtime"

	"C"
)

// StreamReader provides buffered reading of extracted content
type StreamReader struct {
	handle *C.CStreamReader
}

// Read implements io.Reader interface
func (r *StreamReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	var bytesRead C.size_t
	rc := C.extractous_stream_read(
		r.handle,
		(*C.uchar)(&p[0]),
		C.size_t(len(p)),
		&bytesRead,
	)

	n := int(bytesRead)

	if rc != C.ERR_OK {
		if n == 0 {
			return 0, io.EOF
		}
		return n, errorFromCode(rc)
	}

	if n == 0 {
		return 0, io.EOF
	}

	return n, nil
}

// Close releases the StreamReader resources
func (r *StreamReader) Close() error {
	if r.handle != nil {
		C.extractous_stream_free(r.handle)
		r.handle = nil
	}
	return nil
}

func newStreamReader(handle *C.CStreamReader) *StreamReader {
	reader := &StreamReader{handle: handle}
	runtime.SetFinalizer(reader, (*StreamReader).Close)
	return reader
}
