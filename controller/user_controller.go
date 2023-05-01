package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/RuhullahReza/SecondHand/service"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/google/uuid"
)

type UserController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RegisterAdmin(c *gin.Context)
	Update(c *gin.Context)
	MyProfile(c *gin.Context)
	Profile(c *gin.Context)
	UpdateImage(c *gin.Context)
}

type UserControllerImpl struct {
	Service 	service.UserService
	Translator	ut.Translator
}

func NewUserController(service service.UserService, translator ut.Translator) UserController {
	return &UserControllerImpl{
		Service: service,
		Translator: translator,
	}
}

func (u *UserControllerImpl) Login(c *gin.Context) {
	var req web.LoginRequest

	if ok := helper.BindData(c, u.Translator, &req); !ok {
		return
	}
 
	var res web.LoginResponse 

	err := u.Service.Login(c, req, &res)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.SetCookie("token", res.Token, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, res)
}

func (u *UserControllerImpl) Register(c *gin.Context) {
	var req web.CreateUserRequest

	if ok := helper.BindData(c, u.Translator, &req); !ok {
		return
	}
	
	err := u.Service.Register(c, req, "USER")
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("account with email : %s successfully created", req.Email),
	})
}

func (u *UserControllerImpl) RegisterAdmin(c *gin.Context) {
	var req web.CreateUserRequest

	if ok := helper.BindData(c, u.Translator, &req); !ok {
		return
	}

	err := u.Service.Register(c, req, "ADMIN")
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("account with email : %s successfully created", req.Email),
	})
}

func (u *UserControllerImpl) Update(c *gin.Context) {
	
	payload:= c.MustGet("payload")

	var req web.UpdateProfileRequest

	if ok := helper.BindData(c, u.Translator, &req); !ok {
		return
	}

	req.Id = payload.(helper.Payload).UserId
 
	err := u.Service.UpdateProfile(c, req)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status" : "Success",
		"Message" : "profile successfully updated",
	})
}

func (u *UserControllerImpl) MyProfile(c *gin.Context) {
	
	payload := c.MustGet("payload")

	id := payload.(helper.Payload).UserId

	var res web.GetProfileResponse

	err := u.Service.GetProfile(c, &res, id)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status" : "Success",
		"Message" : res,
	})
}

func (u *UserControllerImpl) Profile(c *gin.Context) {

	var req web.GetAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return
	}
	
	id, err := uuid.Parse(req.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var res web.GetProfileResponse

	err = u.Service.GetProfile(c, &res, id)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status" : "Success",
		"Message" : res,
	})
}

func (u *UserControllerImpl) UpdateImage(c *gin.Context) {
	payload := c.MustGet("payload")

	id := payload.(helper.Payload).UserId

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4194304)

	imageFileHeader, err := c.FormFile("imageFile")

	if err != nil {
		log.Printf("Unable parse multipart/form-data: %+v", err)

		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("Max request body size is %v bytes\n", 4194304),
			})
			return
		}
		e := helper.NewBadRequest("Unable to parse multipart/form-data")
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	err = u.Service.UpdateImage(c, id, imageFileHeader)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status" : "Success",
		"Message" : "Successfuly update image",
	})
}