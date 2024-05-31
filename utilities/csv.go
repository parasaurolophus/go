// Copyright Kirk Rader 2024

package utilities

import (
	"encoding/csv"
	"io"
)

type (

	// Type of function used to process the first row of a CSV file.
	CSVHeadersHandler func(headers []string) ([]string, error)

	// Type of function used to process each non-header row in a CSV file.
	CSVRowHandler func(headers, columns []string) error
)

// Apply the given handlers to each row in the given CSV file's contents. If
// headersHandler is non-nil, it will be applied to the first row of the CSV
// file and what it returns will be passed as the first argument to each
// subsequent invocation of rowHandler. If headersHandler is nil, the first
// argument to each invocation of rowHandler will be nil.
func ForCSVReader(headersHandler CSVHeadersHandler, rowHandler CSVRowHandler, reader io.Reader) (err error) {
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
	if rowHandler != nil {
		for {
			var columns []string
			columns, err = csv.Read()
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}
			err = rowHandler(headers, columns)
			if err != nil {
				break
			}
		}
	}
	return
}
