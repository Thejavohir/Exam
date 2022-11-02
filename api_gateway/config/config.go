package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment string // develop, staging, production

	// context timeout in seconds
	CtxTimeout int

	CustomerServiceHost string
	CustomerServicePort int

	PostServiceHost string
	PostServicePort int

	ReviewServiceHost string
	ReviewServicePort int

	RedisHost string
	RedisPort int

	LogLevel string
	HTTPPort string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":9090"))
	
	c.CustomerServiceHost = cast.ToString(getOrReturnDefault("CUSTOMER_SERVICE_HOST", "localhost"))
	c.CustomerServicePort = cast.ToInt(getOrReturnDefault("CUSTOMER_SERVICE_PORT", 3000))

	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "localhost"))
	c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_PORT", 4000))

	c.ReviewServiceHost = cast.ToString(getOrReturnDefault("REVIEW_SERVICE_HOST", "localhost"))
	c.ReviewServicePort = cast.ToInt(getOrReturnDefault("REVIEW_SERVICE_PORT", 2000))
	
	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))


	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
