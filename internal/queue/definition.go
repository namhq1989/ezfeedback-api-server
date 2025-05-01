package queue

var TypeNames = struct {
	DeleteExpiredVerificationCodes string
	SendVerificationCodeEmail      string

	FeedbackCreated string
}{
	DeleteExpiredVerificationCodes: "iam.deleteExpiredVerificationCodes",
	SendVerificationCodeEmail:      "iam.sendVerificationCodeEmail",

	FeedbackCreated: "feedback.created",
}
