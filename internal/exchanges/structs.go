package exchanges

import (
	"gitlab.com/fiahub/bot/internal/exchanges/binance"
	"gitlab.com/fiahub/bot/internal/exchanges/ftx"
)

type ExchangeClient struct {
	Ftx *ftx.FtxClient
	Bn  *binance.Binance
}

type OrderResp struct {
	ID          int64
	ClientID    string
	OriginQty   float64
	ExecutedQty float64
	Status      string
	Side        string
	Price       float64
}

func (or OrderResp) IsCanceled() bool {
	return or.Status == binance.ORDER_CANCELED
}

func (or OrderResp) IsFilled() bool {
	return or.Status == binance.ORDER_FILLED
}

func (or OrderResp) IsPartiallyFilled() bool {
	return or.Status == binance.ORDER_PARTIALLY_FILLED
}
