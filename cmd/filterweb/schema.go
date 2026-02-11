package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/goccy/go-yaml"
	"github.com/invopop/jsonschema"
	"github.com/wtnb75/go-filterweb"
)

type Schema struct {
	Format string `long:"format" choice:"yaml" choice:"json" default:"json"`
}

func (s *Schema) Execute(args []string) error {
	init_log()
	for _, name := range filterweb.ListFilters() {
		f, err := filterweb.GetFilter(name)
		slog.Info("filter", "name", name, "accept", f.Accepts())
		if err != nil {
			slog.Error("failed to get filter", "name", name, "error", err)
			return err
		} else {
			r := jsonschema.Reflector{}
			scm := r.Reflect(f)
			var res []byte
			var err error
			switch s.Format {
			case "yaml":
				res, err = yaml.Marshal(scm)
			case "json":
				res, err = json.Marshal(scm)
			}
			if err != nil {
				slog.Error("failed to marshal schema", "name", name, "error", err)
				return err
			} else {
				fmt.Println(string(res))
			}
		}
	}
	return nil
}
