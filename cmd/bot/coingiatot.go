package main

import (
	"encoding/json"
	"fmt"

	"gitlab.com/fiahub/bot/internal/fiahub"
)

func validateCoinGiaTotParams(params *fiahub.CoinGiaTotParams) bool {
	result := true
	if params.AutoMode != 0 && params.AutoMode != 1 {
		text := "@ndtan AutoMode: Out of range"
		go teleClient.SendMessage(text, -465055332)
		result = false
	}
	if params.ProfitMax < 0 || params.ProfitMax > 1 {
		text := "@ndtan ProfitMax: Out of range"
		go teleClient.SendMessage(text, -465055332)
		result = false
	}
	if params.ProfitPerThousand < 0 || params.ProfitPerThousand > 0.004 {
		text := "@ndtan ProfitPerThousand: Out of range"
		go teleClient.SendMessage(text, -465055332)
		result = false
	}
	if params.Spread <= 0 || params.Spread > 0.1 {
		text := "@ndtan Spead: Out of range"
		go teleClient.SendMessage(text, -465055332)
		result = false
	}
	if params.USDTMax < 0 || params.USDTMax > 60000 {
		text := "@ndtan USDTMax: Out of range"
		go teleClient.SendMessage(text, -465055332)
		result = false
	}
	if params.USDTMidPoint < 0 || params.USDTMidPoint > 60000 {
		text := "@ndtan USDTMidPoint: Out of range"
		go teleClient.SendMessage(text, -465055332)
		result = false
	}
	if params.USDTOffset2 < -30000 || params.USDTOffset2 > 240000 {
		text := "@ndtan USDTOffset2: Out of range"
		go teleClient.SendMessage(text, -465055332)
		result = false
	}
	return result
}

func renewCoinGiaTotParams(params *fiahub.CoinGiaTotParams) bool {
	isChange := false
	redisValue := redisClient.Get("coingiatot_params").(string)
	var oldParams *fiahub.CoinGiaTotParams
	_ = json.Unmarshal([]byte(redisValue), oldParams)

	if params.AutoMode != oldParams.AutoMode ||
		params.ProfitMax != oldParams.ProfitMax ||
		params.ProfitPerThousand != oldParams.ProfitPerThousand ||
		params.Spread != oldParams.Spread ||
		params.USDTMax != oldParams.USDTMax ||
		params.USDTMidPoint != oldParams.USDTMax ||
		params.USDTOffset2 != oldParams.USDTOffset2 {
		isChange = true
	}
	jsonParams, _ := json.Marshal(params)
	redisClient.Set("coingiatot_params", string(jsonParams))
	if isChange {
		autoMode := fmt.Sprintf("AutoMode: %d -> %d", oldParams.AutoMode, params.AutoMode)
		profitMax := fmt.Sprintf("ProfitMax: %v -> %v", oldParams.ProfitMax, params.ProfitMax)
		profitPerThousand := fmt.Sprintf("ProfitPerThousand: %v -> %v", oldParams.ProfitPerThousand, params.ProfitPerThousand)
		spread := fmt.Sprintf("Spread: %v -> %v", oldParams.Spread, params.Spread)
		usdtMax := fmt.Sprintf("USDTMax: %v -> %v", oldParams.USDTMax, params.USDTMax)
		usdtMidPoint := fmt.Sprintf("USDTMidPoint: %v -> %v", oldParams.USDTMidPoint, params.USDTMidPoint)
		offset := fmt.Sprintf("USDTOffset2: %v -> %v", oldParams.USDTOffset2, params.USDTOffset2)

		text := fmt.Sprintf("@ndtan AutoMode Params: \n %s\n %s\n %s\n %s\n %s\n %s\n %s",
			autoMode, profitMax, profitPerThousand, spread, usdtMax, usdtMidPoint, offset)
		go teleClient.SendMessage(text, -465055332)
	}
	return isChange
}
