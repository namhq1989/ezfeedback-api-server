package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type VerifyVerificationCodeHandler struct {
	userRepository             domain.UserRepository
	verificationCodeRepository domain.VerificationCodeRepository
	jwtRepository              domain.JwtRepository
}

func NewVerifyVerificationCodeHandler(
	userRepository domain.UserRepository,
	verificationCodeRepository domain.VerificationCodeRepository,
	jwtRepository domain.JwtRepository,
) VerifyVerificationCodeHandler {
	return VerifyVerificationCodeHandler{
		userRepository:             userRepository,
		verificationCodeRepository: verificationCodeRepository,
		jwtRepository:              jwtRepository,
	}
}

// VerifyVerificationCode godoc
// @tags     IAM
// @summary  Verify verification code
// @id       user-verify-verification-code
// @accept   json
// @produce  json
// @param    payload body    dto.VerifyVerificationCodeRequest true "Body"
// @success  200     {object} dto.VerifyVerificationCodeResponse
// @router   /api/user/verify-verification-code [post]
func (h VerifyVerificationCodeHandler) VerifyVerificationCode(ctx *appcontext.AppContext, ip string, req dto.VerifyVerificationCodeRequest) (*dto.VerifyVerificationCodeResponse, error) {
	ctx.Logger().Info("new verify verification code request", appcontext.Fields{"ip": ip, "email": req.Email, "code": req.Code})

	ctx.Logger().Text("find verification record in db")
	verificationCode, err := h.verificationCodeRepository.Find(ctx, ip, req.Email, req.Code)
	if err != nil {
		ctx.Logger().Error("failed to find verification record in db", err, appcontext.Fields{})
		return nil, err
	}
	if verificationCode == nil {
		ctx.Logger().ErrorText("verification record not found")
		return nil, apperrors.User.InvalidVerificationCode
	}

	ctx.Logger().Info("verify code", appcontext.Fields{
		"isUsed":    verificationCode.IsUsed,
		"code":      verificationCode.Code,
		"expiresAt": verificationCode.ExpiresAt.String(),
	})
	if verificationCode.IsUsed || !verificationCode.IsValidCode(req.Code) || verificationCode.IsExpired() {
		ctx.Logger().ErrorText("code not the same or expired")
		return nil, apperrors.User.InvalidVerificationCode
	}

	defer func() {
		ctx.Logger().Text("set used and update verification record in db")
		verificationCode.MarkAsUsed()
		_ = h.verificationCodeRepository.Update(ctx, *verificationCode)
	}()

	ctx.Logger().Text("find user by email")
	user, err := h.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		ctx.Logger().Error("failed to find user by email", err, appcontext.Fields{})
		return nil, err
	}
	if user != nil {
		ctx.Logger().Text("user found, generate jwt token and return")
		accessToken, gErr := h.jwtRepository.GenerateAccessToken(ctx, user.ID)
		if gErr != nil {
			ctx.Logger().Error("failed to generate access token", err, appcontext.Fields{})
			return nil, apperrors.Common.BadRequest
		}

		ctx.Logger().Text("done verify verification code request")
		return &dto.VerifyVerificationCodeResponse{
			IsNewUser: false,
			Token:     accessToken,
		}, nil
	}

	ctx.Logger().Text("user not found, create new user")
	user, err = domain.NewUser(req.Email, req.Email)
	if err != nil {
		ctx.Logger().Error("failed to create new user", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("persist user in db")
	if err = h.userRepository.Create(ctx, *user); err != nil {
		ctx.Logger().Error("failed to persist user in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("generate jwt token and return")
	accessToken, err := h.jwtRepository.GenerateAccessToken(ctx, user.ID)
	if err != nil {
		ctx.Logger().Error("failed to generate access token", err, appcontext.Fields{})
		return nil, apperrors.Common.BadRequest
	}

	ctx.Logger().Text("done verify verification code request")
	return &dto.VerifyVerificationCodeResponse{
		IsNewUser: true,
		Token:     accessToken,
	}, nil
}
