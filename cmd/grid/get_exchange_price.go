package main

import (
	"fmt"
	"log"

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
	text := fmt.Sprintf("%s Updated currentAskPrice to: %f", coin, currentAskPrice)
	log.Println(text)
}

func updateCurrentBidPrice() {
	exchangeBidPrice, err := exchanges.GetBidPriceByQuantity(coin, quantityToGetPrice)
	if err != nil {
		text := fmt.Sprintf("%s Err GetPriceByQuantity: %s", coin, err.Error())
		go teleClient.SendMessage(text, chatErrorID)
		return
	}
	text := fmt.Sprintf("%s Updated currentBidPrice to: %f", coin, currentBidPrice)
	log.Println(text)
	currentBidPrice = exchangeBidPrice
}
