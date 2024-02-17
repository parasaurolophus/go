// Copright Kirk Rader 2024

package logging

import (
	"fmt"
	"parasaurolophus/go/stacktraces"
	"testing"
)

func TestTraceString(t *testing.T) {
	s := TRACE.String()
	if s != "TRACE" {
		t.Fatalf("expected 'TRACE', got '%s'", s)
	}
}

func TestFineString(t *testing.T) {
	s := FINE.String()
	if s != "FINE" {
		t.Fatalf("expected 'FINE', got '%s'", s)
	}
}

func TestOptionalString(t *testing.T) {
	s := OPTIONAL.String()
	if s != "OPTIONAL" {
		t.Fatalf("expected 'OPTIONAL', got '%s'", s)
	}
}

func TestAlwaysString(t *testing.T) {
	s := ALWAYS.String()
	if s != "ALWAYS" {
		t.Fatalf("expected 'ALWAYS', got '%s'", s)
	}
}

func TestUnrecognizedVerbosityString(t *testing.T) {
	s := Verbosity(100).String()
	if s != "100" {
		t.Fatalf("expected '100', got '%s'", s)
	}
}

func TestTraceScan(t *testing.T) {
	var verbosity Verbosity
	n, err := fmt.Sscanf("foo TRACE bar", "foo %v bar", &verbosity)
	if err != nil {
		switch e := err.(type) {
		case stacktraces.StackTrace:
			t.Fatalf("sscanf returned error: %s (%s)", e.Error(), e.ShortTrace())
		default:
			t.Fatalf("sscanf returned error: %s", e.Error())
		}
	}
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
	if verbosity != TRACE {
		t.Fatalf("expected %s, got %s", TRACE, verbosity)
	}
}

func TestFineScan(t *testing.T) {
	var verbosity Verbosity
	n, err := fmt.Sscanf("foo FINE bar", "foo %v bar", &verbosity)
	if err != nil {
		switch e := err.(type) {
		case stacktraces.StackTrace:
			t.Fatalf("sscanf returned error: %s (%s)", e.Error(), e.ShortTrace())
		default:
			t.Fatalf("sscanf returned error: %s", e.Error())
		}
	}
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
	if verbosity != FINE {
		t.Fatalf("expected %s, got %s", FINE, verbosity)
	}
}

func TestOptionalScan(t *testing.T) {
	var verbosity Verbosity
	n, err := fmt.Sscanf("foo OPTIONAL bar", "foo %v bar", &verbosity)
	if err != nil {
		switch e := err.(type) {
		case stacktraces.StackTrace:
			t.Fatalf("sscanf returned error: %s (%s)", e.Error(), e.ShortTrace())
		default:
			t.Fatalf("sscanf returned error: %s", e.Error())
		}
	}
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
	if verbosity != OPTIONAL {
		t.Fatalf("expected %s, got %s", OPTIONAL, verbosity)
	}
}

func TestAlwaysScan(t *testing.T) {
	var verbosity Verbosity
	n, err := fmt.Sscanf("foo ALWAYS bar", "foo %v bar", &verbosity)
	if err != nil {
		switch e := err.(type) {
		case stacktraces.StackTrace:
			t.Fatalf("sscanf returned error: %s (%s)", e.Error(), e.ShortTrace())
		default:
			t.Fatalf("sscanf returned error: %s", e.Error())
		}
	}
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
	if verbosity != ALWAYS {
		t.Fatalf("expected %s, got %s", ALWAYS, verbosity)
	}
}

func TestIntScan(t *testing.T) {
	var verbosity Verbosity
	n, err := fmt.Sscanf("foo 100 bar", "foo %v bar", &verbosity)
	if err != nil {
		switch e := err.(type) {
		case stacktraces.StackTrace:
			t.Fatalf("sscanf returned error: %s (%s)", e.Error(), e.ShortTrace())
		default:
			t.Fatalf("sscanf returned error: %s", e.Error())
		}
	}
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
	if verbosity != Verbosity(100) {
		t.Fatalf("expected %s, got %s", Verbosity(100), verbosity)
	}
}

func TestUnsupportedScan(t *testing.T) {
	verbosity := ALWAYS
	n, err := fmt.Sscanf("foo FUBAR bar", "foo %v bar", &verbosity)
	if err == nil {
		t.Fatalf("expected error to be signaled")
	}
	if n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
	if verbosity != ALWAYS {
		t.Fatalf("expected verbosity not to have been changed, got %s", verbosity)
	}
}
