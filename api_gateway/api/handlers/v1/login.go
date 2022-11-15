package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/Exam/api_gateway/api/handlers/models"
	"github.com/Exam/api_gateway/genproto/customer"
	"github.com/Exam/api_gateway/pkg/etc"
	"github.com/Exam/api_gateway/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// customer login
// @Summary 		Login Customer
// @Description 	login customer
// @Tags 			Customer
// @Accept 			json
// @Produce			json
// @Param 			email 		path string true "email"
// @Param 			password 	path string true "password"
// @Success 		200 {object} 	customer.LoginResp
// @Failure			500 {object} 	models.Error
// @Failure			400 {object} 	models.Error
// @Router			/v1/login/{email}/{password} [get]
func (h *handlerV1) Login(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions

	jspbMarshal.UseEnumNumbers = true
	var (
		password = c.Param("password")
		email    = c.Param("email")
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	loginResp, err := h.serviceManager.CustomerService().Login(ctx, &customer.LoginReq{
		Email: email,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Error:       err,
			Description: "Couln't find matching information, Have you registered before?",
		})
		h.log.Error("Error while logging in by email", logger.Any("post", err))
		return
	}

	if !etc.CheckPasswordHash(password, loginResp.Password) {
		c.JSON(http.StatusNotFound, models.Error{
			Description: "Wrong email or password",
			Code:        http.StatusBadRequest,
		})
		return
	}

	h.jwthandler.Iss = "user"
	h.jwthandler.Sub = loginResp.Id
	h.jwthandler.Role = "authorized"
	h.jwthandler.Aud = []string{"exam-app"}
	h.jwthandler.SigninKey = h.cfg.SignInKey
	h.jwthandler.Log = h.log
	tokens, err := h.jwthandler.GenerateAuthJWT()
	accessToken := tokens[0]
	refreshToken := tokens[1]

	if err != nil {
		h.log.Error("error occured while generating tokens")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}
	loginResp.AccessToken = accessToken
	loginResp.RefreshToken = refreshToken
	loginResp.Password = ""
	c.JSON(http.StatusOK, loginResp)
}
