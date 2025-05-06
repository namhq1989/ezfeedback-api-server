package domain

import "github.com/namhq1989/go-utilities/appcontext"

type QueueRepository interface {
	SendVerificationCodeEmail(ctx *appcontext.AppContext, payload QueueSendVerificationCodeEmailPayload) error
}

type QueueCleanupStaleVerificationCodesPayload struct{}

type QueueSendVerificationCodeEmailPayload struct {
	Email string
	Code  string
}
