// Copyright Kirk Rader 2024

package stacktraces

import (
	"strings"
	"testing"
)

func TestFunctionInfo(t *testing.T) {

	const expectedName = "parasaurolophus/go/stacktraces.TestFunctionInfo"
	actualName, actualFile, actuaLine := FunctionInfo()

	if actualName != expectedName {
		t.Fatalf("expected name to be '%s', got '%s'", expectedName, actualName)
	}

	if !strings.HasSuffix(actualFile, "/function_info_test.go") {
		t.Fatalf("expected file name to be 'function_info_test.go', got '%s'", actualFile)
	}

	if actuaLine != 13 {
		t.Fatalf("expected line to be 13, got %d", actuaLine)
	}
}

func TestFunctionName(t *testing.T) {

	const expectedName = "parasaurolophus/go/stacktraces.TestFunctionName"
	actualName := FunctionName()

	if actualName != expectedName {
		t.Fatalf("expected name to be '%s', got '%s'", expectedName, actualName)
	}
}
