package redis

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"rplss/service/config"
	"time"
)

type marshaller struct {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	v interface{}
}

func (b *marshaller) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, b.v)
}

func (b *marshaller) MarshalBinary() (data []byte, err error) {
	return json.Marshal(b.v)
}

type Redis interface {
	Get(ctx context.Context, key string, value interface{}) (error, bool)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	LRange(ctx context.Context, key string, f GetListCallBack) error
	LPush(ctx context.Context, key string, value interface{}) error
	LTrim(ctx context.Context, key string, n int) error
	Close() error
}

type RedisDb struct {
	client *redis.Client
}

func NewRedisConnection(config config.RedisConfigs) (Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Host, config.Port),
		DB:   config.Db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &RedisDb{
		client: client,
	}, nil
}

func (r *RedisDb) Get(ctx context.Context, key string, value interface{}) (error, bool) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	m := marshaller{v: value}
	err := r.client.Get(ctx, key).Scan(&m)
	if err != nil {
		if err == redis.Nil {
			return nil, false
		}
		return err, false
	}

	return nil, true
}

func (r *RedisDb) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	m := marshaller{v: value}
	err := r.client.Set(ctx, key, &m, ttl).Err()
	return err
}

func (r *RedisDb) Del(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := r.client.Del(ctx, key).Err()
	return err
}

type GetListCallBack func(item string) error

func (r *RedisDb) LRange(ctx context.Context, key string, f GetListCallBack) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	m, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	for _, v := range m {
		err = f(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisDb) LPush(ctx context.Context, key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	m := marshaller{v: value}
	err := r.client.LPush(ctx, key, &m).Err()
	return err
}

func (r *RedisDb) LTrim(ctx context.Context, key string, n int) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := r.client.LTrim(ctx, key, 0, int64(n)).Err()
	return err
}

func (r *RedisDb) Close() error {
	return r.client.Close()
}
