package filterweb

import (
	"log/slog"

	"github.com/go-viper/mapstructure/v2"
	"github.com/itchyny/gojq"
)

type JqConfig struct {
	Filter
	Expression string
	query      *gojq.Query
}

func (jc *JqConfig) New() Filter {
	return &JqConfig{}
}

func (jc *JqConfig) Name() string {
	return "jq"
}

func (jc *JqConfig) Accepts() []string {
	return []string{}
}

func (jc *JqConfig) Prep(config Config, data Data) error {
	err := mapstructure.Decode(config.Params, jc)
	if err != nil {
		return err
	}
	if jc.Expression == "" {
		slog.Error("jq filter requires 'expression' parameter")
		return ErrMissingParams
	}
	jc.query, err = gojq.Parse(jc.Expression)
	return nil
}

func (jc *JqConfig) Process(data Data) (Data, error) {
	iter := jc.query.Run(data.Data)
	res := []any{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			slog.Error("jq processing error", "error", err)
			return Data{}, err
		}
		res = append(res, v)
	}
	return Data{ContentType: "application/json", Data: res}, nil
}

func (jc *JqConfig) Post(config Config, data Data) error {
	return nil
}

func init() {
	RegisterFilter(&JqConfig{})
}
