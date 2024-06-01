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
	headersHandler := func(headers []string) ([]string, error) {
		initialized = true
		return headers, nil
	}
	rowHandler := func(headers, columns []string) error {
		rowCount += 1
		return nil
	}
	zipHandler := func(entry *zip.File) (err error) {
		reader, err := entry.Open()
		if err != nil {
			return
		}
		defer reader.Close()
		zipCount += 1
		err = ForCSVReader(headersHandler, rowHandler, reader)
		return
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
