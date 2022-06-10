package rpcx

import (
	"net"
	"time"

	prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	metrics "github.com/rcrowley/go-metrics"
	etcdplugin "github.com/rpcxio/rpcx-etcd/serverplugin"
	rpcx "github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Options struct {
	BasePath       string
	UpdateInterval time.Duration
	Address        []string
}

type InitServers func(server *rpcx.Server)

type Server struct {
	opt      *Options
	app      string
	server   *rpcx.Server
	logger   *zap.Logger
	initFunc InitServers
}

func NewServerOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("rpcx.server", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal rpcx server option error")
	}
	logger.Info("load rpcx options success", zap.Any("rpcx options", o))
	return o, nil
}

func NewServer(opt *Options, logger *zap.Logger, tp trace.TracerProvider, init InitServers) (*Server, error) {
	rx := rpcx.NewServer()
	// rx.DisableHTTPGateway = true

	rxMetrics := serverplugin.NewMetricsPlugin(metrics.DefaultRegistry)
	rx.Plugins.Add(rxMetrics)
	rpcxMetrics()
	tracer := tp.Tracer("rpcx")
	serverplugin.NewOpenTelemetryPlugin(tracer, nil)
	return &Server{
		opt:      opt,
		logger:   logger.With(zap.String("type", "rpcx")),
		server:   rx,
		initFunc: init,
	}, nil
}

func (s *Server) ApplicationName(name string) {
	s.app = name
}

func (s *Server) Start(ln net.Listener) error {
	s.logger.Info("rpc server starting ...")

	if err := s.register(ln.Addr().String()); err != nil {
		return errors.Wrap(err, "register rpc server error")
	}

	//初始化rpc服务器的provider
	s.initFunc(s.server)

	go func() {
		if err := s.server.ServeListener("tcp", ln); err != nil {
			s.logger.Fatal("failed to serve rpc: %v", zap.Error(err))
		}
	}()

	return nil
}

func rpcxMetrics() {
	metrics.RegisterRuntimeMemStats(metrics.DefaultRegistry)
	go metrics.CaptureRuntimeMemStats(metrics.DefaultRegistry, time.Second)
	prometheusClient := prometheusmetrics.NewPrometheusProvider(
		metrics.DefaultRegistry, "whatever", "something", prometheus.DefaultRegisterer, 1*time.Second)
	go prometheusClient.UpdatePrometheusMetrics()
	//addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
	//go graphite.Graphite(metrics.DefaultRegistry, 1e9, "rpcx.services.host.127_0_0_1", addr)
}

func (s *Server) register(addr string) error {
	r := &etcdplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + addr,
		EtcdServers:    s.opt.Address,
		BasePath:       s.opt.BasePath,
		UpdateInterval: s.opt.UpdateInterval * time.Millisecond,
		Metrics:        metrics.NewRegistry(),
	}

	if err := r.Start(); err != nil {
		return errors.Wrap(err, "register center error")
	}
	s.server.Plugins.Add(r)

	s.logger.Info("register rpc service success", zap.String("id", "tcp@"+addr))
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("rpc server stopping ... ")
	if err := s.server.Close(); err != nil {
		return errors.Wrap(err, "stop rpc server error")
	}
	return nil
}
