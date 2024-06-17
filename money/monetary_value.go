// Copyright Kirk Rader 2024

package money

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"unicode"
)

// Implementation of Money interface.
type monetaryValue struct {
	value  float64
	digits int
}

///////////////////////////////////////////////////////////////////////////////
// Implement utilities.Money

// Return the numeric value of m.
func (m monetaryValue) Get() float64 {
	return m.value
}

// Update the number of digits to the right of the decimal point when
// representing m in a text-based format.
func (m *monetaryValue) SetDigits(digits int) (err error) {
	if digits < 0 {
		err = fmt.Errorf("negative number of digits specified: %d", digits)
		return
	}
	m.digits = digits
	return
}

// Update the numeric value of m.
func (m *monetaryValue) SetValue(value float64) {
	m.value = value
}

///////////////////////////////////////////////////////////////////////////////
// Implement fmt.Scanner

// Set m's value to the float64 represented by the given fmt.ScanState.
//
// Returns nil if m's value was successfully parsed, an error if parsing the
// contents of the given fmt.ScanState as a float64 fails. The value of m is set
// to 0.0 if error is non-nil.
func (m *monetaryValue) Scan(state fmt.ScanState, _ rune) (err error) {
	makeTokenTest := func() func(rune) bool {
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
	token, err := state.Token(true, makeTokenTest())
	if err != nil {
		return
	}
	m.value, err = strconv.ParseFloat(string(token), 64)
	return
}

///////////////////////////////////////////////////////////////////////////////
// Implement fmt.Stringer

// Return the string representation of m's value.
func (m monetaryValue) String() string {
	return strconv.FormatFloat(m.value, 'f', m.digits, 64)
}

///////////////////////////////////////////////////////////////////////////////
// Implement json.Marshaler

// Return the JSON representation of m and nil.
func (m monetaryValue) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

///////////////////////////////////////////////////////////////////////////////
// Implement xml.Marshaler

// Return the result of calling encoder.EncodeElement(m.value, start).
func (m monetaryValue) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(m.String(), start)
}

///////////////////////////////////////////////////////////////////////////////
// Implement json.Unmarshaler

// Set m's value to the result of converting the given sequence of bytes to a
// float64.
//
// Returns nil if m's value was successfully parsed, an error if parsing the
// contents of the given sequence of bytes as a float64 fails. The value of m is
// left unchanged if error is non-nil.
func (m *monetaryValue) UnmarshalJSON(b []byte) error {
	var f float64
	_, err := fmt.Sscan(string(b), &f)
	if err == nil {
		m.value = f
	}
	return err
}

///////////////////////////////////////////////////////////////////////////////
// Implement xml.Unmarshaler

// Return the result of calling decoder.Decode(&m.value, &start).
func (m *monetaryValue) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return decoder.DecodeElement(&(m.value), &start)
}

///////////////////////////////////////////////////////////////////////////////
// Implement xml.UnmarshalerAttr

// Return the result of calling fmt.Sscan(attr.Value, m).
func (m *monetaryValue) UnmarshalXMLAttr(attr xml.Attr) error {
	_, err := fmt.Sscan(attr.Value, m)
	return err
}
