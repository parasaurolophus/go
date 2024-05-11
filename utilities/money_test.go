// Copyright Kirk Rader 2024

package utilities

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	type S struct {
		M Money `json:"m,omitempty"`
	}
	s := S{
		M: New(4.2, 2),
	}
	b, err := json.Marshal(&s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != `{"m":4.20}` {
		t.Errorf(`expected {"m":4.20}, got %s`, string(b))
	}
	s = S{
		M: New(4.2, 0),
	}
	b, err = json.Marshal(&s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != `{"m":4}` {
		t.Errorf(`expected {"m":4}, got %s`, string(b))
	}
}

func TestScan(t *testing.T) {
	m1 := New(0.0, 2)
	m2 := New(0.0, 2)
	n, err := fmt.Sscan(" 4.20 -0.01 ", &m1, &m2)
	if err != nil {
		t.Error(err.Error())
	}
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
	if m1.Value() != 4.2 {
		t.Errorf("expected 4.2, got %f", m1.Value())
	}
	if m2.Value() != -0.01 {
		t.Errorf("expected -0.01, got %f", m1.Value())
	}
	n, err = fmt.Sscan("12.34.", &m1)
	if err != nil {
		t.Error(err.Error())
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if m1.Value() != 12.34 {
		t.Errorf("expected 12.34, got %f", m1.Value())
	}
	n, err = fmt.Sscan("567.80-", &m1)
	if err != nil {
		t.Error(err.Error())
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if m1.Value() != 567.8 {
		t.Errorf("expected 567.8, got %f", m1.Value())
	}
}

func TestString(t *testing.T) {
	m := New(78.9, 2)
	s := m.String()
	if s != "78.90" {
		t.Errorf(`expected "78.90", got "%s"`, s)
	}
	m = New(78.9, 1)
	s = m.String()
	if s != "78.9" {
		t.Errorf(`expected "78.9", got "%s"`, s)
	}
	m = New(78.9, -1)
	s = m.String()
	if s != "79" {
		t.Errorf(`expected "79", got "%s"`, s)
	}
}

func TestUnmarshal(t *testing.T) {
	type S struct {
		M Money `json:"m,omitempty"`
	}
	var s S
	err := json.Unmarshal([]byte(`{}`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M.Value() != 0.0 {
		t.Errorf("expected 0.0, got %f", s.M.Value())
	}
	err = json.Unmarshal([]byte(`{"m": 4.200}`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M.Value() != 4.2 {
		t.Errorf("expected 4.2, got %f", s.M.Value())
	}
}
