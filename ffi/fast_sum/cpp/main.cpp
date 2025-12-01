#include <iostream>
#include <chrono>
#include <iomanip>
#include "../hotpath.h"

using namespace std::chrono;

int main() {
    const int FAST_ITERS = 1000000;
    volatile long long result = 0;
    
    auto start = high_resolution_clock::now();
    for (int i = 0; i < FAST_ITERS; i++) {
        result = fast_sum8(1, 2, 3, 4, 5, 6, 7, 8);
    }
    auto end = high_resolution_clock::now();
    
    double fast_ns = duration_cast<nanoseconds>(end - start).count() / (double)FAST_ITERS;
    double fast_total_ms = duration_cast<microseconds>(end - start).count() / 1000.0;
    
    std::cout << std::fixed << std::setprecision(2);
    std::cout << "fast_sum8: " << fast_total_ms << " ms total, " << fast_ns << " ns/call\n";

    return 0;
}
