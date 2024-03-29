package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/utils"
)

func buy_worker(id string, coin string, step int, results chan<- bool) {
	for {
		autoMode := redisClient.GetBool(currentExchange + "_auto_mode")
		coinRunable := redisClient.GetBool(currentExchange + coin + "_worker_runable")
		workerRunable := redisClient.GetBool(coin + "_buy_worker_runable")
		if !autoMode || !coinRunable || !workerRunable {
			time.Sleep(30 * time.Second)
			continue
		}

		totalBuySize, _ := redisClient.GetFloat64(coin + "_total_buy_size")
		totalSellSize, _ := redisClient.GetFloat64(coin + "_total_sell_size")
		diff := totalBuySize - totalSellSize
		if diff > buySellDiffSize {
			text := fmt.Sprintf("%s Ignore buy, due to buy too much, diff: %.3f", coin, totalBuySize-totalSellSize)
			log.Println(text)
			go teleClient.SendMessage(text, chatID)
			time.Sleep(1 * time.Minute)
			continue
		}

		exchangeBidPrice, err := exchanges.GetBidPriceByQuantity(coin, fiat, quantityToGetPrice)

		jumpPrice := exchangeBidPrice * jumpPricePercentage / 100

		key := fmt.Sprintf("%s_up_trend_percentage", coin)
		upTrendPercentage, _ := redisClient.GetFloat64(key)
		upTrendPriceAdjust := jumpPrice * upTrendPercentage / 100

		// When market is up trend, upTrendPercentage > 0 => upTrendPriceAdjust > 0, Buy order price should be closed to the current market price
		// When market is down trend, upTrendPercentage < 0 => upTrendPriceAdjust < 0, Buy order price should be distance from the current market price
		finalPrice := utils.RoundTo(exchangeBidPrice-jumpPrice*float64(step)+upTrendPriceAdjust, decimalsToRound)

		// if totalSellSize-totalBuySize > buySellDiffSize { // over sell, place a buy order with price less then last_sell_price
		// 	lastSellPrice, _ := redisClient.GetFloat64(coin + "_last_sell_price")
		// 	jumpPrice = lastSellPrice * jumpPricePercentage / 100
		// 	finalPrice = utils.RoundTo(lastSellPrice-jumpPrice/2, decimalsToRound)
		// }

		// Contraint order maximum quanity
		maxOrderQuantity := maximumOrderUsdt / currentBidPrice
		if orderQuantity > maxOrderQuantity {
			orderQuantity = utils.RoundTo(maxOrderQuantity, orderQuantityToRound)
		}
		order, err := placeOrder(id, orderQuantity, finalPrice, "buy")
		if err != nil {
			log.Printf("%s Err Can not place buy order %s", coin, err.Error())
			time.Sleep(30 * time.Second)
			continue
		}
		time.Sleep(10 * time.Second)

		for {
			orderDetails, err := exchangeClient.GetOrder(coin, fiat, order.ID, order.ClientID)
			if err != nil {
				text := fmt.Sprintf("%s %s Err getOrderDetails: %s", coin, id, err)
				log.Println(text)
				go teleClient.SendMessage(text, chatErrorID)
				time.Sleep(60 * time.Second)
				continue
			}

			log.Printf("%s %s Check Order %d status: %s", coin, id, orderDetails.ID, orderDetails.Status)
			if orderDetails.IsFilled() {
				calculateProfit(orderDetails.ID, orderDetails.ExecutedQty, orderDetails.Price, "buy")
				time.Sleep(30 * time.Second)
				break
			} else if orderDetails.IsCanceled() {
				log.Printf("%s %s Order %d is canceled at price %.3f", coin, id, orderDetails.ID, orderDetails.Price)
				break
			}

			if currentBidPrice > exchangeBidPrice+jumpPrice+upTrendPriceAdjust || currentBidPrice < exchangeBidPrice-jumpPrice+upTrendPriceAdjust {
				if orderDetails.IsPartiallyFilled() {
					calculateProfit(orderDetails.ID, orderDetails.ExecutedQty, orderDetails.Price, "buy")
				}
				_, err := exchangeClient.CancelOrder(coin, fiat, orderDetails.ID, orderDetails.ClientID)
				if err != nil {
					text := fmt.Sprintf("%s %s Err CancelOrder: %s", coin, id, err)
					log.Println(text)
					go teleClient.SendMessage(text, chatErrorID)
				} else {
					text := fmt.Sprintf("%s %s CancelOrder %d due to price change: currentPrice: %.3f, lastPrice: %.3f", coin, id, orderDetails.ID, currentBidPrice, exchangeBidPrice)
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
