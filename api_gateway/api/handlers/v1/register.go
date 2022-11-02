package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Exam/api_gateway/api/handlers/models"
	"github.com/Exam/api_gateway/email"
	"github.com/Exam/api_gateway/genproto/customer"
	"github.com/Exam/api_gateway/pkg/etc"
	"github.com/Exam/api_gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// register customer
// @Summary		registeration customer
// @Description registers customer
// @Tags		Customer
// @Accept		json
// @Produce 	json
// @Param 		body	body  models.CustomerRegister true "Register customer"
// @Success		201 	{object} models.Error
// @Failure		500 	{object} models.Error
// @Router		/v1/customer/register 	[post]
func (h *handlerV1) Register(c *gin.Context) {
	var body customer.CustomerReq

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: err,
		})
		h.log.Error("Error while binding json", logger.Any("json", err))
		return
	}
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	body.Password, err = etc.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error: err,
		})
		h.log.Error("failed to hash the password")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	emailExists, err := h.serviceManager.CustomerService().CheckField(ctx, &customer.CheckFieldRequest{
		Field: "email",
		Value: body.Email,
	})

	fmt.Println(emailExists, "#######", body.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:       err,
			Description: "email not available",
		})
		h.log.Error("Error while cheking email uniqeness", logger.Any("check", err))
		return
	}

	if emailExists.Exists {
		c.JSON(http.StatusConflict, models.Error{
			Error: err,
		})
		return
	}
	exists, err := h.redis.Exists(body.Email)
	if err != nil {
		h.log.Error("Error while checking email from redis")
		c.JSON(http.StatusConflict, models.Error{
			Error: err,
		})
		return
	}
	if emailExists.Exists {
		c.JSON(http.StatusConflict, models.Error{
			Error: err,
		})
		return
	}
	if cast.ToInt(exists) == 1 {
		c.JSON(http.StatusConflict, models.Error{
			Error: err,
		})
		return
	}

	customerToSaved := &customer.CustomerReq{
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Bio:         body.Bio,
		Email:       body.Email,
		Password:    body.Password,
		PhoneNumber: body.PhoneNumber,
	}
	for _, address := range body.Adresses {
		customerToSaved.Adresses = append(customerToSaved.Adresses, &customer.AddressReq{
			Street: address.Street,
		})
	}
	customerToSaved.Code = etc.GenerateCode(6)
	msg := "Subject: Customer email verification\n Your verification code: " + customerToSaved.Code
	err = email.SendMail([]string{body.Email}, []byte(msg))

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:       nil,
			Code:        http.StatusAccepted,
			Description: "Your Email is not valid, Please recheck it",
		})
		return
	}
	c.JSON(http.StatusAccepted, models.Error{
		Error:       nil,
		Code:        http.StatusAccepted,
		Description: "Your request successfuly accepted we have send code to your email, Your code is : " + customerToSaved.Code,
	})

	bodyByte, err := json.Marshal(customerToSaved)
	if err != nil {
		h.log.Error("Error while marshaling to json", logger.Any("json", err))
		return
	}
	err = h.redis.SetWithTTL(customerToSaved.Email, string(bodyByte), 600)
	if err != nil {
		h.log.Error("Error while marshaling to json", logger.Any("json", err))
		return
	}
}
