package queue

var TypeNames = struct {
	DeleteExpiredVerificationCodes string
	SendVerificationCodeEmail      string
}{
	DeleteExpiredVerificationCodes: "iam.deleteExpiredVerificationCodes",
	SendVerificationCodeEmail:      "iam.sendVerificationCodeEmail",
}
