// Copyright Kirk Rader 2024

package utilities

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode"
)

type (
	// Wrap a float64 in a struct which specifies the number of digits to emit
	// for its fractional part when converting to a string or marshaling JSON.
	//
	// This type's methods only affect how the underlying floatint-point value
	// is represented in text-based formats. It does not alter the mathematical
	// precision of monetary values nor perform any scaling based on the
	// denominations of particular currencies. For example, the appropriate
	// number of digits to use for USD, CAD, GBP, EUR etc. is 2. The appropriate
	// number of digits to use for JPY is 0. 100 JPY would be represented by the
	// float64 value 100.0.
	Money interface {
		fmt.Scanner
		fmt.Stringer
		json.Marshaler
		json.Unmarshaler
		Value() float64
	}

	// Implementation of Money interface.
	money struct {
		value  float64
		digits int
	}
)

// Create a Money structure initialized to the given values.
func New(value float64, digits int) Money {
	if digits < 0 {
		digits = 0
	}
	return &money{
		value:  value,
		digits: digits,
	}
}

// Return the JSON representation of m and nil.
func (m money) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

// Set m's value to the float64 represented by the given fmt.ScanState.
//
// Returns nil if m's value was successfully parsed, an error if parsing the
// contents of the given fmt.ScanState as a float64 fails. The value of m is
// left unchanged if error is non-nil.
func (m *money) Scan(state fmt.ScanState, _ rune) error {
	token, err := state.Token(true, m.makeTokenizer())
	if err != nil {
		return err
	}
	f, err := strconv.ParseFloat(string(token), 64)
	if err != nil {
		return err
	}
	m.value = f
	return nil
}

// Return the string representation of m's value.
func (m money) String() string {
	f := "%." + strconv.Itoa(m.digits) + "f"
	return fmt.Sprintf(f, m.Value())
}

// Set m's value to the result of converting the given sequence of bytes to a
// float64.
//
// Returns nil if m's value was successfully parsed, an error if parsing the
// contents of the given sequence of bytes as a float64 fails. The value of m is
// left unchanged if error is non-nil.
func (m *money) UnmarshalJSON(b []byte) error {
	var f float64
	_, err := fmt.Sscan(string(b), &f)
	if err == nil {
		m.value = f
	}
	return err
}

// Return the numeric value of m.
func (m money) Value() float64 {
	return m.value
}

// Create a function for use by Money's Scan() method to tokenize a
// floating-point value.
func (m money) makeTokenizer() func(rune) bool {
	firstRune := true
	decimalPointSeen := false
	return func(r rune) bool {
		defer func() { firstRune = false }()
		if r == '-' {
			return firstRune
		}
		if r == '.' {
			if decimalPointSeen {
				return false
			}
			decimalPointSeen = true
			return true
		}
		return unicode.IsDigit(r)
	}
}
