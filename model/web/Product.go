package web

import (
	"time"

	"github.com/RuhullahReza/SecondHand/model/entity"
	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required" conform:"trim"`
	Price       int64  `json:"price" binding:"required,number,gte=1000" conform:"trim"`
	Category    string `json:"category" binding:"required" conform:"name"`
	Description string `json:"description" conform:"trim,!html,!js"`
}

type ProductResponse struct {
	Id   		uuid.UUID 	`db:"id" json:"id"`
	Name 		string    	`db:"name" json:"name"`
	Price       int64  		`db:"price" json:"price"`
	Category    string 		`db:"category" json:"category"`
	Thumbnail 	string    	`db:"thumbnail" json:"thumbnail"`
}

type ProductDetailResponse struct {
	OwnerId   	uuid.UUID 	`db:"owner_id" json:"owner_id"`	
	Owner 		string    	`db:"owner" json:"owner"`
	City 		string		`db:"city" json:"city"`
	ImageUrl 	string		`db:"image_url" json:"image_url"`
	Id   		uuid.UUID 	`db:"id" json:"id"`
	Name 		string    	`db:"name" json:"name"`
	Price       int64  		`db:"price" json:"price"`
	Category    string 		`db:"category" json:"category"`
	Description string    	`db:"description" json:"description"`
	Sold		bool      `db:"sold" json:"sold"`
	Published	bool      `db:"published" json:"published"`
	UpdatedAt   time.Time 	`db:"updated_at" json:"updated_at"`
	ProductImages []entity.ProductImage `json:"product_image"`
}

type UpdateProductRequest struct {
	Id   		uuid.UUID 	`json:"id"`
	AccountId   uuid.UUID 	`json:"account_id"`
	Name        string 		`json:"name" binding:"required" conform:"trim"`
	Price       int64  		`json:"price" binding:"required,number,gte=1000" conform:"trim"`
	Category    string 		`json:"category" binding:"required" conform:"name"`
	Description string 		`json:"description" conform:"trim,!html,!js"`
}

type ProductImageRequest struct {
	ImageId   	uuid.UUID 	`json:"image_id"`
	ProductId   uuid.UUID 	`json:"product_id"`
}