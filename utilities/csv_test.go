// Copyright 2024 Kirk Rader

package utilities_test

import (
	"bytes"
	"embed"
	"encoding/csv"
	"fmt"
	"parasaurolophus/utilities"
	"strconv"
	"testing"
)

//go:embed all:embedded
var embedded embed.FS

func TestProcessBatchCSVHappyPath(t *testing.T) {
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

func TestProcessBatchCSVInconsistent(t *testing.T) {
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

func TestProcessBatchCSVMalformed(t *testing.T) {
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
