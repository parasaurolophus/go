// Copyright Kirk Rader 2024

package utilities

import "testing"

func TestFilterMapKeys(t *testing.T) {

	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	filtered := FilterMapKeys(m, "one", "three")

	if len(filtered) != 2 {
		t.Errorf("expected 2, got %d", len(filtered))
	}

	if filtered["one"] != 1 {
		t.Errorf("expected 1, got %d", filtered["one"])
	}

	if filtered["three"] != 3 {
		t.Errorf("expected 3, got %d", filtered["three"])
	}
}

func TestMergeMaps(t *testing.T) {

	map1 := map[string]int{
		"one": 1,
		"two": 2,
	}

	map2 := map[string]int{
		"three": 3,
		"four":  4,
	}

	merged := MergeMaps(map1, map2)

	if len(merged) != 4 {
		t.Errorf("expected 4, got %d", len(merged))
	}

	if merged["one"] != 1 {
		t.Errorf("expected 1, got %d", merged["one"])
	}

	if merged["two"] != 2 {
		t.Errorf("expected 2, got %d", merged["two"])
	}

	if merged["three"] != 3 {
		t.Errorf("expected 3, got %d", merged["three"])
	}

	if merged["four"] != 4 {
		t.Errorf("expected 4, got %d", merged["four"])
	}
}
