Copyright &copy; Kirk Rader 2024

# Go Examples

Examples of how to use the utility libraries in this repository.

## Usage

```bash
cd example
go run example.go
```

writes the following to `stdout` (the panic is deliberate and demonstrates the
use `logging.OnPanic()`):

```
FunctionName(): main.main

ShortStackTrace("runtime.main"): 4:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]

ShortStackTrace(-1): 3:main.main [/source/go/example/example.go:33] < 4:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]

ShortStackTrace(0): 0:runtime.Callers [/usr/local/go/src/runtime/extern.go:308] < 1:parasaurolophus/go/stacktraces.formatStackTrace [/source/go/stacktraces/impl.go:127] < 2:parasaurolophus/go/stacktraces.ShortStackTrace [/source/go/stacktraces/stacktraces.go:149] < 3:main.main [/source/go/example/example.go:37] < 4:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 5:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]

stackTrace.Error(): StackTrace as error
stackTrace.LongTrace():
3:main.main.func1
/source/go/example/example.go:41
0xbf86b
---
4:main.main
/source/go/example/example.go:42
0xbf848
---
5:runtime.main
/usr/local/go/src/runtime/proc.go:267
0x442cb
---
6:runtime.goexit
/usr/local/go/src/runtime/asm_arm64.s:1197
0x6f693

you will see this
{"time":"2024-02-09T06:32:06.982153871-06:00","verbosity":"TRACE","msg":"n is 42","counters":{"error1":0,"error2":0},"stacktrace":"5:main.main [/source/go/example/example.go:86] < 6:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 7:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]","tags":["example"]}

{"time":"2024-02-09T06:32:06.982568276-06:00","verbosity":"ALWAYS","msg":"note that counters in the log entry reflects the current value of error1","counters":{"error1":1,"error2":0},"tags":["example","reference_test"]}

deliberately panicing
{"time":"2024-02-09T06:32:06.982660942-06:00","verbosity":"ALWAYS","msg":"panic: \"example\"","counters":{"error1":1,"error2":0},"stacktrace":"5:parasaurolophus/go/logging.(*Logger).OnPanic [/source/go/logging/logging.go:182] < 6:runtime.gopanic [/usr/local/go/src/runtime/panic.go:914] < 7:main.main.func6 [/source/go/example/example.go:110] < 8:main.main [/source/go/example/example.go:113] < 9:runtime.main [/usr/local/go/src/runtime/proc.go:267] < 10:runtime.goexit [/usr/local/go/src/runtime/asm_arm64.s:1197]","tags":["example","ERROR","PANIC","SEVERE"]}
panic: example [recovered]
	panic: example

goroutine 1 [running]:
parasaurolophus/go/logging.(*Logger).OnPanic(0x4000147ad8?, {0x108258, 0x1d2fc0}, 0xeb030, {0x4000147af8, 0x4, 0x4})
	/source/go/logging/logging.go:185 +0xb8
panic({0xc8bc0?, 0x107988?})
	/usr/local/go/src/runtime/panic.go:914 +0x218
main.main.func6(...)
	/source/go/example/example.go:110
main.main()
	/source/go/example/example.go:113 +0x634
exit status 2
```

