package cronjob

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type DefaultCronJobService struct {
	logger *zap.Logger
	v      *viper.Viper
}

func NewDefualteCronJobService(logger *zap.Logger, v *viper.Viper) *DefaultCronJobService {
	return &DefaultCronJobService{
		logger: logger.With(zap.String("type", "DefaultCronJobService")),
		v:      v,
	}
}

func (s *DefaultCronJobService) Hello() (string, error) {
	fmt.Println("Hello")
	return "", nil
}
