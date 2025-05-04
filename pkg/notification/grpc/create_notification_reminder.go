package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/notificationpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CreateNotificationReminderHandler struct {
	notificationReminderRepository domain.NotificationReminderRepository
}

func NewCreateNotificationReminderHandler(notificationReminderRepository domain.NotificationReminderRepository) CreateNotificationReminderHandler {
	return CreateNotificationReminderHandler{
		notificationReminderRepository: notificationReminderRepository,
	}
}

func (h CreateNotificationReminderHandler) CreateNotificationReminder(ctx *appcontext.AppContext, req *notificationpb.CreateNotificationReminderRequest) (*notificationpb.CreateNotificationReminderResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new create notification reminder request", appcontext.Fields{"projectID": req.GetProjectId()})

	ctx.Logger().Text("find reminder in db")
	reminder, err := h.notificationReminderRepository.FindByProjectID(ctx, req.GetProjectId())
	if err != nil {
		ctx.Logger().Error("failed to find reminder in db", err, appcontext.Fields{})
		return nil, err
	}
	if reminder != nil {
		ctx.Logger().Text("reminder found in db, skip")
		return &notificationpb.CreateNotificationReminderResponse{}, nil
	}

	ctx.Logger().Text("reminder not found, create new model")
	reminder, err = domain.NewNotificationReminder(req.GetProjectId())
	if err != nil {
		ctx.Logger().Error("failed to create new model", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("persist reminder to db")
	if err = h.notificationReminderRepository.Create(ctx, *reminder); err != nil {
		ctx.Logger().Error("failed to persist reminder to db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done create notification reminder request")
	return &notificationpb.CreateNotificationReminderResponse{}, nil
}
