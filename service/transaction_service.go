package service

import (
	"context"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/entity"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/RuhullahReza/SecondHand/repository"
	"github.com/google/uuid"
)

type TransactionService interface {
	Create(ctx context.Context, req web.TransactionRequest) error
	GetTransactionDetail(ctx context.Context, payload helper.Payload, id uuid.UUID, res *web.TransactionDetailResponse) error
	GetOfferByProduct(ctx context.Context, payload helper.Payload, id uuid.UUID, res *web.OfferByProduct) error
	GetOfferByBuyer(ctx context.Context, payload helper.Payload, id uuid.UUID, res *web.OfferByBuyer) error
	GetMyTransaction(ctx context.Context, id uuid.UUID, res *[]web.OfferWithAccount) error
	GetOfferByAccount(ctx context.Context, id uuid.UUID, res *[]web.OfferWithAccount) error
	UpdatePriceOffer(ctx context.Context, req web.TransactionUpdateRequest) error
	UpdateStatus(ctx context.Context, productId uuid.UUID, sellerId uuid.UUID, res *bool) error
}

type TransactionServiceImpl struct {
	ProductRepository		repository.ProductRepository
	ProfileRepository		repository.ProfileRepository
	TransactionRepository 	repository.TransactionRepository
}

func NewTransactionSerive(
	productRepository repository.ProductRepository, 
	profileRepository repository.ProfileRepository,
	transactionRepository repository.TransactionRepository,
	) TransactionService {
	return &TransactionServiceImpl{
		ProductRepository: productRepository,
		ProfileRepository: profileRepository,
		TransactionRepository: transactionRepository,
	}
}

func (service *TransactionServiceImpl) Create(ctx context.Context, req web.TransactionRequest) error {

	validProfile, err := service.ProfileRepository.CheckProfile(ctx, req.BuyerId)
	if err != nil {
		return err
	}

	if !validProfile {
		return helper.NewBadRequest("complete your profile first")
	}
	
	sellerId, err := service.ProductRepository.GetOwnerId(ctx, req.ProductId)
	if err != nil {
		return err
	}

	if sellerId == req.BuyerId {
		return helper.NewBadRequest("you cannot buy your own product")
	}

	newTransaction := &entity.Transaction{
		SellerId: sellerId,
		BuyerId: req.BuyerId,
		ProductId: req.ProductId,
		PriceOffer: req.Price,
	}

	err = service.TransactionRepository.Create(ctx, newTransaction)
	if err != nil {
		return err
	}

	return nil
}

func (service *TransactionServiceImpl) GetTransactionDetail(ctx context.Context, payload helper.Payload, id uuid.UUID, res *web.TransactionDetailResponse) error {
	
	transactions, err := service.TransactionRepository.GetTransactionDetail(ctx, id)
	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return helper.NewNotFound("id",id.String())
	}

	if transactions[0].BuyerId != payload.UserId && transactions[0].SellerId != payload.UserId && payload.Role != "ADMIN" {
		return helper.NewAuthorization("only buyer, seller, and admin can access")
	}

	*res = transactions[0]

	return nil
}

func (service *TransactionServiceImpl) GetOfferByProduct(ctx context.Context, payload helper.Payload, id uuid.UUID, res *web.OfferByProduct) error {
	
	products, err := service.ProductRepository.GetProduct(ctx, id)
	if err != nil {
		return err
	}

	if len(products) == 0 {
		return helper.NewNotFound("id",id.String())
	}

	if products[0].OwnerId != payload.UserId && payload.Role != "ADMIN" {
		return helper.NewAuthorization("seller and admin can access")
	}

	offers, err := service.TransactionRepository.GetOfferByProduct(ctx, id)
	if err != nil {
		return err
	}

	*res = products[0]
	res.Offer = offers

	return nil
}

func (service *TransactionServiceImpl) GetOfferByBuyer(ctx context.Context, payload helper.Payload, id uuid.UUID, res *web.OfferByBuyer) error {
	
	profile, err := service.ProfileRepository.GetProfileById(ctx, id)
	if err != nil {
		return err
	}

	res.BuyerId = id
	res.BuyerName = profile.Name
	res.BuyerCity = profile.City
	res.BuyerImage = profile.ImageUrl

	offers, err := service.TransactionRepository.GetOfferByBuyer(ctx, id, payload.UserId)
	if err != nil {
		return err
	}

	res.Offer = offers

	return nil
}

func (service *TransactionServiceImpl) GetOfferByAccount(ctx context.Context, id uuid.UUID, res *[]web.OfferWithAccount) error {
	
	transactions, err := service.TransactionRepository.GetOfferByAccount(ctx, id)
	if err != nil {
		return err
	}
	
	*res = transactions

	return nil
}

func (service *TransactionServiceImpl) GetMyTransaction(ctx context.Context, id uuid.UUID, res *[]web.OfferWithAccount) error {
	
	transactions, err := service.TransactionRepository.GetMyTransaction(ctx, id)
	if err != nil {
		return err
	}
	
	*res = transactions

	return nil
}

func (service *TransactionServiceImpl) UpdatePriceOffer(ctx context.Context, req web.TransactionUpdateRequest) error {
	
	err := service.TransactionRepository.UpdatePrice(ctx, req.Price, req.TransactionId, req.BuyerId)
	if err != nil {
		return err
	}

	return nil
}

func (service *TransactionServiceImpl) UpdateStatus(ctx context.Context, productId uuid.UUID, sellerId uuid.UUID, res *bool) error {
	
	accepted, err := service.TransactionRepository.CheckStatus(ctx, productId)
	if err != nil {
		return err
	}

	*res = !accepted
	err = service.TransactionRepository.SetStatus(ctx, !accepted, productId, sellerId)
	if err != nil {
		return err
	}

	return nil
}

func (service *TransactionServiceImpl) DeleteTransaction(ctx context.Context, payload helper.Payload, product_id uuid.UUID) error {

	transaction, err := service.TransactionRepository.GetTransactionById(ctx, product_id)
	if err != nil {
		return err
	}

	if payload.UserId != transaction.BuyerId || payload.UserId != transaction.SellerId || payload.Role != "ADMIN" {
		return helper.NewAuthorization("only buyer, seller, and admin can access")
	}

	err = service.TransactionRepository.DeleteOne(ctx, product_id)
	if err != nil {
		return err
	}

	return nil
}