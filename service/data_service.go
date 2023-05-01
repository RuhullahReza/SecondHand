package service

import (
	"context"
	"log"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/RuhullahReza/SecondHand/repository"
	"github.com/google/uuid"
	"github.com/leebenson/conform"
)

type DataService interface {
	CreateCity(ctx context.Context, req web.CreateDataRequest) error
	DeleteCity(ctx context.Context, id uuid.UUID) error
	GetAllCity(ctx context.Context, res *[]web.DataResponse) error
	CreateCategory(ctx context.Context, req web.CreateDataRequest) error
	GetAllCategory(ctx context.Context, res *[]web.DataResponse) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type DataServiceImpl struct {
	DataRepository repository.DataRepository
}

func NewDataService(repo repository.DataRepository) DataService {
	return &DataServiceImpl{
		DataRepository: repo,
	}
}

func (service *DataServiceImpl) CreateCity(ctx context.Context, req web.CreateDataRequest) error {

	err := conform.Strings(&req)
	if err != nil {
		log.Printf("error while sanitize on service create city, err : %v\n", err)
		return helper.NewInternal()
	}

	err = service.DataRepository.CreateCity(ctx, req.Name)
	if err != nil {
		return err
	}

	return nil
}

func (service *DataServiceImpl) DeleteCity(ctx context.Context, id uuid.UUID) error {

	err := service.DataRepository.DeleteCity(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *DataServiceImpl) GetAllCity(ctx context.Context, res *[]web.DataResponse) error {
	
	data, err := service.DataRepository.GetAllCity(ctx)
	if err != nil {
		return err
	}

	*res = data

	return nil
}

func (service *DataServiceImpl) CreateCategory(ctx context.Context, req web.CreateDataRequest) error {
	
	err := conform.Strings(&req)
	if err != nil {
		log.Printf("error while sanitize on service create category, err : %v\n", err)
		return helper.NewInternal()
	}

	err = service.DataRepository.CreateCategory(ctx, req.Name)
	if err != nil {
		return err
	}

	return nil
}

func (service *DataServiceImpl) GetAllCategory(ctx context.Context, res *[]web.DataResponse) error {
	
	data, err := service.DataRepository.GetAllCategory(ctx)
	if err != nil {
		return err
	}

	*res = data

	return nil
}

func (service *DataServiceImpl) DeleteCategory(ctx context.Context, id uuid.UUID) error {

	err := service.DataRepository.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}

	return nil
}