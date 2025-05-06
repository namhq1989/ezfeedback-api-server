package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/namhq1989/ezfeedback-api-server/internal/queue"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Handlers interface {
		FeedbackCreated(ctx *appcontext.AppContext, payload domain.QueueFeedbackCreatedPayload) error
	}

	Instance interface {
		Handlers
	}

	workerHandlers struct {
		FeedbackCreatedHandler
	}
	Worker struct {
		queue queue.Operations
		workerHandlers
	}
)

var _ Instance = (*Worker)(nil)

func New(
	queue queue.Operations,
	projectHub domain.ProjectHub,
	notificationHub domain.NotificationHub,
) Worker {
	return Worker{
		queue: queue,
		workerHandlers: workerHandlers{
			FeedbackCreatedHandler: NewFeedbackCreatedHandler(projectHub, notificationHub),
		},
	}
}

func (w Worker) Start() {
	server := w.queue.GetServer()

	// immediately
	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.FeedbackCreated), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueFeedbackCreatedPayload](bgCtx, t, queue.ParsePayload[domain.QueueFeedbackCreatedPayload], w.FeedbackCreated)
	})
}
