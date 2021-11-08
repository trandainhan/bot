package exchanges

import (
	"os"

	"gitlab.com/fiahub/bot/internal/exchanges/binance"
	"gitlab.com/fiahub/bot/internal/exchanges/ftx"
)

func GetAskPriceByQuantity(coin, fiat string, quantity float64) (float64, error) {
	_, ask, err := getPriceByQuantity(coin, fiat, quantity)
	return ask, err
}

func GetBidPriceByQuantity(coin, fiat string, quantity float64) (float64, error) {
	bid, _, err := getPriceByQuantity(coin, fiat, quantity)
	return bid, err
}

func getPriceByQuantity(coin, fiat string, quantity float64) (float64, float64, error) {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		return ftx.GetPriceByQuantity(coin+"/"+fiat, quantity)
	}
	return binance.GetPriceByQuantity(coin+fiat, quantity)
}
