package v1

import (
	j "github.com/Exam/api_gateway/api/token"
	"github.com/Exam/api_gateway/config"
	"github.com/Exam/api_gateway/pkg/logger"
	"github.com/Exam/api_gateway/services"
	"github.com/Exam/api_gateway/storage/repo"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redis          repo.NewRepo
	jwthandler     j.JWTHandler
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.NewRepo
	JWTHandler     j.JWTHandler
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redis:          c.Redis,
		jwthandler:     c.JWTHandler,
	}
}

func GetClaims(h handlerV1, c *gin.Context) (*j.CustomClaims, error) {

	var (
		claims = j.CustomClaims{}
	)

	strToken := c.GetHeader("Authorization")

	token, err := jwt.Parse(strToken, func(t *jwt.Token) (interface{}, error) { return []byte(h.cfg.SignInKey), nil })

	if err != nil {
		h.log.Error("invalid access token")
		return nil, err
	}
	rawClaims := token.Claims.(jwt.MapClaims)

	claims.Sub = rawClaims["sub"].(string)
	claims.Exp = rawClaims["exp"].(float64)
	// fmt.Printf("%T type of value in map %v", rawClaims["exp"], rawClaims["exp"])
	// fmt.Printf("%T type of value in map %v", rawClaims["iat"], rawClaims["iat"])
	aud := cast.ToStringSlice(rawClaims["aud"])
	claims.Aud = aud
	claims.Role = rawClaims["role"].(string)
	claims.Sub = rawClaims["sub"].(string)
	claims.Token = token
	return &claims, nil
}
