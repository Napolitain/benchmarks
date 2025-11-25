import time

# Try to import SWIG bindings, skip FFI tests if not available
try:
    import hotpath
    HAS_FFI = True
except ImportError:
    HAS_FFI = False
    print("Note: hotpath module not found, skipping FFI tests\n")

# Native Python implementations
def fast_sum8_native(a, b, c, d, e, f, g, h):
    return a + b + c + d + e + f + g + h

def slow_compute_native(seed, iterations):
    h = seed
    for _ in range(iterations):
        h ^= h >> 33
        h = (h * 0xff51afd7ed558ccd) & 0xFFFFFFFFFFFFFFFF
        h ^= h >> 33
        h = (h * 0xc4ceb9fe1a85ec53) & 0xFFFFFFFFFFFFFFFF
        h ^= h >> 33
    return h

def main():
    print("Python FFI Benchmark")
    print("====================\n")

    FAST_ITERS = 1000000
    SLOW_ITERS = 100
    COMPUTE_ITERS = 1000000

    # === FAST function ===
    print("fast_sum8 (1M calls):")
    
    # Native Python
    start = time.perf_counter_ns()
    for _ in range(FAST_ITERS):
        result = fast_sum8_native(1, 2, 3, 4, 5, 6, 7, 8)
    end = time.perf_counter_ns()
    native_ns = (end - start) / FAST_ITERS
    native_total = (end - start) / 1_000_000
    print(f"  Native:  {native_total:8.2f} ms total, {native_ns:8.2f} ns/call")

    # FFI (C++)
    if HAS_FFI:
        start = time.perf_counter_ns()
        for _ in range(FAST_ITERS):
            result = hotpath.fast_sum8(1, 2, 3, 4, 5, 6, 7, 8)
        end = time.perf_counter_ns()
        ffi_ns = (end - start) / FAST_ITERS
        ffi_total = (end - start) / 1_000_000
        print(f"  FFI:     {ffi_total:8.2f} ms total, {ffi_ns:8.2f} ns/call")
        print(f"  Speedup: {native_ns/ffi_ns:.1f}x {'(FFI wins)' if ffi_ns < native_ns else '(Native wins)'}")

    # === SLOW function ===
    print(f"\nslow_compute (100 calls, 1M iters each):")
    
    # Native Python
    start = time.perf_counter_ns()
    for i in range(SLOW_ITERS):
        result = slow_compute_native(i, COMPUTE_ITERS)
    end = time.perf_counter_ns()
    native_ms = (end - start) / 1_000_000 / SLOW_ITERS
    native_total = (end - start) / 1_000_000
    print(f"  Native:  {native_total:8.2f} ms total, {native_ms:8.2f} ms/call")

    # FFI (C++)
    if HAS_FFI:
        start = time.perf_counter_ns()
        for i in range(SLOW_ITERS):
            result = hotpath.slow_compute(i, COMPUTE_ITERS)
        end = time.perf_counter_ns()
        ffi_ms = (end - start) / 1_000_000 / SLOW_ITERS
        ffi_total = (end - start) / 1_000_000
        print(f"  FFI:     {ffi_total:8.2f} ms total, {ffi_ms:8.2f} ms/call")
        print(f"  Speedup: {native_ms/ffi_ms:.1f}x {'(FFI wins)' if ffi_ms < native_ms else '(Native wins)'}")

if __name__ == "__main__":
    main()
