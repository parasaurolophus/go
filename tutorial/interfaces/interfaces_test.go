// Copyright Kirk Rader 2024

package interfaces

import (
	"fmt"
	"testing"
)

func TestIncrement(t *testing.T) {

	counter := BasicCounter()

	n := counter.Value()

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}

	n = counter.Increment()

	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	n = counter.Value()

	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
}

func TestDecrement(t *testing.T) {

	counter := BasicCounter()

	n := counter.Value()

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}

	n = counter.Decrement()

	if n != -1 {
		t.Errorf("expected -1, got %d", n)
	}

	n = counter.Value()

	if n != -1 {
		t.Errorf("expected -1, got %d", n)
	}
}

func TestString(t *testing.T) {

	counter := BasicCounter()
	s := fmt.Sprint(counter)

	if s != "0" {
		t.Errorf("expected \"0\", got \"%s\"", s)
	}
}

func TestScan(t *testing.T) {

	counter1 := BasicCounter()
	counter2 := BasicCounter()

	// Note that you can't use &counter1 or &counter2 here for similar reasons
	// to those noted by a comment in the body of BasicCounter().
	n, err := fmt.Sscanf(" 100, 2 ", " %v, %v ", counter1, counter2)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}

	n = counter1.Value()

	if n != 100 {
		t.Errorf("expected 100, got %d", n)
	}

	n = counter2.Value()

	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}
