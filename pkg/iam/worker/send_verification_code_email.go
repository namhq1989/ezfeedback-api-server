package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type SendSignInVerificationCodeEmailHandler struct {
	mailerRepository domain.MailerRepository
}

func NewSendSignInVerificationCodeEmailHandler(mailerRepository domain.MailerRepository) SendSignInVerificationCodeEmailHandler {
	return SendSignInVerificationCodeEmailHandler{
		mailerRepository: mailerRepository,
	}
}

func (w SendSignInVerificationCodeEmailHandler) SendSignInVerificationCodeEmail(ctx *appcontext.AppContext, payload domain.QueueSendVerificationCodeEmailPayload) error {
	tracer := otel.Tracer("[tracer] send verification code email")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] send verification code email")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Text("send email")
	return w.mailerRepository.SendVerificationCodeEmail(ctx, payload.Email, payload.Code)
}
