// Copyright 2024 Kirk Rader

package utilities_test

import (
	"parasaurolophus/utilities"
	"sync"
	"testing"
)

func TestStartWorkers(t *testing.T) {
	actual := 0
	lock := sync.Mutex{}
	handler := func(n int) {
		defer lock.Unlock()
		lock.Lock()
		actual += n
	}
	func() {
		values, await := utilities.StartWorkers(3, 1, handler)
		defer utilities.CloseAllAndWait(values, await)
		for i := range len(values) {
			values[i] <- i
		}
	}()
	if actual != 3 {
		t.Errorf("expected 3, got %d", actual)
	}
}
