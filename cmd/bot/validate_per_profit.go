package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/utils"
)

func validatePerProfit() bool {
	redisValue := redisClient.Get(currentExchange + "_coingiatot_params")
	var params fiahub.CoinGiaTotParams
	_ = json.Unmarshal([]byte(redisValue), &params)

	teleHanlder := os.Getenv("TELEGRAM_HANDLER")

	usdtFund, err := exchangeClient.CheckFund("USDT")
	if err != nil {
		text := fmt.Sprintf("%s %s", teleHanlder, err)
		go teleClient.SendMessage(text, chatErrorID)
		fia.CancelAllOrder()
		return false
	}

	perProfitBid := params.GetSpread()/2 + (usdtFund-params.GetUSDTOffset2()-params.GetUSDTMidPoint())/1000*params.GetProfitPerThousand()
	perProfitAsk := params.GetSpread() - perProfitBid

	perProfitAsk = utils.RoundTo(perProfitAsk, 6)
	perProfitBid = utils.RoundTo(perProfitBid, 6)

	askRedisKey := fmt.Sprintf("%s_%s_per_profit_ask", coin, currentExchange)
	bidRedisKey := fmt.Sprintf("%s_%s_per_profit_bid", coin, currentExchange)
	oldPerProfitAsk := redisClient.GetFloat64(askRedisKey)
	oldPerProfitBid := redisClient.GetFloat64(bidRedisKey)

	var text string
	minUSDTFund, _ := strconv.ParseFloat(os.Getenv("MIN_USDT_FUND"), 64)
	maxUSDTFund, _ := strconv.ParseFloat(os.Getenv("MAX_USDT_FUND"), 64)
	if usdtFund < minUSDTFund || usdtFund > maxUSDTFund {
		text = fmt.Sprintf("%s %s %s USDTFund: Out of range %v", currentExchange, coin, teleHanlder, usdtFund)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}

	if perProfitBid < -0.16 || perProfitBid > 0.16 {
		text = fmt.Sprintf("%s %s PerProfitBid: Out of range %v", coin, teleHanlder, perProfitBid)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}

	if perProfitAsk < -0.16 || perProfitAsk > 0.16 {
		text = fmt.Sprintf("%s %s PerProfitAsk: Out of range %v", coin, teleHanlder, perProfitBid)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}

	isChange := false
	if perProfitAsk != oldPerProfitAsk {
		isChange = true
		redisClient.Set(askRedisKey, perProfitAsk)
	}

	if perProfitBid != oldPerProfitBid {
		isChange = true
		redisClient.Set(bidRedisKey, perProfitBid)
	}

	if isChange {
		perProfitchange := math.Abs(oldPerProfitAsk - perProfitAsk)
		if perProfitchange < 0.001 {
			text = fmt.Sprintf("%s Per Profit Change: %.6f < 0.1%%", coin, perProfitchange)
		} else {
			fia.CancelAllOrder()
			text = fmt.Sprintf("%s CancelAllOrder Per Profit Change: %.6f > 0.1%%", coin, perProfitchange)
		}
		text = fmt.Sprintf("%s\n USDTFund: %.6f\n PerProfitAsk: %.6f -> %.6f\n PerProfitBid: %.6f -> %.6f",
			text, usdtFund, oldPerProfitAsk, perProfitAsk, oldPerProfitBid, perProfitBid)
		go teleClient.SendMessage(text, chatID)
	}
	return true
}
