package exchanges

import (
	"os"

	"gitlab.com/fiahub/bot/internal/exchanges/binance"
	"gitlab.com/fiahub/bot/internal/exchanges/ftx"
)

func GetAskPriceByQuantity(coin string, quantity float64) (float64, error) {
	_, ask, err := getPriceByQuantity(coin, quantity)
	return ask, err
}

func GetBidPriceByQuantity(coin string, quantity float64) (float64, error) {
	bid, _, err := getPriceByQuantity(coin, quantity)
	return bid, err
}

func getPriceByQuantity(coin string, quantity float64) (float64, float64, error) {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		return ftx.GetPriceByQuantity(coin+"/USDT", quantity)
	}
	return binance.GetPriceByQuantity(coin+"USDT", quantity)
}
