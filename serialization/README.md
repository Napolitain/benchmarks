# Serialization Benchmarks

This Rust CLI tool benchmarks various serialization formats to compare their performance and encoding efficiency.

## Formats Tested

- **Apache Arrow**: A columnar memory format optimized for analytics
- **Protocol Buffers v3**: Google's language-neutral data serialization format
- **Cap'n Proto**: A fast data interchange format with zero-copy deserialization

## Metrics

The benchmarks measure:
1. **Encoding Rate**: Throughput in MB/s or GB/s
2. **Encoding Size**: Size of the serialized data in bytes

## Usage

### Quick Test

Run the CLI tool for a quick benchmark comparison:

```bash
cargo run --release
```

### Detailed Benchmarks

Run comprehensive benchmarks using Criterion:

```bash
cargo bench
```

This will:
- Test with different data sizes (100, 1000, 10000 records)
- Measure encoding throughput
- Compare encoding sizes
- Generate detailed reports in `target/criterion/`

## Test Data Structure

Each benchmark uses a consistent data structure:
- `id`: 64-bit integer
- `name`: String (variable length)
- `value`: 64-bit floating point
- `numbers`: List of 32-bit integers
- `data`: Binary blob (64 bytes)

## Building

```bash
cargo build --release
```

## Development

The project structure:
- `src/main.rs`: CLI tool for quick benchmarks
- `benches/serialization_bench.rs`: Criterion benchmarks
- `proto/`: Protocol definitions for Protobuf and Cap'n Proto
- `build.rs`: Build script for compiling schemas
