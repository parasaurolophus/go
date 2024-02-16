_Copyright &copy; Kirk Rader 2024_

# Wrapper for `log/slog`

Note: the primary source of documentation for this library is in the form of
inline comments that can be viewed using `go doc`.

```
$ go doc -all
package logging // import "parasaurolophus/go/logging"


CONSTANTS

const (

	// Logger.Defer() and Logger.DeferContext() will include the value returned
	// by recover() when logging a panic.
	RECOVERED = "recovered"

	// Values of "stacktrace" attributes will be replaced with one-line stack
	// traces for the function that called the given logging method.
	STACKTRACE = "stacktrace"

	// Value will be merged with the currently configured
	// LoggerOptions.BaseTags.
	TAGS = "tags"
)
    Specially handled attributes.

const (

	// Only emit a log entry when extremely verbose output is specified.
	//
	// Intended for use in development environments for focused debugging
	// sessions. This should never be enabled outside of development
	// environments. Any logging that might potentially reveal PII, SPI or
	// critically sensitive security information must only be written at TRACE
	// level in environments where only synthetic or redacted data is in use.
	TRACE = Verbosity(slog.LevelDebug)

	// Only emit a log entry when unusually verbose output is specified.
	//
	// Intended for use in development environments for everyday testing and
	// troubleshooting prior to a release candidate being deployed.
	FINE = Verbosity(slog.LevelInfo)

	// Only emit a log entry when moderately verbose output is specified.
	//
	// Intended for use in testing and staging environments, e.g. during
	// acceptance and regression tests before release to production.
	OPTIONAL = Verbosity(slog.LevelWarn)

	// Always emit a log entry.
	//
	// Intended for production environments to drive monitoring, alerting and
	// reporting.
	ALWAYS = Verbosity(slog.LevelError)
)
    Mapping of Verbosity to slog.Level values.

    Generally, assume that only ALWAYS will be enabled in production
    environments and that TRACE will never be enabled outside of development
    environments.


TYPES

type Finally func()
    Type of function passed as first argument to Logger.Defer() and
    Logger.DeferContext().

type Logger struct {
	// Has unexported fields.
}
    Wrapper for an instance of slog.Logger.

func New(writer io.Writer, options *LoggerOptions) *Logger
    Returns a newly created, wrapped instance of slog.Logger.

    Log entries written using the returned Logger instance will have
    "verbosity" attributes instead of "level" attributes and the values of
    their "stacktrace" attributes, if present, will be replaced as if by an
    invocation of ShortStackTrace(skipFrames) where skipFrames is the value
    of the "stacktrace" attribute passed to a logging method. The final set of
    attributes for each log entry will be the result of combining the value of
    LoggerOptions.BaseAttributes and LoggerOptions.BaseTags with the attributes
    passed to the given logging method.

    For example:

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

    produces a log entry like:

        {"time":"2024-02-11T06:16:41.852302853-06:00","verbosity":"OPTIONAL","msg":"n = 42","counters":{"error1":1,"error2":0},"baz":"waka","stacktrace":"5:main.main [/source/go/scratch/scratch.go:29] < 6:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 7:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]","tags":["foo","bar","hoo"]}

    Note that if LoggerOptions.ReplaceAttr is not nil, it will be called
    as described by the documentation for slog.HandlerOptions.ReplaceAttr
    indirectly through a custom replacer function that replaces "level" with
    "verbosity" as just described.

    Note also that the values for attributes in LoggerOptions.BaseAttributes may
    be passed by value or reference. Passing by reference allows for cases where
    each log entry should include the current value for that attribute rather
    than a copy of the value at the time the Logger was created.

func (l *Logger) Always(message MessageBuilder, attributes ...any)
    Log at ALWAYS verbosity.

func (l *Logger) AlwaysContext(ctx context.Context, message MessageBuilder, attributes ...any)
    Log at ALWAYS verbosity using the supplied context.

func (l *Logger) BaseAttributes() []any
    Return the current base attributes.

func (l *Logger) BaseTags() []string
    Return the current base tags.

func (l *Logger) Defer(panicAgain bool, finally Finally, recoverHandler RecoverHandler, attributes ...any)
    See documentation for Logger.DeferContext().

func (l *Logger) DeferContext(panicAgain bool, finally Finally, ctx context.Context, recoverHandler RecoverHandler, attributes ...any)
    For use with defer to log if a panic occurs.

    If recover() returns non-nil, its value will be passed to handler.

    Handler's return value will be used as the msg string in writing a log entry
    using l.AlwaysContext().

    If panicAgain is true, any panics that occur while this deferred method is
    in effect will be passed to panic() so as to cause the process to terminate
    abnormally.

    For example, if the following is invoked in a goroutine that was passed a
    channel named ch:

        name := stacktraces.FunctionName()
        defer logger.DeferContext(

            // don't cause process to exit abnormally even if a panic occurs
            false,

            // clean-up function is always invoked
            func() { close(ch) },

            // remaining parameters are passed to logger.AlwaysContext() when
            // recover() returns non-nil

            ctx,
            func(r any) (string, any) {
                // second value will be used to resume panicing if non-nil
                // (typically this would be r to continue the now tidied
                // and logged panic in main.main or nil in a goroutine
                // so as to allow other goroutines to complete)
                return fmt.Sprintf("%s recovered from %v", name, r), nil
            },
        )

    the goroutine will close ch on exit and, if a panic occurs, write a log
    entry whose msg is the string representation of the value returned by
    recover()while allowing other goroutines to continue running normally.
    If panicAgain were passed true, recovered value would be passed to panic()
    after the clean up and logging functions were invoked. The value of
    panicAgain is also used to determine whether or not panics in the clean-up
    or message builder functions cause an abnormal exit. [See the documentation
    for panic() and recover() for more information.]

func (l *Logger) Enabled(verbosity Verbosity) bool
    Return true or false depending on whether or not the given verbosity is
    currently enabled for the given logger.

func (l *Logger) EnabledContext(ctx context.Context, verbosity Verbosity) bool
    Return true or false depending on whether or not the given verbosity is
    currently enabled for the given logger.

func (l *Logger) Fine(message MessageBuilder, attributes ...any)
    Log at FINE verbosity.

func (l *Logger) FineContext(ctx context.Context, message MessageBuilder, attributes ...any)
    Log at FINE verbosity using the supplied context.

func (l *Logger) Optional(message MessageBuilder, attributes ...any)
    Log at OPTIONAL verbosity.

func (l *Logger) OptionalContext(ctx context.Context, message MessageBuilder, attributes ...any)
    Log at OPTIONAL verbosity using the supplied context.

func (l *Logger) SetBaseAttributes(attributes ...any)
    Update the base attributes.

func (l *Logger) SetBaseTags(tags ...string)
    Update the base tags.

func (l *Logger) SetContext(ctx context.Context)
    Deprecated hack for backwards compatibility.

func (l *Logger) SetVerbosity(verbosity Verbosity)
    Update the verbosity

func (l *Logger) Trace(message MessageBuilder, attributes ...any)
    Log at TRACE verbosity.

func (l *Logger) TraceContext(ctx context.Context, message MessageBuilder, attributes ...any)
    Log at TRACE verbosity using the supplied context.

func (l *Logger) Verbosity() Verbosity
    Return the current verbosity

type LoggerOptions struct {

	// Allow panics in Logger.log() to cause abnormal process termination.
	AllowPanics bool

	// Initial set of attributes that will be added to every log entry.
	BaseAttributes []any

	// Initial set of tags that will be added to every log entry.
	BaseTags []string

	// Pass through to HandlerOptions for the wrapped slog.Logger.
	AddSource bool

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

type RecoverHandler func(recovered any) string
    Type of function passed to Logger.Defer() and Logger.DeferContext() to allow
    for including the value returned by recover() in the log entry.

type Verbosity int
    Verbosity-based nomenclature used in place of slog.Level.

func (v *Verbosity) Scan(state fmt.ScanState, verb rune) error

func (v Verbosity) String() string
    Implement fmt.Stringer interface for Verbosity.
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
slogger.Debug("...")
slogger.Info("...")
slogger.Warn("...")
slogger.Error("...")
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
