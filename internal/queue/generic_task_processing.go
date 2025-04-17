package queue

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/namhq1989/go-utilities/appcontext"
)

func ParsePayload[T any](ctx *appcontext.AppContext, t *asynq.Task) (T, error) {
	var payload T

	ctx.Logger().Info("unmarshal task payload for queue", appcontext.Fields{"type": t.Type()})
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		ctx.Logger().Error("failed to unmarshal task payload", err, appcontext.Fields{})
		return payload, err
	}

	return payload, nil
}

type ProcessorHandler[T any] func(*appcontext.AppContext, T) error

func ProcessTask[T any](bgCtx context.Context, t *asynq.Task, parse func(*appcontext.AppContext, *asynq.Task) (T, error), process ProcessorHandler[T]) error {
	ctx := appcontext.NewWorker(bgCtx)
	ctx.Logger().Info("[worker] process new task", appcontext.Fields{"type": t.Type(), "payload": string(t.Payload())})

	ctx.Logger().Text("unmarshal task payload")
	payload, err := parse(ctx, t)
	if err != nil {
		ctx.Logger().Error("failed to unmarshal task payload", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("process task")
	if err = process(ctx, payload); err != nil {
		ctx.Logger().Error("failed to process task", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Info("[worker] done task", appcontext.Fields{"type": t.Type()})
	return nil
}

func EnqueueTask[T any](ctx *appcontext.AppContext, q Operations, typename string, payload T, retryTimes int) error {
	n := q.GenerateTypename(typename)
	t, err := q.RunTask(n, payload, retryTimes)
	if err != nil {
		ctx.Logger().Error("failed to enqueue task", err, appcontext.Fields{"typename": typename, "payload": payload})
		return err
	}

	ctx.Logger().Info("enqueued task", appcontext.Fields{"taskId": t.ID, "typename": typename})
	return nil
}
