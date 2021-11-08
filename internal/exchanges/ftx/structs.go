package ftx

import (
	"time"
)

type FtxClient struct{}

type NewOrder struct {
	Market     string  `json:"market"`
	Side       string  `json:"side"`
	Price      float64 `json:"price"`
	Type       string  `json:"type"`
	Size       float64 `json:"size"`
	ReduceOnly bool    `json:"reduceOnly"`
	Ioc        bool    `json:"ioc"`
	PostOnly   bool    `json:"postOnly"`
	// ClientID                string  `json:"clientId"`
	// ExternalReferralProgram string  `json:"externalReferralProgram"`
}

type Order struct {
	CreatedAt     time.Time `json:"createdAt"`
	FilledSize    float64   `json:"filledSize"`
	Future        string    `json:"future"`
	ID            int64     `json:"id"`
	Market        string    `json:"market"`
	Price         float64   `json:"price"`
	AvgFillPrice  float64   `json:"avgFillPrice"`
	RemainingSize float64   `json:"remainingSize"`
	Side          string    `json:"side"`
	Size          float64   `json:"size"`
	Status        string    `json:"status"`
	Type          string    `json:"type"`
	ReduceOnly    bool      `json:"reduceOnly"`
	Ioc           bool      `json:"ioc"`
	PostOnly      bool      `json:"postOnly"`
	ClientID      string    `json:"clientId"`
}

type NewOrderResponse struct {
	Success bool  `json:"success"`
	Result  Order `json:"result"`
}

type CancelOrderResponse struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
}

type OpenOrderResponse struct {
	Success bool    `json:"success"`
	Result  []Order `json:"result"`
}

type Balance struct {
	Coin                   string  `json:"coin"`
	Free                   float64 `json:"free"`
	SpotBorrow             float64 `json:"spotBorrow"`
	Total                  float64 `json:"total"`
	UsdValue               float64 `json:"usdValue"`
	AvailableWithoutBorrow float64 `json:"availableWithoutBorrow"`
}

type WalletResponse struct {
	Success bool      `json:"success"`
	Result  []Balance `json:"result"`
}
