#include "hotpath.h"

extern "C" {

// FAST: ~10ns - just sum 8 values, exposes FFI call overhead
long long fast_sum8(long long a, long long b, long long c, long long d,
                    long long e, long long f, long long g, long long h) {
    return a + b + c + d + e + f + g + h;
}

// SLOW: ~10ms - heavy compute, amortizes FFI overhead
// Uses a mixing function similar to MurmurHash
long long slow_compute(long long seed, int iterations) {
    long long h = seed;
    for (int i = 0; i < iterations; i++) {
        h ^= h >> 33;
        h *= 0xff51afd7ed558ccdULL;
        h ^= h >> 33;
        h *= 0xc4ceb9fe1a85ec53ULL;
        h ^= h >> 33;
    }
    return h;
}

}
