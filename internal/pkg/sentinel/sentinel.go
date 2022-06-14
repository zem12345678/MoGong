package sentinel

import (
	"fmt"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/util"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Options struct {
	QpsLimit struct {
		StatIntervalMs uint32
		Threshold      float64
		Resource       string
	}
	Circuitbreaker struct {
		Resource                     string
		RetryTimeoutMs               uint32
		MinRequestAmount             uint64
		StatIntervalMs               uint32
		StatSlidingWindowBucketCount uint32
		MaxAllowedRtMs               uint64
		Threshold                    float64
	}
}

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

type Entry struct {
	QpsLimit       string
	Circuitbreaker string
	QLch           chan struct{}
	Cbch           chan struct{}
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)

	if err = v.UnmarshalKey("sentinel", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal sentinel configuration error")
	}
	logger.Info("load sentinel configuration success")

	return o, nil
}

func New(o *Options) (*Entry, error) {
	if err := sentinel.InitDefault(); err != nil {
		zap.S().Error()
	}

	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               o.QpsLimit.Resource,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              o.QpsLimit.Threshold,
			StatIntervalInMs:       o.QpsLimit.StatIntervalMs,
		},
	})
	if err != nil {
		zap.S().Errorf("Unexpected error: %+v", err)
		return nil, err
	}

	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, slowRtUpperBound=50ms, maxSlowRequestRatio=50%
		{
			Resource:                     o.Circuitbreaker.Resource,
			Strategy:                     circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:               o.Circuitbreaker.RetryTimeoutMs,
			MinRequestAmount:             o.Circuitbreaker.MinRequestAmount,
			StatIntervalMs:               o.Circuitbreaker.StatIntervalMs,
			StatSlidingWindowBucketCount: o.Circuitbreaker.StatSlidingWindowBucketCount,
			MaxAllowedRtMs:               o.Circuitbreaker.MaxAllowedRtMs,
			Threshold:                    o.Circuitbreaker.Threshold,
		},
	})
	if err != nil {
		zap.S().Errorf("Unexpected error: %+v", err)
		return nil, err
	}
	return &Entry{
		o.QpsLimit.Resource,
		o.Circuitbreaker.Resource,
		make(chan struct{}),
		make(chan struct{}),
	}, nil
}
