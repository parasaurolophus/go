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

verbosity: ALWAYS, panicInMain: false, panicAgain: true

{"time":"2024-02-17T11:22:01.362040903-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:96] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

$ go run example.go OPTIONAL

verbosity: OPTIONAL, panicInMain: false, panicAgain: true

{"time":"2024-02-17T11:22:08.034775245-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:96] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

{"time":"2024-02-17T11:22:08.035249923-06:00","verbosity":"OPTIONAL","msg":"","sent":5,"received":5,"tags":["EXAMPLE","DEBUG"]}
$ go run example.go FINE    

verbosity: FINE, panicInMain: false, panicAgain: true

{"time":"2024-02-17T11:22:13.49688685-06:00","verbosity":"FINE","msg":"0","sent":2,"received":1,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:13.497264492-06:00","verbosity":"FINE","msg":"1","sent":2,"received":2,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:13.497312214-06:00","verbosity":"FINE","msg":"2","sent":4,"received":3,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:13.497353139-06:00","verbosity":"FINE","msg":"3","sent":4,"received":4,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:13.497453119-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":4,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:96] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}
{"time":"2024-02-17T11:22:13.497549358-06:00","verbosity":"FINE","msg":"4","sent":5,"received":5,"tags":["EXAMPLE"]}

5 sent, 5 received

{"time":"2024-02-17T11:22:13.497639764-06:00","verbosity":"OPTIONAL","msg":"","sent":5,"received":5,"tags":["EXAMPLE","DEBUG"]}
$ go run example.go TRACE

verbosity: TRACE, panicInMain: false, panicAgain: true

{"time":"2024-02-17T11:22:19.230981987-06:00","verbosity":"TRACE","msg":"main.main starting sender goroutine","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-17T11:22:19.231428425-06:00","verbosity":"TRACE","msg":"main.main consuming output from sender goroutine","sent":0,"received":0,"tags":["EXAMPLE","DEBUG"]}
{"time":"2024-02-17T11:22:19.231510738-06:00","verbosity":"FINE","msg":"0","sent":2,"received":1,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:19.231542534-06:00","verbosity":"FINE","msg":"1","sent":2,"received":2,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:19.231568089-06:00","verbosity":"FINE","msg":"2","sent":4,"received":3,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:19.231591144-06:00","verbosity":"FINE","msg":"3","sent":4,"received":4,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:19.231639292-06:00","verbosity":"TRACE","msg":"sender deliberately causing a panic","sent":5,"received":5,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:19.231679384-06:00","verbosity":"FINE","msg":"4","sent":5,"received":5,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:19.231791289-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":5,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:96] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}
{"time":"2024-02-17T11:22:19.231826493-06:00","verbosity":"TRACE","msg":"main.sender goroutine closing channel","sent":5,"received":5,"tags":["EXAMPLE"]}

5 sent, 5 received

{"time":"2024-02-17T11:22:19.231867862-06:00","verbosity":"TRACE","msg":"main.main exiting normally","sent":5,"received":5,"tags":["EXAMPLE"]}
{"time":"2024-02-17T11:22:19.231892899-06:00","verbosity":"OPTIONAL","msg":"","sent":5,"received":5,"tags":["EXAMPLE","DEBUG"]}
$ go run example.go ALWAYS true

verbosity: ALWAYS, panicInMain: true, panicAgain: true

{"time":"2024-02-17T11:22:53.308865816-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from: deliberate panic","sent":5,"received":4,"recovered":"deliberate panic","stacktrace":"8:main.sender [/source/go/example/example.go:96] < 9:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","MEDIUM"]}

5 sent, 5 received

{"time":"2024-02-17T11:22:53.309389623-06:00","verbosity":"ALWAYS","msg":"main.main panicing: another deliberate panic","sent":5,"received":5,"recovered":"another deliberate panic","stacktrace":"8:main.main [/source/go/example/example.go:70] < 9:runtime.main [/usr/local/go/src/runtime/proc.go:271] < 10:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","SEVERE"]}
panic: another deliberate panic [recovered]
	panic: another deliberate panic

goroutine 1 [running]:
parasaurolophus/go/logging.(*Logger).finallyCommon(0x400006a180, 0x1, {0x12fa78, 0x235800}, 0x4000074c18?, 0x4000074ec0, {0x4000074c38, 0x4, 0x4})
	/source/go/logging/logger.go:338 +0x1f8
parasaurolophus/go/logging.(*Logger).Finally(...)
	/source/go/logging/logger.go:148
panic({0xe9800?, 0x12ee18?})
	/usr/local/go/src/runtime/panic.go:770 +0x124
main.main()
	/source/go/example/example.go:70 +0x59c
exit status 2
```