// Copyright Kirk Rader 2024

package main

import (
	"fmt"
)

type (

	// Type of closure used in a demonstration of basic functional programming
	// principles.
	//
	// A Transformer accepts any type of paramter and returns the same or a
	// different value. It returns nil as a second result or else an error
	// object providing diagnostic information if a transformation fails.
	Transformer func(input any) (any, error)
)

// Return a Transformer which returns the result of adding n to its input, when
// passed an int.
//
// This demonstrates that a closure provides an encapsulation of the state of
// lexical variables that are in scope when it is created.
func IntAdder(n int) Transformer {

	return func(input any) (any, error) {

		switch v := input.(type) {

		case int:
			return v + n, nil

		default:
			return input, fmt.Errorf("unsupported value %v", v)
		}
	}
}

// Return a Transformer that returns the result of invoking the set of
// Transformers passed to it, in succession.
//
// This demonstrates the use of functional composition to share implementations
// in the absence of inheritance between instances of classes.
func CompositeTransformer(transformers ...Transformer) Transformer {

	return func(input any) (any, error) {

		var err error
		output := input

		for _, transformer := range transformers {

			output, err = transformer(output)

			if err != nil {
				return input, err
			}
		}

		return output, nil
	}
}

func main() {

	failed := 0
	oneAdder := IntAdder(1)
	twoAdder := IntAdder(2)
	composite := CompositeTransformer(twoAdder, oneAdder, twoAdder)

	output, err := composite(0)

	if err != nil {

		fmt.Println(err.Error())
		failed += 1
	}

	switch v := output.(type) {

	case int:
		if v != 5 {
			fmt.Printf("expected 5, got %d\n", v)
			failed += 1
		}

	default:
		fmt.Printf("expected 5, got %v\n", v)
		failed += 1
	}

	if failed == 0 {
		fmt.Println("ok")
	} else {
		fmt.Printf("%d failed\n", failed)
	}
}
