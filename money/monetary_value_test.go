// Copyright Kirk Rader 2024

package money

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	type S struct {
		M Money `json:"m,omitempty"`
	}
	s := S{}
	b, err := json.Marshal(s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != "{}" {
		t.Errorf(`expected "{}", got "%s"`, string(b))
	}
	s.M = NewMoney(0.0, 2)
	b, err = json.Marshal(s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != `{"m":0.00}` {
		t.Errorf(`expected {"m":0.00}, got %s`, string(b))
	}
	s.M = NewMoney(4.2, 0)
	b, err = json.Marshal(s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != `{"m":4}` {
		t.Errorf(`expected {"m":4}, got %s`, string(b))
	}
}

func TestMarshalXML(t *testing.T) {
	type S struct {
		XMLName xml.Name `xml:"s"`
		M       Money    `xml:"m"`
	}
	s := S{}
	b, err := xml.Marshal(s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != "<s></s>" {
		t.Errorf(`expected "<s></s>", got "%s"`, string(b))
	}
	s.M = NewMoney(0.0, 2)
	b, err = xml.Marshal(s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != `<s><m>0.00</m></s>` {
		t.Errorf(`expected <s><m>0.00</m></s>, got %s`, string(b))
	}
	s.M = NewMoney(4.2, 0)
	b, err = xml.Marshal(s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != `<s><m>4</m></s>` {
		t.Errorf(`expected <s><m>4</m></s>, got %s`, string(b))
	}
}

func TestScan(t *testing.T) {
	m1 := NewMoney(0.0, 2)
	m2 := NewMoney(0.0, 2)
	n, err := fmt.Sscan(" 4.20 -0.01 ", m1, m2)
	if err != nil {
		t.Error(err.Error())
	}
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
	if m1.Get() != 4.2 {
		t.Errorf("expected 4.2, got %f", m1.Get())
	}
	if m2.Get() != -0.01 {
		t.Errorf("expected -0.01, got %f", m1.Get())
	}
	n, err = fmt.Sscan("12.34.", m1)
	if err != nil {
		t.Error(err.Error())
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if m1.Get() != 12.34 {
		t.Errorf("expected 12.34, got %f", m1.Get())
	}
	n, err = fmt.Sscan("567.80-", m1)
	if err != nil {
		t.Error(err.Error())
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if m1.Get() != 567.8 {
		t.Errorf("expected 567.8, got %f", m1.Get())
	}
}

func TestString(t *testing.T) {
	m := NewMoney(78.9, 2)
	s := m.String()
	if s != "78.90" {
		t.Errorf(`expected "78.90", got "%s"`, s)
	}
	m = NewMoney(78.9, 1)
	s = m.String()
	if s != "78.9" {
		t.Errorf(`expected "78.9", got "%s"`, s)
	}
	m = NewMoney(78.9, -1)
	s = m.String()
	if s != "79" {
		t.Errorf(`expected "79", got "%s"`, s)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type S struct {
		M Money `json:"m,omitempty"`
	}
	var s S
	err := json.Unmarshal([]byte(`{}`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M != nil {
		t.Errorf(`expected nil, got %v`, s.M)
	}
	s.M = NewMoney(0.0, 2)
	err = json.Unmarshal([]byte(`{"m": 4.200}`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M.Get() != 4.2 {
		t.Errorf("expected 4.2, got %f", s.M.Get())
	}
}

func TestUnmarshalXMLAttribute(t *testing.T) {
	type S struct {
		XMLName xml.Name `xml:"s"`
		M       Money    `xml:"m,attr"`
	}
	var s S
	err := xml.Unmarshal([]byte(`<s></s>`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M != nil {
		t.Errorf("expected nil, got %v", s.M)
	}
	s.M = NewMoney(0.0, 2)
	err = xml.Unmarshal([]byte(`<s m="4.20"/>`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M.Get() != 4.2 {
		t.Errorf("expected 4.2, got %f", s.M.Get())
	}
	err = xml.Unmarshal([]byte(`<s m="invalid"/>`), &s)
	if err == nil {
		t.Error("expected err not to be nil")
	}
	if s.M.Get() != 0.0 {
		t.Errorf("expected 0.0, got %f", s.M.Get())
	}
}

func TestUnmarshalXMLElement(t *testing.T) {
	type S struct {
		XMLName xml.Name `xml:"s"`
		M       Money    `xml:"m"`
	}
	var s S
	err := xml.Unmarshal([]byte(`<s></s>`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M != nil {
		t.Errorf("expected nil, got %v", s.M)
	}
	s.M = NewMoney(0.0, 2)
	err = xml.Unmarshal([]byte(`<s><m>4.20</m></s>`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M.Get() != 4.2 {
		t.Errorf("expected 4.2, got %f", s.M.Get())
	}
}
