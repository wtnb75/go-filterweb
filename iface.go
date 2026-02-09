package filterweb

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Name   string
	Params map[string]any
}

type Data struct {
	ContentType string
	Data        any
}

type ConfigSchema struct {
	Path    string
	Method  string
	Filters []Config
}

type Filter interface {
	Name() string
	New() Filter
	Accepts() []string
	Prep(config Config, data Data) error
	Process(data Data) (Data, error)
	Post(config Config, data Data) error
}

var filters map[string]Filter = make(map[string]Filter)

func RegisterFilter(f Filter) {
	filters[f.Name()] = f
}

func GetFilter(name string) (Filter, error) {
	f, ok := filters[name]
	if !ok {
		return nil, ErrFilterNotFound
	}
	return f.New(), nil
}

func ListFilters() []string {
	var names []string
	for name := range filters {
		names = append(names, name)
	}
	return names
}

func ProcessFilters(configs []Config) (Data, error) {
	var data = Data{}
	for _, config := range configs {
		slog.Debug("processing filter", "config", config, "data", data)
		filter, err := GetFilter(config.Name)
		if err != nil {
			return data, err
		}
		accepts := filter.Accepts()
		accepted := false
		if len(accepts) == 0 {
			accepted = true
		}
		for _, accept := range accepts {
			if accept == data.ContentType || accept == "*" {
				accepted = true
				break
			}
		}
		if !accepted {
			slog.Error("filter does not accept content type", "filter", filter, "data", data)
			return data, ErrContentTypeMismatch
		}
		slog.Debug("Prep", "name", filter.Name(), "filter", filter, "data", data)
		err = filter.Prep(config, data)
		if err != nil {
			return data, err
		}
		slog.Debug("Process", "name", filter.Name(), "filter", filter, "data", data)
		data, err = filter.Process(data)
		if err != nil {
			return data, err
		}
		slog.Debug("Post", "name", filter.Name(), "filter", filter, "data", data)
		err = filter.Post(config, data)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

func DecodeContentType(contentType string, data []byte) (any, error) {
	slog.Debug("decodeContentType", "contentType", contentType, "data", data)
	var res any
	switch contentType {
	case "application/json":
		err := json.Unmarshal(data, res)
		if err != nil {
			return res, err
		}
	case "application/yaml", "text/yaml":
		err := yaml.Unmarshal(data, &res)
		if err != nil {
			return res, err
		}
	case "text/xml", "application/xml":
		err := xml.Unmarshal(data, &res)
		if err != nil {
			return res, err
		}
	case "text/csv":
		buf := bytes.NewReader(data)
		reader := csv.NewReader(buf)
		hdr, err := reader.Read()
		if err != nil {
			slog.Error("csv read error", "error", err)
			return res, err
		}
		resultdata := make([]map[string]any, 0)
		for {
			row, err := reader.Read()
			if err != nil {
				break
			}
			m := make(map[string]any)
			for i, v := range row {
				m[hdr[i]] = v
			}
			resultdata = append(resultdata, m)
		}
		res = resultdata
	default:
		res = data
	}
	slog.Debug("decoded", "contentType", contentType, "data", res)
	return res, nil
}

func EncodeContentType(contentType string, data any) ([]byte, error) {
	slog.Debug("encodeContentType", "contentType", contentType, "data", data)
	var res []byte
	var err error
	switch contentType {
	case "application/json":
		res, err = json.Marshal(data)
		if err != nil {
			return res, err
		}
	case "application/yaml", "text/yaml":
		res, err = yaml.Marshal(data)
		if err != nil {
			return res, err
		}
	case "text/xml", "application/xml":
		res, err = xml.Marshal(data)
		if err != nil {
			return res, err
		}
	case "text/csv":
		buf := &bytes.Buffer{}
		writer := csv.NewWriter(buf)
		records, ok := data.([]map[string]any)
		if !ok || len(records) == 0 {
			return nil, fmt.Errorf("data is not []map[string]any or empty")
		}
		// write header
		var header []string
		for k := range records[0] {
			header = append(header, k)
		}
		err = writer.Write(header)
		if err != nil {
			slog.Error("csv write error", "error", err)
			return nil, err
		}
		// write records
		for _, record := range records {
			var row []string
			for _, k := range header {
				v, ok := record[k]
				if !ok {
					v = ""
				}
				row = append(row, fmt.Sprintf("%v", v))
			}
			err = writer.Write(row)
			if err != nil {
				slog.Error("csv write error", "error", err)
				return nil, err
			}
		}
		writer.Flush()
		res = buf.Bytes()
	default:
		res = fmt.Appendf(nil, "%v", data)
	}
	slog.Debug("encoded", "contentType", contentType, "data", string(res))
	return res, nil
}
