package queue

var TypeNames = struct {
	CleanupStaleVerificationCodes string
	SendVerificationCodeEmail     string

	FeedbackCreated string

	CleanupStaleNotifications   string
	ScanNotificationReminders   string
	ProcessNotificationReminder string
}{
	CleanupStaleVerificationCodes: "iam.cleanupStaleVerificationCodes",
	SendVerificationCodeEmail:     "iam.sendVerificationCodeEmail",

	FeedbackCreated: "feedback.created",

	CleanupStaleNotifications:   "notification.cleanupStaleNotifications",
	ScanNotificationReminders:   "notification.scanNotificationReminders",
	ProcessNotificationReminder: "notification.processNotificationReminder",
}
