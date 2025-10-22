package config

type Option func(*Options)

type Options struct {
	consulAddress string
	consulToken   string
}

func WithConsul(address, token string) Option {
	return func(o *Options) {
		o.consulAddress = address
		o.consulToken = token
	}
}
