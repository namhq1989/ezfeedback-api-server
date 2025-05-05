package grpc

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/genproto/feedbackpb"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
	"github.com/namhq1989/go-utilities/appcontext"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetProjectStatsForNotificationReminderHandler struct {
	feedbackRepository domain.FeedbackRepository
	projectHub         domain.ProjectHub
}

func NewGetProjectStatsForNotificationReminderHandler(feedbackRepository domain.FeedbackRepository, projectHub domain.ProjectHub) GetProjectStatsForNotificationReminderHandler {
	return GetProjectStatsForNotificationReminderHandler{
		feedbackRepository: feedbackRepository,
		projectHub:         projectHub,
	}
}

func (h GetProjectStatsForNotificationReminderHandler) GetProjectStatsForNotificationReminder(ctx *appcontext.AppContext, req *feedbackpb.GetProjectStatsForNotificationReminderRequest) (*feedbackpb.GetProjectStatsForNotificationReminderResponse, error) {
	ctx.SetTraceID(req.GetTraceId())
	ctx.Logger().Info("new get project stats for notification reminder request", appcontext.Fields{
		"projectId": req.GetProjectId(), "timestamp": req.GetTimestamp(), "limit": req.GetLimit(),
	})

	ctx.Logger().Text("count total created")
	totalFeedbacks, err := h.feedbackRepository.CountFeedbackForProjectSinceTimestamp(ctx, req.GetProjectId(), req.GetTimestamp().AsTime())
	if err != nil {
		ctx.Logger().Error("failed to count total created", err, appcontext.Fields{})
		return nil, err
	}
	if totalFeedbacks == 0 {
		ctx.Logger().Text("total created is 0, respond")
		return &feedbackpb.GetProjectStatsForNotificationReminderResponse{
			Total:     0,
			Feedbacks: make([]*feedbackpb.Feedback, 0),
		}, nil
	}

	ctx.Logger().Text("find feedbacks in db")
	feedbacks, err := h.feedbackRepository.FindFeedbackForProjectSinceTimestamp(ctx, req.GetProjectId(), req.GetTimestamp().AsTime(), req.GetLimit())
	if err != nil {
		ctx.Logger().Error("failed to find feedbacks in db", err, appcontext.Fields{})
		return nil, err
	}

	var result []*feedbackpb.Feedback
	for _, f := range feedbacks {
		var feedback = &feedbackpb.Feedback{
			Id:           f.ID,
			Content:      f.Content,
			Rating:       f.Rating,
			CampaignType: f.CampaignType.String(),
			CreatedAt:    timestamppb.New(f.CreatedAt),
		}
		if f.AppUserID != nil {
			feedback.AppUserId = *f.AppUserID
		}
		if f.Email != nil {
			feedback.Email = *f.Email
		}

		project, pErr := h.projectHub.GetProjectByID(ctx, f.ProjectID)
		if pErr != nil || project == nil {
			ctx.Logger().Error("failed to get project by id", pErr, appcontext.Fields{
				"projectId": f.ProjectID,
			})
			continue
		}
		feedback.Project = &feedbackpb.Project{
			Id:    project.ID,
			Title: project.Title,
		}

		result = append(result, feedback)
	}

	ctx.Logger().Text("done get project stats for notification reminder request")
	return &feedbackpb.GetProjectStatsForNotificationReminderResponse{
		Total:     totalFeedbacks,
		Feedbacks: result,
	}, nil
}
