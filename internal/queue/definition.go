package queue

var TypeNames = struct {
	CleanupStaleVerificationCodes string
	SendVerificationCodeEmail     string

	OnFeedbackCreated string

	FeedbackCreated string

	CleanupStaleNotifications   string
	ScanNotificationReminders   string
	ProcessNotificationReminder string
}{
	CleanupStaleVerificationCodes: "iam.cleanupStaleVerificationCodes",
	SendVerificationCodeEmail:     "iam.sendVerificationCodeEmail",

	OnFeedbackCreated: "project.onFeedbackCreated",

	FeedbackCreated: "feedback.created",

	CleanupStaleNotifications:   "notification.cleanupStaleNotifications",
	ScanNotificationReminders:   "notification.scanNotificationReminders",
	ProcessNotificationReminder: "notification.processNotificationReminder",
}
