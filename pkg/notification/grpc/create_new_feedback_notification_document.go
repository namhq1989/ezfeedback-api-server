package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/notificationpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
)

type CreateNewFeedbackNotificationDocumentHandler struct {
	notificationRepository domain.NotificationRepository
}

func NewCreateNewFeedbackNotificationDocumentHandler(notificationRepository domain.NotificationRepository) CreateNewFeedbackNotificationDocumentHandler {
	return CreateNewFeedbackNotificationDocumentHandler{
		notificationRepository: notificationRepository,
	}
}

func (h CreateNewFeedbackNotificationDocumentHandler) CreateNewFeedbackNotificationDocument(ctx *appcontext.AppContext, req *notificationpb.CreateNewFeedbackNotificationDocumentRequest) (*notificationpb.CreateNewFeedbackNotificationDocumentResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new create new feedback notification document request", appcontext.Fields{
		"userID": req.GetUserId(), "metadata": req.GetMetadata(),
	})

	ctx.Logger().Text("new notification model")
	var metadata = req.GetMetadata()
	notification, err := domain.NewFeedbackNotification(req.GetUserId(), domain.NotificationMetadata{
		ProjectID:    metadata.GetProjectId(),
		ProjectTitle: metadata.GetProjectTitle(),
	})
	if err != nil {
		ctx.Logger().Error("failed to create new notification model", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("persist notification to db")
	if err = h.notificationRepository.Create(ctx, *notification); err != nil {
		ctx.Logger().Error("failed to persist notification to db", err, appcontext.Fields{})
		return nil, err
	}

	ctx.Logger().Text("done create notification document request")
	return &notificationpb.CreateNewFeedbackNotificationDocumentResponse{}, nil
}
