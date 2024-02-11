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
{"time":"2024-02-10T16:20:06.729367361-06:00","verbosity":"TRACE","msg":"main.main starting goroutine","tags":["EXAMPLE"]}
{"time":"2024-02-11T06:04:52.587630266-06:00","verbosity":"TRACE","msg":"main.main starting goroutine","tags":["EXAMPLE"]}
{"time":"2024-02-11T06:04:52.58806089-06:00","verbosity":"TRACE","msg":"main.main consuming output from goroutine","tags":["EXAMPLE"]}
0
1
2
3
{"time":"2024-02-11T06:04:52.588156667-06:00","verbosity":"TRACE","msg":"sender deliberately causing a panic","tags":["EXAMPLE"]}
4
{"time":"2024-02-11T06:04:52.58829935-06:00","verbosity":"ALWAYS","msg":"main.sender recovered from 'deliberate panic'","stacktrace":"7:main.sender [/source/go/example/example.go:103] < 8:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]","tags":["EXAMPLE","PANIC","ERROR","SEVERE"]}
{"time":"2024-02-11T06:04:52.588404108-06:00","verbosity":"TRACE","msg":"main.sender closing channel","tags":["EXAMPLE"]}
{"time":"2024-02-11T06:04:52.588465533-06:00","verbosity":"TRACE","msg":"main.main exiting normally","tags":["EXAMPLE"]}
```
