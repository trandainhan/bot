package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
)

func ask_worker(id string, coin string, askB float64, perProfitStep float64) {
	marketParam := coin + "USDT"
	_, askPriceByQuantity := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)
	for {
		runable := redisClient.Get("runable").(bool)
		riki1_runable := redisClient.Get("riki1_runable").(bool)
		perFeeBinance := redisClient.Get("per_fee_binance").(float64) // 0.075 / 100
		perProfitAsk := redisClient.Get("per_profit_ask").(float64)
		if runable && riki1_runable {
			askF, isOutRange := calculateAskFFromAskB(askPriceByQuantity, perFeeBinance, perProfitAsk, minPrice, maxPrice)
			if isOutRange {
				text := fmt.Sprintf("%s @ndtan Error! Price out of range. PriceF: %v PriceB: %v Range: %v - %v", coin, askF, askPriceByQuantity, minPrice, maxPrice)
				chatID, _ := strconv.ParseInt(os.Getenv("chat_id"), 10, 64)
				go teleClient.SendMessage(text, chatID)
				time.Sleep(2 * time.Second)
			} else {
				id := "riki1"
				perProfitStep := 1.0
				worker(id, coin, askF, askPriceByQuantity, perProfitStep)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
