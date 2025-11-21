package main

import (
	"fmt"
	"net/http"
	"runtime"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"Hello, World!"}`)
}

func main() {
	runtime.GOMAXPROCS(1)
	http.HandleFunc("/", helloHandler)
	fmt.Println("Go server listening on :8080 (single-threaded)")
	http.ListenAndServe(":8080", nil)
}
