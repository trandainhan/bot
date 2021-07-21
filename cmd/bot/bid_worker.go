package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
)

func bid_worker(id string, coin string, perProfitStep float64, cancalFactor int, results chan<- bool) {
	for {
		randNumber := rand.Intn(1000)
		time.Sleep(time.Duration(randNumber) * time.Millisecond)
		runableKey := fmt.Sprintf("%s_bid_runable", coin)
		runable := redisClient.GetBool(runableKey)
		perFeeBinance := redisClient.GetFloat64("per_fee_" + currentExchange)
		perProfitBid := redisClient.GetFloat64(coin + "_" + currentExchange + "_per_profit_bid")
		exchangeBidPrice, err := exchanges.GetBidPriceByQuantity(coin, quantityToGetPrice)
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
		fiahubBidPrice, isOutRange := calculateBidFFromBidB(exchangeBidPrice, perFeeBinance, perProfitBid, minPrice, maxPrice)
		if isOutRange {
			text := fmt.Sprintf("%s %s Err Price out of range. PriceF: %v PriceBidB: %v Range: %v - %v",
				coin, os.Getenv("TELEGRAM_HANDLER"), fiahubBidPrice, exchangeBidPrice, minPrice, maxPrice)
			log.Println(text)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("Trade bid order with coin: %s fiahubBidPrice: %.6f exchangeBidPrice: %.6f", coin, fiahubBidPrice, exchangeBidPrice)
			trade_bid(id, coin, fiahubBidPrice, exchangeBidPrice, cancalFactor)
		}

		time.Sleep(3 * time.Second)
	}
	results <- true
}
