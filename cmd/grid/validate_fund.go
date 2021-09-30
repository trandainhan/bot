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
	}
	if usdtFund < minUSDTFund || usdtFund > maxUSDTFund {
		text = fmt.Sprintf("%s %s %s USDTFund: Out of range %v", currentExchange, coin, teleHanlder, usdtFund)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}
	redisClient.Set(coin+"_buy_worker_runable", true, 0)
	log.Printf("%s %s USDTFund: %v", currentExchange, coin, usdtFund)
	return true
}
