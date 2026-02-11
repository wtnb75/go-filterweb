package filterweb

import (
	"context"
	"io"
	"log/slog"
	"mime"
	"net"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
)

type HTTPConfig struct {
	Filter
	Url         string            // request URL
	UnixSocket  string            // unix socket path
	ContentType string            // override content type
	Verify      bool              // verify TLS certificates
	Method      string            // HTTP method
	Headers     map[string]string // HTTP headers
	ExpectCode  []int             // expected HTTP status code
}

func (hc *HTTPConfig) New() Filter {
	return &HTTPConfig{}
}

func (hc *HTTPConfig) Name() string {
	return "http"
}

func (hc *HTTPConfig) Accepts() []string {
	return []string{}
}

func (hc *HTTPConfig) Prep(config Config, data Data) error {
	// defaults
	hc.Method = http.MethodGet
	hc.Verify = true
	hc.ExpectCode = []int{200}
	err := mapstructure.Decode(config.Params, hc)
	if err != nil {
		return err
	}
	// mandatory
	if hc.Url == "" {
		slog.Error("http filter requires 'url' parameter")
		return ErrMissingParams
	}
	return nil
}

func (hc *HTTPConfig) Process(data Data) (Data, error) {
	res := Data{}
	var client *http.Client
	if hc.UnixSocket != "" {
		client = &http.Client{Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", hc.UnixSocket)
			},
		}}
	} else {
		client = &http.Client{}
	}
	httpreq, err := http.NewRequest(hc.Method, hc.Url, nil)
	if err != nil {
		slog.Error("http request(prep)", "method", hc.Method, "url", hc.Url, "err", err)
		return res, ErrHTTPRequestFailed
	}
	for key, value := range hc.Headers {
		httpreq.Header.Add(key, value)
	}
	httpres, err := client.Do(httpreq)
	if err != nil {
		slog.Error("http request(do)", "method", hc.Method, "url", hc.Url, "err", err)
		return res, ErrHTTPRequestFailed
	}
	defer httpres.Body.Close()
	success := false
	for _, v := range hc.ExpectCode {
		if httpres.StatusCode == v {
			success = true
			break
		}
	}
	if !success {
		slog.Error("status code", "expected", hc.ExpectCode, "actual", httpres.StatusCode)
		return res, ErrHTTPStatusNotOK
	}
	if hc.ContentType == "" {
		ct := httpres.Header.Get("Content-Type")
		if ct != "" {
			mediaType, _, err := mime.ParseMediaType(ct)
			slog.Debug("Parsed media type", "mediaType", mediaType, "err", err)
			if err != nil {
				return res, err
			}
			res.ContentType = mediaType
		}
	} else {
		res.ContentType = hc.ContentType
	}
	buf, err := io.ReadAll(httpres.Body)
	if err != nil {
		slog.Error("read body", "method", hc.Method, "url", hc.Url, "err", err)
		return res, ErrHTTPRequestFailed
	}
	res.Data, err = DecodeContentType(res.ContentType, buf)
	if err != nil {
		slog.Error("decode", "method", hc.Method, "url", hc.Url, "contenttype", res.ContentType, "err", err, "data", res.Data)
	}
	return res, err
}

func (hc *HTTPConfig) Post(config Config, data Data) error {
	return nil
}

func init() {
	RegisterFilter(&HTTPConfig{})
}
