// Copyright Kirk Rader 2024

package utilities

// Invoke asyncFunction for each value received on in, sending the result on
// out. If asyncFunction panics, invoke panicHandler and then send T's zero
// value to out.
func Async[T any](asyncFunction func(T) T, in chan T, out chan T, panicHandler func(recovered any)) {
	defer close(out)
	for value := range in {
		result := func(x T) (y T) {
			defer func() {
				if r := recover(); r != nil {
					if panicHandler != nil {
						panicHandler(r)
					}
				}
			}()
			y = asyncFunction(x)
			return
		}(value)
		out <- result
	}
}