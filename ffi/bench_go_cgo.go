package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo LDFLAGS: -L. -lhotpath
#include "hotpath.h"
*/
import "C"
import (
	"fmt"
	"time"
)

// Native Go implementations (for comparison)
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
var csink C.longlong

func main() {
	fmt.Println("Go CGO FFI Benchmark")
	fmt.Println("====================\n")

	const FAST_ITERS = 1000000
	const SLOW_ITERS = 100
	const COMPUTE_ITERS = 1000000

	// === FAST function ===
	fmt.Println("fast_sum8 (1M calls):")

	// Native Go
	start := time.Now()
	for i := 0; i < FAST_ITERS; i++ {
		sink = fastSum8Native(1, 2, 3, 4, 5, 6, 7, 8)
	}
	elapsed := time.Since(start)
	nativeNs := float64(elapsed.Nanoseconds()) / FAST_ITERS
	nativeTotalMs := float64(elapsed.Microseconds()) / 1000.0
	fmt.Printf("  Native:  %8.2f ms total, %8.2f ns/call\n", nativeTotalMs, nativeNs)

	// FFI (C++)
	start = time.Now()
	for i := 0; i < FAST_ITERS; i++ {
		csink = C.fast_sum8(1, 2, 3, 4, 5, 6, 7, 8)
	}
	elapsed = time.Since(start)
	ffiNs := float64(elapsed.Nanoseconds()) / FAST_ITERS
	ffiTotalMs := float64(elapsed.Microseconds()) / 1000.0
	fmt.Printf("  FFI:     %8.2f ms total, %8.2f ns/call\n", ffiTotalMs, ffiNs)

	winner := "(FFI wins)"
	if nativeNs < ffiNs {
		winner = "(Native wins)"
	}
	fmt.Printf("  Speedup: %.1fx %s\n", nativeNs/ffiNs, winner)

	// === SLOW function ===
	fmt.Printf("\nslow_compute (100 calls, 1M iters each):\n")

	// Native Go
	start = time.Now()
	for i := 0; i < SLOW_ITERS; i++ {
		sink = slowComputeNative(int64(i), COMPUTE_ITERS)
	}
	elapsed = time.Since(start)
	nativeMs := float64(elapsed.Microseconds()) / 1000.0 / SLOW_ITERS
	nativeTotalMs = float64(elapsed.Microseconds()) / 1000.0
	fmt.Printf("  Native:  %8.2f ms total, %8.2f ms/call\n", nativeTotalMs, nativeMs)

	// FFI (C++)
	start = time.Now()
	for i := 0; i < SLOW_ITERS; i++ {
		csink = C.slow_compute(C.longlong(i), COMPUTE_ITERS)
	}
	elapsed = time.Since(start)
	ffiMs := float64(elapsed.Microseconds()) / 1000.0 / SLOW_ITERS
	ffiTotalMs = float64(elapsed.Microseconds()) / 1000.0
	fmt.Printf("  FFI:     %8.2f ms total, %8.2f ms/call\n", ffiTotalMs, ffiMs)

	winner = "(FFI wins)"
	if nativeMs < ffiMs {
		winner = "(Native wins)"
	}
	fmt.Printf("  Speedup: %.1fx %s\n", nativeMs/ffiMs, winner)
}
