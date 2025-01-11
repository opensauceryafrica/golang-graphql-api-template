package config

import (
	"fmt"
	"os"
	"time"

	"cendit.io/garage/logger"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	AppSrvName = "cendit" // service name

	//Headers

	HeaderRequestSource = "X-Request-Source"
	HeaderRequestID     = "X-Request-ID"

	// HTTP Client

	HTTPClientTimeout             = 10 * time.Second
	HTTPClientMaxIdleConns        = 100
	HTTPClientMaxIdleConnsPerHost = 100
)

type EmailConfig struct {
	ClientID         string `env:"MAIL_CLIENT_ID,required,notEmpty,unset"`
	Secret           string `env:"MAIL_CLIENT_SECRET,required,notEmpty,unset"`
	SendPulseBaseURL string `env:"SEND_PULSE_BASE_URL,required,notEmpty,unset"`
	Name             string `env:"EMAIL_NAME" envDefault:"Cendit"`
	From             string `env:"EMAIL_ADDRESS" envDefault:"support@cendit.io"`
}

type SmsConfig struct {
	SendChampKey     string `env:"SENDCHAMP_KEY,required,notEmpty,unset"`
	SendChampSender  string `env:"SENDCHAMP_SENDER,required,notEmpty,unset" envDefault:"cendit"`
	SendChampBaseURL string `env:"SENDCHAMP_BASE_URL,required,notEmpty,unset"`
}

type Variable struct {
	ServiceName             string
	Port                    int     `env:"PORT"`
	Environment             string  `env:"ENVIRONMENT" envDefault:"development"`
	CenditDatabaseURL       string  `env:"CENDIT_DATABASE_URL,required,notEmpty,unset"`
	RedisURL                string  `env:"REDIS_URL,required,notEmpty,unset"`
	SentryDSN               string  `env:"SENTRY_DSN,required,notEmpty,unset"`
	SentryDebug             bool    `env:"SENTRY_DEBUG" envDefault:"false"`
	SentrySampleRate        float64 `env:"SENTRY_SAMPLE_RATE" envDefault:"0.1"`
	DebugDatabase           bool    `env:"DEBUG_DATABASE" envDefault:"true"`
	DatabaseConnectionLimit int     `env:"DATABASE_CONNECTION_LIMIT" envDefault:"10"`
	TOTPSecret              string  `env:"TOTP_SECRET,required" envDefault:""`
	AppSecret               string  `env:"APP_SECRET,required" envDefault:""`
	Email                   EmailConfig
	SMS                     SmsConfig
}

var (
	// Variable is the global configuration object
	Env *Variable
)

// New initializes loads the environment from the .env file
// and parses them to the Variable struct returning a pointer to Variable
func Environment() *Variable {

	if Env != nil {
		return Env
	}

	envPath := ".env"
	if os.Getenv("ENV_PATH") != "" {
		envPath = os.Getenv("ENV_PATH")
	}

	if loadErr := godotenv.Load(envPath); loadErr != nil {
		logger.GetLogger().Sugar().Warnf("Unable to loan .env file: %v. Ignore if in production environment", loadErr.Error())
	}

	Env = &Variable{}

	if err := env.Parse(Env); err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("Failed to parse environment variables: %v", err.Error()))
	}
	Env.ServiceName = AppSrvName

	return Env
}
