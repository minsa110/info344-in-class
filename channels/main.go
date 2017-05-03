package main

import (
	"fmt"
	"math/rand"
	"time"
)

//someLongFunc is a function that might
//take a while to complete, so we want
//to run it on its own go routine
func someLongFunc(ch chan int) { // channel of ints
	r := rand.Intn(2000)
	d := time.Duration(r)
	time.Sleep(time.Millisecond * d)
	ch <- r // write the value into the channel
}

func main() {
	//create a channel and call
	//someLongFunc() on a go routine
	//passing the channel so that
	//someLongFunc() can communicate
	//its results
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Starting long-running func...")
	n := 10

	// ch := make(chan int) // creates unbuffered channel
	ch := make(chan int, n) // making buffered channel with capacity of n
	start := time.Now()
	for i := 0; i < n; i++ {
		go someLongFunc(ch)
		// if just "someLongFunc(ch)", then runs in serial, and doesn't actually
		// print until everything is done
	}
	for i := 0; i < n; i++ {
		result := <-ch
		fmt.Printf("Result was %d\n", result)
	}
	fmt.Printf("Took %v\n", time.Since(start))
}
