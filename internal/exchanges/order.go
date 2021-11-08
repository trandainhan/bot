package exchanges

import (
	"os"
)

func (ex ExchangeClient) BuyLimit(coin, fiat string, price float64, quantity float64) (*OrderResp, error) {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		ftxOrderResp, err := ex.Ftx.BuyLimit(coin+"/"+fiat, price, quantity)
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
	binanceOrderDetails, err := ex.Bn.BuyLimit(coin+fiat, price, quantity)
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

func (ex ExchangeClient) SellLimit(coin, fiat string, price float64, quantity float64) (*OrderResp, error) {
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		ftxOrderResp, err := ex.Ftx.SellLimit(coin+"/"+fiat, price, quantity)
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
	binanceOrderDetails, err := ex.Bn.SellLimit(coin+fiat, price, quantity)
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

func (ex ExchangeClient) GetOrder(coin, fiat string, orderID int64, clientID string) (*OrderResp, error) {
	var order OrderResp
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		ftxOrderResp, err := ex.Ftx.GetOrder(coin+"/"+fiat, orderID)
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
	binanceOrderDetails, err := ex.Bn.GetOrder(coin+fiat, orderID, clientID)
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

func (ex ExchangeClient) CancelOrder(coin, fiat string, orderID int64, clientID string) (*OrderResp, error) {
	var order OrderResp
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		_, err := ex.Ftx.CancelOrder(coin+"/"+fiat, orderID)
		if err != nil {
			return nil, err
		}
		return nil, nil // ftx doesn't return order details
	}
	binanceOrderDetails, err := ex.Bn.CancelOrder(coin+fiat, orderID, clientID)
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

func (ex ExchangeClient) CancelAllOrder(coin string) ([]OrderResp, error) {
	var result []OrderResp
	orders, err := ex.Bn.CancelAllOrder(coin + "USDT")
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		temp := OrderResp{
			ID:          order.OrderID,
			ClientID:    order.ClientOrderID,
			OriginQty:   order.GetOriginQty(),
			ExecutedQty: order.GetExecutedQty(),
			Price:       order.GetPrice(),
			Status:      order.Status,
			Side:        order.Side,
		}
		result = append(result, temp)
	}
	return result, nil
}

func (ex ExchangeClient) GetAllOpenOrder(coin, fiat string) ([]OrderResp, error) {
	var result []OrderResp
	exchangeClient := os.Getenv("EXCHANGE_CLIENT")
	if exchangeClient == "FTX" {
		orders, err := ex.Ftx.GetAllOpenOrder(coin + "/" + fiat)
		if err != nil {
			return nil, err
		}
		for _, order := range orders {
			temp := OrderResp{
				ID:          order.ID,
				ClientID:    order.ClientID,
				OriginQty:   order.Size,
				ExecutedQty: order.FilledSize,
				Price:       order.Price,
				Status:      order.Status,
				Side:        order.Side,
			}
			result = append(result, temp)
		}
		return result, nil
	}
	orders, err := ex.Bn.GetAllOpenOrder(coin + fiat)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		temp := OrderResp{
			ID:          order.OrderID,
			ClientID:    order.ClientOrderID,
			OriginQty:   order.GetOriginQty(),
			ExecutedQty: order.GetExecutedQty(),
			Price:       order.GetPrice(),
			Status:      order.Status,
			Side:        order.Side,
		}
		result = append(result, temp)
	}
	return result, nil
}
