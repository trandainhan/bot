package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"gitlab.com/fiahub/bot/internal/rediswrapper"
	"gitlab.com/fiahub/bot/internal/telegram"
)

func main() {
	log.Println("Start audit")
	chatProfitID, err := strconv.ParseInt(os.Getenv("CHAT_PROFIT_ID"), 10, 64)
	if err != nil {
		log.Panic("Missing ChatProfitID ENV Variable")
	}
	teleClient := telegram.NewTeleBot(os.Getenv("TELE_BOT_TOKEN"))

	coins := []string{"ALICE", "ATOM", "ADA"}
	ctx := context.Background()
	redisURL := os.Getenv("REDIS_URL")
	redisDBNum, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	redisClient := rediswrapper.NewRedisClient(ctx, redisURL, redisDBNum)

	for _, coin := range coins {
		log.Println(coin)

		totalBuySize, _ := redisClient.GetFloat64(coin + "_total_buy_size")
		totalBuyValue, _ := redisClient.GetFloat64(coin + "_total_buy_value")

		totalSellSize, _ := redisClient.GetFloat64(coin + "_total_sell_size")
		totalSellValue, _ := redisClient.GetFloat64(coin + "_total_sell_value")

		averageBuyPrice := totalBuyValue / totalBuySize
		averageSellPrice := totalSellValue / totalSellSize

		diffAvgPrice := averageSellPrice - averageBuyPrice
		unrealizedProfit := diffAvgPrice * (totalBuyValue + totalSellValue) / 2

		fee := (totalBuyValue + totalSellValue) * 0.00075
		lastFee, err := redisClient.GetFloat64(coin + "_total_fee")
		if err != nil { // if total fee is not in redis yet
			lastFee = fee
		}

		todayFee := fee - lastFee
		text := fmt.Sprintf("%s\nToday fee: %.4f", coin, todayFee)
		text = fmt.Sprintf("%s\nTotal fee: %.4f", text, fee)
		text = fmt.Sprintf("%s\nUnrealized profit: %.4f", text, unrealizedProfit)
		teleClient.SendMessage(text, chatProfitID)
		redisClient.Set(coin+"_total_fee", fee, 0)
	}
	log.Println("=============")
	log.Println("Done")
}
