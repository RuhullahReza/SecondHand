package entity

import (
	"time"
	"github.com/google/uuid"
)

type Profile struct {
	Id		    uuid.UUID 	`db:"id" json:"id"`
	Name    	string   	`db:"name" json:"name"`
	City 		string   	`db:"city" json:"city"`
	Address 	string   	`db:"address" json:"address"`
	PhoneNumber string   	`db:"phone_number" json:"phone_number"`
	ImageUrl 	string   	`db:"image_url" json:"image_url"`
	CreatedAt 	time.Time 	`db:"created_at" json:"created_at"`
	UpdatedAt 	time.Time 	`db:"updated_at" json:"updated_at"`
}