// Copyright Kirk Rader 2024

package utilities

import (
	"bytes"
	"text/template"
)

type (
	TemplateResultHandler func(string, error)

	TemplateParameters struct {
		Name    string
		Text    string
		Datum   any
		Handler TemplateResultHandler
	}

	TemplateParametersChannel chan TemplateParameters
)

func ExecuteTemplate(name, text string, datum any) (string, error) {
	tmplt, err := template.New(name).Parse(text)
	if err != nil {
		return "", err
	}
	buffer := bytes.Buffer{}
	err = tmplt.Execute(&buffer, datum)
	return buffer.String(), err
}

func ExecuteTemplates(channel TemplateParametersChannel) {
	for parameters := range channel {
		result, err := ExecuteTemplate(parameters.Name, parameters.Text, parameters.Datum)
		parameters.Handler(result, err)
	}
}
