// Copyright Kirk Rader 2024

package csv

import (
	"bytes"
	"fmt"
	"testing"
)

func TestForEachCSVRow(t *testing.T) {
	reader := bytes.NewReader([]byte("foo,bar\n1,3\n2,4"))
	initialized := false
	rowCount := 0
	headersHandler := func(headers []string) ([]string, error) {
		initialized = true
		return headers, nil
	}
	rowHandler := func(headers, columns []string) error {
		rowCount += 1
		return nil
	}
	err := ForEachCSVRow(headersHandler, rowHandler, reader)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !initialized {
		t.Error(`not initialized`)
	}
	if rowCount < 2 {
		t.Errorf(`expected at least 2, got %d`, rowCount)
	}
}

func TestForEachCSVRowPanic(t *testing.T) {
	reader := bytes.NewReader([]byte("foo,bar\n1,3\n2,4"))
	headersHandler := func(headers []string) ([]string, error) {
		return headers, nil
	}
	rowHandler := func(headers, columns []string) error {
		panic("deliberate")
	}
	err := ForEachCSVRow(headersHandler, rowHandler, reader)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if err.Error() != "panic in CSV handler: deliberate" {
		t.Errorf(`expected "panic in CSV handler: deliberate", got "%s"`, err.Error())
	}
}

func TestForEachCSVRowHeadersError(t *testing.T) {
	reader := bytes.NewReader([]byte("foo,bar\n1,3\n2,4"))
	headersHandler := func(headers []string) ([]string, error) {
		return nil, fmt.Errorf("deliberate")
	}
	rowHandler := func(headers, columns []string) error {
		return nil
	}
	err := ForEachCSVRow(headersHandler, rowHandler, reader)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if err.Error() != "deliberate" {
		t.Errorf(`expected "deliberate", got "%s"`, err.Error())
	}
}

func TestForEachCSVRowRowError(t *testing.T) {
	reader := bytes.NewReader([]byte("foo,bar\n1,3\n2,4"))
	headersHandler := func(headers []string) ([]string, error) {
		return headers, nil
	}
	rowHandler := func(headers, columns []string) error {
		return fmt.Errorf("deliberate")
	}
	err := ForEachCSVRow(headersHandler, rowHandler, reader)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if err.Error() != "deliberate" {
		t.Errorf(`expected "deliberate", got "%s"`, err.Error())
	}
}

func TestForEachCSVRowRowEOF(t *testing.T) {
	reader := bytes.NewReader([]byte(""))
	headersHandler := func(headers []string) ([]string, error) {
		return headers, nil
	}
	rowHandler := func(headers, columns []string) error {
		return nil
	}
	err := ForEachCSVRow(headersHandler, rowHandler, reader)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if err.Error() != "EOF" {
		t.Errorf(`expected "EOF", got "%s"`, err.Error())
	}
}
