// Copyright 2024 Kirk Rader

package utilities

// Start a goroutine which will invoke the given handler for each item sent to
// the returned values channel, until it is closed, at which time it will close
// the await channel before exiting.
//
//	{
//	  values, await := StartWorker(bufferSize, handler)
//	  defer CloseAndWait(values, await)
//	  for _, value := range data {
//	    values <- value
//	  }
//	}
//
// See CloseAndWait, StartWorkers
func StartWorker[V any](

	bufferSize int,
	handler func(V),

) (

	values chan<- V,
	await <-chan any,

) {

	v := make(chan V, bufferSize)
	values = v
	a := make(chan any)
	await = a
	go func() {
		defer close(a)
		for value := range v {
			handler(value)
		}
	}()
	return
}
