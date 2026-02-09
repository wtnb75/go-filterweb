package main

import (
	"fmt"
	"log/slog"

	"github.com/wtnb75/go-filterweb"
)

type CheckFilter struct{}

func (cf *CheckFilter) Execute(args []string) error {
	init_log()
	slog.Debug("config file", "path", globalOption.Config)
	configData, err := load_config(string(globalOption.Config))
	if err != nil {
		slog.Error("fail to load config file", "path", globalOption.Config, "error", err)
		return err
	}
	for _, config := range configData {
		fdata, err := filterweb.ProcessFilters(config.Filters)
		if err != nil {
			slog.Error("failed to process filters", "error", err)
			return err
		}
		fmt.Printf("Content-Type: %s\n\n%s", fdata.ContentType, fdata.Data)
	}
	return nil
}
