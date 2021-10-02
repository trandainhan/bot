package main

import (
	"fmt"
	"math"
	"time"

	"gitlab.com/fiahub/bot/internal/utils"
)

func adjustUpTrendPercentage() {
	now := time.Now()

	oneLastTime := now.Add(time.Duration(-1) * time.Hour)
	key := fmt.Sprintf("%s_price_%d_%d_%d", coin, oneLastTime.Day(), oneLastTime.Hour(), oneLastTime.Minute())
	oneHourAgoPrice, err := redisClient.GetFloat64(key)
	if err != nil {
		return
	}

	threeLastTime := now.Add(time.Duration(-3) * time.Hour)
	key = fmt.Sprintf("%s_price_%d_%d_%d", coin, threeLastTime.Day(), threeLastTime.Hour(), threeLastTime.Minute())
	threeHourAgoPrice, err := redisClient.GetFloat64(key)
	if err != nil {
		return
	}

	sixLastTime := now.Add(time.Duration(-6) * time.Hour)
	key = fmt.Sprintf("%s_price_%d_%d_%d", coin, sixLastTime.Day(), sixLastTime.Hour(), sixLastTime.Minute())
	sixHourAgoPrice, err := redisClient.GetFloat64(key)

	if err != nil {
		return
	}

	percentage1 := (currentBidPrice - oneHourAgoPrice) * 100 / oneHourAgoPrice
	percentage2 := (currentBidPrice - threeHourAgoPrice) * 100 / threeHourAgoPrice
	percentage3 := (currentBidPrice - sixHourAgoPrice) * 100 / sixHourAgoPrice

	upTrendfactor, _ := redisClient.GetFloat64(coin + "_up_trend_percentage_factor")
	finalPercentage := utils.RoundTo((percentage1+percentage2+percentage3)*upTrendfactor/3, decimalsToRound)

	maximumUpTrendPercentage := 50.0
	if math.Abs(finalPercentage) > maximumUpTrendPercentage {
		finalPercentage = maximumUpTrendPercentage
	}

	upTrendKey := coin + "_up_trend_percentage"
	redisClient.Set(upTrendKey, fmt.Sprintf("%.2f", finalPercentage), 0)
	text := fmt.Sprintf("Update %s to %.2f", upTrendKey, finalPercentage)
	go teleClient.SendMessage(text, chatID)
}
