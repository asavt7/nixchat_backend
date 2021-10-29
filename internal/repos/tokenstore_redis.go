package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/asavt7/nixchat_backend/internal/config"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitRedisOpts(redisConfig *config.RedisConfig) *redis.Options {
	op := &redis.Options{
		Addr: redisConfig.Host + ":" + redisConfig.Port,
	}
	log.Debugf("Redis config %v", op)
	return op
}

func InitRedisClient(opt *redis.Options) *redis.Client {
	client := redis.NewClient(opt)

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to run redis client : %s", err)
	}

	return client
}

// RedisTokenStorage struct
type RedisTokenStorage struct {
	client            *redis.Client
	autoLogoffMinutes time.Duration
}

func NewRedisTokenStorage(client *redis.Client, cfg *config.AuthConfig) *RedisTokenStorage {
	return &RedisTokenStorage{client: client, autoLogoffMinutes: cfg.AutoLogoffTime}
}

func (r *RedisTokenStorage) Get(userID string) (model.CachedTokens, error) {
	key := fmt.Sprintf("token-%s", userID)
	res, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return model.CachedTokens{}, err
	}

	err = r.client.Expire(context.Background(), key, r.autoLogoffMinutes).Err()
	if err != nil {
		return model.CachedTokens{}, err
	}

	cachedTokens := new(model.CachedTokens)
	err = json.Unmarshal([]byte(res), cachedTokens)
	return *cachedTokens, err

}

func (r *RedisTokenStorage) Delete(userID string) error {
	key := fmt.Sprintf("token-%s", userID)
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisTokenStorage) Save(userID string, tokens model.CachedTokens) (model.CachedTokens, error) {
	key := fmt.Sprintf("token-%s", userID)
	cacheJSON, err := json.Marshal(tokens)
	if err != nil {
		return tokens, err
	}
	err = r.client.Set(context.Background(), key, cacheJSON, r.autoLogoffMinutes).Err()
	return tokens, err
}
