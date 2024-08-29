// Copyright 2024 Kirk Rader

package utilities

import (
	"sync"
	"time"
)

// Close all of the given values channels then wait for the given group to
// signal that they all have exited cleanly.
func CloseAllAndWait[V any](values []chan<- V, await *sync.WaitGroup) {
	for _, v := range values {
		close(v)
	}
	await.Wait()
}

// Close the given values channel then wait for the given await channel to be
// closed.
func CloseAndWait[V any](values chan<- V, await <-chan any) {
	close(values)
	<-await
}

// Start n+1 goroutines and wait for them all to complete after invoking the
// given generator function. The generator function must send values in a
// round-robin fashion to the set of values it is passed. The first n worker
// goroutines will send the result of invoking the given transformer function
// to the last goroutine, which invokes the given consumer function.
//
// ```mermaid
// ```
//
// See CloseAndWait, CloseAllAndWait, StartWorker, StartWorkers
func ProcessBatch[Input any, Output any](
	n int,
	generate func([]chan<- Input),
	transform func(Input) Output,
	consume func(Output),
) {

	consumer, awaitConsumer := StartWorker(consume)
	defer CloseAndWait(consumer, awaitConsumer)
	func() {
		produce := func(request Input) {
			consumer <- transform(request)
		}
		producers, awaitProducers := StartWorkers(n, produce)
		defer CloseAllAndWait(producers, awaitProducers)
		generate(producers)
	}()
}

// Start a goroutine which will invoke the given handler for each item sent to
// the returned values channel until it is closed at which time it will close
// the await channel before exiting.
//
// See:
// CloseAndWait[V any](chan<- V, <-chan any)
// StartWorkers[V any](int, func(V))
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

// Start the specified number of goroutines, each of which will invoke the
// given handler for each item sent to one of the returned values channels. The
// returned wait group's counter will be set initially to the number of
// goroutines specified by n. Each goroutine will decrement the returned wait
// group's counter before terminating when the value channel to which it is
// listening is closed. group's count
//
// See:
// CloseAllAndWait[V any]([]chan<- V, *sync.WaitGroup)
// StartWorker[V any](func(V))
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

// Invoke fn asynchronously. Return its value if it completes within the
// specified duration. Otherwise, return the value of calling the timeout
// function.
//
// Warning! the goroutine used to invoke fn constitutes a resource leak if it
// never completes, so use this function with caution. For example, it would be
// reasonable to invoke WithTimeLimit in a console application, a lambda's
// request handler function or any similar "one and done" flow. But Go provides
// no mechanism for forcibly terminating a goroutine, so long-running services
// should use this function (or any that involve goroutines that could possibly
// hang) with caution.
func WithTimeLimit[V any](fn func() V, timeout func() V, timeLimit time.Duration) V {
	value := make(chan V)
	timer := time.NewTimer(timeLimit)
	go func() { value <- fn() }()
	select {
	case v := <-value:
		return v
	case <-timer.C:
		return timeout()
	}
}
