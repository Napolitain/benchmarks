# Deployment Guide for Linux Benchmarking

This guide covers deploying and running benchmarks on a Linux system.

## Prerequisites

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y build-essential gcc g++ make git golang python3 python3-pip nodejs npm nginx

# Check versions
go version        # Should be 1.21+
python3 --version # Should be 3.8+
node --version    # Should be 14+
gcc --version
g++ --version
```

## Setup

1. **Clone and Build:**

```bash
cd /path/to/benchmarks

# Install Go dependencies
go mod download

# Build the benchrunner CLI
go build -o bin/benchrunner ./cmd/benchrunner

# Build all required binaries (this will build uSockets and uWebSockets)
./bin/benchrunner build
```

2. **Install Python Dependencies:**

```bash
cd api/python-fastapi
pip3 install -r requirements.txt
cd ../..
```

## Running Benchmarks

### Quick Start

```bash
# List available servers
./bin/benchrunner list

# Run all benchmarks with default settings (100 connections, 10s duration)
./bin/benchrunner run

# Results will be saved to results/benchmark_YYYYMMDD_HHMMSS.json
```

### Custom Benchmark Parameters

```bash
# High load test
./bin/benchrunner run --connections 500 --duration 30 --pipeline 10

# Test specific servers
./bin/benchrunner run --servers go-http,cpp-uwebsockets --duration 60

# Low load baseline
./bin/benchrunner run --connections 10 --duration 5
```

### Benchmark Parameters

- `--connections` / `-c`: Number of concurrent connections (default: 100)
- `--pipeline` / `-p`: HTTP pipelining factor - requests per connection (default: 1)
- `--duration` / `-d`: Duration in seconds (default: 10)
- `--servers` / `-s`: Comma-separated list of servers to test (default: all)

## Available Servers

| Server | Language | Framework |
|--------|----------|-----------|
| go-http | Go | net/http (stdlib) |
| go-fasthttp | Go | fasthttp |
| python-fastapi | Python | FastAPI + Uvicorn |
| node-http | Node.js | http (stdlib) |
| nginx-static | - | Nginx (static file) |
| cpp-uwebsockets | C++ | uWebSockets |

## Manual Testing

Run servers individually for manual testing:

```bash
# Go HTTP
cd api/go-http && go run main.go

# Go FastHTTP
cd api/go-fasthttp && go run main.go

# Python FastAPI
cd api/python-fastapi && python3 main.py

# Node.js
cd api/node-http && node index.js

# Nginx
cd api/nginx-static && nginx -p . -c nginx.conf
# Stop: nginx -p . -c nginx.conf -s stop

# C++ uWebSockets
./bin/HelloWorldBenchmark
```

Test endpoint:
```bash
curl http://localhost:8080/
# Expected: {"message":"Hello, World!"}
```

## Build System Details

The benchrunner CLI orchestrates building:

1. **uSockets Library**: Static library (`uSockets.a`) built first with `WITH_OPENSSL=0`
2. **http_load_test**: Load testing tool compiled from uSockets examples
3. **HelloWorldBenchmark**: uWebSockets C++ server compiled using their build system

All binaries are placed in `bin/` directory.

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
sudo lsof -i :8080
# Kill it
sudo kill -9 <PID>
```

### Build Failures

```bash
# Clean build artifacts
cd api/uWebSockets/uSockets
make clean
cd ../
rm -f *.o build HelloWorldBenchmark

# Rebuild
cd /path/to/benchmarks
rm -rf bin/
./bin/benchrunner build
```

### Python Module Not Found

```bash
# Use virtual environment
cd api/python-fastapi
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

## Results

Benchmark results are saved as JSON in `results/` directory with timestamp:
- `req_per_sec`: Average requests per second
- `connections`: Number of connections used
- `pipeline`: Pipeline factor
- `duration_seconds`: Test duration
- `timestamp`: When test was run
- `error`: Any error message (if failed)

Results are also printed as a summary table to stdout.
