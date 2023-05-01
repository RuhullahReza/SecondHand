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

type TransactionController interface {
	AddTransaction(c *gin.Context)
	GetTransactionDetail(c *gin.Context)
	GetTransactionByProduct(c *gin.Context)
	GetTransactionByBuyer(c *gin.Context)
	GetTransactionByAccount(c *gin.Context)
	GetMyTransaction(c *gin.Context)
	UpdatePriceOffer(c *gin.Context)
	UpdateTransactionStatus(c *gin.Context)
}

type TransactionControllerImpl struct {
	TransactionService service.TransactionService
	Translator ut.Translator
}

func NewTransactionController(transactionService service.TransactionService, translator ut.Translator) TransactionController {
	return &TransactionControllerImpl{
		TransactionService: transactionService,
		Translator: translator,
	}
}

func (t *TransactionControllerImpl) AddTransaction(c *gin.Context) {
	payload := c.MustGet("payload").(helper.Payload)

	var req web.TransactionRequest

	if ok := helper.BindData(c, t.Translator, &req); !ok {
		return
	}

	req.BuyerId = payload.UserId

	err := t.TransactionService.Create(c, req)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Transaction for product id: %s successfully created", req.ProductId),
	})
}

func (t *TransactionControllerImpl) GetTransactionDetail(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	transaction_id, err := uuid.Parse(param.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var res web.TransactionDetailResponse
	err = t.TransactionService.GetTransactionDetail(c, payload, transaction_id, &res)
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

func (t *TransactionControllerImpl) GetTransactionByProduct(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	id, err := uuid.Parse(param.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var res web.OfferByProduct
	err = t.TransactionService.GetOfferByProduct(c, payload, id, &res)
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

func (t *TransactionControllerImpl) GetTransactionByBuyer(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	id, err := uuid.Parse(param.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var res web.OfferByBuyer
	err = t.TransactionService.GetOfferByBuyer(c, payload, id, &res)
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


func (t *TransactionControllerImpl) GetTransactionByAccount(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	var res []web.OfferWithAccount
	err := t.TransactionService.GetOfferByAccount(c, payload.UserId, &res)
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

func (t *TransactionControllerImpl) GetMyTransaction(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	var res []web.OfferWithAccount
	err := t.TransactionService.GetMyTransaction(c, payload.UserId, &res)
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

func (t *TransactionControllerImpl) UpdatePriceOffer(c *gin.Context) {
	payload := c.MustGet("payload").(helper.Payload)

	var req web.TransactionUpdateRequest

	if ok := helper.BindData(c, t.Translator, &req); !ok {
		return
	}

	req.BuyerId = payload.UserId

	err := t.TransactionService.UpdatePriceOffer(c, req)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Transaction id: %s successfully updated", req.TransactionId),
	})
}

func (t *TransactionControllerImpl) UpdateTransactionStatus(c *gin.Context) {

	payload := c.MustGet("payload").(helper.Payload)

	var param web.GetByIdRequest
	if err := c.ShouldBindUri(&param); err != nil {
		return
	}

	transactionId, err := uuid.Parse(param.ID)
	if err != nil {
		err := helper.NewBadRequest("invalid uuid")

		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	var res bool
	err = t.TransactionService.UpdateStatus(c, transactionId, payload.UserId, &res)
	if err != nil {
		c.JSON(helper.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":"Success",
		"Message":fmt.Sprintf("Transaction id: %s status set to: %v", transactionId, res),
	})
}
