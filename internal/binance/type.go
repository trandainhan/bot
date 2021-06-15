package binance

import "gitlab.com/fiahub/bot/internal/rediswrapper"

type Binance struct {
	RedisClient     *rediswrapper.MyRedis
	TimeDifferences int64
}
