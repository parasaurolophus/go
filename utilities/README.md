# parasaurolophus/go/utilities

Output of `go doc -all`:

```
package utilities // import "parasaurolophus/go/utilities"


FUNCTIONS

func Async[T any](asyncFunction func(T) T, in chan T, out chan T, panicHandler func(recovered any))
    Invoke asyncFunction for each value received on in, sending the result on
    out. If asyncFunction panics, invoke panicHandler and then send T's zero
    value to out.

func Fetch(url string) (readCloser io.ReadCloser, err error)
    Fetch a document from the given URL.

func FilterMapKeys[K comparable, V any](m map[K]V, keepKeys ...K) map[K]V
    Return a shallow copy of m, with only the keys specified by keepKeys.

func MakeJSONNumberTokenTest() func(rune) bool
    Return a closure that can be used with fmt.ScanState.Token to convert
    the text representation of a number to a float64 according to JSON number
    syntax.

func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V
    Return a map that combines the key / value pairs from all the given ones.
    The maps are processed in the given order. If the same key appears more than
    once, the value in the result will be the last one from the parameter list.
```
