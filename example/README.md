Copyright &copy; Kirk Rader 2024

# Go Examples

[example.go](./example.go) contains examples of how to use the utility libraries in this repository.

## Usage

```bash
go run example.go
```

writes the following to `stdout` (the panic is deliberate and demonstrates the
use of `logging.Defer()` to log and optionally recover):

```
{"time":"2024-02-15T04:06:09.32851209-06:00","verbosity":"TRACE","msg":"main.main starting goroutine","sent":0,"tags":["EXAMPLE"]}
{"time":"2024-02-15T04:06:09.32906717-06:00","verbosity":"TRACE","msg":"main.main consuming output from goroutine","sent":0,"tags":["EXAMPLE"]}
0
1
2
3
{"time":"2024-02-15T04:06:09.329214037-06:00","verbosity":"TRACE","msg":"sender deliberately causing a panic","sent":5,"tags":["EXAMPLE"]}
4
{"time":"2024-02-15T04:06:09.3293842-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from 'deliberate panic'","sent":5,"recovered":"deliberate panic","stacktrace":"7:main.sender [/source/go/example/example.go:112] < 8:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1222]","tags":["EXAMPLE","PANIC","ERROR","SEVERE"]}
{"time":"2024-02-15T04:06:09.329432513-06:00","verbosity":"TRACE","msg":"main.sender closing channel","sent":5,"tags":["EXAMPLE"]}
{"time":"2024-02-15T04:06:09.329465938-06:00","verbosity":"TRACE","msg":"main.main exiting normally","sent":5,"tags":["EXAMPLE"]}
```
