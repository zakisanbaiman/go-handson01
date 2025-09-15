package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/zakisanbaiman/go-handson01/config"
	"github.com/zakisanbaiman/go-handson01/entity"
)

func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &KVS{
		Cli: cli,
	}, nil
}

type KVS struct {
	Cli *redis.Client
}

func (kvs *KVS) Save(ctx context.Context, key string, userID entity.UserID) error {
	return kvs.Cli.Set(ctx, key, userID, 0).Err()
}

func (kvs *KVS) Load(ctx context.Context, key string) (entity.UserID, error) {
	userID, err := kvs.Cli.Get(ctx, key).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return entity.UserID(userID), nil
}
