Copyright 2024 Kirk Rader

# Automation Library

- [./trigger/](./trigger/) exposes `automation.SendTriggerEvents`, a function
  which, as the name implies, launches a worker goroutine to periodically send
  `automation.Trigger` events at specific times over the course of the current
  day.

- [./hue/](./hue/) contains a wrapper for the V2 API exposed by Philips Hue
  Bridges, including support for receiving SSE messages asynchronously.

- [./powerview/](./powerview/) contains a wrapper for the API exposed by a
  Hunter-Douglas (PowerView) Hub for controlling motroized window shades.

- [./main/](./main/) contains a console application for manual integration
  testing all of the above.

## Requirements

This began as an idle experiment to see what it would be like to replace any
part of my fairly elaborate, Node-RED based [home automation
system](https://github.com/parasaurolophus/home-automation) with code written
in Go. I began by focusing on idiomatic Go implementations of the core back end
features of the existing system:

- Access to the data models provided by some specific "smart" device
  manufacturers' API's (Philips Hue and Hunter-Douglas PowerView)

- Subscribe to asynchronous SSE messages that are part of one such API (Philips Hue)

- Send asyncrhonous messages at specific times of day to drive lighting and
  window covering automation

## General Principles

By "idiomatic" Go code, I mean:

- Use goroutines and channels for asynchronous operations

- Rely on communication via channels rather than access to shared memory in
  goroutines to the greatest degree possible

This led inevitably to a functional programming paradigm rather than, for
example, a collection of "wrapper" structs with methods (i.e. the closest Go
comes to an object-oriented approach in the sense of languages like C++, Java
etc.)

For example, starting the time-of-day event stream is accomplished by calling
the `trigger.SendTriggerEvents` function.

That function takes only scalar argument types (two `float64`'s and an `int`),
all passed by value. It returns three channels used  to interact with the
goroutine it launches as a side-effect.

Two of the returned channels are used only for goroutine control and
synchronization and are never used to actually transmit any data between
goroutines.

The other channel returned by `trigger.SendTriggerEvents` is used to deliver
the actual trigger messages asynchronously. But the data transmitted is, again,
simply an alias for `string`, passed by value (i.e. not using pointers,
interfaces or other reference types).

Meanwhile, the implementation of the goroutine launched as a side-effect of
calling `trigger.SendTriggerEvents` is a closure that accesses no global
variables or other state that is maintained outside of its own lexical scope.

Note that this is a much higher degree of encapsulation than could be achieved
using, for example, private fields of a `struct`, where the containing data
structure exists outside of the lexical scope of any of its methods. This does
not by itself guarantee that the internal state is immune from race conditions,
deadly embrace and the like, but it makes it far easier to enforce the
constraints necessary to make it so.

In the case of `trigger.SendTriggerEvents`, note that other than the channels
(which, by definition, must be shared across goroutine boundaries in order to
be useful) all of the state managed by the worker goroutine (i.e. the `suncalc`
"times" map and the pointers to multiple `time.Timer` structures) is enclosed
within its own innermost lexical scope. Since only one goroutine can see this
data, it could only become subject to multi-tasking synchronization issues if
it were shared by some other means (such as being stored in fields of a
`struct`). Since the code was written explicitly _not_ to share that state in
any such way, no mutexes or similar synchronization constructs are required.

I.e. the hallmark of good Go style is passing data by value as arguments to
functions, returned values or across channels between goroutines while avoiding
global variables, shared access to data via pointers, interfaces or similar
reference types.

The same approach was used for the Hue and PowerView wrappers. Note that the
only `struct` definitions that appear anywhere in this library represent
various parts of the data models defined by the respective device API's. The
functions in this library construct instances of these data models, but once
they are returned by a given function or transmitted by a goroutine over a
given channel, that particular instance is "forgotten" by the function that
constructed it. No such data structures are part of the internal state of any
of this library's goroutines or lexical closures.

The most complex (not to say fraught) example of these principles is the
`hue.SubscribeSSE` function. It launches not one but two goroutines per
invocation.
