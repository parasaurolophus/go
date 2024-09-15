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
	if bar, err := utilities.GetJSONPathString[int]("foo/bar", m); err != nil {
		t.Error(err.Error())
	} else if bar != 42 {
		t.Errorf("expected 42 but got %d", bar)
	}

	// empty path
	if _, err := utilities.GetJSONPath[float64]([]string{}, m); err == nil {
		t.Errorf("error expected")
	}

	// incorrect leaf type
	if _, err := utilities.GetJSONPath[float64]([]string{"foo", "bar"}, m); err == nil {
		t.Errorf("error expected")
	}

	// incorrect intermediate type
	if _, err := utilities.GetJSONPath[float64]([]string{"foo", "bar", "baz"}, m); err == nil {
		t.Errorf("error expected")
	}

	// missing key
	if _, err := utilities.GetJSONPath[float64]([]string{"foo", "baz"}, m); err == nil {
		t.Errorf("error expected")
	}
}
