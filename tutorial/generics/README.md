_Copyright &copy; Kirk Rader 2024_

# Generic Functions in Go

Go has no mechanism for defining parameterized types. It does provide a
mechanism for defining generic functions, which is a partial work-around for its
lack of support for function overloading. The `generics` package in this
directory demonstrates a simple example of a generic `Sum()` function that
operates on any numeric type. In particular:

- Use type constraints to declare requirements for parameters to a generic
  function.

- Declare a generic function whose parameters thus guaranteed to support the
  required operations.

```
$ go doc -all
package generics // import "parasaurolophus/tutorial/generics"


FUNCTIONS

func Sum[N Number](numbers ...N) N
    Generic function for summing any type that matches the Number constraint.


TYPES

type Complex interface {
	complex64 | complex128
}
    Type constraint for either type of imaginary number.

type Float interface {
	float32 | float64
}
    Type constraint for either type of floating-point number.

type Integer interface {
	int | int8 | int16 | int32 | int64
}
    Type constraint for any type of signed integer.

type Number interface {
	byte | rune | Integer | Unsigned | Float | Complex
}
    Type constraint for any type of number.

type Unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}
    Type constraint for any type of unsigned integer.
```
