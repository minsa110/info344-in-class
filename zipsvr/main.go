package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type zip struct {
	Zip   string `json:"zip"`   // field names are upper-case bc JSON coming in
	City  string `json:"city"`  // must be exported so that it can be decoded0
	State string `json:"state"` // <-- name here MUST be what it is in JSON
}

// * for creating pointers, & for using pointers (returns pointer of the element (i.e. &zips))
type zipSlice []*zip
type zipIndex map[string]zipSlice // key is string (i.e. city name) && value is zipSlice

// if `json:"-"` then ignores from exporting (i.e. passwords)

// * is a pointer (passing as reference instead of copying the request)
// slice is a reference point - never work with arrays in Go, work with slices
func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name") // r is a request object

	w.Header().Add("Content-Type", "text/plain") // header name and header value as arguments for Add

	w.Write([]byte("hello " + name)) // Write method takes in slice of bytes, not strings
}

// ResponseWrite is an interface
// pointer to Request because that struct is super big
// with zi as a receiver (can only have ONE parameter) (it's "this" in Java)
func (zi zipIndex) zipsForCityHandler(w http.ResponseWriter, r *http.Request) {
	// /zips/city/seattle
	_, city := path.Split(r.URL.Path)
	lcity := strings.ToLower(city)

	w.Header().Add("Content-Type", "application/json; charset=utf-8") // for the client :)
	w.Header().Add("Access-Control-Allow-Origin", "*")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(zi[lcity]); err != nil {
		http.Error(w, "error encoding json: "+err.Error(), http.StatusInternalServerError)
		// our fault, so internal, should be 5XX
	}
}

func main() {
	addr := os.Getenv("ADDR") // same as var addr string = os...
	if len(addr) == 0 {
		log.Fatal("please set ADDR environment variable")
	}

	// JSON decoding
	f, err := os.Open("../data/zips.json") // this opens the file
	if err != nil {
		log.Fatal("error oepning zips file: " + err.Error())
	}

	// this receives the data by creating slices
	// make slices of pointers to zips && set number of slices so we don't have to allocate more spaces
	// Go knows the type "zip" because we created a struct above
	zips := make(zipSlice, 0, 43000) // make(type/pointer, length, capacity)
	decoder := json.NewDecoder(f)    // decoder decodes json into slices
	if err := decoder.Decode(&zips); err != nil {
		log.Fatal("error decoding zips json: " + err.Error())
	}
	fmt.Printf("loaded %d zips\n", len(zips))

	zi := make(zipIndex)
	// pointer of slice to a pointer of slices ???????????????????????????????????????????????????????????????????????

	// kind of like for each loop
	for _, z := range zips { // using range operator (index (here it's ignored by using "_"), value)
		lower := strings.ToLower(z.City)
		zi[lower] = append(zi[lower], z)
	}

	fmt.Printf("there are %d zips in Seattle\n", len(zi["seattle"])) // test
	// this prints the above statement on the command line

	// first argument is the resource path
	// second argument is passing a pointer to the function here, instead of calling it
	http.HandleFunc("/hello", helloHandler)

	http.HandleFunc("/zips/city/", zi.zipsForCityHandler) // since it ends with a slash, it means it starts with "/zip/city"

	fmt.Printf("server is listening at %s ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
