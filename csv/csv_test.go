// Copyright Kirk Rader 2024

package csv

import (
	"archive/zip"
	"fmt"
	"parasaurolophus/go/common_test"
	z "parasaurolophus/go/zip"
	"testing"
)

func TestForEachCSVRow(t *testing.T) {
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
		err = ForEachCSVRow(headersHandler, rowHandler, reader)
		return
	}
	err = z.ForEachZipEntryFromReader(zipHandler, embedded)
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

func TestForEachCSVRowPanic(t *testing.T) {
	embedded, err := common_test.TestData.Open("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer embedded.Close()
	headersHandler := func(headers []string) ([]string, error) {
		return headers, nil
	}
	rowHandler := func(headers, columns []string) error {
		panic("deliberate")
	}
	zipHandler := func(entry *zip.File) (err error) {
		reader, err := entry.Open()
		if err != nil {
			return
		}
		defer reader.Close()
		err = ForEachCSVRow(headersHandler, rowHandler, reader)
		return
	}
	err = z.ForEachZipEntryFromReader(zipHandler, embedded)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if err.Error() != "panic in CSV handler: deliberate" {
		t.Errorf(`expected "panic in CSV handler: deliberate", got "%s"`, err.Error())
	}
}

func TestForEachCSVRowHeadersError(t *testing.T) {
	embedded, err := common_test.TestData.Open("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer embedded.Close()
	headersHandler := func(headers []string) ([]string, error) {
		return nil, fmt.Errorf("deliberate")
	}
	rowHandler := func(headers, columns []string) error {
		return nil
	}
	zipHandler := func(entry *zip.File) (err error) {
		reader, err := entry.Open()
		if err != nil {
			return
		}
		defer reader.Close()
		err = ForEachCSVRow(headersHandler, rowHandler, reader)
		return
	}
	err = z.ForEachZipEntryFromReader(zipHandler, embedded)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if err.Error() != "deliberate" {
		t.Errorf(`expected "deliberate", got "%s"`, err.Error())
	}
}

func TestForEachCSVRowRowError(t *testing.T) {
	embedded, err := common_test.TestData.Open("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer embedded.Close()
	headersHandler := func(headers []string) ([]string, error) {
		return headers, nil
	}
	rowHandler := func(headers, columns []string) error {
		return fmt.Errorf("deliberate")
	}
	zipHandler := func(entry *zip.File) (err error) {
		reader, err := entry.Open()
		if err != nil {
			return
		}
		defer reader.Close()
		err = ForEachCSVRow(headersHandler, rowHandler, reader)
		return
	}
	err = z.ForEachZipEntryFromReader(zipHandler, embedded)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if err.Error() != "deliberate" {
		t.Errorf(`expected "deliberate", got "%s"`, err.Error())
	}
}
