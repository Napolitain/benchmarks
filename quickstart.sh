#!/bin/bash
set -e

echo "=================================="
echo "HTTP Benchmark Suite - Quick Start"
echo "=================================="
echo ""

# Check prerequisites
echo "Checking prerequisites..."
command -v go >/dev/null 2>&1 || { echo "Error: Go is not installed"; exit 1; }
command -v gcc >/dev/null 2>&1 || { echo "Error: GCC is not installed"; exit 1; }
command -v g++ >/dev/null 2>&1 || { echo "Error: G++ is not installed"; exit 1; }
command -v make >/dev/null 2>&1 || { echo "Error: Make is not installed"; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "Error: Python3 is not installed"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "Error: Node.js is not installed"; exit 1; }
command -v nginx >/dev/null 2>&1 || { echo "Error: Nginx is not installed"; exit 1; }

echo "✓ All prerequisites found"
echo ""

# Install Python dependencies
echo "Installing Python dependencies..."
cd api/python-fastapi
pip3 install -q -r requirements.txt
cd ../..
echo "✓ Python dependencies installed"
echo ""

# Download Go modules
echo "Downloading Go modules..."
go mod download
echo "✓ Go modules downloaded"
echo ""

# Build benchrunner
echo "Building benchrunner CLI..."
go build -o bin/benchrunner ./cmd/benchrunner
chmod +x bin/benchrunner
echo "✓ Benchrunner built"
echo ""

# Build all binaries
echo "Building load test tool and servers..."
./bin/benchrunner build
echo "✓ All binaries built"
echo ""

echo "=================================="
echo "Setup complete!"
echo "=================================="
echo ""
echo "Available commands:"
echo "  ./bin/benchrunner list              # List available servers"
echo "  ./bin/benchrunner run               # Run all benchmarks (default: 100 conns, 10s)"
echo "  ./bin/benchrunner run -c 500 -d 30  # Custom benchmark"
echo ""
echo "Example:"
echo "  ./bin/benchrunner run --connections 200 --duration 20"
echo ""
