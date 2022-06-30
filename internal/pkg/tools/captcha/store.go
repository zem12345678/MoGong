package captcha

import (
	"context"
	"mogong/global"
	"mogong/internal/pkg/common/cache/redis"
	"mogong/internal/pkg/log"

	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type RedisStore struct {
	Expiration int
	PreKey     string
	Context    context.Context
}

func NewDefaultRedisStore() *RedisStore {
	return &RedisStore{
		Expiration: global.CaptchaExpireTime,
		PreKey:     "CAPTHCA_",
	}
}

func (rs *RedisStore) UseWithCtx(ctx context.Context) base64Captcha.Store {
	rs.Context = ctx
	return rs
}

func (rs *RedisStore) Set(id string, value string) error {
	err := redis.Set(rs.PreKey+id, value, rs.Expiration)
	if err != nil {
		log.Client.Logger.Error("RedisStoreSetError", zap.Error(err))
		return err
	}
	return err
}

func (rs *RedisStore) Get(key string, clear bool) string {
	val, err := redis.Get(key)
	if err != nil {
		log.Client.Logger.Error("RedisStoreGetError", zap.Error(err))
		return ""
	}
	if clear {
		res, err := redis.Delete(key)
		if res != true && err != nil {
			log.Client.Logger.Error("RedisStoreGetError", zap.Error(err))
			return ""
		}
	}
	return string(val)
}

func (rs *RedisStore) Verify(id, answer string, clear bool) bool {
	key := rs.PreKey + id
	v := rs.Get(key, clear)
	return v == answer
}
