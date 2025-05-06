package worker

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"go.opentelemetry.io/otel"
)

type CleanupStaleVerificationCodesHandler struct {
	verificationCodeRepository domain.VerificationCodeRepository
}

func NewCleanupStaleVerificationCodesHandler(verificationCodeRepository domain.VerificationCodeRepository) CleanupStaleVerificationCodesHandler {
	return CleanupStaleVerificationCodesHandler{
		verificationCodeRepository: verificationCodeRepository,
	}
}

func (h CleanupStaleVerificationCodesHandler) CleanupStaleVerificationCodes(ctx *appcontext.AppContext, _ domain.QueueCleanupStaleVerificationCodesPayload) error {
	tracer := otel.Tracer("[tracer] cleanup stale verification codes")
	spanCtx, span := tracer.Start(ctx.Context(), "[worker] cleanup stale verification codes")
	ctx.SetContext(spanCtx)
	defer span.End()

	ctx.Logger().Text("delete in db")
	return h.verificationCodeRepository.CleanupStale(ctx)
}
