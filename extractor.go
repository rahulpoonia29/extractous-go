package extractous

/*
#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

// Extractor is the main entry point for document extraction.
//
// It provides methods to extract text content and metadata from various document
// formats including PDF, Word, Excel, PowerPoint, HTML, and many others. The
// Extractor can process files, URLs, or byte slices, with options for both
// string-based extraction (loading all content into memory) and streaming
// extraction (for large documents).
//
// # Thread Safety
//
// Extractor instances are NOT safe for concurrent use by multiple goroutines.
// Each goroutine should create its own Extractor instance. The underlying
// implementation requires thread affinity, and sharing extractors across
// goroutines can lead to undefined behavior and crashes.
//
//	// WRONG - unsafe concurrent use
//	extractor := extractous.New()
//	for _, file := range files {
//	    go func(f string) {
//	        extractor.ExtractFileToString(f) // UNSAFE!
//	    }(file)
//	}
//
//	// CORRECT - separate extractor per goroutine
//	for _, file := range files {
//	    go func(f string) {
//	        extractor := extractous.New()
//	        defer extractor.Close()
//	        extractor.ExtractFileToString(f) // Safe
//	    }(file)
//	}
//
// # Memory Management
//
// Extractors use finalizers for automatic cleanup, but calling Close() explicitly
// is strongly recommended for deterministic resource cleanup, especially in
// long-running applications or when processing many documents.
//
// Builder methods (Set*) follow the builder pattern and return a new Extractor
// instance. The old instance is consumed and should not be used:
//
//	// WRONG - old extractor is invalid
//	extractor := extractous.New()
//	extractor.SetXmlOutput(true) // extractor is now invalid!
//
//	// CORRECT - use returned value
//	extractor := extractous.New()
//	extractor = extractor.SetXmlOutput(true) // extractor is valid
//
//	// BEST - chain method calls
//	extractor := extractous.New().
//	    SetExtractStringMaxLength(10000).
//	    SetXmlOutput(true)
//
// # Basic Usage
//
// Extract a document to a string:
//
//	extractor := extractous.New()
//	defer extractor.Close()
//
//	content, metadata, err := extractor.ExtractFileToString("document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Content:", content)
//	fmt.Println("Author:", metadata.Get("author"))
//
// # Configuration
//
// Configure extraction behavior before processing documents:
//
//	pdfConfig := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrAuto).
//	    SetExtractAnnotationText(true)
//
//	extractor := extractous.New().
//	    SetPdfConfig(pdfConfig).
//	    SetExtractStringMaxLength(10_000_000)
//	defer extractor.Close()
//
// # Streaming Large Documents
//
// For large documents, use streaming to avoid loading everything into memory:
//
//	reader, metadata, err := extractor.ExtractFile("large.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close()
//
//	// Read in chunks
//	buf := make([]byte, 4096)
//	for {
//	    n, err := reader.Read(buf)
//	    if err == io.EOF {
//	        break
//	    }
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    fmt.Print(string(buf[:n]))
//	}
type Extractor struct {
	ptr *C.struct_CExtractor
}

// New creates a new Extractor with default configuration.
//
// The default configuration includes:
//   - Character encoding: UTF-8
//   - Maximum string length: 100 MB
//   - XML output: disabled
//   - OCR: disabled (for PDFs)
//   - All Office features: enabled
//
// The returned Extractor must be closed when no longer needed:
//
//	extractor := extractous.New()
//	defer extractor.Close()
//
// Returns nil if the extractor could not be created (rare, usually indicates
// system resource exhaustion).
func New() *Extractor {
	ptr := C.extractous_extractor_new()
	if ptr == nil {
		return nil
	}

	ext := &Extractor{ptr: ptr}
	runtime.SetFinalizer(ext, (*Extractor).Close)
	return ext
}

// SetExtractStringMaxLength sets the maximum length for extracted string content.
//
// This prevents excessive memory consumption when extracting very large documents
// to a string. Content exceeding this limit will be truncated. This setting only
// affects *ToString methods; streaming extraction is not limited.
//
// The default limit is 100 MB (104,857,600 bytes).
//
// Example:
//
//	// Limit to 10 MB
//	extractor := extractous.New().
//	    SetExtractStringMaxLength(10 * 1024 * 1024)
//
// This method consumes the receiver and returns a new Extractor. Always use the
// returned value.
//
// Returns nil if the configuration failed.
func (e *Extractor) SetExtractStringMaxLength(maxLen int) *Extractor {
	newPtr := C.extractous_extractor_set_extract_string_max_length(e.ptr, C.int(maxLen))
	if newPtr == nil {
		return nil
	}
	e.ptr = newPtr
	return e
}

// SetEncoding sets the character encoding for extraction.
//
// Most documents should use UTF-8 (the default). Only change this if you need
// a specific encoding for compatibility with legacy systems.
//
// Supported encodings:
//   - CharSetUTF8 (default, recommended)
//   - CharSetUSASCII
//   - CharSetUTF16BE
//
// Example:
//
//	extractor := extractous.New().
//	    SetEncoding(extractous.CharSetUTF8)
//
// This method consumes the receiver and returns a new Extractor.
//
// Returns nil if the encoding is invalid.
func (e *Extractor) SetEncoding(charset CharSet) *Extractor {
	newPtr := C.extractous_extractor_set_encoding(e.ptr, C.int(charset))
	if newPtr == nil {
		return nil
	}
	e.ptr = newPtr
	return e
}

// SetPdfConfig sets PDF parsing configuration.
//
// Use this to configure PDF-specific behavior such as OCR strategy, inline
// image extraction, and annotation text extraction.
//
// Example:
//
//	pdfConfig := extractous.NewPdfConfig().
//	    SetOcrStrategy(extractous.PdfOcrAuto).
//	    SetExtractAnnotationText(true)
//
//	extractor := extractous.New().
//	    SetPdfConfig(pdfConfig)
//
// The config is consumed by this method and becomes owned by the Extractor.
// Do not use the config after passing it to this method.
//
// This method consumes the receiver and returns a new Extractor.
//
// Returns nil if the config is nil or invalid.
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

// SetOfficeConfig sets Microsoft Office document parsing configuration.
//
// Use this to configure Office-specific behavior such as macro extraction,
// handling of deleted content, and shape-based content extraction.
//
// Example:
//
//	officeConfig := extractous.NewOfficeConfig().
//	    SetIncludeShapeBasedContent(true).
//	    SetExtractMacros(false) // More secure
//
//	extractor := extractous.New().
//	    SetOfficeConfig(officeConfig)
//
// The config is consumed by this method and becomes owned by the Extractor.
//
// This method consumes the receiver and returns a new Extractor.
//
// Returns nil if the config is nil or invalid.
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

// SetOcrConfig sets Tesseract OCR configuration.
//
// Use this to configure OCR behavior for images and scanned documents. Tesseract
// must be installed on the system for OCR to work.
//
// Example:
//
//	ocrConfig := extractous.NewOcrConfig().
//	    SetLanguage("eng").
//	    SetDensity(300).
//	    SetTimeoutSeconds(120)
//
//	extractor := extractous.New().
//	    SetOcrConfig(ocrConfig)
//
// Common language codes:
//   - "eng" - English
//   - "deu" - German
//   - "fra" - French
//   - "spa" - Spanish
//   - "eng+fra" - Multiple languages
//
// The config is consumed by this method and becomes owned by the Extractor.
//
// This method consumes the receiver and returns a new Extractor.
//
// Returns nil if the config is nil or invalid.
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

// SetXmlOutput sets whether to output XML structure instead of plain text.
//
// When enabled, the extracted content will be in XML format with structural
// information (paragraphs, headings, etc.). When disabled (default), only plain
// text is extracted.
//
// Example:
//
//	// Enable XML output for structured extraction
//	extractor := extractous.New().
//	    SetXmlOutput(true)
//
// This method consumes the receiver and returns a new Extractor.
//
// Returns nil if the configuration failed.
func (e *Extractor) SetXmlOutput(xmlOutput bool) *Extractor {
	if e == nil || e.ptr == nil {
		return nil
	}
	newPtr := C.extractous_extractor_set_xml_output(e.ptr, C.bool(xmlOutput))
	if newPtr == nil {
		return nil
	}
	e.ptr = newPtr
	return e
}

// ExtractFileToString extracts a file's content to a string.
//
// This method loads the entire document content into memory, which is suitable
// for small to medium documents. For large documents (> 100 MB), consider using
// ExtractFile for streaming extraction.
//
// Parameters:
//   - path: File system path to the document (absolute or relative)
//
// Returns:
//   - content: Extracted text content (may be truncated if exceeds max length)
//   - metadata: Document metadata (author, title, creation date, etc.)
//   - err: Error if extraction failed
//
// Example:
//
//	content, metadata, err := extractor.ExtractFileToString("document.pdf")
//	if err != nil {
//	    if errors.Is(err, extractous.ErrIO) {
//	        fmt.Println("File not found or not readable")
//	    } else {
//	        fmt.Printf("Extraction failed: %v\n", err)
//	    }
//	    return
//	}
//
//	fmt.Println("Content:", content)
//	fmt.Println("Title:", metadata.Get("title"))
//	fmt.Println("Author:", metadata.Get("author"))
//
// Supported file formats: PDF, DOCX, XLSX, PPTX, ODT, HTML, TXT, and many more.
func (e *Extractor) ExtractFileToString(path string) (content string, metadata Metadata, err error) {
	if e == nil || e.ptr == nil {
		return "", nil, ErrNullPointer
	}

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

// ExtractFile extracts a file's content as a streaming reader.
//
// This method returns an io.Reader that streams the document content, making it
// suitable for large documents that shouldn't be loaded entirely into memory.
// The reader must be closed when done.
//
// Parameters:
//   - path: File system path to the document (absolute or relative)
//
// Returns:
//   - reader: StreamReader implementing io.Reader
//   - metadata: Document metadata
//   - err: Error if extraction failed
//
// Example:
//
//	reader, metadata, err := extractor.ExtractFile("large.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close()
//
//	// Copy to stdout
//	io.Copy(os.Stdout, reader)
//
//	// Or read in chunks
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
// The reader is buffered and can be used with any io.Reader-compatible code.
func (e *Extractor) ExtractFile(path string) (reader *StreamReader, metadata Metadata, err error) {
	if e == nil || e.ptr == nil {
		return nil, nil, ErrNullPointer
	}

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

// ExtractBytesToString extracts content from a byte slice to a string.
//
// Use this when you have document data already loaded in memory (e.g., from a
// database, network request, or embedded resource).
//
// Parameters:
//   - data: Document data as a byte slice
//
// Returns:
//   - content: Extracted text content
//   - metadata: Document metadata
//   - err: Error if extraction failed
//
// Example:
//
//	// From HTTP response
//	resp, _ := http.Get("https://example.com/document.pdf")
//	defer resp.Body.Close()
//	data, _ := io.ReadAll(resp.Body)
//
//	content, metadata, err := extractor.ExtractBytesToString(data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// From database
//	var docData []byte
//	db.QueryRow("SELECT data FROM documents WHERE id = ?", id).Scan(&docData)
//	content, metadata, err := extractor.ExtractBytesToString(docData)
//
// For large byte slices, consider using ExtractBytes for streaming extraction.
func (e *Extractor) ExtractBytesToString(data []byte) (content string, metadata Metadata, err error) {
	if e == nil || e.ptr == nil {
		return "", nil, ErrNullPointer
	}

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

// ExtractBytes extracts content from a byte slice to a streaming reader.
//
// This is the streaming version of ExtractBytesToString, suitable for large
// documents already in memory but too large to process as a single string.
//
// Parameters:
//   - data: Document data as a byte slice
//
// Returns:
//   - reader: StreamReader implementing io.Reader
//   - metadata: Document metadata
//   - err: Error if extraction failed
//
// Example:
//
//	var docData []byte
//	// ... load document data ...
//
//	reader, metadata, err := extractor.ExtractBytes(docData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close()
//
//	// Process stream
//	scanner := bufio.NewScanner(reader)
//	for scanner.Scan() {
//	    line := scanner.Text()
//	    processLine(line)
//	}
func (e *Extractor) ExtractBytes(data []byte) (reader *StreamReader, metadata Metadata, err error) {
	if e == nil || e.ptr == nil {
		return nil, nil, ErrNullPointer
	}

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

// ExtractURLToString extracts content from a URL to a string.
//
// This method fetches the document from the URL and extracts its content. The
// URL can point to any supported document format.
//
// Parameters:
//   - url: HTTP(S) URL to the document
//
// Returns:
//   - content: Extracted text content
//   - metadata: Document metadata
//   - err: Error if download or extraction failed
//
// Example:
//
//	content, metadata, err := extractor.ExtractURLToString(
//	    "https://example.com/document.pdf",
//	)
//	if err != nil {
//	    if errors.Is(err, extractous.ErrIO) {
//	        fmt.Println("Download failed or file not found")
//	    } else {
//	        fmt.Printf("Extraction failed: %v\n", err)
//	    }
//	    return
//	}
//
// Note: This method downloads the entire document before extraction. For large
// remote documents, consider downloading to a file first and using ExtractFile,
// or use ExtractURL for streaming.
func (e *Extractor) ExtractURLToString(url string) (content string, metadata Metadata, err error) {
	if e == nil || e.ptr == nil {
		return "", nil, ErrNullPointer
	}

	cUrl := cString(url)
	defer freeString(cUrl)

	var cContent *C.char
	var cMeta *C.struct_CMetadata

	code := C.extractous_extractor_extract_url_to_string(e.ptr, cUrl, &cContent, &cMeta)
	if code != errOK {
		return "", nil, newError(code)
	}

	content = goString(cContent)
	C.extractous_string_free(cContent)

	metadata = newMetadata(cMeta)

	return content, metadata, nil
}

// ExtractURL extracts content from a URL as a streaming reader.
//
// This is the streaming version of ExtractURLToString, suitable for large
// remote documents.
//
// Parameters:
//   - url: HTTP(S) URL to the document
//
// Returns:
//   - reader: StreamReader implementing io.Reader
//   - metadata: Document metadata
//   - err: Error if download or extraction failed
//
// Example:
//
//	reader, metadata, err := extractor.ExtractURL(
//	    "https://example.com/large-document.pdf",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close()
//
//	// Process stream
//	io.Copy(os.Stdout, reader)
func (e *Extractor) ExtractURL(url string) (reader *StreamReader, metadata Metadata, err error) {
	if e == nil || e.ptr == nil {
		return nil, nil, ErrNullPointer
	}

	cUrl := cString(url)
	defer freeString(cUrl)

	var cReader *C.struct_CStreamReader
	var cMeta *C.struct_CMetadata

	code := C.extractous_extractor_extract_url(e.ptr, cUrl, &cReader, &cMeta)
	if code != errOK {
		return nil, nil, newError(code)
	}

	reader = newStreamReader(cReader)
	metadata = newMetadata(cMeta)

	return reader, metadata, nil
}

// Close releases the extractor's resources.
//
// While extractors use finalizers for automatic cleanup, calling Close explicitly
// is recommended for deterministic resource management, especially when processing
// many documents or in long-running applications.
//
// After Close is called, the Extractor should not be used. Calling Close multiple
// times is safe (subsequent calls are no-ops).
//
// Example:
//
//	extractor := extractous.New()
//	defer extractor.Close() // Always close when done
//
//	content, _, err := extractor.ExtractFileToString("document.pdf")
//	// ... use content ...
//
// Always returns nil (implements io.Closer for compatibility).
func (e *Extractor) Close() error {
	if e.ptr != nil {
		C.extractous_extractor_free(e.ptr)
		e.ptr = nil
	}
	return nil
}
