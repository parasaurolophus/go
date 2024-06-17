// Copyright Kirk Rader 2024

package money

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
)

func TestGetDigits(t *testing.T) {
	m, err := NewMoney(4.2, 3)
	if err != nil {
		t.Fatal(err.Error())
	}
	actual := m.GetDigits()
	if actual != 3 {
		t.Fatalf("expected 3, got %d", actual)
	}
}

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
	s.M, err = NewMoney(0.0, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	b, err = json.Marshal(s)
	if err != nil {
		t.Fatal(err.Error())
	}
	if string(b) != `{"m":0.00}` {
		t.Errorf(`expected {"m":0.00}, got %s`, string(b))
	}
	s.M, err = NewMoney(4.2, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
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
	s.M, err = NewMoney(0.0, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	b, err = xml.Marshal(s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != `<s><m>0.00</m></s>` {
		t.Errorf(`expected <s><m>0.00</m></s>, got %s`, string(b))
	}
	s.M, err = NewMoney(4.2, 0)
	if err != nil {
		t.Fatal(err.Error())
	}
	b, err = xml.Marshal(s)
	if err != nil {
		t.Error(err.Error())
	}
	if string(b) != `<s><m>4</m></s>` {
		t.Errorf(`expected <s><m>4</m></s>, got %s`, string(b))
	}
}

func TestScan(t *testing.T) {
	m1, err := NewMoney(0.0, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	m2, err := NewMoney(0.0, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	n, err := fmt.Sscan(" 4.20 -0.01 ", m1, m2)
	if err != nil {
		t.Fatal(err.Error())
	}
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
	if m1.GetValue() != 4.2 {
		t.Errorf("expected 4.2, got %f", m1.GetValue())
	}
	if m2.GetValue() != -0.01 {
		t.Errorf("expected -0.01, got %f", m1.GetValue())
	}
	n, err = fmt.Sscan("12.34.", m1)
	if err != nil {
		t.Error(err.Error())
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if m1.GetValue() != 12.34 {
		t.Errorf("expected 12.34, got %f", m1.GetValue())
	}
	n, err = fmt.Sscan("567.80-", m1)
	if err != nil {
		t.Error(err.Error())
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if m1.GetValue() != 567.8 {
		t.Errorf("expected 567.8, got %f", m1.GetValue())
	}
}

func TestString(t *testing.T) {
	m, err := NewMoney(78.9, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	s := m.String()
	if s != "78.90" {
		t.Errorf(`expected "78.90", got "%s"`, s)
	}
	m, err = NewMoney(78.9, 1)
	if err != nil {
		t.Fatal(err.Error())
	}
	s = m.String()
	if s != "78.9" {
		t.Errorf(`expected "78.9", got "%s"`, s)
	}
	_, err = NewMoney(78.9, -1)
	if err == nil {
		t.Fatal("expected err not to be nil")
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
	s.M, err = NewMoney(0.0, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = json.Unmarshal([]byte(`{"m": 4.200}`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M.GetValue() != 4.2 {
		t.Errorf("expected 4.2, got %f", s.M.GetValue())
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
	s.M, err = NewMoney(0.0, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = xml.Unmarshal([]byte(`<s m="4.20"/>`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M.GetValue() != 4.2 {
		t.Errorf("expected 4.2, got %f", s.M.GetValue())
	}
	err = xml.Unmarshal([]byte(`<s m="invalid"/>`), &s)
	if err == nil {
		t.Error("expected err not to be nil")
	}
	if s.M.GetValue() != 0.0 {
		t.Errorf("expected 0.0, got %f", s.M.GetValue())
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
	s.M, err = NewMoney(0.0, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = xml.Unmarshal([]byte(`<s><m>4.20</m></s>`), &s)
	if err != nil {
		t.Errorf(err.Error())
	}
	if s.M.GetValue() != 4.2 {
		t.Errorf("expected 4.2, got %f", s.M.GetValue())
	}
}
