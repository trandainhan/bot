package main

import (
	"fmt"
	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/rediswrapper"
	"gitlab.com/fiahub/bot/internal/telegram"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	coin                string
	minPrice            float64
	maxPrice            float64
	redisClient         *rediswrapper.MyRedis
	teleClient          *telegram.TeleBot
	decimalsToRound     int
	defaultSleepSeconds int64
	quantityToGetPrice  int
	fiahubToken         string
)

func main() {
	perFeeBinance := redisClient.Get("per_fee_binance").(float64) // 0.075 / 100
	perProfitAsk := redisClient.Get("per_profit_ask").(float64)
	perProfitBid := redisClient.Get("per_profit_bid").(float64)

	marketParam := coin + "USDT"
	bidPriceByQuantity, askPriceByQuantity := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)

	log.Println(bidPriceByQuantity)
	log.Println(askPriceByQuantity)
	log.Println(perProfitBid)

	// go renew params

	var perProfitStep float64
	// go ask_worker
	perProfitStep = 1.0
	go ask_worker("riki1", coin, askPriceByQuantity, perProfitStep)
	// go ask_worker
	// go ask_worker
	// go ask_worker

	// go bid_worker
	// go bid_worker
	// go bid_worker
	// go bid_worker

	for {
		runable := redisClient.Get("runable").(bool)
		if !runable {
			time.Sleep(10 * time.Second)
			continue
		}

		riki1_runable := redisClient.Get("riki1_runable").(bool)
		if riki1_runable {
			askF, isOutRange := calculateAskFFromAskB(askPriceByQuantity, perFeeBinance, perProfitAsk, minPrice, maxPrice)
			if isOutRange {
				text := fmt.Sprintf("%s @ndtan Error! Price out of range. PriceF: %v PriceB: %v Range: %v - %v", coin, askF, askPriceByQuantity, minPrice, maxPrice)
				chatID, _ := strconv.ParseInt(os.Getenv("chat_id"), 10, 64)
				go teleClient.SendMessage(text, chatID)
				time.Sleep(2 * time.Second)
			} else {
				id := "riki1"
				perProfitStep := 1.0
				go worker(id, coin, askF, askPriceByQuantity, perProfitStep)
				redisClient.Set("riki1_runable", false)
			}
		}

		riki2_runable := redisClient.Get("riki2_runable").(bool)
		if riki2_runable {
			askF, isOutRange := calculateAskFFromAskB(askPriceByQuantity, perFeeBinance, perProfitAsk, minPrice, maxPrice)
			if isOutRange {
				text := fmt.Sprintf("%s @ndtan Error! Price out of range. PriceF: %v PriceB: %v Range: %v - %v", coin, askF, askPriceByQuantity, minPrice, maxPrice)
				chatID, _ := strconv.ParseInt(os.Getenv("chat_id"), 10, 64)
				go teleClient.SendMessage(text, chatID)
				time.Sleep(2 * time.Second)
			} else {
				id := "riki2"
				perProfitStep := 2.0
				go worker(id, coin, askF, askPriceByQuantity, perProfitStep)
				redisClient.Set("riki2_runable", false)
			}
		}
	}
}
