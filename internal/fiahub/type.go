package fiahub

import "gitlab.com/fiahub/bot/internal/rediswrapper"

type Fiahub struct {
	RedisClient         *rediswrapper.MyRedis
	Token               string
	latestCancelAllTime int64
}

func (fia *Fiahub) SetToken(token string) bool {
	fia.Token = token
	return true
}

func (fia *Fiahub) SetCancelTime(time int64) bool {
	fia.latestCancelAllTime = time
	return true
}

func (fia Fiahub) GetCancelTime() int64 {
	return fia.latestCancelAllTime
}
