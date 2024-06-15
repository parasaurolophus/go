// Copyright Kirk Rader 2024

package money

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// An interface for representing monetary values in text-based formats.
type Money interface {
	fmt.Scanner
	fmt.Stringer
	json.Marshaler
	json.Unmarshaler
	xml.Marshaler
	xml.Unmarshaler
	xml.UnmarshalerAttr
	Get() float64
}

// Create a monetaryValue structure initialized to the given values. The second
// parameter specifies the number of digits to emit when converting to text
// based representations, e.g. 2 for currencies like USD, EUR; 0 for JPY; etc.
func NewMoney(value float64, digits int) Money {
	if digits < 0 {
		digits = 0
	}
	return &monetaryValue{
		value:  value,
		digits: digits,
	}
}