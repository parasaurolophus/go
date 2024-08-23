// Copyright 2024 Kirk Rader

package utilities_test

import (
	"bytes"
	"parasaurolophus/utilities"
	"sync"
	"testing"
)

func TestNewWorker(t *testing.T) {
	log := new(bytes.Buffer)
	actual := 0
	handler := func(n int) {
		actual += n
	}
	values, await := utilities.NewWorker(handler, log)
	values <- 1
	values <- 1
	close(values)
	<-await
	if actual != 2 {
		t.Errorf("expected 2, got %d", actual)
	}
	s := log.String()
	if len(s) != 0 {
		t.Errorf(`expected no errors to be logged, got "%s"`, s)
	}
}

func TestNewWorkerWaitGroup(t *testing.T) {
	log := new(bytes.Buffer)
	actual := 0
	handler := func(n int) {
		actual += n
	}
	waitGroup := sync.WaitGroup{}
	values := make([]chan<- int, 3)
	for i := range 3 {
		values[i] = utilities.NewWorkerWaitGroup(handler, &waitGroup, log)
	}
	for _, v := range values {
		v <- 1
		close(v)
	}
	waitGroup.Wait()
	if actual != len(values) {
		t.Errorf("expected %d, got %d", len(values), actual)
	}
	s := log.String()
	if len(s) != 0 {
		t.Errorf(`expected no errors to be logged, got "%s"`, s)
	}
}

func TestNeworkerTrapPanic(t *testing.T) {
	log := new(bytes.Buffer)
	actual := 0
	handler := func(n int) {
		panic("deliberate")
	}
	values, await := utilities.NewWorker(handler, log)
	values <- 1
	values <- 1
	close(values)
	<-await
	if actual != 0 {
		t.Errorf("expected 0, got %d", actual)
	}
	s := log.String()
	if len(s) == 0 {
		t.Error(`expected errors to be logged`)
	}
}
