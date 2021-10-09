package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func validateFund() bool {
	teleHanlder := os.Getenv("TELEGRAM_HANDLER")

	usdtFund, err := exchangeClient.CheckFund("USDT")
	if err != nil {
		text := fmt.Sprintf("%s %s", teleHanlder, err)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}

	minUSDTFund, _ := strconv.ParseFloat(os.Getenv("MIN_USDT_FUND"), 64)
	maxUSDTFund, _ := strconv.ParseFloat(os.Getenv("MAX_USDT_FUND"), 64)

	if usdtFund < minUSDTFund || usdtFund > maxUSDTFund {

		redisClient.Set(currentExchange+coin+"_worker_runable", false, 0)
		text := fmt.Sprintf("Update %s_worker_runable to %v due to usdtFund is too low", coin, false)
		go teleClient.SendMessage(text, chatRunableID)

		_, err := redisClient.GetTime(currentExchange + "_usdt_fund_notify_time")
		if err != nil { // mean the key is not existed, Only notify fund if it haven't been notified in 5 minutes
			redisClient.Set(currentExchange+"_usdt_fund_notify_time", time.Now(), time.Duration(5)*time.Minute)
			text = fmt.Sprintf("%s %s %s USDT Fund: Out of range %.3f", teleHanlder, currentExchange, coin, usdtFund)
			teleClient.SendMessage(text, chatErrorID)
		}
		return false
	}
	log.Printf("%s %s USDTFund: %.3f", currentExchange, coin, usdtFund)

	runnable := redisClient.GetBool(currentExchange + coin + "_worker_runable")
	if runnable == false {
		redisClient.Set(currentExchange+coin+"_worker_runable", true, 0)
		text := fmt.Sprintf("Reset %s_worker_runable to %v", coin, true)
		teleClient.SendMessage(text, chatRunableID)
	}
	return true
}
