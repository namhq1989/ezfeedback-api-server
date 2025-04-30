package infrastructure

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/queue"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type QueueRepository struct {
	queue queue.Operations
}

func NewQueueRepository(queue queue.Operations) QueueRepository {
	return QueueRepository{
		queue: queue,
	}
}

func (r QueueRepository) FeedbackCreated(ctx *appcontext.AppContext, payload domain.QueueFeedbackCreatedPayload) error {
	return queue.EnqueueTask(ctx, r.queue, queue.TypeNames.FeedbackCreated, payload, 3)
}
