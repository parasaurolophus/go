// Copyright Kirk Rader 2024

package utilities

import (
	"encoding/csv"
	"io"
)

type (

	// Type of function used to initialize
	CSVHeadersHandler func(headers []string) []string

	// Type of function used to process each row in a CSV file.
	CSVRowHandler func(headers, columns []string)
)

// Apply the given handler to each row in the given CSV file's contents.
func ForCSVReader(headersHandler CSVHeadersHandler, rowHandler CSVRowHandler, reader io.Reader) error {
	csv := csv.NewReader(reader)
	csv.TrimLeadingSpace = true
	headers, err := csv.Read()
	if err != nil {
		return err
	}
	if headersHandler != nil {
		headers = headersHandler(headers)
	}
	for {
		columns, err := csv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if rowHandler != nil {
			rowHandler(headers, columns)
		}
	}
	return nil
}
