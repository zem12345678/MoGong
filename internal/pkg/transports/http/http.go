package http

import (
	"context"
	"fmt"
	nettool "mogong/internal/pkg/tools/net"
	"mogong/internal/pkg/transports/http/middleware/ginprometheus"
	"mogong/internal/pkg/transports/http/middleware/log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Options struct {
	Host string
	Port int
	Mode string
}

type Server struct {
	o          *Options
	app        string
	host       string
	port       int
	logger     *zap.Logger
	router     *gin.Engine
	httpServer http.Server
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal http option error")
	}

	logger.Info("load http options success", zap.Any("http options", o))
	return o, err
}

type InitControllers func(r *gin.Engine)

func NewRouter(o *Options, logger *zap.Logger, init InitControllers, tracer trace.TracerProvider) *gin.Engine {
	gin.SetMode(o.Mode)
	r := gin.New()
	// 跨域
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.Use(cors.Default())
	// panic之后自动恢复
	r.Use(gin.Recovery())
	// 日志格式化
	r.Use(log.Ginzap(logger, time.RFC3339, true))
	// panic日志格式化
	r.Use(ginzap.RecoveryWithZap(logger, true))
	// 添加prometheus 监控
	r.Use(ginprometheus.New(r).Middleware())
	r.Use(otelgin.Middleware("http-sever", otelgin.WithTracerProvider(tracer)))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	pprof.Register(r)
	init(r)
	return r
}

func New(o *Options, logger *zap.Logger, router *gin.Engine) (*Server, error) {
	return &Server{
		logger: logger.With(zap.String("type", "http.server")),
		router: router,
		o:      o,
	}, nil
}

// Application set app name
func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	s.port = s.o.Port
	if s.port == 0 {
		s.port = nettool.GetAvailablePort()
	}

	s.host = s.o.Host

	if s.host == "" {
		// return errors.New("get local ipv4 error")

		s.host = nettool.GetLocalIP4()
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	s.logger.Info("http server starting ...", zap.String("addr", addr))

	s.httpServer = http.Server{Addr: addr, Handler: s.router}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("start http server error", zap.Error(err))
			return
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 平滑关闭,等待5秒钟处理
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}

	return nil
}
