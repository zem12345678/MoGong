package grpc

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientOptions struct {
	Wait            time.Duration
	Tag             string
	GrpcDialOptions []grpc.DialOption
}

type Client struct {
	o *ClientOptions
}

func NewClientOptions(v *viper.Viper, logger *zap.Logger, tracer trace.TracerProvider) (*ClientOptions, error) {
	var (
		err error
		o   = new(ClientOptions)
	)
	if err = v.UnmarshalKey("grpc.client", o); err != nil {
		return nil, err
	}

	logger.Info("load grpc.client options success", zap.Any("grpc.client options", o))
	grpc_prometheus.EnableClientHandlingTimeHistogram()
	o.GrpcDialOptions = append(o.GrpcDialOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_prometheus.StreamClientInterceptor,
			grpc_zap.StreamClientInterceptor(logger),
			otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(tracer)),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_prometheus.UnaryClientInterceptor,
			grpc_zap.UnaryClientInterceptor(logger),
			otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(tracer)),
		)),
	)
	return o, nil
}

// ClientOptional grpc client optional
type ClientOptional func(o *ClientOptions)

// WithTimeout grpc client time out
func WithTimeout(d time.Duration) ClientOptional {
	return func(o *ClientOptions) {
		o.Wait = d
	}
}

// WithTag grpc client tag
func WithTag(tag string) ClientOptional {
	return func(o *ClientOptions) {
		o.Tag = tag
	}
}

// WithGrpcDialOptions grpc dial option
func WithGrpcDialOptions(options ...grpc.DialOption) ClientOptional {
	return func(o *ClientOptions) {
		o.GrpcDialOptions = append(o.GrpcDialOptions, options...)
	}
}

// NewClient new grpc client server
func NewClient(o *ClientOptions) (*Client, error) {
	return &Client{
		o: o,
	}, nil
}

// Dial grpc client dail
func (c *Client) Dial(service string, options ...ClientOptional) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	o := &ClientOptions{
		Wait:            c.o.Wait,
		Tag:             c.o.Tag,
		GrpcDialOptions: c.o.GrpcDialOptions,
	}

	for _, option := range options {
		option(o)
	}
	conn, err := grpc.DialContext(ctx, service, o.GrpcDialOptions...)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial error")
	}
	return conn, nil
}
