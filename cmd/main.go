package main

import (
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namhq1989/ezfeedback-api-server/docs"
	"github.com/namhq1989/ezfeedback-api-server/internal/caching"
	"github.com/namhq1989/ezfeedback-api-server/internal/config"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	appjwt "github.com/namhq1989/ezfeedback-api-server/internal/jwt"
	"github.com/namhq1989/ezfeedback-api-server/internal/monitoring"
	"github.com/namhq1989/ezfeedback-api-server/internal/monolith"
	"github.com/namhq1989/ezfeedback-api-server/internal/queue"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/staticfiles"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/waiter"
	"github.com/namhq1989/ezfeedback-api-server/pkg/common"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam"
	"github.com/namhq1989/go-utilities/logger"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title       EzFeedback - App api
// @version     1.0
// @description Apis for EzFeedback app
// @termsOfService https://easyfeedback.com
// @contact.name   Nam
// @contact.url    https://easyfeedback.com
// @contact.email  namhq.1989@gmail.com
// @basePath       /

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization

func main() {
	var err error

	// config
	cfg := config.Init()

	// logger
	logger.Init(cfg.Environment)

	// app error
	apperrors.Init()

	// static files
	staticfiles.Init(cfg.CDNEndpoint)

	// server
	a := app{}
	a.cfg = cfg

	// jwt
	a.jwt, err = appjwt.Init(cfg.AccessTokenSecret, time.Second*time.Duration(cfg.AccessTokenTTL))
	if err != nil {
		panic(err)
	}

	// rest
	a.rest = initRest(cfg)

	// grpc
	a.rpc = initRPC()

	// database
	a.database = database.NewDatabaseClient(cfg.PostgresConn)

	// caching
	a.caching = caching.NewCachingClient(cfg.CachingRedisURL)

	// monitoring
	a.monitoring = monitoring.NewMonitoringClient(
		a.rest,
		monitoring.OtelConfig{
			Endpoint:   cfg.OpenObserveHttpEndpoint,
			StreamName: cfg.OpenObserveStreamName,
			Token:      cfg.OpenObserveToken,
		},
		monitoring.SentryConfig{
			Dsn:         cfg.SentryDSN,
			MachineName: cfg.SentryMachineName,
		},
		cfg.AppName,
		cfg.Environment,
	)

	// queue
	a.queue = queue.Init(cfg.QueueRedisURL, cfg.QueueConcurrency)

	// init queue's dashboard
	a.rest.Any(fmt.Sprintf("%s/*", queue.DashboardPath), echo.WrapHandler(queue.EnableDashboard(cfg.QueueRedisURL)), middleware.BasicAuth(func(username, password string, _ echo.Context) (bool, error) {
		if !cfg.IsEnvRelease {
			return true, nil
		}
		return subtle.ConstantTimeCompare([]byte(username), []byte(cfg.QueueUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(cfg.QueuePassword)) == 1, nil
	}))

	// waiter
	a.waiter = waiter.New(waiter.CatchSignals())

	// swagger
	if !cfg.IsEnvRelease {
		docs.SwaggerInfo.Host = cfg.SwaggerURL
		a.rest.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// modules
	a.modules = []monolith.Module{
		&common.Module{},
		&iam.Module{},
	}

	// start
	if err = a.startupModules(); err != nil {
		panic(err)
	}

	fmt.Println("--- started ezfeedback-api-server application")
	defer fmt.Println("--- stopped ezfeedback-api-server application")

	// wait for other service starts
	a.waiter.Add(
		a.waitForRest,
		a.waitForRPC,
	)
	if err = a.waiter.Wait(); err != nil {
		panic(err)
	}
}
