package web

import (
	"github.com/google/uuid"
	_ "github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	Name     string `db:"name" json:"name" binding:"required" conform:"name"`
	Email    string `db:"email" json:"email" binding:"required,email" conform:"lower,trim,email"` 
	Password string `db:"password" json:"password" binding:"required,gte=5,lte=255"`
}

type LoginRequest struct {
	Email    string `db:"email" json:"email" binding:"required" conform:"lower,trim,email"`
	Password string `db:"password" json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdateProfileRequest struct {
	Id          uuid.UUID 	`db:"id" json:"id"`
	Name        string		`db:"name" json:"name" binding:"required" conform:"name"`
	City    	string   	`db:"city" json:"city" conform:"name"`
	Address     string   	`db:"address" json:"address" conform:"trim"`
	PhoneNumber string   	`db:"phone_number" json:"phone_number"`
}


type GetProfileResponse struct {
	Id          uuid.UUID 	`db:"id" json:"id"`
	Name        string		`db:"name" json:"name" binding:"required"`
	City    	string   	`db:"city" json:"city"`
	Address     string   	`db:"address" json:"address"`
	PhoneNumber string   	`db:"phone_number" json:"phone_number"`
	ImageUrl 	string  	`db:"image_url" json:"image_url"`
}

type GetAccountRequest struct {
	ID string `uri:"id"`
}