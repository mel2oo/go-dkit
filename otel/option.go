package otel

type Option func(*Options)

type Options struct {
	serviceName    string
	serviceVersion string
	endpoint       string
	setGlobal      bool
	setStandard    bool
}

func newOptions() *Options {
	return &Options{
		serviceName:    "go-otel",
		serviceVersion: "v0.1.0",
	}
}

func WithServiceName(name string) Option {
	return func(o *Options) { o.serviceName = name }
}

func WithServiceVersion(version string) Option {
	return func(o *Options) { o.serviceVersion = version }
}

func WithEndponit(endpoint string) Option {
	return func(o *Options) {
		if endpoint != "" {
			o.endpoint = endpoint
		}
	}
}

func SetGlobal() Option {
	return func(o *Options) { o.setGlobal = true }
}

func SetStandard() Option {
	return func(o *Options) { o.setStandard = true }
}
