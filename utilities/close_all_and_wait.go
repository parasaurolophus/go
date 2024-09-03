// Copyright 2024 Kirk Rader

package utilities

import (
	"sync"
)

// Close all of the given values channels then wait for the given group to
// signal that all workers have exited cleanly. For example:
//
//	{
//	  values, await := StartWorkers(n, handler)
//	  defer CloseAllAndWait(values, await)
//	  for i, value := range data {
//	    values[i%n] <- value
//	  }
//	}
//
// See StartWorkers
func CloseAllAndWait[V any](values []chan<- V, await *sync.WaitGroup) {
	for _, v := range values {
		close(v)
	}
	await.Wait()
}
