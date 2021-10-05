package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/utils"
)

func makeAdditionalBuySell() {
	totalBuySize, _ := redisClient.GetFloat64(coin + "_total_buy_size")
	totalBuyValue, _ := redisClient.GetFloat64(coin + "_total_buy_value")
	averageBuyPrice := totalBuyValue / totalBuySize

	totalSellSize, _ := redisClient.GetFloat64(coin + "_total_sell_size")
	totalSellValue, _ := redisClient.GetFloat64(coin + "_total_sell_value")
	averageSellPrice := totalSellValue / totalSellSize

	var order *exchanges.OrderResp
	var err error
	var side string
	isOrderPlaced := false

	diff1 := totalSellSize - totalBuySize
	if diff1 > buySellDiffSize {
		// make buy
		side = "buy"
		isOrderPlaced = true

		// additionBuy should be equal to averageSellPrice for profit purpose
		// and the buy order could be filled quickly
		finalPrice := averageSellPrice
		if finalPrice > currentAskPrice {
			finalPrice = currentAskPrice
		}
		text := fmt.Sprintf("%s Make additionalBuy Size: %.3f Price: %.3f", coin, diff1, finalPrice)
		log.Println(text)
		go teleClient.SendMessage(text, chatID)
		finalPrice = utils.RoundTo(finalPrice, decimalsToRound)
		diff1 = utils.RoundTo(diff1, orderQuantityToRound)
		order, err = placeOrder("additionalBuy", diff1, finalPrice, "buy")
		if err != nil {
			text := fmt.Sprintf("%s %s Err Can not make order: %s", coin, "additionalBuy", err)
			go teleClient.SendMessage(text, chatErrorID)
		}
	}

	diff2 := totalBuySize - totalSellSize
	if diff2 > buySellDiffSize {
		//make sell
		side = "sell"
		isOrderPlaced = true

		// additionalSell should be equal to averageBuyPrice for profit purpose
		// and the sell order could be filled quickly
		finalPrice := averageBuyPrice
		if finalPrice < currentBidPrice {
			finalPrice = currentBidPrice
		}
		text := fmt.Sprintf("%s Make additionalSell Size: %.2f Price: %.3f", coin, diff2, finalPrice)
		log.Println(text)
		go teleClient.SendMessage(text, chatID)

		finalPrice = utils.RoundTo(finalPrice, decimalsToRound)
		diff2 = utils.RoundTo(diff2, orderQuantityToRound)
		order, err = placeOrder("additionalSell", diff2, finalPrice, "sell")
		if err != nil {
			text := fmt.Sprintf("%s %s Err Can not make order: %s", coin, "additionalSell", err)
			go teleClient.SendMessage(text, chatErrorID)
		}
	}

	if !isOrderPlaced {
		return
	}
	if order == nil {
		return
	}

	for i := 1; i <= 5; i++ {
		if !isStillGoodToMakeAdditionalBuySell() {
			_, err = exchangeClient.CancelOrder(coin, order.ID, order.ClientID)
			if err != nil {
				text := fmt.Sprintf("%s %s Err Can not cancel order: %d err: %s", coin, "additionalBuy/Sell", order.ID, err)
				teleClient.SendMessage(text, chatErrorID)
			}
			log.Printf("%s additional%s Order %d is canceled at price %f due to total buy/sell has changed", coin, side, order.ID, order.Price)
			break
		}
		orderDetails, err := exchangeClient.GetOrder(coin, order.ID, order.ClientID)
		if err != nil {
			text := fmt.Sprintf("%s %s Err getOrderDetails: %s", coin, "additionalBuy/Sell", err)
			log.Println(text)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(60 * time.Second)
			continue
		}

		log.Printf("%s additional%s Check Order %d status: %s", coin, side, orderDetails.ID, orderDetails.Status)
		if orderDetails.IsFilled() {
			calculateProfit(orderDetails.ID, orderDetails.ExecutedQty, orderDetails.Price, side)
			return // just return
		} else if orderDetails.IsCanceled() {
			log.Printf("%s additional%s Order %d is canceled at price %f", coin, side, orderDetails.ID, orderDetails.Price)
			break
		}
		if i == 5 {
			text := fmt.Sprintf("%s additional%s Order %d %.3f is not filled after 5 minutes will cancel it", coin, side, orderDetails.ID, orderDetails.Price)
			teleClient.SendMessage(text, chatID)
			_, err = exchangeClient.CancelOrder(coin, orderDetails.ID, orderDetails.ClientID)
			if err != nil {
				text := fmt.Sprintf("%s %s Err Can not cancel order: %d err: %s", coin, "additionalBuy/Sell", order.ID, err)
				teleClient.SendMessage(text, chatErrorID)
			}
			log.Printf("%s additional%s Order %d is canceled at price %f", coin, side, orderDetails.ID, orderDetails.Price)
		}
		time.Sleep(60 * time.Second)
	}
}

func isStillGoodToMakeAdditionalBuySell() bool {
	totalBuySize, _ := redisClient.GetFloat64(coin + "_total_buy_size")
	totalSellSize, _ := redisClient.GetFloat64(coin + "_total_sell_size")
	return math.Abs(totalSellSize-totalBuySize) > buySellDiffSize
}
