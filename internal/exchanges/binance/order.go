package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"gitlab.com/fiahub/bot/internal/utils"
)

const (
	ORDER_FILLED = "FILLED"
)

type OrderDetailsResp struct {
	OrderID       int64  `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	OriginQty     string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	Status        string `json:"status"`
	Side          string `json:"side"`
	Price         string `json:"price"`
}

func (od OrderDetailsResp) GetPrice() float64 {
	res, _ := strconv.ParseFloat(od.Price, 64)
	return res
}

func (od OrderDetailsResp) GetOriginQty() float64 {
	res, _ := strconv.ParseFloat(od.OriginQty, 64)
	return res
}

func (od OrderDetailsResp) GetExecutedQty() float64 {
	res, _ := strconv.ParseFloat(od.ExecutedQty, 64)
	return res
}

func GetPriceByQuantity(marketParam string, quantity float64) (float64, float64, error) {
	orderBook, err := getOrderBook(marketParam, 100)
	if err != nil {
		return -1.0, -1.0, err // return negative price
	}
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
	return bidPriceByQuantity, askPriceByQuantity, nil
}

type OrderBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func getOrderBook(marketParam string, limit int) (*OrderBook, error) {
	var BASE_URL = os.Getenv("BINANCE_URL")
	url := fmt.Sprintf("%s/api/v3/depth?symbol=%s&limit=%d", BASE_URL, marketParam, limit)
	body, code, err := utils.HttpGet(url, nil)
	if err != nil {
		log.Printf("Err getOrderBook, StatusCode: %d, Err: %s", code, err.Error())
		return nil, err
	}
	var orderBook OrderBook
	err = json.Unmarshal([]byte(body), &orderBook)
	if err != nil {
		log.Printf("Err getOrderBook, can not unmarshal, with body: %s", body)
		return nil, err
	}
	return &orderBook, nil
}

func (binance Binance) GetOrder(marketParam string, orderId int64, originClientOrderID string) (*OrderDetailsResp, error) {
	params := map[string]string{
		"symbol":            marketParam,
		"orderId":           strconv.FormatInt(orderId, 10),
		"origClientOrderId": originClientOrderID,
	}
	body, code, err := binance.makeRequest("GET", params, "/api/v3/order")
	if err != nil {
		log.Printf("Err GetOrder, StatusCode: %d, Err: %s", code, err.Error())
		return nil, err
	}
	var orderDetailsResp OrderDetailsResp
	err = json.Unmarshal([]byte(body), &orderDetailsResp)
	if err != nil {
		log.Printf("Err GetOrder, can not unmarshal, with body: %s", body)
		return nil, err
	}
	return &orderDetailsResp, nil
}
