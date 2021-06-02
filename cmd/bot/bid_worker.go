package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func bid_worker(id string, coin string, bidB float64, perProfitStep float64) {
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

		bidF, isOutRange := calculateBidFFromBidB(bidB, perFeeBinance, perProfitAsk, minPrice, maxPrice)
		if isOutRange {
			text := fmt.Sprintf("%s @ndtan Error! Price out of range. PriceF: %v PriceB: %v Range: %v - %v", coin, bidF, bidB, minPrice, maxPrice)
			go teleClient.SendMessage(text, chatID)
			time.Sleep(2 * time.Second)
		} else {
			trade_bid(id, coin, bidF, bidB, perProfitStep)
		}

		time.Sleep(3 * time.Second)
	}
}
