/**
 * Streaming extraction example
 * 
 * Demonstrates how to:
 * - Extract large files using streaming
 * - Process content in chunks
 * - Avoid loading entire file into memory
 */

#include <extractous.h>
#include <stdio.h>
#include <stdlib.h>

#define BUFFER_SIZE 4096

int main(int argc, char** argv) {
    if (argc != 2) {
        fprintf(stderr, "Usage: %s <file_path>\n", argv[0]);
        return 1;
    }

    const char* file_path = argv[1];
    
    // Create extractor
    CExtractor* extractor = extractous_extractor_new();
    if (!extractor) {
        fprintf(stderr, "Failed to create extractor\n");
        return 1;
    }
    
    // Extract to stream
    CStreamReader* reader = NULL;
    CMetadata* metadata = NULL;
    
    printf("Streaming extraction from: %s\n", file_path);
    
    int err = extractous_extractor_extract_file(
        extractor,
        file_path,
        &reader,
        &metadata
    );
    
    if (err != ERR_OK) {
        char* error_msg = extractous_error_message(err);
        fprintf(stderr, "Failed to start extraction (code %d): %s\n", err, error_msg);
        extractous_string_free(error_msg);
        extractous_extractor_free(extractor);
        return 1;
    }
    
    // Print metadata first
    printf("\n=== Metadata (%zu entries) ===\n", metadata->len);
    for (size_t i = 0; i < metadata->len; i++) {
        printf("%s: %s\n", metadata->keys[i], metadata->values[i]);
    }
    
    // Stream content in chunks
    printf("\n=== Content ===\n");
    
    char buffer[BUFFER_SIZE];
    size_t bytes_read;
    size_t total_bytes = 0;
    
    while (extractous_stream_read(reader, (uint8_t*)buffer, BUFFER_SIZE, &bytes_read) == ERR_OK 
           && bytes_read > 0) {
        // Process chunk (here we just print it)
        fwrite(buffer, 1, bytes_read, stdout);
        total_bytes += bytes_read;
    }
    
    printf("\n\n=== Summary ===\n");
    printf("Total bytes read: %zu\n", total_bytes);
    
    // Cleanup
    extractous_stream_free(reader);
    extractous_metadata_free(metadata);
    extractous_extractor_free(extractor);
    
    printf("Streaming extraction successful!\n");
    return 0;
}
