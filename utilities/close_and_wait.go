// Copyright 2024 Kirk Rader

package utilities

// Close the given values channel then wait for the given await channel to be
// closed. For example:
//
//	{
//	  values, await := StartWorker(handler)
//	  defer CloseAndWait(values, await)
//	  for _, value := range data {
//	    values <- value
//	  }
//	}
//
// See StartWorker
func CloseAndWait[V any](values chan<- V, await <-chan any) {
	close(values)
	<-await
}
