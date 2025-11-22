# TODO - Benchmarks Project

## Completed ✅
- HTTP server benchmark orchestrator with proper output capture
- Memory tracking using `ps` RSS measurement
- Startup time benchmarks using hyperfine
- Proper process group management for clean termination
- Command structure: `run server [api-name]` and `run startup`
- Build system integration with uSockets/uWebSockets Makefiles
- Port standardization (all servers on 8080)

## In Progress 🚧

### 1. Serialization Benchmarks
**Location:** `serialization/`
**Status:** Partial implementation exists but not integrated into benchrunner

**TODO:**
- [ ] Implement JSON serialization benchmarks (serde_json, encoding/json, JSON.parse)
- [ ] Implement Protocol Buffers benchmarks
- [ ] Implement MessagePack benchmarks
- [ ] Implement Avro benchmarks (schema exists in `proto/benchmark.avsc`)
- [ ] Add `benchrunner run serialization` command
- [ ] Create consistent test data across all implementations
- [ ] Measure serialization time, deserialization time, and payload size

**Files:**
- `serialization/Cargo.toml` - Rust dependencies configured
- `serialization/benches/serialization_bench.rs` - Benchmark harness
- `serialization/proto/benchmark.avsc` - Avro schema definition
- `serialization/src/lib.rs` - Library structure

### 2. Additional Server Implementations
**Status:** Configuration exists but servers not implemented

**TODO:**
- [ ] Add Rust Axum HTTP server
- [ ] Add Python Flask/Starlette alternatives
- [ ] Add C++ beast (Boost.Beast) implementation
- [ ] Add Bun.js HTTP server
- [ ] Add Deno HTTP server

### 3. Enhanced HTTP Benchmarks
**TODO:**
- [ ] Add POST/PUT/DELETE benchmarks (currently only GET)
- [ ] Add JSON payload benchmarks
- [ ] Add file upload benchmarks
- [ ] Add WebSocket connection benchmarks
- [ ] Add concurrent connection scaling tests (10, 100, 1000, 10000)
- [ ] Add keepalive vs connection-per-request comparison

### 4. Results Management
**Status:** Basic JSON output implemented

**TODO:**
- [ ] Create results comparison tool
- [ ] Generate HTML/Markdown reports with charts
- [ ] Track benchmark history over time
- [ ] Add statistical analysis (mean, median, p95, p99)
- [ ] Export to CSV for external analysis

### 5. CI/CD Integration
**TODO:**
- [ ] GitHub Actions workflow for automated benchmarks
- [ ] Benchmark regression detection
- [ ] Automated results publishing to GitHub Pages
- [ ] Performance comparison on PRs

### 6. Documentation
**TODO:**
- [ ] Add architecture documentation
- [ ] Document how to add new servers
- [ ] Document how to add new benchmark types
- [ ] Add performance tuning guide
- [ ] Document interpretation of results

### 7. Infrastructure
**TODO:**
- [ ] Docker containers for consistent benchmark environment
- [ ] Support for remote benchmarking (client/server on different machines)
- [ ] CPU affinity/isolation for more consistent results
- [ ] System resource monitoring during benchmarks (CPU, memory, network)

## Nice to Have 💡

### Database Benchmarks
- [ ] PostgreSQL connection pool benchmarks
- [ ] Redis get/set benchmarks
- [ ] MongoDB CRUD benchmarks

### Compression Benchmarks
- [ ] gzip/brotli/zstd comparison
- [ ] Different compression levels

### Cold Start Benchmarks
- [ ] Container startup time
- [ ] Lambda/serverless cold start simulation

### Memory Leak Detection
- [ ] Long-running stability tests
- [ ] Memory growth analysis over time

## Known Issues 🐛
- None currently tracked

## Notes 📝
- All HTTP servers standardized on port 8080
- Uses `stdbuf -oL -eL` to unbuffer http_load_test output
- Process groups are killed to ensure no zombie processes
- 5-second delay between benchmarks to ensure port cleanup
- Hyperfine handles warmup automatically (3 runs before measurement)

---
**Last Updated:** 2025-11-22
