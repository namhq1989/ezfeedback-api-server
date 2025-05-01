package application

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/application/command"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/application/query"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Commands interface {
		CreateFeedback(ctx *appcontext.AppContext, ip, domain string, req dto.CreateFeedbackRequest) (*dto.CreateFeedbackResponse, error)
	}
	Queries interface {
		Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error)
	}
	Instance interface {
		Commands
		Queries
	}

	commandHandlers struct {
		command.CreateFeedbackHandler
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
	feedbackRepository domain.FeedbackRepository,
	queueRepository domain.QueueRepository,
	billingHub domain.BillingHub,
	projectHub domain.ProjectHub,
	service domain.Service,
) *Application {
	return &Application{
		commandHandlers: commandHandlers{
			CreateFeedbackHandler: command.NewCreateFeedbackHandler(
				feedbackRepository,
				queueRepository,
				billingHub,
				projectHub,
				service,
			),
		},
		queryHandlers: queryHandlers{
			PingHandler: query.NewPingHandler(),
		},
	}
}
