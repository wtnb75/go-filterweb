package filterweb

import (
	"strings"
	"testing"
	"time"
)

func TestMatch(t *testing.T) {
	if !match("^a.*", "abc") {
		t.Fatalf("expected match true")
	}
	if match("^b.*", "abc") {
		t.Fatalf("expected match false")
	}
}

func TestCapture(t *testing.T) {
	pat := `(?P<name>\w+)-(\d+)`
	res := capture(pat, "bob-123")
	if got := res["0"]; got != "bob-123" {
		t.Fatalf("0 mismatch: %q", got)
	}
	if got := res["1"]; got != "bob" {
		t.Fatalf("1 mismatch: %q", got)
	}
	if got := res["2"]; got != "123" {
		t.Fatalf("2 mismatch: %q", got)
	}
	if got := res["name"]; got != "bob" {
		t.Fatalf("name mismatch: %q", got)
	}
}

func TestStrptimeAndDoStrftime(t *testing.T) {
	ts := "2020-01-02"
	tm := strptime("%Y-%m-%d", ts)
	if tm.IsZero() {
		t.Fatalf("strptime returned zero time")
	}
	if tm.Year() != 2020 || tm.Month() != 1 || tm.Day() != 2 {
		t.Fatalf("unexpected date: %v", tm)
	}
	s := do_strftime("%Y-%m-%d", tm)
	if s != "2020-01-02" {
		t.Fatalf("strftime mismatch: %q", s)
	}
}

func TestInAndToJSONAlias(t *testing.T) {
	arr := []any{"a", 1, "b"}
	if !in("a", arr) {
		t.Fatalf("in should find 'a'")
	}
	if in("z", arr) {
		t.Fatalf("in should not find 'z'")
	}
	// tojson is implemented same as in; verify behavior
	if !tojson(1, arr) {
		t.Fatalf("tojson should find 1")
	}
}

func TestToYAMLAndToXML(t *testing.T) {
	m := map[string]any{"name": "Alice", "age": 30}
	y := toyaml(m)
	if !strings.Contains(y, "name:") || !strings.Contains(y, "Alice") {
		t.Fatalf("toyaml output unexpected: %q", y)
	}
	type Person struct {
		Name string `xml:"name"`
	}
	p := Person{Name: "Bob"}
	x := toxml(p)
	if !strings.Contains(x, "<name>Bob</name>") {
		t.Fatalf("toxml output unexpected: %q", x)
	}
}

func TestRFC3339(t *testing.T) {
	tm := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	s := rfc3339(tm)
	if s != "2020-01-02T03:04:05Z" {
		t.Fatalf("rfc3339 mismatch: %q", s)
	}
}

func TestMatch_InvalidPattern(t *testing.T) {
	// invalid regex should not panic; match() returns false
	if match("(*", "abc") {
		t.Fatalf("expected match false for invalid pattern")
	}
}

func TestCapture_InvalidPattern_Panic(t *testing.T) {
	// capture uses MustCompile and will panic on invalid pattern
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid pattern")
		}
	}()
	_ = capture("(*", "abc")
}

func TestStrptime_InvalidInput(t *testing.T) {
	tm := strptime("%Y-%m-%d", "not-a-date")
	if !tm.IsZero() {
		t.Fatalf("expected zero time for invalid input, got %v", tm)
	}
}
