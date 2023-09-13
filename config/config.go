package config

import (
	"blacheapi/logger"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	AppSrvName = "blache" // service name

	//Headers

	HeaderRequestSource = "X-Request-Source"
	HeaderRequestID     = "X-Request-ID"

	// HTTP Client

	HTTPClientTimeout             = 10 * time.Second
	HTTPClientMaxIdleConns        = 100
	HTTPClientMaxIdleConnsPerHost = 100
)

type Config struct {
	ServiceName             string
	Port                    int     `env:"PORT"`
	Environment             string  `env:"ENVIRONMENT" envDefault:"development"`
	BlacheDatabaseURL       string  `env:"SAVE_DATABASE_URL,required,notEmpty,unset"`
	RedisURL                string  `env:"REDIS_URL,required,notEmpty,unset"`
	SentryDSN               string  `env:"SENTRY_DSN,required,notEmpty,unset"`
	SentryDebug             bool    `env:"SENTRY_DEBUG" envDefault:"false"`
	SentrySampleRate        float64 `env:"SENTRY_SAMPLE_RATE" envDefault:"0.1"`
	DebugDatabase           bool    `env:"DEBUG_DATABASE" envDefault:"false"`
	DatabaseConnectionLimit int     `env:"DATABASE_CONNECTION_LIMIT" envDefault:"10"`
}

// New initializes loads the environment from the .env file
// and parses them to the Config struct returning a pointer to Config
func New() *Config {

	envPath := ".env"
	if os.Getenv("ENV_PATH") != "" {
		envPath = os.Getenv("ENV_PATH")
	}

	if loadErr := godotenv.Load(envPath); loadErr != nil {
		logger.GetLogger().Sugar().Warnf("Unable to loan .env file: %v. Ignore if in production environment", loadErr.Error())
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("Failed to parse environment variables: %v", err.Error()))
	}
	cfg.ServiceName = AppSrvName

	return &cfg
}
