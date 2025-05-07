package queue

var TypeNames = struct {
	CleanupStaleVerificationCodes string
	SendVerificationCodeEmail     string

	OnProjectCreated  string
	OnFeedbackCreated string

	FeedbackCreated string

	CleanupStaleNotifications   string
	ScanNotificationReminders   string
	ProcessNotificationReminder string
}{
	CleanupStaleVerificationCodes: "iam.cleanupStaleVerificationCodes",
	SendVerificationCodeEmail:     "iam.sendVerificationCodeEmail",

	OnProjectCreated:  "project.onProjectCreated",
	OnFeedbackCreated: "project.onFeedbackCreated",

	FeedbackCreated: "feedback.created",

	CleanupStaleNotifications:   "notification.cleanupStaleNotifications",
	ScanNotificationReminders:   "notification.scanNotificationReminders",
	ProcessNotificationReminder: "notification.processNotificationReminder",
}
