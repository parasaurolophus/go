_Copyright &copy; Kirk Rader 2024_

# Wrapper for `log/slog`

Note: the primary source of documentation for this library is in the form of
inline comments that can be viewed using `go doc`.

```bash
cd logging
go doc -all
```

## Overview

The design of Go's standard `log/slog` package leaves much to be desired. This
is a very thin wrapper that helps address some (but far from all) of its
shortcomings.

- [Streamline syntax](#streamlined-syntax) for logger construction and usage

  - Hard-code use of `slog.JSONHandler`

  - Ensure `slog.LevelVar` exists for each logger

  - Encapsulate `context.Context` for use with each logger

  - Promote `slog.Level.Set` method to `logging.Logger` interface

- [Lazy evaluation](#lazy-evaluation) of message formatting code using closures

- [Verbosity-based nomenclature](#verbosity-based-nomenclature) rather than
  conflating "verbosity of log output" with "severity of issue"

- Support embedding [stack traces](#stack-traces) in log entries in a simple,
  flexible way

- Enhanced support for custom attributes

  - Special handling of ["tags" attributes](#special-handling-of-tags)

  - Declare [default tags and attributes](#default-tags-and-attributes) per
    logger instance

## Streamlined Syntax

Use:

```go
// do this...

logger := logging.New(os.Stdout, nil)
logger.Set(logging.TRACE)
```

Rather than:

```go
// ...instead of this

levelVar := slog.LevelVar{}

handlerOptions := slog.HandlerOptions{
    Level: &levelVar,
}

slogger := slog.New(slog.NewJSONHandler(os.Stdout, &handlerOptions))
levelVar.Set(slog.LevelDebug)
```

## Lazy Evaluation

Use:

```go
// do this...

logger.Fine(
    func() string {
        return fmt.Sprintf("log formatting overhead %#v", counters)
    })
```

Rather than:

```go
// ...instead of this

if slogger.Enabled(ctx, slog.LevelInfo) {
    slogger.Info(fmt.Sprintf("log formatting overhead %#v", counters))
}
```

## Verbosity Based Nomenclature

Use:

```go
// do this...
logger.Trace(func() string { return "..." })
logger.Fine(func() string { return "..." })
logger.Optional(func() string { return "..." })
logger.Always(func() string { return "..." })
```

Rather than:

```go
// ...instead of this
slogger.DebugContext("...")
slogger.InfoContext("...")
slogger.WarnContext("...")
slogger.ErrorContext("...")
```

## Attributes with Special Handling

### Stack Traces

The value of any attribute named "stacktrace" (`logging.STACKTRACE`) will be
replaced by a one-line stack trace that starts at the frame for the function
invoked the given logging method, e.g.:

```go
// include stack traces in log entries
logger.Fine(
    func() string {
        return "something bad happened"
    },
    // nil is replaced by a stack trace in the log entry
    logging.STACKTRACE, nil)
```

### Tags

The value of any attribute named "tags" (`logging.TAGS`) is expected to be a
slice of strings and has special formatting behavior as described below.

## Default Tags and Attributes

`logging.LoggerOptions` extends `slog.HandlerOptions` in various ways, including
by supporting `BaseTags` and `BaseAttributes` keys. Both are used to allow
default attributes to be added to every log entry emitted by a given
`logging.Logger` instance.

- `BaseAttributes` is of type `[]any` and will be appended to the attributes
  supplied exclicitly when a logging method is invoked

- `BaseTags`, is of type `[]string` and will be appended to the value of any
  "tags" attribute supplied when a logging method is invoked

For example:

```go
package main

import (
	"fmt"
	"os"
	"parasaurolophus/go/logging"
)

func main() {

	type Counters struct {
		Error1 int `json:"error1"`
		Error2 int `json:"error2"`
	}

	counters := Counters{}

	options := logging.LoggerOptions{
		BaseAttributes: []any{"counters", &counters},
		BaseTags:       []string{"foo", "bar"},
	}

	logger := logging.New(os.Stdout, &options)
	n := 42
	counters.Error1 += 1

	logger.Optional(
		func() string { return fmt.Sprintf("n = %d", n) },
		logging.STACKTRACE, nil,
		logging.TAGS, []string{"hoo"},
		"baz", "waka")
}
```

will emit a log entry like:

```json
{"time":"2024-02-11T06:51:51.00758035-06:00","verbosity":"OPTIONAL","msg":"n = 42","counters":{"error1":1,"error2":0},"baz":"waka","stacktrace":"5:main.main [/source/go/scratch/scratch.go:29] < 6:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 7:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]","tags":["foo","bar","hoo"]}
```

See [../example/example.go](../example/example.go) for a more complete set of
examples, including how to ensure that panics are logged with stack traces using
`logging.Logger.Defer().
