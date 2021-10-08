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
		coinRunable := redisClient.GetBool(currentExchange + coin + "_worker_runable")
		workerRunable := redisClient.GetBool(coin + "_sell_worker_runable")
		if !autoMode || !coinRunable || !workerRunable {
			time.Sleep(30 * time.Second)
			continue
		}

		totalBuySize, _ := redisClient.GetFloat64(coin + "_total_buy_size")
		totalSellSize, _ := redisClient.GetFloat64(coin + "_total_sell_size")
		if totalSellSize-totalBuySize > buySellDiffSize {
			text := fmt.Sprintf("%s Ignore sell, due to sell too much, diff: %.3f", coin, totalSellSize-totalBuySize)
			log.Println(text)
			go teleClient.SendMessage(text, chatID)
			time.Sleep(1 * time.Minute)
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
			orderQuantity = utils.RoundTo(maxOrderQuantity, orderQuantityToRound)
		}
		order, err := placeOrder(id, orderQuantity, finalPrice, "sell")
		if err != nil {
			log.Printf("%s Err Can not place sell order %s", coin, err.Error())
			time.Sleep(30 * time.Second)
			continue
		}
		time.Sleep(10 * time.Second)

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
				go calculateProfit(orderDetails.ID, orderDetails.ExecutedQty, orderDetails.Price, "sell")
				time.Sleep(30 * time.Second)
				break
			} else if orderDetails.IsCanceled() {
				log.Printf("%s %s Order %d is canceled at price %f", coin, id, orderDetails.ID, orderDetails.Price)
				break
			}

			if currentAskPrice > exchangeAskPrice+jumpPrice+upTrendPriceAdjust || currentAskPrice < exchangeAskPrice-jumpPrice+upTrendPriceAdjust {
				if orderDetails.IsPartiallyFilled() {
					go calculateProfit(orderDetails.ID, orderDetails.ExecutedQty, orderDetails.Price, "sell")
				}
				_, err := exchangeClient.CancelOrder(coin, orderDetails.ID, orderDetails.ClientID)
				if err != nil {
					text := fmt.Sprintf("%s %s Err CancelOrder: %s", coin, id, err)
					go teleClient.SendMessage(text, chatErrorID)
				} else {
					text := fmt.Sprintf("%s %s CancelOrder %d due to price change: currentPrice: %.3f, lastPrice: %.3f", coin, id, orderDetails.ID, currentAskPrice, exchangeAskPrice)
					log.Println(text)
					go teleClient.SendMessage(text, chatID)
				}
				break
			}
			time.Sleep(30 * time.Second)
		}
	}
	results <- true
}
