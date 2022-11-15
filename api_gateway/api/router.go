package api

import (
	_ "github.com/Exam/api_gateway/api/docs"
	v1 "github.com/Exam/api_gateway/api/handlers/v1"
	"github.com/Exam/api_gateway/config"
	"github.com/Exam/api_gateway/pkg/logger"
	"github.com/Exam/api_gateway/services"
	"github.com/Exam/api_gateway/storage/repo"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Redis          repo.NewRepo
}

// New ...
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.Redis,
	})

	api := router.Group("/v1")

	// customer apis
	api.POST("/customer", handlerV1.CreateCust)
	api.POST("/customer/register", handlerV1.Register)
	api.GET("/login/:email/:password", handlerV1.Login)
	api.GET("/customer/:id", handlerV1.GetCustById)
	api.PUT("/customer/update", handlerV1.UpdateCust)
	api.GET("/customer/allcustomers", handlerV1.ListCusts)
	api.DELETE("/customer/delete/:id", handlerV1.DeleteCust)
	api.GET("/verify/:email/:code", handlerV1.Verify)

	// post apis
	api.POST("/post", handlerV1.CreatePost)
	api.GET("/post/:id", handlerV1.GetPost)
	api.PUT("/posts", handlerV1.UpdatePost)
	api.GET("/post/allposts", handlerV1.ListPosts)
	api.DELETE("/posts/:id", handlerV1.DeletePost)

	// review apis
	api.POST("/review", handlerV1.CreateReview)
	api.GET("/review/:id", handlerV1.GetReview)
	api.PUT("/reviews", handlerV1.UpdateReview)
	api.DELETE("/reviews/:id", handlerV1.DeleteReview)
	api.GET("/review/post/:id", handlerV1.PostReview)
	// api.POST()

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
