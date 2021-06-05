package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/utils"
)

func calculatePerProfit() bool {

	redisValue := redisClient.Get("coingiatot_params")
	var params fiahub.CoinGiaTotParams
	_ = json.Unmarshal([]byte(redisValue), &params)

	bn := binance.Binance{
		RedisClient: redisClient,
	}

	usdtFund := bn.CheckFund("USDT")
	perProfitBid := params.Spread/2 + (usdtFund-params.USDTOffset2-params.USDTMidPoint)/1000*params.ProfitPerThousand
	perProfitAsk := params.Spread - perProfitBid

	perProfitAsk = utils.RoundTo(perProfitAsk, 6)
	perProfitBid = utils.RoundTo(perProfitBid, 6)

	oldPerProfitAsk := redisClient.GetFloat64("per_profit_ask")
	oldPerProfitBid := redisClient.GetFloat64("per_profit_bid")

	var text string
	teleHanlder := os.Getenv("TELEGRAM_HANDLER")
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if usdtFund < 2600 || usdtFund > 90000 {
		text = fmt.Sprintf("%s USDTFund: Out of range %v", teleHanlder, usdtFund)
		go teleClient.SendMessage(text, chatID)
		return false
	}

	if perProfitBid < -0.16 || perProfitBid > 0.16 {
		text = fmt.Sprintf("%s PerProfitBid: Out of range %v", teleHanlder, perProfitBid)
		go teleClient.SendMessage(text, chatID)
		return false
	}

	if perProfitAsk < -0.16 || perProfitAsk > 0.16 {
		text = fmt.Sprintf("%s PerProfitAsk: Out of range %v", teleHanlder, perProfitBid)
		go teleClient.SendMessage(text, chatID)
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
			text = fmt.Sprintf("perProfitchange: %v < 0.1%%", perProfitchange)
		} else {
			fiahubToken := redisClient.Get("fiahub_token")
			fia.CancelAllOrder(fiahubToken)
			text = fmt.Sprintf("CancelAllOrder perProfitchange: %v > 0.1%%", perProfitchange)
		}
	}
	return true
}
