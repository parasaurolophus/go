// Copyright Kirk Rader 2024

package enums

import (
	"fmt"
	"testing"
)

func TestOneString(t *testing.T) {

	actual := One.String()

	if actual != one {
		t.Errorf("expected \"%s\", got \"%s\"", one, actual)
	}
}

func TestTwoString(t *testing.T) {

	actual := Two.String()

	if actual != two {
		t.Errorf("expected \"%s\", got \"%s\"", two, actual)
	}
}

func TestThreeString(t *testing.T) {

	actual := Three.String()

	if actual != three {
		t.Errorf("expected \"%s\", got \"%s\"", three, actual)
	}
}

func TestEnum1String(t *testing.T) {

	actual := Enum1(100).String()

	if actual != "" {
		t.Errorf("expected \"\", got \"%s\"", actual)
	}
}

func TestFourString(t *testing.T) {

	actual := Four.String()

	if actual != four {
		t.Errorf("expected \"%s\", got \"%s\"", four, actual)
	}
}

func TestFiveString(t *testing.T) {

	actual := Five.String()

	if actual != five {
		t.Errorf("expected \"%s\", got \"%s\"", five, actual)
	}
}

func TestSixString(t *testing.T) {

	actual := Six.String()

	if actual != six {
		t.Errorf("expected \"%s\", got \"%s\"", six, actual)
	}
}

func TestEnum2String(t *testing.T) {

	actual := Enum2(100).String()

	if actual != "" {
		t.Errorf("expected \"\", got \"%s\"", actual)
	}
}

func TestOneScan(t *testing.T) {

	var e1 Enum1

	n, err := fmt.Sscan(one, &e1)

	if err != nil {
		t.Errorf(err.Error())
	}

	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	if e1 != One {
		t.Errorf("expected One, got %s", e1)
	}
}

func TestTwoScan(t *testing.T) {

	var e1 Enum1

	n, err := fmt.Sscan(two, &e1)

	if err != nil {
		t.Errorf(err.Error())
	}

	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	if e1 != Two {
		t.Errorf("expected Two, got %s", e1)
	}
}

func TestThreeScan(t *testing.T) {

	var e1 Enum1

	n, err := fmt.Sscan(three, &e1)

	if err != nil {
		t.Errorf(err.Error())
	}

	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	if e1 != Three {
		t.Errorf("expected Three, got %s", e1)
	}
}

func TestEnum1Scan(t *testing.T) {

	var e1 Enum1

	n, err := fmt.Sscan(four, &e1)

	if err == nil {
		t.Errorf("expected an error to be signaled")
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestFourScan(t *testing.T) {

	var e2 Enum2

	n, err := fmt.Sscan(four, &e2)

	if err != nil {
		t.Errorf(err.Error())
	}

	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	if e2 != Four {
		t.Errorf("expected Four, got %s", e2)
	}
}

func TestFiveScan(t *testing.T) {

	var e2 Enum2

	n, err := fmt.Sscan(five, &e2)

	if err != nil {
		t.Errorf(err.Error())
	}

	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	if e2 != Five {
		t.Errorf("expected Five, got %s", e2)
	}
}

func TestSixScan(t *testing.T) {

	var e2 Enum2

	n, err := fmt.Sscan(six, &e2)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}

	if e2 != Six {
		t.Fatalf("expected Six, got %s", e2)
	}
}

func TestEnum2Scan(t *testing.T) {

	var e2 Enum2

	n, err := fmt.Sscan(one, &e2)

	if err == nil {
		t.Errorf("expected an error to be signaled")
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestMultiScan(t *testing.T) {

	var e1 Enum1
	var e2 Enum2

	format := fmt.Sprintf("%s, %s", two, five)
	n, err := fmt.Sscanf(format, "%v, %v", &e1, &e2)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}

	if e1 != Two {
		t.Fatalf("expected %v, got %v", Two, e1)
	}

	if e2 != Five {
		t.Fatalf("expected %v, got %v", Five, e2)
	}
}
