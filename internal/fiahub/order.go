package fiahub

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/fiahub/bot/internal/telegram"
	u "gitlab.com/fiahub/bot/internal/utils"
)

var BASE_URL = os.Getenv("fiahub_url")

const (
	ORDER_CANCELLED = "cancelled"
	ORDER_FINISHED  = "finished"
)

type Order struct {
	Coin               string
	OriginalCoinAmount float64
	CoinAmount         float64
	PriceSellRandom    float64
	Currency           string
	Type               string
}

type OrderDetails struct {
	ID         string  `json:"id"`
	State      string  `json:"state"`
	CoinAmount float64 `json:"coin_amount"`
	Matching   bool    `json:"matching"`
}

func CancelAllOrder(token string) (string, int, error) {
	headers := &map[string]string{
		"access-token": token,
	}
	url := fmt.Sprintf("%s/orders/cancel_all", BASE_URL)

	resp, code, err := u.HttpPost(url, nil, headers)
	if err != nil {
		teleClient := telegram.NewTeleBot(os.Getenv("tele_fia_bot_token")) // 549902830:AAFcC-rqU5ErzwvDPfMcIKSJ6f6HzezWeUY
		text := fmt.Sprintf("%s \n resp: %s code: %d", url, resp, code)
		go teleClient.SendMessage(text, -307500490)
	}
	return resp, code, err
}

func CancelOrder(token string, orderID string) (*OrderDetails, int, error) {
	headers := &map[string]string{
		"access-token": token,
	}
	url := fmt.Sprintf("%s/orders/%s/cancel", BASE_URL, orderID)
	body, code, err := u.HttpPost(url, nil, headers)
	if err != nil {
		return nil, code, err
	}
	var order *OrderDetails
	err = json.Unmarshal([]byte(body), order)
	if err != nil {
		return nil, 500, err
	}
	return order, code, nil
}

func CreateAskOrder(token string, askOrder Order) (string, int, error) {
	return "", 200, nil
}

func GetOrderDetails(token string, orderID string) (*OrderDetails, int, error) {
	url := fmt.Sprintf("%s/orders/details/?token=%s&id=%s", BASE_URL, token, orderID)
	var order *OrderDetails
	body, code, err := u.HttpGet(url)
	if err != nil {
		return nil, code, err
	}
	err = json.Unmarshal([]byte(body), order)
	if err != nil {
		return nil, 500, err
	}
	return order, 200, nil
}
