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

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetAll(ctx context.Context) ([]web.ProductResponse, error)
	GetByCategory(ctx context.Context, category string) ([]web.ProductResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) ([]web.OfferByProduct, error)
	GetOne(ctx context.Context, id uuid.UUID) ([]web.ProductDetailResponse, error)
	OwnerGetOne(ctx context.Context, id uuid.UUID) ([]web.ProductDetailResponse, error)
	GetByAccount(ctx context.Context, account_id uuid.UUID, status bool, published bool) ([]web.ProductResponse, error)
	Update(ctx context.Context, product *entity.Product) error
	GetOwnerId(ctx context.Context, productId uuid.UUID) (uuid.UUID, error)
	IsOwner(ctx context.Context, accountId uuid.UUID, productId uuid.UUID) (bool, error)
	CheckOwner(ctx context.Context, accountId uuid.UUID, productId uuid.UUID) error
	CheckThumbnail(ctx context.Context, productId uuid.UUID) (bool, error)
	IsThumbnail(ctx context.Context, productId uuid.UUID, path string) (bool, error)
	SetThumbnail(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID, status bool) error
	CheckPublished(ctx context.Context, productId uuid.UUID) (bool, error)
	CheckSold(ctx context.Context, productId uuid.UUID) (bool, error)
	SetSold(ctx context.Context, id uuid.UUID, status bool) error
}

type ProductRepositoryImpl struct {
	DB *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &ProductRepositoryImpl{
		DB: db,
	}
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {

	query := `
	INSERT INTO 
		products 
		(account_id, name, price, category, description) 
	VALUES 
		($1, $2, $3, $4, $5)
	`
	_, err := r.DB.ExecContext(ctx, query, product.AccountId, product.Name, product.Price, product.Category, product.Description)

	if err != nil {
		log.Printf("failed to query create product, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

func (r *ProductRepositoryImpl) CheckOwner(ctx context.Context, accountId uuid.UUID, productId uuid.UUID) error {

	product := &entity.Product{}

	query := "SELECT id FROM products WHERE account_id=$1 AND id=$2 LIMIT 1"

	if err := r.DB.GetContext(ctx, product, query, accountId, productId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return helper.NewNotFound("product id", productId.String())
		}

		log.Printf("failed to query chek owner, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

func (r *ProductRepositoryImpl) IsOwner(ctx context.Context, accountId uuid.UUID, productId uuid.UUID) (bool, error) {

	product := &entity.Product{}

	query := "SELECT id FROM products WHERE account_id=$1 AND id=$2 LIMIT 1"

	if err := r.DB.GetContext(ctx, product, query, accountId, productId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, helper.NewNotFound("product id", productId.String())
		}

		log.Printf("failed to query chek owner, err : %v\n", err)
		return false, helper.NewInternal()
	}

	return true, nil
}

func (r *ProductRepositoryImpl) GetOwnerId(ctx context.Context, productId uuid.UUID) (uuid.UUID, error) {

	product := &entity.Product{}

	query := `
	SELECT 
		account_id 
	FROM 
		products 
	WHERE 
		id=$1 AND sold=FALSE AND published = TRUE AND deleted = FALSE
	LIMIT 1
	`

	if err := r.DB.GetContext(ctx, product, query, productId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return product.Id, helper.NewNotFound("product id", productId.String())
		}

		log.Printf("failed to query chek owner Id, err : %v\n", err)
		return product.Id, helper.NewInternal()
	}

	return product.AccountId, nil
}

func (r *ProductRepositoryImpl) CheckThumbnail(ctx context.Context, productId uuid.UUID) (bool, error) {

	product := &entity.Product{}

	query := `
	SELECT
		COALESCE(thumbnail, '') as thumbnail 
	FROM products 
	WHERE id=$1 LIMIT 1
	`

	if err := r.DB.GetContext(ctx, product, query, productId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, helper.NewNotFound("product id", productId.String())
		}

		log.Printf("failed to query chek thumbnail, err : %v\n", err)
		return false, helper.NewInternal()
	}

	return product.Thumbnail != "", nil
}

func (r *ProductRepositoryImpl) IsThumbnail(ctx context.Context, productId uuid.UUID, path string) (bool, error) {

	product := &entity.Product{}

	query := `
	SELECT
		COALESCE(thumbnail, '') as thumbnail 
	FROM products 
	WHERE id=$1 LIMIT 1
	`

	if err := r.DB.GetContext(ctx, product, query, productId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, helper.NewNotFound("product id", productId.String())
		}

		log.Printf("failed to query chek thumbnail, err : %v\n", err)
		return false, helper.NewInternal()
	}

	return product.Thumbnail == path, nil
}

func (r *ProductRepositoryImpl) SetThumbnail(ctx context.Context, product *entity.Product) error {
	query := `
	UPDATE 
		products 
	SET 
		thumbnail = $1, updated_at = $2
	WHERE 
		id = $3 AND account_id = $4 AND sold = FALSE AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, product.Thumbnail, product.UpdatedAt, product.Id, product.AccountId)

	if err != nil {
		log.Printf("failed to query set thumbnail, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when set thumbnail, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", product.Id.String())
	}

	return nil
}

func (r *ProductRepositoryImpl) GetAll(ctx context.Context) ([]web.ProductResponse, error) {

	products := []web.ProductResponse{}

	query := `
		SELECT 
			id, name, price, category, thumbnail 
		FROM 
			products
		WHERE 
			deleted=FALSE AND sold=FALSE AND published=TRUE
		LIMIT 50
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		log.Printf("failed to query get all product, err : %v\n", err)
		return products, helper.NewInternal()
	}

	for rows.Next() {
		product := web.ProductResponse{}
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Category, &product.Thumbnail)
		if err != nil {
			log.Printf("failed to scanning product, err : %v\n", err)
			return products, helper.NewInternal()
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) GetByCategory(ctx context.Context, category string) ([]web.ProductResponse, error) {

	products := []web.ProductResponse{}

	query := `
		SELECT 
			id, name, price, category, thumbnail 
		FROM 
			products
		WHERE 
			category = $1 AND deleted=FALSE AND sold=FALSE AND published=TRUE
		LIMIT 50
	`
	rows, err := r.DB.QueryContext(ctx, query, category)
	if err != nil {
		log.Printf("failed to query get product by category, err : %v\n", err)
		return products, helper.NewInternal()
	}

	for rows.Next() {
		product := web.ProductResponse{}
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Category, &product.Thumbnail)
		if err != nil {
			log.Printf("failed to scanning product, err : %v\n", err)
			return products, helper.NewInternal()
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) OwnerGetOne(ctx context.Context, id uuid.UUID) ([]web.ProductDetailResponse, error) {

	products := []web.ProductDetailResponse{}

	query := `
		SELECT 
			products.id as id, products.name as name, products.price as price, products.category as category, 
			products.description as description, products.updated_at as updated_at, products.sold as sold, products.published as published,
			profiles.id as owner_id, profiles.name as owner, profiles.city as city, COALESCE(profiles.image_url,'') as image_url
		FROM 
			products
		LEFT JOIN 
			profiles 
		ON 
			profiles.id = products.account_id
		WHERE 
			products.id = $1 AND products.deleted = FALSE
		LIMIT 1
	`
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to query get all product, err : %v\n", err)
		return products, helper.NewInternal()
	}

	for rows.Next() {
		product := web.ProductDetailResponse{}
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Category, &product.Description,
			&product.UpdatedAt, &product.Sold, &product.Published, &product.OwnerId, &product.Owner, &product.City, &product.ImageUrl)
		if err != nil {
			log.Printf("failed to scanning product, err : %v\n", err)
			return products, helper.NewInternal()
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) GetOne(ctx context.Context, id uuid.UUID) ([]web.ProductDetailResponse, error) {

	products := []web.ProductDetailResponse{}

	query := `
		SELECT 
			products.id as id, products.name as name, products.price as price, products.category as category, 
			products.description as description, products.updated_at as updated_at, products.sold as sold, products.published as published,
			profiles.id as owner_id, profiles.name as owner, profiles.city as city, COALESCE(profiles.image_url,'') as image_url
		FROM 
			products
		LEFT JOIN 
			profiles 
		ON 
			profiles.id = products.account_id
		WHERE 
			products.id = $1 AND products.sold = FALSE AND products.published = TRUE AND products.deleted = FALSE
		LIMIT 1
	`
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to query get all product, err : %v\n", err)
		return products, helper.NewInternal()
	}

	for rows.Next() {
		product := web.ProductDetailResponse{}
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Category, &product.Description,
			&product.UpdatedAt, &product.Sold, &product.Published, &product.OwnerId, &product.Owner, &product.City, &product.ImageUrl)
		if err != nil {
			log.Printf("failed to scanning product, err : %v\n", err)
			return products, helper.NewInternal()
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) GetByAccount(ctx context.Context, account_id uuid.UUID, status bool, published bool) ([]web.ProductResponse, error) {

	products := []web.ProductResponse{}

	query := `
		SELECT 
			id, name, price, category, COALESCE(thumbnail,'') as thumbnail 
		FROM 
			products
		WHERE 
			account_id = $1 AND sold=$2 AND published=$3 AND deleted=FALSE
		LIMIT 50
	`
	rows, err := r.DB.QueryContext(ctx, query, account_id, status, published)
	if err != nil {
		log.Printf("failed to query get all product, err : %v\n", err)
		return products, helper.NewInternal()
	}

	for rows.Next() {
		product := web.ProductResponse{}
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Category, &product.Thumbnail)
		if err != nil {
			log.Printf("failed to scanning product, err : %v\n", err)
			return products, helper.NewInternal()
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) GetProduct(ctx context.Context, id uuid.UUID) ([]web.OfferByProduct, error) {

	products := []web.OfferByProduct{}

	query := `
	SELECT 
		account_id as owner_id, id as product_id, name as product_name, price as product_price, COALESCE(thumbnail,'') as product_image
	FROM
		products
	WHERE 
		id = $1 AND deleted = FALSE
	LIMIT 1
	`
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to query get product, err : %v\n", err)
		return products, helper.NewInternal()
	}

	for rows.Next() {
		offer := web.OfferByProduct{}
		err := rows.Scan(&offer.OwnerId, &offer.ProductId, &offer.ProductName, &offer.ProductPrice, &offer.ProductImage)
		if err != nil {
			log.Printf("failed to scanning on get product, err : %v\n", err)
			return products, helper.NewInternal()
		}

		products = append(products, offer)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {

	query := `
	UPDATE 
		products 
	SET 
		name = $1, price = $2, category = $3, description = $4, updated_at = $5
	WHERE 
		id = $6 AND sold = FALSE AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, product.Name, product.Price, product.Category, product.Description, product.UpdatedAt, product.Id)

	if err != nil {
		log.Printf("failed to query update product, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when update product, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", product.Id.String())
	}

	return nil
}

func (r *ProductRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {

	query := `
	UPDATE 
		products 
	SET 
		updated_at = $1, deleted = TRUE
	WHERE 
		id = $2 AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, time.Now(), id)

	if err != nil {
		log.Printf("failed to query delete product, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when delete product, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", id.String())
	}

	return nil
}

func (r *ProductRepositoryImpl) Publish(ctx context.Context, id uuid.UUID, status bool) error {

	query := `
	UPDATE 
		products 
	SET 
		published = $1, updated_at = $2
	WHERE 
		id = $3 AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, status, time.Now(), id)

	if err != nil {
		log.Printf("failed to query publish product, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when publish product, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", id.String())
	}

	return nil
}

func (r *ProductRepositoryImpl) CheckPublished(ctx context.Context, productId uuid.UUID) (bool, error) {

	product := &entity.Product{}

	query := `
	SELECT
		published
	FROM 
		products 
	WHERE 
		id=$1 LIMIT 1
	`

	if err := r.DB.GetContext(ctx, product, query, productId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, helper.NewNotFound("product id", productId.String())
		}

		log.Printf("failed to query chek published status, err : %v\n", err)
		return false, helper.NewInternal()
	}

	return product.Published, nil
}

func (r *ProductRepositoryImpl) CheckSold(ctx context.Context, productId uuid.UUID) (bool, error) {

	product := &entity.Product{}

	query := `
	SELECT
		sold
	FROM 
		products 
	WHERE 
		id=$1 AND deleted = FALSE 
	LIMIT 1
	`

	if err := r.DB.GetContext(ctx, product, query, productId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, helper.NewNotFound("product id", productId.String())
		}

		log.Printf("failed to query chek sold status, err : %v\n", err)
		return false, helper.NewInternal()
	}

	return product.Sold, nil
}

func (r *ProductRepositoryImpl) SetSold(ctx context.Context, id uuid.UUID, status bool) error {

	query := `
	UPDATE 
		products 
	SET 
		sold = $1, updated_at = $2
	WHERE 
		id = $3 AND deleted = FALSE
	`

	result, err := r.DB.ExecContext(ctx, query, status, time.Now(), id)

	if err != nil {
		log.Printf("failed to query set sold product, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when set sold product, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", id.String())
	}

	return nil
}
