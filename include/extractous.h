/* Extractous Go FFI - Auto-generated */

#ifndef EXTRACTOUS_FFI_H
#define EXTRACTOUS_FFI_H

#pragma once

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

#define ERR_OK 0

#define ERR_NULL_POINTER -1

#define ERR_INVALID_UTF8 -2

#define ERR_INVALID_STRING -3

#define ERR_EXTRACTION_FAILED -4

#define ERR_IO_ERROR -5

#define ERR_INVALID_CONFIG -6

#define ERR_INVALID_ENUM -7

/**
 * UTF-8 encoding (default)
 */
#define CHARSET_UTF_8 0

/**
 * US-ASCII encoding
 */
#define CHARSET_US_ASCII 1

/**
 * UTF-16 Big Endian encoding
 */
#define CHARSET_UTF_16BE 2

/**
 * No OCR, extract existing text only
 */
#define PDF_OCR_NO_OCR 0

/**
 * OCR only, ignore existing text
 */
#define PDF_OCR_OCR_ONLY 1

/**
 * OCR and extract existing text
 */
#define PDF_OCR_OCR_AND_TEXT_EXTRACTION 2

/**
 * Automatically decide based on content
 */
#define PDF_OCR_AUTO 3

/**
 * Opaque handle to a PdfParserConfig instance
 */
typedef struct CPdfParserConfig {
  uint8_t _private[0];
} CPdfParserConfig;

/**
 * Opaque handle to an OfficeParserConfig instance
 */
typedef struct COfficeParserConfig {
  uint8_t _private[0];
} COfficeParserConfig;

/**
 * Opaque handle to a TesseractOcrConfig instance
 */
typedef struct CTesseractOcrConfig {
  uint8_t _private[0];
} CTesseractOcrConfig;

/**
 * Opaque handle to an Extractor instance
 *
 * This is an opaque pointer that should only be used through the FFI functions.
 * The actual Extractor is stored on the heap.
 */
typedef struct CExtractor {
  uint8_t _private[0];
} CExtractor;

/**
 * C-compatible metadata structure
 *
 * Contains parallel arrays of keys and values, with length stored separately.
 * Both keys and values are null-terminated C strings.
 */
typedef struct CMetadata {
  /**
   * Array of key string pointers
   */
  char **keys;
  /**
   * Array of value string pointers (comma-separated if multiple values)
   */
  char **values;
  /**
   * Number of key-value pairs
   */
  size_t len;
} CMetadata;

/**
 * Opaque handle to a StreamReader instance
 */
typedef struct CStreamReader {
  uint8_t _private[0];
} CStreamReader;

#ifdef __cplusplus
extern "C" {
#endif // __cplusplus

/**
 * Create new PDF parser config with default settings.
 */
struct CPdfParserConfig *extractous_pdf_config_new(void);

/**
 * Sets the OCR strategy for PDF parsing.
 *
 * # Safety
 * - `handle` must be a valid PdfParserConfig pointer.
 * - `strategy` must be a valid PDF_OCR_* constant.
 * - Returns a NEW handle; old handle is consumed and must not be used.
 */
struct CPdfParserConfig *extractous_pdf_config_set_ocr_strategy(struct CPdfParserConfig *handle,
                                                                int strategy);

/**
 * If true, extract the literal inline embedded OBXImages. Use with caution.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct CPdfParserConfig *extractous_pdf_config_set_extract_inline_images(struct CPdfParserConfig *handle,
                                                                         bool value);

/**
 * If true, extract each unique inline image only once.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct CPdfParserConfig *extractous_pdf_config_set_extract_unique_inline_images_only(struct CPdfParserConfig *handle,
                                                                                     bool value);

/**
 * If true, try to extract text and its marked structure.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct CPdfParserConfig *extractous_pdf_config_set_extract_marked_content(struct CPdfParserConfig *handle,
                                                                          bool value);

/**
 * If true, try to extract the text of annotations.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct CPdfParserConfig *extractous_pdf_config_set_extract_annotation_text(struct CPdfParserConfig *handle,
                                                                           bool value);

/**
 * Free PDF config.
 *
 * # Safety
 * - `handle` must be a valid PdfParserConfig pointer.
 * - `handle` must not be used after this call.
 */
void extractous_pdf_config_free(struct CPdfParserConfig *handle);

/**
 * Create new Office parser config.
 */
struct COfficeParserConfig *extractous_office_config_new(void);

/**
 * Sets whether MSOffice parsers should extract macros.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct COfficeParserConfig *extractous_office_config_set_extract_macros(struct COfficeParserConfig *handle,
                                                                        bool value);

/**
 * Whether to include deleted content from DOCX files.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct COfficeParserConfig *extractous_office_config_set_include_deleted_content(struct COfficeParserConfig *handle,
                                                                                 bool value);

/**
 * Whether to include content from "moveFrom" sections in DOCX.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct COfficeParserConfig *extractous_office_config_set_include_move_from_content(struct COfficeParserConfig *handle,
                                                                                   bool value);

/**
 * Whether to include text from drawing shapes.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct COfficeParserConfig *extractous_office_config_set_include_shape_based_content(struct COfficeParserConfig *handle,
                                                                                     bool value);

/**
 * Free Office config.
 *
 * # Safety
 * - `handle` must be valid and not used after this call.
 */
void extractous_office_config_free(struct COfficeParserConfig *handle);

/**
 * Create new Tesseract OCR config.
 */
struct CTesseractOcrConfig *extractous_ocr_config_new(void);

/**
 * Sets the OCR language.
 *
 * # Safety
 * - `handle` must be a valid TesseractOcrConfig pointer.
 * - `language` must be a valid null-terminated UTF-8 string.
 * - Returns a NEW handle; old handle is consumed.
 */
struct CTesseractOcrConfig *extractous_ocr_config_set_language(struct CTesseractOcrConfig *handle,
                                                               const char *language);

/**
 * Sets the DPI (dots per inch) for OCR.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct CTesseractOcrConfig *extractous_ocr_config_set_density(struct CTesseractOcrConfig *handle,
                                                              int32_t density);

/**
 * Sets the color depth for OCR.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct CTesseractOcrConfig *extractous_ocr_config_set_depth(struct CTesseractOcrConfig *handle,
                                                            int32_t depth);

/**
 * Sets whether to enable image preprocessing for OCR.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct CTesseractOcrConfig *extractous_ocr_config_set_enable_image_preprocessing(struct CTesseractOcrConfig *handle,
                                                                                 bool value);

/**
 * Sets the timeout in seconds for the OCR process.
 *
 * # Safety
 * - Returns a NEW handle; old handle is consumed.
 */
struct CTesseractOcrConfig *extractous_ocr_config_set_timeout_seconds(struct CTesseractOcrConfig *handle,
                                                                      int32_t seconds);

/**
 * Free OCR config.
 *
 * # Safety
 * - `handle` must be valid and not used after this call.
 */
void extractous_ocr_config_free(struct CTesseractOcrConfig *handle);

char *extractous_error_message(int code);

/**
 * Create a new Extractor with default configuration
 *
 * # Returns
 * Pointer to new Extractor, or NULL on failure.
 * Must be freed with `extractous_extractor_free`.
 */
struct CExtractor *extractous_extractor_new(void);

/**
 * Free an Extractor instance
 *
 * # Safety
 * - `handle` must be a valid pointer returned by `extractous_extractor_new`
 * - `handle` must not be used after this call
 * - Calling this twice on the same pointer causes undefined behavior
 */
void extractous_extractor_free(struct CExtractor *handle);

/**
 * Set maximum length for extracted string content
 *
 * # Safety
 * - `handle` must be a valid Extractor pointer
 * - Returns a NEW handle; old handle is consumed and must not be used
 *
 * # Returns
 * New Extractor handle with updated config, or NULL on error.
 */
struct CExtractor *extractous_extractor_set_extract_string_max_length(struct CExtractor *handle,
                                                                      int max_length);

/**
 * Set character encoding for extraction
 *
 * # Safety
 * - `handle` must be a valid Extractor pointer
 * - `encoding` must be a valid CHARSET_* constant
 * - Returns a NEW handle; old handle is consumed
 *
 * # Returns
 * New Extractor handle, or NULL if encoding is invalid.
 */
struct CExtractor *extractous_extractor_set_encoding(struct CExtractor *handle, int encoding);

/**
 * Set PDF parser configuration
 *
 * # Safety
 * - `handle` must be a valid Extractor pointer
 * - `config` must be a valid PdfParserConfig pointer
 * - Returns a NEW handle; old handle is consumed
 */
struct CExtractor *extractous_extractor_set_pdf_config(struct CExtractor *handle,
                                                       struct CPdfParserConfig *config);

/**
 * Set Office parser configuration
 *
 * # Safety
 * Same safety requirements as `extractous_extractor_set_pdf_config`.
 */
struct CExtractor *extractous_extractor_set_office_config(struct CExtractor *handle,
                                                          struct COfficeParserConfig *config);

/**
 * Set OCR configuration
 *
 * # Safety
 * Same safety requirements as `extractous_extractor_set_pdf_config`.
 */
struct CExtractor *extractous_extractor_set_ocr_config(struct CExtractor *handle,
                                                       struct CTesseractOcrConfig *config);

/**
 * Extract file content to string
 *
 * # Safety
 * - `handle` must be a valid Extractor pointer
 * - `path` must be a valid null-terminated UTF-8 string
 * - `out_content` and `out_metadata` must be valid pointers
 * - Caller must free returned content with `extractous_string_free`
 * - Caller must free returned metadata with `extractous_metadata_free`
 *
 * # Returns
 * ERR_OK on success, error code on failure.
 */
int extractous_extractor_extract_file_to_string(struct CExtractor *handle,
                                                const char *path,
                                                char **out_content,
                                                struct CMetadata **out_metadata);

/**
 * Extract file content to stream
 *
 * # Safety
 * - `handle` must be a valid Extractor pointer
 * - `path` must be a valid null-terminated UTF-8 string
 * - `out_reader` and `out_metadata` must be valid pointers
 * - Caller must free returned reader with `extractous_stream_free`
 * - Caller must free returned metadata with `extractous_metadata_free`
 *
 * # Returns
 * ERR_OK on success, error code on failure.
 */
int extractous_extractor_extract_file(struct CExtractor *handle,
                                      const char *path,
                                      struct CStreamReader **out_reader,
                                      struct CMetadata **out_metadata);

/**
 * Extract from byte array to string
 *
 * # Safety
 * - `handle` must be a valid Extractor pointer
 * - `data` must point to at least `data_len` valid bytes
 * - `out_content` and `out_metadata` must be valid pointers
 *
 * # Returns
 * ERR_OK on success, error code on failure.
 */
int extractous_extractor_extract_bytes_to_string(struct CExtractor *handle,
                                                 const uint8_t *data,
                                                 size_t data_len,
                                                 char **out_content,
                                                 struct CMetadata **out_metadata);

/**
 * Extract from byte array to stream
 *
 * # Safety
 * - `handle` must be a valid Extractor pointer
 * - `data` must point to at least `data_len` valid bytes
 * - `out_reader` and `out_metadata` must be valid pointers
 *
 * # Returns
 * ERR_OK on success, error code on failure.
 */
int extractous_extractor_extract_bytes(struct CExtractor *handle,
                                       const uint8_t *data,
                                       size_t data_len,
                                       struct CStreamReader **out_reader,
                                       struct CMetadata **out_metadata);

/**
 * Free a string allocated by Rust
 *
 * # Safety
 * - `s` must be a pointer returned by an extractous function
 * - `s` must not be used after this call
 * - Calling this twice on the same pointer causes undefined behavior
 */
void extractous_string_free(char *s);

/**
 * Free metadata structure
 *
 * # Safety
 * - `meta` must be a valid pointer returned by an extraction function
 * - `meta` must not be used after this call
 */
void extractous_metadata_free(struct CMetadata *meta);

/**
 * Read from stream into buffer
 *
 * # Safety
 * - `handle` must be a valid StreamReader pointer
 * - `buffer` must point to at least `buffer_size` bytes
 * - `bytes_read` can be NULL, otherwise must be valid pointer
 *
 * # Returns
 * ERR_OK on success, ERR_IO_ERROR on failure, or ERR_OK with 0 bytes_read on EOF.
 */
int extractous_stream_read(struct CStreamReader *handle,
                           uint8_t *buffer,
                           size_t buffer_size,
                           size_t *bytes_read);

/**
 * Free stream reader
 *
 * # Safety
 * - `handle` must be a valid StreamReader pointer
 * - `handle` must not be used after this call
 */
void extractous_stream_free(struct CStreamReader *handle);

#ifdef __cplusplus
}  // extern "C"
#endif  // __cplusplus

#endif  /* EXTRACTOUS_FFI_H */
