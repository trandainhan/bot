package exchanges

import "os"

func (ex ExchangeClient) GetFundsMessages() string {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		return ex.Ftx.GetFundsMessages()
	}
	return ex.Bn.GetFundsMessages()
}
