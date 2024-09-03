// Copyright 2024 Kirk Rader

package utilities

import (
	"encoding/csv"
	"io"
)

type (

	// Data sent to transformers channel by functions created using
	// MakeCSVGenerator.
	//
	// See ProcessBatch, MakeCSVGenerator, MakeCSVConsumer
	CSVTransformerParameters struct {
		Row   int
		Input map[string]string
	}
)

// Return a function for use as the generate parameter to ProcessBatch. The
// returned function reads the rows of the given CSV file, sending each to the
// batch's transformers channels in round-robin fashion. Any errors encountered
// along the way will be passed to the given errorHandler function. Note that
// the column headers and starting row number are passed in here so as to
// support CSV's without a headers row.
//
// See ProcessBatch, MakeCSVConsumer
func MakeCSVGenerator(

	reader *csv.Reader,
	headers []string,
	startRow int,
	errorHandler func(error),

) (

	generator func([]chan<- CSVTransformerParameters),
	err error,

) {

	generator = func(transformers []chan<- CSVTransformerParameters) {
		row := startRow
		n := len(transformers)
		for {
			columns, e := reader.Read()
			if e != nil {
				if e != io.EOF {
					errorHandler(e)
				}
				break
			}
			parameters := CSVTransformerParameters{
				Row:   row,
				Input: map[string]string{},
			}
			for i, h := range headers {
				parameters.Input[h] = columns[i]
			}
			transformers[(row-startRow)%n] <- parameters
			row++
		}
	}
	return
}
