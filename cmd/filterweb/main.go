package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jessevdk/go-flags"
	"github.com/wtnb75/go-filterweb"
)

var globalOption struct {
	Verbose bool           `short:"v" long:"verbose" description:"show verbose logs"`
	Quiet   bool           `short:"q" long:"quiet" description:"suppress logs"`
	TextLog bool           `long:"text-log" description:"use text format for logging"`
	Config  flags.Filename `short:"c" long:"config" description:"config file" env:"FILTERWEB_CONFIG"`
}

func init_log() {
	var level = slog.LevelInfo
	if globalOption.Verbose {
		level = slog.LevelDebug
	} else if globalOption.Quiet {
		level = slog.LevelWarn
	}
	slog.SetLogLoggerLevel(level)
	if !globalOption.TextLog {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
	}
}

func load_config(fn string) ([]filterweb.ConfigSchema, error) {
	var configData []filterweb.ConfigSchema
	slog.Debug("config file", "path", globalOption.Config)
	if f, err := os.Open(string(globalOption.Config)); err == nil {
		defer f.Close()
		if data, err := io.ReadAll(f); err == nil {
			if err = yaml.Unmarshal(data, &configData); err != nil {
				slog.Error("failed to parse config file", "error", err)
				return nil, err
			}
		} else {
			slog.Error("failed to read config file", "error", err)
			return nil, err
		}
	} else {
		slog.Error("failed to open config file", "error", err)
		return nil, err
	}
	return configData, nil
}

type SubCommand struct {
	Name  string
	Short string
	Long  string
	Data  any
}

func realMain() int {
	var err error
	parser := flags.NewParser(&globalOption, flags.Default)
	commands := []SubCommand{
		{Name: "checkfilter", Short: "check", Long: "check filter", Data: &CheckFilter{}},
		{Name: "schema", Short: "schema", Long: "show schema", Data: &Schema{}},
		{Name: "webserver", Short: "webserver", Long: "run webserver", Data: &WebServer{}},
	}
	for _, cmd := range commands {
		_, err = parser.AddCommand(cmd.Name, cmd.Short, cmd.Long, cmd.Data)
		if err != nil {
			slog.Error(cmd.Name, "error", err)
			return -1
		}
	}
	if _, err = parser.Parse(); err != nil {
		if _, ok := err.(*flags.Error); ok {
			return 0
		}
		slog.Error("error exit", "error", err)
		parser.WriteHelp(os.Stdout)
		return 1
	}
	return 0
}

func main() {
	os.Exit(realMain())
}
