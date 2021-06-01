package binance

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/utils"
)

func (binance Binance) SellLimit(marketParam string, price, quantity float64) (*OrderDetails, error) {
	params := map[string]string{
		"symbol":      marketParam,
		"side":        "SELL",
		"type":        "LIMIT",
		"timeInForce": "GTC",
		"quantity":    fmt.Sprintf("%f", quantity),
		"price":       fmt.Sprintf("%f", price),
	}
	return binance.makeTradeRequest(params)
}

func (binance Binance) BuyLimit(marketParam string, price float64, quantity float64) (*OrderDetails, error) {
	params := map[string]string{
		"symbol":      marketParam,
		"side":        "BUY",
		"type":        "LIMIT",
		"timeInForce": "GTC",
		"quantity":    fmt.Sprintf("%f", quantity),
		"price":       fmt.Sprintf("%f", price),
	}
	return binance.makeTradeRequest(params)
}

func (binance Binance) makeTradeRequest(params map[string]string) (*OrderDetails, error) {
	body, _, err := binance.makeRequest("POST", params, "/api/v3/order")
	if err != nil {

	}
	var order *OrderDetails
	err = json.Unmarshal([]byte(body), order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (binance Binance) makeRequest(httpType string, params map[string]string, postURL string) (string, int, error) {
	var BASE_URL = os.Getenv("binance_url") // https://api.binance.com
	now := time.Now()
	sec := now.UnixNano()
	timeDifferences := binance.RedisClient.Get("local_binance_time_difference").(int64)
	mili := (sec-5)/int64(time.Millisecond) + timeDifferences
	params["recvWindow"] = "59000"
	params["timestamp"] = fmt.Sprintf("%v", mili)

	url_with_para, _ := utils.BuildUrlWithParams(fmt.Sprintf("%s%s", BASE_URL, postURL), params)

	binanceAPIKey := os.Getenv("binance_api_key")
	binanceAPISecret := os.Getenv("binance_api_secret")

	hmac := utils.GenerateHmac(url_with_para, binanceAPISecret)
	params["signature"] = hmac

	final_url, _ := utils.BuildUrlWithParams(fmt.Sprintf("%s/order", BASE_URL), params)

	headers := map[string]string{
		"X-MBX-APIKEY": binanceAPIKey,
	}
	if httpType == "POST" {
		return utils.HttpPost(final_url, nil, &headers)
	}
	return utils.HttpGet(final_url, &headers)
}
