package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Balance struct {
	Asset      string `json:"asset"`
	Free       string `json:"free"`
	LockedFund string `json:"locked"`
}

func (balance Balance) GetFree() float64 {
	res, _ := strconv.ParseFloat(balance.Free, 64)
	return res
}

func (balance Balance) GetLockedFund() float64 {
	res, _ := strconv.ParseFloat(balance.LockedFund, 64)
	return res
}

type Fund struct {
	Balances []Balance `json:"balances"`
}

func (binance Binance) CheckFund(name string) (float64, error) {
	fund, err := binance.checkFund()
	if err != nil {
		return -1.0, err
	}
	var result float64
	for _, balance := range fund.Balances {
		if balance.Asset == name {
			result = balance.GetFree()
		}
	}
	return result, nil
}

func (binance Binance) checkFund() (*Fund, error) {
	params := make(map[string]string)
	body, code, err := binance.makeRequest("GET", params, "/api/v3/account")
	if err != nil {
		log.Printf("Err checkFund, statusCode: %d err: %s", code, err.Error())
		return nil, err
	}
	if os.Getenv("LOG_DEBUGGER") == "true" {
		log.Printf("checkFund %s", body)
	}
	var fund Fund
	err = json.Unmarshal([]byte(body), &fund)
	if err != nil {
		panic(err)
	}
	return &fund, nil
}

func (binance Binance) GetFundsMessages() string {
	fund, err := binance.checkFund()
	if err != nil {
		text := fmt.Sprintf("Err GetFundsMessages %s", err)
		log.Println(text)
		return text
	}
	text1 := "Binance Funds:  "
	text2 := "\n Inorder: "
	for _, balance := range fund.Balances {
		asset := balance.Asset
		freeFund := balance.GetFree()
		lockedFund := balance.GetLockedFund()
		if freeFund > 0 || lockedFund > 0 {
			text1 = fmt.Sprintf("%s %v %v - ", text1, freeFund, asset)
			text2 = fmt.Sprintf("%s %v %v - ", text2, lockedFund, asset)
		}
	}
	return text1 + text2
}

type UserAsset struct {
	Asset    string `json:"asset"`
	NetAsset string `json:"netAsset"`
}

func (ua UserAsset) GetNetAsset() float64 {
	res, _ := strconv.ParseFloat(ua.NetAsset, 64)
	return res
}

type MarginDetails struct {
	UserAssets []UserAsset `json:"userAssets"`
}

func (binance Binance) GetMarginDetails() (*MarginDetails, error) {
	params := make(map[string]string)
	body, code, err := binance.makeRequest("GET", params, "/sapi/v1/margin/account")
	if err != nil {
		log.Printf("Err GetMarginDetails, statusCode: %d err: %s", code, err.Error())
		return nil, err
	}
	var marginDetails MarginDetails
	err = json.Unmarshal([]byte(body), &marginDetails)
	if err != nil {
		panic(err)
	}
	return &marginDetails, nil
}
