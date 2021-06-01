package binance

import (
	"encoding/json"
	"fmt"
)

type Balance struct {
	Asset      string  `json:"asset"`
	Free       float64 `json:"fee"`
	LockedFund float64 `json:"locked"`
}

type Fund struct {
	Balances []Balance `json:"balances"`
}

func (binance Binance) CheckFund(name string) float64 {
	fund := binance.checkFund()
	var result float64
	for _, balance := range fund.Balances {
		if balance.Asset == name {
			result = balance.Free
		}
	}
	return result
}

func (binance Binance) checkFund() *Fund {
	params := make(map[string]string)
	body, _, err := binance.makeRequest("GET", params, "/api/v3/account")
	var fund *Fund
	err = json.Unmarshal([]byte(body), fund)
	if err != nil {
		return nil
	}
	return fund
}

func (binance Binance) GetFundsMessages() string {
	fund := binance.checkFund()
	text1 := "Binance Funds:  "
	text2 := "\n Inorder: "
	for _, balance := range fund.Balances {
		asset := balance.Asset
		freeFund := balance.Free
		lockedFund := balance.LockedFund
		if freeFund > 0 || lockedFund > 0 {
			text1 = fmt.Sprintf("%s %v %v - ", text1, freeFund, asset)
			text2 = fmt.Sprintf("%s %v %v - ", text1, lockedFund, asset)
		}
	}
	return text1 + text2
}

type UserAsset struct {
	Name     string
	NetAsset float64
}

type MarginDetails struct {
	UserAssets []UserAsset `json:"userAssets"`
}

func (binance Binance) GetMarginDetails() (*MarginDetails, error) {
	params := make(map[string]string)
	body, _, err := binance.makeRequest("GET", params, "/sapi/v1/margin/account")
	if err != nil {

	}
	var marginDetails *MarginDetails
	err = json.Unmarshal([]byte(body), marginDetails)
	if err != nil {
		return nil, err
	}
	return marginDetails, nil
}
