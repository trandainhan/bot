package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/utils"
)

func sell_worker(id string, coin string, step int, results chan<- bool) {
	for {
		autoMode := redisClient.GetBool(currentExchange + "_auto_mode")
		if !autoMode {
			time.Sleep(30 * time.Second)
			continue
		}

		exchangeAskPrice, err := exchanges.GetAskPriceByQuantity(coin, quantityToGetPrice)
		jumpPrice := exchangeAskPrice * jumpPercentage / 100

		finalPrice := utils.RoundTo(currentAskPrice+jumpPrice*float64(step), decimalsToRound)
		order, err := placeOrder(id, defaultQuantity, finalPrice, "sell")
		if err != nil {
			continue
		}

		time.Sleep(3 * time.Second)

		for {
			orderDetails, err := exchangeClient.GetOrder(coin, order.ID, order.ClientID)
			if err != nil {
				text := fmt.Sprintf("%s %s Err getOrderDetails: %s", coin, id, err)
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
			if currentAskPrice > exchangeAskPrice+jumpPrice || currentAskPrice < exchangeAskPrice-jumpPrice {
				_, err := exchangeClient.CancelOrder(coin, orderDetails.ID, orderDetails.ClientID)
				if err != nil {
					text := fmt.Sprintf("%s %s Err CancelOrder: %s", coin, id, err)
					go teleClient.SendMessage(text, chatErrorID)
				} else {
					text := fmt.Sprintf("%s %s CancelOrder %d due to price change: currentPrice: %f, lastPrice: %f", coin, id, orderDetails.ID, currentAskPrice, exchangeAskPrice)
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
