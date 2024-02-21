Copyright &copy; Kirk Rader 2024

# Go Examples

[example.go](./example.go) contains examples of how to use the utility libraries in this repository.

```bash
$ go doc -cmd -u -all
```

```
package main // import "parasaurolophus/go/example"


VARIABLES

var (

	// The number of values sent by a goroutine.
	sent = 0

	// The number of values received from the goroutine.
	received = 0

	// Logger configuration.
	loggerOptions = logging.LoggerOptions{
		BaseTags: []string{"EXAMPLE"},
		BaseAttributes: []any{
			"sent", &sent,
			"received", &received,
		},
	}

	// Logger.
	logger = logging.New(os.Stdout, &loggerOptions)
)

FUNCTIONS

func main()
    Print the number of values sent by and received from a goroutine to stdout.

    Logging verbosity defaults to OPTIONAL but may be set using a command-line
    argument.

func parseArg(index int, typeName string, val any)
func sender(ch chan int)
    Goroutine that sends int values to a channel.

    This deliberately panics after sending a few values as a demonstration of
    logging.Logger.Defer().
```

## Usage

### User Input Error

```bash
go run example.go DEBUG
```

```
{"time":"2024-02-21T04:05:28.942968146-06:00","verbosity":"OPTIONAL","msg":"unsupported verbosity token: 'DEBUG'","sent":0,"received":0,"stacktrace":"5:main.parseArg [/source/go/example/example.go:133] < 6:main.main [/source/go/example/example.go:44] < 7:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 8:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","ERROR","USER","BAD_ARGS"]}
exit status 1
```

### Least Verbose, One Recovered Panic

```bash
$ go run example.go ALWAYS false false
```

```
verbosity: ALWAYS, panicInMain: false, panicAgain: false

{"time":"2024-02-21T03:47:29.170238189-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:111] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received
```

### Least Verbose, Two Recovered Panics


```bash
$ go run example.go ALWAYS true false
```

```
verbosity: ALWAYS, panicInMain: true, panicAgain: false

{"time":"2024-02-21T03:47:40.480239625-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:111] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

{"time":"2024-02-21T03:47:40.480704452-06:00","verbosity":"ALWAYS","msg":"main.main panicing: another deliberate panic","sent":5,"received":5,"recovered":"another deliberate panic","stacktrace":"8:main.main [/source/go/example/example.go:82] < 9:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 10:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","SEVERE"]}
```

### Least Verbose, Two Panics, One Recovered


```bash
$ go run example.go ALWAYS true true 
```

```
verbosity: ALWAYS, panicInMain: true, panicAgain: true

{"time":"2024-02-21T03:47:53.98154743-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:111] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

{"time":"2024-02-21T03:47:53.981992812-06:00","verbosity":"ALWAYS","msg":"main.main panicing: another deliberate panic","sent":5,"received":5,"recovered":"another deliberate panic","stacktrace":"8:main.main [/source/go/example/example.go:82] < 9:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 10:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","SEVERE"]}

panic: another deliberate panic [recovered]
	panic: another deliberate panic

goroutine 1 [running]:
parasaurolophus/go/logging.(*Logger).finallyCommon(0x4000122120, 0x1, {0x12fad8, 0x235800}, 0x4000116c18?, 0x4000116ec0, {0x4000116c38, 0x4, 0x4})
	/source/go/logging/logger.go:369 +0x1f8
parasaurolophus/go/logging.(*Logger).Finally(...)
	/source/go/logging/logger.go:157
panic({0xe9800?, 0x12ee88?})
	/usr/local/go/src/runtime/panic.go:770 +0x124
main.main()
	/source/go/example/example.go:82 +0x59c
exit status 2
```

### Most Vebose, Two Recovered Panics

```bash
$ go run example.go TRACE true false
```

```
verbosity: TRACE, panicInMain: true, panicAgain: false

{"time":"2024-02-21T03:48:08.458892518-06:00","verbosity":"TRACE","msg":"main.main starting sender goroutine","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-21T03:48:08.459284031-06:00","verbosity":"TRACE","msg":"main.main consuming output from sender goroutine","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-21T03:48:08.459369178-06:00","verbosity":"FINE","msg":"0","sent":2,"received":1,"tags":["EXAMPLE"]}
{"time":"2024-02-21T03:48:08.459411677-06:00","verbosity":"FINE","msg":"1","sent":2,"received":2,"tags":["EXAMPLE"]}
{"time":"2024-02-21T03:48:08.459449251-06:00","verbosity":"FINE","msg":"2","sent":4,"received":3,"tags":["EXAMPLE"]}
{"time":"2024-02-21T03:48:08.459486528-06:00","verbosity":"FINE","msg":"3","sent":4,"received":4,"tags":["EXAMPLE"]}
{"time":"2024-02-21T03:48:08.459522787-06:00","verbosity":"TRACE","msg":"sender deliberately causing a panic","sent":5,"received":4,"tags":["EXAMPLE"]}
{"time":"2024-02-21T03:48:08.459672007-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:111] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}
{"time":"2024-02-21T03:48:08.45972571-06:00","verbosity":"TRACE","msg":"main.sender goroutine closing channel","sent":5,"received":5,"tags":["EXAMPLE"]}
{"time":"2024-02-21T03:48:08.459709154-06:00","verbosity":"FINE","msg":"4","sent":5,"received":5,"tags":["EXAMPLE"]}

5 sent, 5 received

{"time":"2024-02-21T03:48:08.459884115-06:00","verbosity":"ALWAYS","msg":"main.main panicing: another deliberate panic","sent":5,"received":5,"recovered":"another deliberate panic","stacktrace":"8:main.main [/source/go/example/example.go:82] < 9:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 10:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","SEVERE"]}
{"time":"2024-02-21T03:48:08.459916448-06:00","verbosity":"OPTIONAL","msg":"","sent":5,"received":5,"tags":["EXAMPLE","DEBUG"]}
```
