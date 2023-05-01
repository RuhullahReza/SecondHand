package service

import (
	"context"
	"testing"

	"github.com/RuhullahReza/SecondHand/db"
	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/RuhullahReza/SecondHand/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	DB := db.NewPostgresConnection()
	cld := db.NewCloudinaryConnection()
	accountRepository := repository.NewAccountRepository(DB)
	profileRepository := repository.NewProfileRepository(DB)
	dataRepository := repository.NewDataRepository(DB)
	imageRepository := repository.NewImageRepository(cld,DB)

	userService := NewUserService(accountRepository, profileRepository, imageRepository, dataRepository)

	request := &web.CreateUserRequest{
		Email:    "Ozza@gmail.com",
		Name:     "Ozza",
		Password: "Rahasia",
	}

	err := userService.Register(context.Background(), *request, "USER")
	require.NoError(t, err)
}

func TestFailCreate(t *testing.T) {
	DB := db.NewPostgresConnection()
	cld := db.NewCloudinaryConnection()
	accountRepository := repository.NewAccountRepository(DB)
	profileRepository := repository.NewProfileRepository(DB)
	dataRepository := repository.NewDataRepository(DB)
	imageRepository := repository.NewImageRepository(cld,DB)

	userService := NewUserService(accountRepository, profileRepository, imageRepository, dataRepository)

	request := &web.CreateUserRequest{
		Email:    "Ozza@gmail.com",
		Name:     "Ozza",
		Password: "Rahasia",
	}

	err := userService.Register(context.Background(), *request, "USER")
	expErr := helper.NewConflict("email", request.Email)
	require.EqualError(t, err, expErr.Error())
}

func TestLogin(t *testing.T) {
	DB := db.NewPostgresConnection()
	cld := db.NewCloudinaryConnection()
	accountRepository := repository.NewAccountRepository(DB)
	profileRepository := repository.NewProfileRepository(DB)
	dataRepository := repository.NewDataRepository(DB)
	imageRepository := repository.NewImageRepository(cld,DB)

	userService := NewUserService(accountRepository, profileRepository, imageRepository, dataRepository)

	request := &web.LoginRequest{
		Email:    "Ozza@gmail.com",
		Password: "Rahasia",
	}

	response := &web.LoginResponse{}

	err := userService.Login(context.Background(), *request, response)
	require.NoError(t, err)
}

func TestFailLogin(t *testing.T) {
	DB := db.NewPostgresConnection()
	cld := db.NewCloudinaryConnection()
	accountRepository := repository.NewAccountRepository(DB)
	profileRepository := repository.NewProfileRepository(DB)
	dataRepository := repository.NewDataRepository(DB)
	imageRepository := repository.NewImageRepository(cld,DB)

	userService := NewUserService(accountRepository, profileRepository, imageRepository, dataRepository)

	request := &web.LoginRequest{
		Email:    "Ozza123@gmail.com",
		Password: "Rahasiaa",
	}

	response := &web.LoginResponse{}

	err := userService.Login(context.Background(), *request, response)
	expErr := helper.NewAuthorization("Invalid email and password combination")
	require.EqualError(t, err, expErr.Error())
}

func TestUpdateProfile(t *testing.T) {
	DB := db.NewPostgresConnection()
	cld := db.NewCloudinaryConnection()
	accountRepository := repository.NewAccountRepository(DB)
	profileRepository := repository.NewProfileRepository(DB)
	dataRepository := repository.NewDataRepository(DB)
	imageRepository := repository.NewImageRepository(cld,DB)

	userService := NewUserService(accountRepository, profileRepository, imageRepository, dataRepository)

	id, err := uuid.Parse("4b2b723c-00f2-4e8f-91c9-611b3494ba9a")	
	require.NoError(t, err)

	request := &web.UpdateProfileRequest{
		Id: id,
		Name: "ozza Updated",
		City: "Bandung",
		Address: "alamat",
		PhoneNumber: "08123456",
	}

	err = userService.UpdateProfile(context.Background(), *request)
	require.NoError(t, err)
}

func TestNotFoundUpdateProfile(t *testing.T) {
	DB := db.NewPostgresConnection()
	cld := db.NewCloudinaryConnection()
	accountRepository := repository.NewAccountRepository(DB)
	profileRepository := repository.NewProfileRepository(DB)
	dataRepository := repository.NewDataRepository(DB)
	imageRepository := repository.NewImageRepository(cld,DB)

	userService := NewUserService(accountRepository, profileRepository, imageRepository, dataRepository)

	id, err := uuid.Parse("4b2b723c-00f2-4e8f-91c9-611b3494ba9b")	
	require.NoError(t, err)

	request := &web.UpdateProfileRequest{
		Id: id,
		Name: "ozza Updated",
		City: "Bandung",
		Address: "alamat",
		PhoneNumber: "08123456",
	}

	err = userService.UpdateProfile(context.Background(), *request)
	expErr := helper.NewNotFound("id", id.String())
	require.EqualError(t, err, expErr.Error())
}
