package query

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetUserProjectNotificationSettingHandler struct {
	service domain.Service
}

func NewGetUserProjectNotificationSettingHandler(service domain.Service) GetUserProjectNotificationSettingHandler {
	return GetUserProjectNotificationSettingHandler{
		service: service,
	}
}

// GetUserProjectNotificationSetting godoc
// @tags     Notification
// @summary  Get user project notification setting
// @id       notification-get-user-project-notification-setting
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload query    dto.GetUserProjectNotificationSettingRequest true "Query"
// @success  200     {object} dto.GetUserProjectNotificationSettingResponse
// @router   /api/notification/project-setting [get]
func (h GetUserProjectNotificationSettingHandler) GetUserProjectNotificationSetting(ctx *appcontext.AppContext, performerID string, req dto.GetUserProjectNotificationSettingRequest) (*dto.GetUserProjectNotificationSettingResponse, error) {
	ctx.Logger().Info("new get user project notification setting request", appcontext.Fields{
		"performerID": performerID, "projectID": req.ProjectID,
	})

	ctx.Logger().Text("find setting in db")
	setting, err := h.service.GetUserProjectNotificationSetting(ctx, performerID, req.ProjectID)
	if err != nil {
		ctx.Logger().Error("failed to find setting in db", err, appcontext.Fields{})
		return nil, err
	}
	if setting == nil {
		ctx.Logger().ErrorText("setting not found")
		return nil, apperrors.Notification.UserProjectNotificationSettingNotFound
	}

	ctx.Logger().Text("convert to response")
	result := dto.GetUserProjectNotificationSettingResponse{
		Setting: dto.UserProjectNotificationSetting{}.FromDomain(*setting),
	}

	ctx.Logger().Text("done get user project notification setting request")
	return &result, nil
}
