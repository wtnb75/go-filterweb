package filterweb

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTP_Prep_MissingURL(t *testing.T) {
	hc := &HTTPConfig{}
	cfg := Config{Params: map[string]any{}}
	err := hc.Prep(cfg, Data{})
	if err == nil {
		t.Fatalf("expected ErrMissingParams")
	}
	if err != ErrMissingParams {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTP_Process_JSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = io.WriteString(w, `{"name":"Alice"}`)
	}))
	defer srv.Close()

	hc := &HTTPConfig{}
	cfg := Config{Params: map[string]any{"Url": srv.URL}}
	if err := hc.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	out, err := hc.Process(Data{})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	if out.ContentType != "application/json" {
		t.Fatalf("unexpected content type: %v", out.ContentType)
	}
	m, ok := out.Data.(map[string]any)
	if !ok {
		t.Fatalf("unexpected data type: %T", out.Data)
	}
	if m["name"] != "Alice" {
		t.Fatalf("unexpected name: %v", m["name"])
	}
}

func TestHTTP_Process_ContentTypeOverride_YAML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return plain text but the filter should use overridden content type
		w.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(w, "name: Bob")
	}))
	defer srv.Close()

	hc := &HTTPConfig{}
	cfg := Config{Params: map[string]any{"Url": srv.URL, "ContentType": "application/yaml"}}
	if err := hc.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	out, err := hc.Process(Data{})
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	if out.ContentType != "application/yaml" {
		t.Fatalf("unexpected content type: %v", out.ContentType)
	}
	m, ok := out.Data.(map[string]any)
	if !ok {
		t.Fatalf("unexpected data type: %T", out.Data)
	}
	if m["name"] != "Bob" {
		t.Fatalf("unexpected name: %v", m["name"])
	}
}

func TestHTTP_Process_StatusNotOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = io.WriteString(w, "error")
	}))
	defer srv.Close()

	hc := &HTTPConfig{}
	cfg := Config{Params: map[string]any{"Url": srv.URL}}
	if err := hc.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	_, err := hc.Process(Data{})
	if err == nil {
		t.Fatalf("expected ErrHTTPStatusNotOK")
	}
	if err != ErrHTTPStatusNotOK {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTP_Process_InvalidContentTypeHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ";;;invalid;;;")
		_, _ = io.WriteString(w, "data")
	}))
	defer srv.Close()

	hc := &HTTPConfig{}
	cfg := Config{Params: map[string]any{"Url": srv.URL}}
	if err := hc.Prep(cfg, Data{}); err != nil {
		t.Fatalf("Prep failed: %v", err)
	}
	_, err := hc.Process(Data{})
	if err == nil {
		t.Fatalf("expected parse error for invalid Content-Type header")
	}
}
