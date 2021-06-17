package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
)

func ask_worker(id string, coin string, perProfitStep float64, results chan<- bool) {
	marketParam := coin + "USDT"
	for {
		randNumber := rand.Intn(1000)
		time.Sleep(time.Duration(randNumber) * time.Millisecond)
		runableKey := fmt.Sprintf("%s_%s_runable", coin, id)
		runable := redisClient.GetBool(runableKey)
		perFeeBinance := redisClient.GetFloat64("per_fee_binance")
		perProfitAsk := redisClient.GetFloat64("per_profit_ask")
		_, askB := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)
		if askB == -1.0 {
			text := "There is may be a error when get price from binance, skip and wait"
			go teleClient.SendMessage(text, chatID)
			time.Sleep(30 * time.Second)
			continue
		}
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
