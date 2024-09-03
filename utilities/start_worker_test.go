// Copyright 2024 Kirk Rader

package utilities_test

import (
	"parasaurolophus/utilities"
	"testing"
)

func TestStartWorker(t *testing.T) {
	actual := 0
	handler := func(n int) {
		actual += n
	}
	func() {
		values, await := utilities.StartWorker(1, handler)
		defer utilities.CloseAndWait(values, await)
		for i := range 10 {
			values <- i
		}
	}()
	if actual != 45 {
		t.Errorf("expected 45, got %d", actual)
	}
}
