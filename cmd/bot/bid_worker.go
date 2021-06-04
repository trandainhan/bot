package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func bid_worker(id string, coin string, bidB float64, perProfitStep float64, results chan<- bool) {
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	for {
		runable := redisClient.GetBool("runable")
		perFeeBinance := redisClient.GetFloat64("per_fee_binance")
		perProfitAsk := redisClient.GetFloat64("per_profit_ask")
		if !(runable) {
			time.Sleep(30 * time.Second)
			continue
		}

		bidF, isOutRange := calculateBidFFromBidB(bidB, perFeeBinance, perProfitAsk, minPrice, maxPrice)
		if isOutRange {
			text := fmt.Sprintf("%s @ndtan Error! Price out of range. PriceF: %v PriceBidB: %v Range: %v - %v", coin, bidF, bidB, minPrice, maxPrice)
			log.Println(text)
			go teleClient.SendMessage(text, chatID)
			time.Sleep(2 * time.Second)
		} else {
			trade_bid(id, coin, bidF, bidB, perProfitStep)
		}

		time.Sleep(3 * time.Second)
	}
	results <- true
}
