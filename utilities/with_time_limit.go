// Copyright 2024 Kirk Rader

package utilities

import (
	"time"
)

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
func WithTimeLimit[V any](

	fn func() V,
	timeout func(time.Time) V,
	timeLimit time.Duration,

) V {

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
