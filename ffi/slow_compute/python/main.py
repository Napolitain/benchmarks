import time

def slow_compute(seed, iterations):
    h = seed
    for _ in range(iterations):
        h ^= h >> 33
        h = (h * 0xff51afd7ed558ccd) & 0xFFFFFFFFFFFFFFFF
        h ^= h >> 33
        h = (h * 0xc4ceb9fe1a85ec53) & 0xFFFFFFFFFFFFFFFF
        h ^= h >> 33
    return h

def main():
    SLOW_ITERS = 100
    COMPUTE_ITERS = 1000000

    start = time.perf_counter_ns()
    for i in range(SLOW_ITERS):
        result = slow_compute(i, COMPUTE_ITERS)
    end = time.perf_counter_ns()
    
    total_ms = (end - start) / 1_000_000
    per_call_ms = total_ms / SLOW_ITERS
    print(f"slow_compute: {total_ms:.2f} ms total, {per_call_ms:.2f} ms/call")

if __name__ == "__main__":
    main()
