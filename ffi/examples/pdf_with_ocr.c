/**
 * PDF with OCR example
 * 
 * Demonstrates how to:
 * - Configure PDF parser for OCR
 * - Set OCR language and parameters
 * - Extract from scanned PDFs
 */

#include <extractous.h>
#include <stdio.h>
#include <stdlib.h>

int main() {
    // PDF config
    CPdfParserConfig* pdf = extractous_pdf_config_new();
    if (!pdf) {
        fprintf(stderr, "Failed to create PDF config\n");
        return 1;
    }
    pdf = extractous_pdf_config_set_ocr_strategy(pdf, PDF_OCR_STRATEGY_AUTO);
    pdf = extractous_pdf_config_set_extract_annotation_text(pdf, true);
    
    // OCR config
    CTesseractOcrConfig* ocr = extractous_ocr_config_new();
    if (!ocr) {
        fprintf(stderr, "Failed to create OCR config\n");
        extractous_pdf_config_free(pdf);
        return 1;
    }
    ocr = extractous_ocr_config_set_language(ocr, "eng");
    ocr = extractous_ocr_config_set_density(ocr, 300);
    
    // Extractor
    CExtractor* ext = extractous_extractor_new();
    if (!ext) {
        fprintf(stderr, "Failed to create extractor\n");
        extractous_ocr_config_free(ocr);
        extractous_pdf_config_free(pdf);
        return 1;
    }
    ext = extractous_extractor_set_pdf_config(ext, pdf);  // Consumes pdf
    ext = extractous_extractor_set_ocr_config(ext, ocr);  // Consumes ocr
    
    // Extract
    char* content;
    CMetadata* metadata;
    int result = extractous_extractor_extract_file_to_string(
        ext, "document.pdf", &content, &metadata
    );
    
    if (result == ERR_OK) {
        printf("Content: %s\n", content);
        extractous_string_free(content);
        extractous_metadata_free(metadata);
    } else {
        char* msg = extractous_error_message(result);
        fprintf(stderr, "Error: %s\n", msg);
        extractous_string_free(msg);
    }
    
    extractous_extractor_free(ext);
    return 0;
}
