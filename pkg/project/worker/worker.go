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
		OnProjectCreated(ctx *appcontext.AppContext, payload domain.QueueOnProjectCreatedPayload) error
		OnFeedbackCreated(ctx *appcontext.AppContext, payload domain.QueueOnFeedbackCreatedPayload) error
	}

	Instance interface {
		Handlers
	}

	workerHandlers struct {
		OnProjectCreatedHandler
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
	projectCollaboratorRepository domain.ProjectCollaboratorRepository,
	cachingRepository domain.CachingRepository,
) Worker {
	return Worker{
		queue: queue,
		workerHandlers: workerHandlers{
			OnProjectCreatedHandler:  NewOnProjectCreatedHandler(projectCollaboratorRepository),
			OnFeedbackCreatedHandler: NewOnFeedbackCreatedHandler(projectRepository, projectCampaignRepository, cachingRepository),
		},
	}
}

func (w Worker) Start() {
	server := w.queue.GetServer()

	// immediately
	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.OnProjectCreated), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueOnProjectCreatedPayload](bgCtx, t, queue.ParsePayload[domain.QueueOnProjectCreatedPayload], w.OnProjectCreated)
	})

	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.OnFeedbackCreated), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueOnFeedbackCreatedPayload](bgCtx, t, queue.ParsePayload[domain.QueueOnFeedbackCreatedPayload], w.OnFeedbackCreated)
	})
}
