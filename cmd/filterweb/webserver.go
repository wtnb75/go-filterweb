package main

import (
	"log/slog"
	"net/http"

	"github.com/wtnb75/go-filterweb"
)

type WebServer struct {
	Listen     string `long:"listen" description:"listen address" default:":3000"`
	configData []filterweb.ConfigSchema
}

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("received request", "path", r.URL.Path, "method", r.Method)
	for _, cfg := range s.configData {
		if cfg.Path == r.URL.Path && cfg.Method == r.Method {
			slog.Debug("matched config", "path", cfg.Path, "method", cfg.Method)
			fdata, err := filterweb.ProcessFilters(cfg.Filters)
			if err != nil {
				slog.Error("failed to process filters", "error", err)
				http.Error(w, "Internal Server Error", 500)
				return
			}
			if data, ok := fdata.Data.([]byte); ok {
				w.Header().Set("Content-Type", fdata.ContentType)
				w.WriteHeader(200)
				_, err = w.Write(data)
				if err != nil {
					slog.Error("failed to write response data([]byte)", "error", err)
				}
				return
			} else if strdata, ok := fdata.Data.(string); ok {
				w.Header().Set("Content-Type", fdata.ContentType)
				w.WriteHeader(200)
				_, err = w.Write([]byte(strdata))
				if err != nil {
					slog.Error("failed to write response data(string)", "error", err)
				}
				return
			} else {
				buf, err := filterweb.EncodeContentType(fdata.ContentType, fdata.Data)
				if err != nil {
					slog.Error("failed to encode response data", "error", err)
					http.Error(w, "Internal Server Error", 500)
					return
				}
				w.Header().Set("Content-Type", fdata.ContentType)
				w.WriteHeader(200)
				_, err = w.Write(buf)
				if err != nil {
					slog.Error("failed to write response data", "error", err)
				}
			}
			return
		}
	}
	slog.Warn("no matching config found", "path", r.URL.Path, "method", r.Method)
	http.Error(w, "not found", 404)
}

func (s *WebServer) Execute(args []string) error {
	init_log()
	var err error
	s.configData, err = load_config(string(globalOption.Config))
	if err != nil {
		slog.Error("fail to load config file", "path", globalOption.Config, "error", err)
		return err
	}
	hdl := http.NewServeMux()
	hdl.Handle("/", s)
	srv := http.Server{
		Addr:    s.Listen,
		Handler: hdl,
	}
	slog.Info("starting webserver", "address", srv.Addr)
	return srv.ListenAndServe()
}
