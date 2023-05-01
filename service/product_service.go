package service

import (
	"context"
	"log"
	"mime/multipart"
	"time"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/entity"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/RuhullahReza/SecondHand/repository"
	"github.com/google/uuid"
	"github.com/leebenson/conform"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req web.CreateProductRequest, userId uuid.UUID) error
	GetAllProduct(ctx context.Context, res *[]web.ProductResponse) error
	GetByCategory(ctx context.Context, category string, res *[]web.ProductResponse) error
	GetByAccount(ctx context.Context, id uuid.UUID, status bool, published bool, res *[]web.ProductResponse) error
	GetProductById(ctx context.Context, payload helper.Payload, id uuid.UUID, res *web.ProductDetailResponse) error 
	UpdateProduct(ctx context.Context, req web.UpdateProductRequest) error
	AddProductImage(ctx context.Context, payload helper.Payload, id uuid.UUID, imageFileHeader *multipart.FileHeader) error
	SetThumbnail(ctx context.Context, accountId uuid.UUID, productId uuid.UUID, imageId uuid.UUID) error
	DeleteProduct(ctx context.Context, payload helper.Payload, id uuid.UUID) error
	DeleteImageProduct(ctx context.Context, accountId uuid.UUID, productId uuid.UUID, imageId uuid.UUID) error 
	UpdatePublished(ctx context.Context, payload helper.Payload, id uuid.UUID, res *bool) error
	UpdateSold(ctx context.Context, payload helper.Payload, id uuid.UUID, res *bool) error
}

type ProductServiceImpl struct {
	ProductRepository 	repository.ProductRepository
	ProfileRepository 	repository.ProfileRepository 
	DataRepository	  	repository.DataRepository
	ImageRepository		repository.ImageRepository
}

func NewProductService(
	productRepository repository.ProductRepository, 
	profileRepository repository.ProfileRepository, 
	dataRepository repository.DataRepository,
	imageRepository	repository.ImageRepository,
	) ProductService {
	return &ProductServiceImpl{
		ProductRepository: productRepository,
		ProfileRepository: profileRepository,
		DataRepository: dataRepository,
		ImageRepository: imageRepository,
	}
}

func (service *ProductServiceImpl) CreateProduct(ctx context.Context, req web.CreateProductRequest, userId uuid.UUID) error {
	
	validProfile, err := service.ProfileRepository.CheckProfile(ctx,userId)
	if err != nil {
		return err
	}

	if !validProfile {
		return helper.NewBadRequest("complete your profile first")
	}

	err = conform.Strings(&req)
	if err != nil {
		log.Printf("error while sanitize on service create product, err : %v\n", err)
		return helper.NewInternal()
	}

	err = service.DataRepository.CheckCategory(ctx, req.Category)
	if err != nil {
		return err
	}

	product := &entity.Product{
		AccountId: userId,
		Name: req.Name,
		Price: req.Price,
		Category: req.Category,
		Description: req.Description,
	}

	err = service.ProductRepository.Create(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProductServiceImpl) GetProductById(ctx context.Context, payload helper.Payload, id uuid.UUID, res *web.ProductDetailResponse) error {
	
	isOwner, err := service.ProductRepository.IsOwner(ctx, payload.UserId, id)
	if err != nil {
		return err
	}

	var data []web.ProductDetailResponse
	if isOwner || payload.Role == "Admin" {
		data, err = service.ProductRepository.OwnerGetOne(ctx, id)
		if err != nil {
			return err
		}
	} else {
		data, err = service.ProductRepository.GetOne(ctx, id)
		if err != nil {
			return err
		}
	}

	if len(data) == 0 {
		return helper.NewNotFound("id",id.String())
	}

	*res = data[0]

	imageList, err := service.ImageRepository.GetByProductId(ctx, id)
	if err != nil {
		return err
	}
	
	res.ProductImages = imageList

	return nil
}

func (service *ProductServiceImpl) GetAllProduct(ctx context.Context, res *[]web.ProductResponse) error {
	
	data, err := service.ProductRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	*res = data

	return nil
}

func (service *ProductServiceImpl) GetByCategory(ctx context.Context, category string, res *[]web.ProductResponse) error {
	
	data, err := service.ProductRepository.GetByCategory(ctx, category)
	if err != nil {
		return err
	}

	*res = data

	return nil
}

func (service *ProductServiceImpl) GetByAccount(ctx context.Context, id uuid.UUID, status bool, published bool, res *[]web.ProductResponse) error {
	
	data, err := service.ProductRepository.GetByAccount(ctx, id, status, published)
	if err != nil {
		return err
	}

	*res = data

	return nil
}

func (service *ProductServiceImpl) UpdateProduct(ctx context.Context, req web.UpdateProductRequest) error {

	err := service.ProductRepository.CheckOwner(ctx, req.AccountId, req.Id)
	if err != nil {
		return err
	}

	err = conform.Strings(&req)
	if err != nil {
		log.Printf("error while sanitize on service update product, err : %v\n", err)
		return helper.NewInternal()
	}

	err = service.DataRepository.CheckCategory(ctx, req.Category)
	if err != nil {
		return err
	}

	product := &entity.Product{
		Id: req.Id,
		Name: req.Name,
		Price: req.Price,
		Category: req.Category,
		UpdatedAt: time.Now(),
		Description: req.Description,
	}

	err = service.ProductRepository.Update(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProductServiceImpl) AddProductImage(ctx context.Context, payload helper.Payload, id uuid.UUID, imageFileHeader *multipart.FileHeader) error {

	err := service.ProductRepository.CheckOwner(ctx, payload.UserId, id)
	if err != nil {
		return err
	}

	imageFile, err := imageFileHeader.Open()
	if err != nil {
		log.Printf("error while open image on service add product image, err: %v\n", err)
		return helper.NewInternal()
	}

	path, err := service.ImageRepository.Upload(ctx,imageFile,"secondHand-go/product")
	if err != nil {
		return err
	}

	hasThumbnail, err := service.ProductRepository.CheckThumbnail(ctx,id)
	if err != nil {
		return err
	}

	image := &entity.Image{
		ProductId: id,
		Url: path,
	}
	   
	err = service.ImageRepository.Create(ctx, image)
	if err != nil {
		return err
	}

	if !hasThumbnail {
		updatedProduct := &entity.Product{
			Id: id,
			AccountId: payload.UserId,
			Thumbnail: path,
			UpdatedAt: time.Now(),
		}

		err := service.ProductRepository.SetThumbnail(ctx,updatedProduct)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *ProductServiceImpl) SetThumbnail(ctx context.Context, accountId uuid.UUID, productId uuid.UUID, imageId uuid.UUID) error {

	image, err := service.ImageRepository.GetPathById(ctx, imageId)
	if err != nil {
		return err
	}

	updatedProduct := &entity.Product{
		Id: productId,
		AccountId: accountId,
		Thumbnail: image.Url,
		UpdatedAt: time.Now(),
	}

	err = service.ProductRepository.SetThumbnail(ctx, updatedProduct)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProductServiceImpl) DeleteImageProduct(ctx context.Context, accountId uuid.UUID, productId uuid.UUID, imageId uuid.UUID) error {

	err := service.ProductRepository.CheckOwner(ctx, accountId, productId)
	if err != nil {
		return err
	}

	image, err := service.ImageRepository.GetPathById(ctx, imageId)
	if err != nil {
		return err
	}

	isThumbnail, err :=  service.ProductRepository.IsThumbnail(ctx, productId, image.Url)
	if err != nil {
		return err
	}

	if isThumbnail {
		return helper.NewBadRequest("cannot delete thumbnail image")
	}

	err = service.ImageRepository.Delete(ctx, image.Url)
	if err != nil {
		return err
	}

	err = service.ImageRepository.DeleteById(ctx, imageId)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProductServiceImpl) DeleteProduct(ctx context.Context, payload helper.Payload, id uuid.UUID) error {

	if payload.Role != "ADMIN" {
		err := service.ProductRepository.CheckOwner(ctx, payload.UserId, id)
		if err != nil {
			return err
		}
	}
	
	err := service.ProductRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *ProductServiceImpl) UpdatePublished(ctx context.Context, payload helper.Payload, id uuid.UUID, res *bool) error {

	if payload.Role != "ADMIN" {
		err := service.ProductRepository.CheckOwner(ctx, payload.UserId, id)
		if err != nil {
			return err
		}
	}

	hasThumbnail, err := service.ProductRepository.CheckThumbnail(ctx, id)
	if err != nil {
		return err
	}

	if !hasThumbnail {
		return helper.NewBadRequest("add thumbnail before publish product")
	}

	status, err := service.ProductRepository.CheckPublished(ctx, id)
	if err != nil {
		return err
	}

	err = service.ProductRepository.Publish(ctx, id, !status)
	if err != nil {
		return err
	}

	*res = !status

	return nil
}

func (service *ProductServiceImpl) UpdateSold(ctx context.Context, payload helper.Payload, id uuid.UUID, res *bool) error {

	if payload.Role != "ADMIN" {
		err := service.ProductRepository.CheckOwner(ctx, payload.UserId, id)
		if err != nil {
			return err
		}
	}

	status, err := service.ProductRepository.CheckSold(ctx, id)
	if err != nil {
		return err
	}

	err = service.ProductRepository.SetSold(ctx, id, !status)
	if err != nil {
		return err
	}

	*res = !status

	return nil
}