package queue

var TypeNames = struct {
	DeleteExpiredVerificationCodes string
	SendVerificationCodeEmail      string

	FeedbackCreated string

	ScanNotificationReminders   string
	ProcessNotificationReminder string
}{
	DeleteExpiredVerificationCodes: "iam.deleteExpiredVerificationCodes",
	SendVerificationCodeEmail:      "iam.sendVerificationCodeEmail",

	FeedbackCreated: "feedback.created",

	ScanNotificationReminders:   "notification.scanNotificationReminders",
	ProcessNotificationReminder: "notification.processNotificationReminder",
}
