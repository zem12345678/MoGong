package rpcx

import (
	"context"

	rpcxEtcd "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
)

type ClientOptions struct {
	BasePath string
	Address  []string
	tracer   trace.Tracer
}

func NewClientOptions(v *viper.Viper, tp trace.TracerProvider) (*ClientOptions, error) {
	opt := &ClientOptions{}

	if err := v.UnmarshalKey("rpcx.client", opt); err != nil {
		return nil, err
	}
	tracer := tp.Tracer("rpcx")
	opt.tracer = tracer
	return opt, nil
}

type Client struct {
	opt *ClientOptions
}

func NewClient(opt *ClientOptions) (*Client, error) {
	return &Client{opt: opt}, nil
}

func (c *Client) Dial(service string) (*client.XClient, error) {

	discovery, err := rpcxEtcd.NewEtcdV3Discovery(c.opt.BasePath, service, c.opt.Address, false, nil)
	if err != nil {
		return nil, err
	}
	xClient := client.NewXClient(service, client.Failover, client.WeightedRoundRobin, discovery, client.DefaultOption)
	p := client.NewOpenTelemetryPlugin(c.opt.tracer, nil)
	pc := client.NewPluginContainer()
	pc.Add(p)
	xClient.SetPlugins(pc)
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{c.opt.BasePath: "from client"})
	_ = context.WithValue(ctx, share.ResMetaDataKey, make(map[string]string))
	return &xClient, nil
}
