// Copyright 2024 Kirk Rader

package utilities

import (
	"io"
	"sync"
)

// Start a goroutine which will invoke handler for each item sent to the values
// channel until it is closed. The goroutine will close the await channel
// before terminating.
func NewWorker[V any](handler func(V), log io.Writer) (values chan<- V, await <-chan any) {
	v := make(chan V)
	values = v
	a := make(chan any)
	await = a
	go func() {
		defer close(a)
		for value := range v {
			Invoke(handler, value, log)
		}
	}()
	return
}

// Start a goroutine which will invoke handler for each item sent to the values
// channel until it is closed. The await's count will be incremented by for
// this function returns, and the goroutine will decrement await's count before
// terminating.
func NewWorkerWaitGroup[V any](handler func(V), await *sync.WaitGroup, log io.Writer) (values chan<- V) {
	v := make(chan V)
	values = v
	await.Add(1)
	go func() {
		defer await.Done()
		for value := range v {
			Invoke(handler, value, log)
		}
	}()
	return
}
