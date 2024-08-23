// Copyright 2024 Kirk Rader

package utilities

import "time"

// Similar to standard time.Timer, but with the ability to suppress invocations
// of the timeout handler by sending values to a reset channel.
type Watchdog struct {
	terminate chan any
	await     chan any
}

// Construct a Watchdog. The timeout handler will be invoked periodically in
// its own goroutine except when suppressed by sending values to the given
// reset channel. Note that the time until the next timeout is reset to
// time.Now() each time a value is on the reset channel, so the exact frequency
// of timeouts is erratic, as determined by the base interval and the times at
// which a watchdog timer is reset.
func NewWatchdog(interval time.Duration, timeout func()) (watchdog Watchdog, reset chan<- any) {
	r := make(chan any)
	reset = r
	watchdog = Watchdog{
		terminate: make(chan any),
		await:     make(chan any),
	}
	go func() {
		defer close(watchdog.await)
		start := time.Now()

		for {
			select {
			case <-watchdog.terminate:
				return
			case <-r:
				start = time.Now()
			default:
				if time.Since(start) > interval {
					timeout()
					start = time.Now()
				}
				time.Sleep(interval)
			}
		}
	}()
	return
}

// Stop the watchdog and block until it exits.
func (watchdog Watchdog) Stop() {
	defer func() {
		_ = recover()
	}()
	close(watchdog.terminate)
	<-watchdog.await
}
