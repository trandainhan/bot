package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

func checkPriceVolatility() {
	now := time.Now()

	last5 := now.Add(time.Duration(-5) * time.Minute)
	key := fmt.Sprintf("%s_price_%d_%d_%d", coin, last5.Day(), last5.Hour(), last5.Minute())
	last5Price, err := redisClient.GetFloat64(key)
	if err != nil {
		return
	}

	last10 := now.Add(time.Duration(-10) * time.Minute)
	key = fmt.Sprintf("%s_price_%d_%d_%d", coin, last10.Day(), last10.Hour(), last10.Minute())
	last10Price, err := redisClient.GetFloat64(key)
	if err != nil {
		return
	}

	last15 := now.Add(time.Duration(-15) * time.Minute)
	key = fmt.Sprintf("%s_price_%d_%d_%d", coin, last15.Day(), last15.Hour(), last15.Minute())
	last15Price, err := redisClient.GetFloat64(key)

	if err != nil {
		return
	}

	percentage1 := math.Abs((currentBidPrice - last5Price) * 100 / last5Price)
	percentage2 := math.Abs((currentBidPrice - last10Price) * 100 / last10Price)
	percentage3 := math.Abs((currentBidPrice - last15Price) * 100 / last15Price)

	redisClient.Set(coin+"_price_percentage_change_in_5min", percentage1, 0)
	redisClient.Set(coin+"_price_percentage_change_in_10min", percentage2, 0)
	redisClient.Set(coin+"_price_percentage_change_in_15min", percentage3, 0)

	if percentage1 > 3.5 || percentage2 > 4.5 || percentage3 > 5.5 {
		redisClient.Set(currentExchange+coin+"_worker_runable", false, 0)
		teleHanlder := os.Getenv("TELEGRAM_HANDLER")
		text := fmt.Sprintf("%s %s stop buy and sell worker due to high price volatility\n Price changed in: 5min: %.2f, 10min: %.2f, 15min: %.2f",
			teleHanlder, coin, percentage1, percentage2, percentage3)
		go teleClient.SendMessage(text, chatRunableID)
		return
	}

	coinRunable := redisClient.GetBool(currentExchange + coin + "_worker_runable")
	if coinRunable == false {
		redisClient.Set(currentExchange+coin+"_worker_runable", true, 0)
		text := fmt.Sprintf("Reset %s_%s_worker_runable to %v", currentExchange, coin, true)
		teleClient.SendMessage(text, chatRunableID)
	}
}
