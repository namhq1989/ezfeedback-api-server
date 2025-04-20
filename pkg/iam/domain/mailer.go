package domain

import "github.com/namhq1989/go-utilities/appcontext"

type MailerRepository interface {
	SendVerificationCodeEmail(ctx *appcontext.AppContext, toEmail, code string) error
}
