// Copyright Kirk Rader 2024

package ecb

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

//go:embed eurofxref-hist-90d.xml
var ninetyDays string

//go:embed eurofxref-daily.xml
var daily string

func TestFetchNinetyDays(t *testing.T) {
	data, err := Fetch(NinetyDayURL)
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	writer.Flush()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(buffer.String())
}

func TestFetchDaily(t *testing.T) {
	data, err := Fetch(DailyURL)
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	writer.Flush()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(buffer.String())
}

func TestParseNinetyDayString(t *testing.T) {
	data, err := Parse(strings.NewReader(ninetyDays))
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	writer.Flush()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(buffer.String())
}

func TestParseOneDayString(t *testing.T) {
	data, err := Parse(strings.NewReader(daily))
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(data)
	writer.Flush()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(buffer.String())
}
