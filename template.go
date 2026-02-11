package filterweb

import (
	"bytes"
	tmplHtml "html/template"
	"io"
	"log/slog"
	"maps"
	"os"
	tmplText "text/template"

	"github.com/go-viper/mapstructure/v2"
)

type TemplateConfig struct {
	Filter
	Type        string             // template type: text or html
	File        string             // template file path
	Content     string             // template content
	ContentType string             // output content type
	Vars        map[string]any     // template variables
	BaseKey     string             // base key for variables in input data
	tmpltxt     *tmplText.Template // text template
	tmplhtml    *tmplHtml.Template // html template
}

func (tc *TemplateConfig) New() Filter {
	return &TemplateConfig{}
}

func (tc *TemplateConfig) Name() string {
	return "template"
}

func (tc *TemplateConfig) Accepts() []string {
	return []string{"application/json", "application/yaml", "text/yaml", "text/xml", "application/xml", "text/dotenv"}
}

func (tc *TemplateConfig) load(tmpl string) (err error) {
	funcs := makefuncs()
	switch tc.Type {
	case "text":
		tc.tmplhtml = nil
		tc.tmpltxt, err = tmplText.New("template").Funcs(funcs).Parse(tmpl)
	case "html":
		tc.tmpltxt = nil
		tc.tmplhtml, err = tmplHtml.New("template").Funcs(funcs).Parse(tmpl)
	default:
		slog.Error("unsupported template type", "type", tc.Type)
		return ErrReadTemplate
	}
	return
}

func (tc *TemplateConfig) Prep(config Config, data Data) error {
	// defaults
	tc.Type = "text"
	err := mapstructure.Decode(config.Params, tc)
	if err != nil {
		return err
	}
	// defaults
	if tc.ContentType == "" {
		if tc.Type == "html" {
			tc.ContentType = "text/html"
		} else {
			tc.ContentType = "text/plain"
		}
	}
	// mandatory
	if tc.Content != "" {
		return tc.load(tc.Content)
	} else if tc.File != "" {
		f, err := os.Open(tc.File)
		if err != nil {
			return err
		}
		defer f.Close()
		buf, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		content := string(buf)
		return tc.load(content)
	} else {
		slog.Error("template filter requires 'file' or 'content' parameter")
		return ErrMissingParams
	}
}

func (tc *TemplateConfig) Process(data Data) (Data, error) {
	res := Data{ContentType: tc.ContentType}
	wr := &bytes.Buffer{}
	if tc.BaseKey != "" {
		data.Data = map[string]any{tc.BaseKey: data.Data}
	}
	if dataMap, ok := data.Data.(map[string]any); ok {
		maps.Copy(dataMap, tc.Vars)
	} else if len(tc.Vars) > 0 {
		slog.Warn("ignore vars: input data is not a map", "vars", tc.Vars)
	}
	if tc.tmplhtml != nil {
		err := tc.tmplhtml.Execute(wr, data.Data)
		if err != nil {
			return res, err
		}
	} else if tc.tmpltxt != nil {
		err := tc.tmpltxt.Execute(wr, data.Data)
		if err != nil {
			return res, err
		}
	}
	res.Data = wr.String()
	return res, nil
}

func (tc *TemplateConfig) Post(config Config, data Data) error {
	return nil
}

func init() {
	RegisterFilter(&TemplateConfig{})
}
