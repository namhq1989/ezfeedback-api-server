package mailer

import (
	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/sendgrid/sendgrid-go"
)

type Operations interface {
	SendEmail(ctx *appcontext.AppContext, req SendEmailRequest) error
}

const (
	ServiceSendgrid = "sendgrid"
	ServiceBrevo    = "brevo"
)

type Mailer struct {
	service   string
	sendgrid  *sendgrid.Client
	brevo     *brevo.APIClient
	fromEmail string
}

func NewMailerClient(service, sendgridApiKey, brevoApiKey, fromEmail string) *Mailer {
	// brevo
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", brevoApiKey)
	cfg.AddDefaultHeader("partner-key", brevoApiKey)
	br := brevo.NewAPIClient(cfg)

	return &Mailer{
		service:   service,
		sendgrid:  sendgrid.NewSendClient(sendgridApiKey),
		brevo:     br,
		fromEmail: fromEmail,
	}
}
