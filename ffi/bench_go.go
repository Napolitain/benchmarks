package main

import (
	"fmt"
	"time"
)

// Native Go implementations
//go:noinline
func fastSum8Native(a, b, c, d, e, f, g, h int64) int64 {
	return a + b + c + d + e + f + g + h
}

//go:noinline
func slowComputeNative(seed int64, iterations int) int64 {
	h := uint64(seed)
	for i := 0; i < iterations; i++ {
		h ^= h >> 33
		h *= 0xff51afd7ed558ccd
		h ^= h >> 33
		h *= 0xc4ceb9fe1a85ec53
		h ^= h >> 33
	}
	return int64(h)
}

var sink int64

func main() {
	fmt.Println("Go Native Benchmark")
	fmt.Println("===================\n")

	const FAST_ITERS = 1000000
	const SLOW_ITERS = 100
	const COMPUTE_ITERS = 1000000

	// === FAST function ===
	fmt.Println("fast_sum8 (1M calls):")

	start := time.Now()
	for i := 0; i < FAST_ITERS; i++ {
		sink = fastSum8Native(1, 2, 3, 4, 5, 6, 7, 8)
	}
	elapsed := time.Since(start)
	nativeNs := float64(elapsed.Nanoseconds()) / FAST_ITERS
	nativeTotalMs := float64(elapsed.Microseconds()) / 1000.0
	fmt.Printf("  Total time:  %8.2f ms\n", nativeTotalMs)
	fmt.Printf("  Per call:    %8.2f ns\n", nativeNs)

	// === SLOW function ===
	fmt.Printf("\nslow_compute (100 calls, 1M iters each):\n")

	start = time.Now()
	for i := 0; i < SLOW_ITERS; i++ {
		sink = slowComputeNative(int64(i), COMPUTE_ITERS)
	}
	elapsed = time.Since(start)
	nativeMs := float64(elapsed.Microseconds()) / 1000.0 / SLOW_ITERS
	nativeTotalMs = float64(elapsed.Microseconds()) / 1000.0
	fmt.Printf("  Total time:  %8.2f ms\n", nativeTotalMs)
	fmt.Printf("  Per call:    %8.2f ms\n", nativeMs)
}
