package main

import (
	"math"
)

func calculateBidFFromBidB(bidB, feePerBinance, perProfit, minPrice, maxPrice float64) (float64, bool) {
	rate, _ := redisClient.GetFloat64("usdtvnd_rate")
	bidF := (bidB * rate * (1 - feePerBinance)) / (1 + perProfit)
	bidF = math.Round(bidF)
	if bidF < maxPrice && bidF > minPrice {
		return bidF, false
	}
	return 0, true
}

func calculateAskFFromAskB(askB, perFeeBinance, perProfit, minPrice, maxPrice float64) (float64, bool) {
	rate, _ := redisClient.GetFloat64("usdtvnd_rate")
	askF := askB * (1 + perFeeBinance) / (1 - perProfit) * rate
	askF = math.Round(askF)
	if askF > minPrice && askF < maxPrice {
		return askF, false
	}
	return 0, true
}
