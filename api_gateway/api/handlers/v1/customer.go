package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Exam/api_gateway/genproto/customer"
	pbp "github.com/Exam/api_gateway/genproto/customer"

	l "github.com/Exam/api_gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateCustomer
// @Summary creation of customers
// @Description creating customers
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body customer.CustomerReq true "Customer"
// @Success 200 {object} customer.CustomerResp
// @Failure 400 "ErrorResponse"
// @Router /v1/customer [post]
func (h *handlerV1) CreateCust(c *gin.Context) {
	var (
		body        customer.CustomerReq
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.CustomerService().CreateCust(ctx, &body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create customer", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetCustomerById
// @Summary gets the customer info
// @Description getting customer info
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} customer.Customer
// @Failure 400 "ErrorResponse"
// @Router /v1/customer/{id} [get]
func (h *handlerV1) GetCustById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseEnumNumbers = true

	guid := c.Param("id")
	Id, err := strconv.ParseInt(guid, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed parse string to int", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	responseCus, err := h.serviceManager.CustomerService().GetCustById(ctx, &pbp.GetCustByIdReq{
		Id: Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get customer", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, responseCus)
}

// CutomerUpdate
// @Summary updates customers
// @Description updating customers
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body customer.Customer true "Customer"
// @Success 200 {object} customer.Customer
// @Failure 400 "ErrorResponse"
// @Router /v1/customer/update [put]
func (h *handlerV1) UpdateCust(c *gin.Context) {
	var (
		customerBody customer.Customer
		jspbMarshal  protojson.MarshalOptions
	)
	jspbMarshal.UseEnumNumbers = true

	err := c.ShouldBindJSON(&customerBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.CustomerService().UpdateCust(ctx, &customerBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary get all customers
// @Description getting all customers
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} customer.ListCustsResp
// @Router /v1/customer/allcustomers [get]
func (h *handlerV1) ListCusts(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseEnumNumbers = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.CustomerService().ListCusts(ctx, &customer.Empty{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to collect all customers", l.Error(err))
	}
	c.JSON(http.StatusOK, response)
}

// DeleteCustomer
// @Summary deletes customer
// @Description deleting of customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 "success"
// @Failure 400 "ErrorResponse"
// @Router /v1/customer/delete/{id} [delete]
func (h *handlerV1) DeleteCust(c *gin.Context) {

	var (
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseEnumNumbers = true
	guid := c.Param("id")

	id, err := strconv.ParseInt(guid, 10, 64)
	body := &customer.Id{
		Id: id,
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed parse string to int", l.Error(err))
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	customerResp, err := h.serviceManager.CustomerService().DeleteCust(ctx, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
	}

	c.JSON(http.StatusCreated, customerResp)
}
