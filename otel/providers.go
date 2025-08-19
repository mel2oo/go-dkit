package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// newMeterProvider creates a new meter provider with the OTLP gRPC exporter.
func newMeterProvider(ctx context.Context, res *resource.Resource,
	opts *Options) (*metric.MeterProvider, error) {
	newOpts := []otlpmetricgrpc.Option{}
	newOpts = append(newOpts, otlpmetricgrpc.WithInsecure())
	if len(opts.endpoint) > 0 {
		newOpts = append(newOpts, otlpmetricgrpc.WithEndpoint(opts.endpoint))
	}

	exporter, err := otlpmetricgrpc.New(ctx, newOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	return mp, nil
}

// newTracerProvider creates a new tracer provider with the OTLP gRPC exporter.
func newTracerProvider(ctx context.Context, res *resource.Resource,
	opts *Options) (*trace.TracerProvider, error) {
	newOpts := []otlptracegrpc.Option{}
	newOpts = append(newOpts, otlptracegrpc.WithInsecure())
	if len(opts.endpoint) > 0 {
		newOpts = append(newOpts, otlptracegrpc.WithEndpoint(opts.endpoint))
	}

	exporter, err := otlptracegrpc.New(ctx, newOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Create Resource
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	return tp, nil
}
