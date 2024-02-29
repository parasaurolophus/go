// Copyright Kirk Rader 2024

package closures

import "testing"

func TestClosures(t *testing.T) {

	getter1, setter1 := MakeClosures(0)
	getter2, setter2 := MakeClosures(0)

	value1 := getter1()
	value2 := getter2()

	if value1 != 0 {
		t.Errorf("expected 0, got %d", value1)
	}

	if value2 != 0 {
		t.Errorf("expected 0, got %d", value2)
	}

	value1 = setter1(42)

	if value1 != 42 {
		t.Errorf("expected 42, got %d", value1)
	}

	if value2 != 0 {
		t.Errorf("expected 0, got %d", value2)
	}

	setter2(-1)
	value2 = getter2()

	if value1 != 42 {
		t.Errorf("expected 42, got %d", value1)
	}

	if value2 != -1 {
		t.Errorf("expected -1, got %d", value2)
	}
}

func TestPassContinuations(t *testing.T) {

	// create a pair of transformer functions which, respectively, add and
	// subtract the given increment from their own arguments
	makeTransforms := func(increment int) (func(int) int, func(int) int) {

		add := func(i int) int { return i + increment }
		subtract := func(i int) int { return i - increment }
		return add, subtract
	}

	// a will add 1, s will subtract 1
	a, s := makeTransforms(1)

	// p will deliberately panic
	p := func(int) int { panic("deliberate") }

	// note that the parameters to PassContinuations() include p, which
	// explcitly calls panic(), and nil, which will cause a nil pointer
	// dereference panic
	//
	// if you run this test from a console or under a debugger you should see
	// those two panics logged to stderr but execution continue to completion
	// due to the use of defer and recover() inside the implementation of
	// PassContuations()
	n := PassContinuations(0, a, a, nil, s, p, a)

	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}
