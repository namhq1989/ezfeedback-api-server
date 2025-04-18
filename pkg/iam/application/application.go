package application

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/common/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/common/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Queries interface {
		Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error)
	}
	Instance interface {
		Queries
	}

	queryHandlers struct {
		query.PingHandler
	}
	Application struct {
		queryHandlers
	}
)

var _ Instance = (*Application)(nil)

func New() *Application {
	return &Application{
		queryHandlers: queryHandlers{
			PingHandler: query.NewPingHandler(),
		},
	}
}
