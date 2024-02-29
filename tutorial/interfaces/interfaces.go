// Copyright Kirk Rader 2024

package interfaces

import (
	"fmt"
	"strconv"
	"unicode"
)

type (

	// Declare a simple interface.
	Counter interface {
		Value() int
		Increment() int
		Decrement() int
	}

	// Define a private type that will implement Counter.
	basicCounter int
)

// Implement Counter.Value() for basicCounter.
func (counter basicCounter) Value() int {

	return int(counter)
}

// Implement Counter.Increment() for basicCounter.
func (counter *basicCounter) Increment() int {

	*counter += 1
	return int(*counter)
}

// Implement Counter.Decrement() for basicCounter.
func (counter *basicCounter) Decrement() int {

	*counter -= 1
	return int(*counter)
}

// Implement fmt.Stringer.String() for basicCounter.
func (counter basicCounter) String() string {

	return strconv.Itoa(int(counter))
}

// Implement fmt.Scanner.Scan() for basicCounter.
func (counter *basicCounter) Scan(state fmt.ScanState, _ rune) error {

	b, err := state.Token(true, unicode.IsDigit)

	if err != nil {
		return err
	}

	n, err := strconv.Atoi(string(b))

	if err != nil {
		return err
	}

	*counter = basicCounter(n)
	return nil
}

// Construct a new basicCounter.
func BasicCounter() Counter {

	counter := basicCounter(0)

	// Note the use of & here even though the return type is not a pointer. This
	// is because basicCounter has methods with pointer receivers. If you leave
	// off the & here, the compiler will signal a syntax error with a rather
	// cryptic message referring to pointer receivers and no hint that this
	// rather confusing idiom is the fix.
	return &counter
}
