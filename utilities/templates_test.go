// Copyright Kirk Rader 2024

package utilities

import (
	"encoding/json"
	"testing"
)

func TestExecuteTemplate(t *testing.T) {
	type Datum struct {
		Foo int    `json:"foo"`
		Bar string `json:"bar"`
	}
	var datum Datum
	payload := `{"foo":1,"bar":"one"}`
	err := json.Unmarshal([]byte(payload), &datum)
	if err != nil {
		t.Fatal(err.Error())
	}
	text := `Foo is {{.Foo}}, Bar is "{{.Bar}}"`
	actual, err := ExecuteTemplate("", text, datum)
	if err != nil {
		t.Error(err.Error())
	}
	expected := `Foo is 1, Bar is "one"`
	if actual != expected {
		t.Errorf(`expected "%s", got "%s"`, expected, actual)
	}
}

func TestExecuteTemplates(t *testing.T) {
	successCount := 0
	errorCount := 0
	parameters := make(chan TemplateParameters)
	defer close(parameters)
	sync := make(chan bool)
	defer close(sync)
	testHandler := func(_ string, err error) {
		if err != nil {
			errorCount++
			return
		}
		successCount++
	}
	data := []TemplateParameters{
		{
			Name: "",
			Text: `Foo is {{.Foo}}, Bar is "{{.Bar}}"`,
			Datum: struct {
				Foo int
				Bar string
			}{
				Foo: 1,
				Bar: "one",
			},
			Handler: testHandler,
		},
		{
			Name: "",
			Text: `Foo is {{.Foo}}, Bar is "{{.Bar}"`,
			Datum: struct {
				Foo int
				Bar string
			}{
				Foo: 2,
				Bar: "two",
			},
			Handler: testHandler,
		},
		{
			Name: "",
			Text: `The answer is {{.Answer}}`,
			Datum: struct {
				Answer float64
			}{
				Answer: 42.0,
			},
			Handler: testHandler,
		},
		{
			Name:    "",
			Text:    "",
			Datum:   nil,
			Handler: func(string, error) { sync <- true },
		},
	}
	go ExecuteTemplates(parameters)
	for _, param := range data {
		parameters <- param
	}
	<-sync
	if errorCount != 1 {
		t.Errorf("expected 1, got %d", errorCount)
	}
	if successCount != 2 {
		t.Errorf("expected 2, got %d", successCount)
	}
}
