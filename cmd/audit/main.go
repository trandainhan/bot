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
		totalBuyValue, _ := redisClient.GetFloat64(coin + "_total_buy_value")
		totalSellValue, _ := redisClient.GetFloat64(coin + "_total_sell_value")
		fee := (totalBuyValue + totalSellValue) * 0.00075
		lastFee, err := redisClient.GetFloat64(coin + "_total_fee")
		if err != nil { // if total fee is not in redis yet
			lastFee = fee
		}
		todayFee := fee - lastFee
		text := fmt.Sprintf("%s\nToday Fee: %.4f", coin, todayFee)
		text = fmt.Sprintf("%s\nTotal Fee: %.4f", text, fee)
		teleClient.SendMessage(text, chatProfitID)
		redisClient.Set(coin+"_total_fee", fee, 0)
	}
	log.Println("=============")
	log.Println("Done")
}
