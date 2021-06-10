package fiahub

import "gitlab.com/fiahub/bot/internal/rediswrapper"

type Fiahub struct {
	RedisClient *rediswrapper.MyRedis
	Token       string
}

func (fia *Fiahub) SetToken(token string) bool {
	fia.Token = token
	return true
}
