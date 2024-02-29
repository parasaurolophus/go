// Copyright Kirk Rader 2024

package main

import (
	"fmt"
	"parasaurolophus/go/tutorial/interfaces"
)

func printAny(value any) {

	t := "unknown type"

	switch v := value.(type) {

	case int:
		t = "int"

	case string:
		t = "string"

	case complex128:
		t = "complex128"

	default:
		t = fmt.Sprintf("%v is of an unsupported type", v)
	}

	fmt.Printf("%30s: %v\n", t, value)
}

// prints
//
//	int: 42
//	string: forty-two
//	complex128: (4+2i)
//	0 is of an unsupported type: 0
//
// to stdout
func main() {

	printAny(42)
	printAny("forty-two")
	printAny(4 + 2i)
	printAny(interfaces.BasicCounter())
}
