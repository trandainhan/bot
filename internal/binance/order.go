package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"gitlab.com/fiahub/bot/internal/utils"
)

type OrderDetailsResp struct {
	OrderID       *string `json:"orderId"`
	ClientOrderID string  `json:"clientOrderId"`
	OriginQty     float64 `json:"origQty"`
	ExecutedQty   float64 `json:"executedQty"`
	Status        string  `json:"status"`
	Side          string  `json:"side"`
	Price         float64 `json:"price"`
}

func GetPriceByQuantity(marketParam string, quantity float64) (float64, float64) {
	orderBook := getOrderBook(marketParam, 100)
	totalQuantity := 0.0
	bidPriceByQuantity := 0.0
	for _, v := range orderBook.Bids {
		price, _ := strconv.ParseFloat(v[0], 64)
		innerQuantity, _ := strconv.ParseFloat(v[1], 64)
		totalQuantity = totalQuantity + innerQuantity
		if totalQuantity > quantity {
			bidPriceByQuantity = price
			break
		}
	}
	totalQuantity = 0.0
	askPriceByQuantity := 999999999999.0

	for _, v := range orderBook.Asks {
		price, _ := strconv.ParseFloat(v[0], 64)
		innerQuantity, _ := strconv.ParseFloat(v[1], 64)
		totalQuantity = totalQuantity + innerQuantity
		if totalQuantity > quantity {
			askPriceByQuantity = price
			break
		}
	}
	return bidPriceByQuantity, askPriceByQuantity
}

type OrderBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func getOrderBook(marketParam string, limit int) *OrderBook {
	var BASE_URL = os.Getenv("BINANCE_URL")
	url := fmt.Sprintf("%s/api/v3/depth?symbol=%s&limit=%d", BASE_URL, marketParam, limit)
	body, code, err := utils.HttpGet(url, nil)
	if err != nil {
		log.Printf("Err getOrderBook, StatusCode: %d, Err: %s", code, err.Error())
	}
	var orderBook OrderBook
	err = json.Unmarshal([]byte(body), &orderBook)
	if err != nil {
		panic(err)
	}
	return &orderBook
}

func (binance Binance) GetOrder(marketParam string, orderId string, originClientOrderID string) (*OrderDetailsResp, error) {
	params := map[string]string{
		"symbol":            marketParam,
		"orderId":           orderId,
		"origClientOrderId": originClientOrderID,
	}
	body, _, err := binance.makeRequest("GET", params, "/api/v3/order")
	if err != nil {
		return nil, err
	}
	var orderDetailsResp OrderDetailsResp
	err = json.Unmarshal([]byte(body), &orderDetailsResp)
	if err != nil {
		panic(err)
	}
	return &orderDetailsResp, nil
}
