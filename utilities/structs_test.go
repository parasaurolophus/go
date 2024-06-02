// Copyright Kirk Rader 2024

package utilities

import "testing"

func TestFilterStructFields(t *testing.T) {

	type A struct {
		InheritedField int
	}

	type B struct {
		A
		IntField     int
		StringField  string
		privateField float64
	}

	b := B{
		A: A{
			InheritedField: -1,
		},
		IntField:     42,
		StringField:  "forty-two",
		privateField: 4.2,
	}

	filteredValue, err := FilterStructFields(b, "InheritedField", "StringField", "privateField", "InvalidField")

	if err != nil {
		t.Fatal(err.Error())
	}

	if filteredValue == nil {
		t.Fatal("FilterStructFields return nil for a struct")
	}

	filtered := filteredValue.(B)

	if filtered.InheritedField != -1 {

		t.Errorf("expected -1, got %d", filtered.InheritedField)
	}

	if filtered.IntField != 0 {

		t.Errorf("expected 0, got %d", filtered.IntField)
	}

	if filtered.StringField != "forty-two" {

		t.Errorf("expected \"forty-two\", got \"%s\"", filtered.StringField)
	}

	if filtered.privateField != 0.0 {

		t.Errorf("expected 0, got %f", filtered.privateField)
	}

	filteredValue, err = FilterStructFields(42)

	if err == nil {
		t.Errorf("expected err not to be nil")
	}

	if filteredValue != nil {

		t.Error("expected filteredValue to be nil")
	}
}
