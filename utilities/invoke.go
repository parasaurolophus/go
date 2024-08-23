// Copyright 2024 Kirk Rader

package utilities

import (
	"fmt"
	"io"
	"runtime"
)

// Trap and log panics.
func Invoke[T any](handler func(T), value T, log io.Writer) {
	defer func() {
		if r := recover(); r != nil {
			if log == nil {
				return
			}
			b := make([]byte, 1024)
			n := runtime.Stack(b, false)
			s := string(b[:n])
			fmt.Fprintf(log, "panic: %v\n%s\n", r, s)
		}
	}()
	handler(value)
}
