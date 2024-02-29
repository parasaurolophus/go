// Copyright Kirk Rader 2024

package main

import (
	"fmt"
)

// Send 0 ... 9 to ch then close it.
func worker(ch chan int) {

	// Use defer to ensure that ch is always closed so as to reduce the risk of
	// deadlocks.
	defer close(ch)

	for n := 0; n < 10; n += 1 {

		ch <- n
	}
}

// Prints
//
//	0
//	1
//	2
//	3
//	4
//	5
//	6
//	7
//	8
//	9
//
// to stdout.
func main() {

	ch := make(chan int)

	// A function invoked using the go keyword is launched as a concurrently
	// running goroutine.
	go worker(ch)

	// A loop that "ranges over" a channel is Go's standard idiom for consuming
	// values from that channel. The channel will block the main goroutine at
	// each iteration until the worker goroutine sends a message and the loop
	// terminates when the channel is closed. I.e. channels are used both to
	// communicate between goroutines and to synchronize their activities.
	for v := range ch {

		fmt.Printf("\t%2d\n", v)
	}
}
