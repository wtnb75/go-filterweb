package filterweb

import (
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestEncode_Prep_MissingContentType(t *testing.T) {
	ec := &EncodeConfig{}
	cfg := Config{Params: map[string]any{}}
	err := ec.Prep(cfg, Data{})
	if err == nil {
		t.Fatalf("expected ErrMissingParams")
	}
	if err != ErrMissingParams {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEncode_Process_JSON(t *testing.T) {
	ec := &EncodeConfig{}
	cfg := Config{Params: map[string]any{"ContentType": "application/json"}}
	if err := ec.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	in := map[string]any{"name": "Alice", "age": 30}
	out, err := ec.Process(Data{Data: in})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	b, ok := out.Data.([]byte)
	if !ok {
		t.Fatalf("expected []byte, got %T", out.Data)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}
	if m["name"] != "Alice" {
		t.Fatalf("unexpected name: %v", m["name"])
	}
}

func TestEncode_Process_YAML(t *testing.T) {
	ec := &EncodeConfig{}
	cfg := Config{Params: map[string]any{"ContentType": "application/yaml"}}
	if err := ec.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	in := map[string]any{"name": "Bob", "age": 25}
	out, err := ec.Process(Data{Data: in})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	b, ok := out.Data.([]byte)
	if !ok {
		t.Fatalf("expected []byte, got %T", out.Data)
	}
	var m map[string]any
	if err := yaml.Unmarshal(b, &m); err != nil {
		t.Fatalf("yaml unmarshal failed: %v", err)
	}
	if m["name"] != "Bob" {
		t.Fatalf("unexpected name: %v", m["name"])
	}
}

func TestEncode_Process_XML(t *testing.T) {
	ec := &EncodeConfig{}
	cfg := Config{Params: map[string]any{"ContentType": "text/xml"}}
	if err := ec.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	type Person struct {
		Name string `xml:"name"`
	}
	p := Person{Name: "Carol"}
	out, err := ec.Process(Data{Data: p})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	b, ok := out.Data.([]byte)
	if !ok {
		t.Fatalf("expected []byte, got %T", out.Data)
	}
	s := string(b)
	if !strings.Contains(s, "<name>Carol</name>") {
		t.Fatalf("unexpected xml output: %q", s)
	}
}

func TestEncode_Process_CSV(t *testing.T) {
	ec := &EncodeConfig{}
	cfg := Config{Params: map[string]any{"ContentType": "text/csv"}}
	if err := ec.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	records := []map[string]any{{"col1": "a", "col2": "1"}, {"col1": "b", "col2": "2"}}
	out, err := ec.Process(Data{Data: records})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	b, ok := out.Data.([]byte)
	if !ok {
		t.Fatalf("expected []byte, got %T", out.Data)
	}
	r := csv.NewReader(strings.NewReader(string(b)))
	hdr, err := r.Read()
	if err != nil {
		t.Fatalf("csv read header failed: %v", err)
	}
	if len(hdr) != 2 {
		t.Fatalf("unexpected header length: %v", hdr)
	}
	// read rows
	row1, err := r.Read()
	if err != nil {
		t.Fatalf("csv read row1 failed: %v", err)
	}
	if row1[0] != "a" && row1[1] != "1" {
		t.Fatalf("unexpected row1: %v", row1)
	}
}
