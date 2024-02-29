// Copyright Kirk Rader 2024

package generics

type (

	// Type constraint for any type of signed integer.
	Integer interface {
		int | int8 | int16 | int32 | int64
	}

	// Type constraint for any type of unsigned integer.
	Unsigned interface {
		uint | uint8 | uint16 | uint32 | uint64
	}

	// Type constraint for either type of floating-point number.
	Float interface {
		float32 | float64
	}

	// Type constraint for either type of imaginary number.
	Complex interface {
		complex64 | complex128
	}

	// Type constraint for any type of number.
	Number interface {
		byte | rune | Integer | Unsigned | Float | Complex
	}
)

// Generic function for summing any type that matches the Number constraint.
func Sum[N Number](numbers ...N) N {

	var n N

	for _, operand := range numbers {
		n += operand
	}

	return n
}
