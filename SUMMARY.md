# Benchmark Orchestrator - Implementation Summary

## Overview

Complete HTTP 1.1 benchmark suite with automated orchestration for comparing performance across multiple languages and frameworks.

## What Was Built

### 1. HTTP Server Implementations (6 total)

All servers serve the same response: `{"message":"Hello, World!"}`

| Server | Language | Framework | Notes |
|--------|----------|-----------|-------|
| **go-http** | Go | net/http | Standard library, `GOMAXPROCS=1` |
| **go-fasthttp** | Go | valyala/fasthttp | High-performance library, `GOMAXPROCS=1` |
| **python-fastapi** | Python | FastAPI + Uvicorn | ASGI with async, `workers=1` |
| **node-http** | Node.js | http | Native module, single-threaded by design |
| **nginx-static** | - | Nginx | Static file serving, `worker_processes=1` |
| **cpp-uwebsockets** | C++ | uWebSockets | Low-level C++ library, single-threaded |

### 2. Load Testing Tool

- **http_load_test**: Compiled from uSockets/examples
- Features:
  - Configurable connections
  - HTTP pipelining support
  - Reports req/sec metrics
  - Built without SSL for simplicity

### 3. Orchestration CLI (`benchrunner`)

Go-based CLI tool with Cobra that automates:
- Building all required binaries
- Starting/stopping servers sequentially
- Running load tests
- Collecting and reporting results

**Commands:**
```bash
benchrunner list                    # List available servers
benchrunner build                   # Build all binaries
benchrunner run                     # Run all benchmarks
benchrunner run --servers go-http   # Run specific server
benchrunner run -c 500 -d 30        # Custom parameters
```

## Architecture

```
benchmarks/
├── cmd/benchrunner/          # CLI application entry point
├── internal/
│   ├── builder/             # Builds uSockets, http_load_test, uWebSockets
│   ├── server/              # Server lifecycle (start/stop/wait)
│   ├── benchmark/           # Benchmark execution and result parsing
│   └── config/              # Server configurations
├── api/                     # All server implementations
│   ├── go-http/
│   ├── go-fasthttp/
│   ├── python-fastapi/
│   ├── node-http/
│   ├── nginx-static/
│   └── uWebSockets/        # Submodule with uSockets
├── bin/                     # Compiled binaries (gitignored)
└── results/                 # JSON results with timestamps
```

## Build System Integration

### uSockets (C library)
- Built via Makefile in `uSockets/`
- Produces `uSockets.a` static library
- Compiled with `WITH_OPENSSL=0` (no SSL)

### http_load_test (C)
- Source: `uSockets/examples/http_load_test.c`
- Compiled with gcc against `uSockets.a`
- No SSL support (`-DLIBUS_NO_SSL`)

### uWebSockets Server (C++)
- Source: `uWebSockets/examples/HelloWorldBenchmark.cpp` (newly added)
- Built using uWebSockets' `build.c` system
- Compiled with: `WITH_OPENSSL=0`, `WITH_ZLIB=0`, `WITH_LTO=0`
- Integrated into example build list in `build.c`

## Key Design Decisions

1. **Single-threaded**: All servers forced to single thread for fair comparison
   - Go: `runtime.GOMAXPROCS(1)`
   - Python: `workers=1`
   - Node: Single-threaded by default
   - Nginx: `worker_processes=1`
   - uWebSockets: Uses `App()` not `LocalCluster`

2. **No SSL**: Simplified builds and testing, focus on pure HTTP performance

3. **Same Port (8080)**: Servers run sequentially, not concurrently

4. **Linux Target**: Optimized for Linux deployment (production benchmark environment)

5. **Build System Reuse**: 
   - Uses uWebSockets' existing build.c tool
   - Uses uSockets' Makefile
   - No custom build scripts for C/C++ code

## Workflow

```
benchrunner run
    ├─> Build Phase
    │   ├─> make (uSockets.a)
    │   ├─> gcc (http_load_test)
    │   └─> make examples (HelloWorldBenchmark)
    │
    └─> For each server:
        ├─> Start server process
        ├─> Wait for port 8080 ready (TCP probe)
        ├─> Run http_load_test <connections> localhost 8080 <pipeline>
        ├─> Parse "Req/sec: X" output lines
        ├─> Calculate average req/sec
        ├─> Stop server (SIGKILL)
        └─> Wait 2s cooldown
    
    ├─> Print summary table
    └─> Save JSON results
```

## Result Format

```json
[
  {
    "server_name": "cpp-uwebsockets",
    "req_per_sec": 125432.50,
    "connections": 100,
    "pipeline": 1,
    "duration_seconds": 10,
    "timestamp": "2025-11-21T07:00:00Z",
    "error": ""
  }
]
```

## Files Modified in uWebSockets

1. **`examples/HelloWorldBenchmark.cpp`** (new file)
   - Simple HTTP server on port 8080
   - Returns JSON response
   - Single-threaded using `uWS::App()`

2. **`build.c`** (modified)
   - Added "HelloWorldBenchmark" to EXAMPLE_FILES array
   - Now builds when `make examples` is called

## Testing

```bash
# Development (local)
go run ./cmd/benchrunner list
go run ./cmd/benchrunner build
go run ./cmd/benchrunner run --servers go-http --duration 5

# Production (Linux)
./bin/benchrunner run --connections 500 --duration 30
```

## Future Enhancements

Potential additions:
- [ ] Rust (Actix, Axum)
- [ ] Bun HTTP server
- [ ] Deno HTTP server
- [ ] Java (Netty)
- [ ] Add latency percentiles (p50, p95, p99)
- [ ] Multi-threaded comparison mode
- [ ] SSL/TLS benchmarks
- [ ] Chart generation from results
- [ ] CSV export option
