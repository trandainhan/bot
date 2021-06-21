package fiahub

import (
	"log"
	"time"

	"github.com/go-pg/pg/v10"
)

type Transaction struct {
	ID    int
	State string
}

func (fia Fiahub) GetSelfMatchingTransaction(user_id int, order_id int) (*Transaction, error) {
	tx := &Transaction{}
	now := time.Now().UTC().AddDate(0, 0, -3)
	created_at := now.Format("2006-01-02 15:04:05")
	err := fia.DB.Model(tx).
		Where("user_id = ? AND orderrer_id = ?", user_id, user_id).
		Column("id", "state").
		WhereGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.WhereOr("source_type = 'Order' and source_id = ? and type in ('BuyTransaction', 'SellTransaction') and created_at > ?", order_id, created_at).
				WhereOr("origin_type = 'Order' and origin_id = ? and type in ('BuyTransaction', 'SellTransaction') and created_at > ?", order_id, created_at)
			return q, nil
		}).Limit(1).Select()
	if err != pg.ErrNoRows {
		log.Printf("Err GetSelfMatchingTransaction: %s", err.Error())
	}
	if err != nil {
		return nil, err
	}
	return tx, nil
}
