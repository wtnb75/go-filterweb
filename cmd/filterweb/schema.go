package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/invopop/jsonschema"
	"github.com/wtnb75/go-filterweb"
)

type Schema struct{}

func (s *Schema) Execute(args []string) error {
	init_log()
	for _, name := range filterweb.ListFilters() {
		f, err := filterweb.GetFilter(name)
		slog.Info("filter", "name", name, "accept", f.Accepts())
		if err != nil {
			slog.Error("failed to get filter", "name", name, "error", err)
			return err
		} else {
			scm := jsonschema.Reflect(f)
			res, err := json.Marshal(scm)
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
