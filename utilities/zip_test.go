package utilities

import (
	"encoding/csv"
	"parasaurolophus/go/common_test"
	"testing"
)

func TestUnzip(t *testing.T) {
	archive, err := common_test.TestData.Open("testdata/eurofxref.zip")
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
