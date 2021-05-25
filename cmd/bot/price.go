package main

import (
	// "gitlab.com/fiahub/bot/internal/rediswrapper"
	"math"
)

func calculateBidFFromBidB(bidB, feePerBinance, perProfit, minPrice, maxPrice float64) bool {
	rate := redisClient.Get("usdtvnd_rate").(float64)
	bidF := (bidB * rate * (1 - feePerBinance)) / (1 + perProfit)
	bidF = math.Round(bidF)
	if bidF < maxPrice && bidF > minPrice {
		return true
	}
	return false
}

func calculateAskFFromAskB(askB, perFeeBinance, perProfit, minPrice, maxPrice float64) (float64, bool) {
	rate := redisClient.Get("usdtvnd_rate").(float64)
	askF := askB * (1 + perFeeBinance) / (1 - perProfit) * rate
	askF = math.Round(askF)
	if askF > minPrice && askF < maxPrice {
		return askF, false
	}
	return 0, true
}
