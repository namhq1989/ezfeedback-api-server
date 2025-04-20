package application

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Commands interface {
		RequestVerificationCode(ctx *appcontext.AppContext, ip string, req dto.RequestVerificationCodeRequest) (*dto.RequestVerificationCodeResponse, error)
	}
	Queries interface {
		Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error)
	}
	Instance interface {
		Commands
		Queries
	}

	commandHandlers struct {
		command.RequestVerificationCodeHandler
	}
	queryHandlers struct {
		query.PingHandler
	}
	Application struct {
		commandHandlers
		queryHandlers
	}
)

var _ Instance = (*Application)(nil)

func New(
	verificationCodeRepository domain.VerificationCodeRepository,
	queueRepository domain.QueueRepository,
) *Application {
	return &Application{
		commandHandlers: commandHandlers{
			RequestVerificationCodeHandler: command.NewRequestVerificationCodeHandler(
				verificationCodeRepository,
				queueRepository,
			),
		},
		queryHandlers: queryHandlers{
			PingHandler: query.NewPingHandler(),
		},
	}
}
