// Copyright Kirk Rader 2024

package enums

import (
	"fmt"
	"unicode"
)

type (

	// Define a public type alias for int.
	Enum1 int

	// Define another public type alias for int.
	Enum2 int
)

// Define public named constants of type Enum1.
const (
	One Enum1 = iota
	Two
	Three
)

// Define public named constants of type Enum2.
const (
	Four Enum2 = iota
	Five
	Six
)

// Implement the fmt.Stringer interface for Enum1.
func (e1 Enum1) String() string {

	switch e1 {

	case One:
		return one

	case Two:
		return two

	case Three:
		return three

	default:
		return ""
	}
}

// Implement the fmt.Stringer interface for Enum1.
func (e2 Enum2) String() string {

	switch e2 {

	case Four:
		return four

	case Five:
		return five

	case Six:
		return six

	default:
		return ""
	}
}

// Implement the fmt.Scanner interface for Enum1.
func (e1 *Enum1) Scan(state fmt.ScanState, _ rune) error {

	b, err := state.Token(true, unicode.IsLetter)

	if err != nil {
		return err
	}

	s := string(b)

	switch s {

	case one:
		*e1 = One

	case two:
		*e1 = Two

	case three:
		*e1 = Three

	default:
		return fmt.Errorf("not an Enum1: \"%s\"", s)
	}

	return nil
}

// Implement the fmt.Scanner interface for Enum2.
func (e2 *Enum2) Scan(state fmt.ScanState, _ rune) error {

	b, err := state.Token(true, unicode.IsLetter)

	if err != nil {
		return err
	}

	s := string(b)

	switch s {

	case four:
		*e2 = Four

	case five:
		*e2 = Five

	case six:
		*e2 = Six

	default:
		return fmt.Errorf("not an Enum2: \"%s\"", s)
	}

	return nil
}

// Define private named constants for use when implementing the fmt.Stringer and
// fmt.Scanner interfaces.
const (
	one   = "One"
	two   = "Two"
	three = "Three"
	four  = "Four"
	five  = "Five"
	six   = "Six"
)
