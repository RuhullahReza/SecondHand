package entity

import (
	"time"
	"github.com/google/uuid"
)

type City struct {
	Id		    uuid.UUID 	`db:"id" json:"id"`
	Name    	string   	`db:"name" json:"name"`
	CreatedAt 	time.Time 	`db:"created_at" json:"created_at"`
}

type Category struct {
	Id		    uuid.UUID 	`db:"id" json:"id"`
	Name    	string   	`db:"name" json:"name"`
	CreatedAt 	time.Time 	`db:"created_at" json:"created_at"`
}