// Copyright Kirk Rader 2024

package ecb

import (
	"parasaurolophus/go/common_test"
	"parasaurolophus/go/utilities"
	"testing"
)

func TestFetchDailyCSV(t *testing.T) {
	reader, err := utilities.Fetch(DAILY_CSV_URL)
	if err != nil {
		t.Fatal(err.Error())
	}
	files, err := utilities.Unzip(reader)
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
	reader, err := utilities.Fetch(DAILY_XML_URL)
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
	zipFile, err := common_test.TestData.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer zipFile.Close()
	csvFiles, err := utilities.Unzip(zipFile)
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
	zipFile, err := common_test.TestData.Open("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer zipFile.Close()
	csvFiles, err := utilities.Unzip(zipFile)
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
	xmlFile, err := common_test.TestData.Open("testdata/eurofxref-daily.xml")
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
	xmlFile, err := common_test.TestData.Open("testdata/eurofxref-hist.xml")
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
	xmlFile, err := common_test.TestData.Open("testdata/eurofxref-hist-90d.xml")
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
