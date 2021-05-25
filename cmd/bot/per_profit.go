package main

import (
	"encoding/json"
	"fmt"
	"math"

	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/utils"
)

func calculatePerProfit() bool {

	redisValue := redisClient.Get("coingiatot_params").(string)
	var params *fiahub.CoinGiaTotParams
	_ = json.Unmarshal([]byte(redisValue), params)

	usdtFund := binance.CheckFund("USDT")
	perProfitBid := params.Spread/2 + (usdtFund-params.USDTOffset2-params.USDTMidPoint)/1000*params.ProfitPerThousand
	perProfitAsk := params.Spread - perProfitBid

	perProfitAsk = utils.RoundTo(perProfitAsk, 6)
	perProfitBid = utils.RoundTo(perProfitBid, 6)

	oldPerProfitAsk := redisClient.Get("per_profit_ask").(float64)
	oldPerProfitBid := redisClient.Get("per_profit_bid").(float64)

	var text string
	if usdtFund < 2600 || usdtFund > 90000 {
		text = fmt.Sprintf("@ndtan USDTFund: Out of range %v", usdtFund)
		go teleClient.SendMessage(text, -465055332)
		return false
	}

	if perProfitBid < -0.16 || perProfitBid > 0.16 {
		text = fmt.Sprintf("@ndtan PerProfitBid: Out of range %v", perProfitBid)
		go teleClient.SendMessage(text, -465055332)
		return false
	}

	if perProfitAsk < -0.16 || perProfitAsk > 0.16 {
		text = fmt.Sprintf("@ndtan PerProfitAsk: Out of range %v", perProfitBid)
		go teleClient.SendMessage(text, -465055332)
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
			fiahubToken := redisClient.Get("fiahub_token").(string)
			fiahub.CancelAllOrder(fiahubToken)
			text = fmt.Sprintf("CancelAllOrder perProfitchange: %v > 0.1%%", perProfitchange)
		}
	}
	return true
}
