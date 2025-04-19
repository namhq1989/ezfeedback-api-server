package mailer

import (
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/sendgrid/sendgrid-go"
)

type Operations interface {
	SendEmail(ctx *appcontext.AppContext, req SendEmailRequest) error
}

type Mailer struct {
	sendgrid  *sendgrid.Client
	fromEmail string
}

func NewMailerClient(apiKey, fromEmail string) *Mailer {
	return &Mailer{
		sendgrid:  sendgrid.NewSendClient(apiKey),
		fromEmail: fromEmail,
	}
}
