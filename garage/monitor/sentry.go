package monitor

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

// InitSentry initialize a sentry instance on startup
func InitSentry(sentryDSN string, sentryDebug bool, sentrySampleRate float64, environment string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDSN,
		Debug:            sentryDebug,
		TracesSampleRate: sentrySampleRate,
		AttachStacktrace: true,
		Environment:      fmt.Sprintf("ams-%s", environment),
	})
	if err != nil {
		logger, _ := zap.NewProduction()
		logger.Sugar().Fatalf("[Sentry]: failed to initialize sentry: %v", err.Error())
	}
	defer sentry.Flush(2 * time.Second)
}

// SendScopeLocalizedError contains the monitor context as well as the emitting the event preventing mixing
// error context. Additional context for the error can be store in the errorContext variable
func SendScopeLocalizedError(err error, errorContext map[string]interface{}, userEmail string, ID int32, sentryLevel sentry.Level) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			Email: userEmail,
			ID:    fmt.Sprintf("%v", ID),
		})
		scope.SetContext("additional_information", errorContext)
		scope.SetLevel(sentryLevel)
		sentry.CaptureException(err)
	})
}
