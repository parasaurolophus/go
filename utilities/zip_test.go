// Copyright Kirk Rader 2024

package utilities

import (
	"archive/zip"
	"io"
	"os"
	"parasaurolophus/go/common_test"
	"strings"
	"testing"
)

func TestForZipFile(t *testing.T) {
	archive, err := common_test.TestData.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer archive.Close()
	tempFile, err := os.CreateTemp(t.TempDir(), "TestForZipFile")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()
	_, err = io.Copy(tempFile, archive)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	entryCount := 0
	totalSize := 0
	handler := func(entry *zip.File) {
		if strings.HasSuffix(entry.Name, "/") {
			t.Fatalf("unsupported entry type %s", entry.Name)
		}
		reader, err := entry.Open()
		if err != nil {
			t.Fatal(err.Error())
		}
		defer reader.Close()
		entryCount += 1
		b, err := io.ReadAll(reader)
		if err != nil {
			t.Fatal(err.Error())
		}
		totalSize += len(b)
	}
	err = ForZipFile(handler, tempFile)
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

func TestForZipReadCloser(t *testing.T) {
	archive, err := common_test.TestData.Open("testdata/eurofxref.zip")
	if err != nil {
		t.Fatal(err.Error())
	}
	entryCount := 0
	totalSize := 0
	handler := func(entry *zip.File) {
		if strings.HasSuffix(entry.Name, "/") {
			t.Fatalf("unsupported entry type %s", entry.Name)
		}
		reader, err := entry.Open()
		if err != nil {
			t.Fatal(err.Error())
		}
		defer reader.Close()
		entryCount += 1
		b, err := io.ReadAll(reader)
		if err != nil {
			t.Fatal(err.Error())
		}
		totalSize += len(b)
	}
	defer archive.Close()
	err = ForZipReader(handler, archive)
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
