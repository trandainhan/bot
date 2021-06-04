package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func ask_worker(id string, coin string, askB float64, perProfitStep float64, results chan<- bool) {
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	for {
		runable := redisClient.GetBool("runable")
		perFeeBinance := redisClient.GetFloat64("per_fee_binance")
		perProfitAsk := redisClient.GetFloat64("per_profit_ask")
		if !(runable) {
			time.Sleep(30 * time.Second)
			continue
		}

		askF, isOutRange := calculateAskFFromAskB(askB, perFeeBinance, perProfitAsk, minPrice, maxPrice)
		if isOutRange {
			text := fmt.Sprintf("%s @ndtan Error! Price out of range. PriceF: %v PriceB: %v Range: %v - %v", coin, askF, askB, minPrice, maxPrice)
			go teleClient.SendMessage(text, chatID)
			time.Sleep(2 * time.Second)
		} else {
			trade_ask(id, coin, askF, askB, perProfitStep)
		}

		time.Sleep(3 * time.Second)
	}
	results <- true
}
