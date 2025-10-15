package extractous

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

// StreamReader implements io.Reader for streaming document content.
//
// StreamReader provides efficient streaming access to extracted document content,
// allowing you to process large documents without loading everything into memory.
// It implements the standard io.Reader interface and can be used with any Go code
// that works with readers.
//
// # Interface Compliance
//
// StreamReader implements:
//   - io.Reader: Read(p []byte) (n int, err error)
//   - io.Closer: Close() error
//
// This means it can be used with:
//   - io.Copy, io.ReadAll, io.ReadFull
//   - bufio.NewReader, bufio.NewScanner
//   - io.TeeReader, io.LimitReader
//   - Any function accepting io.Reader or io.ReadCloser
//
// # Resource Management
//
// StreamReaders must be closed when done to free underlying resources. While they
// use finalizers for automatic cleanup, calling Close() explicitly is strongly
// recommended:
//
//	reader, metadata, err := extractor.ExtractFile("document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close() // Always close
//
// # Usage Patterns
//
// Copy to stdout:
//
//	reader, _, _ := extractor.ExtractFile("document.pdf")
//	defer reader.Close()
//	io.Copy(os.Stdout, reader)
//
// Read in chunks:
//
//	reader, _, _ := extractor.ExtractFile("document.pdf")
//	defer reader.Close()
//	buf := make([]byte, 4096)
//	for {
//	    n, err := reader.Read(buf)
//	    if err == io.EOF {
//	        break
//	    }
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    process(buf[:n])
//	}
//
// Use with bufio.Scanner:
//
//	reader, _, _ := extractor.ExtractFile("document.pdf")
//	defer reader.Close()
//	scanner := bufio.NewScanner(reader)
//	for scanner.Scan() {
//	    line := scanner.Text()
//	    fmt.Println(line)
//	}
//
// Read all at once (for moderate-sized documents):
//
//	reader, _, _ := extractor.ExtractFile("document.pdf")
//	defer reader.Close()
//	content, err := io.ReadAll(reader)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(content))
//
// # Performance Considerations
//
// StreamReader is buffered at the FFI layer, so you don't need to wrap it with
// bufio.Reader for basic read operations. However, bufio.Scanner can still be
// useful for line-oriented processing.
//
// Typical buffer sizes:
//   - Small reads (< 512 bytes): May have overhead, prefer larger reads
//   - Medium reads (4KB - 64KB): Optimal for most use cases
//   - Large reads (> 1MB): Generally no advantage over medium reads
//
// # Thread Safety
//
// StreamReader is NOT safe for concurrent use. Do not call Read() from multiple
// goroutines simultaneously.
type StreamReader struct {
	ptr    *C.struct_CStreamReader // FFI stream pointer
	closed bool                    // Whether Close() has been called
}

// newStreamReader creates a StreamReader from a C pointer.
//
// This is an internal function used to wrap FFI stream pointers. It sets up
// a finalizer for automatic resource cleanup.
//
// Returns nil if the C pointer is nil.
//
// Internal use only.
func newStreamReader(ptr *C.struct_CStreamReader) *StreamReader {
	if ptr == nil {
		return nil
	}

	reader := &StreamReader{ptr: ptr}
	runtime.SetFinalizer(reader, (*StreamReader).Close)
	return reader
}

// Read reads up to len(p) bytes into p.
//
// This implements the io.Reader interface. It reads extracted content from the
// document into the provided byte slice.
//
// Parameters:
//   - p: Byte slice to read into (must not be nil or empty for meaningful reads)
//
// Returns:
//   - n: Number of bytes read (0 <= n <= len(p))
//   - err: Error if read failed, or io.EOF when stream is exhausted
//
// # Behavior
//
// Read may return fewer bytes than requested (0 < n < len(p)) without error.
// This is normal io.Reader behavior and does not indicate an error or EOF.
//
// Read returns io.EOF when no more data is available. After receiving io.EOF,
// all subsequent calls will return (0, io.EOF).
//
// If the reader has been closed, Read returns (0, io.EOF).
//
// # Example
//
//	reader, _, err := extractor.ExtractFile("document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close()
//
//	buf := make([]byte, 4096)
//	for {
//	    n, err := reader.Read(buf)
//	    if n > 0 {
//	        // Process buf[:n]
//	        fmt.Print(string(buf[:n]))
//	    }
//	    if err == io.EOF {
//	        break
//	    }
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
//
// # Error Handling
//
// Errors other than io.EOF indicate actual read failures and should be handled:
//
//	n, err := reader.Read(buf)
//	if err != nil && err != io.EOF {
//	    log.Printf("Read error: %v", err)
//	    return
//	}
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

// Close closes the stream and releases underlying resources.
//
// This implements the io.Closer interface. After calling Close, the StreamReader
// should not be used. Subsequent calls to Read will return (0, io.EOF).
//
// Calling Close multiple times is safe - subsequent calls are no-ops and return
// nil.
//
// While StreamReaders use finalizers for automatic cleanup, calling Close
// explicitly is strongly recommended for deterministic resource management,
// especially when processing many documents.
//
// Returns:
//   - Always returns nil (implements io.Closer)
//
// Example:
//
//	reader, _, err := extractor.ExtractFile("document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close() // Ensure cleanup
//
//	// Use reader...
//	io.Copy(os.Stdout, reader)
//
//	// Explicit close (defer will call again, which is safe)
//	reader.Close()
//
// # Resource Management Best Practices
//
// Always use defer:
//
//	reader, _, err := extractor.ExtractFile("doc.pdf")
//	if err != nil {
//	    return err
//	}
//	defer reader.Close() // Cleanup even if function panics
//
// For long-running processes, close explicitly in loops:
//
//	for _, file := range files {
//	    reader, _, err := extractor.ExtractFile(file)
//	    if err != nil {
//	        log.Printf("Error: %v", err)
//	        continue
//	    }
//
//	    processStream(reader)
//	    reader.Close() // Don't wait for defer in loop
//	}
func (r *StreamReader) Close() error {
	if r.closed || r.ptr == nil {
		return nil
	}

	C.extractous_stream_free(r.ptr)
	r.ptr = nil
	r.closed = true
	return nil
}
