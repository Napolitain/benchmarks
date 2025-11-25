# FFI Overhead Benchmark

Compares **native** implementations vs **FFI** (calling C++ via SWIG/CGO).

## Files

| File | Description |
|------|-------------|
| `bench_native.cpp` | C++ baseline |
| `bench_rust.rs` | Rust native |
| `bench_go.go` | Go native (no CGO) |
| `bench_go_cgo.go` | Go with CGO (native vs FFI) |
| `bench_python.py` | Python (native vs FFI) |
| `bench_node.ts` | Node.js/TypeScript native |
| `BenchJava.java` | Java native |

## Functions

| Function | Work per call | Calls | Purpose |
|----------|---------------|-------|---------|
| `fast_sum8()` | ~1ns (sum 8 ints) | 1M | FFI overhead dominates |
| `slow_compute()` | ~2ms (1M hash ops) | 100 | Compute dominates, FFI worth it |

## Build & Run

```bash
# C++ baseline
g++ -O3 -o bench_native bench_native.cpp hotpath.cpp
./bench_native

# Rust native
rustc -O -o bench_rust bench_rust.rs
./bench_rust

# Go native (no CGO)
go run bench_go.go

# Go with CGO (requires shared library)
g++ -O3 -fPIC -shared -o libhotpath.so hotpath.cpp   # Linux
g++ -O3 -shared -o hotpath.dll hotpath.cpp           # Windows
LD_LIBRARY_PATH=. go run bench_go_cgo.go             # Linux
set PATH=%CD%;%PATH% && go run bench_go_cgo.go       # Windows

# Node.js/TypeScript
npx ts-node bench_node.ts
# or
npx tsc bench_node.ts && node bench_node.js

# Java
javac BenchJava.java && java BenchJava

# Python native (works without SWIG)
python bench_python.py

# Python with FFI (requires SWIG)
swig -python -c++ hotpath.i
g++ -O3 -fPIC -shared $(python3-config --includes) hotpath_wrap.cxx hotpath.cpp -o _hotpath.so
python bench_python.py
```

## Expected Results

```
fast_sum8 (1M calls):
                Total       Per-call
C++             0.7ms       0.7ns
Rust            1.5ms       1.5ns
Go              1.0ms       1.0ns
Java            4.0ms       4.0ns
Node.js         25ms        25ns      <- BigInt overhead
Go+CGO Native   1.0ms       1.0ns
Go+CGO FFI      50ms        50ns      <- CGO overhead!
Python Native   55ms        55ns
Python FFI      200ms       200ns     <- SWIG overhead

slow_compute (100 calls × 1M iters):
                Total       Per-call
C++             220ms       2.2ms
Rust            160ms       1.6ms
Go              220ms       2.2ms
Java            235ms       2.4ms
Node.js         TBD         TBD       <- BigInt is slow
Go+CGO FFI      220ms       2.2ms     <- FFI overhead negligible
Python Native   19000ms     190ms
Python FFI      220ms       2.2ms     <- 100x faster!
```

## Conclusion

- **Cheap operations**: Stay native, FFI/CGO overhead kills performance
- **Expensive operations**: FFI to C/C++ massively benefits Python
- **Go**: Native Go is fast enough, CGO adds ~50ns overhead per call
- **Rust**: Comparable to C++, no FFI needed
