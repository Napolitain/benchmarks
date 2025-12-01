#ifndef HOTPATH_H
#define HOTPATH_H

#ifdef __cplusplus
extern "C" {
#endif

// FAST: ~10ns work per call - exposes FFI latency cost
// Sum 8 integers (fits in registers, minimal work)
long long fast_sum8(long long a, long long b, long long c, long long d,
                    long long e, long long f, long long g, long long h);

// SLOW: ~10ms work per call - amortizes FFI overhead
// Compute N iterations of a hash-like mixing function
long long slow_compute(long long seed, int iterations);

#ifdef __cplusplus
}
#endif

#endif // HOTPATH_H
