import time

def fast_sum8(a, b, c, d, e, f, g, h):
    return a + b + c + d + e + f + g + h

def main():
    FAST_ITERS = 1000000

    start = time.perf_counter_ns()
    for _ in range(FAST_ITERS):
        result = fast_sum8(1, 2, 3, 4, 5, 6, 7, 8)
    end = time.perf_counter_ns()
    
    total_ms = (end - start) / 1_000_000
    per_call_ns = (end - start) / FAST_ITERS
    print(f"fast_sum8: {total_ms:.2f} ms total, {per_call_ns:.2f} ns/call")

if __name__ == "__main__":
    main()
