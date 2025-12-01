package main

/*
#cgo CXXFLAGS: -std=c++11 -O3
#cgo LDFLAGS: -L.. -lhotpath
#include "../hotpath.h"
*/
import "C"
import (
	"fmt"
	"time"
)

var csink C.longlong

func main() {
	const FAST_ITERS = 1000000

	start := time.Now()
	for i := 0; i < FAST_ITERS; i++ {
		csink = C.fast_sum8(1, 2, 3, 4, 5, 6, 7, 8)
	}
	elapsed := time.Since(start)
	
	ffiNs := float64(elapsed.Nanoseconds()) / FAST_ITERS
	ffiTotalMs := float64(elapsed.Microseconds()) / 1000.0
	fmt.Printf("fast_sum8: %.2f ms total, %.2f ns/call\n", ffiTotalMs, ffiNs)
}
