package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
)

func ask_worker(id string, coin string, perProfitStep float64, cancelFactor int, results chan<- bool) {
	for {
		randNumber := rand.Intn(1000)
		time.Sleep(time.Duration(randNumber) * time.Millisecond)
		runableKey := fmt.Sprintf("%s_ask_runable", coin)
		runable := redisClient.GetBool(runableKey)
		perFeeBinance := redisClient.GetFloat64("per_fee_binance")
		perProfitAsk := redisClient.GetFloat64(coin + "_per_profit_ask")
		askB, err := exchanges.GetAskPriceByQuantity(coin, quantityToGetPrice)
		if err != nil {
			text := fmt.Sprintf("%s Err GetPriceByQuantity: %s", coin, err.Error())
			go teleClient.SendMessage(text, chatErrorID)
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
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("Trade ask order with coin: %s askf: %v askB: %v", coin, askF, askB)
			trade_ask(id, coin, askF, askB, cancelFactor)
		}

		time.Sleep(3 * time.Second)
	}
	results <- true
}
