package yaml

import (
	"github.com/drone/envsubst/v2"
	"github.com/ghodss/yaml"
	"go-micro.dev/v5/config/encoder"
)

type yamlEncoder struct{}

func (y yamlEncoder) Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y yamlEncoder) Decode(d []byte, v interface{}) error {
	res, err := envsubst.EvalEnv(string(d))
	if err != nil {
		return err
	}
	return yaml.Unmarshal([]byte(res), v)
}

func (y yamlEncoder) String() string {
	return "yaml"
}

func NewEncoder() encoder.Encoder {
	return yamlEncoder{}
}
