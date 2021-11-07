package ftx

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gitlab.com/fiahub/bot/internal/utils"
)

const (
	ORDER_CLOSED = "closed"
)

func GetPriceByQuantity(marketParam string, quantity float64) (float64, float64, error) {
	orderBook, err := getOrderBook(marketParam, 100)
	if err != nil {
		return -1.0, -1.0, err // return negative price
	}
	totalQuantity := 0.0
	bidPriceByQuantity := 0.0
	for _, v := range orderBook.Bids {
		price := v[0]
		innerQuantity := v[1]
		totalQuantity = totalQuantity + innerQuantity
		if totalQuantity > quantity {
			bidPriceByQuantity = price
			break
		}
	}
	totalQuantity = 0.0
	askPriceByQuantity := 999999999999.0

	for _, v := range orderBook.Asks {
		price := v[0]
		innerQuantity := v[1]
		totalQuantity = totalQuantity + innerQuantity
		if totalQuantity > quantity {
			askPriceByQuantity = price
			break
		}
	}
	return bidPriceByQuantity, askPriceByQuantity, nil
}

type OrderBookResp struct {
	Success bool      `json:"success"`
	Result  OrderBook `json:"result:`
}

type OrderBook struct {
	Bids [][]float64 `json:"bids"`
	Asks [][]float64 `json:"asks"`
}

func getOrderBook(marketParam string, limit int) (*OrderBook, error) {
	var BASE_URL = os.Getenv("FTX_URL")
	url := fmt.Sprintf("%s/api/markets/%s/orderbook?depth=%d", BASE_URL, marketParam, limit)
	body, code, err := utils.HttpGet(url, nil)
	if err != nil {
		log.Printf("Err getOrderBook, StatusCode: %d, Err: %s", code, err.Error())
		return nil, err
	}
	var orderBookResp OrderBookResp
	err = json.Unmarshal([]byte(body), &orderBookResp)
	if err != nil {
		log.Printf("Err getOrderBook, can not unmarshal, with body: %s", body)
		return nil, err
	}
	return &orderBookResp.Result, nil
}

func (ftx FtxClient) GetOrder(marketParam string, orderId int64) (*Order, error) {
	path := fmt.Sprintf("/orders/%d", orderId)
	body, code, err := ftx.makeRequest("GET", path, "")
	if err != nil {
		log.Printf("Err GetOrder, StatusCode: %d, Err: %s", code, err.Error())
		return nil, err
	}
	var resp NewOrderResponse
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		log.Printf("Err GetOrder, can not unmarshal, with body: %s", body)
		return nil, err
	}
	return &resp.Result, nil
}

func (ftx FtxClient) GetAllOpenOrder(marketParam string) ([]Order, error) {
	path := fmt.Sprintf("/orders?market=%s", marketParam)
	body, code, err := ftx.makeRequest("GET", path, "")
	if err != nil {
		log.Printf("Err GetOrder, StatusCode: %d, Err: %s", code, err.Error())
		return nil, err
	}
	var resp OpenOrderResponse
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		log.Printf("Err GetOrder, can not unmarshal, with body: %s", body)
		return nil, err
	}
	return resp.Result, nil
}
