// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"strconv"
	"unicode"
)

type (
	Money struct {
		value  float64
		digits int
	}
)

func New(value float64, digits int) Money {
	if digits < 0 {
		digits = 0
	}
	return Money{
		value:  value,
		digits: digits,
	}
}

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Money) Scan(state fmt.ScanState, _ rune) error {
	token, err := state.Token(true, m.tokenizer())
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

func (m Money) String() string {
	f := "%." + strconv.Itoa(m.digits) + "f"
	return fmt.Sprintf(f, m.Value())
}

func (m *Money) UnmarshalJSON(b []byte) error {
	var f float64
	_, err := fmt.Sscan(string(b), &f)
	if err == nil {
		m.value = f
	}
	return err
}

func (m Money) Value() float64 {
	return m.value
}

func (m Money) tokenizer() func(rune) bool {
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
