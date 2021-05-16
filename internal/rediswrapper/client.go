package rediswrapper

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type MyRedis struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(ctx context.Context, redisURL string) *MyRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &MyRedis{client, ctx}
}

func (myRedis *MyRedis) Set(key, value interface{}) bool {
	err := myRedis.Client.Set(myRedis.Ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	return true
}

func (myRedis *MyRedis) Get(key string) interface{} {
	val, err := myRedis.Client.Get(myRedis.Ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		panic(err)
	}
	return val
}
