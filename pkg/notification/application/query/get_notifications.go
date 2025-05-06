package query

import (
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type GetNotificationsHandler struct {
	notificationRepository domain.NotificationRepository
}

func NewGetNotificationsHandler(notificationRepository domain.NotificationRepository) GetNotificationsHandler {
	return GetNotificationsHandler{
		notificationRepository: notificationRepository,
	}
}

// GetNotifications godoc
// @tags     Notification
// @summary  Get notifications
// @id       notification-get-list
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    Accept-Language  header  string  true "en|vi"
// @param    payload query    dto.GetNotificationsRequest true "Query"
// @success  200     {object} dto.GetNotificationResponse
// @router   /api/notification [get]
func (h GetNotificationsHandler) GetNotifications(ctx *appcontext.AppContext, performerID string, req dto.GetNotificationsRequest) (*dto.GetNotificationResponse, error) {
	ctx.Logger().Info("new get notifications request", appcontext.Fields{
		"performerID": performerID, "page": req.Page,
	})

	ctx.Logger().Text("create new filter")
	filter, err := domain.NewNotificationFilter(performerID, req.Page)
	if err != nil {
		ctx.Logger().Error("failed to create new filter", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("find notifications in db")
	notifications, err := h.notificationRepository.FindWithFilter(ctx, *filter)
	if err != nil {
		ctx.Logger().Error("failed to find notifications in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("convert response data")
	var (
		result = make([]dto.Notification, 0)
		lang   = ctx.GetLang().String()
	)
	for _, notification := range notifications {
		var content = h.notificationRepository.GenerateNewFeedbackContent(ctx, lang, notification.Metadata.ProjectTitle)
		var action, params = domain.GenerateActionDataFromNotificationType(notification.Type, notification.Metadata)
		result = append(result, dto.Notification{}.FromDomain(notification, content, action, params))
	}

	ctx.Logger().Text("done get notifications request")
	return &dto.GetNotificationResponse{
		Notifications: result,
		Limit:         filter.Limit,
	}, nil
}
