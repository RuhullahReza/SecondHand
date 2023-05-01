package repository

import (
	"context"
	"log"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProfileRepository interface {
	Create(ctx context.Context, p *entity.Profile) error
	Update(ctx context.Context, p *entity.Profile) error
	GetNameById(ctx context.Context, id uuid.UUID) (*entity.Profile, error)
	CheckProfile(ctx context.Context, id uuid.UUID) (bool, error)
	GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, error)
	GetProfileImage(ctx context.Context, id uuid.UUID) (*entity.Profile, error)
	UpdateImage(ctx context.Context, p *entity.Profile) error
}

type ProfileRepositoryImpl struct {
	DB *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) ProfileRepository {
	return &ProfileRepositoryImpl{
		DB: db,
	}
}

func (r *ProfileRepositoryImpl) Create(ctx context.Context, p *entity.Profile) error {

	query := "INSERT INTO profiles (id, name) VALUES ($1, $2)"
	_, err := r.DB.ExecContext(ctx, query, p.Id, p.Name)

	if err != nil {
		log.Printf("failed to query create profile, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

func (r *ProfileRepositoryImpl) Update(ctx context.Context, p *entity.Profile) error {
	
	query := "UPDATE profiles SET name = $1, city = $2, address = $3, phone_number = $4, updated_at = $5 WHERE id = $6"

	result, err := r.DB.ExecContext(ctx, query, p.Name, p.City, p.Address, p.PhoneNumber, p.UpdatedAt, p.Id)

	if err != nil {
		log.Printf("failed to query update profile, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when update profile, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", p.Id.String())
	}

	return nil
}

func (r *ProfileRepositoryImpl) GetNameById(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
	
	profile := &entity.Profile{}

	query := "SELECT name FROM profiles WHERE id=$1"

	if err := r.DB.GetContext(ctx, profile, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return profile, helper.NewNotFound("id", id.String())
		}
		
		log.Printf("failed to query get name by Id, err : %v\n", err)
		return profile, helper.NewInternal()
	}

	return profile, nil
}

func (r *ProfileRepositoryImpl) GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
	
	profile := &entity.Profile{}

	query := `
	SELECT 
		COALESCE(name, '') as name, 
		COALESCE(city, '') as city,
		COALESCE(address, '') as address, 
		COALESCE(phone_number, '') as phone_number, 
		COALESCE(image_url, '') as image_url 
	FROM profiles 
	WHERE id=$1
	`

	if err := r.DB.GetContext(ctx, profile, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return profile, helper.NewNotFound("id", id.String())
		}
		
		log.Printf("failed to query get profile by Id, err : %v\n", err)
		return profile, helper.NewInternal()
	}

	return profile, nil
}

func (r *ProfileRepositoryImpl) CheckProfile(ctx context.Context, id uuid.UUID) (bool, error) {
	
	profile := &entity.Profile{}

	query := `
	SELECT 
		COALESCE(name, '') as name, 
		COALESCE(city, '') as city,
		COALESCE(address, '') as address, 
		COALESCE(phone_number, '') as phone_number, 
		COALESCE(image_url, '') as image_url 
	FROM profiles 
	WHERE id=$1
	`

	if err := r.DB.GetContext(ctx, profile, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, helper.NewNotFound("id", id.String())
		}
		
		log.Printf("failed to query check profile by Id, err : %v\n", err)
		return false, helper.NewInternal()
	}

	if profile.Name == "" {
		return false, nil
	}

	if profile.City == "" {
		return false, nil
	}

	if profile.Address == "" {
		return false, nil
	}

	if profile.PhoneNumber == "" {
		return false, nil
	}

	return true, nil
}

func (r *ProfileRepositoryImpl) GetProfileImage(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
	
	profile := &entity.Profile{}

	query := `
	SELECT 
		COALESCE(image_url, '') as image_url 
	FROM profiles 
	WHERE id=$1
	`

	if err := r.DB.GetContext(ctx, profile, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return profile, helper.NewNotFound("id", id.String())
		}
		
		log.Printf("failed to query get profile image, err : %v\n", err)
		return profile, helper.NewInternal()
	}

	return profile, nil
}

func (r *ProfileRepositoryImpl) UpdateImage(ctx context.Context, p *entity.Profile) error {
	
	query := "UPDATE profiles SET image_url = $1, updated_at = $2 WHERE id = $3"

	result, err := r.DB.ExecContext(ctx, query, p.ImageUrl, p.UpdatedAt, p.Id)

	if err != nil {
		log.Printf("failed to query update profile image, err : %v\n", err)
		return helper.NewInternal()
	}

	row, err := result.RowsAffected()

	if err != nil {
		log.Printf("failed to get rows affected when update profile image, err : %v\n", err)
		return helper.NewInternal()
	}

	if row == 0 {
		return helper.NewNotFound("id", p.Id.String())
	}

	return nil
}

