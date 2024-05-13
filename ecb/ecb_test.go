// Copyright Kirk Rader 2024

package ecb

import (
	"embed"
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

func TestFetchDailyCSV(t *testing.T) {
	reader, err := Fetch(DailyCSV)
	if err != nil {
		t.Fatal(err.Error())
	}
	files, err := Unzip(reader)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(files) == 0 {
		t.Fatal("empty zip file")
	}
	for _, file := range files {
		data, err := ParseCSV(file)
		if err != nil {
			t.Error(err.Error())
		}
		if len(data) == 0 {
			t.Error("empty CSV file")
		}
	}
}

func TestFetchDailyXML(t *testing.T) {
	reader, err := Fetch(DailyXML)
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := ParseXML(reader)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty XML document")
	}
}

func TestParseDailyCSV(t *testing.T) {
	zipFile, err := dailyZip.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer zipFile.Close()
	csvFiles, err := Unzip(zipFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, csvFile := range csvFiles {
		data, err := ParseCSV(csvFile)
		if err != nil {
			t.Error(err.Error())
		}
		if len(data) < 1 {
			t.Errorf("no data returned")
		}
	}
}

func TestParseHistoricalCSV(t *testing.T) {
	zipFile, err := historicalZip.Open("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer zipFile.Close()
	csvFiles, err := Unzip(zipFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, csvFile := range csvFiles {
		data, err := ParseCSV(csvFile)
		if err != nil {
			t.Error(err.Error())
		}
		if len(data) < 1 {
			t.Errorf("no data returned")
		}
	}
}

func TestParseDailyXML(t *testing.T) {
	xmlFile, err := dailyXML.Open("testdata/eurofxref-daily.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer xmlFile.Close()
	data, err := ParseXML(xmlFile)
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
}

func TestParseHistoricalXML(t *testing.T) {
	xmlFile, err := historicalXml.Open("testdata/eurofxref-hist.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer xmlFile.Close()
	data, err := ParseXML(xmlFile)
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
}

func TestParseNinetyDayXML(t *testing.T) {
	xmlFile, err := ninetyDaysXML.Open("testdata/eurofxref-hist-90d.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer xmlFile.Close()
	data, err := ParseXML(xmlFile)
	if err != nil {
		t.Error(err.Error())
	}
	if len(data) == 0 {
		t.Error("empty response")
	}
}
