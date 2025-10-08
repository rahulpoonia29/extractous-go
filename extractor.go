package extractous

import "C"

import (
	"runtime"
	"unsafe"
	
)

// Extractor is the main entry point for document extraction
type Extractor struct {
	handle *C.CExtractor
}

// New creates a new Extractor with default configuration
func New() (*Extractor, error) {
	h := C.extractous_extractor_new()
	if h == nil {
		return nil, ErrNullPointer
	}

	ex := &Extractor{handle: h}
	runtime.SetFinalizer(ex, (*Extractor).Close)

	return ex, nil
}

// SetExtractStringMaxLength sets the maximum length of extracted text
func (e *Extractor) SetExtractStringMaxLength(maxLength int) *Extractor {
	newHandle := C.extractous_extractor_set_extract_string_max_length(e.handle, C.int(maxLength))
	if newHandle != nil {
		e.handle = newHandle
	}
	return e
}

// SetEncoding sets the character encoding for extraction
func (e *Extractor) SetEncoding(encoding CharSet) *Extractor {
	newHandle := C.extractous_extractor_set_encoding(e.handle, C.int(encoding))
	if newHandle != nil {
		e.handle = newHandle
	}
	return e
}

// SetPdfConfig sets PDF-specific extraction configuration
func (e *Extractor) SetPdfConfig(config *PdfParserConfig) *Extractor {
	if config == nil {
		return e
	}
	newHandle := C.extractous_extractor_set_pdf_config(e.handle, config.handle)
	if newHandle != nil {
		e.handle = newHandle
	}
	return e
}

// ExtractFileToString extracts text from a file to a string
func (e *Extractor) ExtractFileToString(path string) (string, Metadata, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var cContent *C.char
	var cMeta *C.CMetadata

	rc := C.extractous_extractor_extract_file_to_string(e.handle, cPath, &cContent, &cMeta)
	if rc != C.ERR_OK {
		return "", nil, errorFromCode(rc)
	}

	defer C.extractous_string_free(cContent)
	defer C.extractous_metadata_free(cMeta)

	content := C.GoString(cContent)
	metadata := cMetadataToGo(cMeta)

	return content, metadata, nil
}

// Close releases the Extractor resources
func (e *Extractor) Close() {
	if e.handle != nil {
		C.extractous_extractor_free(e.handle)
		e.handle = nil
	}
}
