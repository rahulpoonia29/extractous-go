package src

/*
#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

// Extractor is the main entry point for document extraction
type Extractor struct {
	ptr *C.struct_CExtractor
}

// New creates a new Extractor with default configuration
func New() *Extractor {
	ptr := C.extractous_extractor_new()
	if ptr == nil {
		return nil
	}

	ext := &Extractor{ptr: ptr}
	runtime.SetFinalizer(ext, (*Extractor).Close)
	return ext
}

// SetExtractStringMaxLength sets the maximum length for extracted string content
func (e *Extractor) SetExtractStringMaxLength(maxLen int) *Extractor {
	newPtr := C.extractous_extractor_set_extract_string_max_length(e.ptr, C.int(maxLen))
	if newPtr == nil {
		return nil
	}
	e.ptr = newPtr
	return e
}

// SetEncoding sets the character encoding for extraction
func (e *Extractor) SetEncoding(charset CharSet) *Extractor {
	newPtr := C.extractous_extractor_set_encoding(e.ptr, C.int(charset))
	if newPtr == nil {
		return nil
	}
	e.ptr = newPtr
	return e
}

// SetPdfConfig sets PDF parsing configuration
func (e *Extractor) SetPdfConfig(config *PdfConfig) *Extractor {
	if config == nil || config.ptr == nil {
		return nil
	}
	newPtr := C.extractous_extractor_set_pdf_config(e.ptr, config.ptr)
	if newPtr == nil {
		return nil
	}
	e.ptr = newPtr
	return e
}

// SetOfficeConfig sets Office document parsing configuration
func (e *Extractor) SetOfficeConfig(config *OfficeConfig) *Extractor {
	if config == nil || config.ptr == nil {
		return nil
	}
	newPtr := C.extractous_extractor_set_office_config(e.ptr, config.ptr)
	if newPtr == nil {
		return nil
	}
	e.ptr = newPtr
	return e
}

// SetOcrConfig sets OCR configuration
func (e *Extractor) SetOcrConfig(config *OcrConfig) *Extractor {
	if config == nil || config.ptr == nil {
		return nil
	}
	newPtr := C.extractous_extractor_set_ocr_config(e.ptr, config.ptr)
	if newPtr == nil {
		return nil
	}
	e.ptr = newPtr
	return e
}

// ExtractFileToString extracts a file's content to a string
func (e *Extractor) ExtractFileToString(path string) (content string, metadata Metadata, err error) {
	cPath := cString(path)
	defer freeString(cPath)

	var cContent *C.char
	var cMeta *C.struct_CMetadata

	code := C.extractous_extractor_extract_file_to_string(e.ptr, cPath, &cContent, &cMeta)
	if code != errOK {
		return "", nil, newError(code)
	}

	content = goString(cContent)
	C.extractous_string_free(cContent)

	metadata = newMetadata(cMeta)
	return content, metadata, nil
}

// ExtractFile extracts a file's content to a streaming reader
func (e *Extractor) ExtractFile(path string) (reader *StreamReader, metadata Metadata, err error) {
	cPath := cString(path)
	defer freeString(cPath)

	var cReader *C.struct_CStreamReader
	var cMeta *C.struct_CMetadata

	code := C.extractous_extractor_extract_file(e.ptr, cPath, &cReader, &cMeta)
	if code != errOK {
		return nil, nil, newError(code)
	}

	reader = newStreamReader(cReader)
	metadata = newMetadata(cMeta)
	return reader, metadata, nil
}

// ExtractBytesToString extracts content from a byte slice to a string
func (e *Extractor) ExtractBytesToString(data []byte) (content string, metadata Metadata, err error) {
	if len(data) == 0 {
		return "", make(Metadata), nil
	}

	var cContent *C.char
	var cMeta *C.struct_CMetadata

	code := C.extractous_extractor_extract_bytes_to_string(
		e.ptr,
		(*C.uint8_t)(&data[0]),
		C.size_t(len(data)),
		&cContent,
		&cMeta,
	)

	if code != errOK {
		return "", nil, newError(code)
	}

	content = goString(cContent)
	C.extractous_string_free(cContent)

	metadata = newMetadata(cMeta)
	return content, metadata, nil
}

// ExtractBytes extracts content from a byte slice to a streaming reader
func (e *Extractor) ExtractBytes(data []byte) (reader *StreamReader, metadata Metadata, err error) {
	if len(data) == 0 {
		return nil, make(Metadata), nil
	}

	var cReader *C.struct_CStreamReader
	var cMeta *C.struct_CMetadata

	code := C.extractous_extractor_extract_bytes(
		e.ptr,
		(*C.uint8_t)(&data[0]),
		C.size_t(len(data)),
		&cReader,
		&cMeta,
	)

	if code != errOK {
		return nil, nil, newError(code)
	}

	reader = newStreamReader(cReader)
	metadata = newMetadata(cMeta)
	return reader, metadata, nil
}

// Close releases the extractor's resources
func (e *Extractor) Close() error {
	if e.ptr != nil {
		C.extractous_extractor_free(e.ptr)
		e.ptr = nil
	}
	return nil
}
