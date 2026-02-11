package main

import (
	"fmt"
	"log/slog"

	"github.com/wtnb75/go-filterweb"
)

type CheckFilter struct {
	HideCt bool `long:"hide-content-type"`
}

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
		if !cf.HideCt {
			fmt.Printf("%s %s\n", config.Method, config.Path)
		}
		if err != nil {
			slog.Error("failed to process filters", "error", err)
			return err
		}
		if !cf.HideCt {
			fmt.Printf("Content-Type: %s\n\n", fdata.ContentType)
		}
		fmt.Print(fdata)
	}
	return nil
}
