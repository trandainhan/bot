package exchanges

import (
	"strings"

	"gitlab.com/fiahub/bot/internal/exchanges/binance"
	"gitlab.com/fiahub/bot/internal/exchanges/ftx"
)

type ExchangeClient struct {
	Ftx *ftx.FtxClient
	Bn  *binance.Binance
}

type OrderResp struct {
	ID          int64   `json:"ID"`
	ClientID    string  `json:"ClientID"`
	OriginQty   float64 `json:"OriginQty"`
	ExecutedQty float64 `json:"ExecutedQty"`
	Status      string  `json:"Status"`
	Side        string  `json:"Side"`
	Price       float64 `json:"Price"`
}

func (or OrderResp) IsCanceled() bool {
	return or.Status == binance.ORDER_CANCELED || or.Status == ftx.ORDER_CLOSED
}

func (or OrderResp) IsFilled() bool {
	return or.Status == binance.ORDER_FILLED || or.Status == ftx.ORDER_CLOSED
}

func (or OrderResp) IsPartiallyFilled() bool {
	return or.Status == binance.ORDER_PARTIALLY_FILLED
}

func (or OrderResp) GetSide() string {
	return strings.ToLower(or.Side)
}
