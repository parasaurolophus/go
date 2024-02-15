// Copright Kirk Rader 2024

package logging

import "testing"

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
