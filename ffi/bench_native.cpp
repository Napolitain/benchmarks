#include <iostream>
#include <chrono>
#include <iomanip>
#include "hotpath.h"

using namespace std::chrono;

int main() {
    std::cout << std::fixed << std::setprecision(2);
    std::cout << "C++ Native Baseline\n";
    std::cout << "===================\n\n";

    // FAST function: 1M iterations
    const int FAST_ITERS = 1000000;
    volatile long long result = 0;
    
    auto start = high_resolution_clock::now();
    for (int i = 0; i < FAST_ITERS; i++) {
        result = fast_sum8(1, 2, 3, 4, 5, 6, 7, 8);
    }
    auto end = high_resolution_clock::now();
    double fast_ns = duration_cast<nanoseconds>(end - start).count() / (double)FAST_ITERS;
    double fast_total_ms = duration_cast<microseconds>(end - start).count() / 1000.0;
    
    std::cout << "fast_sum8 (1M calls):\n";
    std::cout << "  Total time:    " << fast_total_ms << " ms\n";
    std::cout << "  Per call:      " << fast_ns << " ns\n\n";

    // SLOW function: 100 iterations, 1M compute each
    const int SLOW_ITERS = 100;
    const int COMPUTE_ITERS = 1000000;
    
    start = high_resolution_clock::now();
    for (int i = 0; i < SLOW_ITERS; i++) {
        result = slow_compute(i, COMPUTE_ITERS);
    }
    end = high_resolution_clock::now();
    double slow_ms = duration_cast<microseconds>(end - start).count() / 1000.0 / SLOW_ITERS;
    double slow_total_ms = duration_cast<microseconds>(end - start).count() / 1000.0;
    
    std::cout << "slow_compute (100 calls, 1M iters each):\n";
    std::cout << "  Total time:    " << slow_total_ms << " ms\n";
    std::cout << "  Per call:      " << slow_ms << " ms\n";

    return 0;
}
