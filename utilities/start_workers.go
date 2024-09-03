// Copyright 2024 Kirk Rader

package utilities

import (
	"sync"
)

// Start the specified number of goroutines, each of which will invoke the
// given handler for each item sent to one of the returned values channels. The
// returned wait group's counter will be set initially to the number of
// goroutines specified by n. Each goroutine will decrement the returned wait
// group's counter before terminating when the value channel to which it is
// listening is closed. group's count. For example:
//
//	{
//	  values, await := StartWorkers(numWorkers, bufferSize, handler)
//	  defer CloseAllAndWait(values, await)
//	  for i, value := range data {
//	    values[i%numWorkers] <- value
//	  }
//	}
//
// See CloseAllAndWait, StartWorker
func StartWorkers[V any](

	numWorkers int,
	bufferSize int,
	handler func(V),

) (

	values []chan<- V,
	await *sync.WaitGroup,

) {

	v := make([]chan V, numWorkers)
	values = make([]chan<- V, numWorkers)
	await = &sync.WaitGroup{}
	await.Add(numWorkers)
	for i := range numWorkers {
		v[i] = make(chan V, bufferSize)
		values[i] = v[i]
		go func() {
			defer await.Done()
			for value := range v[i] {
				handler(value)
			}
		}()
	}
	return
}
