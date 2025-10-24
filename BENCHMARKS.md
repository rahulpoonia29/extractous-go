
## BENCHMARKS.md

```markdown
# Performance Benchmarks

Benchmark results from [extractous-benchmarks](https://github.com/yobix-ai/extractous-benchmarks) dataset.

## Test Environment

- Dataset: Research papers and financial reports from [extractous-benchmarks/dataset](https://github.com/yobix-ai/extractous-benchmarks/tree/main/dataset)
- Libraries tested: extractous-go (string mode, streaming mode), ledongthuc/pdf
- Files: 40 PDFs ranging from 0.25 MB to 14.7 MB

## Summary Statistics

### extractous-string (In-Memory Mode)

- **Average Throughput**: 4.67 MB/s
- **Average Memory**: 11.27 MB
- **Average Accuracy**: 96.76%
- **Average F1 Score**: 98.63%

### extractous-stream-64KB (Streaming Mode)

- **Average Throughput**: 2.28 MB/s
- **Average Memory**: 21.28 MB
- **Average Accuracy**: 93.82%
- **Average F1 Score**: 96.14%

### ledongthuc-pdf

- **Average Throughput**: 23.68 MB/s
- **Average Memory**: 12.25 MB
- **Average Accuracy**: 83.72%
- **Average F1 Score**: 94.58%

## Observations

1. **Accuracy vs Speed**: extractous-go prioritizes accuracy with comprehensive text extraction (96%+ accuracy) while ledongthuc-pdf is faster but less accurate (83%).

2. **Memory Usage**: String mode uses less memory for small files, while streaming mode maintains consistent memory usage for large files.

3. **Best Use Cases**:
   - **extractous-string**: Small to medium files (<5 MB) where accuracy matters
   - **extractous-stream**: Large files (>10 MB) to avoid loading entire content into memory
   - **ledongthuc-pdf**: Speed-critical applications where text quality is less important

## Sample Results

### Large File Performance (14.7 MB PDF)

| Library | Duration (ms) | Throughput (MB/s) | Memory (MB) | Accuracy (%) |
|---------|--------------|------------------|-------------|--------------|
| extractous-string | 400 | 36.70 | 15.78 | 86.95 |
| extractous-stream | 1037 | 14.16 | 21.83 | 87.74 |
| ledongthuc-pdf | 185 | 79.38 | 44.67 | 82.02 |

### Small File Performance (0.78 MB PDF)

| Library | Duration (ms) | Throughput (MB/s) | Memory (MB) | Accuracy (%) |
|---------|--------------|------------------|-------------|--------------|
| extractous-string | 729 | 1.07 | 6.34 | 96.61 |
| extractous-stream | 729 | 1.07 | 7.37 | 96.90 |
| ledongthuc-pdf | 1 | 709.49 | 1.46 | 0.54 |

Note: ledongthuc-pdf fails silently on some complex PDFs, resulting in near-zero accuracy.

## Full Results

Complete benchmark data is available in `benchmark_results_20251018_034604.csv`.

## Running Benchmarks

```
# Clone benchmarks repository
git clone https://github.com/yobix-ai/extractous-benchmarks
cd extractous-benchmarks

# Run benchmarks
go test -bench=. -benchmem
```

## Methodology

Benchmarks measure:
- **Duration**: Time to extract text from file
- **Throughput**: File size divided by duration
- **Memory**: Peak memory usage during extraction
- **Accuracy**: Edit distance compared to ground truth
- **Precision/Recall/F1**: Standard information retrieval metrics

Each file is tested multiple times and results are averaged.
```

These files provide a professional, developer-focused documentation structure:[1][2][3]

[1](https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/84579331/cca3bd07-1380-4c4b-a89f-a7be0003a5f8/DEVELOPER.md-Complete-Developer-Guide.md)
[2](https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/84579331/5378c626-0f89-4d2a-a61a-9c7ecc8e1e2d/benchmark_results_20251018_034604.csv)
[3](https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/84579331/72e5f134-4845-4ced-8882-4f1a11ee5e07/rahulpoonia29-extractous-go-8a5edab282632443-3.txt)