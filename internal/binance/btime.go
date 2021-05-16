package binance // binance time

import (
	"encoding/json"
	"time"

	"gitlab.com/fiahub/bot/internal/utils"
)

type TimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}

func GetOffsetTimeUnix() int64 {
	now := time.Now()
	sec := now.UnixNano()
	url := "https://api.binance.com/api/v1/time"
	body, _, _ := utils.HttpGet(url)
	var result TimeResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		panic(err)
	}

	offset := result.ServerTime - sec/int64(time.Millisecond)
	return offset
}
