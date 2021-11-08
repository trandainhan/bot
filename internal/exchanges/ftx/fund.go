package ftx

import (
	"errors"
	"fmt"
	"log"

	"gitlab.com/fiahub/bot/internal/utils"
)

func (ftx FtxClient) CheckFund(name string) (float64, error) {
	balances, err := ftx.checkFund()
	if err != nil {
		return -1.0, err
	}
	var result float64
	for _, balance := range balances {
		if balance.Coin == name {
			result = balance.Total
		}
	}
	return result, nil
}

func (ftx FtxClient) checkFund() ([]Balance, error) {
	body, code, err := ftx.makeRequest("GET", "/wallet/balances", "")
	if err != nil {
		log.Printf("Err checkFund, statusCode: %d err: %s", code, err.Error())
		return nil, err
	}
	if code >= 400 {
		text := fmt.Sprintf("Err checkFund, statusCode: %d err: %s", code, body)
		return nil, errors.New(text)
	}
	var resp WalletResponse
	err = utils.ProcessResponse(body, &resp)
	if err != nil {
		panic(err)
	}
	return resp.Result, nil
}

func (ftx FtxClient) GetFundsMessages() string {
	balances, err := ftx.checkFund()
	if err != nil {
		text := fmt.Sprintf("Err GetFundsMessages %s", err)
		log.Println(text)
		return text
	}
	text1 := "FTX Funds:  "
	text2 := "\n Inorder: "
	for _, balance := range balances {
		asset := balance.Coin
		freeFund := balance.Free
		lockedFund := balance.Total - balance.Free
		if freeFund > 0 || lockedFund > 0 {
			text1 = fmt.Sprintf("%s %.6f %s - ", text1, freeFund, asset)
			text2 = fmt.Sprintf("%s %.6f %s - ", text2, lockedFund, asset)
		}
	}
	return text1 + text2
}
