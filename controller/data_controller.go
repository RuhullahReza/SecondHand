package controller

import (
	"fmt"
	"net/http"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/RuhullahReza/SecondHand/service"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/google/uuid"
)

type DataController interface {
	AddCity(c *gin.Context)
	GetAllCity(c *gin.Context)
	DeleteCity(c *gin.Context)
	AddCategory(c *gin.Context)
	GetAllCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

type DataControllerImpl struct {
	DataService service.DataService
	Translator	ut.Translator
}

func NewDataController(service service.DataService, translator ut.Translator) DataController {
	return &DataControllerImpl{
		DataService: service,
		Translator: translator,
	}
}

func (d *DataControllerImpl) AddCity(c *gin.Context) {
	var req web.CreateDataRequest

	if ok := helper.BindData(c, d.Translator, &req); !ok {
		return
	}

	err := d.DataService.CreateCity(c, req)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("data City : %s successfully created", req.Name),
	})
}

func (d *DataControllerImpl) GetAllCity(c *gin.Context) {

	var res []web.DataResponse
	err := d.DataService.GetAllCity(c, &res)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Data":res,
	})
}

func (d *DataControllerImpl) DeleteCity(c *gin.Context) {

	var req web.GetByIdRequest
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

	err = d.DataService.DeleteCity(c, id)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Successfully delete city with Id %s", id),
	})
}

func (d *DataControllerImpl) AddCategory(c *gin.Context) {
	var req web.CreateDataRequest

	if ok := helper.BindData(c, d.Translator, &req); !ok {
		return
	}

	err := d.DataService.CreateCategory(c, req)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("data Category : %s successfully created", req.Name),
	})
}

func (d *DataControllerImpl) GetAllCategory(c *gin.Context) {

	var res []web.DataResponse
	err := d.DataService.GetAllCategory(c, &res)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Data":res,
	})
}

func (d *DataControllerImpl) DeleteCategory(c *gin.Context) {

	var req web.GetByIdRequest
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

	err = d.DataService.DeleteCategory(c, id)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Successfully delete category with Id %s", id),
	})
}
