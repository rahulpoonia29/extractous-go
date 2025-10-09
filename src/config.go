package src

/*
#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import "runtime"

// PdfConfig configures PDF document parsing
type PdfConfig struct {
	ptr *C.struct_CPdfParserConfig
}

// NewPdfConfig creates a new PDF configuration with default settings
func NewPdfConfig() *PdfConfig {
	ptr := C.extractous_pdf_config_new()
	if ptr == nil {
		return nil
	}

	cfg := &PdfConfig{ptr: ptr}
	runtime.SetFinalizer(cfg, (*PdfConfig).free)
	return cfg
}

// SetOcrStrategy sets the OCR strategy for PDF parsing
func (c *PdfConfig) SetOcrStrategy(strategy PdfOcrStrategy) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_ocr_strategy(c.ptr, C.int(strategy))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetExtractInlineImages sets whether to extract inline embedded images
func (c *PdfConfig) SetExtractInlineImages(value bool) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_extract_inline_images(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetExtractUniqueInlineImagesOnly sets whether to extract each unique inline image only once
func (c *PdfConfig) SetExtractUniqueInlineImagesOnly(value bool) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_extract_unique_inline_images_only(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetExtractMarkedContent sets whether to extract text and marked structure
func (c *PdfConfig) SetExtractMarkedContent(value bool) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_extract_marked_content(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetExtractAnnotationText sets whether to extract annotation text
func (c *PdfConfig) SetExtractAnnotationText(value bool) *PdfConfig {
	newPtr := C.extractous_pdf_config_set_extract_annotation_text(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

func (c *PdfConfig) free() {
	if c.ptr != nil {
		C.extractous_pdf_config_free(c.ptr)
		c.ptr = nil
	}
}

// OfficeConfig configures Microsoft Office document parsing
type OfficeConfig struct {
	ptr *C.struct_COfficeParserConfig
}

// NewOfficeConfig creates a new Office configuration with default settings
func NewOfficeConfig() *OfficeConfig {
	ptr := C.extractous_office_config_new()
	if ptr == nil {
		return nil
	}

	cfg := &OfficeConfig{ptr: ptr}
	runtime.SetFinalizer(cfg, (*OfficeConfig).free)
	return cfg
}

// SetExtractMacros sets whether to extract VBA macros
func (c *OfficeConfig) SetExtractMacros(value bool) *OfficeConfig {
	newPtr := C.extractous_office_config_set_extract_macros(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetIncludeDeletedContent sets whether to include deleted content from track changes
func (c *OfficeConfig) SetIncludeDeletedContent(value bool) *OfficeConfig {
	newPtr := C.extractous_office_config_set_include_deleted_content(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetIncludeMoveFromContent sets whether to include "move from" content
func (c *OfficeConfig) SetIncludeMoveFromContent(value bool) *OfficeConfig {
	newPtr := C.extractous_office_config_set_include_move_from_content(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetIncludeShapeBasedContent sets whether to include text from shapes
func (c *OfficeConfig) SetIncludeShapeBasedContent(value bool) *OfficeConfig {
	newPtr := C.extractous_office_config_set_include_shape_based_content(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

func (c *OfficeConfig) free() {
	if c.ptr != nil {
		C.extractous_office_config_free(c.ptr)
		c.ptr = nil
	}
}

// OcrConfig configures Tesseract OCR settings
type OcrConfig struct {
	ptr *C.struct_CTesseractOcrConfig
}

// NewOcrConfig creates a new OCR configuration with default settings
func NewOcrConfig() *OcrConfig {
	ptr := C.extractous_ocr_config_new()
	if ptr == nil {
		return nil
	}

	cfg := &OcrConfig{ptr: ptr}
	runtime.SetFinalizer(cfg, (*OcrConfig).free)
	return cfg
}

// SetLanguage sets the OCR language (e.g., "eng", "fra", "deu")
func (c *OcrConfig) SetLanguage(lang string) *OcrConfig {
	cLang := cString(lang)
	defer freeString(cLang)

	newPtr := C.extractous_ocr_config_set_language(c.ptr, cLang)
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetDensity sets the DPI density for OCR (default: 300)
func (c *OcrConfig) SetDensity(dpi int) *OcrConfig {
	newPtr := C.extractous_ocr_config_set_density(c.ptr, C.int32_t(dpi))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetDepth sets the color depth for OCR
func (c *OcrConfig) SetDepth(depth int) *OcrConfig {
	newPtr := C.extractous_ocr_config_set_depth(c.ptr, C.int32_t(depth))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetEnableImagePreprocessing enables/disables image preprocessing for better OCR accuracy
func (c *OcrConfig) SetEnableImagePreprocessing(value bool) *OcrConfig {
	newPtr := C.extractous_ocr_config_set_enable_image_preprocessing(c.ptr, C.bool(value))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

// SetTimeoutSeconds sets the timeout for OCR processing in seconds
func (c *OcrConfig) SetTimeoutSeconds(seconds int) *OcrConfig {
	newPtr := C.extractous_ocr_config_set_timeout_seconds(c.ptr, C.int32_t(seconds))
	if newPtr == nil {
		return nil
	}
	c.ptr = newPtr
	return c
}

func (c *OcrConfig) free() {
	if c.ptr != nil {
		C.extractous_ocr_config_free(c.ptr)
		c.ptr = nil
	}
}
