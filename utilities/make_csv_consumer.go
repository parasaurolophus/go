// Copyright 2024 Kirk Rader

package utilities

import (
	"encoding/csv"
)

type (

	// Data required on consuper channel by functions created using
	// MakeCSVConsumer.
	//
	//
	// See ProcessBatch, MakeCSVGenerator, MakeCSVConsumer
	CSVConsumerParamters struct {
		CSVTransformerParameters
		Output map[string]string
	}
)

// Return a function for use as the consume parameter to ProcessBatch. The
// returned function will write each received row to the given CSV file. Any
// errors encountered along the way will be passed to the given errorHandler
// function.
//
// See ProcessBatch, MakeCSVGenerator
func MakeCSVConsumer(

	writer *csv.Writer,
	headers []string,
	errorHandler func(error),

) (

	consumer func(CSVConsumerParamters),

) {

	consumer = func(parameters CSVConsumerParamters) {
		columns := make([]string, len(parameters.Output))
		for i, h := range headers {
			columns[i] = parameters.Output[h]
		}
		e := writer.Write(columns)
		if e != nil {
			errorHandler(e)
		}
	}
	return
}
