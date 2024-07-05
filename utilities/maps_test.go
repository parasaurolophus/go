// Copyright Kirk Rader 2024

package utilities

import (
	"testing"

	"github.com/google/uuid"
)

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

func TestGetRequiredFloat(t *testing.T) {
	m := map[string]any{
		"field1": 4.2,
		"field2": "value2",
	}
	f, err := GetRequiredFloat("field1", m)
	if err != nil {
		t.Error(err.Error())
	}
	if f != 4.2 {
		t.Errorf("expected 4.2, got %g", f)
	}
	_, err = GetRequiredFloat("field2", m)
	if err == nil {
		t.Error("expected err not to be nil")
	}
	_, err = GetRequiredFloat("missing", m)
	if err == nil {
		t.Error("expected err not to be nil")
	}
}

func TestGetRequiredString(t *testing.T) {
	m := map[string]any{
		"field1": "value1",
		"field2": 42,
	}
	s, err := GetRequiredString("field1", m)
	if err != nil {
		t.Error(err.Error())
	}
	if s != "value1" {
		t.Errorf(`expected "value1", got "%s"`, s)
	}
	_, err = GetRequiredString("field2", m)
	if err == nil {
		t.Error("expected err not to be nil")
	}
	_, err = GetRequiredString("missing", m)
	if err == nil {
		t.Error("expected err not to be nil")
	}
}

type uuidString struct {
	id uuid.UUID
}

func (u uuidString) String() string {
	return u.id.String()
}

func TestGetRequiredUUID(t *testing.T) {
	expected := uuid.New()
	m := map[string]any{
		"field1": expected,
		"field2": expected.String(),
		"field3": uuidString{id: expected},
		"field4": 42,
	}
	actual, err := GetRequiredUUID("field1", m)
	if err != nil {
		t.Error(err.Error())
	}
	if actual != expected {
		t.Errorf(`expected "%v", got "%v"`, expected, actual)
	}
	actual, err = GetRequiredUUID("field2", m)
	if err != nil {
		t.Error(err.Error())
	}
	if actual != expected {
		t.Errorf(`expected "%v", got "%v"`, expected, actual)
	}
	actual, err = GetRequiredUUID("field3", m)
	if err != nil {
		t.Error(err.Error())
	}
	if actual != expected {
		t.Errorf(`expected "%v", got "%v"`, expected, actual)
	}
	_, err = GetRequiredUUID("field4", m)
	if err == nil {
		t.Error("expected err not to be nil")
	}
	_, err = GetRequiredUUID("missing", m)
	if err == nil {
		t.Error("expected err not to be nil")
	}
}

func TestGetValue(t *testing.T) {
	m := map[int]any{
		-1: "value -1",
		2:  4.2,
	}
	s, found, err := GetValue[int, string](-1, m)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !found {
		t.Fatalf("expected value for -1 to be found")
	}
	if s != "value -1" {
		t.Errorf(`expected "value -1", got "%s"`, s)
	}
	f, found, err := GetValue[int, float64](2, m)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !found {
		t.Fatalf("expected value for 2 to be found")
	}
	if f != 4.2 {
		t.Errorf("expected 4.2, got %g", f)
	}
	_, found, err = GetValue[int, string](0, m)
	if err != nil {
		t.Fatal(err.Error())
	}
	if found {
		t.Errorf("expected value for 0 not to be found")
	}
	_, _, err = GetValue[int, string](2, m)
	if err == nil {
		t.Errorf("expected err not to be nil")
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
