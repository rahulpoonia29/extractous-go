package extractous

import "C"
import "runtime"

// PdfParserConfig holds PDF-specific extraction configuration
type PdfParserConfig struct {
	handle *C.CPdfParserConfig
}

// NewPdfParserConfig creates a new PDF parser configuration
func NewPdfParserConfig() *PdfParserConfig {
	h := C.extractous_pdf_config_new()
	config := &PdfParserConfig{handle: h}
	runtime.SetFinalizer(config, (*PdfParserConfig).Close)
	return config
}

// SetOcrStrategy sets the OCR strategy for PDF extraction
func (c *PdfParserConfig) SetOcrStrategy(strategy PdfOcrStrategy) *PdfParserConfig {
	newHandle := C.extractous_pdf_config_set_ocr_strategy(c.handle, C.int(strategy))
	if newHandle != nil {
		c.handle = newHandle
	}
	return c
}

// SetExtractInlineImages enables/disables extraction of inline images
func (c *PdfParserConfig) SetExtractInlineImages(value bool) *PdfParserConfig {
	newHandle := C.extractous_pdf_config_set_extract_inline_images(c.handle, C.bool(value))
	if newHandle != nil {
		c.handle = newHandle
	}
	return c
}

// Close releases the config resources
func (c *PdfParserConfig) Close() {
	if c.handle != nil {
		C.extractous_pdf_config_free(c.handle)
		c.handle = nil
	}
}

// OfficeParserConfig and TesseractOcrConfig follow the same pattern...
