package grpc

import (
	"fmt"
	"log"
	discovery "mogong/internal/pkg/etcd"
	nettools "mogong/internal/pkg/tools/net"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServerOptions struct {
	Port        int
	EtcdAddr    []string
	ServiceName string
}

type Server struct {
	o      *ServerOptions
	app    string
	host   string
	port   int
	logger *zap.Logger
	server *grpc.Server
}

func NewServerOptions(v *viper.Viper, logger *zap.Logger) (*ServerOptions, error) {
	var (
		err error
		o   = new(ServerOptions)
	)
	if err = v.UnmarshalKey("grpc", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal grpc server option error")
	}

	logger.Info("load grpc options success", zap.Any("grpc options", o))

	return o, nil
}

type InitServers func(server *grpc.Server)

func NewServer(o *ServerOptions, logger *zap.Logger, init InitServers, tracer trace.TracerProvider) (*Server, error) {
	var gs *grpc.Server
	logger = logger.With(zap.String("type", "grpc"))
	gs = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tracer)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tracer)),
		)),
	)
	init(gs)
	grpc_health_v1.RegisterHealthServer(gs, health.NewServer())
	return &Server{
		o:      o,
		logger: logger.With(zap.String("type", "grpc.Server")),
		server: gs,
	}, nil
}

// Application 服务应用
func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	s.port = s.o.Port
	if s.port == 0 {
		s.port = nettools.GetAvailablePort()
	}

	s.host = nettools.GetLocalIP4()

	if s.host == "" {
		return errors.New("get local ipv4 error")
	}
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	s.logger.Info("grpc server starting ...", zap.String("addr", addr))

	//将服务地址注册到etcd中
	etcdRegister := discovery.NewRegister(s.o.EtcdAddr, s.logger)
	defer etcdRegister.Stop()
	node := discovery.Server{
		Name: s.app,
		Addr: addr,
	}
	go etcdRegister.Register(node, 10)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sig := <-ch
		etcdRegister.Unregister()
		if i, ok := sig.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()
	go func() {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Fatal("failed to listen: %v", zap.Error(err))
		}
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}

		if err := s.server.Serve(lis); err != nil {
			s.logger.Fatal("failed to serve: %v", zap.Error(err))
		}
	}()

	return nil
}

// Stop  停止GRPC服务
func (s *Server) Stop() error {
	s.logger.Info("grpc server stopping ...")
	s.server.GracefulStop()
	return nil
}
