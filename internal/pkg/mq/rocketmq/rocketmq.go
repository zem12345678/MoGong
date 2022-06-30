package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Options struct {
	Url []string
}

type RocketMqOptions struct {
	ConsumerOption consumer.Option
	ProducerOption producer.Option
}

var RocketmqOptions *RocketMqOptions

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("rocketmq", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal es option error")
	}

	logger.Info("load rocketmq options success", zap.Any("rocketmq options", o))
	return o, err
}

func New(o *Options) (*RocketMqOptions, error) {
	op := new(RocketMqOptions)
	if len(o.Url) == 0 {
		return nil, errors.New("缺少rocketmq配置")
	}
	ConsumerOption := consumer.WithNsResolver(primitive.NewPassthroughResolver(o.Url))
	ProducerOption := producer.WithNsResolver(primitive.NewPassthroughResolver(o.Url))

	op.ConsumerOption = ConsumerOption
	op.ProducerOption = ProducerOption
	RocketmqOptions = op
	return op, nil
}
