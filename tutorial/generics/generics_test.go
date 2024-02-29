// Copyright Kirk Rader 2024

package generics

import (
	"testing"
)

func TestSumOfInts(t *testing.T) {

	n := Sum(1, 2, -1)

	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}

func TestSumOfUInts(t *testing.T) {

	var u uint = Sum(uint(1), 2, 3)

	if u != 6 {
		t.Errorf("expected 6, got %d", u)
	}
}

func TestSumOfFloat64s(t *testing.T) {

	f := Sum(1.1, 2, -0.9)

	if f != 2.2 {
		t.Errorf("expected 6, got %f", f)
	}
}

func TestSumOfNone(t *testing.T) {

	f := Sum[float64]()

	if f != 0 {
		t.Errorf("expected 0, got %f", f)
	}
}
