package repository

import (
	"context"
	"log"
	"time"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/entity"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetTransactionById(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetTransactionDetail(ctx context.Context, id uuid.UUID) ([]web.TransactionDetailResponse, error)
	GetOfferByProduct(ctx context.Context, id uuid.UUID) ([]web.Offer, error)
	GetOfferByBuyer(ctx context.Context, buyer_id uuid.UUID, seller_id uuid.UUID) ([]web.OfferWithProduct, error)
	GetMyTransaction(ctx context.Context, buyer_id uuid.UUID) ([]web.OfferWithAccount, error)
	GetOfferByAccount(ctx context.Context, seller_id uuid.UUID) ([]web.OfferWithAccount, error)
	UpdatePrice(ctx context.Context, price int64, id uuid.UUID, buyer_id uuid.UUID) error
	CheckStatus(ctx context.Context, transactionId uuid.UUID) (bool, error)
	SetStatus(ctx context.Context, status bool, id uuid.UUID, sellerId uuid.UUID) error
	DeleteOne(ctx context.Context, id uuid.UUID) error
	DeleteOnSold(ctx context.Context, product_id uuid.UUID) error
}

type TransactionRepositoryImpl struct {
	DB *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &TransactionRepositoryImpl{
		DB: db,
	}
}

func (r *TransactionRepositoryImpl) Create(ctx context.Context, transaction *entity.Transaction) error {

	query := `
	INSERT INTO 
		transactions
		(buyer_id, seller_id, product_id, price_offer) 
	VALUES 
		($1, $2, $3, $4)
	`
	_, err := r.DB.ExecContext(ctx, query, transaction.BuyerId, transaction.SellerId, transaction.ProductId, transaction.PriceOffer)

	if err != nil {
		log.Printf("failed to query create transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

func (r *TransactionRepositoryImpl) GetTransactionById(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	
	transaction := &entity.Transaction{}

	query := `
	SELECT *
	FROM 
		transactions
	WHERE 
		id=$1 AND deleted = FALSE
	LIMIT 1
	`

	if err := r.DB.GetContext(ctx, transaction, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return transaction, helper.NewNotFound("id", id.String())
		}
		
		log.Printf("failed to query get transaction by Id, err : %v\n", err)
		return transaction, helper.NewInternal()
	}

	return transaction, nil
}

func (r *TransactionRepositoryImpl) GetTransactionDetail(ctx context.Context, id uuid.UUID) ([]web.TransactionDetailResponse, error) {

	transactions := []web.TransactionDetailResponse{}

	query := `
	SELECT 
		products.id as product_id, products.name as product_name, products.price as product_price, COALESCE(products.thumbnail,'') as product_image,
		buyer.id as buyer_id, buyer.name as buyer_name, buyer.city as buyer_city, COALESCE(buyer.image_url,'') as buyer_image,
		seller.id as seller_id, seller.name as seller_name, seller.city as seller_city, COALESCE(seller.image_url,'') as seller_image,
		t.price_offer as price_offer, t.accepted as accepted, t.updated_at as updated_at
	FROM
		transactions t
	JOIN 
		products ON t.product_id = products.id
	LEFT jOIN 
		profiles buyer ON t.buyer_id = buyer.id
	LEFT JOIN 
		profiles seller ON t.seller_id= seller.id
	WHERE 
		t.id = $1 AND t.deleted = FALSE
	`
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to query get transaction detail, err : %v\n", err)
		return transactions, helper.NewInternal()
	}

	for rows.Next(){
		transaction := web.TransactionDetailResponse{}
		err := rows.Scan(
			&transaction.ProductId, &transaction.ProductName, &transaction.ProductPrice, &transaction.ProductImage,
			&transaction.BuyerId, &transaction.BuyerName, &transaction.BuyerCity, &transaction.BuyerImage,
			&transaction.SellerId, &transaction.SellerName, &transaction.SellerCity, &transaction.SellerImage,
			&transaction.PriceOffer, &transaction.Accepted, &transaction.UpdatedAt,
		)
		if err != nil {
			log.Printf("failed to scanning on transaction detail, err : %v\n", err)
			return transactions, helper.NewInternal()
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepositoryImpl) GetOfferByProduct(ctx context.Context, id uuid.UUID) ([]web.Offer, error) {

	offers := []web.Offer{}

	query := `
	SELECT 
		buyer.id as buyer_id, buyer.name as buyer_name, buyer.city as buyer_city, COALESCE(buyer.image_url,'') as buyer_image,
		t.price_offer as price_offer, t.accepted as accepted, t.created_at as created_at
	FROM
		transactions t
	JOIN 
		products ON t.product_id = products.id
	LEFT jOIN 
		profiles buyer ON t.buyer_id = buyer.id
	WHERE 
		t.product_id = $1 AND t.deleted = FALSE
	`
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to query get offer by product, err : %v\n", err)
		return offers, helper.NewInternal()
	}

	for rows.Next(){
		offer := web.Offer{}
		err := rows.Scan(
			&offer.BuyerId, &offer.BuyerName, &offer.BuyerCity, &offer.BuyerImage,
			&offer.PriceOffer, &offer.Accepted, &offer.UpdatedAt,
		)
		if err != nil {
			log.Printf("failed to scanning on get offer by product, err : %v\n", err)
			return offers, helper.NewInternal()
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func (r *TransactionRepositoryImpl) GetOfferByBuyer(ctx context.Context, buyer_id uuid.UUID, seller_id uuid.UUID) ([]web.OfferWithProduct, error) {

	offers := []web.OfferWithProduct{}

	query := `
	SELECT 
		p.id as product_id, p.name as product_name, p.price as product_price, COALESCE(p.thumbnail,'') as product_image,
		t.accepted, t.price_offer, t.updated_at
	FROM
		transactions t
	JOIN 
		products p ON t.product_id = p.id
	WHERE
		t.buyer_id = $1 AND t.seller_id = $2 AND t.deleted = FALSE
	`
	rows, err := r.DB.QueryContext(ctx, query, buyer_id, seller_id)
	if err != nil {
		log.Printf("failed to query get offer by buyer, err : %v\n", err)
		return offers, helper.NewInternal()
	}

	for rows.Next(){
		offer := web.OfferWithProduct{}
		err := rows.Scan(&offer.ProductId, &offer.ProductName, &offer.ProductPrice, &offer.ProductImage, &offer.Accepted, &offer.PriceOffer, &offer.UpdatedAt)
		if err != nil {
			log.Printf("failed to scanning on get offer by buyer, err : %v\n", err)
			return offers, helper.NewInternal()
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func (r *TransactionRepositoryImpl) GetOfferByAccount(ctx context.Context, seller_id uuid.UUID) ([]web.OfferWithAccount, error) {

	offers := []web.OfferWithAccount{}

	query := `
	SELECT 
		u.id as id, u.name as name, u.city as city, COALESCE(u.image_url,'') as image,
		p.id as product_id, p.name as product_name, p.price as product_price, COALESCE(p.thumbnail,'') as product_image,
		t.accepted as accepted, t.price_offer as price_offer, t.updated_at as updated_at	
	FROM
		transactions t
	JOIN 
		products p ON t.product_id = p.id
	JOIN
		profiles u ON t.buyer_id = u.id
	WHERE
		t.seller_id = $1 AND t.deleted = FALSE
	`
	rows, err := r.DB.QueryContext(ctx, query, seller_id)
	if err != nil {
		log.Printf("failed to query get offer by account, err : %v\n", err)
		return offers, helper.NewInternal()
	}

	for rows.Next(){
		offer := web.OfferWithAccount{}
		err := rows.Scan(
			&offer.Id, &offer.Name, &offer.City, &offer.Image, 
			&offer.ProductId, &offer.ProductName, &offer.ProductPrice, &offer.ProductImage,
			&offer.Accepted, &offer.PriceOffer, &offer.UpdatedAt,
		)
		if err != nil {
			log.Printf("failed to scanning on get offer by account, err : %v\n", err)
			return offers, helper.NewInternal()
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func (r *TransactionRepositoryImpl) GetMyTransaction(ctx context.Context, buyer_id uuid.UUID) ([]web.OfferWithAccount, error) {

	offers := []web.OfferWithAccount{}

	query := `
	SELECT 
		u.id as id, u.name as name, u.city as city, coalesce(u.image_url,'') as image,
		p.id as product_id, p.name as product_name, p.price as product_price, COALESCE(p.thumbnail,'') as product_image,
		t.accepted as accepted, t.price_offer as price_offer, t.updated_at as updated_at	
	FROM
		transactions t
	JOIN 
		products p ON t.product_id = p.id
	JOIN
		profiles u ON t.buyer_id = u.id
	WHERE
		t.buyer_id = $1 AND t.deleted = FALSE
	`
	rows, err := r.DB.QueryContext(ctx, query, buyer_id)
	if err != nil {
		log.Printf("failed to query get my transaction, err : %v\n", err)
		return offers, helper.NewInternal()
	}

	for rows.Next(){
		offer := web.OfferWithAccount{}
		err := rows.Scan(
			&offer.Id, &offer.Name, &offer.City, &offer.Image, 
			&offer.ProductId, &offer.ProductName, &offer.ProductPrice, &offer.ProductImage,
			&offer.Accepted, &offer.PriceOffer, &offer.UpdatedAt,
		)
		if err != nil {
			log.Printf("failed to scanning on get my transaction, err : %v\n", err)
			return offers, helper.NewInternal()
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func (r *TransactionRepositoryImpl) UpdatePrice(ctx context.Context, price int64, id uuid.UUID, buyer_id uuid.UUID) error {
	
	query := `
	UPDATE 
		transactions 
	SET 
		price_offer = $1, updated_at = $2
	WHERE 
		id = $3 AND buyer_id = $4 AND accepted = FALSE AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, price, time.Now(), id, buyer_id)

	if err != nil {
		log.Printf("failed to query update transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when update transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", id.String())
	}

	return nil
}

func (r *TransactionRepositoryImpl) CheckStatus(ctx context.Context, transactionId uuid.UUID) (bool, error) {

	transaction := &entity.Transaction{}

	query := `
	SELECT
		accepted
	FROM 
		transactions
	WHERE 
		id=$1 AND deleted = FALSE 
	LIMIT 1
	`

	if err := r.DB.GetContext(ctx, transaction, query, transactionId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, helper.NewNotFound("transaction id", transactionId.String())
		}

		log.Printf("failed to query chek transaction status, err : %v\n", err)
		return false, helper.NewInternal()
	}

	return transaction.Accepted, nil
}

func (r *TransactionRepositoryImpl) SetStatus(ctx context.Context, status bool, id uuid.UUID, sellerId uuid.UUID) error {
	
	query := `
	UPDATE 
		transactions 
	SET 
		accepted = $1
	WHERE 
		id = $2 AND seller_id = $3 AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, status, id, sellerId)

	if err != nil {
		log.Printf("failed to query update transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when update transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", id.String())
	}

	return nil
}

func (r *TransactionRepositoryImpl) DeleteOne(ctx context.Context, id uuid.UUID) error {
	
	query := `
	UPDATE 
		transactions 
	SET 
		deleted = TRUE
	WHERE 
		id = 1  AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, id)

	if err != nil {
		log.Printf("failed to query delete one transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when delete one transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", id.String())
	}

	return nil
}

func (r *TransactionRepositoryImpl) DeleteOnSold(ctx context.Context, product_id uuid.UUID) error {
	
	query := `
	UPDATE 
		transactions 
	SET 
		deleted = TRUE
	WHERE 
		product_id = 1 AND accepted=FALSE AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, product_id)

	if err != nil {
		log.Printf("failed to query delete on sold transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when delete on sold transaction, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", product_id.String())
	}

	return nil
}