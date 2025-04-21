package query

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetMeHandler struct {
	service domain.Service
}

func NewGetMeHandler(service domain.Service) GetMeHandler {
	return GetMeHandler{
		service: service,
	}
}

// GetMe godoc
// @tags     IAM
// @summary  Get me
// @id       iam-get-me
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload query    dto.GetMeRequest true "Query"
// @success  200     {object} dto.GetMeResponse
// @router   /api/iam/me [get]
func (h GetMeHandler) GetMe(ctx *appcontext.AppContext, performerID string, _ dto.GetMeRequest) (*dto.GetMeResponse, error) {
	ctx.Logger().Info("new get me request", appcontext.Fields{"performerID": performerID})

	ctx.Logger().Text("get user by id")
	user, err := h.service.GetUserByID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to get user by id", err, appcontext.Fields{})
		return nil, err
	}
	if user == nil {
		ctx.Logger().ErrorText("user not found")
		return nil, apperrors.User.UserNotFound
	}

	ctx.Logger().Text("get me response")
	result := dto.User{}.FromDomain(*user)

	ctx.Logger().Text("done get me request")
	return &dto.GetMeResponse{
		Me: result,
	}, nil
}
