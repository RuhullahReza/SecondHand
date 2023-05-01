package repository

import (
	"context"
	"log"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) (*entity.Account, error)
	FindByEmail(ctx context.Context, email string) (*entity.Account, error) 
	FindById(ctx context.Context, id uuid.UUID) (*entity.Account, error)
}

type AccountRepositoryImpl struct {
	DB *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &AccountRepositoryImpl{
		DB: db,
	}
}

func (r *AccountRepositoryImpl) Create(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	
	query := "INSERT INTO accounts (email, password, role) VALUES ($1, $2, $3) RETURNING *"
	err := r.DB.GetContext(ctx, account, query, account.Email, account.Password, account.Role)

	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return account, helper.NewConflict("email", account.Email)
		}

		log.Printf("failed to query create account, err : %v\n", err)
		return account, helper.NewInternal()
	}

	return account, nil
}

func (r *AccountRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.Account, error) {
	
	account := &entity.Account{}

	query := "SELECT * FROM accounts WHERE email=$1"

	if err := r.DB.GetContext(ctx, account, query, email); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return account, helper.NewNotFound("email", email)
		}

		log.Printf("failed to query find by email, err : %v\n", err)
		return account, helper.NewInternal()
	}

	return account, nil
}

func (r *AccountRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	
	account := &entity.Account{}

	query := "SELECT * FROM accounts WHERE id=$1"

	if err := r.DB.GetContext(ctx, account, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return account, helper.NewNotFound("id", id.String())
		}

		log.Printf("failed to query find by email, err : %v\n", err)
		return account, helper.NewInternal()
	}

	return account, nil
}