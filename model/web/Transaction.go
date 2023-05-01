package web

import (
	"time"

	"github.com/google/uuid"
)

type TransactionRequest struct {
	BuyerId 	uuid.UUID 	`json:"buyer_id"`
	ProductId 	uuid.UUID	`json:"product_id" binding:"required"`
	Price 		int64		`json:"price" binding:"required"`
}

type TransactionUpdateRequest struct {
	BuyerId 		uuid.UUID 	`json:"buyer_id"`
	TransactionId 	uuid.UUID	`json:"transaction_id" binding:"required"`
	Price 			int64		`json:"price" binding:"required"`
}

type TransactionDetailResponse struct {
	ProductId   	uuid.UUID 	`db:"product_id" json:"product_id"`	
	ProductName 	string		`db:"product_name" json:"product_name"`	
	ProductPrice 	int64		`db:"product_price" json:"product_price"`
	ProductImage 	string		`db:"product_image" json:"product_image"`	
	BuyerId   		uuid.UUID 	`db:"buyer_id" json:"buyer_id"`		
	BuyerName 		string		`db:"buyer_name" json:"buyer_name"`	
	BuyerCity 		string		`db:"buyer_city" json:"buyer_city"`	
	BuyerImage 		string		`db:"buyer_image" json:"buyer_image"`	
	SellerId   		uuid.UUID 	`db:"seller_id" json:"seller_id"`		
	SellerName 		string		`db:"seller_name" json:"seller_name"`	
	SellerCity 		string		`db:"seller_city" json:"seller_city"`	
	SellerImage 	string		`db:"seller_image" json:"seller_image"`		
	Accepted	 	bool		`db:"accepted" json:"accepted"`
	PriceOffer	 	int64		`db:"price_offer" json:"price_offer"`
	UpdatedAt   	time.Time 	`db:"updated_at" json:"updated_at"`
}

type Offer struct {
	BuyerId   		uuid.UUID 	`db:"buyer_id" json:"buyer_id"`		
	BuyerName 		string		`db:"buyer_name" json:"buyer_name"`	
	BuyerCity 		string		`db:"buyer_city" json:"buyer_city"`	
	BuyerImage 		string		`db:"buyer_image" json:"buyer_image"`		
	Accepted	 	bool		`db:"accepted" json:"accepted"`
	PriceOffer	 	int64		`db:"price_offer" json:"price_offer"`
	UpdatedAt   	time.Time 	`db:"updated_at" json:"updated_at"`
}

type OfferByProduct struct {
	OwnerId			uuid.UUID 	`db:"owner_id" json:"owner_id"`	
	ProductId   	uuid.UUID 	`db:"product_id" json:"product_id"`	
	ProductName 	string		`db:"product_name" json:"product_name"`	
	ProductPrice 	int64		`db:"product_price" json:"product_price"`
	ProductImage 	string		`db:"product_image" json:"product_image"`	
	Offer			[]Offer		`json:"offer"`
}

type OfferWithAccount struct {	
	Id   			uuid.UUID 	`db:"id" json:"id"`	
	Name 			string		`db:"name" json:"name"`	
	City 			string		`db:"city" json:"city"`	
	Image 			string		`db:"image" json:"image"`	
	ProductId   	uuid.UUID 	`db:"product_id" json:"product_id"`	
	ProductName 	string		`db:"product_name" json:"product_name"`	
	ProductPrice 	int64		`db:"product_price" json:"product_price"`
	ProductImage 	string		`db:"product_image" json:"product_image"`	
	Accepted	 	bool		`db:"accepted" json:"accepted"`
	PriceOffer	 	int64		`db:"price_offer" json:"price_offer"`
	UpdatedAt   	time.Time 	`db:"updated_at" json:"updated_at"`
}

type OfferWithProduct struct {
	ProductId   	uuid.UUID 	`db:"product_id" json:"product_id"`	
	ProductName 	string		`db:"product_name" json:"product_name"`	
	ProductPrice 	int64		`db:"product_price" json:"product_price"`
	ProductImage 	string		`db:"product_image" json:"product_image"`		
	Accepted	 	bool		`db:"accepted" json:"accepted"`
	PriceOffer	 	int64		`db:"price_offer" json:"price_offer"`
	UpdatedAt   	time.Time 	`db:"updated_at" json:"updated_at"`
}

type OfferByBuyer struct {
	BuyerId   	uuid.UUID 	`db:"buyer_id" json:"buyer_id"`	
	BuyerName 	string		`db:"buyer_name" json:"buyer_name"`	
	BuyerCity 	string		`db:"buyer_city" json:"buyer_city"`	
	BuyerImage 	string		`db:"buyer_image" json:"buyer_image"`	
	Offer		[]OfferWithProduct		`json:"offer"`
}


