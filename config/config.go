package config

import (
	"github.com/mel2oo/go-dkit/config/consul"
	"github.com/mel2oo/go-dkit/config/yaml"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/config/reader"
	"go-micro.dev/v5/config/reader/json"
	"go-micro.dev/v5/config/source"
	"go-micro.dev/v5/config/source/file"
)

func New(filename string, cbRefresh func([]byte) error, opts ...Option) error {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}

	var src source.Source
	var err error
	var encoder = yaml.NewEncoder()

	if len(opt.consulAddress) > 0 {
		src = consul.NewSource(
			consul.WithAddress(opt.consulAddress),
			consul.WithToken(opt.consulToken),
			consul.WithPrefix(filename),
			consul.StripPrefix(true),
			source.WithEncoder(encoder),
		)
	} else {
		src = file.NewSource(
			file.WithPath(filename),
			source.WithEncoder(encoder),
		)
	}

	config, err := config.NewConfig(
		config.WithReader(
			json.NewReader(
				reader.WithEncoder(yaml.NewEncoder()),
			),
		),
	)
	if err != nil {
		return err
	}

	if err := config.Load(src); err != nil {
		return err
	}

	if err := cbRefresh(config.Bytes()); err != nil {
		return err
	}

	watcher, err := config.Watch()
	if err != nil {
		return err
	}

	go func() {
		for {
			v, err := watcher.Next()
			if err != nil {
				logrus.Errorf("watch config error: %s", err)
				return
			}

			if string(v.Bytes()) == "null" {
				continue
			}

			if err := cbRefresh(v.Bytes()); err != nil {
				logrus.Errorf("refresh config error: %s", err)
			}
		}
	}()

	return nil
}
