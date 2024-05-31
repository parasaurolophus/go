// Copyright Kirk Rader 2024

package utilities

import (
	"archive/zip"
	"parasaurolophus/go/common_test"
	"testing"
)

func TestForCSVReader(t *testing.T) {
	embedded, err := common_test.TestData.Open("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer embedded.Close()
	initialized := false
	zipCount := 0
	rowCount := 0
	headersHandler := func(headers []string) []string {
		initialized = true
		return headers
	}
	rowHandler := func(headers, columns []string) {
		rowCount += 1
	}
	zipHandler := func(entry *zip.File) {
		reader, zipErr := entry.Open()
		if zipErr != nil {
			t.Fatal(zipErr.Error())
		}
		defer reader.Close()
		zipCount += 1
		csvErr := ForCSVReader(headersHandler, rowHandler, reader)
		if csvErr != nil {
			t.Fatal(csvErr.Error())
		}
	}
	err = ForZipReader(zipHandler, embedded)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !initialized {
		t.Error(`not initialized`)
	}
	if zipCount != 1 {
		t.Errorf(`expected 1, got %d`, zipCount)
	}
	if rowCount < 2 {
		t.Errorf(`expected at least 2, got %d`, rowCount)
	}
}
