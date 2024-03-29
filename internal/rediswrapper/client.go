package rediswrapper

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type MyRedis struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(ctx context.Context, redisURL string, dbNum int, pass string) *MyRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: pass,  // no password set
		DB:       dbNum, // use default DB
	})

	return &MyRedis{Client: client, Ctx: ctx}
}

func (myRedis *MyRedis) Set(key string, value interface{}, expiration time.Duration) bool {
	err := myRedis.Client.Set(myRedis.Ctx, key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
	return true
}

func (myRedis *MyRedis) Get(key string) (string, error) {
	val, err := myRedis.Client.Get(myRedis.Ctx, key).Result()
	if err == redis.Nil {
		return "", err
	}
	return val, nil
}

func (myRedis *MyRedis) GetBool(key string) bool {
	val, err := myRedis.Client.Get(myRedis.Ctx, key).Bool()
	if err != nil {
		panic(err)
	}
	return val
}

func (myRedis *MyRedis) GetFloat64(key string) (float64, error) {
	val, err := myRedis.Client.Get(myRedis.Ctx, key).Float64()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (myRedis *MyRedis) GetInt64(key string) int64 {
	val, err := myRedis.Client.Get(myRedis.Ctx, key).Int64()
	if err != nil {
		panic(err)
	}
	return val
}

func (myRedis *MyRedis) GetInt(key string) int {
	val, err := myRedis.Client.Get(myRedis.Ctx, key).Int()
	if err != nil {
		panic(err)
	}
	return val
}

func (myRedis *MyRedis) GetTime(key string) (*time.Time, error) {
	val, err := myRedis.Client.Get(myRedis.Ctx, key).Time()
	if err != nil {
		return nil, err
	}
	return &val, nil
}

func (myRedis *MyRedis) Del(key string) bool {
	err := myRedis.Client.Del(myRedis.Ctx, key).Err()
	if err != nil {
		panic(err)
	}
	return true
}
