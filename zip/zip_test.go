// Copyright Kirk Rader 2024

package zip

import (
	"archive/zip"
	"io"
	"os"
	"parasaurolophus/go/common_test"
	"strings"
	"testing"
)

func TestForEachZipEntryFromFile(t *testing.T) {
	embedded, err := common_test.TestData.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer embedded.Close()
	file, err := os.CreateTemp(t.TempDir(), "TestForZipFile")
	if err != nil {
		t.Fatal(err.Error())
	}
	// deferred functions are invoked in reverse order
	defer os.Remove(file.Name()) // invoked second
	defer file.Close()           // invoked first
	_, err = io.Copy(file, embedded)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	entryCount := 0
	totalSize := 0
	handler := func(entry *zip.File) (err error) {
		if strings.HasSuffix(entry.Name, "/") {
			t.Fatalf("unsupported entry type %s", entry.Name)
		}
		reader, err := entry.Open()
		if err != nil {
			return
		}
		defer reader.Close()
		entryCount += 1
		b, err := io.ReadAll(reader)
		if err != nil {
			return
		}
		totalSize += len(b)
		return
	}
	err = ForEachZipEntryFromFile(handler, file)
	if err != nil {
		t.Fatal(err.Error())
	}
	if entryCount != 1 {
		t.Fatalf("expected 1, got %d", entryCount)
	}
	if totalSize < 1 {
		t.Fatalf("expected at least 1, got %d", totalSize)
	}
}

func TestForEachZipEntryFromReader(t *testing.T) {
	embedded, err := common_test.TestData.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer embedded.Close()
	entryCount := 0
	totalSize := 0
	handler := func(entry *zip.File) (err error) {
		if strings.HasSuffix(entry.Name, "/") {
			t.Fatalf("unsupported entry type %s", entry.Name)
		}
		reader, err := entry.Open()
		if err != nil {
			return
		}
		defer reader.Close()
		entryCount += 1
		b, err := io.ReadAll(reader)
		if err != nil {
			return
		}
		totalSize += len(b)
		return
	}
	err = ForEachZipEntryFromReader(handler, embedded)
	if err != nil {
		t.Fatal(err.Error())
	}
	if entryCount != 1 {
		t.Fatalf("expected 1, got %d", entryCount)
	}
	if totalSize < 1 {
		t.Fatalf("expected at least 1, got %d", totalSize)
	}
}

func TestForEachZipEntryPanic(t *testing.T) {
	embedded, err := common_test.TestData.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer embedded.Close()
	handler := func(entry *zip.File) error {
		panic("deliberate")
	}
	err = ForEachZipEntryFromReader(handler, embedded)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if err.Error() != "ZipHandler panic: deliberate" {
		t.Errorf(`expected "ZipHandler panic: deliberate", got "%s"`, err.Error())
	}
}
