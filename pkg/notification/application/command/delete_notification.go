package command

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/dto"
	"github.com/namhq1989/go-utilities/appcontext"
)

type DeleteNotificationHandler struct {
	notificationRepository domain.NotificationRepository
}

func NewDeleteNotificationHandler(notificationRepository domain.NotificationRepository) DeleteNotificationHandler {
	return DeleteNotificationHandler{
		notificationRepository: notificationRepository,
	}
}

// DeleteNotification godoc
// @tags     Notification
// @summary  Delete notification
// @id       notification-delete
// @security ApiKeyAuth
// @accept   json
// @produce  json
// @param    id  	 path     string true "Notification id"
// @param    payload body    dto.DeleteNotificationRequest true "Body"
// @success  200     {object} dto.DeleteNotificationResponse
// @router   /api/notification/{id} [delete]
func (h DeleteNotificationHandler) DeleteNotification(ctx *appcontext.AppContext, performerID, notificationID string, _ dto.DeleteNotificationRequest) (*dto.DeleteNotificationResponse, error) {
	ctx.Logger().Info("new delete notification request", appcontext.Fields{
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

	ctx.Logger().Text("delete notification in db")
	if err = h.notificationRepository.Delete(ctx, *notification); err != nil {
		ctx.Logger().Error("failed to delete notification in db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done delete notification request")
	return &dto.DeleteNotificationResponse{}, nil
}
