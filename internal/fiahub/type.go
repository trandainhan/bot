package fiahub

import "gitlab.com/fiahub/bot/internal/rediswrapper"

type Fiahub struct {
	RedisClient *rediswrapper.MyRedis
}
