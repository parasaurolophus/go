// Copyright Kirk Rader 2024

package utilities

import (
	"fmt"
	"strconv"
	"testing"
)

type MyFloat float64

func (m *MyFloat) Scan(state fmt.ScanState, _ rune) error {
	token, err := state.Token(true, MakeJSONNumberTokenTest())
	if err != nil {
		return err
	}
	f, err := strconv.ParseFloat(string(token), 64)
	if err != nil {
		return err
	}
	*m = MyFloat(f)
	return nil
}
func TestMakeJSONNumberTokenTest(t *testing.T) {
	tokenTest := MakeJSONNumberTokenTest()
	result := tokenTest('-')
	if !result {
		t.Fatal("expected true for '-'")
	}
	result = tokenTest('1')
	if !result {
		t.Fatal("expected true for '1'")
	}
	result = tokenTest('2')
	if !result {
		t.Fatal("expected true for '2'")
	}
	result = tokenTest('.')
	if !result {
		t.Fatal("expected true for '.'")
	}
	result = tokenTest('3')
	if !result {
		t.Fatal("expected true for '3'")
	}
	result = tokenTest('e')
	if !result {
		t.Fatal("expected true for 'e'")
	}
	result = tokenTest('-')
	if !result {
		t.Fatal("expected true for '-'")
	}
	result = tokenTest('4')
	if !result {
		t.Fatal("expected true for '4'")
	}
	result = tokenTest('-')
	if result {
		t.Fatal("expected false for '-'")
	}
	result = tokenTest('.')
	if result {
		t.Fatal("expected false for '.'")
	}
	result = tokenTest('e')
	if result {
		t.Fatal("expected false for 'e'")
	}
	result = tokenTest(',')
	if result {
		t.Fatal("expected false for ','")
	}
}

func TestScanJSONNumber(t *testing.T) {
	var m1, m2, m3 MyFloat
	n, err := fmt.Sscan("12.3e-1-5.678 9.10.", &m1, &m2, &m3)
	if err != nil {
		t.Fatal(err.Error())
	}
	if n != 3 {
		t.Fatalf("expected 3, got %d", n)
	}
	if m1 != 1.23 {
		t.Errorf("expected 1.23, got %f", m1)
	}
	if m2 != -5.678 {
		t.Errorf("expected -5.678, got %f", m2)
	}
	if m3 != 9.1 {
		t.Errorf("expected 9.1, got %f", m3)
	}
	n, err = fmt.Sscan(" ", &m1)
	if err == nil {
		t.Fatal("expected err not to be nil")
	}
	if n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
}
