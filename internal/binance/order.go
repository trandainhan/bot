package binance

import (
	"os"
)

var binanceURL = os.Getenv("binance_url")

type OrderDetails struct {
	ID            *string
	ClientOrderID *string
}

func GetPriceByQuantity(marketParam string, quantity int) (float64, float64) {
	return 0.0, 0.0

	// content := BinanceAPI_GetOrderBook(Marketpara,limit := 100)
	// parsed := JSON.Load(content)
	//
	//
	// Maxindex := parsed.bids.Maxindex()
	// TotalQuant := 0
	// BidPriceByQuantity := 0
	// Loop, %Maxindex%
	// {
	// 	price := parsed.bids[a_index][1]
	// 	quant := parsed.bids[a_index][2]
	// 	TotalQuant += quant
	//
	// 	If(TotalQuant >= Quantity)
	// 	{
	// 		BidPriceByQuantity := price
	// 		break
	// 	}
	// }
	//
	//
	// Maxindex := parsed.asks.Maxindex()
	// TotalQuant := 0
	// AskPriceByQuantity := 999999999999
	// Loop, %Maxindex%
	// {
	// 	price := parsed.asks[a_index][1]
	// 	quant := parsed.asks[a_index][2]
	// 	TotalQuant += quant
	//
	// 	If(TotalQuant >= Quantity)
	// 	{
	// 		AskPriceByQuantity := price
	// 		break
	// 	}
	// }
	//
	// return 1
}

func getOrderBook(marketParam string, limit int) {
	// link := "https://api.binance.com/api/v1/depth?symbol=" . Marketpara . "&limit=" . limit
	// content := GetContent(link)
	// return content
}

type OrderDetailsResp struct {
	OriginQty   float64 `json:"origQty"`
	ExecutedQty float64 `json:"executedQty"`
	Status      string  `json:"status"`
	Side        string  `json:"side"`
	Price       float64 `json:"price"`
}

func GetOrder(marketParam string, ID string, originClientID string) (OrderDetailsResp, error) {
	return OrderDetailsResp{}, nil
}
