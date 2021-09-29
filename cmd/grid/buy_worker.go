package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/exchanges/binance"
	"gitlab.com/fiahub/bot/internal/exchanges/ftx"
	"gitlab.com/fiahub/bot/internal/utils"
)

func buy_worker(id string, coin string, step int, results chan<- bool) {
	for {
		autoMode := redisClient.GetBool(currentExchange + "_auto_mode")
		if !autoMode {
			time.Sleep(30 * time.Second)
			continue
		}

		exchangeBidPrice, err := exchanges.GetBidPriceByQuantity(coin, quantityToGetPrice)
		log.Println(exchangeBidPrice)

		jumpPrice := exchangeBidPrice * jumpPercentage / 100

		finalPrice := utils.RoundTo(exchangeBidPrice-jumpPrice*float64(step), decimalsToRound)
		order, err := placeOrder(id, defaultQuantity, finalPrice, "buy")
		if err != nil {
			continue
		}
		time.Sleep(5 * time.Second)

		for {
			orderDetails, err := exchangeClient.GetOrder(coin, order.ID, order.ClientID)
			if err != nil {
				text := fmt.Sprintf("%s %s Err getOrderDetails: %s", coin, id, err)
				log.Println(text)
				go teleClient.SendMessage(text, chatErrorID)
			} else {
				text := fmt.Sprintf("%s %s Check Order %d status: %s", coin, id, orderDetails.ID, orderDetails.Status)
				log.Println(text)
				if isFilledStatus(orderDetails.Status) {
					text := fmt.Sprintf("%s %s Order %d is filled at price %f", coin, id, orderDetails.ID, orderDetails.Price)
					go teleClient.SendMessage(text, chatID)
					log.Println(text)
					break
				}
			}
			if currentBidPrice > exchangeBidPrice+jumpPrice || currentBidPrice < exchangeBidPrice-jumpPrice {
				_, err := exchangeClient.CancelOrder(coin, orderDetails.ID, orderDetails.ClientID)
				if err != nil {
					text := fmt.Sprintf("%s %s Err CancelOrder: %s", coin, id, err)
					log.Println(text)
					go teleClient.SendMessage(text, chatErrorID)
				} else {
					text := fmt.Sprintf("%s %s CancelOrder %d due to price change: currentPrice: %f, lastPrice: %f", coin, id, orderDetails.ID, currentBidPrice, exchangeBidPrice)
					log.Println(text)
					go teleClient.SendMessage(text, chatID)
				}
				break
			}
			time.Sleep(60 * time.Second)
		}
	}
	results <- true
}

func isFilledStatus(status string) bool {
	if currentExchange == "FTX" {
		return status == ftx.ORDER_CLOSED
	}
	return status == binance.ORDER_FILLED
}
