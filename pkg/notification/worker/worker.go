package worker

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/namhq1989/ezfeedback-api-server/internal/queue"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type (
	Handlers interface {
		ProcessNotificationReminder(ctx *appcontext.AppContext, payload domain.QueueProcessNotificationReminderPayload) error
	}
	Cronjob interface {
		CleanupStaleNotifications(ctx *appcontext.AppContext, _ domain.QueueCleanupStaleNotificationsPayload) error
		ScanNotificationReminders(ctx *appcontext.AppContext, _ domain.QueueScanNotificationRemindersPayload) error
	}

	Instance interface {
		Handlers
		Cronjob
	}

	workerHandlers struct {
		ProcessNotificationReminderHandler
	}
	cronjobHandlers struct {
		CleanupStaleNotificationsHandler
		ScanNotificationRemindersHandler
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
	notificationRepository domain.NotificationRepository,
	notificationReminderRepository domain.NotificationReminderRepository,
	queueRepository domain.QueueRepository,
	mailerRepository domain.MailerRepository,
	feedbackHub domain.FeedbackHub,
	projectHub domain.ProjectHub,
) Worker {
	return Worker{
		queue: queue,
		workerHandlers: workerHandlers{
			ProcessNotificationReminderHandler: NewProcessNotificationReminderHandler(notificationReminderRepository, mailerRepository, feedbackHub, projectHub),
		},
		cronjobHandlers: cronjobHandlers{
			CleanupStaleNotificationsHandler: NewCleanupStaleNotificationsHandler(notificationRepository),
			ScanNotificationRemindersHandler: NewScanNotificationRemindersHandler(notificationReminderRepository, queueRepository),
		},
	}
}

func (w Worker) Start() {
	w.addCronjob()

	server := w.queue.GetServer()

	// immediately
	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.ProcessNotificationReminder), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueProcessNotificationReminderPayload](bgCtx, t, queue.ParsePayload[domain.QueueProcessNotificationReminderPayload], w.ProcessNotificationReminder)
	})

	// cronjob
	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.CleanupStaleNotifications), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueCleanupStaleNotificationsPayload](bgCtx, t, queue.ParsePayload[domain.QueueCleanupStaleNotificationsPayload], w.CleanupStaleNotifications)
	})
	server.HandleFunc(w.queue.GenerateTypename(queue.TypeNames.ScanNotificationReminders), func(bgCtx context.Context, t *asynq.Task) error {
		return queue.ProcessTask[domain.QueueScanNotificationRemindersPayload](bgCtx, t, queue.ParsePayload[domain.QueueScanNotificationRemindersPayload], w.ScanNotificationReminders)
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
				Task:       w.queue.GenerateTypename(queue.TypeNames.CleanupStaleNotifications),
				CronSpec:   "0 */1 * * *", // every day
				Payload:    domain.QueueCleanupStaleNotificationsPayload{},
				RetryTimes: 3,
			},
			{
				Task: w.queue.GenerateTypename(queue.TypeNames.ScanNotificationReminders),
				// CronSpec:   "*/10 * * * *", // every 10m
				CronSpec:   "0 */5 * * *", // every 1m
				Payload:    domain.QueueScanNotificationRemindersPayload{},
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
