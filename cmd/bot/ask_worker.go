package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
)

func ask_worker(id string, coin string, perProfitStep float64, results chan<- bool) {
	marketParam := coin + "USDT"
	for {
		runableKey := fmt.Sprintf("%s_%s_runable", coin, id)
		runable := redisClient.GetBool(runableKey)
		perFeeBinance := redisClient.GetFloat64("per_fee_binance")
		perProfitAsk := redisClient.GetFloat64("per_profit_ask")
		_, askB := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)
		if !(runable) {
			time.Sleep(30 * time.Second)
			continue
		}

		perProfitAsk = perProfitAsk + perProfitStep*0.6/100
		askF, isOutRange := calculateAskFFromAskB(askB, perFeeBinance, perProfitAsk, minPrice, maxPrice)
		if isOutRange {
			text := fmt.Sprintf("%s %s Error! Price out of range. PriceF: %v PriceAskB: %v Range: %v - %v",
				coin, os.Getenv("TELEGRAM_HANDLER"), askF, askB, minPrice, maxPrice)
			log.Println(text)
			go teleClient.SendMessage(text, chatID)
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("Trade ask order with coin: %s askf: %v askB: %v", coin, askF, askB)
			trade_ask(id, coin, askF, askB)
		}

		time.Sleep(3 * time.Second)
	}
	results <- true
}
