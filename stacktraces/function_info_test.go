// Copyright Kirk Rader 2024

package stacktraces

import (
	"strings"
	"testing"
)

func TestFunctionInfoNil(t *testing.T) {
	frame, sourceInfo, found := FunctionInfo(nil)
	if !found {
		t.Fatalf("function info for nil returned false as its last value")
	}
	if frame != defaultSkip {
		t.Fatalf("expected frame to be %d, got %d", defaultSkip, frame)
	}
	if !strings.HasSuffix(sourceInfo.Function, ".TestFunctionInfoNil") {
		t.Fatalf("expected name to end with '.TestFunctionInfoNil', got '%s'", sourceInfo.Function)
	}
	if !strings.HasSuffix(sourceInfo.File, "/function_info_test.go") {
		t.Fatalf("expected file name to be 'function_info_test.go', got '%s'", sourceInfo.File)
	}
	if sourceInfo.Line != 11 {
		t.Fatalf("expected line to be 11, got %d", sourceInfo.Line)
	}
}

func TestFunctionInfoNegative(t *testing.T) {
	frame, _, found := FunctionInfo(-1)
	if !found {
		t.Fatalf("function info for nil returned false as its last value")
	}
	if frame != defaultSkip+1 {
		t.Fatalf("expected frame to be %d, got %d", defaultSkip+1, frame)
	}
}

func TestFunctionInfoPositive(t *testing.T) {
	_, sourceInfo, found := FunctionInfo(1)
	if !found {
		t.Fatalf("function info for nil returned false as its last value")
	}
	if !strings.HasSuffix(sourceInfo.Function, ".functionInfo") {
		t.Fatalf("expected name to end with '.functionInfo', got '%s'", sourceInfo.Function)
	}
	if !strings.HasSuffix(sourceInfo.File, "/function_info.go") {
		t.Fatalf("expected file name to be 'function_info.go', got '%s'", sourceInfo.File)
	}
}

func TestFunctionInfoString(t *testing.T) {
	expected := FunctionName()
	frame, sourceInfo, found := FunctionInfo(expected)
	if !found {
		t.Fatalf("function info for nil returned false as its last value")
	}
	if frame != defaultSkip {
		t.Fatalf("expected frame to be %d, got %d", defaultSkip, frame)
	}
	if sourceInfo.Function != expected {
		t.Fatalf("expected name to be '%s', got '%s'", expected, sourceInfo.Function)
	}
	if !strings.HasSuffix(sourceInfo.File, "/function_info_test.go") {
		t.Fatalf("expected file name to be 'function_info_test.go', got '%s'", sourceInfo.File)
	}
	if sourceInfo.Line != 54 {
		t.Fatalf("expected line to be 54, got %d", sourceInfo.Line)
	}
}

func TestFunctionInfoStringNotFound(t *testing.T) {
	frame, sourceInfo, found := FunctionInfo("stacktraces.New")
	if found {
		t.Fatalf(
			"expected 'stacktraces.New' not to be on call stack, got %v at frame %d",
			sourceInfo,
			frame)
	}
}

func TestFunctionInfoDepth(t *testing.T) {
	_, _, found := FunctionInfo(100)
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
