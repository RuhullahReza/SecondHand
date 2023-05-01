package entity

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	Id        uuid.UUID `db:"id" json:"id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	Url       string    `db:"image_url" json:"image_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ProductImage struct {
	Id        uuid.UUID `db:"id" json:"id"`
	Url       string    `db:"image_url" json:"image_url"`
}