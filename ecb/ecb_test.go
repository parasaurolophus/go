// Copyright Kirk Rader 2024

package ecb

import (
	"embed"
	"io"
	"testing"
)

//go:embed testdata/eurofxref.zip
var dailyZip embed.FS

//go:embed testdata/eurofxref-daily.xml
var dailyXML embed.FS

//go:embed testdata/eurofxref-hist.xml
var historicalXml embed.FS

//go:embed testdata/eurofxref-hist.zip
var historicalZip embed.FS

//go:embed testdata/eurofxref-hist-90d.xml
var ninetyDaysXML embed.FS

func TestFetchBadURL(t *testing.T) {
	_, err := Fetch("bad", ParseCSV)
	if err == nil {
		t.Error("expected error when Fetch() was passed an ill-formed URL")
	}
	_, err = Fetch("http://google.com", ParseXML)
	if err == nil {
		t.Error("expected error when Fetch() was passed a URL that responds with non-XML data")
	}
}

func TestFetchDailyCSV(t *testing.T) {
	data, err := Fetch(DailyCSV, ParseCSV)
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
}

func TestFetchDailyXML(t *testing.T) {
	data, err := Fetch(DailyXML, ParseXML)
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
}

func TestParseDailyCSV(t *testing.T) {
	source, err := dailyZip.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer source.Close()
	data, err := ParseCSV(source)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(data) < 1 {
		t.Errorf("no data returned")
	}
}

func TestParseHistoricalCSV(t *testing.T) {
	source, err := historicalZip.Open("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer source.Close()
	data, err := ParseCSV(source)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(data) < 1 {
		t.Errorf("no data returned")
	}
}

func TestParseDailyXML(t *testing.T) {
	source, err := dailyXML.Open("testdata/eurofxref-daily.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer source.Close()
	data, err := ParseXML(source.(io.ReadCloser))
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
}

func TestParseHistoricalXML(t *testing.T) {
	source, err := historicalXml.Open("testdata/eurofxref-hist.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer source.Close()
	data, err := ParseXML(source.(io.ReadCloser))
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
}

func TestParseNinetyDayXML(t *testing.T) {
	source, err := ninetyDaysXML.Open("testdata/eurofxref-hist-90d.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer source.Close()
	data, err := ParseXML(source.(io.Reader))
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
}
