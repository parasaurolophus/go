// Copyright 2024 Kirk Rader

package utilities

import (
	"encoding/csv"
	"io"
	"sync"
	"time"
)

// ----------------------------------------------------------------------------
//
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
// ----------------------------------------------------------------------------
func CloseAllAndWait[V any](values []chan<- V, await *sync.WaitGroup) {
	for _, v := range values {
		close(v)
	}
	await.Wait()
}

// ----------------------------------------------------------------------------
//
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
// ----------------------------------------------------------------------------
func CloseAndWait[V any](values chan<- V, await <-chan any) {
	close(values)
	<-await
}

// ----------------------------------------------------------------------------
//
// Return a function for use as the generate parameter to ProcessBatch. The
// returned function will invoke the given parse for each row of the given CSV
// file, sending the result to the batch's transformers channels. Any errors
// encountered along the way will be passed to the given errorHandler function.
//
// [See] ProcessBatch
//
// ----------------------------------------------------------------------------
func MakeCSVGenerator(

	reader *csv.Reader,
	headers []string,
	errorHandler func(error),

) (

	generator func([]chan<- map[string]string),
	err error,

) {

	generator = func(transformers []chan<- map[string]string) {
		var err error
		row := 1
		n := len(transformers)
		for {
			var columns []string
			columns, err = reader.Read()
			if err != nil {
				if err != io.EOF {
					errorHandler(err)
				}
				break
			}
			m := map[string]string{}
			for i, h := range headers {
				m[h] = columns[i]
			}
			transformers[(row-1)%n] <- m
			row++
		}
	}
	return
}

// ----------------------------------------------------------------------------
//
// Process items in a set of data concurrently. Specifically, start n+1
// goroutines and wait for them all to complete after invoking the given
// generator function. The generator function must send input values in a
// round-robin fashion to the set of transformer channels it is passed. The
// transformer goroutines will send the result of invoking the given transform
// function to the consumer goroutine, which invokes the given consume
// function:
//
//	                   +-----------+
//	              +-->>| transform |----+
//	              |    +-----------+    |
//	              |          .          |
//	+----------+  |          .          |    +---------+
//	| generate |--+     concurrent      |-->>| consume |
//	+----------+  |     goroutines      |    +---------+
//	              |          .          |
//	              |          .          |
//	              |    +-----------+    |
//	              +-->>| transform |----+
//	                   +-----------+
//
// ProcessBatch will call the finish function after all itsworker goroutines
// have terminated.
//
// Note that this function will hang if any of the generate, transform or
// consume functions do not return. If your transform function invokes some SDK
// function or API that can hang, consier the use of WithTimeLimit to allow the
// batch to run to completion even if some operations would otherwise block it
// (but then be aware of the consequences of resulting resource leaks).
//
// [See] CloseAndWait, CloseAllAndWait, StartWorker, StartWorkers,
// TestProcessBatch
//
// ----------------------------------------------------------------------------
func ProcessBatch[Input any, Output any](
	n int,
	generate func(transformers []chan<- Input),
	transform func(Input) Output,
	consume func(Output),
) {

	// start a goroutine that will apply the consume function to each value
	// sent to its channel
	consumer, awaitConsumer := StartWorker(consume)
	defer CloseAndWait(consumer, awaitConsumer)

	// wrap the transform function in a closure that will send a given
	// transformed input to the consumer channel
	produce := func(request Input) {
		consumer <- transform(request)
	}

	// start n goroutines each of which will call the produce closure for each
	// value sent to its channel
	transformers, awaitTransformers := StartWorkers(n, produce)
	defer CloseAllAndWait(transformers, awaitTransformers)

	// generate must send values of type Input to the channels it is
	// passed,then return
	generate(transformers)
}

// ----------------------------------------------------------------------------
//
// Start a goroutine which will invoke the given handler for each item sent to
// the returned values channel, until it is closed, at which time it will close
// the await channel before exiting.
//
//	{
//	  values, await := StartWorker(handler)
//	  defer CloseAndWait(values, await)
//	  for _, value := range data {
//	    values <- value
//	  }
//	}
//
// [See] CloseAndWait, StartWorkers
//
// ----------------------------------------------------------------------------
func StartWorker[V any](handler func(V)) (values chan<- V, await <-chan any) {
	v := make(chan V)
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

// ---------------------------------------------------------------------------
//
// Start the specified number of goroutines, each of which will invoke the
// given handler for each item sent to one of the returned values channels. The
// returned wait group's counter will be set initially to the number of
// goroutines specified by n. Each goroutine will decrement the returned wait
// group's counter before terminating when the value channel to which it is
// listening is closed. group's count. For example:
//
//	{
//	  values, await := StartWorkers(n, handler)
//	  defer CloseAllAndWait(values, await)
//	  for i, value := range data {
//	    values[i%n] <- value
//	  }
//	}
//
// [See] CloseAllAndWait, StartWorker
//
// ---------------------------------------------------------------------------
func StartWorkers[V any](n int, handler func(V)) (values []chan<- V, await *sync.WaitGroup) {
	v := make([]chan V, n)
	values = make([]chan<- V, n)
	await = &sync.WaitGroup{}
	await.Add(n)
	for i := range n {
		v[i] = make(chan V)
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

// ---------------------------------------------------------------------------
//
// Invoke fn asynchronously. Return its value if it completes within the
// specified duration. Otherwise, return the value of calling the timeout
// function.
//
// Warning! the goroutine used to invoke fn constitutes a resource leak if it
// never completes, so use this function with caution. For example, it would be
// reasonable to invoke WithTimeLimit in a console application, a lambda's
// request handler function or any similar "one and done" flow. But Go provides
// no mechanism for forcibly terminating a goroutine, so long-running processes
// should not use this function (or any that involve goroutines that could
// possibly hang).
//
// ---------------------------------------------------------------------------
func WithTimeLimit[V any](fn func() V, timeout func(time.Time) V, timeLimit time.Duration) V {
	value := make(chan V)
	timer := time.NewTimer(timeLimit)
	go func() { value <- fn() }()
	select {
	case v := <-value:
		return v
	case t := <-timer.C:
		return timeout(t)
	}
}
