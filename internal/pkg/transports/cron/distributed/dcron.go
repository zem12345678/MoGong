package distributed

import (
	"github.com/libi/dcron"
	"github.com/libi/dcron/driver/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Options struct {
	Server        string
	Projects      map[string]string
	RedisHost     string
	RedisPort     int
	RedisPassword string
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("cron", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal cron option error")
	}
	logger.Info("load cron options success", zap.Any("cron options", o))
	return o, err
}

type ServerOptional struct {
	spec string
	f    func()
}

type Server struct {
	app    string
	o      *Options
	logger *zap.Logger
	cron   *dcron.Dcron
	jobs   map[string]ServerOptional
}

type InitServers map[string]func()

func New(o *Options, logger *zap.Logger, init InitServers) (*Server, error) {
	optionals := make(map[string]ServerOptional)
	for name, spec := range o.Projects {
		if jobFunc, ok := init[name]; ok {
			optionals[name] = ServerOptional{
				spec: spec,
				f:    jobFunc,
			}
		} else {
			logger.Error("定时任务不存在", zap.String("name", name))
			return nil, errors.New("定时任务不存在")
		}
	}
	drv, _ := redis.NewDriver(&redis.Conf{
		Host:     o.RedisHost,
		Port:     o.RedisPort,
		Password: o.RedisPassword,
	})
	cron := dcron.NewDcron(o.Server, drv)
	return &Server{
		o:      o,
		logger: logger.With(zap.String("type", "cronServer")),
		cron:   cron,
		jobs:   optionals,
	}, nil
}

func (s *Server) Application(name string) {
	s.app = name
}

func (s *Server) Start() error {
	go func() {
		if err := s.register(); err != nil {
			s.logger.Fatal("failed to register cron: %v", zap.Error(err))
		}
		s.cron.Start()
	}()
	return nil
}

func (s *Server) register() error {
	for name, job := range s.jobs {
		err := s.cron.AddFunc(name, job.spec, job.f)
		if err != nil {
			s.logger.Error("注册job失败", zap.Error(err))
			return err
		}
		s.logger.Info("注册cron任务成功", zap.String("name", name))
	}
	return nil
}

func (s *Server) deRegister() error {
	s.cron.Stop()
	for name := range s.jobs {
		s.cron.Remove(name)
		s.logger.Info("deregister cron services success", zap.String("name", name))
	}
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("cron server stopping ...")
	if err := s.deRegister(); err != nil {
		return errors.Wrap(err, "deregister cron server error")
	}
	return nil
}
