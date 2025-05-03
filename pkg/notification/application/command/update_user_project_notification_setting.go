package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type UpdateUserProjectNotificationSettingHandler struct {
	userProjectNotificationSettingRepository domain.UserProjectNotificationSettingRepository
	cachingRepository                        domain.CachingRepository
	service                                  domain.Service
}

func NewUpdateUserProjectNotificationSettingHandler(
	userProjectNotificationSettingRepository domain.UserProjectNotificationSettingRepository,
	cachingRepository domain.CachingRepository,
	service domain.Service,
) UpdateUserProjectNotificationSettingHandler {
	return UpdateUserProjectNotificationSettingHandler{
		userProjectNotificationSettingRepository: userProjectNotificationSettingRepository,
		cachingRepository:                        cachingRepository,
		service:                                  service,
	}
}

// UpdateUserProjectNotificationSetting godoc
// @tags     Notification
// @summary  Update user project notification setting
// @id       notification-update-user-project-notification-setting
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload body    dto.UpdateUserProjectNotificationSettingRequest true "Body"
// @success  200     {object} dto.UpdateUserProjectNotificationSettingResponse
// @router   /api/notification/project-setting [put]
func (h UpdateUserProjectNotificationSettingHandler) UpdateUserProjectNotificationSetting(ctx *appcontext.AppContext, performerID string, req dto.UpdateUserProjectNotificationSettingRequest) (*dto.UpdateUserProjectNotificationSettingResponse, error) {
	ctx.Logger().Info("new update user project notification setting request", appcontext.Fields{
		"performerID": performerID, "projectID": req.ProjectID,
		"receiveNewFeedback": req.ReceiveNewFeedback, "receiveDailySummary": req.ReceiveDailySummary, "receiveWeeklySummary": req.ReceiveWeeklySummary,
	})

	ctx.Logger().Text("find user project notification setting in db")
	setting, err := h.service.GetUserProjectNotificationSetting(ctx, performerID, req.ProjectID)
	if err != nil {
		ctx.Logger().Error("failed to find user project notification setting in db", err, appcontext.Fields{})
		return nil, err
	}
	if setting == nil {
		ctx.Logger().ErrorText("user project notification setting not found")
		return nil, apperrors.Notification.UserProjectNotificationSettingNotFound
	}

	ctx.Logger().Text("update user project notification setting data")
	setting.SetReceiveNewFeedback(req.ReceiveNewFeedback)
	setting.SetReceiveDailySummary(req.ReceiveDailySummary)
	setting.SetReceiveWeeklySummary(req.ReceiveWeeklySummary)

	ctx.Logger().Text("update user project notification setting in db")
	if err = h.userProjectNotificationSettingRepository.Update(ctx, *setting); err != nil {
		ctx.Logger().Error("failed to update user project notification setting in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("update user project notification setting in caching")
	if err = h.cachingRepository.SetUserProjectNotificationSetting(ctx, performerID, req.ProjectID, *setting); err != nil {
		ctx.Logger().Error("failed to update user project notification setting in caching", err, appcontext.Fields{})
	}

	ctx.Logger().Text("done update user project notification setting")
	return &dto.UpdateUserProjectNotificationSettingResponse{}, nil
}
