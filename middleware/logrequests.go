package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func logReq(r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	// use "log" to log the date & time
}

func logReqs(hfn http.HandlerFunc) http.HandlerFunc {
	// closure, so can access "hfn" in the inner function
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		start := time.Now()
		hfn(w, r)
		fmt.Printf("%v\n", time.Since(start))
	}
}

// http.Hanlder is an interface
// func logRequests(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s", r.Method, r.URL.Path) // pre-processing
// 		start := time.Now()
// 		handler.ServeHTTP(w, r)
// 		fmt.Printf("%v\n", time.Since(start)) // post-processing
// 	})
// }

func logRequests(logger *log.Logger) Adapter {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s", r.Method, r.URL.Path)
			start := time.Now()
			handler.ServeHTTP(w, r)
			logger.Printf("%v\n", time.Since(start))
		})
	}
}
