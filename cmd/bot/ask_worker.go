package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func ask_worker(id string, coin string, askB float64, perProfitStep float64) {
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	for {
		runable := redisClient.Get("runable").(bool)
		riki1_runable := redisClient.Get("riki1_runable").(bool)
		perFeeBinance := redisClient.Get("per_fee_binance").(float64) // 0.075 / 100
		perProfitAsk := redisClient.Get("per_profit_ask").(float64)
		if !(runable && riki1_runable) {
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
}
