package monitoring

import (
	"github.com/labstack/echo/v4"
)

type Operations interface {
}

type Monitoring struct {
}

type OtelConfig struct {
	Endpoint   string
	StreamName string
	Token      string
}

type SentryConfig struct {
	Dsn         string
	MachineName string
}

func NewMonitoringClient(e *echo.Echo, otelCfg OtelConfig, sentryCfg SentryConfig, appName, env string) *Monitoring {
	initOtel(otelCfg, appName, env)
	initSentry(e, sentryCfg, env)
	return &Monitoring{}
}
