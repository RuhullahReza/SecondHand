package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id          uuid.UUID `db:"id" json:"id"`
	SellerId   	uuid.UUID `db:"seller_id" json:"seller_id"`
	BuyerId   	uuid.UUID `db:"buyer_id" json:"buyer_id"`
	ProductId   uuid.UUID `db:"product_id" json:"product_id"`
	PriceOffer  int64     `db:"price_offer" json:"price_offer"`
	Accepted   	bool      `db:"accepted" json:"accepted"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Deleted     bool      `db:"deleted" json:"deleted"`
}