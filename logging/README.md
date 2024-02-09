_Copyright &copy; Kirk Rader 2024_

# Wrapper for `log/slog`

Note: the Go code in this repository contains extensive inline comments that can
be accessed using the `go doc` command.

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
    ctx,
    func() string {
        return fmt.Sprintf("log formatting overhead %#v", counters)
    })
```

Rather than:

```go
// ...instead of this

if slogger.Enabled(ctx, slog.LevelInfo) {
    slogger.InfoContext(ctx, fmt.Sprintf("log formatting overhead %#v", counters))
}
```

## Verbosity Based Nomenclature

Use:

```go
// do this...
logger.Trace(ctx, func() string { return "..." })
logger.Fine(ctx, func() string { return "..." })
logger.Optional(ctx, func() string { return "..." })
logger.Always(ctx, func() string { return "..." })
```

Rather than:

```go
// ...instead of this
slogger.DebugContext(ctx, "...")
slogger.InfoContext(ctx, "...")
slogger.WarnContext(ctx, "...")
slogger.ErrorContext(ctx, "...")
```

## Logging Panics

```go
defer logger.OnPanic(
    ctx,
    func(r any) string {
        return fmt.Sprintf("panic: %v", r)
    },
    logging.STACKTRACE, nil,
    logging.TAGS, []string{"ERROR", "PANIC", "SEVERE"},
)
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
type (
    errorCounters struct {
        Error1 int `json:"error1"`
        Error2 int `json:"error2"`
    }
)

counters := errorCounters{}

options := logging.LoggerOptions{
    BaseAttributes: []any{"counters", &counters},
    BaseTags:       []string{"EXAMPLE"},
}

logger := logging.New(os.Stdout, &options)
counters.Error1 += 1

logger.Fine(
    func() string {
        return fmt.Sprintf("error 1 has occured %d times", counters.Error1)
    },
    logging.TAGS, []string{"FOR_README"})
```

will emit a log entry like:

```json
{"time":"2024-02-03T04:41:20.757272556-06:00","verbosity":"FINE","msg":"error 1 has occured 1 times","counters":{"error1":1,"error2":0},"tags":["EXAMPLE","FOR_README"]}
```

## Go Docs

```bash
logging $ go doc -all
```

```
package logging // import "parasaurolophus/go/logging"


CONSTANTS

const (
	TRACE    = Verbosity(slog.LevelDebug)
	FINE     = Verbosity(slog.LevelInfo)
	OPTIONAL = Verbosity(slog.LevelWarn)
	ALWAYS   = Verbosity(slog.LevelError)
)
    Mapping of Verbosity to slog.Level values.

const (

	// Value of "stacktrace" attributes will be replaced with one-line stack
	// traces for the function that called the given logging method.
	STACKTRACE = "stacktrace"

	// Value will be merged with the currently configured base tags.
	TAGS = "tags"
)
    Specially handled attributes


TYPES

type Logger struct {
	// Has unexported fields.
}
    Wrapper for an instance of slog.Logger.

func New(writer io.Writer, options *LoggerOptions) *Logger
    Returns a newly created, wrapped instance of slog.Logger.

    Log entries written using the returned Logger instance will have
    "verbosity" attributes instead of "level" attributes and the values of their
    "stacktrace" attributes, if present, will be replaced by an invocation of
    ShortStackTrace(caller) where caller is the name of the function that calls
    a logging method. The final set of attributes for each log entry will be
    the result of combining the value of LoggerOptions.BaseAttributes with the
    attributes passed to the given logging method.

    For example:

        type Counters struct{
            Error1: int `json:"error1"`,
            Error2: int `json:"error2"`,
        }
        counters := Counters{}
        options := logging.LoggerOptions{
          BaseAttributes: []any{"counters", &counters},
          BaseTags: []string{"foo", "bar"},
        }
        logger := logging.New(os.Stdout, &options)
        logger.Always(
            nil,
            logging.STACKTRACE, nil,
            logging.TAGS, []string{"hoo"},
            "baz", "waka")

    will include a one-line trace of the call stack for the currently executing
    function as the value of the log entry's "stacktrace" attribute in addition
    to having a "verbosity" attribute whose value is "ALWAYS" instead of a
    "level" attribute whose value is "ERROR". It will have a "counter" attribute
    whose values is the JSON representation of the current value of the counters
    variable each time a log entry is written, a "baz" attribute whose value
    is "waka" and a "tags" attribute whose value is a JSON array containing the
    strings "foo", "bar" and "hoo".

    Note that if LoggerOptions.ReplaceAttr is not nil, it will be called
    as described by the documentation for slog.HandlerOptions.ReplaceAttr
    indirectly through a custom replacer function that replaces "level" with
    "verbosity" as just described.

    Note also that the values for attributes in LoggerOptions.BaseAttributes may
    be passed by value or reference. Passing by reference allows for cases where
    each log entry should include the current value for that attribute rather
    than a copy of the value at the time the Logger was created.

func (l *Logger) Always(ctx context.Context, message MessageBuilder, attributes ...any)

func (l *Logger) BaseAttributes() []any

func (l *Logger) BaseTags() []string

func (l *Logger) Enabled(ctx context.Context, verbosity Verbosity) bool

func (l *Logger) Fine(ctx context.Context, message MessageBuilder, attributes ...any)

func (l *Logger) OnPanic(ctx context.Context, handler PanicHandler, attributes ...any)

func (l *Logger) Optional(ctx context.Context, message MessageBuilder, attributes ...any)

func (l *Logger) SetBaseAttributes(attributes ...any)

func (l *Logger) SetBaseTags(tags ...string)

func (l *Logger) SetVerbosity(verbosity Verbosity)

func (l *Logger) Trace(ctx context.Context, message MessageBuilder, attributes ...any)

func (l *Logger) Verbosity() Verbosity

type LoggerOptions struct {

	// Pass through to HandlerOptions for the wrapped slog.Logger.
	AddSource bool

	// Initial set of attributes that will be added to every log entry.
	BaseAttributes []any

	// Initial set of tags that will be added to every log entry.
	BaseTags []string

	// Shared slog.LevelVar, if desired; a Leveler will be created if this
	// is nil.
	Level *slog.LevelVar

	// If not nil, an attribute replacer function that will be called in
	// addition to replacing "level" attributes with "verbosiy" and other
	// special attribute handling.
	ReplaceAttr func([]string, slog.Attr) slog.Attr
}
    Configuration parameters for an instance of Logger.

type MessageBuilder func() string
    Type of function passed to logging methods for lazy evaluation of message
    formatting.

    The returned string becomes the value of the log entry's msg attribute.

    Such a function is invoked only if a given verbosity is enabled for a given
    logger.

type PanicHandler func(r any) (string, any)
    Type of function passed to Logger.OnPanic()

      - r: value returned from recover()
      - first return value: message to log
      - second return value: r or nil

    Logger.OnPanic() will re-invocke panic() if the second return value is not
    nil.

type Verbosity int
    Verbosity-based nomenclature used in place of slog.Level.
```
