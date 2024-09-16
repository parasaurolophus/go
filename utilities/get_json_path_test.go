// Copyright 2024 Kirk Rader

package utilities_test

import (
	"parasaurolophus/utilities"
	"testing"
)

func TestGetJSONPath(t *testing.T) {

	m := map[string]any{
		"foo": map[string]any{
			"bar": 42,
		},
	}

	// happy path
	if bar, err := utilities.GetJSONPath[int](m, "foo", "bar"); err != nil {
		t.Error(err.Error())
	} else if bar != 42 {
		t.Errorf("expected 42 but got %d", bar)
	}

	// empty path
	if _, err := utilities.GetJSONPath[float64](m); err == nil {
		t.Errorf("error expected")
	}

	// incorrect leaf type
	if _, err := utilities.GetJSONPath[float64](m, "foo", "bar"); err == nil {
		t.Errorf("error expected")
	}

	// incorrect intermediate type
	if _, err := utilities.GetJSONPath[float64](m, "foo", "bar", "baz"); err == nil {
		t.Errorf("error expected")
	}

	// missing key
	if _, err := utilities.GetJSONPath[float64](m, "foo", "baz"); err == nil {
		t.Errorf("error expected")
	}
}
