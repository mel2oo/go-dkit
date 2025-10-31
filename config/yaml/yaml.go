package yaml

import (
	"os"
	"regexp"

	"github.com/ghodss/yaml"
	"go-micro.dev/v5/config/encoder"
)

type yamlEncoder struct{}

func (y yamlEncoder) Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y yamlEncoder) Decode(d []byte, v interface{}) error {
	res := expandEnvDefaults(string(d))
	return yaml.Unmarshal([]byte(res), v)
}

func (y yamlEncoder) String() string {
	return "yaml"
}

func NewEncoder() encoder.Encoder {
	return yamlEncoder{}
}

var envRe = regexp.MustCompile(`\$\{([^:}]+)(:(.*))?\}`)

func expandEnvDefaults(input string) string {
	return envRe.ReplaceAllStringFunc(input, func(m string) string {
		sub := envRe.FindStringSubmatch(m)
		name := sub[1]
		def := ""
		if len(sub) >= 4 {
			def = sub[3]
		}
		if val, ok := os.LookupEnv(name); ok && val != "" {
			return val
		}
		return def
	})
}
