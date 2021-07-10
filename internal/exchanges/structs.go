package exchanges

import (
	"gitlab.com/fiahub/bot/internal/exchanges/binance"
	"gitlab.com/fiahub/bot/internal/exchanges/ftx"
)

type ExchangeClient struct {
	Ftx ftx.FtxClient
	Bn  *binance.Binance
}
