// Copyright 2024 Kirk Rader

package utilities_test

import (
	"parasaurolophus/utilities"
	"testing"
	"time"
)

func TestWithTimeLimit(t *testing.T) {
	fn1 := func() int {
		return 1
	}
	fn2 := func() int {
		time.Sleep(time.Millisecond * 100)
		return 2
	}
	timeout := func(time.Time) int {
		return 0
	}
	if v := utilities.WithTimeLimit(fn1, timeout, time.Millisecond*50); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
	if v := utilities.WithTimeLimit(fn2, timeout, time.Millisecond*50); v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}
