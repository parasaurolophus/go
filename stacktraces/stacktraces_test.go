// Copyright Kirk Rader 2024

package stacktraces

import (
	"parasaurolophus/go/stacktraces_test"
	"strings"
	"testing"
)

func TestStackTrace(t *testing.T) {

	_, err := func() (any, error) {
		return nil, New("test", 0)
	}()

	switch v := err.(type) {
	case StackTrace:
		trace := v.LongTrace()
		if !strings.HasPrefix(trace, "0:runtime.Callers") {
			t.Fatalf("expected long trace to start with '0:runtime.Callers()', got '%s'", trace)
		}
		trace = v.ShortTrace()
		if !strings.HasPrefix(trace, "0:runtime.Callers") {
			t.Fatalf("expected short trace to start with '0:runtime.Callers()', got '%s'", trace)
		}

	default:
		t.Fatalf("expected err to be a StackTrace")
	}

	msg := err.Error()
	if msg != "test" {
		t.Fatalf(
			"expected Error() to return 'test', got '%s'",
			msg)
	}
}

func TestStackFrameNameNotFound(t *testing.T) {

	trace := New("", "fubar").ShortTrace()

	if trace != "" {
		t.Fatalf("expected empty stack trace, got '%s", trace)
	}
}

func TestStackTraceString(t *testing.T) {

	var expected string

	_, err := func() (any, error) {
		expected, _, _ = FunctionInfo()
		return nil, New("test", expected)
	}()

	switch v := err.(type) {

	case StackTrace:
		trace := v.LongTrace()
		functionName, _, err := stacktraces_test.FirstFunctionLong(trace)

		if err != nil {
			t.Fatalf("error parsing stack frames; %s", err.Error())
		}

		if functionName != expected {
			t.Fatalf("expected long trace to start with '%s', got '%s'", expected, functionName)
		}

		trace = v.ShortTrace()
		functionName, _, err = stacktraces_test.FirstFunctionShort(trace)

		if err != nil {
			t.Fatalf("TestAlways: error parsing stack frames; %s", err.Error())
		}

		if functionName != expected {
			t.Fatalf("expected short trace to start with '%s', got '%s'", expected, functionName)
		}

	default:
		t.Fatalf("expected err to be a StackTrace")
	}

	msg := err.Error()

	if msg != "test" {
		t.Fatalf(
			"expected Error() to return 'test', got '%s'",
			msg)
	}
}

func TestStackTraceNameNotFound(t *testing.T) {

	_, err := func() (any, error) {
		return nil, New("test", "fubar")
	}()

	switch v := err.(type) {
	case StackTrace:
		trace := v.LongTrace()
		if trace != "" {
			t.Fatalf(
				"expected long trace to be empty, got '%s'",
				trace)
		}
		trace = v.ShortTrace()
		if trace != "" {
			t.Fatalf(
				"expected short trace to be empty, got '%s'",
				trace)
		}

	default:
		t.Fatalf("expected err to be a StackTrace")
	}

	msg := err.Error()

	if msg != "test" {
		t.Fatalf(
			"expected Error() to return 'test', got '%s'",
			msg)
	}
}

func TestStackTraceDepth(t *testing.T) {

	trace := New("test", 100)

	if trace.LongTrace() != "" {
		t.Fatalf(
			"expected long stack trace to be empty, got '%s'",
			trace)
	}

	if trace.ShortTrace() != "" {
		t.Fatalf(
			"expected short stack trace to be empty, got '%s'",
			trace)
	}
}

func TestStackTraceFloat(t *testing.T) {

	functionName, _, _ := FunctionInfo()
	expected := functionName + ".func1"

	_, stacktrace := func() (any, StackTrace) {
		return nil, New("test", 1.1)
	}()

	trace := stacktrace.LongTrace()
	functionName, _, err := stacktraces_test.FirstFunctionLong(trace)

	if err != nil {
		t.Fatalf("error parsing long stack trace: %s", err.Error())
	}

	if functionName != expected {
		t.Fatalf("expected long trace to start with '%s', got '%s'", expected, trace)
	}

	trace = stacktrace.ShortTrace()

	functionName, _, err = stacktraces_test.FirstFunctionShort(trace)

	if err != nil {
		t.Fatalf("error parsing short stack trace: %s", err.Error())
	}

	if functionName != expected {
		t.Fatalf("expected short trace to start with '%s', got '%s'", expected, trace)
	}

	msg := stacktrace.Error()
	if msg != "test" {
		t.Fatalf(
			"expected Error() to return 'test', got '%s'",
			msg)
	}
}

func TestStackTraceAuto(t *testing.T) {

	functionName, _, _ := FunctionInfo()
	trace := New("test", -1)
	long := trace.LongTrace()
	short := trace.ShortTrace()
	firstLong, m, err := stacktraces_test.FirstFunctionLong(long)

	if err != nil {
		t.Fatalf("error parsing long frames; %s", err.Error())
	}

	if m <= 0 {
		t.Fatalf("expected frame > 0, got %d", m)
	}

	firstShort, n, err := stacktraces_test.FirstFunctionShort(short)

	if err != nil {
		t.Fatalf("error parsing short frames; %s", err.Error())
	}

	if m != n {
		t.Fatalf("expected frame numbers to be the same, got %d, %d", m, n)
	}

	if firstLong != functionName {
		t.Fatalf("expected long trace to start with '%s', got '%s'", functionName, firstLong)
	}

	if firstShort != functionName {
		t.Fatalf("expected short trace to start with '%s', got '%s'", functionName, firstShort)
	}
}

func TestLongStackTraceAuto(t *testing.T) {

	functionName, _, _ := FunctionInfo()
	trace := LongStackTrace(-1)
	name, n, err := stacktraces_test.FirstFunctionLong(trace)

	if err != nil {
		t.Fatalf("error parsing truncated stack trace; %s", err.Error())
	}

	if n <= 0 {
		t.Fatalf("expected frame number to be greater than 0, got %d", n)
	}

	if name != functionName {
		t.Fatalf("expected first function to be '%s', got '%s'", functionName, name)
	}
}

func TestLongStackTraceInt(t *testing.T) {

	trace := LongStackTrace(0)
	name, n, err := stacktraces_test.FirstFunctionLong(trace)

	if err != nil {
		t.Fatalf("error parsing full stack trace; %s", err.Error())
	}

	if n != 0 {
		t.Fatalf("expected first frame in full stack trace to be 0, got %d", n)
	}

	if name != "runtime.Callers" {
		t.Fatalf("expected first function to be 'runtime.Callers', got '%s'", name)
	}
}

func TestLongStackTraceString(t *testing.T) {

	name, _, _ := FunctionInfo()
	trace := LongStackTrace(name)

	if trace == "" {
		t.Fatalf("expected trace not to be empty")
	}
}

func TestLongStackTraceFloat(t *testing.T) {

	trace := LongStackTrace(12.7)

	if trace == "" {
		t.Fatalf("expected trace not to be empty")
	}
}

func TestShortStackTraceInt(t *testing.T) {

	trace1 := ShortStackTrace(4)
	_, m, err := stacktraces_test.FirstFunctionShort(trace1)

	if err != nil {
		t.Fatalf("error parsing full stack trace; %s", err.Error())
	}

	if m != 4 {
		t.Fatalf("expected first frame to be 4, got %d", m)
	}

	trace2 := ShortStackTrace(5)
	_, n, err := stacktraces_test.FirstFunctionShort(trace2)

	if err != nil {
		t.Fatalf("error parsing truncated stack trace; %s", err.Error())
	}

	if n != 5 {
		t.Fatalf("expected first frame to be 5, got %d", n)
	}

	diff := n - m

	if diff != 1 {
		t.Fatalf("expected 1 more frame in trace1 than trace2, got %d", diff)
	}

	if !strings.HasSuffix(trace1, trace2) {

		t.Fatalf("expected '%s' to be a suffix of '%s'", trace2, trace1)
	}
}

func TestShortStackTraceString(t *testing.T) {

	expected, _, _ := FunctionInfo()
	trace := ShortStackTrace(expected)
	name, _, err := stacktraces_test.FirstFunctionShort(trace)

	if err != nil {
		t.Fatalf("error parsing frame; %s", err.Error())
	}

	if name != expected {
		t.Fatalf("expected trace to start with '%s', got '%s'", expected, name)
	}
}

func TestShortStackTraceFloat(t *testing.T) {

	trace := ShortStackTrace(12.7)
	functionName, _, _ := FunctionInfo()
	name, _, err := stacktraces_test.FirstFunctionShort(trace)

	if err != nil {
		t.Fatalf("error parsing short stack trace: %s", err.Error())
	}

	if name != functionName {
		t.Fatalf("expected trace to start with '%s', got '%s'", functionName, trace)
	}
}
