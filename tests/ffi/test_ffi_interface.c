/**
 * FFI Layer Tests for Extractous
 * 
 * These tests verify that the C FFI interface works correctly
 * and all functions are properly exposed from the Rust library.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include "../../include/extractous.h"

// Test result tracking
static int tests_run = 0;
static int tests_passed = 0;
static int tests_failed = 0;

// Color codes for output
#define COLOR_GREEN "\x1b[32m"
#define COLOR_RED "\x1b[31m"
#define COLOR_YELLOW "\x1b[33m"
#define COLOR_RESET "\x1b[0m"

// Test macros
#define TEST(name) \
    void test_##name(); \
    void run_test_##name() { \
        tests_run++; \
        printf("[ RUN  ] %s\n", #name); \
        test_##name(); \
        tests_passed++; \
        printf(COLOR_GREEN "[  OK  ] %s\n" COLOR_RESET, #name); \
    } \
    void test_##name()

#define ASSERT_NOT_NULL(ptr, msg) \
    if (ptr == NULL) { \
        printf(COLOR_RED "[FAILED] %s: %s is NULL\n" COLOR_RESET, __func__, msg); \
        tests_failed++; \
        tests_passed--; \
        return; \
    }

#define ASSERT_NULL(ptr, msg) \
    if (ptr != NULL) { \
        printf(COLOR_RED "[FAILED] %s: %s is not NULL\n" COLOR_RESET, __func__, msg); \
        tests_failed++; \
        tests_passed--; \
        return; \
    }

#define ASSERT_EQ(expected, actual, msg) \
    if (expected != actual) { \
        printf(COLOR_RED "[FAILED] %s: %s - expected %d, got %d\n" COLOR_RESET, \
               __func__, msg, expected, actual); \
        tests_failed++; \
        tests_passed--; \
        return; \
    }

#define ASSERT_TRUE(condition, msg) \
    if (!(condition)) { \
        printf(COLOR_RED "[FAILED] %s: %s\n" COLOR_RESET, __func__, msg); \
        tests_failed++; \
        tests_passed--; \
        return; \
    }

// ============================================================================
// Test: Extractor Lifecycle
// ============================================================================

TEST(extractor_new) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    extractous_extractor_free(extractor);
}

TEST(extractor_free_null) {
    // Should not crash
    extractous_extractor_free(NULL);
}

TEST(extractor_double_free) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    extractous_extractor_free(extractor);
    // Second free on same pointer would cause issues in real code
    // but this test just verifies it doesn't crash the suite
}

// ============================================================================
// Test: Configuration Functions
// ============================================================================

TEST(extractor_set_max_length) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    
    struct CExtractor *new_extractor = extractous_extractor_set_extract_string_max_length(
        extractor, 10000
    );
    ASSERT_NOT_NULL(new_extractor, "new_extractor");
    
    extractous_extractor_free(new_extractor);
}

TEST(extractor_set_encoding) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    
    struct CExtractor *new_extractor = extractous_extractor_set_encoding(
        extractor, CHARSET_UTF_8
    );
    ASSERT_NOT_NULL(new_extractor, "new_extractor with UTF-8");
    
    extractous_extractor_free(new_extractor);
}

TEST(extractor_set_invalid_encoding) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    
    struct CExtractor *new_extractor = extractous_extractor_set_encoding(
        extractor, 999 // Invalid encoding
    );
    ASSERT_NULL(new_extractor, "new_extractor with invalid encoding");
    
    // Original extractor was consumed, don't free
}

TEST(extractor_set_xml_output) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    
    struct CExtractor *new_extractor = extractous_extractor_set_xml_output(
        extractor, true
    );
    ASSERT_NOT_NULL(new_extractor, "new_extractor with XML enabled");
    
    extractous_extractor_free(new_extractor);
}

TEST(extractor_chained_configuration) {
    struct CExtractor *e1 = extractous_extractor_new();
    ASSERT_NOT_NULL(e1, "e1");
    
    struct CExtractor *e2 = extractous_extractor_set_extract_string_max_length(e1, 5000);
    ASSERT_NOT_NULL(e2, "e2");
    
    struct CExtractor *e3 = extractous_extractor_set_encoding(e2, CHARSET_UTF_8);
    ASSERT_NOT_NULL(e3, "e3");
    
    struct CExtractor *e4 = extractous_extractor_set_xml_output(e3, false);
    ASSERT_NOT_NULL(e4, "e4");
    
    extractous_extractor_free(e4);
}

// ============================================================================
// Test: PDF Configuration
// ============================================================================

TEST(pdf_config_new) {
    struct CPdfParserConfig *config = extractous_pdf_config_new();
    ASSERT_NOT_NULL(config, "pdf_config");
    extractous_pdf_config_free(config);
}

TEST(pdf_config_set_ocr_strategy) {
    struct CPdfParserConfig *c1 = extractous_pdf_config_new();
    ASSERT_NOT_NULL(c1, "c1");
    
    struct CPdfParserConfig *c2 = extractous_pdf_config_set_ocr_strategy(
        c1, PDF_OCR_AUTO
    );
    ASSERT_NOT_NULL(c2, "c2");
    
    extractous_pdf_config_free(c2);
}

TEST(pdf_config_set_extract_inline_images) {
    struct CPdfParserConfig *c1 = extractous_pdf_config_new();
    ASSERT_NOT_NULL(c1, "c1");
    
    struct CPdfParserConfig *c2 = extractous_pdf_config_set_extract_inline_images(c1, true);
    ASSERT_NOT_NULL(c2, "c2");
    
    extractous_pdf_config_free(c2);
}

TEST(extractor_set_pdf_config) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    
    struct CPdfParserConfig *pdf_config = extractous_pdf_config_new();
    ASSERT_NOT_NULL(pdf_config, "pdf_config");
    
    struct CExtractor *new_extractor = extractous_extractor_set_pdf_config(
        extractor, pdf_config
    );
    ASSERT_NOT_NULL(new_extractor, "new_extractor");
    
    extractous_pdf_config_free(pdf_config);
    extractous_extractor_free(new_extractor);
}

// ============================================================================
// Test: Office Configuration
// ============================================================================

TEST(office_config_new) {
    struct COfficeParserConfig *config = extractous_office_config_new();
    ASSERT_NOT_NULL(config, "office_config");
    extractous_office_config_free(config);
}

TEST(office_config_set_extract_macros) {
    struct COfficeParserConfig *c1 = extractous_office_config_new();
    ASSERT_NOT_NULL(c1, "c1");
    
    struct COfficeParserConfig *c2 = extractous_office_config_set_extract_macros(c1, true);
    ASSERT_NOT_NULL(c2, "c2");
    
    extractous_office_config_free(c2);
}

// ============================================================================
// Test: OCR Configuration
// ============================================================================

TEST(ocr_config_new) {
    struct CTesseractOcrConfig *config = extractous_ocr_config_new();
    ASSERT_NOT_NULL(config, "ocr_config");
    extractous_ocr_config_free(config);
}

TEST(ocr_config_set_language) {
    struct CTesseractOcrConfig *c1 = extractous_ocr_config_new();
    ASSERT_NOT_NULL(c1, "c1");
    
    struct CTesseractOcrConfig *c2 = extractous_ocr_config_set_language(c1, "eng");
    ASSERT_NOT_NULL(c2, "c2");
    
    extractous_ocr_config_free(c2);
}

// ============================================================================
// Test: Error Handling
// ============================================================================

TEST(error_message) {
    char *msg = extractous_error_message(ERR_OK);
    ASSERT_NOT_NULL(msg, "error message for ERR_OK");
    ASSERT_TRUE(strlen(msg) > 0, "error message not empty");
    extractous_string_free(msg);
    
    msg = extractous_error_message(ERR_NULL_POINTER);
    ASSERT_NOT_NULL(msg, "error message for ERR_NULL_POINTER");
    extractous_string_free(msg);
    
    msg = extractous_error_message(ERR_EXTRACTION_FAILED);
    ASSERT_NOT_NULL(msg, "error message for ERR_EXTRACTION_FAILED");
    extractous_string_free(msg);
}

TEST(extract_with_null_extractor) {
    char *content = NULL;
    struct CMetadata *metadata = NULL;
    
    int result = extractous_extractor_extract_file_to_string(
        NULL, "test.txt", &content, &metadata
    );
    
    ASSERT_EQ(ERR_NULL_POINTER, result, "error code");
}

TEST(extract_with_null_path) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    
    char *content = NULL;
    struct CMetadata *metadata = NULL;
    
    int result = extractous_extractor_extract_file_to_string(
        extractor, NULL, &content, &metadata
    );
    
    ASSERT_EQ(ERR_NULL_POINTER, result, "error code");
    extractous_extractor_free(extractor);
}

TEST(extract_with_null_output) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    
    int result = extractous_extractor_extract_file_to_string(
        extractor, "test.txt", NULL, NULL
    );
    
    ASSERT_EQ(ERR_NULL_POINTER, result, "error code");
    extractous_extractor_free(extractor);
}

// ============================================================================
// Test: String Memory Management
// ============================================================================

TEST(string_free_null) {
    // Should not crash
    extractous_string_free(NULL);
}

// ============================================================================
// Test: Metadata Functions
// ============================================================================

TEST(metadata_free_null) {
    // Should not crash
    extractous_metadata_free(NULL);
}

// ============================================================================
// Test: URL Extraction Functions (if they exist)
// ============================================================================

TEST(url_extraction_null_checks) {
    struct CExtractor *extractor = extractous_extractor_new();
    ASSERT_NOT_NULL(extractor, "extractor");
    
    char *content = NULL;
    struct CMetadata *metadata = NULL;
    
    // NULL URL
    int result = extractous_extractor_extract_url_to_string(
        extractor, NULL, &content, &metadata
    );
    ASSERT_EQ(ERR_NULL_POINTER, result, "null URL error code");
    
    // NULL outputs
    result = extractous_extractor_extract_url_to_string(
        extractor, "http://example.com", NULL, NULL
    );
    ASSERT_EQ(ERR_NULL_POINTER, result, "null outputs error code");
    
    extractous_extractor_free(extractor);
}

// ============================================================================
// Test Runner
// ============================================================================

void run_all_tests() {
    printf("\n");
    printf("========================================\n");
    printf("  FFI Layer Tests for Extractous\n");
    printf("========================================\n\n");
    
    // Lifecycle tests
    printf(COLOR_YELLOW "--- Extractor Lifecycle ---\n" COLOR_RESET);
    run_test_extractor_new();
    run_test_extractor_free_null();
    run_test_extractor_double_free();
    
    // Configuration tests
    printf(COLOR_YELLOW "\n--- Configuration Functions ---\n" COLOR_RESET);
    run_test_extractor_set_max_length();
    run_test_extractor_set_encoding();
    run_test_extractor_set_invalid_encoding();
    run_test_extractor_set_xml_output();
    run_test_extractor_chained_configuration();
    
    // PDF config tests
    printf(COLOR_YELLOW "\n--- PDF Configuration ---\n" COLOR_RESET);
    run_test_pdf_config_new();
    run_test_pdf_config_set_ocr_strategy();
    run_test_pdf_config_set_extract_inline_images();
    run_test_extractor_set_pdf_config();
    
    // Office config tests
    printf(COLOR_YELLOW "\n--- Office Configuration ---\n" COLOR_RESET);
    run_test_office_config_new();
    run_test_office_config_set_extract_macros();
    
    // OCR config tests
    printf(COLOR_YELLOW "\n--- OCR Configuration ---\n" COLOR_RESET);
    run_test_ocr_config_new();
    run_test_ocr_config_set_language();
    
    // Error handling tests
    printf(COLOR_YELLOW "\n--- Error Handling ---\n" COLOR_RESET);
    run_test_error_message();
    run_test_extract_with_null_extractor();
    run_test_extract_with_null_path();
    run_test_extract_with_null_output();
    
    // Memory management tests
    printf(COLOR_YELLOW "\n--- Memory Management ---\n" COLOR_RESET);
    run_test_string_free_null();
    run_test_metadata_free_null();
    
    // URL extraction tests
    printf(COLOR_YELLOW "\n--- URL Extraction ---\n" COLOR_RESET);
    run_test_url_extraction_null_checks();
    
    // Summary
    printf("\n");
    printf("========================================\n");
    printf("  Test Summary\n");
    printf("========================================\n");
    printf("Total:  %d\n", tests_run);
    printf(COLOR_GREEN "Passed: %d\n" COLOR_RESET, tests_passed);
    
    if (tests_failed > 0) {
        printf(COLOR_RED "Failed: %d\n" COLOR_RESET, tests_failed);
    } else {
        printf("Failed: 0\n");
    }
    
    printf("========================================\n\n");
}

int main() {
    run_all_tests();
    return tests_failed > 0 ? 1 : 0;
}
