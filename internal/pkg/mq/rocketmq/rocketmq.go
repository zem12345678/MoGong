package rocketmq

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Options struct{}

type RocketMq struct{}

func NewOptions(v viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)

	if err = v.UnmarshalKey("rocketmq", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal rocketmq configuration error")
	}
	logger.Info("load rocketmq configuration success")

	return o, nil

}

func New(o *Options) (*RocketMq, error) {
	return &RocketMq{}, nil
}
