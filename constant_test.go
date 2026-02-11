package filterweb

import (
	"bytes"
	"testing"
)

func TestConstant_PrepAndProcess_DefaultText(t *testing.T) {
	hc := &ConstantConfig{}
	cfg := Config{Params: map[string]any{"Data": "hello"}}
	if err := hc.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	if hc.ContentType != "text/plain" {
		t.Fatalf("unexpected default ContentType: %v", hc.ContentType)
	}
	out, err := hc.Process(Data{})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	b, ok := out.Data.([]byte)
	if !ok {
		t.Fatalf("expected []byte data, got: %T", out.Data)
	}
	if !bytes.Equal(b, []byte("hello")) {
		t.Fatalf("unexpected data: %q", string(b))
	}
}

func TestConstant_Prep_MissingParams(t *testing.T) {
	hc := &ConstantConfig{}
	cfg := Config{Params: map[string]any{}}
	err := hc.Prep(cfg, Data{})
	if err == nil {
		t.Fatalf("expected ErrMissingParams")
	}
	if err != ErrMissingParams {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConstant_Process_YAML_String(t *testing.T) {
	hc := &ConstantConfig{}
	cfg := Config{Params: map[string]any{"ContentType": "application/yaml", "Data": "name: Alice"}}
	if err := hc.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	out, err := hc.Process(Data{})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	m, ok := out.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", out.Data)
	}
	if v, ok := m["name"]; !ok || v != "Alice" {
		t.Fatalf("unexpected map contents: %#v", m)
	}
}

func TestConstant_Process_YAML_Bytes(t *testing.T) {
	hc := &ConstantConfig{}
	cfg := Config{Params: map[string]any{"ContentType": "application/yaml", "Data": []byte("name: Bob")}}
	if err := hc.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	out, err := hc.Process(Data{})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	m, ok := out.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", out.Data)
	}
	if v, ok := m["name"]; !ok || v != "Bob" {
		t.Fatalf("unexpected map contents: %#v", m)
	}
}

func TestConstant_Process_NonString(t *testing.T) {
	hc := &ConstantConfig{}
	cfg := Config{Params: map[string]any{"Data": 42}}
	if err := hc.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	out, err := hc.Process(Data{})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	if out.Data != 42 {
		t.Fatalf("unexpected data: %#v", out.Data)
	}
}
