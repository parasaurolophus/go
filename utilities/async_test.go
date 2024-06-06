// Copyright Kirk Rader 2024

package utilities

import (
	"bytes"
	"encoding/json"
	"parasaurolophus/go/logging"
	"testing"
)

func TestAsync(t *testing.T) {
	buffer := bytes.Buffer{}
	options := logging.LoggerOptions{
		BaseTags: []string{t.Name()},
	}
	logger := logging.New(&buffer, &options)
	in := make(chan int)
	out := make(chan int)
	defer close(out)
	asyncFunction := func(n int) int {
		if n%2 == 0 {
			return n + 1
		}
		panic("odd")
	}
	panicHandler := func(r any) {
		logger.Always(
			nil,
			logging.RECOVERED, r,
			logging.TAGS, t.Name(),
			logging.STACKTRACE, nil,
		)
	}
	go Async(asyncFunction, out, in, panicHandler)
	out <- 0
	v := <-in
	if v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
	out <- 1
	v = <-in
	if v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
	b := buffer.Bytes()
	entry := map[string]any{}
	err := json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal(err.Error())
	}
	if entry["recovered"] != "odd" {
		t.Errorf(`expected entry to have "recovered" field value of "odd", got %v`, entry["recovered"])
	}
}
