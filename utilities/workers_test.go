// Copyright 2024 Kirk Rader

package utilities_test

import (
	"parasaurolophus/utilities"
	"sync"
	"testing"
	"time"
)

func TestProcessBatch(t *testing.T) {
	actual := 0.0
	generate := func(producers []chan<- int) {
		for i := range 9 {
			producers[i%len(producers)] <- i
		}
	}
	transfrom := func(input int) float64 {
		return float64(input) / 2.0
	}
	consume := func(output float64) {
		actual += output
	}
	utilities.ProcessBatch(3, generate, transfrom, consume)
	if actual != 18 {
		t.Errorf("expected 28, got %f", actual)
	}
}

func TestStartWorker(t *testing.T) {
	actual := 0
	lock := sync.Mutex{}
	handler := func(n int) {
		defer lock.Unlock()
		lock.Lock()
		actual += n
	}
	func() {
		values, await := utilities.StartWorker(handler)
		defer utilities.CloseAndWait(values, await)
		for i := range 3 {
			values <- i
		}
	}()
	if actual != 3 {
		t.Errorf("expected 3, got %d", actual)
	}
}

func TestStartWorkers(t *testing.T) {
	actual := 0
	lock := sync.Mutex{}
	handler := func(n int) {
		defer lock.Unlock()
		lock.Lock()
		actual += n
	}
	func() {
		values, await := utilities.StartWorkers(3, handler)
		defer utilities.CloseAllAndWait(values, await)
		for i := range len(values) {
			values[i] <- i
		}
	}()
	if actual != 3 {
		t.Errorf("expected 3, got %d", actual)
	}
}

func TestWithTimeLimit(t *testing.T) {
	fn1 := func() int {
		return 1
	}
	fn2 := func() int {
		time.Sleep(time.Millisecond * 100)
		return 2
	}
	timeout := func() int {
		return 0
	}
	if v := utilities.WithTimeLimit(fn1, timeout, time.Millisecond*50); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
	if v := utilities.WithTimeLimit(fn2, timeout, time.Millisecond*50); v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}
