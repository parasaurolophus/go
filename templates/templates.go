package templates

import (
	"bytes"
	"context"
	"parasaurolophus/go/logging"
	"text/template"
)

type (
	TemplateParameters struct {
		Context   context.Context
		Logger    logging.Logger
		Name      string
		Text      string
		Arguments []any
		Handler   func(name string, text string, arguments []any, results []string)
	}
)

func ExecuteTemplate(ctx context.Context, logger logging.Logger, name string, text string, arguments ...any) (results []string) {
	results = make([]string, len(arguments))
	template, err := template.New(name).Parse(text)
	if err != nil {
		logError(ctx, logger, err)
		return
	}
	buffer := bytes.Buffer{}
	for index, argument := range arguments {
		buffer.Reset()
		err = template.Execute(&buffer, argument)
		if err != nil {
			logError(ctx, logger, err)
			continue
		}
		results[index] = buffer.String()
	}
	return
}

func ExecuteTemplates(parameters chan TemplateParameters, sync chan any) {
	defer close(sync)
	for parameter := range parameters {
		results := ExecuteTemplate(parameter.Context, parameter.Logger, parameter.Name, parameter.Text, parameter.Arguments...)
		parameter.Handler(parameter.Name, parameter.Text, parameter.Arguments, results)
	}
}

func logError(ctx context.Context, logger logging.Logger, err error) {
	logger.Always(
		ctx,
		func() string {
			return err.Error()
		},
		logging.TAGS,
		[]string{"templates.ExecuteTemplate", logging.ERROR},
	)
}
