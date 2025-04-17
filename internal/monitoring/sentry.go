package monitoring

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

func initSentry(e *echo.Echo, cfg SentryConfig, envName string) {
	// skip if the "machine" is not set
	if cfg.MachineName == "" {
		fmt.Println("⚡️ [monitoring]: machine is not set, skipping sentry")
		return
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.Dsn,
		Environment: fmt.Sprintf("%s-%s", envName, cfg.MachineName),
		// Debug: true,
	}); err != nil {
		panic(err)
	}

	e.Use(sentryecho.New(sentryecho.Options{}))

	// recover
	defer sentry.Recover()

	fmt.Printf("⚡️ [monitoring]: sentry initialized \n")
}
