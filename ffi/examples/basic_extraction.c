/**
 * Basic extraction example
 * 
 * Demonstrates how to:
 * - Create an extractor
 * - Extract text from a file
 * - Access metadata
 * - Proper memory cleanup
 */

#include <extractous.h> // Replace with the correct path to extractous.h
#include <stdio.h>
#include <stdlib.h>

int main() {
    CExtractor* ext = extractous_extractor_new();
    if (!ext) {
        fprintf(stderr, "Failed to create extractor\n");
        return 1;
    }
    
    char* content;
    CMetadata* metadata;
    
    int result = extractous_extractor_extract_file_to_string(
        ext, "document.pdf", &content, &metadata
    );
    
    if (result == ERR_OK) {
        printf("Content: %s\n", content);
        
        for (size_t i = 0; i < metadata->len; i++) {
            printf("%s: %s\n", metadata->keys[i], metadata->values[i]);
        }
        
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
