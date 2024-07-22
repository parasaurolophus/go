// Copyright Kirk Rader 2024

package utilities

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestExecuteTemplate(t *testing.T) {
	type Datum struct {
		Foo int    `json:"foo"`
		Bar string `json:"bar"`
	}
	var arguments []Datum
	payload := `[{"foo":1,"bar":"one"},{"foo":2,"bar":"two"}]`
	err := json.Unmarshal([]byte(payload), &arguments)
	if err != nil {
		t.Fatal(err.Error())
	}
	text := `Foo is {{.Foo}}, Bar is "{{.Bar}}"`
	actual, err := ExecuteTemplate("", text, arguments...)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(actual) != 2 {
		t.Fatalf("expected 2, got %d", len(actual))
	}
	expected1 := `Foo is 1, Bar is "one"`
	if actual[0] != expected1 {
		t.Errorf(`expected "%s", got "%s"`, expected1, actual[0])
	}
	expected2 := `Foo is 2, Bar is "two"`
	if actual[1] != expected2 {
		t.Errorf(`expected "%s", got "%s"`, expected2, actual[1])
	}
	actual, err = ExecuteTemplate("", text, struct{ Baz float64 }{Baz: 4.2})
	fmt.Printf("\"%s\"\n", actual[0])
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestExecuteTemplates(t *testing.T) {
	successCount := 0
	errorCount := 0
	parameters := make(chan TemplateParameters)
	defer close(parameters)
	sync := make(chan bool)
	defer close(sync)
	testHandler := func(result string, err error) {
		if err != nil {
			errorCount++
			fmt.Fprintf(
				os.Stderr,
				"Result: \"%s\", Error: \"%s\"\n",
				result, err.Error(),
			)
			return
		}
		fmt.Printf("Result: \"%s\"\n", result)
		successCount++
	}
	data := []TemplateParameters{
		{
			Name: "",
			Text: `Foo is {{.Foo}}, Bar is "{{.Bar}}"`,
			Argument: struct {
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
			Argument: struct {
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
			Argument: struct {
				Answer float64
			}{
				Answer: 42.0,
			},
			Handler: testHandler,
		},
		{
			Name:     "",
			Text:     "",
			Argument: nil,
			Handler:  func(string, error) { sync <- true },
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
