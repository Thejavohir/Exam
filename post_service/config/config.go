package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment         string // develop, staging, production
	PostgresHost        string
	PostgresPort        int
	PostgresDatabase    string
	PostgresUser        string
	PostgresPassword    string
	LogLevel            string
	RPCPort             string
	ReviewServiceHost   string
	ReviewServicePort   int
	CustomerServiceHost string
	CustomerServicePort int
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "post"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1234"))

	c.ReviewServiceHost = cast.ToString(getOrReturnDefault("REVIEW_SERVICE_HOST", "review"))
	c.ReviewServicePort = cast.ToInt(getOrReturnDefault("REVIEW_SERVICE_PORT", 2000))

	c.CustomerServiceHost = cast.ToString(getOrReturnDefault("CUSTOMER_SERVICE_HOST", "customer"))
	c.CustomerServicePort = cast.ToInt(getOrReturnDefault("CUSTOMER_SERVICE_PORT", 3000))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))

	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":4000"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
