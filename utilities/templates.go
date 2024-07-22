// Copyright Kirk Rader 2024

package utilities

import (
	"bytes"
	"text/template"
)

type (
	TemplateResultHandler func(string, error)

	TemplateParameters struct {
		Name     string
		Text     string
		Argument any
		Handler  TemplateResultHandler
	}

	TemplateParametersChannel chan TemplateParameters
)

func ExecuteTemplate[T any](name, text string, arguments ...T) (results []string, err error) {
	results = make([]string, len(arguments))
	tmplt, err := template.New(name).Parse(text)
	if err != nil {
		return
	}
	buffer := bytes.Buffer{}
	for index, argument := range arguments {
		buffer.Reset()
		err = tmplt.Execute(&buffer, argument)
		if err != nil {
			return
		}
		results[index] = buffer.String()
	}
	return
}

func ExecuteTemplates(channel TemplateParametersChannel) {
	for parameters := range channel {
		results, err := ExecuteTemplate(parameters.Name, parameters.Text, parameters.Argument)
		parameters.Handler(results[0], err)
	}
}
