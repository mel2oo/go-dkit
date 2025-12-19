package register

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/micro/plugins/v5/registry/consul"
	"go-micro.dev/v5/registry"
)

const (
	RegisterTypeConsul = "consul"
)

type Config struct {
	Type     string `yaml:"type" json:"type,omitempty"`
	Host     string `yaml:"host" json:"host,omitempty"`
	Port     int    `yaml:"port" json:"port,omitempty"`
	Token    string `yaml:"token" json:"token,omitempty"`
	Username string `yaml:"username" json:"username,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
}

func NewMicroRegister(cfg Config) registry.Registry {
	switch cfg.Type {
	case RegisterTypeConsul:
		config := api.DefaultNonPooledConfig()
		config.Address = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		config.Token = cfg.Token
		return consul.NewRegistry(consul.Config(config))
	default:
		return registry.NewMemoryRegistry()
	}
}
