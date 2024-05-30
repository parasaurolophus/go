// Copyright Kirk Rader 2024

package ecb

import (
	"archive/zip"
	"parasaurolophus/go/common_test"
	"parasaurolophus/go/utilities"
	"testing"
)

func TestFetchDailyCSV(t *testing.T) {
	readCloser, err := utilities.Fetch(DAILY_CSV_URL)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer readCloser.Close()
	handler := func(entry *zip.File) {
		file, err := entry.Open()
		if err != nil {
			t.Fatal(err.Error())
		}
		data, err := ParseCSV(file)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(data) == 0 {
			t.Error("empty CSV file")
		}
	}
	err = utilities.ForZipReader(handler, readCloser)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFetchDailyXML(t *testing.T) {
	readCloser, err := utilities.Fetch(DAILY_XML_URL)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer readCloser.Close()
	data, err := ParseXML(readCloser)
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
	handler := func(entry *zip.File) {
		file, err := entry.Open()
		if err != nil {
			t.Fatal(err.Error())
		}
		data, err := ParseCSV(file)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(data) == 0 {
			t.Error("empty CSV file")
		}
	}
	err = utilities.ForZipReader(handler, zipFile)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestParseHistoricalCSV(t *testing.T) {
	zipFile, err := common_test.TestData.Open("testdata/eurofxref-hist.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer zipFile.Close()
	handler := func(entry *zip.File) {
		file, err := entry.Open()
		if err != nil {
			t.Fatal(err.Error())
		}
		data, err := ParseCSV(file)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(data) == 0 {
			t.Error("empty CSV file")
		}
	}
	err = utilities.ForZipReader(handler, zipFile)
	if err != nil {
		t.Fatal(err.Error())
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
