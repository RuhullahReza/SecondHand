package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/web"
	"github.com/RuhullahReza/SecondHand/service"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/google/uuid"
)

type ProductController interface {
	AddProduct(c *gin.Context)
	GetAllProduct(c *gin.Context)
	GetProductById(c *gin.Context)
	GetMyProduct(c *gin.Context)
	GetProductByAccount(c *gin.Context)
	GetByCategory(c *gin.Context)
	UpdateProduct(c *gin.Context)
	AddProductImage(c *gin.Context)
	UpdateProductThumbnail(c *gin.Context)
	DeleteProduct(c *gin.Context)
	DeleteProductImage(c *gin.Context)
	UpdatePublishStatus(c *gin.Context)
	UpdateSoldStatus(c *gin.Context)
	
}

type ProductControllerImpl struct {
	ProductService service.ProductService
	Translator ut.Translator
}

func NewProductController(service service.ProductService, translator ut.Translator) ProductController {
	return &ProductControllerImpl{
		ProductService: service,
		Translator: translator,
	}
}

func (p *ProductControllerImpl) AddProduct(c *gin.Context) {

	payload := c.MustGet("payload")

	var req web.CreateProductRequest

	if ok := helper.BindData(c, p.Translator, &req); !ok {
		return
	}

	account_id := payload.(helper.Payload).UserId

	err := p.ProductService.CreateProduct(c, req, account_id)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Product : %s successfully created", req.Name),
	})
}

func (p *ProductControllerImpl) GetAllProduct(c *gin.Context) {

	var res []web.ProductResponse
	err := p.ProductService.GetAllProduct(c, &res)
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

func (p *ProductControllerImpl) GetByCategory(c *gin.Context) {

	var param web.GetByPath
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	var res []web.ProductResponse
	err := p.ProductService.GetByCategory(c, param.Path, &res)
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

func (p *ProductControllerImpl) GetMyProduct(c *gin.Context) {

	payload := c.MustGet("payload")

	id := payload.(helper.Payload).UserId

	
	sold, err :=strconv.ParseBool(c.Query("sold")) 
	if err != nil {
		sold = false
	}

	published, err :=strconv.ParseBool(c.Query("published")) 
	if err != nil {
		published = true
	}
	
	var res []web.ProductResponse
	err = p.ProductService.GetByAccount(c, id, sold, published, &res)
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


func (p *ProductControllerImpl) GetProductByAccount(c *gin.Context) {

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	account_id, err := uuid.Parse(param.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var res []web.ProductResponse
	err = p.ProductService.GetByAccount(c, account_id, false, true, &res)
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

func (p *ProductControllerImpl) GetProductById(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	product_id, err := uuid.Parse(param.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var res web.ProductDetailResponse
	err = p.ProductService.GetProductById(c, payload, product_id, &res)
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

func (p *ProductControllerImpl) UpdateProduct(c *gin.Context) {

	payload := c.MustGet("payload")

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	productId, err := uuid.Parse(param.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var req web.UpdateProductRequest

	if ok := helper.BindData(c, p.Translator, &req); !ok {
		return
	}

	req.AccountId = payload.(helper.Payload).UserId
	req.Id = productId

	err = p.ProductService.UpdateProduct(c, req)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Product id: %s successfully updated", req.Id),
	})
}

func (p *ProductControllerImpl) AddProductImage(c *gin.Context) {
	payload := c.MustGet("payload").(helper.Payload)

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	productId, err := uuid.Parse(param.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

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
	
	err = p.ProductService.AddProductImage(c, payload, productId, imageFileHeader)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("successfuly upload image for product id : %s", productId),
	})
}

func (p *ProductControllerImpl) UpdateProductThumbnail(c *gin.Context) {

	payload := c.MustGet("payload")


	var req web.ProductImageRequest

	if ok := helper.BindData(c, p.Translator, &req); !ok {
		return
	}

	AccountId := payload.(helper.Payload).UserId

	err := p.ProductService.SetThumbnail(c, AccountId, req.ProductId, req.ImageId)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Product id: %s thumbnail successfully updated", req.ProductId),
	})
}

func (p *ProductControllerImpl) DeleteProductImage(c *gin.Context) {

	payload := c.MustGet("payload")


	var req web.ProductImageRequest

	if ok := helper.BindData(c, p.Translator, &req); !ok {
		return
	}

	AccountId := payload.(helper.Payload).UserId

	err := p.ProductService.DeleteImageProduct(c, AccountId, req.ProductId, req.ImageId)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Image id: %s from Product Id: %s successfully deleted", req.ImageId, req.ProductId),
	})
}

func (p *ProductControllerImpl) DeleteProduct(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

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

	err = p.ProductService.DeleteProduct(c, payload, id)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Successfully delete product with Id : %s", id),
	})
}

func (p *ProductControllerImpl) UpdatePublishStatus(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

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

	var res bool
	err = p.ProductService.UpdatePublished(c, payload, id, &res)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Successfully updated publish status product with Id : %s to : %v", id, res),
	})
}

func (p *ProductControllerImpl) UpdateSoldStatus(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

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

	var res bool
	err = p.ProductService.UpdateSold(c, payload, id, &res)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Successfully updated sold status product with Id : %s to : %v", id, res),
	})
}