package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func validateFund() bool {
	teleHanlder := os.Getenv("TELEGRAM_HANDLER")

	usdtFund, err := exchangeClient.CheckFund("USDT")
	if err != nil {
		text := fmt.Sprintf("%s %s", teleHanlder, err)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}

	var text string
	minUSDTFund, _ := strconv.ParseFloat(os.Getenv("MIN_USDT_FUND"), 64)
	maxUSDTFund, _ := strconv.ParseFloat(os.Getenv("MAX_USDT_FUND"), 64)
	if usdtFund < minUSDTFund {
		redisClient.Set(coin+"_buy_worker_runable", false, 0)
		text := fmt.Sprintf("Update %s_buy_worker_runable to %v due to usdtFund is too low", coin, false)
		go teleClient.SendMessage(text, chatRunableID)
	}

	if usdtFund > maxUSDTFund {
		redisClient.Set(coin+"_sell_worker_runable", false, 0)
		text := fmt.Sprintf("Update %s_sell_worker_runable to %v due to usdtFund excceed limit", coin, false)
		teleClient.SendMessage(text, chatRunableID)
	}
	if usdtFund < minUSDTFund || usdtFund > maxUSDTFund {
		text = fmt.Sprintf("%s %s %s USDTFund: Out of range %.2f", teleHanlder, currentExchange, coin, usdtFund)
		teleClient.SendMessage(text, chatErrorID)
		return false
	}
	log.Printf("%s %s USDTFund: %.2f", currentExchange, coin, usdtFund)

	buyRunnable := redisClient.GetBool(coin + "_buy_worker_runable")
	if buyRunnable == false {
		redisClient.Set(coin+"_buy_worker_runable", true, 0)
		text := fmt.Sprintf("Reset %s_sell_worker_runable to %v", coin, true)
		teleClient.SendMessage(text, chatRunableID)
	}
	sellRunnable := redisClient.GetBool(coin + "_sell_worker_runable")
	if sellRunnable == false {
		redisClient.Set(coin+"_sell_worker_runable", true, 0)
		text := fmt.Sprintf("Reset %s_sell_worker_runable to %v", coin, false)
		teleClient.SendMessage(text, chatRunableID)
	}

	return true
}
