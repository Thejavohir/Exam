package v1

import (
	"github.com/Exam/api_gateway/config"
	"github.com/Exam/api_gateway/pkg/logger"
	"github.com/Exam/api_gateway/services"
	"github.com/Exam/api_gateway/storage/repo"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redis          repo.NewRepo
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.NewRepo
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redis:          c.Redis,
	}
}
