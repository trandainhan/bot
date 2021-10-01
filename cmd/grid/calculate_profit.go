package main

import (
	"fmt"
	"log"
)

func calculate_profit(orderSize float64, price float64, side string) {
	if side == "buy" {
		totalBuySize, _ := redisClient.GetFloat64(coin + "_total_buy_size")
		totalBuySize = totalBuySize + orderSize
		redisClient.Set(coin+"_total_buy_size", totalBuySize, 0)

		totalBuyValue, _ := redisClient.GetFloat64(coin + "_total_buy_value")
		totalBuyValue = totalBuyValue + orderSize*price
		redisClient.Set(coin+"_total_buy_value", totalBuyValue, 0)
		text := fmt.Sprintf("%s totalBuySize: %.2f, totalBuyValue: %.2f", coin, totalBuySize, totalBuyValue)
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
		text := fmt.Sprintf("%s totalSellSize: %.2f, totalSellValue: %.2f", coin, totalSellSize, totalSellValue)
		log.Println(text)
		teleClient.SendMessage(text, chatProfitID)
	}
}
