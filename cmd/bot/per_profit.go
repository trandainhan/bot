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

func calculatePerProfit() bool {

	redisValue := redisClient.Get("coingiatot_params")
	var params fiahub.CoinGiaTotParams
	_ = json.Unmarshal([]byte(redisValue), &params)

	usdtFund := bn.CheckFund("USDT")
	perProfitBid := params.GetSpread()/2 + (usdtFund-params.GetUSDTOffset2()-params.GetUSDTMidPoint())/1000*params.GetProfitPerThousand()
	perProfitAsk := params.GetSpread() - perProfitBid

	perProfitAsk = utils.RoundTo(perProfitAsk, 6)
	perProfitBid = utils.RoundTo(perProfitBid, 6)

	oldPerProfitAsk := redisClient.GetFloat64("per_profit_ask")
	oldPerProfitBid := redisClient.GetFloat64("per_profit_bid")

	var text string
	teleHanlder := os.Getenv("TELEGRAM_HANDLER")
	minUSDTFund, _ := strconv.ParseFloat(os.Getenv("MIN_USDT_FUND"), 64)
	maxUSDTFund, _ := strconv.ParseFloat(os.Getenv("MAX_USDT_FUND"), 64)
	if usdtFund < minUSDTFund || usdtFund > maxUSDTFund {
		text = fmt.Sprintf("%s USDTFund: Out of range %v", teleHanlder, usdtFund)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}

	if perProfitBid < -0.16 || perProfitBid > 0.16 {
		text = fmt.Sprintf("%s PerProfitBid: Out of range %v", teleHanlder, perProfitBid)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}

	if perProfitAsk < -0.16 || perProfitAsk > 0.16 {
		text = fmt.Sprintf("%s PerProfitAsk: Out of range %v", teleHanlder, perProfitBid)
		go teleClient.SendMessage(text, chatErrorID)
		return false
	}

	isChange := false
	if perProfitAsk != oldPerProfitAsk {
		isChange = true
		redisClient.Set("per_profit_ask", perProfitAsk)
	}

	if perProfitBid != oldPerProfitBid {
		isChange = true
		redisClient.Set("per_profit_bid", perProfitBid)
	}

	if isChange {
		perProfitchange := math.Abs(oldPerProfitAsk - perProfitAsk)
		if perProfitchange < 0.001 {
			text = fmt.Sprintf("Per Profit Change: %.6f < 0.1%%", perProfitchange)
		} else {
			fia.CancelAllOrder()
			text = fmt.Sprintf("CancelAllOrder Per Profit Change: %.6f > 0.1%%", perProfitchange)
		}
		text = fmt.Sprintf("%s\n USDTFund: %.6f\n PerProfitAsk: %.6f -> %.6f\n PerProfitBid: %.6f -> %.6f",
			text, usdtFund, oldPerProfitAsk, perProfitAsk, oldPerProfitBid, perProfitBid)
		go teleClient.SendMessage(text, chatID)
	}
	return true
}
