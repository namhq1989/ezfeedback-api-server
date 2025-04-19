package mailer

import (
	"strings"

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
	ctx.Logger().Info("new send email request", appcontext.Fields{"to": req.To, "subject": req.Subject})

	from := mail.NewEmail("", m.fromEmail)
	to := mail.NewEmail("", req.To)
	message := mail.NewSingleEmail(from, req.Subject, to, htmlToText(req.Content), req.Content)

	_, err := m.sendgrid.Send(message)
	if err != nil {
		ctx.Logger().Error("[mailer] error send email", err, appcontext.Fields{})
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
