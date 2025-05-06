package query

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CountNotificationsHandler struct {
	notificationRepository domain.NotificationRepository
}

func NewCountNotificationsHandler(notificationRepository domain.NotificationRepository) CountNotificationsHandler {
	return CountNotificationsHandler{
		notificationRepository: notificationRepository,
	}
}

// CountNotifications godoc
// @tags     Notification
// @summary  Count notifications
// @id       notification-count
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    payload query    dto.CountNotificationsRequest true "Query"
// @success  200     {object} dto.CountNotificationResponse
// @router   /api/notification/count [get]
func (h CountNotificationsHandler) CountNotifications(ctx *appcontext.AppContext, performerID string, _ dto.CountNotificationsRequest) (*dto.CountNotificationResponse, error) {
	ctx.Logger().Info("new count notifications request", appcontext.Fields{"performerID": performerID})

	ctx.Logger().Text("create new filter")
	filter, err := domain.NewNotificationFilter(performerID, 0)
	if err != nil {
		ctx.Logger().Error("failed to create new filter", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("count notifications in db")
	total, err := h.notificationRepository.CountWithFilter(ctx, *filter)
	if err != nil {
		ctx.Logger().Error("failed to count notifications in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done get notifications request")
	return &dto.CountNotificationResponse{
		Total: total,
	}, nil
}
