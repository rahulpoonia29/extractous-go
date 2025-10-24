/* 
 * Extractous FFI - C Interface
 * 
 * This header file provides a C-compatible interface to the Extractous
 * document extraction library. It is safe for use with Go via cgo or any
 * C-compatible FFI system.
 *
 * License: Apache-2.0
 * Repository: https://github.com/rahulpoonia229/extractous-go
 *
 * MEMORY MANAGEMENT:
 * All pointers returned by Extractous functions must be freed using the function extractous_free_string.
 * Failure to do so will result in memory leaks.
 *
 *
 * CGO USAGE:
 *   // #cgo CFLAGS: -I${SRCDIR}/include
 *   // #cgo LDFLAGS: -L${SRCDIR}/lib -lextractous_ffi
 *   // #cgo linux LDFLAGS: -Wl,-rpath,$ORIGIN
 *   // #cgo darwin LDFLAGS: -Wl,-rpath,@loader_path
 *   // #include "extractous.h"
 *   import "C"
 */


#ifndef EXTRACTOUS_H
#define EXTRACTOUS_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#define ERR_OK 0

#define ERR_NULL_POINTER -1

#define ERR_INVALID_UTF8 -2

#define ERR_INVALID_STRING -3

#define ERR_EXTRACTION_FAILED -4

#define ERR_IO_ERROR -5

#define ERR_INVALID_CONFIG -6

#define ERR_INVALID_ENUM -7

#define ERR_UNSUPPORTED_FORMAT -8

#define ERR_OUT_OF_MEMORY -9

#define ERR_OCR_FAILED -10

#define CHARSET_UTF_8 0

#define CHARSET_US_ASCII 1

#define CHARSET_UTF_16BE 3

#define PDF_OCR_STRATEGY_NO_OCR 0

#define PDF_OCR_STRATEGY_OCR_ONLY 1

#define PDF_OCR_STRATEGY_OCR_AND_TEXT_EXTRACTION 2

#define PDF_OCR_STRATEGY_AUTO 3

typedef struct CPdfParserConfig {
  uint8_t _private[0];
} CPdfParserConfig;

typedef struct COfficeParserConfig {
  uint8_t _private[0];
} COfficeParserConfig;

typedef struct CTesseractOcrConfig {
  uint8_t _private[0];
} CTesseractOcrConfig;

typedef struct CExtractor {
  uint8_t _private[0];
} CExtractor;

typedef struct CMetadata {
  /*
   Array of pointers to null-terminated key strings
   */
  char **keys;
  /*
   Array of pointers to null-terminated value strings
   */
  char **values;
  /*
   The number of key-value pairs in the arrays
   */
  size_t len;
} CMetadata;

typedef struct CStreamReader {
  uint8_t _private[0];
} CStreamReader;

/*
 Returns the FFI wrapper version as a null-terminated UTF-8 string.
 The returned pointer is to a static string and must not be freed.
 */
const char *extractous_ffi_version(void);

/*
 Returns the underlying Extractous core library version.
 The returned pointer is to a static string and must not be freed.
 */
const char *extractous_core_version(void);

/*
 Creates a new PDF parser configuration with default settings.
 The returned handle must be freed with `extractous_pdf_config_free()`
 unless passed to an extractor, which will take ownership.
 */
auto struct CPdfParserConfig *extractous_pdf_config_new(void);

/*
 Frees the memory associated with a PDF parser configuration.
 Do not call this if the config has been attached to an extractor.
 */
void extractous_pdf_config_free(struct CPdfParserConfig *handle);

/*
 Sets the OCR strategy for PDF parsing. Modifies the config in-place.
 */
void extractous_pdf_config_set_ocr_strategy(struct CPdfParserConfig *handle, int strategy);

/*
 Enables or disables extraction of inline images. Modifies the config in-place.
 */
void extractous_pdf_config_set_extract_inline_images(struct CPdfParserConfig *handle, bool value);

/*
 If enabled, only unique inline images (by digest) will be extracted.
 */
void extractous_pdf_config_set_extract_unique_inline_images_only(struct CPdfParserConfig *handle,
                                                                 bool value);

/*
 Enables or disables extraction of text from marked content sections.
 */
void extractous_pdf_config_set_extract_marked_content(struct CPdfParserConfig *handle, bool value);

/*
 Enables or disables extraction of text from annotations.
 */
void extractous_pdf_config_set_extract_annotation_text(struct CPdfParserConfig *handle, bool value);

/*
 Creates a new Office parser configuration with default settings.
 */
auto struct COfficeParserConfig *extractous_office_config_new(void);

/*
 Frees the memory associated with an Office parser configuration.
 */
void extractous_office_config_free(struct COfficeParserConfig *handle);

/*
 Enables or disables macro extraction. Modifies the config in-place.
 */
void extractous_office_config_set_extract_macros(struct COfficeParserConfig *handle, bool value);

/*
 Enables or disables inclusion of deleted content (track changes).
 */
void extractous_office_config_set_include_deleted_content(struct COfficeParserConfig *handle,
                                                          bool value);

/*
 Enables or disables inclusion of moved-from content (track changes).
 */
void extractous_office_config_set_include_move_from_content(struct COfficeParserConfig *handle,
                                                            bool value);

/*
 Enables or disables inclusion of content from shapes.
 */
void extractous_office_config_set_include_shape_based_content(struct COfficeParserConfig *handle,
                                                              bool value);

/*
 Creates a new Tesseract OCR configuration with default settings.
 */
auto struct CTesseractOcrConfig *extractous_ocr_config_new(void);

/*
 Frees the memory associated with a Tesseract OCR configuration.
 */
void extractous_ocr_config_free(struct CTesseractOcrConfig *handle);

/*
 Sets the OCR language. Modifies the config in-place.
 */
void extractous_ocr_config_set_language(struct CTesseractOcrConfig *handle, const char *language);

/*
 Sets the DPI for OCR processing. Modifies the config in-place.
 */
void extractous_ocr_config_set_density(struct CTesseractOcrConfig *handle, int32_t density);

/*
 Sets the bit depth for OCR processing.
 */
void extractous_ocr_config_set_depth(struct CTesseractOcrConfig *handle, int32_t depth);

/*
 Enables or disables image preprocessing for OCR.
 */
void extractous_ocr_config_set_enable_image_preprocessing(struct CTesseractOcrConfig *handle,
                                                          bool value);

/*
 Sets the timeout for the Tesseract process in seconds.
 */
void extractous_ocr_config_set_timeout_seconds(struct CTesseractOcrConfig *handle, int32_t seconds);

char *extractous_error_message(int code);

/*
 Retrieves a detailed debug report for the last error on this thread
 full error chain and a backtrace if RUST_BACKTRACE=1
 */
char *extractous_error_get_last_debug(void);

/*
 Checks if debug information is available for the current thread
 */
int extractous_error_has_debug(void);

void extractous_error_clear_last(void);

/*
 Creates a new `Extractor` with a default configuration.
 The returned handle must be freed with `extractous_extractor_free`.
 */
auto struct CExtractor *extractous_extractor_new(void);

/*
 Frees the memory associated with an `Extractor` handle.
 */
void extractous_extractor_free(struct CExtractor *handle);

/*
 Sets the maximum length for extracted string content.
 */
void extractous_extractor_set_extract_string_max_length_mut(struct CExtractor *handle,
                                                            int max_length);

/*
 Sets the character encoding for the extracted text.
 */
void extractous_extractor_set_encoding_mut(struct CExtractor *handle, int encoding);

/*
 Sets the configuration for the PDF parser.
 */
void extractous_extractor_set_pdf_config_mut(struct CExtractor *handle,
                                             const struct CPdfParserConfig *config);

/*
 Sets the configuration for the Office document parser.
 */
void extractous_extractor_set_office_config_mut(struct CExtractor *handle,
                                                const struct COfficeParserConfig *config);

/*
 Sets the configuration for Tesseract OCR.
 */
void extractous_extractor_set_ocr_config_mut(struct CExtractor *handle,
                                             const struct CTesseractOcrConfig *config);

/*
 Sets whether to output structured XML instead of plain text.
 */
void extractous_extractor_set_xml_output_mut(struct CExtractor *handle, bool xml_output);

/*
 Extracts content and metadata from a local file path into a string.

 Output strings must be freed with `extractous_string_free`.
 Output metadata must be freed with `extractous_metadata_free`.
 */
int extractous_extractor_extract_file_to_string(struct CExtractor *handle,
                                                const char *path,
                                                char **out_content,
                                                struct CMetadata **out_metadata);

/*
 Extracts content and metadata from a local file path into a stream.
 */
int extractous_extractor_extract_file(struct CExtractor *handle,
                                      const char *path,
                                      struct CStreamReader **out_reader,
                                      struct CMetadata **out_metadata);

/*
 Extracts content and metadata from a byte slice into a string.
 */
int extractous_extractor_extract_bytes_to_string(struct CExtractor *handle,
                                                 const uint8_t *data,
                                                 size_t data_len,
                                                 char **out_content,
                                                 struct CMetadata **out_metadata);

/*
 Extracts content and metadata from a byte slice into a stream.
 */
int extractous_extractor_extract_bytes(struct CExtractor *handle,
                                       const uint8_t *data,
                                       size_t data_len,
                                       struct CStreamReader **out_reader,
                                       struct CMetadata **out_metadata);

/*
 Extracts content and metadata from a URL into a string.
 */
int extractous_extractor_extract_url_to_string(struct CExtractor *handle,
                                               const char *url,
                                               char **out_content,
                                               struct CMetadata **out_metadata);

/*
 Extracts content and metadata from a URL into a stream.
 */
int extractous_extractor_extract_url(struct CExtractor *handle,
                                     const char *url,
                                     struct CStreamReader **out_reader,
                                     struct CMetadata **out_metadata);

/*
 Frees a C-style string that was allocated by this library.
 */
void extractous_string_free(char *s);

/*
 Frees a metadata structure and all associated memory.
 */
void extractous_metadata_free(struct CMetadata *metadata);

/*
 Reads data from a stream into a user-provided buffer.

 Returns the actual number of bytes read via the `bytes_read` output parameter.
 Reaching the end of the stream is indicated by `ERR_OK` and `*bytes_read == 0`.
 */
int extractous_stream_read(struct CStreamReader *handle,
                           uint8_t *buffer,
                           size_t buffer_size,
                           size_t *bytes_read);

/*
 Reads exactly `buffer_size` bytes from the stream.

 Function will continue reading until the buffer is full, or the end of
 the stream is reached, or an error occurs.
 */
int extractous_stream_read_exact(struct CStreamReader *handle,
                                 uint8_t *buffer,
                                 size_t buffer_size,
                                 size_t *bytes_read);

/*
 Reads the remaining stream into a newly allocated buffer.
 */
auto
int extractous_stream_read_all(struct CStreamReader *handle,
                               uint8_t **out_buffer,
                               size_t *out_size);

/*
 Frees a buffer allocated by `extractous_stream_read_all`.
 */
void extractous_buffer_free(uint8_t *buffer, size_t size);

/*
 Frees a stream reader and releases its resources.
 */
void extractous_stream_free(struct CStreamReader *handle);

#endif  /* EXTRACTOUS_H */
