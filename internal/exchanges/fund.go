package exchanges

import "os"

func (ex ExchangeClient) GetFundsMessages() string {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		return ex.Ftx.GetFundsMessages()
	}
	return ex.Bn.GetFundsMessages()
}

func (ex ExchangeClient) CheckFund(name string) (float64, error) {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		return ex.Ftx.CheckFund(name)
	}
	return ex.Bn.CheckFund(name)
}
