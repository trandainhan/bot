package ftx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/utils"
)

func (ftx FtxClient) SellLimit(marketParam string, price, quantity float64) (*Order, error) {
	newOrder := NewOrder{
		Market:     marketParam,
		Side:       "sell",
		Type:       "limit",
		Price:      price,
		Size:       quantity,
		ReduceOnly: false,
		Ioc:        false,
		PostOnly:   false,
	}
	requestBody, _ := json.Marshal(newOrder)
	return post("orders", requestBody)
}

func (ftx FtxClient) BuyLimit(marketParam string, price float64, quantity float64) (*Order, error) {
	newOrder := NewOrder{
		Market:     marketParam,
		Side:       "buy",
		Type:       "limit",
		Price:      price,
		Size:       quantity,
		ReduceOnly: false,
		Ioc:        false,
		PostOnly:   false,
	}
	requestBody, _ := json.Marshal(newOrder)
	return post("orders", requestBody)
}

func (ftx FtxClient) makeTradeRequest(data NewOrder) (*Order, error) {
	requestBody, _ := json.Marshal(data)
	body, _, err := ftx.makeRequest("POST", "/orders", string(requestBody))
	log.Println(body)
	if err != nil {
		log.Printf("Error ftx makeTradeRequest: %s", err.Error())
		return nil, err
	}
	var resp NewOrderResponse
	err = utils.ProcessResponse(body, &resp)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully make an order in ftx: %v", resp)
	return &resp.Result, nil
}

func (ftx FtxClient) makeRequest(method, path, requestBody string) (string, int, error) {
	BASE_URL := os.Getenv("FTX_URL")

	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api" + path + requestBody
	log.Println(signaturePayload)
	ftxAPIKey := os.Getenv("FTX_API_KEY")
	ftxAPISecret := os.Getenv("FTX_API_SECRET")

	hmac := utils.GenerateHmac(signaturePayload, ftxAPISecret)

	final_url := fmt.Sprintf("%s/api%s", BASE_URL, path)

	headers := map[string]string{
		"FTX-KEY":  ftxAPIKey,
		"FTX-TS":   ts,
		"FTX-SIGN": hmac,
	}
	if method == "POST" {
		return utils.HttpPost(final_url, requestBody, &headers)
	}
	return utils.HttpGet(final_url, &headers)
}

func post(path string, body []byte) (*Order, error) {
	BASE_URL := os.Getenv("FTX_URL")
	ftxAPIKey := os.Getenv("FTX_API_KEY")
	ftxAPISecret := os.Getenv("FTX_API_SECRET")

	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + "POST" + "/api/" + path + string(body)
	signature := utils.GenerateHmac(signaturePayload, ftxAPISecret)

	final_url := fmt.Sprintf("%s/api/%s", BASE_URL, path)

	req, _ := http.NewRequest("POST", final_url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", ftxAPIKey)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error processing response: %s", err.Error())
		return nil, err
	}
	var orderResp NewOrderResponse
	err = json.Unmarshal(body, &orderResp)
	if err != nil {
		return nil, err
	}
	return &orderResp.Result, nil
}
