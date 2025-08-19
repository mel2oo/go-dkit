package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	nm "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	nt "go.opentelemetry.io/otel/trace/noop"
)

var std *Provider = &Provider{
	opts:           newOptions(),
	MeterProvider:  nm.NewMeterProvider(),
	TracerProvider: nt.NewTracerProvider(),
	Propagators: propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	),
}

func Standard() *Provider {
	return std
}

type Provider struct {
	opts *Options

	MeterProvider  metric.MeterProvider
	TracerProvider trace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

func NewProvider(ctx context.Context, opts ...Option) (*Provider, error) {
	opt := newOptions()
	for _, v := range opts {
		v(opt)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(opt.serviceName),
		semconv.ServiceVersion(opt.serviceVersion),
	)

	meterProvider, err := newMeterProvider(ctx, res, opt)
	if err != nil {
		return nil, err
	}

	tracerProvider, err := newTracerProvider(ctx, res, opt)
	if err != nil {
		return nil, err
	}

	op := &Provider{
		opts:           opt,
		MeterProvider:  meterProvider,
		TracerProvider: tracerProvider,
		Propagators: propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	}

	if opt.setGlobal {
		otel.SetMeterProvider(op.MeterProvider)
		otel.SetTracerProvider(op.TracerProvider)
		otel.SetTextMapPropagator(op.Propagators)
	}

	if opt.setStandard {
		std = op
	}

	return op, nil
}
