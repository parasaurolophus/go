package utilities

import (
	"embed"
	"encoding/csv"
	"testing"
)

//go:embed testdata
var testData embed.FS

func TestUnzip(t *testing.T) {
	archive, err := testData.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	entries, err := Unzip(archive)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		for _, entry := range entries {
			entry.Close()
		}
	}()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	csvReader := csv.NewReader(entries[0])
	sheet, err := csvReader.ReadAll()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(sheet) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(sheet))
	}
}
