package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GenerateTokenHandler struct {
	jwtRepository domain.JwtRepository
}

func NewGenerateTokenHandler(jwtRepository domain.JwtRepository) GenerateTokenHandler {
	return GenerateTokenHandler{
		jwtRepository: jwtRepository,
	}
}

// GenerateToken godoc
// @tags     IAM
// @summary  [Dev only] Generate platform token
// @id       iam-generate-token
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload body    dto.GenerateTokenRequest true "Body"
// @success  200     {object} dto.GenerateTokenResponse
// @router   /api/iam/generate-token [post]
func (h GenerateTokenHandler) GenerateToken(ctx *appcontext.AppContext, req dto.GenerateTokenRequest) (*dto.GenerateTokenResponse, error) {
	ctx.Logger().Info("new generate token request", appcontext.Fields{"userID": req.UserID})

	ctx.Logger().Text("generate access token")
	accessToken, err := h.jwtRepository.GenerateAccessToken(ctx, req.UserID)
	if err != nil {
		ctx.Logger().Error("failed to generate access token", err, appcontext.Fields{})
		return nil, apperrors.Common.BadRequest
	}

	ctx.Logger().Text("done generate token request")
	return &dto.GenerateTokenResponse{Token: accessToken}, nil
}
