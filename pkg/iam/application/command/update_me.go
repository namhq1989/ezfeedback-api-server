package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/iam/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type UpdateMeHandler struct {
	userRepository    domain.UserRepository
	cachingRepository domain.CachingRepository
	service           domain.Service
}

func NewUpdateMeHandler(userRepository domain.UserRepository, cachingRepository domain.CachingRepository, service domain.Service) UpdateMeHandler {
	return UpdateMeHandler{
		userRepository:    userRepository,
		cachingRepository: cachingRepository,
		service:           service,
	}
}

// UpdateMe godoc
// @tags     IAM
// @summary  Update me
// @id       iam-update-me
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload body    dto.UpdateMeRequest true "Body"
// @success  200     {object} dto.UpdateMeResponse
// @router   /api/iam/me [put]
func (h UpdateMeHandler) UpdateMe(ctx *appcontext.AppContext, performerID string, req dto.UpdateMeRequest) (*dto.UpdateMeResponse, error) {
	ctx.Logger().Info("new update me request", appcontext.Fields{"performerID": performerID, "name": req.Name})

	ctx.Logger().Text("find user in db")
	user, err := h.service.GetUserByID(ctx, performerID)
	if err != nil {
		ctx.Logger().Error("failed to find user in db", err, appcontext.Fields{})
		return nil, err
	}
	if user == nil {
		ctx.Logger().ErrorText("user not found")
		return nil, apperrors.User.UserNotFound
	}

	ctx.Logger().Text("update user data")
	if err = user.SetName(req.Name); err != nil {
		ctx.Logger().Error("failed to update user data", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update user in db")
	if err = h.userRepository.Update(ctx, *user); err != nil {
		ctx.Logger().Error("failed to update user in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("delete user in caching")
	if err = h.cachingRepository.DeleteUserByID(ctx, performerID); err != nil {
		ctx.Logger().Error("failed to delete user in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done update me request")
	return &dto.UpdateMeResponse{}, nil
}
