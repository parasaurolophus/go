// Copyright Kirk Rader 2024

package utilities

import (
	"encoding/json"
	"fmt"
)

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
type Money interface {
	fmt.Scanner
	fmt.Stringer
	json.Marshaler
	json.Unmarshaler
	Value() float64
}
