package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Exam/api_gateway/api/handlers/models"
	emailS "github.com/Exam/api_gateway/email"
	"github.com/Exam/api_gateway/genproto/customer"
	"github.com/Exam/api_gateway/pkg/etc"
	"github.com/Exam/api_gateway/pkg/logger"
	"github.com/Exam/api_gateway/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/encoding/protojson"
)

// register customer
// @Summary		registeration customer
// @Description registers customer
// @Tags		Customer
// @Accept		json
// @Produce 	json
// @Param 		body	body  models.CustomerRegister true "Register customer"
// @Success		200 "success"
// @Router		/v1/customer/register 	[post]
func (h *handlerV1) Register(c *gin.Context) {
	newCust := &models.CustomerRegister{}
	c.ShouldBindJSON(newCust)

	email, err := utils.IsValidMail(newCust.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email address",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	exists, err := h.serviceManager.CustomerService().CheckField(ctx, &customer.CheckFieldRequest{
		Field: "email",
		Value: newCust.Email,
	})

	if err != nil {
		h.log.Error("failed checking email existanse", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while checking email existance",
		})
		return
	}

	if exists.Exists {
		c.JSON(http.StatusOK, gin.H{
			"message": "such email is already exists",
		})
		return
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(newCust.Password), 6)
	if err != nil {
		h.log.Error("error while hashing password", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	newCust.Password = string(hashPass)
	code := etc.GenerateCode(6)
	customerData := models.CustomerData{
		FirstName: newCust.FirstName,
		LastName:  newCust.LastName,
		Bio:       newCust.Bio,
		Password:  newCust.Password,
		Email:     newCust.Email,
		Code:      code,
	}
	_, err = h.redis.Get(fmt.Sprint(customerData.Email))
	if err != nil {
		h.log.Error("error marshaling new customer", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while creating customer",
		})
		return
	}
	marshNewCust, err := json.Marshal(customerData)
	if err != nil {
		h.log.Error("error while marshaling new user, inorder to insert it to redis", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while creating user",
		})
		return
	}
	fmt.Println(code)

	if err = h.redis.SetWithTTL(fmt.Sprint(customerData.Email), string(marshNewCust), 86000); err != nil {
		h.log.Error("error while inserting new user into redis")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong, please try again",
		})
		return
	}

	customerData.Email = email
	customerData.Code = code
	msg := "Exam email verification\nVerification code: " + customerData.Code
	errs := emailS.SendMail([]string{customerData.Email}, []byte(msg))

	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "email is invalid",
		})
		h.log.Error("failed to send email")
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Code is successfully sent to your email!",
	})
}

// Verify customer
// @Summary      Verify customer
// @Description  Verifies customer
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        email  path string true "email"
// @Param        code   path string true "code"
// @Success      200  {object}  models.VerifyResponse
// @Router      /v1/verify/{email}/{code} [get]
func (h *handlerV1) Verify(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	var (
		code  = c.Param("code")
		email = c.Param("email")
	)

	customerBody, err := h.redis.Get(email)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"info":  "Your time has expired",
			"error": err.Error(),
		})
		h.log.Error("Error while getting customer from redis", logger.Any("redis", err))
	}
	customerBodys := cast.ToString(customerBody)
	body := customer.CustomerReq{}
	err = json.Unmarshal([]byte(customerBodys), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while unmarshaling from json to customer body", logger.Any("json", err))
		return
	}
	if body.Code != code {
		fmt.Println(body.Code)
		c.JSON(http.StatusConflict, gin.H{
			"info": "Wrong code",
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	// // Genrating refresh and jwt tokens
	// h.jwthandler.Iss = "user"
	// h.jwthandler.Sub = body.Id
	// h.jwthandler.Role = "authorized"
	// h.jwthandler.Aud = []string{"exam-app"}
	// h.jwthandler.SigninKey = h.cfg.SignInKey
	// h.jwthandler.Log = h.log
	// tokens, err := h.jwthandler.GenerateAuthJWT()
	// accessToken := tokens[0]
	// refreshToken := tokens[1]

	// if err != nil {
	// 	h.log.Error("error occured while generating tokens")
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "something went wrong,please try again",
	// 	})
	// 	return
	// }
	// body.RefreshToken = refreshToken
	res, err := h.serviceManager.CustomerService().CreateCust(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while creating customer", logger.Any("post", err))
		return
	}
	response := &models.VerifyResponse{
		Id:          res.Id,
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		Email:       res.Email,
		Bio:         res.Bio,
		PhoneNumber: res.PhoneNumber,
	}
	// response.JWT = accessToken
	// response.RefreshToken = refreshToken

	c.JSON(http.StatusOK, response)
}
