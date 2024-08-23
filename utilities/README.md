_Copyright 2024 Kirk Rader_

# parasaurolophus/utilities

```
package utilities // import "parasaurolophus/utilities"


FUNCTIONS

func Invoke[T any](handler func(T), value T, log io.Writer)
    Trap and log panics.

func NewWorker[V any](handler func(V), log io.Writer) (values chan<- V, await <-chan any)
    Start a goroutine which will invoke handler for each item sent to the values
    channel until it is closed. The goroutine will close the await channel
    before terminating.

func NewWorkerWaitGroup[V any](handler func(V), await *sync.WaitGroup, log io.Writer) (values chan<- V)
    Start a goroutine which will invoke handler for each item sent to the values
    channel until it is closed. The await's count will be incremented by for
    this function returns, and the goroutine will decrement await's count before
    terminating.


TYPES

type Watchdog struct {
        // Has unexported fields.
}
    Similar to standard time.Timer, but with the ability to suppress invocations
    of the timeout handler by sending values to a reset channel.

func NewWatchdog(interval time.Duration, timeout func()) (watchdog Watchdog, reset chan<- any)
    Construct a Watchdog. The timeout handler will be invoked periodically in
    its own goroutine except when suppressed by sending values to the given
    reset channel. Note that the time until the next timeout is reset to
    time.Now() each time a value is on the reset channel, so the exact frequency
    of timeouts is erratic, as determined by the base interval and the times at
    which a watchdog timer is reset.

func (watchdog Watchdog) Stop()
    Stop the watchdog and block until it exits.
```
