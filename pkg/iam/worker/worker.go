package worker

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/namhq1989/ezfeedback-api-server/internal/queue"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Handlers interface {
		SendSignInVerificationCodeEmail(ctx *appcontext.AppContext, payload domain.QueueSendVerificationCodeEmailPayload) error
	}
	Cronjob interface {
		DeleteExpiredVerificationCodes(ctx *appcontext.AppContext, _ domain.QueueDeleteExpiredVerificationCodesPayload) error
	}

	Instance interface {
		Handlers
		Cronjob
	}

	workerHandlers struct {
		SendSignInVerificationCodeEmailHandler
	}
	cronjobHandlers struct {
		DeleteExpiredVerificationCodesHandler
	}
	Worker struct {
		queue queue.Operations
		workerHandlers
		cronjobHandlers
	}
)

var _ Instance = (*Worker)(nil)

func New(
	queue queue.Operations,
	verificationCodeRepository domain.VerificationCodeRepository,
	mailerRepository domain.MailerRepository,
) Worker {
	return Worker{
		queue: queue,
		workerHandlers: workerHandlers{
			SendSignInVerificationCodeEmailHandler: NewSendSignInVerificationCodeEmailHandler(mailerRepository),
		},
		cronjobHandlers: cronjobHandlers{
			DeleteExpiredVerificationCodesHandler: NewDeleteExpiredVerificationCodesHandler(verificationCodeRepository),
		},
	}
}

func (w Worker) Start() {
	w.addCronjob()

	server := w.queue.GetServer()

	// immediately
	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.SendVerificationCodeEmail), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueSendVerificationCodeEmailPayload](bgCtx, t, queue.ParsePayload[domain.QueueSendVerificationCodeEmailPayload], w.SendSignInVerificationCodeEmail)
	})

	// cronjob
	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.DeleteExpiredVerificationCodes), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueDeleteExpiredVerificationCodesPayload](bgCtx, t, queue.ParsePayload[domain.QueueDeleteExpiredVerificationCodesPayload], w.DeleteExpiredVerificationCodes)
	})
}

type cronjobData struct {
	Task       string      `json:"task"`
	CronSpec   string      `json:"cronSpec"`
	Payload    interface{} `json:"payload"`
	RetryTimes int         `json:"retryTimes"`
}

func (w Worker) addCronjob() {
	var (
		ctx  = appcontext.NewWorker(context.Background())
		jobs = []cronjobData{
			{
				Task:       w.queue.GenerateTypename(queue.TypeNames.DeleteExpiredVerificationCodes),
				CronSpec:   "0 */1 * * *", // every day
				Payload:    domain.QueueDeleteExpiredVerificationCodesPayload{},
				RetryTimes: 3,
			},
		}
	)

	for _, job := range jobs {
		entryID, err := w.queue.ScheduleTask(job.Task, job.Payload, job.CronSpec, job.RetryTimes)
		if err != nil {
			ctx.Logger().Error("error when initializing cronjob", err, appcontext.Fields{"job": job})
			panic(err)
		}

		ctx.Logger().Info(fmt.Sprintf("[cronjob] cronjob '%s' initialize successfully with cronSpec '%s' and retryTimes '%d'", job.Task, job.CronSpec, job.RetryTimes), appcontext.Fields{
			"entryId": entryID,
		})
	}
}
