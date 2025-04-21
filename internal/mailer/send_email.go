package mailer

import (
	"errors"
	"strings"

	brevo "github.com/getbrevo/brevo-go/lib"

	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/net/html"
)

type SendEmailRequest struct {
	To      string
	Subject string
	Content string
}

func (m Mailer) SendEmail(ctx *appcontext.AppContext, req SendEmailRequest) error {
	if m.service == ServiceSendgrid {
		return m.sendEmailWithSendgrid(ctx, req)
	} else if m.service == ServiceBrevo {
		return m.sendEmailWithBrevo(ctx, req)
	}

	ctx.Logger().Error("no valid mailer service", nil, appcontext.Fields{"service": m.service})
	return errors.New("no valid mailer service")
}

func (m Mailer) sendEmailWithSendgrid(ctx *appcontext.AppContext, req SendEmailRequest) error {
	ctx.Logger().Info("new send email with Sendgrid request", appcontext.Fields{"to": req.To, "subject": req.Subject})

	from := mail.NewEmail("", m.fromEmail)
	to := mail.NewEmail("", req.To)
	message := mail.NewSingleEmail(from, req.Subject, to, htmlToText(req.Content), req.Content)

	_, err := m.sendgrid.Send(message)
	if err != nil {
		ctx.Logger().Error("[mailer] error send Sendgrid email", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("mail sent successfully")
	return nil
}

func (m Mailer) sendEmailWithBrevo(ctx *appcontext.AppContext, req SendEmailRequest) error {
	ctx.Logger().Info("new send email with Brevo request", appcontext.Fields{"to": req.To, "subject": req.Subject})

	_, _, err := m.brevo.TransactionalEmailsApi.SendTransacEmail(ctx.Context(), brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  m.fromEmail,
			Email: m.fromEmail,
		},
		To:          []brevo.SendSmtpEmailTo{{Email: req.To}},
		HtmlContent: req.Content,
		TextContent: htmlToText(req.Content),
		Subject:     req.Subject,
	})
	if err != nil {
		ctx.Logger().Error("[mailer] error send Brevo email", err, appcontext.Fields{})
		return err
	}

	ctx.Logger().Text("mail sent successfully")
	return nil
}

func htmlToText(htmlString string) string {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return ""
	}

	var text string
	extractText(doc, &text)
	return strings.TrimSpace(text)
}

func extractText(n *html.Node, text *string) {
	if n.Type == html.TextNode {
		*text += n.Data + " "
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, text)
	}
}
