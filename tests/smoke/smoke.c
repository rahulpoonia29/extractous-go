#include "../../include/extractous.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <dirent.h>
#include <sys/stat.h>

#define TEST_DIR "test_files"

#define REQ(c, m) do { if (!(c)) { fprintf(stderr, "FAIL: %s\n", m); return 1; } } while (0)

static void free_cstr(char* s) { if (s) extractous_string_free(s); }
static void free_meta(struct CMetadata* m) { if (m) extractous_metadata_free(m); }

static void print_metadata(const struct CMetadata* meta) {
    if (!meta) return;
    for (size_t i = 0; i < meta->len; i++) {
        printf("[meta] %s: %s\n", meta->keys[i], meta->values[i]);
    }
}

static int read_file_bytes(const char* path, uint8_t** out_data, size_t* out_len) {
    FILE* f = fopen(path, "rb");
    if (!f) return 0;
    fseek(f, 0, SEEK_END);
    long len = ftell(f);
    rewind(f);

    uint8_t* buf = malloc(len);
    if (!buf) {
        fclose(f);
        return 0;
    }

    fread(buf, 1, len, f);
    fclose(f);
    *out_data = buf;
    *out_len = len;
    return 1;
}

int main() {
    printf("[smoke] begin\n");

    struct CPdfParserConfig* pcfg = extractous_pdf_config_new();
    pcfg = extractous_pdf_config_set_ocr_strategy(pcfg, PDF_OCR_AUTO);
    pcfg = extractous_pdf_config_set_extract_inline_images(pcfg, true);
    pcfg = extractous_pdf_config_set_extract_unique_inline_images_only(pcfg, true);
    pcfg = extractous_pdf_config_set_extract_marked_content(pcfg, false);
    pcfg = extractous_pdf_config_set_extract_annotation_text(pcfg, true);

    struct COfficeParserConfig* ocfg = extractous_office_config_new();
    ocfg = extractous_office_config_set_extract_macros(ocfg, true);
    ocfg = extractous_office_config_set_include_deleted_content(ocfg, true);
    ocfg = extractous_office_config_set_include_move_from_content(ocfg, false);
    ocfg = extractous_office_config_set_include_shape_based_content(ocfg, true);

    struct CTesseractOcrConfig* ocr = extractous_ocr_config_new();
    ocr = extractous_ocr_config_set_language(ocr, "eng");
    ocr = extractous_ocr_config_set_density(ocr, 300);
    ocr = extractous_ocr_config_set_depth(ocr, 8);
    ocr = extractous_ocr_config_set_enable_image_preprocessing(ocr, true);
    ocr = extractous_ocr_config_set_timeout_seconds(ocr, 30);

    struct CExtractor* ex = extractous_extractor_new();
    ex = extractous_extractor_set_extract_string_max_length(ex, 4096);
    ex = extractous_extractor_set_encoding(ex, CHARSET_UTF_8);
    ex = extractous_extractor_set_pdf_config(ex, pcfg);
    ex = extractous_extractor_set_office_config(ex, ocfg);
    ex = extractous_extractor_set_ocr_config(ex, ocr);

    DIR* dir = opendir(TEST_DIR);
    if (!dir) {
        fprintf(stderr, "[error] Could not open test directory: %s\n", TEST_DIR);
        return 1;
    }

    struct dirent* entry;
    while ((entry = readdir(dir)) != NULL) {
        if (entry->d_type != DT_REG) continue; // skip non-regular files

        char path[512];
        snprintf(path, sizeof(path), "%s/%s", TEST_DIR, entry->d_name);
        printf("\n[==] Testing file: %s\n", path);

        // --- Extract to string ---
        char* out = NULL;
        struct CMetadata* meta = NULL;
        int rc = extractous_extractor_extract_file_to_string(ex, path, &out, &meta);
        if (rc == ERR_OK && out) {
            printf("[extract_file_to_string] Output:\n%s\n", out);
            print_metadata(meta);
        } else {
            char* em = extractous_error_message(rc);
            fprintf(stderr, "[error] extract_file_to_string rc=%d msg=%s\n", rc, em ? em : "(null)");
            free_cstr(em);
        }
        free_cstr(out);
        free_meta(meta);

        // --- Extract from bytes ---
        uint8_t* data = NULL;
        size_t data_len = 0;
        if (read_file_bytes(path, &data, &data_len)) {
            char* out2 = NULL;
            struct CMetadata* meta2 = NULL;
            rc = extractous_extractor_extract_bytes_to_string(ex, data, data_len, &out2, &meta2);
            if (rc == ERR_OK && out2) {
                printf("[extract_bytes_to_string] Output:\n%s\n", out2);
                print_metadata(meta2);
            } else {
                char* em = extractous_error_message(rc);
                fprintf(stderr, "[error] extract_bytes_to_string rc=%d msg=%s\n", rc, em ? em : "(null)");
                free_cstr(em);
            }
            free(data);
            free_cstr(out2);
            free_meta(meta2);
        } else {
            fprintf(stderr, "[warn] Failed to read file bytes: %s\n", path);
        }
    }

    closedir(dir);

    extractous_extractor_free(ex);
    printf("\n[smoke] success\n");
    return 0;
}