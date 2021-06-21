package fiahub

import "log"

type Order struct {
	ID         int
	State      string
	CoinAmount float64
	UserID     int
}

func (fia Fiahub) GetOrderDetails(orderID int) (*Order, error) {
	order := &Order{ID: orderID}
	err := fia.DB.Model(order).WherePK().Select()
	if err != nil {
		log.Printf("Err GetOrderDetails %s", err.Error())
		return nil, err
	}
	return order, nil
}
