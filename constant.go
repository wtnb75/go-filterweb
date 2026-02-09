package filterweb

import (
	"log/slog"

	"github.com/go-viper/mapstructure/v2"
)

type ConstantConfig struct {
	Filter
	ContentType string // target content type
	Data        any    // constant data
}

func (hc *ConstantConfig) New() Filter {
	return &ConstantConfig{}
}

func (hc *ConstantConfig) Name() string {
	return "constant"
}

func (hc *ConstantConfig) Accepts() []string {
	return []string{}
}

func (hc *ConstantConfig) Prep(config Config, data Data) error {
	// defaults
	hc.ContentType = "text/plain"
	err := mapstructure.Decode(config.Params, hc)
	if err != nil {
		return err
	}
	// mandatory
	if hc.Data == "" || hc.Data == nil {
		slog.Error("constant filter requires 'data' parameter")
		return ErrMissingParams
	}
	return nil
}

func (hc *ConstantConfig) Process(data Data) (Data, error) {
	res := Data{}
	res.ContentType = hc.ContentType
	var err error
	if strdata, ok := hc.Data.(string); ok {
		res.Data, err = DecodeContentType(hc.ContentType, []byte(strdata))
	} else if bindata, ok := hc.Data.([]byte); ok {
		res.Data, err = DecodeContentType(hc.ContentType, bindata)
	} else {
		res.Data = hc.Data
	}
	return res, err
}

func (hc *ConstantConfig) Post(config Config, data Data) error {
	return nil
}

func init() {
	RegisterFilter(&ConstantConfig{})
}
