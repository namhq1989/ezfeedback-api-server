package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type DeleteExpiredVerificationCodesHandler struct {
	verificationCodeRepository domain.VerificationCodeRepository
}

func NewDeleteExpiredVerificationCodesHandler(verificationCodeRepository domain.VerificationCodeRepository) DeleteExpiredVerificationCodesHandler {
	return DeleteExpiredVerificationCodesHandler{
		verificationCodeRepository: verificationCodeRepository,
	}
}

func (h DeleteExpiredVerificationCodesHandler) DeleteExpiredVerificationCodes(ctx *appcontext.AppContext, _ domain.QueueDeleteExpiredVerificationCodesPayload) error {
	tracer := otel.Tracer("[tracer] delete expired verification codes")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] delete expired verification codes")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Text("delete in db")
	return h.verificationCodeRepository.DeleteExpired(ctx)
}
