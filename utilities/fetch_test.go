package utilities

import (
	"io"
	"testing"
)

func TestFetch(t *testing.T) {
	readCloser, err := Fetch("https://www.google.com")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer readCloser.Close()
	buffer, err := io.ReadAll(readCloser)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(buffer) < 1 {
		t.Errorf("ReadAll returned %d", len(buffer))
	}
}
