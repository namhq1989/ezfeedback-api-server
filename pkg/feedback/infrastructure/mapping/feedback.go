package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
)

type FeedbackMapper struct{}

func (FeedbackMapper) FromModelToDomain(feedback model.Feedbacks) (*domain.Feedback, error) {
	var result = &domain.Feedback{
		ID:                feedback.ID,
		ProjectID:         feedback.ProjectID,
		CampaignID:        feedback.CampaignID,
		AppUserID:         feedback.AppUserID,
		Email:             feedback.Email,
		CategoryID:        feedback.CategoryID,
		Content:           feedback.Content,
		Rating:            feedback.Rating,
		IsAnonymous:       feedback.IsAnonymous,
		State:             domain.ToFeedbackState(feedback.State.String()),
		CampaignType:      domain.ToProjectCampaignType(feedback.CampaignType.String()),
		Ip:                feedback.IP,
		CountryCode:       feedback.CountryCode,
		StatsTotalReplies: feedback.StatsTotalReplies,
		CreatedAt:         feedback.CreatedAt,
		UpdatedAt:         feedback.UpdatedAt,
	}

	return result, nil
}

func (FeedbackMapper) FromDomainToModel(feedback domain.Feedback) (*model.Feedbacks, error) {
	var result = &model.Feedbacks{
		ID:                feedback.ID,
		ProjectID:         feedback.ProjectID,
		CampaignID:        feedback.CampaignID,
		AppUserID:         feedback.AppUserID,
		Email:             feedback.Email,
		CategoryID:        feedback.CategoryID,
		Content:           feedback.Content,
		Rating:            feedback.Rating,
		IsAnonymous:       feedback.IsAnonymous,
		State:             model.FeedbackState(feedback.State.String()),
		CampaignType:      model.CampaignType(feedback.CampaignType.String()),
		IP:                feedback.Ip,
		CountryCode:       feedback.CountryCode,
		StatsTotalReplies: feedback.StatsTotalReplies,
		CreatedAt:         feedback.CreatedAt,
		UpdatedAt:         feedback.UpdatedAt,
	}

	return result, nil
}
