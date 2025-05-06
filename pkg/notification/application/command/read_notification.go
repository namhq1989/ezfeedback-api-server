package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type ReadNotificationHandler struct {
	notificationRepository domain.NotificationRepository
}

func NewReadNotificationHandler(notificationRepository domain.NotificationRepository) ReadNotificationHandler {
	return ReadNotificationHandler{
		notificationRepository: notificationRepository,
	}
}

// ReadNotification godoc
// @tags     Notification
// @summary  Read notification
// @id       notification-read
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Notification id"
// @param    payload body    dto.ReadNotificationRequest true "Body"
// @success  200     {object} dto.ReadNotificationResponse
// @router   /api/notification/{id}/read [patch]
func (h ReadNotificationHandler) ReadNotification(ctx *appcontext.AppContext, performerID, notificationID string, _ dto.ReadNotificationRequest) (*dto.ReadNotificationResponse, error) {
	ctx.Logger().Info("new read notification request", appcontext.Fields{
		"performerID": performerID, "notificationID": notificationID,
	})

	ctx.Logger().Text("find notification in db")
	notification, err := h.notificationRepository.FindByID(ctx, notificationID)
	if err != nil {
		ctx.Logger().Error("failed to find notification in db", err, appcontext.Fields{})
		return nil, err
	}
	if notification == nil {
		ctx.Logger().Text("notification not found")
		return nil, apperrors.Notification.NotificationNotFound
	}
	if !notification.IsOwner(performerID) {
		ctx.Logger().ErrorText("user is not notification owner")
		return nil, apperrors.Notification.NotificationNotFound
	}

	if notification.IsRead {
		ctx.Logger().Text("notification already read, skip")
		return &dto.ReadNotificationResponse{}, nil
	}

	ctx.Logger().Text("mark notification as read")
	notification.MarkAsRead()

	ctx.Logger().Text("update notification in db")
	if err = h.notificationRepository.Update(ctx, *notification); err != nil {
		ctx.Logger().Error("failed to update notification in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done read notification request")
	return &dto.ReadNotificationResponse{}, nil
}
