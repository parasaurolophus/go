// Copyright 2024 Kirk Rader

package utilities_test

import (
	"parasaurolophus/utilities"
	"strconv"
	"testing"
)

func TestProcessBatch(t *testing.T) {
	actual := []string{}
	generate := func(transformers []chan<- int) {
		n := len(transformers)
		for i := range 10 {
			transformers[i%n] <- i
		}
	}
	transform := func(input int) string {
		return strconv.Itoa(input)
	}
	consume := func(output string) {
		actual = append(actual, output)
	}
	utilities.ProcessBatch(3, 1, 1, generate, transform, consume)
	if len(actual) != 10 {
		t.Errorf("expected 10, got %d", len(actual))
	}
}
