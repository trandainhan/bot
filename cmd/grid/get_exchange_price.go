package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
)

func updateCurrentAskPrice() {
	exchangeAskPrice, err := exchanges.GetAskPriceByQuantity(coin, quantityToGetPrice)
	if err != nil {
		text := fmt.Sprintf("%s Err GetPriceByQuantity: %s", coin, err.Error())
		go teleClient.SendMessage(text, chatErrorID)
		return
	}
	currentAskPrice = exchangeAskPrice
	log.Printf("%s Updated currentAskPrice to: %f", coin, currentAskPrice)
}

func updateCurrentBidPrice() {
	exchangeBidPrice, err := exchanges.GetBidPriceByQuantity(coin, quantityToGetPrice)
	if err != nil {
		text := fmt.Sprintf("%s Err GetPriceByQuantity: %s", coin, err.Error())
		go teleClient.SendMessage(text, chatErrorID)
		return
	}
	currentBidPrice = exchangeBidPrice
	now := time.Now()
	key := fmt.Sprintf("%s_price_%d_%d_%d", coin, now.Day(), now.Hour(), now.Minute())
	redisClient.Set(key, currentBidPrice, time.Duration(24)*time.Hour)

	log.Printf("%s Updated currentBidPrice to: %f", coin, currentBidPrice)
}
