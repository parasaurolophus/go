Copyright &copy; Kirk Rader 2024

# Go Examples

[example.go](./example.go) contains examples of how to use the utility libraries in this repository.

```
$ go doc -cmd -all -u
package main // import "parasaurolophus/go/example"


VARIABLES

var (
	// The number of values sent by a goroutine.
	sent = 0
	// The number of values received from the goroutine.
	received = 0

	loggerOptions = logging.LoggerOptions{
		BaseTags: []string{"EXAMPLE"},
		BaseAttributes: []any{
			"sent", &sent,
			"received", &received,
		},
	}

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

```
$ go run example.go ALWAYS
{"time":"2024-02-16T05:26:01.903934673-06:00","verbosity":"FINE","msg":"optional argument 2 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-16T05:26:01.904361262-06:00","verbosity":"FINE","msg":"optional argument 3 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}

verbosity: ALWAYS, panicInMain: false, panicAgain: true

{"time":"2024-02-16T05:26:01.904606093-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:95] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

$ go run example.go ALWAYS true false

verbosity: ALWAYS, panicInMain: true, panicAgain: false

{"time":"2024-02-16T05:26:24.54384461-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:95] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

{"time":"2024-02-16T05:26:24.544306773-06:00","verbosity":"ALWAYS","msg":"main.main panicing: another deliberate panic","sent":5,"received":5,"recovered":"another deliberate panic","stacktrace":"8:main.main [/source/go/example/example.go:69] < 9:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 10:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","SEVERE"]}
$ go run example.go ALWAYS true      
{"time":"2024-02-16T05:26:31.015966465-06:00","verbosity":"FINE","msg":"optional argument 3 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}

verbosity: ALWAYS, panicInMain: true, panicAgain: true

{"time":"2024-02-16T05:26:31.016531237-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:95] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

{"time":"2024-02-16T05:26:31.016664403-06:00","verbosity":"ALWAYS","msg":"main.main panicing: another deliberate panic","sent":5,"received":5,"recovered":"another deliberate panic","stacktrace":"8:main.main [/source/go/example/example.go:69] < 9:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 10:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","SEVERE"]}
panic: another deliberate panic [recovered]
	panic: another deliberate panic

goroutine 1 [running]:
parasaurolophus/go/logging.(*Logger).finally(0x40000a8120, 0x1, {0x12fa78, 0x235800}, 0x400008fc18?, 0x400008fec0, {0x400008fc38, 0x4, 0x4})
	/source/go/logging/logging.go:241 +0x1c0
parasaurolophus/go/logging.(*Logger).Defer(...)
	/source/go/logging/logging.go:176
panic({0xe97e0?, 0x12ee20?})
	/usr/local/go/src/runtime/panic.go:770 +0x124
main.main()
	/source/go/example/example.go:69 +0x584
exit status 2
$ go run example.go OPTIONAL   
{"time":"2024-02-16T05:26:47.642566119-06:00","verbosity":"FINE","msg":"optional argument 2 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-16T05:26:47.642953375-06:00","verbosity":"FINE","msg":"optional argument 3 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}

verbosity: OPTIONAL, panicInMain: false, panicAgain: true

{"time":"2024-02-16T05:26:47.643159854-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":4,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:95] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

{"time":"2024-02-16T05:26:47.643232224-06:00","verbosity":"OPTIONAL","msg":"","sent":5,"received":5,"tags":["EXAMPLE","DEBUG"]}
$ go run example.go FINE    
{"time":"2024-02-16T05:26:54.307132506-06:00","verbosity":"FINE","msg":"optional argument 2 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-16T05:26:54.307515058-06:00","verbosity":"FINE","msg":"optional argument 3 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}

verbosity: FINE, panicInMain: false, panicAgain: true

{"time":"2024-02-16T05:26:54.307637427-06:00","verbosity":"FINE","msg":"0","sent":2,"received":1,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:26:54.307673667-06:00","verbosity":"FINE","msg":"1","sent":2,"received":2,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:26:54.307700223-06:00","verbosity":"FINE","msg":"2","sent":4,"received":3,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:26:54.307726074-06:00","verbosity":"FINE","msg":"3","sent":4,"received":4,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:26:54.307829499-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:95] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}
{"time":"2024-02-16T05:26:54.307909702-06:00","verbosity":"FINE","msg":"4","sent":5,"received":5,"tags":["EXAMPLE"]}

5 sent, 5 received

{"time":"2024-02-16T05:26:54.308039294-06:00","verbosity":"OPTIONAL","msg":"","sent":5,"received":5,"tags":["EXAMPLE","DEBUG"]}
$ go run example.go TRACE
{"time":"2024-02-16T05:27:03.285020249-06:00","verbosity":"FINE","msg":"optional argument 2 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-16T05:27:03.285590355-06:00","verbosity":"FINE","msg":"optional argument 3 (of type bool) not supplied","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}

verbosity: TRACE, panicInMain: false, panicAgain: true

{"time":"2024-02-16T05:27:03.285684502-06:00","verbosity":"TRACE","msg":"main.main starting sender goroutine","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-16T05:27:03.285735743-06:00","verbosity":"TRACE","msg":"main.main consuming output from sender goroutine","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-16T05:27:03.285800279-06:00","verbosity":"FINE","msg":"0","sent":2,"received":1,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:27:03.285829871-06:00","verbosity":"FINE","msg":"1","sent":2,"received":2,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:27:03.285856167-06:00","verbosity":"FINE","msg":"2","sent":4,"received":3,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:27:03.285878778-06:00","verbosity":"FINE","msg":"3","sent":4,"received":4,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:27:03.285901056-06:00","verbosity":"FINE","msg":"4","sent":5,"received":5,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:27:03.285924611-06:00","verbosity":"TRACE","msg":"sender deliberately causing a panic","sent":5,"received":5,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:27:03.286045555-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:95] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}
{"time":"2024-02-16T05:27:03.286077517-06:00","verbosity":"TRACE","msg":"main.sender goroutine closing channel","sent":5,"received":5,"tags":["EXAMPLE"]}

5 sent, 5 received

{"time":"2024-02-16T05:27:03.28612472-06:00","verbosity":"TRACE","msg":"main.main exiting normally","sent":5,"received":5,"tags":["EXAMPLE"]}
{"time":"2024-02-16T05:27:03.286149109-06:00","verbosity":"OPTIONAL","msg":"","sent":5,"received":5,"tags":["EXAMPLE","DEBUG"]}
$ go run example.go DEBUG
{"time":"2024-02-16T05:27:16.057378246-06:00","verbosity":"OPTIONAL","msg":"unsupported verbosity token: 'DEBUG'","sent":0,"received":0,"stacktrace":"5:main.parseArg [/source/go/example/example.go:111] < 6:main.main [/source/go/example/example.go:38] < 7:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 8:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","ERROR","USER","BAD_ARGS"]}
exit status 1
```