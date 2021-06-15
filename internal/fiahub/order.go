package fiahub

import (
	"encoding/json"
	"errors"
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

type CancelOrderResp struct {
	Order OrderDetails `json:"order"`
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

func (fiahub *Fiahub) CancelAllOrder() (string, int, error) {
	headers := &map[string]string{
		"access-token": fiahub.Token,
	}
	url := fmt.Sprintf("%s/orders/cancel_all", BASE_URL)

	now := time.Now()
	miliTime := now.UnixNano() / int64(time.Millisecond)
	fiahub.SetCancelTime(miliTime)
	resp, code, err := u.HttpPost(url, nil, headers)
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ERROR_ID"), 10, 64)
	if err != nil {
		teleClient := telegram.NewTeleBot(os.Getenv("TELE_BOT_TOKEN"))
		text := fmt.Sprintf("%s \n resp: %s code: %d", url, resp, code)
		go teleClient.SendMessage(text, chatID)
	}
	return resp, code, err
}

func (fia Fiahub) CancelOrder(orderID int) (*OrderDetails, int, error) {
	headers := &map[string]string{
		"access-token": fia.Token,
	}
	url := fmt.Sprintf("%s/orders/%d/cancel", BASE_URL, orderID)
	body, code, err := u.HttpPost(url, nil, headers)
	if err != nil {
		return nil, code, err
	}
	if code >= 400 {
		return nil, code, errors.New(fmt.Sprintf("Status code %d %s", code, body))
	}
	var resp CancelOrderResp
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return nil, 0, err
	}
	return &resp.Order, code, nil
}

func (fia Fiahub) CreateAskOrder(askOrder Order) (*OrderDetails, int, error) {
	headers := &map[string]string{
		"access-token": fia.Token,
	}
	url := fmt.Sprintf("%s/ask_orders", BASE_URL)

	data := map[string]Order{
		"ask_order": askOrder,
	}

	body, code, err := u.HttpPost(url, data, headers)
	if err != nil {
		return nil, code, err
	}
	if code >= 400 {
		return nil, code, errors.New(fmt.Sprintf("Status code %d %s", code, body))
	}
	var resp CreateAskOrderResp
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return nil, 0, err
	}
	log.Printf("Successfully create fiahub ask order: %v", resp.AskOrder)
	return &resp.AskOrder, code, nil
}

func (fia Fiahub) CreateBidOrder(bidOrder Order) (*OrderDetails, int, error) {
	headers := &map[string]string{
		"access-token": fia.Token,
	}
	url := fmt.Sprintf("%s/bid_orders", BASE_URL)
	data := map[string]Order{
		"bid_order": bidOrder,
	}
	body, code, err := u.HttpPost(url, data, headers)
	if err != nil {
		return nil, code, err
	}
	if code >= 400 {
		return nil, code, errors.New("Status code >= 200" + body)
	}
	var resp CreateBidOrderResp
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return nil, 0, err
	}
	log.Printf("Successfully create fiahub bid order: %v", resp.BidOrder)
	return &resp.BidOrder, code, nil
}

func (fia Fiahub) GetAskOrderDetails(orderID int) (*OrderDetails, int, error) {
	body, code, err := getOrderDetails(fia.Token, orderID)
	if err != nil {
		return nil, code, err
	}
	var order CreateAskOrderResp
	err = json.Unmarshal([]byte(body), &order)
	if err != nil {
		return nil, 0, err
	}
	return &order.AskOrder, code, nil
}

func (fia Fiahub) GetBidOrderDetails(orderID int) (*OrderDetails, int, error) {
	body, code, err := getOrderDetails(fia.Token, orderID)
	if err != nil {
		return nil, code, err
	}
	var order CreateBidOrderResp
	err = json.Unmarshal([]byte(body), &order)
	if err != nil {
		return nil, 0, err
	}
	return &order.BidOrder, code, nil
}

func getOrderDetails(token string, orderID int) (string, int, error) {
	headers := &map[string]string{
		"access-token": token,
	}
	url := fmt.Sprintf("%s/orders/detail?id=%d", BASE_URL, orderID)
	body, code, err := u.HttpGet(url, headers)
	if err != nil {
		return "", code, err
	}
	if code >= 400 {
		return "", code, errors.New(fmt.Sprintf("Error GetOrderDetails, body: %s, code: %d", body, code))
	}
	return body, code, err
}
