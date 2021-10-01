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
		runable := redisClient.GetBool(coin + "_sell_worker_runable")
		if !autoMode || !runable {
			time.Sleep(30 * time.Second)
			continue
		}

		exchangeAskPrice, err := exchanges.GetAskPriceByQuantity(coin, quantityToGetPrice)
		jumpPrice := exchangeAskPrice * jumpPricePercentage / 100

		key := fmt.Sprintf("%s_up_trend_percentage", coin)
		upTrendPercentage, _ := redisClient.GetFloat64(key)
		upTrendPriceAdjust := jumpPrice * upTrendPercentage / 100

		// When market is up trend, upTrendPercentage > 0 => upTrendPriceAdjust > 0, Sell order price should be distanced from the current market price
		// When market is down trend, upTrendPercentage < 0 => upTrendPriceAdjust < 0, Buy order price should be closed to the current market price
		finalPrice := utils.RoundTo(currentAskPrice+jumpPrice*float64(step)+upTrendPriceAdjust, decimalsToRound)

		maxOrderQuantity := maximumOrderUsdt / currentAskPrice
		if orderQuantity > maxOrderQuantity {
			orderQuantity = utils.RoundTo(maxOrderQuantity, 1)
		}
		order, err := placeOrder(id, orderQuantity, finalPrice, "sell")
		if err != nil {
			continue
		}

		time.Sleep(3 * time.Second)

		for {
			orderDetails, err := exchangeClient.GetOrder(coin, order.ID, order.ClientID)
			if err != nil {
				text := fmt.Sprintf("%s %s Err getOrderDetails: %s", coin, id, err)
				log.Println(text)
				go teleClient.SendMessage(text, chatErrorID)
				time.Sleep(60 * time.Second)
				continue
			}

			log.Printf("%s %s Check Order %d status: %s", coin, id, orderDetails.ID, orderDetails.Status)
			if orderDetails.IsFilled() {
				text := fmt.Sprintf("%s %s Order %d is filled at price %f", coin, id, orderDetails.ID, orderDetails.Price)
				go teleClient.SendMessage(text, chatProfitID)
				log.Println(text)
				go calculate_profit(orderDetails.ExecutedQty, orderDetails.Price, "sell")
				time.Sleep(10 * time.Second)
				break
			} else if orderDetails.IsCanceled() {
				log.Printf("%s %s Order %d is canceled at price %f", coin, id, orderDetails.ID, orderDetails.Price)
				break
			}

			if currentAskPrice > exchangeAskPrice+jumpPrice+upTrendPriceAdjust || currentAskPrice < exchangeAskPrice-jumpPrice+upTrendPriceAdjust {
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
