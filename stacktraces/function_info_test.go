// Copyright Kirk Rader 2024

package stacktraces

import (
	"strings"
	"testing"
)

func TestFunctionInfoNil(t *testing.T) {
	frame, actualName, actualFile, actuaLine, found := FunctionInfo(nil)
	if !found {
		t.Fatalf("function info for nil returned false as its last value")
	}
	if frame != defaultSkip {
		t.Fatalf("expected frame to be %d, got %d", defaultSkip, frame)
	}
	if !strings.HasSuffix(actualName, ".TestFunctionInfoNil") {
		t.Fatalf("expected name to end with '.TestFunctionInfoNil', got '%s'", actualName)
	}
	if !strings.HasSuffix(actualFile, "/function_info_test.go") {
		t.Fatalf("expected file name to be 'function_info_test.go', got '%s'", actualFile)
	}
	if actuaLine != 11 {
		t.Fatalf("expected line to be 11, got %d", actuaLine)
	}
}

func TestFunctionInfoNegative(t *testing.T) {
	frame, _, _, _, found := FunctionInfo(-1)
	if !found {
		t.Fatalf("function info for nil returned false as its last value")
	}
	if frame != defaultSkip+1 {
		t.Fatalf("expected frame to be %d, got %d", defaultSkip+1, frame)
	}
}

func TestFunctionInfoPositive(t *testing.T) {
	_, actualName, actualFile, _, found := FunctionInfo(1)
	if !found {
		t.Fatalf("function info for nil returned false as its last value")
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
	frame, actualName, actualFile, actualLine, found := FunctionInfo(expected)
	if !found {
		t.Fatalf("function info for nil returned false as its last value")
	}
	if frame != defaultSkip {
		t.Fatalf("expected frame to be %d, got %d", defaultSkip, frame)
	}
	if actualName != expected {
		t.Fatalf("expected name to be '%s', got '%s'", expected, actualName)
	}
	if !strings.HasSuffix(actualFile, "/function_info_test.go") {
		t.Fatalf("expected file name to be 'function_info_test.go', got '%s'", actualFile)
	}
	if actualLine != 54 {
		t.Fatalf("expected line to be 54, got %d", actualLine)
	}
}

func TestFunctionInfoStringNotFound(t *testing.T) {
	frame, name, file, line, found := FunctionInfo("stacktraces.New")
	if found {
		t.Fatalf(
			"expected 'stacktraces.New' not to be on call stack, got (%d, %s, %s, %d, %v)",
			frame,
			name,
			file,
			line,
			found)
	}
}

func TestFunctionInfoDepth(t *testing.T) {
	_, _, _, _, found := FunctionInfo(100)
	if found {
		t.Fatalf("expected found to return false")
	}
}

func TestFunctionName(t *testing.T) {
	actual := FunctionName()
	if !strings.HasSuffix(actual, ".TestFunctionName") {
		t.Fatalf("expected name to end with '.TestFunctionName', got '%s'", actual)
	}
}
