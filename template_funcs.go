package filterweb

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"log/slog"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/ncruces/go-strftime"
)

func match(pattern, str string) bool {
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		slog.Error("match function error", "pattern", pattern, "string", str, "error", err)
		return false
	}
	return matched
}

func capture(pattern, str string) map[string]string {
	result := make(map[string]string)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(str)
	if matches == nil {
		return result
	}
	for i, name := range re.SubexpNames() {
		result[strconv.Itoa(i)] = matches[i]
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}
	return result
}

func strptime(format, timestr string) time.Time {
	t, err := strftime.Parse(format, timestr)
	if err != nil {
		slog.Error("strptime function error", "format", format, "timestr", timestr, "error", err)
		return time.Time{}
	}
	return t
}

func in(key any, target []any) bool {
	for _, v := range target {
		if v == key {
			return true
		}
	}
	return false
}
func tojson(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		slog.Error("toJSON function error", "value", v, "error", err)
		return ""
	}
	return string(b)
}

func toyaml(v any) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		slog.Error("toYAML function error", "value", v, "error", err)
		return ""
	}
	return string(b)
}

func toxml(v any) string {
	b, err := xml.Marshal(v)
	if err != nil {
		slog.Error("toXML function error", "value", v, "error", err)
		return ""
	}
	return string(b)
}

func rfc3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

func do_strftime(format string, t time.Time) string {
	return strftime.Format(format, t)
}

func do_hex(data []byte) string {
	return hex.EncodeToString(data)
}

func do_unhex(data string) []byte {
	if bout, err := hex.DecodeString(data); err != nil {
		return bout
	} else {
		slog.Error("unhex", "error", err, "input", data)
		return nil
	}
}

func do_base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func do_unbase64(data string) []byte {
	if bout, err := base64.StdEncoding.DecodeString(data); err != nil {
		return bout
	} else {
		slog.Error("unbase64", "error", err, "input", data)
		return nil
	}
}

func makefuncs() template.FuncMap {
	return template.FuncMap{
		"match":    match,
		"capture":  capture,
		"now":      time.Now,
		"rfc3339":  rfc3339,
		"strftime": do_strftime,
		"strptime": strptime,
		"in":       in,
		"toJSON":   tojson,
		"toYAML":   toyaml,
		"toXML":    toxml,
		"hex":      do_hex,
		"unhex":    do_unhex,
		"base64":   do_base64,
		"unbase64": do_unbase64,
	}
}
