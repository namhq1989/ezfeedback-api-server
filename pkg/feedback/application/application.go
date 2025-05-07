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
		ChangeFeedbackState(ctx *appcontext.AppContext, performerID, feedbackID string, req dto.ChangeFeedbackStateRequest) (*dto.ChangeFeedbackStateResponse, error)
	}
	Queries interface {
		Ping(ctx *appcontext.AppContext, _ dto.PingRequest) (*dto.PingResponse, error)

		GetFeedbacks(ctx *appcontext.AppContext, performerID string, req dto.GetFeedbacksRequest) (*dto.GetFeedbacksResponse, error)
		CountFeedbacks(ctx *appcontext.AppContext, performerID string, req dto.CountFeedbacksRequest) (*dto.CountFeedbacksResponse, error)
	}
	Instance interface {
		Commands
		Queries
	}

	commandHandlers struct {
		command.CreateFeedbackHandler
		command.ChangeFeedbackStateHandler
	}
	queryHandlers struct {
		query.PingHandler

		query.GetFeedbacksHandler
		query.CountFeedbacksHandler
	}
	Application struct {
		commandHandlers
		queryHandlers
	}
)

var _ Instance = (*Application)(nil)

func New(
	feedbackRepository domain.FeedbackRepository,
	cachingRepository domain.CachingRepository,
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
			ChangeFeedbackStateHandler: command.NewChangeFeedbackStateHandler(
				feedbackRepository,
				cachingRepository,
				service,
			),
		},
		queryHandlers: queryHandlers{
			PingHandler: query.NewPingHandler(),

			GetFeedbacksHandler:   query.NewGetFeedbacksHandler(feedbackRepository, projectHub),
			CountFeedbacksHandler: query.NewCountFeedbacksHandler(feedbackRepository),
		},
	}
}
