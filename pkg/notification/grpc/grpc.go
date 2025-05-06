package grpc

import (
	"context"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/notificationpb"
	"github.com/namhq1989/go-utilities/appcontext"
	"google.golang.org/grpc"
)

type server struct {
	hub Application
	notificationpb.UnimplementedNotificationServiceServer
}

var _ notificationpb.NotificationServiceServer = (*server)(nil)

func RegisterServer(_ *appcontext.AppContext, registrar grpc.ServiceRegistrar, hub *Application) error {
	notificationpb.RegisterNotificationServiceServer(registrar, server{hub: *hub})
	return nil
}

func (s server) CreateNewFeedbackNotificationDocument(bgCtx context.Context, req *notificationpb.CreateNewFeedbackNotificationDocumentRequest) (*notificationpb.CreateNewFeedbackNotificationDocumentResponse, error) {
	return s.hub.CreateNewFeedbackNotificationDocument(appcontext.NewGRPC(bgCtx), req)
}

func (s server) CreateNotificationReminder(bgCtx context.Context, req *notificationpb.CreateNotificationReminderRequest) (*notificationpb.CreateNotificationReminderResponse, error) {
	return s.hub.CreateNotificationReminder(appcontext.NewGRPC(bgCtx), req)
}
