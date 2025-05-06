package domain

import "github.com/namhq1989/go-utilities/appcontext"

type QueueRepository interface {
	ProcessNotificationReminder(ctx *appcontext.AppContext, payload QueueProcessNotificationReminderPayload) error
}

type QueueScanNotificationRemindersPayload struct{}

type QueueProcessNotificationReminderPayload struct {
	Reminder NotificationReminder
}
