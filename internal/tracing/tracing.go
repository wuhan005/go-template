// Copyright 2025 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tracing

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/wuhan005/go-template/internal/conf"
)

func SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)
	return
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(conf.Tracing.Endpoint),
		otlptracegrpc.WithInsecure(),
	}
	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	r, err := resource.New(ctx, []resource.Option{
		resource.WithAttributes(
			attribute.KeyValue{Key: "token", Value: attribute.StringValue(conf.Tracing.Token)},
			attribute.KeyValue{Key: "service.name", Value: attribute.StringValue(conf.Tracing.ServiceName)},
			attribute.KeyValue{Key: "host.name", Value: attribute.StringValue(conf.Tracing.HostName)}, // <hostName>替换为IP地址
		),
	}...)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}
