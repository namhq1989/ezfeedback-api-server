package monolith

import (
	"github.com/labstack/echo/v4"
	"github.com/namhq1989/ezfeedback-api-server/internal/caching"
	"github.com/namhq1989/ezfeedback-api-server/internal/config"
	"github.com/namhq1989/ezfeedback-api-server/internal/database"
	appjwt "github.com/namhq1989/ezfeedback-api-server/internal/jwt"
	"github.com/namhq1989/ezfeedback-api-server/internal/mailer"
	"github.com/namhq1989/ezfeedback-api-server/internal/monitoring"
	"github.com/namhq1989/ezfeedback-api-server/internal/queue"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/waiter"
	"github.com/namhq1989/go-utilities/appcontext"
	"google.golang.org/grpc"
)

type Monolith interface {
	Config() config.Server
	Database() *database.Database
	Caching() *caching.Caching
	JWT() *appjwt.JWT
	Queue() *queue.Queue
	Mailer() *mailer.Mailer
	Monitoring() *monitoring.Monitoring
	Rest() *echo.Echo
	RPC() *grpc.Server
	Waiter() waiter.Waiter
}

type Module interface {
	Name() string
	Startup(ctx *appcontext.AppContext, monolith Monolith) error
}
