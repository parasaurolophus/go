// Copyright 2024 Kirk Rader

package utilities_test

import (
	"parasaurolophus/utilities"
	"testing"
	"time"
)

func TestWatchdog(t *testing.T) {
	count := 0
	timeout := func() {
		count++
	}
	watchdog, reset := utilities.NewWatchdog(time.Millisecond*10, timeout)
	for range 3 {
		time.Sleep(time.Millisecond * 5)
		reset <- true
	}
	time.Sleep(time.Millisecond * 15)
	watchdog.Stop()
	if count != 1 {
		t.Errorf("expected count to be 1, got %d", count)
	}
}
