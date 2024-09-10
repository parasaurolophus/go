// Copyright 2024 Kirk Rader

package utilities_test

import (
	"parasaurolophus/utilities"
	"strconv"
	"testing"
)

type aStringer int

func (as aStringer) String() string {
	return strconv.Itoa(int(as))
}

func TestGetAttribute(t *testing.T) {

	m := map[string]any{
		"string":   "a string",
		"stringer": aStringer(42),
		"int":      -1,
		"int8":     int8(-8),
	}

	if _, err := utilities.GetAttribute[string](m, "missing"); err == nil {
		t.Error("expected an error due to missing key")
	}

	if _, err := utilities.GetAttribute[int](m, "string"); err == nil {
		t.Error("expected an error due to incorrect type of value")
	}

	if actual, err := utilities.GetAttribute[string](m, "string"); err == nil {
		if actual != "a string" {
			t.Errorf(`expected "a string" but got "%s"`, actual)
		}
	} else {
		t.Error(err.Error())
	}

	if actual, err := utilities.GetAttribute[aStringer](m, "stringer"); err == nil {
		if actual != 42 {
			t.Errorf(`expected 42 but got %s`, actual)
		}
	} else {
		t.Error(err.Error())
	}

	if actual, err := utilities.GetAttribute[int](m, "int"); err == nil {
		if actual != -1 {
			t.Errorf(`expected -1 but got "%d"`, actual)
		}
	} else {
		t.Error(err.Error())
	}

	if actual, err := utilities.GetAttribute[int8](m, "int8"); err == nil {
		if actual != int8(-8) {
			t.Errorf(`expected -8 but got "%d"`, actual)
		}
	} else {
		t.Error(err.Error())
	}
}

func TestGetNumericAttribute(t *testing.T) {
	m := map[string]any{
		"string":   "4.2",
		"stringer": aStringer(42),
		"int8":     int8(-8),
		"int16":    int16(-16),
		"int32":    int32(-32),
		"int64":    int64(-64),
		"int":      -1,
		"uint8":    uint8(8),
		"uint16":   uint16(16),
		"uint32":   uint32(32),
		"uint64":   uint64(64),
		"uint":     uint(1),
		"float32":  float32(4.2),
		"float64":  4.2,
	}

	if _, err := utilities.GetNumericAttribute[int](m, "missing"); err == nil {
		t.Error("expected an error due to missing key")
	}

	if _, err := utilities.GetNumericAttribute[int](m, "string"); err == nil {
		t.Error("expected a parsing error")
	}

	if i, err := utilities.GetNumericAttribute[int](m, "stringer"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[int8](m, "int8"); err == nil {
		if i != -8 {
			t.Errorf("expected -8 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[int16](m, "int16"); err == nil {
		if i != -16 {
			t.Errorf("expected -16 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[int32](m, "int32"); err == nil {
		if i != -32 {
			t.Errorf("expected -32 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[int64](m, "int64"); err == nil {
		if i != -64 {
			t.Errorf("expected -64 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[int](m, "int"); err == nil {
		if i != -1 {
			t.Errorf("expected -1 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[uint8](m, "uint8"); err == nil {
		if i != 8 {
			t.Errorf("expected 8 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[uint16](m, "uint16"); err == nil {
		if i != 16 {
			t.Errorf("expected 16 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[uint32](m, "uint32"); err == nil {
		if i != 32 {
			t.Errorf("expected 32 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[uint64](m, "uint64"); err == nil {
		if i != 64 {
			t.Errorf("expected 64 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.GetNumericAttribute[uint](m, "uint"); err == nil {
		if i != 1 {
			t.Errorf("expected 1 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if f, err := utilities.GetNumericAttribute[float32](m, "float32"); err == nil {
		if f != 4.2 {
			t.Errorf("expected 4.2 but got %f", f)
		}
	} else {
		t.Error(err.Error())
	}

	if f, err := utilities.GetNumericAttribute[float64](m, "float64"); err == nil {
		if f != 4.2 {
			t.Errorf("expected 4.2 but got %f", f)
		}
	} else {
		t.Error(err.Error())
	}
}

func TestParseNumber(t *testing.T) {

	if _, err := utilities.ParseNumber[int]("forty-two"); err == nil {
		t.Error("expected a parsing error")
	}

	if _, err := utilities.ParseNumber[uint]("-42"); err == nil {
		t.Error("expected a parsing error")
	}

	if i, err := utilities.ParseNumber[int8]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[int16]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[int32]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[int64]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[int]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[uint8]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[uint16]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[uint32]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[uint64]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if i, err := utilities.ParseNumber[uint]("42"); err == nil {
		if i != 42 {
			t.Errorf("expected 42 but got %d", i)
		}
	} else {
		t.Error(err.Error())
	}

	if f, err := utilities.ParseNumber[float32]("4.2"); err == nil {
		if f != 4.2 {
			t.Errorf("expected 4.2 but got %f", f)
		}
	} else {
		t.Error(err.Error())
	}

	if f, err := utilities.ParseNumber[float64]("4.2"); err == nil {
		if f != 4.2 {
			t.Errorf("expected 4.2 but got %f", f)
		}
	} else {
		t.Error(err.Error())
	}
}
