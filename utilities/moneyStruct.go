// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"strconv"
	"unicode"
)

// Implementation of Money interface.
type moneyStruct struct {
	value  float64
	digits int
}

// Create a Money structure initialized to the given values.
func NewMoney(value float64, digits int) Money {
	if digits < 0 {
		digits = 0
	}
	return &moneyStruct{
		value:  value,
		digits: digits,
	}
}

// Return the JSON representation of m and nil.
func (m moneyStruct) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

// Set m's value to the float64 represented by the given fmt.ScanState.
//
// Returns nil if m's value was successfully parsed, an error if parsing the
// contents of the given fmt.ScanState as a float64 fails. The value of m is
// left unchanged if error is non-nil.
func (m *moneyStruct) Scan(state fmt.ScanState, _ rune) error {
	makeTokenizer := func() func(rune) bool {
		firstRune := true
		decimalPointSeen := false
		return func(r rune) bool {
			defer func() { firstRune = false }()
			if r == '-' {
				return firstRune
			}
			if r == '.' {
				defer func() { decimalPointSeen = true }()
				return !decimalPointSeen
			}
			return unicode.IsDigit(r)
		}
	}
	token, err := state.Token(true, makeTokenizer())
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
func (m moneyStruct) String() string {
	f := "%." + strconv.Itoa(m.digits) + "f"
	return fmt.Sprintf(f, m.Value())
}

// Set m's value to the result of converting the given sequence of bytes to a
// float64.
//
// Returns nil if m's value was successfully parsed, an error if parsing the
// contents of the given sequence of bytes as a float64 fails. The value of m is
// left unchanged if error is non-nil.
func (m *moneyStruct) UnmarshalJSON(b []byte) error {
	var f float64
	_, err := fmt.Sscan(string(b), &f)
	if err == nil {
		m.value = f
	}
	return err
}

// Return the numeric value of m.
func (m moneyStruct) Value() float64 {
	return m.value
}
