// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"testing"
)

func TestAsync(t *testing.T) {

	out := make(chan int)
	in := make(chan int)
	adder := func(n int) int {

		if n >= 3 {
			panic("out of range")
		}

		return n + 1
	}
	result := []int{}

	go Async(adder, out, in)

	go func() {

		defer close(out)

		for n := range 2 {

			out <- n
		}

		out <- 3
		out <- 2
	}()

	for n := range in {

		result = append(result, n)
	}

	if len(result) != 3 {
		t.Errorf("expected 3, got %d", len(result))
	}

	if result[0] != 1 {
		t.Errorf("expected 1, got %d", result[0])
	}

	if result[1] != 2 {
		t.Errorf("expected 2, got %d", result[1])
	}

	if result[2] != 3 {
		t.Errorf("expected 3, got %d", result[2])
	}
}

func TestCompose(t *testing.T) {
	adder := func(n int) int { return n + 1 }
	subtracter := func(n int) int { return n - 1 }
	result, err := Compose(0, adder, adder, subtracter, adder)
	if err != nil {
		t.Error(err.Error())
	}
	if result != 2 {
		t.Errorf("expected 2, got %d", result)
	}
	panicer := func(n int) int { panic("deliberate") }
	result, err = Compose(0, adder, adder, panicer, adder)
	if err == nil {
		t.Error("expected err to be non-nil")
	}
	fmt.Println(err.Error())
	if result != 2 {
		t.Errorf("expected 2, got %d", result)
	}
}

func TestMap(t *testing.T) {

	slice := []int{0, 1, 2}
	adder := func(n int) int { return n + 1 }
	result := Map(adder, slice)

	if len(result) != 3 {
		t.Errorf("expected 3, got %d", len(result))
	}

	if result[0] != 1 {
		t.Errorf("expected 1, got %d", result[0])
	}

	if result[1] != 2 {
		t.Errorf("expected 2, got %d", result[1])
	}

	if result[2] != 3 {
		t.Errorf("expected 3, got %d", result[2])
	}
}
