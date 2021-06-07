package fiahub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/telegram"
	u "gitlab.com/fiahub/bot/internal/utils"
)

var BASE_URL = os.Getenv("FIAHUB_URL")

const (
	ORDER_CANCELLED = "cancelled"
	ORDER_FINISHED  = "finished"
)

type Order struct {
	Coin               string  `json:"coin"`
	OriginalCoinAmount float64 `json:"original_coin_amount"`
	CoinAmount         float64 `json:"coin_amount"`
	PricePerUnitCents  float64 `json:"price_per_unit_cents"`
	Currency           string  `json:"currency"`
	Type               string  `json:"type"`
}

type OrderDetails struct {
	ID         int    `json:"id"`
	State      string `json:"state"`
	CoinAmount string `json:"coin_amount"`
	Matching   bool   `json:"matching"`
}

func (od OrderDetails) GetCoinAmount() float64 {
	res, _ := strconv.ParseFloat(od.CoinAmount, 64)
	return res
}

type CreateAskOrderResp struct {
	AskOrder OrderDetails `json:"ask_order"`
}

type CreateBidOrderResp struct {
	BidOrder OrderDetails `json:"bid_order"`
}

func (fiahub Fiahub) CancelAllOrder(token string) (string, int, error) {
	headers := &map[string]string{
		"access-token": token,
	}
	url := fmt.Sprintf("%s/orders/cancel_all", BASE_URL)

	now := time.Now()
	miliTime := now.UnixNano() / int64(time.Millisecond)
	fiahub.RedisClient.Set("lastest_cancel_all_time", miliTime)
	resp, code, err := u.HttpPost(url, nil, headers)
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ERROR_ID"), 10, 64)
	if err != nil {
		teleClient := telegram.NewTeleBot(os.Getenv("TELE_BOT_TOKEN"))
		text := fmt.Sprintf("%s \n resp: %s code: %d", url, resp, code)
		go teleClient.SendMessage(text, chatID)
	}
	log.Println("Successfully cancel all fiahub orders")
	return resp, code, err
}

func CancelOrder(token string, orderID int) (*OrderDetails, int, error) {
	headers := &map[string]string{
		"access-token": token,
	}
	url := fmt.Sprintf("%s/orders/%d/cancel", BASE_URL, orderID)
	body, code, err := u.HttpPost(url, nil, headers)
	if err != nil {
		return nil, code, err
	}
	var order OrderDetails
	err = json.Unmarshal([]byte(body), &order)
	if err != nil {
		return nil, 500, err
	}
	return &order, code, nil
}

func CreateAskOrder(token string, askOrder Order) (*OrderDetails, int, error) {
	headers := &map[string]string{
		"access-token": token,
	}
	url := fmt.Sprintf("%s/ask_orders", BASE_URL)

	data := map[string]Order{
		"ask_order": askOrder,
	}

	body, code, err := u.HttpPost(url, data, headers)
	if err != nil { // TODO: Improve it
		log.Printf("Err Fiahub Create Ask Order: %s", err.Error())
	}
	var resp CreateAskOrderResp
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return nil, 500, err
	}
	log.Printf("Successfully create fiahub ask order: %v", resp.AskOrder)
	return &resp.AskOrder, code, nil
}

func CreateBidOrder(token string, bidOrder Order) (*OrderDetails, int, error) {
	headers := &map[string]string{
		"access-token": token,
	}
	url := fmt.Sprintf("%s/bid_orders", BASE_URL)
	data := map[string]Order{
		"bid_order": bidOrder,
	}
	body, code, err := u.HttpPost(url, data, headers)
	if err != nil || code != 200 { // TODO: Improve it
		log.Printf("Err Fiahub Create Bid Order: %s Body: %s", err.Error(), body)
		return nil, code, err
	}
	var resp CreateBidOrderResp
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return nil, 500, err
	}
	log.Printf("Successfully create fiahub bid order: %v", resp.BidOrder)
	return &resp.BidOrder, code, nil
}

func GetOrderDetails(token string, orderID int) (*OrderDetails, int, error) {
	url := fmt.Sprintf("%s/orders/details/?token=%s&id=%d", BASE_URL, token, orderID)
	body, code, err := u.HttpGet(url, nil)
	if err != nil {
		return nil, code, err
	}
	var order OrderDetails
	err = json.Unmarshal([]byte(body), &order)
	if err != nil {
		return nil, 500, err
	}
	return &order, 200, nil
}
