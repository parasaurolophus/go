// Copyright Kirk Rader 2024

package csv

import (
	"encoding/csv"
	"fmt"
	"io"
)

type (

	// Type of function used to process the first row of a CSV file.
	CSVHeadersHandler func(headers []string) ([]string, error)

	// Type of function used to process each subsequent row in a CSV file.
	CSVRowHandler func(headers, columns []string) error
)

// Apply the given handlers to each row in the given CSV file's contents. If
// headersHandler is non-nil, it will be applied to the first row of the CSV
// file and what it returns will be passed as the first argument to each
// subsequent invocation of rowHandler. If headersHandler is nil, the first
// argument to each invocation of rowHandler will be nil.
func ForEachCSVRow(headersHandler CSVHeadersHandler, rowHandler CSVRowHandler, reader io.Reader) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in CSV handler: %v", r)
		}
	}()
	csv := csv.NewReader(reader)
	csv.TrimLeadingSpace = true
	var headers []string
	if headersHandler != nil {
		headers, err = csv.Read()
		if err != nil {
			return
		}
		headers, err = headersHandler(headers)
		if err != nil {
			return
		}
	}
	for {
		var columns []string
		columns, err = csv.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		if rowHandler != nil {
			err = rowHandler(headers, columns)
			if err != nil {
				return
			}
		}
	}
}
