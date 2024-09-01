// Copyright 2024 Kirk Rader

package utilities_test

import (
	"bytes"
	"embed"
	"encoding/csv"
	"fmt"
	"parasaurolophus/utilities"
	"strconv"
	"sync"
	"testing"
	"time"
)

//go:embed all:embedded
var embedded embed.FS

func TestProcessBatchHappyPath(t *testing.T) {
	b, err := embedded.ReadFile("embedded/input.csv")
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(b)
	csvReader := csv.NewReader(reader)
	headers, err := csvReader.Read()
	if err != nil {
		t.Fatal(err)
	}
	errors := 0
	actual := 0
	errorHandler := func(error) {
		errors += 1
	}
	generate, err := utilities.MakeCSVGenerator(csvReader, headers, errorHandler)
	if err != nil {
		t.Fatal(err)
	}
	transform := func(input map[string]string) (n int) {
		col, ok := input["number"]
		if !ok {
			errorHandler(fmt.Errorf("%v has no value for \"number\"", input))
			return
		}
		var err error
		n, err = strconv.Atoi(col)
		if err != nil {
			errorHandler(err)
		}
		return
	}
	consume := func(output int) {
		actual += output
	}
	utilities.ProcessBatch(3, generate, transform, consume)
	if actual != 15 {
		t.Errorf("expected 15, got %d", actual)
	}
	if errors != 0 {
		t.Errorf("expected no errors, got %d", errors)
	}
}

func TestProcessBatchMalformed(t *testing.T) {
	b, err := embedded.ReadFile("embedded/malformed.csv")
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(b)
	csvReader := csv.NewReader(reader)
	headers, err := csvReader.Read()
	if err != nil {
		t.Fatal(err)
	}
	errors := 0
	actual := 0
	errorHandler := func(error) {
		errors += 1
	}
	generate, err := utilities.MakeCSVGenerator(csvReader, headers, errorHandler)
	if err != nil {
		t.Fatal(err)
	}
	transform := func(input map[string]string) (n int) {
		col, ok := input["number"]
		if !ok {
			errorHandler(fmt.Errorf("%v has no value for \"number\"", input))
			return
		}
		var err error
		n, err = strconv.Atoi(col)
		if err != nil {
			errorHandler(err)
		}
		return
	}
	consume := func(output int) {
		actual += output
	}
	utilities.ProcessBatch(3, generate, transform, consume)
	if actual != 1 {
		t.Errorf("expected 1, got %d", actual)
	}
	if errors != 1 {
		t.Errorf("expected 1 error, got %d", errors)
	}
}

func TestStartWorkers(t *testing.T) {
	actual := 0
	lock := sync.Mutex{}
	handler := func(n int) {
		defer lock.Unlock()
		lock.Lock()
		actual += n
	}
	func() {
		values, await := utilities.StartWorkers(3, handler)
		defer utilities.CloseAllAndWait(values, await)
		for i := range len(values) {
			values[i] <- i
		}
	}()
	if actual != 3 {
		t.Errorf("expected 3, got %d", actual)
	}
}

func TestWithTimeLimit(t *testing.T) {
	fn1 := func() int {
		return 1
	}
	fn2 := func() int {
		time.Sleep(time.Millisecond * 100)
		return 2
	}
	timeout := func(time.Time) int {
		return 0
	}
	if v := utilities.WithTimeLimit(fn1, timeout, time.Millisecond*50); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
	if v := utilities.WithTimeLimit(fn2, timeout, time.Millisecond*50); v != 0 {
		t.Errorf("expected 0, got %d", v)
	}
}
