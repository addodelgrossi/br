package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrintObjectJSON(t *testing.T) {
	var buf bytes.Buffer
	err := PrintObject(&buf, Options{JSON: true}, []Field{{Key: "cep", Label: "CEP", Value: "01001000"}}, map[string]string{"cep": "01001000"}, "01001000")
	if err != nil {
		t.Fatalf("PrintObject() error = %v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded["cep"] != "01001000" {
		t.Fatalf("decoded JSON = %v", decoded)
	}
}

func TestPrintListCSV(t *testing.T) {
	var buf bytes.Buffer
	err := PrintList(
		&buf,
		Options{CSV: true},
		[]Column{{Key: "date", Label: "Data"}, {Key: "name", Label: "Nome"}},
		[]map[string]string{{"date": "2026-01-01", "name": "Confraternizacao universal"}},
		nil,
		"name",
	)
	if err != nil {
		t.Fatalf("PrintList() error = %v", err)
	}
	if !strings.HasPrefix(buf.String(), "date,name\n") {
		t.Fatalf("CSV output = %q", buf.String())
	}
}

func TestPrintQuiet(t *testing.T) {
	var buf bytes.Buffer
	err := PrintObject(&buf, Options{Quiet: true}, nil, nil, "true")
	if err != nil {
		t.Fatalf("PrintObject() error = %v", err)
	}
	if buf.String() != "true\n" {
		t.Fatalf("quiet output = %q", buf.String())
	}
}
