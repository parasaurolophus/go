Copyright &copy; Kirk Rader 2024

# Stacktraces Test Helpers

Output of `go doc -all`:

```
package stacktraces_test // import "parasaurolophus/go/stacktraces_test"


FUNCTIONS

func FirstFunctionLong(stackTrace string) (string, int, error)
    Return the function name and frame number from the first entry in the given
    output of stacktraces.LongStackTrace(any).

func FirstFunctionShort(strackTrace string) (string, int, error)
    Return the function name and frame number from the first entry in the given
    output of stacktraces.ShortStackTrace(any).
```
