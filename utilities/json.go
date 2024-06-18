// Copyright Kirk Rader 2024

package utilities

import "unicode"

// Return a closure that can be used with fmt.ScanState.Token to convert the
// text representation of a number to a float64 according to JSON number syntax.
func MakeJSONNumberTokenTest() func(rune) bool {
	firstRune := true
	digitSeen := false
	decimalSeen := false
	exponentSeen := false
	return func(r rune) bool {
		defer func() { firstRune = false }()
		if r == '-' {
			return firstRune || (exponentSeen && !digitSeen)
		}
		if r == 'e' || r == 'E' {
			defer func() {
				exponentSeen = true
				digitSeen = false
			}()
			return digitSeen && !exponentSeen
		}
		if r == '.' {
			defer func() { decimalSeen = true }()
			return digitSeen && !decimalSeen
		}
		if unicode.IsDigit(r) {
			defer func() { digitSeen = true }()
			return true
		}
		return false
	}
}
