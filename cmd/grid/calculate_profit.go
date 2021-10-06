package main

import (
	"fmt"
	"log"
)

func calculateProfit(orderID int64, orderSize float64, price float64, side string) {
	text := fmt.Sprintf("%s %s Order %d size: %.3f is filled at price %.3f", coin, side, orderID, orderSize, price)

	totalBuySize, _ := redisClient.GetFloat64(coin + "_total_buy_size")
	totalBuyValue, _ := redisClient.GetFloat64(coin + "_total_buy_value")

	totalSellSize, _ := redisClient.GetFloat64(coin + "_total_sell_size")
	totalSellValue, _ := redisClient.GetFloat64(coin + "_total_sell_value")

	averageBuyPrice := totalBuyValue / totalBuySize
	averageSellPrice := totalSellValue / totalSellSize

	if side == "buy" {
		totalBuySize = totalBuySize + orderSize
		redisClient.Set(coin+"_total_buy_size", totalBuySize, 0)
		totalBuyValue = totalBuyValue + orderSize*price
		redisClient.Set(coin+"_total_buy_value", totalBuyValue, 0)
		averageBuyPrice = totalBuyValue / totalBuySize

		log.Println(text)
	}

	if side == "sell" {
		totalSellSize = totalSellSize + orderSize
		redisClient.Set(coin+"_total_sell_size", totalSellSize, 0)

		totalSellValue = totalSellValue + orderSize*price
		redisClient.Set(coin+"_total_sell_value", totalSellValue, 0)
		averageSellPrice = totalSellValue / totalSellSize

		log.Println(text)
	}

	text = text + "\nTotal Buy"
	text = fmt.Sprintf("%s\nSize: %.3f, Value: %.3f, avgPrice: %.4f", text, totalBuySize, totalBuyValue, averageBuyPrice)
	text = text + "\nTotal Sell"
	text = fmt.Sprintf("%s\nSize: %.3f, Value: %.3f, avgPrice: %.4f", text, totalSellSize, totalSellValue, averageSellPrice)
	text = fmt.Sprintf("%s\n Diff Avg Price: %.4f", text, averageSellPrice-averageBuyPrice)
	teleClient.SendMessage(text, chatProfitID)
}
