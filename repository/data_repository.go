package repository

import (
	"context"
	"log"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/entity"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DataRepository interface {
	CreateCity(ctx context.Context, name string) error
	DeleteCity(ctx context.Context, id uuid.UUID) error
	CheckCityName(ctx context.Context, name string) error
	GetAllCity(ctx context.Context) ([]web.DataResponse, error)
	CreateCategory(ctx context.Context, name string) error
	CheckCategory(ctx context.Context, name string) error
	GetAllCategory(ctx context.Context) ([]web.DataResponse, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type DataRepositoryImpl struct {
	DB *sqlx.DB
}

func NewDataRepository(db *sqlx.DB) DataRepository {
	return &DataRepositoryImpl{
		DB: db,
	}
}

func (r *DataRepositoryImpl) CreateCity(ctx context.Context, name string) error {

	query := "INSERT INTO cities (name) VALUES ($1)"
	_, err := r.DB.ExecContext(ctx, query, name)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return helper.NewConflict("city",name)
		}

		log.Printf("failed to query create city, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

func (r *DataRepositoryImpl) DeleteCity(ctx context.Context, id uuid.UUID) error {

	query := "DELETE FROM cities WHERE id = $1"
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to query delete city, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected when delete city, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("City",id.String())
	}
	
	return nil
}

func (r *DataRepositoryImpl) CheckCityName(ctx context.Context, name string) error {

	city := &entity.City{}

	query := "SELECT name FROM cities WHERE name=$1 LIMIT 1"

	if err := r.DB.GetContext(ctx, city, query, name); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return helper.NewNotFound("city name", name)
		}

		log.Printf("failed to query check city name, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

func (r *DataRepositoryImpl) GetAllCity(ctx context.Context) ([]web.DataResponse, error) {

	cities := []web.DataResponse{}

	query := "SELECT id, name FROM cities"
	rows, err := r.DB.QueryContext(ctx,query)
	if err != nil {
		log.Printf("failed to query get all city, err : %v\n", err)
		return cities, helper.NewInternal()
	}

	for rows.Next(){
		city := web.DataResponse{}
		err := rows.Scan(&city.Id,&city.Name)
		if err != nil {
			log.Printf("failed to scanning city, err : %v\n", err)
			return cities, helper.NewInternal()
		}

		cities = append(cities, city)
	}
	
	return cities, nil
}

func (r *DataRepositoryImpl) CreateCategory(ctx context.Context, name string) error {

	query := "INSERT INTO categories (name) VALUES ($1) "
	_, err := r.DB.ExecContext(ctx, query, name)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return helper.NewConflict("category",name)
		}

		log.Printf("failed to create category: %v, err: %v\n", name, err)
		return helper.NewInternal()
	}

	return nil
}

func (r *DataRepositoryImpl) CheckCategory(ctx context.Context, name string) error {

	category := &entity.Category{}

	query := "SELECT name FROM categories WHERE name=$1 LIMIT 1"

	if err := r.DB.GetContext(ctx, category, query, name); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return helper.NewNotFound("category name", name)
		}

		log.Printf("failed to query get check category, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

func (r *DataRepositoryImpl) DeleteCategory(ctx context.Context, id uuid.UUID) error {

	query := "DELETE FROM categories WHERE id = $1"
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to query delete category, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected when delete category, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("Category",id.String())
	}
	
	return nil
}

func (r *DataRepositoryImpl) GetAllCategory(ctx context.Context) ([]web.DataResponse, error) {

	categories := []web.DataResponse{}

	query := "SELECT id, name FROM categories"
	rows, err := r.DB.QueryContext(ctx,query)
	if err != nil {
		log.Printf("failed to query getting all category, err : %v\n", err)
		return categories, helper.NewInternal()
	}

	for rows.Next(){
		category := web.DataResponse{}
		err := rows.Scan(&category.Id,&category.Name)
		if err != nil {
			log.Printf("failed to scanning category, err : %v\n", err)
			return categories, helper.NewInternal()
		}
		categories = append(categories, category)
	}
	
	return categories, nil
}