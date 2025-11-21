package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/valyala/fasthttp"
)

func helloHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(`{"message":"Hello, World!"}`)
}

func main() {
	runtime.GOMAXPROCS(1)
	fmt.Println("Go fasthttp server listening on :8080 (single-threaded)")
	if err := fasthttp.ListenAndServe(":8080", helloHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
