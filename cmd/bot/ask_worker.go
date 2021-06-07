package main

import (
	"fmt"
	"log"
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
			text := fmt.Sprintf("%s %s Error! Price out of range. PriceF: %v PriceAskB: %v Range: %v - %v",
				coin, os.Getenv("TELEGRAM_HANDLER"), askF, askB, minPrice, maxPrice)
			log.Println(text)
			go teleClient.SendMessage(text, chatID)
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("Trade ask order with coin: %s bidf: %v bidB: %v", coin, askF, askB)
			trade_ask(id, coin, askF, askB, perProfitStep)
		}

		time.Sleep(3 * time.Second)
	}
	results <- true
}
