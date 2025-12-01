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
	const SLOW_ITERS = 100
	const COMPUTE_ITERS = 1000000

	start := time.Now()
	for i := 0; i < SLOW_ITERS; i++ {
		csink = C.slow_compute(C.longlong(i), COMPUTE_ITERS)
	}
	elapsed := time.Since(start)
	
	ffiMs := float64(elapsed.Microseconds()) / 1000.0 / SLOW_ITERS
	ffiTotalMs := float64(elapsed.Microseconds()) / 1000.0
	fmt.Printf("slow_compute: %.2f ms total, %.2f ms/call\n", ffiTotalMs, ffiMs)
}
