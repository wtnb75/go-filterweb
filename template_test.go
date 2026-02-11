package filterweb

import (
    "os"
    "path/filepath"
    "testing"
)

func TestTemplate_PrepAndProcess_Text(t *testing.T) {
    tc := &TemplateConfig{}
    cfg := Config{Params: map[string]any{"Type": "text", "Content": "Hello {{.name}}"}}
    if err := tc.Prep(cfg, Data{}); err != nil {
        t.Fatalf("Prep failed: %v", err)
    }
    out, err := tc.Process(Data{Data: map[string]any{"name": "Alice"}})
    if err != nil {
        t.Fatalf("Process failed: %v", err)
    }
    if out.ContentType != "text/plain" {
        t.Fatalf("unexpected content type: %v", out.ContentType)
    }
    s, ok := out.Data.(string)
    if !ok {
        t.Fatalf("output data is not string: %T", out.Data)
    }
    if s != "Hello Alice" {
        t.Fatalf("unexpected output: %q", s)
    }
}

func TestTemplate_Prep_FileAndVarsAndBaseKey(t *testing.T) {
    // create temp dir/file
    dir := t.TempDir()
    fpath := filepath.Join(dir, "tmpl.txt")
    content := "Greeting: {{.base.name}} {{.extra}}"
    if err := os.WriteFile(fpath, []byte(content), 0644); err != nil {
        t.Fatalf("write temp file: %v", err)
    }

    tc := &TemplateConfig{}
    cfg := Config{Params: map[string]any{
        "Type":    "text",
        "File":    fpath,
        "BaseKey": "base",
        "Vars":    map[string]any{"extra": "World"},
    }}
    if err := tc.Prep(cfg, Data{}); err != nil {
        t.Fatalf("Prep from file failed: %v", err)
    }
    out, err := tc.Process(Data{Data: map[string]any{"name": "Bob"}})
    if err != nil {
        t.Fatalf("Process failed: %v", err)
    }
    s, ok := out.Data.(string)
    if !ok {
        t.Fatalf("output data is not string: %T", out.Data)
    }
    expected := "Greeting: Bob World"
    if s != expected {
        t.Fatalf("unexpected output: got=%q want=%q", s, expected)
    }
}

func TestTemplate_Prep_MissingParams(t *testing.T) {
    tc := &TemplateConfig{}
    cfg := Config{Params: map[string]any{}}
    err := tc.Prep(cfg, Data{})
    if err == nil {
        t.Fatalf("expected error for missing params")
    }
    if err != ErrMissingParams {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestTemplate_VarsMerge(t *testing.T) {
    tc := &TemplateConfig{}
    cfg := Config{Params: map[string]any{
        "Type":    "text",
        "Content": "Hello {{.name}} {{.extra}}",
        "Vars":    map[string]any{"extra": "Everybody"},
    }}
    if err := tc.Prep(cfg, Data{}); err != nil {
        t.Fatalf("Prep failed: %v", err)
    }
    out, err := tc.Process(Data{Data: map[string]any{"name": "Jane"}})
    if err != nil {
        t.Fatalf("Process failed: %v", err)
    }
    s, ok := out.Data.(string)
    if !ok {
        t.Fatalf("output data is not string: %T", out.Data)
    }
    if s != "Hello Jane Everybody" {
        t.Fatalf("unexpected output: %q", s)
    }
}

func TestTemplate_Prep_UnsupportedType(t *testing.T) {
    tc := &TemplateConfig{}
    cfg := Config{Params: map[string]any{"Type": "unknown", "Content": "tmpl"}}
    err := tc.Prep(cfg, Data{})
    if err == nil {
        t.Fatalf("expected error for unsupported template type")
    }
    if err != ErrReadTemplate {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestTemplate_HTMLType_Escaping(t *testing.T) {
    tc := &TemplateConfig{}
    // html/template should escape variable content
    cfg := Config{Params: map[string]any{"Type": "html", "Content": "<div>{{.name}}</div>"}}
    if err := tc.Prep(cfg, Data{}); err != nil {
        t.Fatalf("Prep failed: %v", err)
    }
    out, err := tc.Process(Data{Data: map[string]any{"name": "<b>Alice</b>"}})
    if err != nil {
        t.Fatalf("Process failed: %v", err)
    }
    if out.ContentType != "text/html" {
        t.Fatalf("unexpected content type: %v", out.ContentType)
    }
    s, ok := out.Data.(string)
    if !ok {
        t.Fatalf("output data is not string: %T", out.Data)
    }
    expected := "<div>&lt;b&gt;Alice&lt;/b&gt;</div>"
    if s != expected {
        t.Fatalf("unexpected output: got=%q want=%q", s, expected)
    }
}
