// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"io"
	"testing"
)

func TestFetchGoodURL(t *testing.T) {
	readCloser, err := Fetch("http://rader.us")
	if err != nil {
		t.Fatal(err.Error())
	}
	if readCloser == nil {
		t.Fatal("Fetch() returned nil readCloser")
	}
	defer readCloser.Close()
	buffer, err := io.ReadAll(readCloser)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(buffer) < 1 {
		t.Errorf("ReadAll returned %d", len(buffer))
	}
	fmt.Println(string(buffer))
}

func TestFetchBadURL(t *testing.T) {
	badFetch(t, "http://invalid")
	badFetch(t, "http://rader.us/invalid")
}

func badFetch(t *testing.T, url string) {
	readCloser, err := Fetch(url)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if readCloser != nil {
		readCloser.Close()
		t.Fatal("expected readCloser to be nil")
	}
}
