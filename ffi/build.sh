#!/bin/bash
set -e

echo "=== Building FFI Benchmark ==="

# Build shared library
g++ -O3 -fPIC -shared -o libhotpath.so hotpath.cpp
echo "Built libhotpath.so"

# Build native benchmark
g++ -O3 -o bench_native bench_native.cpp hotpath.cpp
echo "Built bench_native"

# Python bindings (if SWIG available)
if command -v swig &> /dev/null; then
    swig -python -c++ hotpath.i
    g++ -O3 -fPIC -shared \
        $(python3-config --includes) \
        hotpath_wrap.cxx hotpath.cpp \
        -o _hotpath.so
    echo "Built Python bindings"
fi

# Build Zig benchmark (if zig available)
if command -v zig &> /dev/null; then
    zig build-exe -OReleaseFast bench_zig.zig -femit-bin=bench_zig
    echo "Built bench_zig"
fi

# Build Rust benchmark (if rustc available)
if command -v rustc &> /dev/null; then
    rustc -O -o bench_rust bench_rust.rs
    echo "Built bench_rust"
fi

echo ""
echo "Run: ./bench_native"
echo "Run: python3 bench_python.py"
echo "Run: CGO_LDFLAGS='-L.' go run bench_go.go"
echo "Run: ./bench_zig"
echo "Run: ./bench_rust"
