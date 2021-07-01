package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
)

func bid_worker(id string, coin string, perProfitStep float64, cancalFactor int, results chan<- bool) {
	marketParam := coin + "USDT"
	for {
		randNumber := rand.Intn(1000)
		time.Sleep(time.Duration(randNumber) * time.Millisecond)
		runableKey := fmt.Sprintf("%s_bid_runable", coin)
		runable := redisClient.GetBool(runableKey)
		perFeeBinance := redisClient.GetFloat64("per_fee_binance")
		perProfitBid := redisClient.GetFloat64(coin + "_per_profit_bid")
		bidB, _, err := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)
		if err != nil {
			text := fmt.Sprintf("%s Err GetPriceByQuantity: %s", coin, err.Error())
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(30 * time.Second)
			continue
		}
		if !runable {
			time.Sleep(30 * time.Second)
			continue
		}
		perProfitBid = perProfitBid + perProfitStep*0.6/100
		bidF, isOutRange := calculateBidFFromBidB(bidB, perFeeBinance, perProfitBid, minPrice, maxPrice)
		if isOutRange {
			text := fmt.Sprintf("%s %s Error! Price out of range. PriceF: %v PriceBidB: %v Range: %v - %v",
				coin, os.Getenv("TELEGRAM_HANDLER"), bidF, bidB, minPrice, maxPrice)
			log.Println(text)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("Trade bid order with coin: %s bidf: %v bidB: %v", coin, bidF, bidB)
			trade_bid(id, coin, bidF, bidB, cancalFactor)
		}

		time.Sleep(3 * time.Second)
	}
	results <- true
}
