#include <cstdarg>
#include <cstdint>
#include <cstdlib>
#include <ostream>
#include <new>

/// Success - operation completed without errors
///
/// This is the only non-error return value. All operations that complete
/// successfully will return this code.
constexpr static const int ERR_OK = 0;

/// Null pointer provided as argument
///
/// Returned when a required pointer argument is NULL.
/// Check all pointer arguments before calling FFI functions.
///
/// Common causes:
/// - Forgot to allocate output parameter
/// - Accidentally passed NULL instead of valid pointer
/// - Double-free caused pointer to become invalid
constexpr static const int ERR_NULL_POINTER = -1;

/// Invalid UTF-8 string encoding
///
/// Returned when a C string argument contains invalid UTF-8 sequences.
/// All string arguments must be valid UTF-8.
///
/// Common causes:
/// - String contains binary data
/// - Wrong encoding used (e.g., Latin-1 instead of UTF-8)
/// - Corrupted string data
constexpr static const int ERR_INVALID_UTF8 = -2;

/// String conversion or allocation failed
///
/// Returned when internal string operations fail, typically due to:
/// - Null bytes in unexpected positions
/// - Memory allocation failure
/// - String contains invalid characters for the operation
constexpr static const int ERR_INVALID_STRING = -3;

/// Document extraction failed
///
/// General extraction error when the specific cause is unknown or internal.
/// The document may be:
/// - Corrupted or malformed
/// - Using an unsupported format variant
/// - Encrypted without proper credentials
/// - Too complex for the parser
constexpr static const int ERR_EXTRACTION_FAILED = -4;

/// File system or network I/O error
///
/// Returned when file or network operations fail.
///
/// Common causes:
/// - File not found
/// - Permission denied
/// - Network timeout
/// - Disk full
/// - Path too long
constexpr static const int ERR_IO_ERROR = -5;

/// Invalid configuration value
///
/// Returned when configuration parameters are invalid.
///
/// Common causes:
/// - Out of range values
/// - Incompatible configuration combinations
/// - Invalid enum constants
constexpr static const int ERR_INVALID_CONFIG = -6;

/// Invalid enumeration value
///
/// Returned when an enum constant (like charset or OCR strategy) is invalid.
/// Only use the documented constants.
constexpr static const int ERR_INVALID_ENUM = -7;

/// Unsupported file format
///
/// The file format is not supported by extractous or the parser
/// for this format is not available.
constexpr static const int ERR_UNSUPPORTED_FORMAT = -8;

/// Memory allocation failed
///
/// Extremely rare - indicates the system is out of memory.
constexpr static const int ERR_OUT_OF_MEMORY = -9;

/// OCR operation failed
///
/// OCR processing failed, possibly because:
/// - Tesseract is not installed
/// - Invalid language data
/// - Image format not supported
constexpr static const int ERR_OCR_FAILED = -10;

/// UTF-8 encoding (default, recommended)
///
/// Universal character encoding supporting all languages and emojis.
/// This is the default and recommended encoding for most use cases.
constexpr static const int CHARSET_UTF_8 = 0;

/// US-ASCII encoding
///
/// 7-bit ASCII encoding. Use only if you're certain the content
/// contains only basic ASCII characters (0-127).
constexpr static const int CHARSET_US_ASCII = 1;

/// UTF-16 Big Endian encoding
///
/// 16-bit Unicode encoding with big-endian byte order.
/// Less common, use only if specifically required.
constexpr static const int CHARSET_UTF_16BE = 2;

/// No OCR processing - extract only embedded text
///
/// Fastest option. Extracts only text that is already present in the PDF.
/// Images and scanned pages will not be processed.
///
/// Use when:
/// - PDF contains searchable text
/// - OCR is not needed
/// - Performance is critical
constexpr static const int PDF_OCR_NO_OCR = 0;

/// OCR only - ignore embedded text
///
/// Renders pages as images and performs OCR.
/// Ignores any embedded text in the PDF.
///
/// Use when:
/// - PDF text layer is corrupted or unreliable
/// - You need consistent OCR processing
constexpr static const int PDF_OCR_OCR_ONLY = 1;

/// Combined OCR and text extraction
///
/// Extracts embedded text AND performs OCR on images.
/// Provides most comprehensive extraction but is slower.
///
/// Use when:
/// - PDF has both text and scanned images
/// - Maximum content extraction is needed
constexpr static const int PDF_OCR_OCR_AND_TEXT_EXTRACTION = 2;

/// Automatic OCR strategy selection
///
/// Analyzes the PDF and automatically decides whether to use OCR.
/// Good balance between performance and coverage.
///
/// Use when:
/// - Processing mixed PDFs (some with text, some scanned)
/// - Want automatic optimization
constexpr static const int PDF_OCR_AUTO = 3;

/// Default buffer size for stream reading (4KB)
///
/// Recommended buffer size for efficient stream reading.
/// Balances memory usage and I/O performance.
constexpr static const size_t DEFAULT_BUFFER_SIZE = 4096;

/// Maximum recommended buffer size (1MB)
///
/// Large buffer for high-performance scenarios.
/// Use when processing very large documents.
constexpr static const size_t MAX_BUFFER_SIZE = (1024 * 1024);

/// Default string extraction limit (100MB)
///
/// Default maximum length for extracted strings to prevent
/// excessive memory usage on very large documents.
constexpr static const int DEFAULT_STRING_MAX_LENGTH = ((100 * 1024) * 1024);

/// Opaque handle to a PdfParserConfig instance
///
/// Configuration for PDF document parsing. Create with `extractous_pdf_config_new()`,
/// configure with setter functions, and free with `extractous_pdf_config_free()`.
///
/// Note: Setters consume the old handle and return a new one (builder pattern).
struct CPdfParserConfig {
  uint8_t _private[0];
};

/// Opaque handle to an OfficeParserConfig instance
///
/// Configuration for Microsoft Office document parsing. Create with
/// `extractous_office_config_new()` and free with `extractous_office_config_free()`.
struct COfficeParserConfig {
  uint8_t _private[0];
};

/// Opaque handle to a TesseractOcrConfig instance
///
/// Configuration for Tesseract OCR engine. Create with `extractous_ocr_config_new()`
/// and free with `extractous_ocr_config_free()`.
///
/// Note: Requires Tesseract to be installed on the system.
struct CTesseractOcrConfig {
  uint8_t _private[0];
};

/// Opaque handle to an Extractor instance
///
/// Represents the main extraction engine. Create with `extractous_extractor_new()`
/// and destroy with `extractous_extractor_free()`.
///
/// ## Thread Safety
/// Not thread-safe. Use separate instances per thread or external synchronization.
///
/// ### Example
/// ```c
/// CExtractor* extractor = extractous_extractor_new();
/// // ... use extractor ...
/// extractous_extractor_free(extractor);
/// ```
struct CExtractor {
  uint8_t _private[0];
};

/// C-compatible metadata structure
///
/// Contains document metadata as parallel arrays of keys and values.
/// Multiple values for the same key are comma-separated.
///
/// ### Memory Layout
/// ```text
/// keys[0] -> "author\0"      values[0] -> "John Doe\0"
/// keys[1] -> "title\0"       values[1] -> "My Document\0"
/// keys[2] -> "keywords\0"    values[2] -> "pdf,test,sample\0"
/// ```
///
/// ### Memory Management
/// Must be freed with `extractous_metadata_free()` which will:
/// 1. Free all individual key strings
/// 2. Free all individual value strings
/// 3. Free the key array
/// 4. Free the value array
/// 5. Free the structure itself
///
/// ### Safety
/// - All strings are valid null-terminated UTF-8
/// - Arrays contain exactly `len` elements
/// - Do not modify the structure directly from C
/// - Do not free individual strings; use `extractous_metadata_free()`
struct CMetadata {
  /// Array of pointers to key strings (null-terminated UTF-8)
  char **keys;
  /// Array of pointers to value strings (null-terminated UTF-8, comma-separated if multiple)
  char **values;
  /// Number of key-value pairs in the arrays
  size_t len;
};

/// Opaque handle to a StreamReader instance
///
/// Represents a buffered stream of extracted content. Read data using
/// `extractous_stream_read()` and free with `extractous_stream_free()`.
///
/// ### Example
/// ```c
/// CStreamReader* reader;
/// CMetadata* metadata;
/// extractous_extractor_extract_file(extractor, "doc.pdf", &reader, &metadata);
///
/// char buffer[4096];
/// size_t bytes_read;
/// while (extractous_stream_read(reader, buffer, sizeof(buffer), &bytes_read) == ERR_OK
///        && bytes_read > 0) {
///     // Process buffer...
/// }
/// extractous_stream_free(reader);
/// ```
struct CStreamReader {
  uint8_t _private[0];
};

extern "C" {

/// Returns the FFI wrapper version in semver format.
const char *extractous_ffi_version();

/// Returns the underlying Extractous core library version.
const char *extractous_core_version();

/// Create a new PDF parser configuration with default settings
///
/// ### Default configuration:
/// - OCR strategy: NO_OCR (fastest, text extraction only)
/// - Extract inline images: false
/// - Extract unique inline images only: true
/// - Extract marked content: false
/// - Extract annotation text: false
///
/// Returns
/// Pointer to new PdfParserConfig. Must be freed with `extractous_pdf_config_free()`
/// unless attached to an extractor.
///
/// ```c
/// CPdfParserConfig* config = extractous_pdf_config_new();
/// if (config == NULL) {
///     // Handle allocation error
/// }
/// ```
CPdfParserConfig *extractous_pdf_config_new();

/// Set the OCR strategy for PDF parsing
///
/// Determines how OCR is applied to PDF documents.
///
/// ##### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `strategy` - PDF_OCR_NO_OCR, PDF_OCR_OCR_ONLY, PDF_OCR_OCR_AND_TEXT_EXTRACTION, PDF_OCR_AUTO
///
/// ##### Returns
/// New PdfParserConfig handle with updated strategy, or NULL if invalid.
/// The input handle is consumed and must not be used.
///
/// ### Strategy Guide
/// - `PDF_OCR_NO_OCR`: Fastest, text-based PDFs only
/// - `PDF_OCR_OCR_ONLY`: Scanned documents, ignore existing text
/// - `PDF_OCR_OCR_AND_TEXT_EXTRACTION`: Mixed content, thorough extraction
/// - `PDF_OCR_AUTO`: Let the library decide (recommended)
///
/// ##### Safety
/// - Input handle is consumed; do not use after this call
/// - Returns NULL if handle is NULL or strategy is invalid
CPdfParserConfig *extractous_pdf_config_set_ocr_strategy(CPdfParserConfig *handle, int strategy);

/// Enable or disable extraction of inline images from PDF
///
/// When enabled, extracts embedded image data from the PDF.
/// Can significantly increase memory usage and processing time.
///
/// ##### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `value` - true to extract inline images, false to skip
///
/// ##### Returns
/// New PdfParserConfig handle or NULL if handle is invalid.
///
/// ### Performance Impact
/// - Disabled (default): Fast, minimal memory
/// - Enabled: Slower, higher memory usage
///
/// ##### Safety
/// Input handle is consumed; do not use after this call.
CPdfParserConfig *extractous_pdf_config_set_extract_inline_images(CPdfParserConfig *handle,
                                                                  bool value);

/// Extract each unique inline image only once
///
/// When enabled with inline image extraction, deduplicates repeated images.
///
/// ### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `value` - true for deduplication (recommended), false to extract all
///
/// ### Returns
/// New PdfParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
CPdfParserConfig *extractous_pdf_config_set_extract_unique_inline_images_only(CPdfParserConfig *handle,
                                                                              bool value);

/// Extract text with marked content structure
///
/// Attempts to preserve document structure markers from the PDF.
///
/// ### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `value` - true to extract marked content, false otherwise
///
/// ### Returns
/// New PdfParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
CPdfParserConfig *extractous_pdf_config_set_extract_marked_content(CPdfParserConfig *handle,
                                                                   bool value);

/// Extract text from PDF annotations
///
/// Includes comments, highlights, and other annotation content.
///
/// ### Arguments
/// * `handle` - Valid PdfParserConfig pointer
/// * `value` - true to extract annotations, false to skip
///
/// ### Returns
/// New PdfParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
CPdfParserConfig *extractous_pdf_config_set_extract_annotation_text(CPdfParserConfig *handle,
                                                                    bool value);

/// Free PDF parser configuration
///
/// ### Safety
/// - `handle` must be a valid PdfParserConfig pointer
/// - `handle` must not be used after this call
/// - Do not call this if config was attached to an extractor (it will be freed automatically)
///
/// ### Example
/// ```c
/// CPdfParserConfig* config = extractous_pdf_config_new();
/// // Use config...
/// extractous_pdf_config_free(config);  // Only if not attached to extractor
/// ```
void extractous_pdf_config_free(CPdfParserConfig *handle);

/// Create a new Office parser configuration with default settings
///
/// Default configuration:
/// - Extract macros: false
/// - Include deleted content: false
/// - Include move-from content: false
/// - Include shape-based content: true
///
/// ### Returns
/// Pointer to new OfficeParserConfig. Must be freed with `extractous_office_config_free()`
/// unless attached to an extractor.
COfficeParserConfig *extractous_office_config_new();

/// Enable or disable macro extraction from Office documents
///
/// **Security Warning**: Macros can contain malicious code. Only enable this
/// if you trust the document source and need macro content.
///
/// ### Arguments
/// * `handle` - Valid OfficeParserConfig pointer
/// * `value` - true to extract macros (security risk), false to skip (safer)
///
/// ### Returns
/// New OfficeParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
COfficeParserConfig *extractous_office_config_set_extract_macros(COfficeParserConfig *handle,
                                                                 bool value);

/// Include deleted content from DOCX track changes
///
/// When enabled, extracts text that was deleted but is still present in
/// the document's revision history.
///
/// ### Arguments
/// * `handle` - Valid OfficeParserConfig pointer
/// * `value` - true to include deleted text, false to skip
///
/// ### Returns
/// New OfficeParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
COfficeParserConfig *extractous_office_config_set_include_deleted_content(COfficeParserConfig *handle,
                                                                          bool value);

/// Include "move-from" content in DOCX documents
///
/// Extracts text that was moved from one location to another during editing.
///
/// ### Arguments
/// * `handle` - Valid OfficeParserConfig pointer
/// * `value` - true to include moved text, false to skip
///
/// ### Returns
/// New OfficeParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
COfficeParserConfig *extractous_office_config_set_include_move_from_content(COfficeParserConfig *handle,
                                                                            bool value);

/// Include text from drawing shapes and text boxes
///
/// When enabled, extracts text from shapes, text boxes, and other drawing objects.
///
/// ### Arguments
/// * `handle` - Valid OfficeParserConfig pointer
/// * `value` - true to include shape text (recommended), false to skip
///
/// ### Returns
/// New OfficeParserConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
COfficeParserConfig *extractous_office_config_set_include_shape_based_content(COfficeParserConfig *handle,
                                                                              bool value);

/// Free Office parser configuration
///
/// ### Safety
/// - `handle` must be valid and not used after this call
/// - Do not call this if config was attached to an extractor
void extractous_office_config_free(COfficeParserConfig *handle);

/// Create a new Tesseract OCR configuration with default settings
///
/// Default configuration:
/// - Language: "eng" (English)
/// - Density: 300 DPI
/// - Depth: 32 bits
/// - Image preprocessing: true
/// - Timeout: 300 seconds
///
/// ### Prerequisites
/// Tesseract must be installed on the system with appropriate language data files.
///
/// ### Returns
/// Pointer to new TesseractOcrConfig. Must be freed with `extractous_ocr_config_free()`
/// unless attached to an extractor.
CTesseractOcrConfig *extractous_ocr_config_new();

/// Set the OCR language
///
/// Specifies which language(s) Tesseract should use for recognition.
/// Multiple languages can be specified with '+' separator (e.g., "eng+fra").
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `language` - Null-terminated UTF-8 language code (e.g., "eng", "deu", "eng+fra")
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle or language is invalid.
///
/// ### Common Language Codes
/// - "eng" - English
/// - "deu" - German
/// - "fra" - French
/// - "spa" - Spanish
///
/// ### Requirements
/// The specified language data must be installed on the system.
/// On Debian/Ubuntu: `apt install tesseract-ocr-[lang]`
///
/// ### Safety
/// Input handle is consumed. Language string must be valid UTF-8.
CTesseractOcrConfig *extractous_ocr_config_set_language(CTesseractOcrConfig *handle,
                                                        const char *language);

/// Set the DPI (dots per inch) for OCR processing
///
/// Higher DPI values can improve accuracy but increase processing time.
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `density` - DPI value (recommended: 150-600, default: 300)
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle is invalid.
///
/// ### Recommendations
/// - 150 DPI: Fast, lower quality
/// - 300 DPI: Balanced (default)
/// - 400-600 DPI: High quality, slower
///
/// ### Safety
/// Input handle is consumed.
CTesseractOcrConfig *extractous_ocr_config_set_density(CTesseractOcrConfig *handle,
                                                       int32_t density);

/// Set the color depth for OCR processing
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `depth` - Bit depth (typically 8, 24, or 32)
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
CTesseractOcrConfig *extractous_ocr_config_set_depth(CTesseractOcrConfig *handle, int32_t depth);

/// Enable or disable image preprocessing for OCR
///
/// Preprocessing can improve OCR accuracy by normalizing image quality,
/// adjusting contrast, removing noise, etc.
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `value` - true to enable preprocessing (recommended), false to disable
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle is invalid.
///
/// ### Safety
/// Input handle is consumed.
CTesseractOcrConfig *extractous_ocr_config_set_enable_image_preprocessing(CTesseractOcrConfig *handle,
                                                                          bool value);

/// Set timeout for OCR processing
///
/// Prevents OCR from running indefinitely on problematic images.
///
/// ### Arguments
/// * `handle` - Valid TesseractOcrConfig pointer
/// * `seconds` - Timeout in seconds (0 = no timeout, default: 300)
///
/// ### Returns
/// New TesseractOcrConfig handle or NULL if handle is invalid.
///
/// # Recommendations
/// - 60-120 seconds: Fast processing, may timeout on complex images
/// - 300 seconds: Default, handles most documents
/// - 600+ seconds: Very complex documents
///
/// ### Safety
/// Input handle is consumed.
CTesseractOcrConfig *extractous_ocr_config_set_timeout_seconds(CTesseractOcrConfig *handle,
                                                               int32_t seconds);

/// Free Tesseract OCR configuration
///
/// ### Safety
/// - `handle` must be valid and not used after this call
/// - Do not call this if config was attached to an extractor
void extractous_ocr_config_free(CTesseractOcrConfig *handle);

/// Get human-readable error message for error code
///
/// Returns a newly allocated string containing the error description.
/// The caller must free the returned string with `extractous_string_free()`.
///
/// ### Arguments
/// * `code` - Error code returned by an extractous function
///
/// ### Returns
/// Pointer to null-terminated UTF-8 string, or NULL if allocation fails.
///
/// ### Example
/// ```c
/// int err = extractous_extractor_extract_file(...);
/// if (err != ERR_OK) {
///     char* msg = extractous_error_message(err);
///     printf("Error: %s\n", msg);
///     extractous_string_free(msg);
/// }
/// ```
///
/// ### Safety
/// - Return value must be freed with `extractous_string_free()`
/// - Do not modify the returned string
char *extractous_error_message(int code);

/// Get error category description
///
/// Returns a high-level categorization of the error.
///
/// ### Arguments
/// * `code` - Error code
///
/// ### Returns
/// Static string pointer (do not free)
///
/// ### Safety
/// Return value points to static memory and must not be freed.
const char *extractous_error_category(int code);

/// Create a new Extractor with default configuration
///
/// ### Returns
/// Pointer to new Extractor, or NULL on failure.
/// Must be freed with `extractous_extractor_free`.
CExtractor *extractous_extractor_new();

/// Free an Extractor instance
///
/// ### Safety
/// - `handle` must be a valid pointer returned by `extractous_extractor_new`
/// - `handle` must not be used after this call
/// - Calling this twice on the same pointer causes undefined behavior
void extractous_extractor_free(CExtractor *handle);

/// Set maximum length for extracted string content
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - Returns a NEW handle; old handle is consumed and must not be used
///
/// ### Returns
/// New Extractor handle with updated config, or NULL on error.
CExtractor *extractous_extractor_set_extract_string_max_length(CExtractor *handle, int max_length);

/// Set character encoding for extraction
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - `encoding` must be a valid CHARSET_* constant
/// - Returns a NEW handle; old handle is consumed
///
/// ### Returns
/// New Extractor handle, or NULL if encoding is invalid.
CExtractor *extractous_extractor_set_encoding(CExtractor *handle, int encoding);

/// Set PDF parser configuration
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - `config` must be a valid PdfParserConfig pointer
/// - Returns a NEW handle; old handle is consumed
CExtractor *extractous_extractor_set_pdf_config(CExtractor *handle, CPdfParserConfig *config);

/// Set Office parser configuration
///
/// ### Safety
/// Same safety requirements as `extractous_extractor_set_pdf_config`.
CExtractor *extractous_extractor_set_office_config(CExtractor *handle, COfficeParserConfig *config);

/// Set OCR configuration
///
/// ### Safety
/// Same safety requirements as `extractous_extractor_set_pdf_config`.
CExtractor *extractous_extractor_set_ocr_config(CExtractor *handle, CTesseractOcrConfig *config);

/// Set whether to output XML structure
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - Returns a NEW handle; old handle is consumed
CExtractor *extractous_extractor_set_xml_output(CExtractor *handle, bool xml_output);

/// Extract file content to string
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - `path` must be a valid null-terminated UTF-8 string
/// - `out_content` and `out_metadata` must be valid pointers
/// - Caller must free returned content with `extractous_string_free`
/// - Caller must free returned metadata with `extractous_metadata_free`
///
/// ### Returns
/// ERR_OK on success, error code on failure.
int extractous_extractor_extract_file_to_string(CExtractor *handle,
                                                const char *path,
                                                char **out_content,
                                                CMetadata **out_metadata);

/// Extract file content to stream
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - `path` must be a valid null-terminated UTF-8 string
/// - `out_reader` and `out_metadata` must be valid pointers
/// - Caller must free returned reader with `extractous_stream_free`
/// - Caller must free returned metadata with `extractous_metadata_free`
///
/// ### Returns
/// ERR_OK on success, error code on failure.
int extractous_extractor_extract_file(CExtractor *handle,
                                      const char *path,
                                      CStreamReader **out_reader,
                                      CMetadata **out_metadata);

/// Extract from byte array to string
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - `data` must point to at least `data_len` valid bytes
/// - `out_content` and `out_metadata` must be valid pointers
///
/// ### Returns
/// ERR_OK on success, error code on failure.
int extractous_extractor_extract_bytes_to_string(CExtractor *handle,
                                                 const uint8_t *data,
                                                 size_t data_len,
                                                 char **out_content,
                                                 CMetadata **out_metadata);

/// Extract from byte array to stream
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - `data` must point to at least `data_len` valid bytes
/// - `out_reader` and `out_metadata` must be valid pointers
///
/// ### Returns
/// ERR_OK on success, error code on failure.
int extractous_extractor_extract_bytes(CExtractor *handle,
                                       const uint8_t *data,
                                       size_t data_len,
                                       CStreamReader **out_reader,
                                       CMetadata **out_metadata);

/// Free a string allocated by Rust
///
/// ### Safety
/// - `s` must be a pointer returned by an extractous function
/// - `s` must not be used after this call
/// - Calling this twice on the same pointer causes undefined behavior
void extractous_string_free(char *s);

/// Extract URL content to string
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - `url` must be a valid null-terminated UTF-8 string
/// - `out_content` and `out_metadata` must be valid pointers
/// - Caller must free returned content with `extractous_string_free`
/// - Caller must free returned metadata with `extractous_metadata_free`
///
/// ### Returns
/// ERR_OK on success, error code on failure.
int extractous_extractor_extract_url_to_string(CExtractor *handle,
                                               const char *url,
                                               char **out_content,
                                               CMetadata **out_metadata);

/// Extract URL content to stream
///
/// ### Safety
/// - `handle` must be a valid Extractor pointer
/// - `url` must be a valid null-terminated UTF-8 string
/// - `out_reader` and `out_metadata` must be valid pointers
/// - Caller must free returned reader with `extractous_stream_free`
/// - Caller must free returned metadata with `extractous_metadata_free`
///
/// ### Returns
/// ERR_OK on success, error code on failure.
int extractous_extractor_extract_url(CExtractor *handle,
                                     const char *url,
                                     CStreamReader **out_reader,
                                     CMetadata **out_metadata);

/// Free a metadata structure and all associated memory.
///
/// Frees:
/// 1. All individual key strings
/// 2. All individual value strings
/// 3. The key array
/// 4. The value array
/// 5. The `CMetadata` structure itself
///
/// ### Arguments
/// * `metadata` - Pointer to a `CMetadata` structure to free
///
/// ### Safety
/// - `metadata` must be a pointer returned by an extraction function
/// - `metadata` must not be used after this call
/// - Safe to call with NULL (no-op)
///
/// ### Example
/// ```c
/// CMetadata* metadata;
/// extractous_extractor_extract_file_to_string(extractor, path, &content, &metadata);
/// // ... use metadata ...
/// extractous_metadata_free(metadata);
/// ```
void extractous_metadata_free(CMetadata *metadata);

/// Read data from stream into buffer
///
/// Reads up to `buffer_size` bytes from the stream into the provided buffer.
/// Returns the actual number of bytes read in `bytes_read`.
///
/// ### Arguments
/// * `handle` - Valid StreamReader pointer
/// * `buffer` - Pointer to buffer where data will be written
/// * `buffer_size` - Size of the buffer in bytes
/// * `bytes_read` - Output pointer for number of bytes actually read (can be NULL)
///
/// ### Returns
/// * `ERR_OK` - Read successful (check bytes_read for amount)
/// * `ERR_NULL_POINTER` - handle or buffer is NULL
/// * `ERR_IO_ERROR` - Read operation failed
///
/// ### End of Stream
/// When end of stream is reached:
/// - Function returns `ERR_OK`
/// - `bytes_read` is set to 0
///
/// ### Example
/// ```c
/// char buffer[4096];
/// size_t n;
/// int result = extractous_stream_read(reader, buffer, sizeof(buffer), &n);
///
/// if (result == ERR_OK) {
///     if (n > 0) {
///         // Process n bytes of data in buffer
///     } else {
///         // End of stream reached
///     }
/// } else {
///     // Error occurred
/// }
/// ```
///
/// ### Safety
/// - `handle` must be a valid StreamReader pointer
/// - `buffer` must point to at least `buffer_size` writable bytes
/// - `bytes_read` must be NULL or point to valid size_t
/// - Buffer content is undefined if function returns error
int extractous_stream_read(CStreamReader *handle,
                           uint8_t *buffer,
                           size_t buffer_size,
                           size_t *bytes_read);

/// Read exactly the requested number of bytes or fail
///
/// Similar to `extractous_stream_read()` but guarantees that exactly
/// `buffer_size` bytes are read unless end of stream or error occurs.
///
/// ### Arguments
/// * `handle` - Valid StreamReader pointer
/// * `buffer` - Pointer to buffer where data will be written
/// * `buffer_size` - Exact number of bytes to read
/// * `bytes_read` - Output pointer for number of bytes actually read (can be NULL)
///
/// ### Returns
/// * `ERR_OK` - Successfully read exactly `buffer_size` bytes
/// * `ERR_OK` with bytes_read < buffer_size - End of stream reached
/// * `ERR_NULL_POINTER` - handle or buffer is NULL
/// * `ERR_IO_ERROR` - Read operation failed
///
/// ### Use Cases
/// - Reading fixed-size headers or chunks
/// - When partial reads are not acceptable
///
/// ### Safety
/// - Same safety requirements as `extractous_stream_read()`
int extractous_stream_read_exact(CStreamReader *handle,
                                 uint8_t *buffer,
                                 size_t buffer_size,
                                 size_t *bytes_read);

/// Read entire remaining stream into a newly allocated buffer
///
/// Reads all remaining data from the stream and returns it in a newly
/// allocated buffer. Useful for smaller streams where you want all data at once.
///
/// **Warning**: This loads all content into memory. For large documents,
/// prefer using `extractous_stream_read()` in a loop.
///
/// ### Arguments
/// * `handle` - Valid StreamReader pointer
/// * `out_buffer` - Output pointer for allocated buffer
/// * `out_size` - Output pointer for buffer size
///
/// ### Returns
/// * `ERR_OK` - Success, buffer allocated and filled
/// * `ERR_NULL_POINTER` - Invalid pointer argument
/// * `ERR_IO_ERROR` - Read operation failed
/// * `ERR_OUT_OF_MEMORY` - Memory allocation failed
///
/// ### Memory Management
/// The caller must free the returned buffer using `extractous_buffer_free()`.
///
/// ### Example
/// ```c
/// uint8_t* data;
/// size_t size;
/// int result = extractous_stream_read_all(reader, &data, &size);
///
/// if (result == ERR_OK) {
///     // Use data (size bytes)
///     process_data(data, size);
///
///     // Free the buffer
///     extractous_buffer_free(data);
/// }
/// ```
///
/// ### Safety
/// - `handle` must be a valid StreamReader pointer
/// - `out_buffer` and `out_size` must be valid pointers
/// - Returned buffer must be freed with `extractous_buffer_free()`
int extractous_stream_read_all(CStreamReader *handle, uint8_t **out_buffer, size_t *out_size);

/// Free a buffer allocated by extractous_stream_read_all
///
/// ### Arguments
/// * `buffer` - Pointer to buffer returned by `extractous_stream_read_all()`
/// * `size` - Size of the buffer in bytes
///
/// ### Safety
/// - `buffer` must be a pointer returned by `extractous_stream_read_all()`
/// - `size` must match the size returned by that function
/// - `buffer` must not be used after this call
/// - Do not call this function twice on the same buffer
///
/// ### Example
/// ```c
/// uint8_t* data;
/// size_t size;
/// extractous_stream_read_all(reader, &data, &size);
/// // ... use data ...
/// extractous_buffer_free(data, size);
/// ```
void extractous_buffer_free(uint8_t *buffer, size_t size);

/// Free stream reader and release associated resources
///
/// ### Arguments
/// * `handle` - StreamReader pointer to free
///
/// ### Safety
/// - `handle` must be a valid StreamReader pointer
/// - `handle` must not be used after this call
/// - Calling this function twice on the same pointer causes undefined behavior
/// - Safe to call with NULL (no-op)
///
/// ### Example
/// ```c
/// CStreamReader* reader;
/// // ... extract to stream and use reader ...
/// extractous_stream_free(reader);  // Always free when done
/// ```
void extractous_stream_free(CStreamReader *handle);

}  // extern "C"
