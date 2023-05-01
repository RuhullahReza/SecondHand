package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id		    uuid.UUID `db:"id" json:"id"`
	Email    	string    `db:"email" json:"email"`
	Password 	string    `db:"password" json:"password"`
	Role 		string    `db:"role" json:"role"`
	CreatedAt 	time.Time `db:"created_at" json:"created_at"`
	UpdatedAt 	time.Time `db:"updated_at" json:"updated_at"`
}
