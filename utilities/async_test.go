// Copyright Kirk Rader 2024

package utilities

import (
	"bytes"
	"context"
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
	type S struct {
		Value int
	}
	asyncFunction := func(s *S) *S {
		if s.Value%2 == 0 {
			s.Value += 1
			return s
		}
		panic("odd")
	}
	in := make(chan *S)
	out := make(chan *S)
	defer close(out)
	panicHandler := func(r any) {
		logger.Always(
			context.Background(),
			nil,
			logging.RECOVERED, r,
			logging.TAGS, t.Name(),
			logging.STACKTRACE, nil,
		)
	}
	go Async(asyncFunction, out, in, panicHandler)
	s := S{Value: 0}
	out <- &s
	v := <-in
	if v == nil {
		t.Fatalf("expected result not to be nil")
	}
	if v != &s {
		t.Fatalf("expected v to be &s")
	}
	if s.Value != 1 {
		t.Errorf("expected 1, got %d", s.Value)
	}
	out <- &s
	v = <-in
	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}
	if s.Value != 1 {
		t.Errorf("expected 1, got %d", s.Value)
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
