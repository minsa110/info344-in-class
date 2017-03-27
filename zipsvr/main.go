package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// * is a pointer (passing as reference instead of copying the request)
// slice is a reference point - never work with arrays in Go, work with slices
func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name") // r is a request object

	w.Header().Add("Content-Type", "text/plain") // header name and header value as arguments for Add

	w.Write([]byte("hello " + name)) // Write method takes in slice of bytes, not strings
}

func main() {
	addr := os.Getenv("ADDR") // same as var addr string = os...
	if len(addr) == 0 {
		log.Fatal("please set ADDR environment variable")
	}

	// first argument is the resource path
	// second argument is passing a pointer to the function here, instead of calling it
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("server is listening at %s ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
