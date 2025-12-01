#include <iostream>
#include <chrono>
#include <iomanip>
#include "../hotpath.h"

using namespace std::chrono;

int main() {
    const int SLOW_ITERS = 100;
    const int COMPUTE_ITERS = 1000000;
    volatile long long result = 0;
    
    auto start = high_resolution_clock::now();
    for (int i = 0; i < SLOW_ITERS; i++) {
        result = slow_compute(i, COMPUTE_ITERS);
    }
    auto end = high_resolution_clock::now();
    
    double slow_ms = duration_cast<microseconds>(end - start).count() / 1000.0 / SLOW_ITERS;
    double slow_total_ms = duration_cast<microseconds>(end - start).count() / 1000.0;
    
    std::cout << std::fixed << std::setprecision(2);
    std::cout << "slow_compute: " << slow_total_ms << " ms total, " << slow_ms << " ms/call\n";

    return 0;
}
