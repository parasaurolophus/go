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
	inputData, err := embedded.ReadFile("embedded/input.csv")
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(inputData)
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
	generate, err := utilities.MakeCSVGenerator(csvReader, headers, 1, errorHandler)
	if err != nil {
		t.Fatal(err)
	}
	transform := func(input utilities.CSVTransformerParameters) (output utilities.CSVConsumerParamters) {
		output = utilities.CSVConsumerParamters{}
		output.Row = input.Row
		output.Input = input.Input
		// just copying input to output for this unit test
		output.Output = input.Input
		col, ok := input.Input["number"]
		if !ok {
			errorHandler(fmt.Errorf("%v has no value for \"number\"", input))
			return
		}
		n, err := strconv.Atoi(col)
		if err != nil {
			errorHandler(err)
			return
		}
		actual += n
		return
	}
	buffer := bytes.Buffer{}
	writer := csv.NewWriter(&buffer)
	err = writer.Write(headers)
	if err != nil {
		t.Fatal(err)
	}
	consume := utilities.MakeCSVConsumer(writer, headers, errorHandler)
	func() {
		defer writer.Flush()
		utilities.ProcessBatch(3, 1, 1, generate, transform, consume)
	}()
	if actual != 15 {
		t.Errorf("expected 15, got %d", actual)
	}
	if errors != 0 {
		t.Errorf("expected no errors, got %d", errors)
	}
	outputData := buffer.Bytes()
	inputLen := len(inputData)
	outputLen := len(outputData)
	if inputLen != outputLen {
		t.Errorf(
			"expected CSV input and output to be of same size; input is %d, output is %d",
			inputLen,
			outputLen,
		)
	}
	fmt.Println(string(outputData))
}

func TestProcessBatchInconsistent(t *testing.T) {
	b, err := embedded.ReadFile("embedded/inconsistent.csv")
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
	errorHandler := func(err error) {
		errors += 1
	}
	generate, err := utilities.MakeCSVGenerator(csvReader, headers, 1, errorHandler)
	if err != nil {
		t.Fatal(err)
	}
	transform := func(input utilities.CSVTransformerParameters) (output utilities.CSVConsumerParamters) {
		output = utilities.CSVConsumerParamters{}
		output.Row = input.Row
		output.Input = input.Input
		output.Output = input.Input
		col, ok := input.Input["number"]
		if !ok {
			errorHandler(fmt.Errorf("%v has no value for \"number\"", input))
			return
		}
		n, e := strconv.Atoi(col)
		if e != nil {
			errorHandler(e)
			return
		}
		actual += n
		return
	}
	buffer := bytes.Buffer{}
	writer := csv.NewWriter(&buffer)
	err = writer.Write(headers)
	if err != nil {
		t.Fatal(err)
	}
	consume := utilities.MakeCSVConsumer(writer, headers, errorHandler)
	func() {
		defer writer.Flush()
		utilities.ProcessBatch(3, 1, 1, generate, transform, consume)
	}()
	if actual != 13 {
		t.Errorf("expected 13, got %d", actual)
	}
	if errors != 1 {
		t.Errorf("expected 1 error, got %d", errors)
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
	generate, err := utilities.MakeCSVGenerator(csvReader, headers, 1, errorHandler)
	if err != nil {
		t.Fatal(err)
	}
	transform := func(input utilities.CSVTransformerParameters) (output utilities.CSVConsumerParamters) {
		output = utilities.CSVConsumerParamters{}
		output.Row = input.Row
		output.Input = input.Input
		output.Output = input.Input
		col, ok := input.Input["number"]
		if !ok {
			errorHandler(fmt.Errorf("%v has no value for \"number\"", input))
			return
		}
		n, e := strconv.Atoi(col)
		if e != nil {
			errorHandler(e)
			return
		}
		actual += n
		return
	}
	buffer := bytes.Buffer{}
	writer := csv.NewWriter(&buffer)
	err = writer.Write(headers)
	if err != nil {
		t.Fatal(err)
	}
	consume := utilities.MakeCSVConsumer(writer, headers, errorHandler)
	func() {
		defer writer.Flush()
		utilities.ProcessBatch(3, 1, 1, generate, transform, consume)
	}()
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
		values, await := utilities.StartWorkers(3, 1, handler)
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
