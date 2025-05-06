package grpc

import (
	"context"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/feedbackpb"
	"github.com/namhq1989/go-utilities/appcontext"
	"google.golang.org/grpc"
)

type server struct {
	hub Application
	feedbackpb.UnimplementedFeedbackServiceServer
}

var _ feedbackpb.FeedbackServiceServer = (*server)(nil)

func RegisterServer(_ *appcontext.AppContext, registrar grpc.ServiceRegistrar, hub *Application) error {
	feedbackpb.RegisterFeedbackServiceServer(registrar, server{hub: *hub})
	return nil
}

func (s server) GetProjectStatsForNotificationReminder(ctx context.Context, req *feedbackpb.GetProjectStatsForNotificationReminderRequest) (*feedbackpb.GetProjectStatsForNotificationReminderResponse, error) {
	return s.hub.GetProjectStatsForNotificationReminder(appcontext.NewGRPC(ctx), req)
}
