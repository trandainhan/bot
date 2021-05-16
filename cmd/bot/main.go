package main

import (
	"context"
	"flag"
	"fmt"
	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/fiahub"
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
)

func init() {
	flag.StringVar(&coin, "coin", "btc", "Coin")
	flag.Float64Var(&minPrice, "minPrice", 0, "Min Price")
	flag.Float64Var(&maxPrice, "maxPrice", 1000000, "Min Price")
	flag.Int64Var(&defaultSleepSeconds, "defaultSleepSeconds", 18, "Sleep in second then restart")
	flag.IntVar(&decimalsToRound, "decimalsToRound", 3, "Decimal to round")
	flag.IntVar(&quantityToGetPrice, "quantityToGetPrice", 8, "Quantity To Get Price")
	flag.Parse()
}

func main() {
	// setup client
	ctx := context.Background()
	redisURL := os.Getenv("redis_url")
	redisClient = rediswrapper.NewRedisClient(ctx, redisURL)
	teleClient = telegram.NewTeleBot(os.Getenv("tele_bot_token"))

	// get environment for login
	email := os.Getenv("email")
	password := os.Getenv("password")
	fiahubToken := fiahub.Login(email, password)

	// Cancel all order before starting
	fiahub.CancelAllOrder(fiahubToken)
	time.Sleep(2 * time.Second)

	fiahub.RenewParam()
	fiahub.CalculatePerProfit()

	rate := fiahub.GetUSDVNDRate()
	redisClient.Set("usdtvnd_rate", rate)

	binance.GetOffsetTimeUnix()

	marketParam := coin + "USDT"
	perFeeBinance := redisClient.Get("per_fee_binance").(float64) // 0.075 / 100
	perProfit := redisClient.Get("per_profit").(float64)
	bidPriceByQuantity, askPriceByQuantity := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)
	log.Println(bidPriceByQuantity)
	log.Println(askPriceByQuantity)

	for {
		runable := redisClient.Get("runable").(bool)
		if !runable {
			time.Sleep(10 * time.Second)
			continue
		}

		riki1_runable := redisClient.Get("riki1_runable").(bool)
		if riki1_runable {
			askF, isOutRange := calculateAskFFromAskB(askPriceByQuantity, perFeeBinance, perProfit, minPrice, maxPrice)
			if isOutRange {
				text := fmt.Sprintf("%s @ndtan Error! Price out of range. PriceF: %v PriceB: %v Range: %v - %v", coin, askF, askPriceByQuantity, minPrice, maxPrice)
				chatID, _ := strconv.ParseInt(os.Getenv("chat_id"), 10, 64)
				go teleClient.SendMessage(text, chatID)
				time.Sleep(2 * time.Second)
			} else {
				id := "riki1"
				go worker(id, coin, askF, askPriceByQuantity, fiahubToken)
				// set riki1 worker to runing
			}
		}
		// riki2_runable := redisClient.Get("riki2_runable").(bool)
		// if riki2_runable {
		// 	go worker("riki2", 1)
		// }
		// riki3_runable := redisClient.Get("riki3_runable").(bool)
		// if riki3_runable {
		// 	go worker("riki3", 2)
		// }
		// riki4_runable := redisClient.Get("riki4_runable").(bool)
		// if riki4_runable {
		// 	go worker("riki4", 3)
		// }
	}
}
