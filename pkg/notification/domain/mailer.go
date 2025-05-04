package domain

import (
	"github.com/namhq1989/go-utilities/appcontext"
)

type MailerRepository interface {
	SendNewFeedbackEmail(ctx *appcontext.AppContext, toEmail string, feedbacks []Feedback) error
}
