package infrastructure

import (
	_ "embed"
	"strings"

	"github.com/namhq1989/ezfeedback-api-server/internal/mailer"
	"github.com/namhq1989/go-utilities/appcontext"
)

//go:embed template/verification_code_email.html
var verificationCodeEmailTemplate []byte

type MailerRepository struct {
	mailer                        mailer.Operations
	verificationCodeEmailSubject  string
	verificationCodeEmailTemplate string
}

func NewMailerRepository(mailer mailer.Operations) MailerRepository {
	var m = MailerRepository{
		mailer: mailer,

		verificationCodeEmailSubject:  "[EasyFeedback] Sign in verification code",
		verificationCodeEmailTemplate: string(verificationCodeEmailTemplate),
	}

	return m
}

func (r MailerRepository) SendVerificationCodeEmail(ctx *appcontext.AppContext, toEmail, code string) error {
	var (
		subject = r.verificationCodeEmailSubject
		content = r.verificationCodeEmailTemplate
	)

	content = strings.Replace(content, "{{code}}", code, 1)
	content = strings.Replace(content, "{{toEmail}}", toEmail, 1)

	return r.mailer.SendEmail(ctx, mailer.SendEmailRequest{
		To:      toEmail,
		Subject: subject,
		Content: content,
	})
}
