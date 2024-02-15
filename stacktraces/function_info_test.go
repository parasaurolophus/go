// Copyright Kirk Rader 2024

package stacktraces

import (
	"parasaurolophus/go/stacktraces_test"
	"strings"
	"testing"
)

func TestFunctionInfoNil(t *testing.T) {

	actualName, actualFile, actuaLine, err := FunctionInfo(nil)

	if err != nil {
		t.Fatalf("error getting function info for nil, %s (%s)", err.Error(), err.(StackTrace).ShortTrace())
	}

	if !strings.HasSuffix(actualName, ".TestFunctionInfoNil") {
		t.Fatalf("expected name to end with '.TestFunctionInfoNil', got '%s'", actualName)
	}

	if !strings.HasSuffix(actualFile, "/function_info_test.go") {
		t.Fatalf("expected file name to be 'function_info_test.go', got '%s'", actualFile)
	}

	if actuaLine != 13 {
		t.Fatalf("expected line to be 13, got %d", actuaLine)
	}
}

func TestFunctionInfoAuto(t *testing.T) {

	actualName, actualFile, actuaLine, err := FunctionInfo(-1)

	if err != nil {
		t.Fatalf("error getting function info for -1, %s (%s)", err.Error(), err.(StackTrace).ShortTrace())
	}

	if !strings.HasSuffix(actualName, ".TestFunctionInfoAuto") {
		t.Fatalf("expected name to end with '.TestFunctionInfoAuto', got '%s'", actualName)
	}

	if !strings.HasSuffix(actualFile, "/function_info_test.go") {
		t.Fatalf("expected file name to be 'function_info_test.go', got '%s'", actualFile)
	}

	if actuaLine != 34 {
		t.Fatalf("expected line to be 34, got %d", actuaLine)
	}
}

func TestFunctionInfoSkip(t *testing.T) {

	actualName, actualFile, _, err := FunctionInfo(1)

	if err != nil {
		t.Fatalf("error getting function info for 1, %s (%s)", err.Error(), err.(StackTrace).ShortTrace())
	}

	if !strings.HasSuffix(actualName, ".functionInfo") {
		t.Fatalf("expected name to end with '.functionInfo', got '%s'", actualName)
	}

	if !strings.HasSuffix(actualFile, "/function_info.go") {
		t.Fatalf("expected file name to be 'function_info.go', got '%s'", actualFile)
	}
}

func TestFunctionInfoString(t *testing.T) {

	expected := FunctionName()
	actualName, actualFile, actualLine, err := FunctionInfo(expected)

	if err != nil {
		t.Fatalf("error getting function info for 1, %s (%s)", err.Error(), err.(StackTrace).ShortTrace())
	}

	if actualName != expected {
		t.Fatalf("expected name to be '%s', got '%s'", expected, actualName)
	}

	if !strings.HasSuffix(actualFile, "/function_info_test.go") {
		t.Fatalf("expected file name to be 'function_info_test.go', got '%s'", actualFile)
	}

	if actualLine != 73 {
		t.Fatalf("expected line to be 73, got %d", actualLine)
	}
}

func TestFunctionInfoStringNotFound(t *testing.T) {

	_, _, _, err := FunctionInfo("stacktraces.New")

	if err == nil {
		t.Fatalf("expected an error getting function info for 'stacktraces.New'")
	}

	if err.Error() != "no frame found for stacktraces.New" {
		t.Fatalf("expected error message to be 'no frame found for stacktraces.New', got '%s'", err.Error())
	}

	_, _, e := stacktraces_test.FirstFunctionLong(err.(StackTrace).LongTrace())

	if e != nil {
		t.Fatalf("error parsing stack trace: %s", e.Error())
	}
}

func TestFunctionName(t *testing.T) {

	actual := FunctionName()

	if !strings.HasSuffix(actual, ".TestFunctionName") {
		t.Fatalf("expected name to end with '.TestFunctionName', got '%s'", actual)
	}
}

func TestFunctionNameAtNil(t *testing.T) {

	actual, err := FunctionNameAt(nil)

	if err != nil {
		t.Fatalf("error getting function name for nil, %s (%s)", err.Error(), err.(StackTrace).ShortTrace())
	}

	if !strings.HasSuffix(actual, ".TestFunctionNameAtNil") {
		t.Fatalf("expected name to end with '.TestFunctionNameAtNil', got '%s'", actual)
	}
}
