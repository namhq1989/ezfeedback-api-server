package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type RequestVerificationCodeHandler struct {
	verificationCodeRepository domain.VerificationCodeRepository
	queueRepository            domain.QueueRepository
}

func NewRequestVerificationCodeHandler(
	verificationCodeRepository domain.VerificationCodeRepository,
	queueRepository domain.QueueRepository,
) RequestVerificationCodeHandler {
	return RequestVerificationCodeHandler{
		verificationCodeRepository: verificationCodeRepository,
		queueRepository:            queueRepository,
	}
}

// RequestVerificationCode godoc
// @tags     IAM
// @summary  Request verification code
// @id       iam-request-verification-code
// @accept   json
// @produce  json
// @param    payload body    dto.RequestVerificationCodeRequest true "Body"
// @success  200     {object} dto.RequestVerificationCodeResponse
// @router   /api/iam/request-verification-code [post]
func (h RequestVerificationCodeHandler) RequestVerificationCode(ctx *appcontext.AppContext, ip string, req dto.RequestVerificationCodeRequest) (*dto.RequestVerificationCodeResponse, error) {
	ctx.Logger().Info("new request verification code request", appcontext.Fields{"ip": ip, "email": req.Email})

	ctx.Logger().Text("count total sent today by ip")
	totalSent, err := h.verificationCodeRepository.CountTotalSentTodayByIp(ctx, ip)
	if err != nil {
		ctx.Logger().Error("failed to count total sent today by ip", err, appcontext.Fields{})
		return nil, err
	}
	if domain.IsDailyIpOtpLimitExceeded(totalSent) {
		ctx.Logger().ErrorText("total sent today by ip exceeded")
		return nil, apperrors.User.DailyIpOtpLimitExceeded
	} else {
		ctx.Logger().Info("total sent today by ip not exceeded", appcontext.Fields{"totalSent": totalSent})
	}

	ctx.Logger().Text("create new verification code")
	code, err := domain.NewVerificationCode(ip, req.Email)
	if err != nil {
		ctx.Logger().Error("failed to create new verification code", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("persist code in db")
	if err = h.verificationCodeRepository.Create(ctx, *code); err != nil {
		ctx.Logger().Error("failed to persist verification code", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("enqueue send email task")
	if err = h.queueRepository.SendVerificationCodeEmail(ctx, domain.QueueSendVerificationCodeEmailPayload{
		Email: req.Email,
		Code:  code.Code,
	}); err != nil {
		ctx.Logger().Error("failed to enqueue send email task", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done request verification code request")
	return &dto.RequestVerificationCodeResponse{}, nil
}
