// Copyright 2024 Kirk Rader

package utilities_test

import (
	"bytes"
	"parasaurolophus/utilities"
	"testing"
)

func TestInvoke(t *testing.T) {
	log := new(bytes.Buffer)
	actual := 0
	utilities.Invoke(
		func(n int) {
			actual += n
		},
		1,
		log)
	if actual != 1 {
		t.Errorf("expected 1, got %d", actual)
	}
	s := log.String()
	if len(s) != 0 {
		t.Errorf(`expected no panic to have been loggged, got "%s"`, s)
	}
}

func TestInvokePanic(t *testing.T) {
	log := new(bytes.Buffer)
	utilities.Invoke(
		func(_ int) {
			panic("deliberate")
		},
		1,
		log)
	s := log.String()
	if len(s) == 0 {
		t.Error("expected panic to have been loggged")
	}
}

func TestInvokePanicNil(t *testing.T) {
	actual := 0
	utilities.Invoke(
		func(_ int) {
			actual++
			panic("deliberate")
		},
		1,
		nil)
	if actual != 1 {
		t.Errorf("expected 1, got %d", actual)
	}
}
