package infrastructure

import (
	"time"

	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/feedbackpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/notification/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FeedbackHub struct {
	client feedbackpb.FeedbackServiceClient
}

func NewFeedbackHub(client feedbackpb.FeedbackServiceClient) FeedbackHub {
	return FeedbackHub{
		client: client,
	}
}

func (r FeedbackHub) GetProjectStatsForNotificationReminder(ctx *appcontext.AppContext, projectID string, timestamp time.Time, limit int64) (int64, []domain.Feedback, error) {
	resp, err := r.client.GetProjectStatsForNotificationReminder(ctx.Context(), &feedbackpb.GetProjectStatsForNotificationReminderRequest{
		TraceId:   ctx.GetTraceID(),
		ProjectId: projectID,
		Timestamp: timestamppb.New(timestamp),
		Limit:     limit,
	})
	if err != nil {
		return 0, nil, err
	}
	if resp.GetTotal() == 0 || resp.GetFeedbacks() == nil || len(resp.GetFeedbacks()) == 0 {
		return 0, nil, nil
	}

	var feedbacks = make([]domain.Feedback, 0)
	for _, f := range resp.GetFeedbacks() {
		var feedback = domain.Feedback{
			ID:           f.GetId(),
			Content:      f.GetContent(),
			Rating:       f.GetRating(),
			CampaignType: domain.ToProjectCampaignType(f.GetCampaignType()),
			CreatedAt:    f.GetCreatedAt().AsTime(),
		}

		var (
			appUserID = f.GetAppUserId()
			email     = f.GetEmail()
			project   = f.GetProject()
		)

		feedback.AppUserID = &appUserID
		feedback.Email = &email
		if project != nil {
			feedback.Project = domain.Project{
				ID:    project.GetId(),
				Title: project.GetTitle(),
			}
		}

		feedbacks = append(feedbacks, feedback)
	}

	return resp.GetTotal(), feedbacks, nil
}
