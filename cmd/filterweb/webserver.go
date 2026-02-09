package main

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wtnb75/go-filterweb"
)

type WebServer struct {
	Listen     string `long:"listen" description:"listen address" default:":3000"`
	configData []filterweb.ConfigSchema
}

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug("received request", "path", r.URL.Path, "method", r.Method)
	statuscode := http.StatusOK
	start := time.Now()
	defer func() {
		headers := []any{
			"remote", r.RemoteAddr, "elapsed", time.Since(start),
			"method", r.Method, "path", r.URL.Path,
			"status", statuscode, "protocol", r.Proto,
		}
		if r.URL.User.Username() != "" {
			headers = append(headers, "user", r.URL.User.Username())
		}
		for k, v := range w.Header() {
			switch strings.ToLower(k) {
			case "etag", "content-type", "content-encoding", "location":
				headers = append(headers, strings.ToLower(k), v[0])
			case "content-length":
				if val, err := strconv.Atoi(v[0]); err != nil {
					headers = append(headers, "length", v[0])
				} else {
					headers = append(headers, "length", val)
				}
			case "last-modified":
				if ts, err := time.Parse(http.TimeFormat, v[0]); err != nil {
					headers = append(headers, "last-modified", v[0])
				} else {
					headers = append(headers, "last-modified", ts)
				}
			}
		}
		for k, v := range r.Header {
			switch strings.ToLower(k) {
			case "x-forwarded-for", "x-forwarded-host", "x-forwarded-proto":
				headers = append(headers, strings.TrimPrefix(strings.ToLower(k), "x-"), v[0])
			case "forwarded", "user-agent", "if-none-match", "referer", "accept-encoding", "range":
				headers = append(headers, strings.ToLower(k), v[0])
			case "if-modified-since":
				if ts, err := time.Parse(http.TimeFormat, v[0]); err != nil {
					headers = append(headers, "if-modified-since", v[0])
				} else {
					headers = append(headers, "if-modified-since", ts)
				}
			}
		}
		slog.Info(
			http.StatusText(statuscode), headers...)
	}()
	for _, cfg := range s.configData {
		if cfg.Path == r.URL.Path && cfg.Method == r.Method {
			slog.Debug("matched config", "path", cfg.Path, "method", cfg.Method)
			fdata, err := filterweb.ProcessFilters(cfg.Filters)
			if err != nil {
				statuscode = http.StatusInternalServerError
				slog.Error("failed to process filters", "error", err)
				http.Error(w, "Internal Server Error", statuscode)
				return
			}
			if data, ok := fdata.Data.([]byte); ok {
				w.Header().Set("Content-Type", fdata.ContentType)
				w.WriteHeader(statuscode)
				_, err = w.Write(data)
				if err != nil {
					slog.Error("failed to write response data([]byte)", "error", err)
				}
				return
			} else if strdata, ok := fdata.Data.(string); ok {
				w.Header().Set("Content-Type", fdata.ContentType)
				w.WriteHeader(statuscode)
				_, err = w.Write([]byte(strdata))
				if err != nil {
					slog.Error("failed to write response data(string)", "error", err)
				}
				return
			} else {
				buf, err := filterweb.EncodeContentType(fdata.ContentType, fdata.Data)
				if err != nil {
					statuscode = http.StatusInternalServerError
					slog.Error("failed to encode response data", "error", err)
					http.Error(w, "Internal Server Error", statuscode)
					return
				}
				w.Header().Set("Content-Type", fdata.ContentType)
				w.WriteHeader(statuscode)
				_, err = w.Write(buf)
				if err != nil {
					slog.Error("failed to write response data", "error", err)
				}
			}
			return
		}
	}
	statuscode = http.StatusNotFound
	slog.Warn("no matching config found", "path", r.URL.Path, "method", r.Method)
	http.Error(w, "not found", statuscode)
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
