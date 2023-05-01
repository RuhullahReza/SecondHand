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
	"github.com/RuhullahReza/SecondHand/util"
	"github.com/google/uuid"
	"github.com/leebenson/conform"
)

type UserService interface {
	Register(ctx context.Context, req web.CreateUserRequest, role string) error
	Login(ctx context.Context, req web.LoginRequest, res *web.LoginResponse) error
	GetProfile(ctx context.Context, res *web.GetProfileResponse, id uuid.UUID) error
	UpdateProfile(ctx context.Context, req web.UpdateProfileRequest) error
	UpdateImage(ctx context.Context, id uuid.UUID, imageFileHeader *multipart.FileHeader) error
}

type UserServiceImpl struct {
	AccountRepository  	repository.AccountRepository
	ProfileRepository  	repository.ProfileRepository
	ImageRepository		repository.ImageRepository
	DataRepository		repository.DataRepository
}

func NewUserService(
	accountRepository repository.AccountRepository, 
	profileRepository repository.ProfileRepository, 
	imageRepository repository.ImageRepository, 
	dataRepository repository.DataRepository, 
	) UserService {
	return &UserServiceImpl{
		AccountRepository: accountRepository,
		ProfileRepository: profileRepository,
		ImageRepository: imageRepository,
		DataRepository: dataRepository,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, req web.CreateUserRequest, role string) error  {

	err := conform.Strings(&req)
	if err != nil {
		log.Printf("error while sanitize on service register, err : %v\n", err)
		return helper.NewInternal()
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Printf("error while hashing password on service register, err : %v\n", err)
		return helper.NewInternal()
	}

	account := &entity.Account{
		Email: req.Email,
		Password: hashedPassword,
		Role: role,
	}

	newAccount, err := service.AccountRepository.Create(context.Background(),account)
	if err != nil {
		return err
	}
	
	profile :=  &entity.Profile{
		Id: newAccount.Id,
		Name: req.Name,
	}

	err = service.ProfileRepository.Create(ctx, profile)
	if err != nil {
		return err
	}
	
	return nil
}

func (service *UserServiceImpl) Login(ctx context.Context, req web.LoginRequest, res *web.LoginResponse) error {

	err := conform.Strings(&req)
	if err != nil {
		log.Printf("error while sanitize on service login, err : %v\n", err)
		return helper.NewInternal()
	}

	foundUser, err := service.AccountRepository.FindByEmail(ctx, req.Email)

	if err != nil {
		return err
	}

	matchedPassword, err :=  util.ComparePasswords(foundUser.Password, req.Password)

	if err != nil {
		log.Printf("error while compare password on service login, err : %v\n", err)
		return helper.NewInternal()
	}

	if !matchedPassword {
		return helper.NewAuthorization("Invalid email and password combination")
	}

	foundProfile, err := service.ProfileRepository.GetNameById(ctx, foundUser.Id)

	if err != nil {
		return err
	}

	token, err := helper.SignToken(foundUser.Id, foundProfile.Name, foundUser.Role)

	if err != nil {
		log.Printf("error while sign JWT on service login, err : %v\n", err)
		return helper.NewInternal()
	}

	res.Token = token

	return nil
}

func (service *UserServiceImpl) GetProfile(ctx context.Context, res *web.GetProfileResponse, id uuid.UUID) error {

	profile, err := service.ProfileRepository.GetProfileById(ctx, id)

	if err != nil {
		return err
	}

	res.Id = id
	res.Name = profile.Name
	res.City = profile.City
	res.Address = profile.Address
	res.PhoneNumber = profile.PhoneNumber
	res.ImageUrl = profile.ImageUrl

	return nil
}

func (service *UserServiceImpl) UpdateProfile(ctx context.Context, req web.UpdateProfileRequest) error {

	err := conform.Strings(&req)
	if err != nil {
		log.Printf("error while sanitize on service update profile, err : %v\n", err)
		return helper.NewInternal()
	}

	err = service.DataRepository.CheckCityName(ctx, req.City)
	if err != nil {
		return err
	}

	newProfile := &entity.Profile{
		Id: req.Id,
		Name: req.Name,
		City: req.City,
		Address: req.Address,
		PhoneNumber: req.PhoneNumber,
		UpdatedAt: time.Now(),
	}

	err = service.ProfileRepository.Update(ctx, newProfile)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserServiceImpl) UpdateImage(ctx context.Context, id uuid.UUID, imageFileHeader *multipart.FileHeader) error {

	profile, err := service.ProfileRepository.GetProfileImage(ctx, id)
	if err != nil {
		return err
	}

	if profile.ImageUrl != "" {

		err := service.ImageRepository.Delete(ctx, profile.ImageUrl)	
		if err != nil {
			return err
		}	
	}

	imageFile, err := imageFileHeader.Open()
	if err != nil {
		log.Printf("error while open image on service update image, err: %v\n", err)
		return helper.NewInternal()
	}

	path, err := service.ImageRepository.Upload(ctx,imageFile,"secondHand-go/profile")
	if err != nil {
		return err
	}

	profile.Id = id
	profile.ImageUrl = path

	err = service.ProfileRepository.UpdateImage(ctx, profile)
	if err != nil {
		return err
	}
	
	return nil
}


