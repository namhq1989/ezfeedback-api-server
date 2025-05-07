package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/namhq1989/ezfeedback-api-server/internal/queue"
	"github.com/namhq1989/ezfeedback-api-server/pkg/project/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Handlers interface {
		OnFeedbackCreated(ctx *appcontext.AppContext, payload domain.QueueOnFeedbackCreatedPayload) error
	}

	Instance interface {
		Handlers
	}

	workerHandlers struct {
		OnFeedbackCreatedHandler
	}
	Worker struct {
		queue queue.Operations
		workerHandlers
	}
)

var _ Instance = (*Worker)(nil)

func New(
	queue queue.Operations,
	projectRepository domain.ProjectRepository,
	projectCampaignRepository domain.ProjectCampaignRepository,
	cachingRepository domain.CachingRepository,
) Worker {
	return Worker{
		queue: queue,
		workerHandlers: workerHandlers{
			OnFeedbackCreatedHandler: NewOnFeedbackCreatedHandler(projectRepository, projectCampaignRepository, cachingRepository),
		},
	}
}

func (w Worker) Start() {
	server := w.queue.GetServer()

	// immediately
	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.OnFeedbackCreated), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueOnFeedbackCreatedPayload](bgCtx, t, queue.ParsePayload[domain.QueueOnFeedbackCreatedPayload], w.OnFeedbackCreated)
	})
}
