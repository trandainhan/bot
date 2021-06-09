package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"gitlab.com/fiahub/bot/internal/fiahub"
)

func validateCoinGiaTotParams(params *fiahub.CoinGiaTotParams) bool {
	result := true
	teleHanlder := os.Getenv("TELEGRAM_HANDLER")
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if params.GetAutoMode() != 0 && params.GetAutoMode() != 1 {
		text := fmt.Sprintf("%s AutoMode: Out of range", teleHanlder)
		go teleClient.SendMessage(text, -465055332)
		result = false
	}
	if params.GetProfitMax() < 0 || params.GetProfitMax() > 1 {
		text := fmt.Sprintf("%s ProfitMax: Out of range", teleHanlder)
		go teleClient.SendMessage(text, chatID)
		result = false
	}
	if params.GetProfitPerThousand() < 0 || params.GetProfitPerThousand() > 0.004 {
		text := fmt.Sprintf("%s ProfitPerThousand: Out of range", teleHanlder)
		go teleClient.SendMessage(text, chatID)
		result = false
	}
	if params.GetSpread() <= 0 || params.GetSpread() > 0.1 {
		text := fmt.Sprintf("%s Spead: Out of range", teleHanlder)
		go teleClient.SendMessage(text, chatID)
		result = false
	}
	if params.GetUSDTMax() < 0 || params.GetUSDTMax() > 60000 {
		text := fmt.Sprintf("%s USDTMax: Out of range", teleHanlder)
		go teleClient.SendMessage(text, chatID)
		result = false
	}
	if params.GetUSDTMidPoint() < 0 || params.GetUSDTMidPoint() > 60000 {
		text := fmt.Sprintf("%s USDTMidPoint: Out of range", teleHanlder)
		go teleClient.SendMessage(text, chatID)
		result = false
	}
	if params.GetUSDTOffset2() < -30000 || params.GetUSDTOffset2() > 240000 {
		text := fmt.Sprintf("%s USDTOffset2: Out of range", teleHanlder)
		go teleClient.SendMessage(text, chatID)
		result = false
	}
	return result
}

func renewCoinGiaTotParams(params *fiahub.CoinGiaTotParams) bool {
	isChange := false
	redisValue := redisClient.Get("coingiatot_params")
	var oldParams fiahub.CoinGiaTotParams
	_ = json.Unmarshal([]byte(redisValue), &oldParams)

	if params.AutoMode != oldParams.AutoMode ||
		params.ProfitMax != oldParams.ProfitMax ||
		params.ProfitPerThousand != oldParams.ProfitPerThousand ||
		params.Spread != oldParams.Spread ||
		params.USDTMax != oldParams.USDTMax ||
		params.USDTMidPoint != oldParams.USDTMidPoint ||
		params.USDTOffset2 != oldParams.USDTOffset2 {
		isChange = true
	}
	if isChange {
		chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
		jsonParams, _ := json.Marshal(params)
		redisClient.Set("coingiatot_params", string(jsonParams))

		autoMode := fmt.Sprintf("AutoMode: %s -> %s", oldParams.AutoMode, params.AutoMode)
		profitMax := fmt.Sprintf("ProfitMax: %v -> %v", oldParams.ProfitMax, params.ProfitMax)
		profitPerThousand := fmt.Sprintf("ProfitPerThousand: %v -> %v", oldParams.ProfitPerThousand, params.ProfitPerThousand)
		spread := fmt.Sprintf("Spread: %v -> %v", oldParams.Spread, params.Spread)
		usdtMax := fmt.Sprintf("USDTMax: %v -> %v", oldParams.USDTMax, params.USDTMax)
		usdtMidPoint := fmt.Sprintf("USDTMidPoint: %v -> %v", oldParams.USDTMidPoint, params.USDTMidPoint)
		offset := fmt.Sprintf("USDTOffset2: %v -> %v", oldParams.USDTOffset2, params.USDTOffset2)

		text := fmt.Sprintf("%s AutoMode Params: \n %s\n %s\n %s\n %s\n %s\n %s\n %s", os.Getenv("TELEGRAM_HANDLER"),
			autoMode, profitMax, profitPerThousand, spread, usdtMax, usdtMidPoint, offset)
		go teleClient.SendMessage(text, chatID)
	}
	return isChange
}
