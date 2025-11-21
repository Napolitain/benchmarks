# Benchmarks Repository

This repository contains various benchmarking tools and tests.

## HTTP 1.1 API Benchmarks

See [api/README.md](api/README.md) for HTTP server benchmarks.

### Quick Start

```bash
# Install and build
go mod download
go build -o bin/benchrunner.exe ./cmd/benchrunner

# Run benchmarks
./bin/benchrunner run
```

## Benchmark Orchestrator CLI

The `benchrunner` tool automates HTTP server benchmarking:

**List available servers:**
```bash
benchrunner list
```

**Run all benchmarks:**
```bash
benchrunner run
```

**Run with custom parameters:**
```bash
benchrunner run --connections 200 --duration 30 --pipeline 10
```

**Run specific servers:**
```bash
benchrunner run --servers go-http,node-http
```

Results are saved to `results/` directory.

## Requirements

- **Go 1.21+**
- **Python 3.8+** with pip
- **Node.js 14+**
- **Nginx**
- **GCC/G++** (for compiling C/C++ binaries)
- **Make**
- **Linux** (recommended for benchmarking)
