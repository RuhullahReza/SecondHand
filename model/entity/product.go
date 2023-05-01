package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID `db:"id" json:"id"`
	AccountId   uuid.UUID `db:"account_id" json:"account_id"`
	Name        string    `db:"name" json:"name"`
	Price       int64     `db:"price" json:"price"`
	Category    string    `db:"category" json:"category"`
	Description string    `db:"description" json:"description"`
	Thumbnail 	string    `db:"thumbnail" json:"thumbnail"`
	Sold		bool      `db:"sold" json:"sold"`
	Published	bool      `db:"published" json:"published"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Deleted		bool      `db:"deleted" json:"deleted"`
}