package filterweb

import (
	"log/slog"

	"github.com/go-viper/mapstructure/v2"
)

type EncodeConfig struct {
	Filter
	ContentType string // target content type
}

func (ec *EncodeConfig) New() Filter {
	return &EncodeConfig{}
}

func (ec *EncodeConfig) Name() string {
	return "encode"
}

func (ec *EncodeConfig) Accepts() []string {
	return []string{}
}

func (ec *EncodeConfig) Prep(config Config, data Data) error {
	err := mapstructure.Decode(config.Params, ec)
	if err != nil {
		return err
	}
	// mandatory
	if ec.ContentType == "" {
		slog.Error("encode filter requires 'content_type' parameter")
		return ErrMissingParams
	}
	return nil
}

func (ec *EncodeConfig) Process(data Data) (Data, error) {
	res := Data{}
	var err error
	res.ContentType = ec.ContentType
	res.Data, err = EncodeContentType(ec.ContentType, data.Data)
	return res, err
}

func (ec *EncodeConfig) Post(config Config, data Data) error {
	return nil
}

func init() {
	RegisterFilter(&EncodeConfig{})
}
