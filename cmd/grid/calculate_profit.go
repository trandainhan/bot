package main

import (
	"fmt"
	"log"
)

func calculateProfit(orderSize float64, price float64, side string) {
	text := fmt.Sprintf("%s %s Order is filled at price %.2f", coin, side, price)

	if side == "buy" {
		totalBuySize, _ := redisClient.GetFloat64(coin + "_total_buy_size")
		totalBuySize = totalBuySize + orderSize
		redisClient.Set(coin+"_total_buy_size", totalBuySize, 0)

		totalBuyValue, _ := redisClient.GetFloat64(coin + "_total_buy_value")
		totalBuyValue = totalBuyValue + orderSize*price
		redisClient.Set(coin+"_total_buy_value", totalBuyValue, 0)
		averageBuyPrice := totalBuyValue / totalBuySize

		text = fmt.Sprintf("%s\n totalBuySize: %.2f, totalBuyValue: %.2f, averageBuyPrice: %.2f", text, totalBuySize, totalBuyValue, averageBuyPrice)
		log.Println(text)
		teleClient.SendMessage(text, chatProfitID)
	}

	if side == "sell" {
		totalSellSize, _ := redisClient.GetFloat64(coin + "_total_sell_size")
		totalSellSize = totalSellSize + orderSize
		redisClient.Set(coin+"_total_sell_size", totalSellSize, 0)

		totalSellValue, _ := redisClient.GetFloat64(coin + "_total_sell_value")
		totalSellValue = totalSellValue + orderSize*price
		redisClient.Set(coin+"_total_sell_value", totalSellValue, 0)
		averageSellPrice := totalSellValue / totalSellSize
		text = fmt.Sprintf("%s\n totalSellSize: %.2f, totalSellValue: %.2f, averageSellPrice: %.2f", text, totalSellSize, totalSellValue, averageSellPrice)
		log.Println(text)
		teleClient.SendMessage(text, chatProfitID)
	}
}
