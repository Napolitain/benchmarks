# Benchmarks Repository

This repository contains various benchmarking tools and tests for comparing language performance across different scenarios.

## Quick Start

```bash
# Install and build
go mod download
go build -o benchrunner ./cmd/benchrunner

# List available servers
./benchrunner list

# Run benchmarks
./benchrunner run server
./benchrunner run helloworld
./benchrunner run compute
./benchrunner run cli
./benchrunner run ffi
```

## Benchmark Orchestrator CLI

The `benchrunner` tool orchestrates benchmarks across multiple server and language implementations.

### Commands

| Command | Description |
|---------|-------------|
| `benchrunner build` | Build the http_load_test binary from uSockets |
| `benchrunner list` | List all available server implementations |
| `benchrunner run <type>` | Run different types of benchmarks |
| `benchrunner completion` | Generate autocompletion scripts for your shell |
| `benchrunner help [command]` | Help about any command |

### Benchmark Types

#### HTTP Server Benchmarks

```bash
benchrunner run server [api-name]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-c, --connections` | Number of connections | 100 |
| `-d, --duration` | Duration in seconds | 10 |
| `-p, --pipeline` | Pipeline factor | 1 |

**Examples:**
```bash
benchrunner run server                           # Run all servers
benchrunner run server go-http                   # Run specific server
benchrunner run server -c 200 -d 30 -p 10        # Custom parameters
```

#### CLI Benchmarks (Rectangle YAML Parsing)

```bash
benchrunner run cli [language]
```

#### Compute Benchmarks (Bubblesort)

```bash
benchrunner run compute [language]
```

#### FFI Benchmarks

```bash
benchrunner run ffi [language]
```

Runs two sub-benchmarks:
- **fast_sum**: Measures FFI call overhead (1M calls of a simple sum function)
- **slow_compute**: Measures compute-heavy FFI (100 calls with 1M iterations each)

#### Helloworld Benchmarks

```bash
benchrunner run helloworld [language]
```

### Benchmark Modes

The `cli`, `compute`, `ffi`, and `helloworld` benchmarks support the following modes and flags:

| Flag | Description | Default |
|------|-------------|---------|
| `-m, --mode` | Benchmark mode (see below) | exec |
| `-r, --runs` | Number of benchmark runs | 10 |
| `-w, --warmup` | Number of warmup runs | 3 |

| Mode | Description |
|------|-------------|
| `compile` | Benchmark compilation time only (cold builds) |
| `full-cold` | Benchmark compilation + execution (cold builds, no cache) |
| `full-hot` | Benchmark compilation + execution (hot builds, cache allowed) |
| `exec` | Benchmark execution time only (pre-compiled) |

**Examples:**
```bash
benchrunner run helloworld                       # Run all languages (exec mode)
benchrunner run helloworld go                    # Run specific language
benchrunner run compute -m compile               # Benchmark compilation only
benchrunner run cli -m full-cold -r 20 -w 5      # Full cold benchmark with custom runs
```

Results are saved to the `results/` directory.

## Requirements

- **Go 1.21+**
- **Python 3.8+** with pip
- **Node.js 14+**
- **Nginx**
- **GCC/G++** (for compiling C/C++ binaries)
- **Make**
- **poop** or **hyperfine** (for non-server benchmarks)
- **Linux** (recommended for benchmarking)
