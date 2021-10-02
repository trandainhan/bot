package main

import (
	"fmt"
	"log"
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

		finalPrice := averageBuyPrice
		if finalPrice > currentAskPrice {
			finalPrice = currentAskPrice
		}
		text := fmt.Sprintf("Make additionBuy Size: %.2f Price: %.2f", diff1, finalPrice)
		log.Println(text)
		go teleClient.SendMessage(text, chatID)
		finalPrice = utils.RoundTo(finalPrice, 2)
		diff1 = utils.RoundTo(diff1, 2)
		order, err = placeOrder("additionBuy", diff1, finalPrice, "buy")
		if err != nil {
			text := fmt.Sprintf("%s %s Err Can not make order: %s", coin, "additionBuy", err)
			go teleClient.SendMessage(text, chatErrorID)
		}
	}

	diff2 := totalBuySize - totalSellSize
	if diff2 > buySellDiffSize {
		//make sell
		side = "sell"
		isOrderPlaced = true

		finalPrice := averageSellPrice
		if finalPrice < currentBidPrice {
			finalPrice = currentBidPrice
		}
		text := fmt.Sprintf("Make additionSell Size: %.2f Price: %.2f", diff2, finalPrice)
		log.Println(text)
		go teleClient.SendMessage(text, chatID)

		finalPrice = utils.RoundTo(finalPrice, 2)
		diff2 = utils.RoundTo(diff2, 2)
		order, err = placeOrder("additionSell", diff2, finalPrice, "sell")
		if err != nil {
			text := fmt.Sprintf("%s %s Err Can not make order: %s", coin, "additionSell", err)
			go teleClient.SendMessage(text, chatErrorID)
		}
	}

	if !isOrderPlaced {
		return
	}
	if order == nil {
		return
	}

	for {
		orderDetails, err := exchangeClient.GetOrder(coin, order.ID, order.ClientID)
		if err != nil {
			text := fmt.Sprintf("%s %s Err getOrderDetails: %s", coin, "additionBuy/Sell", err)
			log.Println(text)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(60 * time.Second)
			continue
		}

		log.Printf("%s addition%s Check Order %d status: %s", coin, side, orderDetails.ID, orderDetails.Status)
		if orderDetails.IsFilled() {
			calculateProfit(orderDetails.ExecutedQty, orderDetails.Price, side)
		} else if orderDetails.IsCanceled() {
			log.Printf("%s addition%s Order %d is canceled at price %f", coin, side, orderDetails.ID, orderDetails.Price)
			break
		}
		time.Sleep(30 * time.Second)
	}
}
