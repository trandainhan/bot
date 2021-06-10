package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
)

func bid_worker(id string, coin string, perProfitStep float64, results chan<- bool) {
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	marketParam := coin + "USDT"
	for {
		runable := redisClient.GetBool("runable")
		perFeeBinance := redisClient.GetFloat64("per_fee_binance")
		perProfitBid := redisClient.GetFloat64("per_profit_bid")
		bidB, _ := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)
		if !(runable) {
			time.Sleep(30 * time.Second)
			continue
		}
		perProfitBid = perProfitBid + perProfitStep*0.6/100
		bidF, isOutRange := calculateBidFFromBidB(bidB, perFeeBinance, perProfitBid, minPrice, maxPrice)
		if isOutRange {
			text := fmt.Sprintf("%s %s Error! Price out of range. PriceF: %v PriceBidB: %v Range: %v - %v",
				coin, os.Getenv("TELEGRAM_HANDLER"), bidF, bidB, minPrice, maxPrice)
			log.Println(text)
			go teleClient.SendMessage(text, chatID)
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("Trade bid order with coin: %s bidf: %v bidB: %v", coin, bidF, bidB)
			trade_bid(id, coin, bidF, bidB)
		}

		time.Sleep(3 * time.Second)
	}
	results <- true
}
