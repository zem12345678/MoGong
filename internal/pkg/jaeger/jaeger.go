package jaeger

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Options struct {
	url         string
	service     string
	environment string
	id          int64
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)

	if err = v.UnmarshalKey("jaeger", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal jaeger configuration error")
	}
	o.service = v.GetString("app.name")
	o.environment = v.GetString("app.environment")
	logger.Info("load jaeger configuration success")

	return o, nil
}

func New(o *Options) (trace.TracerProvider, error) {

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.url)))
	if err != nil {
		return nil, errors.Wrap(err, "create jaeger tracer error")
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(o.service),
			attribute.String("environment", o.environment),
			attribute.Int64("ID", o.id),
		)),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tp)
	return tp, nil
}
