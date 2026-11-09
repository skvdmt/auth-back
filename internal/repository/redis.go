package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/skvdmt/auth-back/internal/model"
)

const (
	REFRESH_TOKEN_USER_LIMIT = 5
	REFRESH_TOKEN_TTL        = 1 * time.Hour
	ACCESS_TOKEN_TTL         = time.Minute * 15
	REDIS_PASSWORD_ENV       = "REDIS_PASSWORD"
)

// Репозитория Redis.
type Redis struct {
	client *redis.Client
}

// NewRedis Констрктор.
func NewRedis() (*Redis, error) {
	model.Log.Info.Info("repository redis creating")
	p, o := os.LookupEnv(REDIS_PASSWORD_ENV)
	if !o {
		return nil, fmt.Errorf("env %s not set", REDIS_PASSWORD_ENV)
	}
	r := &Redis{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf(
				"%s:%d", model.Conf.Redis.Host, model.Conf.Redis.Port),
			Password: p,
			DB:       0,
		}),
	}
	r.client.Ping(context.Background()).Err()
	model.Log.Info.Info("repository redis connection success")
	return r, nil
}

// RefreshTokenSaving Сохранение токена обновления.
func (r *Redis) RefreshTokenSaving(ctx context.Context, username, token string) error {
	p := r.client.TxPipeline()
	p.ZAdd(ctx, username, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: token,
	})
	p.ZRemRangeByRank(ctx, username, 0, -1-REFRESH_TOKEN_USER_LIMIT)
	p.Expire(ctx, username, REFRESH_TOKEN_TTL)
	if _, err := p.Exec(ctx); err != nil {
		return err
	}
	return nil
}

// RefreshTokenValidation Валидация токена обновления.
func (r *Redis) RefreshTokenValidation(ctx context.Context, username, token string) (bool, error) {
	c, err := r.client.ZRem(ctx, username, token).Result()
	if err != nil {
		return false, fmt.Errorf("redis zrem failed: %w", err)
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

// Stop Отстановка.
func (r *Redis) Stop() error {
	if err := r.client.Close(); err != nil {
		return err
	}
	model.Log.Info.Info("repository redis stopped")
	return nil
}
