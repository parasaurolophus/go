package templates_test

import (
	"bytes"
	"context"
	"parasaurolophus/go/logging"
	"parasaurolophus/go/templates"
	"testing"
)

func TestExecuteTemplates(t *testing.T) {
	loggerOutput := bytes.Buffer{}
	loggerOptions := logging.LoggerOptions{}
	logger := logging.New(&loggerOutput, &loggerOptions)
	type (
		Test1 struct {
			Foo float64
			Bar string
		}
		Test2 struct {
			Answer int
		}
	)
	arguments1 := []any{
		struct {
			SomeThingWicked string
		}{
			SomeThingWicked: "this way comes",
		},
		Test1{
			Foo: 4.2,
			Bar: "not the answer",
		},
	}
	arguments2 := []any{
		Test2{
			Answer: 42,
		},
	}
	arguments3 := []any{
		Test2{},
	}
	sync := make(chan any)
	func() {
		handler1 := func(name string, text string, arguments []any, results []string) {
			if len(arguments1) != len(results) {
				t.Errorf(`expected %d results, got %d`, len(arguments1), len(results))
			}
			for index, result := range results {
				switch index {
				case 0:
					if result != "" {
						t.Errorf(`expected result %d to be empty, got "%s"`, index, result)
					}
				case 1:
					if result != "Foo is 4.2, Bar is not the answer." {
						t.Errorf(`expected "Foo is 4.2, Bar is not the answer.", got "%s"`, result)
					}
				default:
					t.Errorf(`unexpected result %d`, index)
				}
			}
		}
		handler2 := func(name string, text string, arguments []any, results []string) {
			if len(arguments2) != len(results) {
				t.Errorf(`expected %d results, got %d`, len(arguments2), len(results))
			}
			for index, result := range results {
				switch index {
				case 0:
					if result != "The answer is 42." {
						t.Errorf(`expected "The answer is 42.", got "%s"`, result)
					}
				default:
					t.Errorf(`unexpected result %d`, index)
				}
			}
		}
		handler3 := func(name string, text string, arguments []any, results []string) {
			if len(arguments2) != len(results) {
				t.Errorf(`expected %d results, got %d`, len(arguments2), len(results))
			}
			for index, result := range results {
				switch index {
				case 0:
					if result != "" {
						t.Errorf(`expected result %d to be empty, got "%s"`, index, result)
					}
				default:
					t.Errorf(`unexpected result %d`, index)
				}
			}
		}
		parameters := make(chan templates.TemplateParameters)
		defer close(parameters)
		go templates.ExecuteTemplates(parameters, sync)
		parameters <- templates.TemplateParameters{
			Context:   context.Background(),
			Logger:    logger,
			Name:      "",
			Text:      `Foo is {{.Foo}}, Bar is {{.Bar}}.`,
			Arguments: arguments1,
			Handler:   handler1,
		}
		parameters <- templates.TemplateParameters{
			Context:   context.Background(),
			Logger:    logger,
			Name:      "",
			Text:      `The answer is {{.Answer}}.`,
			Arguments: arguments2,
			Handler:   handler2,
		}
		parameters <- templates.TemplateParameters{
			Context:   context.Background(),
			Logger:    logger,
			Name:      "",
			Text:      `The answer is {{.Answer}.`,
			Arguments: arguments3,
			Handler:   handler3,
		}
	}()
	<-sync
	if loggerOutput.String() == "" {
		t.Error(`expected errors to have been logged`)
	}
}
