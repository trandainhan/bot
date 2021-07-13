package exchanges

import (
	"os"
)

func (ex ExchangeClient) BuyLimit(coin string, price float64, quantity float64) (*OrderResp, error) {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		ftxOrderResp, err := ex.Ftx.BuyLimit(coin+"/USDT", price, quantity)
		if err != nil {
			return nil, err
		}
		order := OrderResp{
			ID:          ftxOrderResp.ID,
			ClientID:    ftxOrderResp.ClientID,
			OriginQty:   ftxOrderResp.Size,
			ExecutedQty: ftxOrderResp.FilledSize,
			Price:       ftxOrderResp.Price,
			Status:      ftxOrderResp.Status,
			Side:        ftxOrderResp.Side,
		}
		return &order, nil
	}
	binanceOrderDetails, err := ex.Bn.BuyLimit(coin+"USDT", price, quantity)
	if err != nil {
		return nil, err
	}
	order := OrderResp{
		ID:          binanceOrderDetails.OrderID,
		ClientID:    binanceOrderDetails.ClientOrderID,
		OriginQty:   binanceOrderDetails.GetOriginQty(),
		ExecutedQty: binanceOrderDetails.GetExecutedQty(),
		Price:       binanceOrderDetails.GetPrice(),
		Status:      binanceOrderDetails.Status,
		Side:        binanceOrderDetails.Side,
	}
	return &order, nil
}

func (ex ExchangeClient) SellLimit(coin string, price float64, quantity float64) (*OrderResp, error) {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		ftxOrderResp, err := ex.Ftx.SellLimit(coin+"/USDT", price, quantity)
		if err != nil {
			return nil, err
		}
		order := OrderResp{
			ID:          ftxOrderResp.ID,
			ClientID:    ftxOrderResp.ClientID,
			OriginQty:   ftxOrderResp.Size,
			ExecutedQty: ftxOrderResp.FilledSize,
			Price:       ftxOrderResp.Price,
			Status:      ftxOrderResp.Status,
			Side:        ftxOrderResp.Side,
		}
		return &order, nil
	}
	binanceOrderDetails, err := ex.Bn.SellLimit(coin+"USDT", price, quantity)
	if err != nil {
		return nil, err
	}
	order := OrderResp{
		ID:          binanceOrderDetails.OrderID,
		ClientID:    binanceOrderDetails.ClientOrderID,
		OriginQty:   binanceOrderDetails.GetOriginQty(),
		ExecutedQty: binanceOrderDetails.GetExecutedQty(),
		Price:       binanceOrderDetails.GetPrice(),
		Status:      binanceOrderDetails.Status,
		Side:        binanceOrderDetails.Side,
	}
	return &order, nil
}

func (ex ExchangeClient) GetOrder(coin string, orderID int64, clientID string) (*OrderResp, error) {
	var order OrderResp
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		ftxOrderResp, err := ex.Ftx.GetOrder(coin+"/USDT", orderID)
		if err != nil {
			return nil, err
		}
		order = OrderResp{
			ID:          ftxOrderResp.ID,
			ClientID:    ftxOrderResp.ClientID,
			OriginQty:   ftxOrderResp.Size,
			ExecutedQty: ftxOrderResp.FilledSize,
			Price:       ftxOrderResp.Price,
			Status:      ftxOrderResp.Status,
			Side:        ftxOrderResp.Side,
		}
		return &order, nil
	}
	binanceOrderDetails, err := ex.Bn.GetOrder(coin+"USDT", orderID, clientID)
	if err != nil {
		return nil, err
	}
	order = OrderResp{
		ID:          binanceOrderDetails.OrderID,
		ClientID:    binanceOrderDetails.ClientOrderID,
		OriginQty:   binanceOrderDetails.GetOriginQty(),
		ExecutedQty: binanceOrderDetails.GetExecutedQty(),
		Price:       binanceOrderDetails.GetPrice(),
		Status:      binanceOrderDetails.Status,
		Side:        binanceOrderDetails.Side,
	}
	return &order, nil
}
