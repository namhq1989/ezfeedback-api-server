package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
)

type FeedbackReplyMapper struct{}

func (FeedbackReplyMapper) FromModelToDomain(feedbackReply model.FeedbackReplies) (*domain.FeedbackReply, error) {
	var result = &domain.FeedbackReply{
		ID:         feedbackReply.ID,
		FeedbackID: feedbackReply.FeedbackID,
		UserID:     feedbackReply.UserID,
		Content:    feedbackReply.Content,
		IsEdited:   feedbackReply.IsEdited,
		CreatedAt:  feedbackReply.CreatedAt,
		UpdatedAt:  feedbackReply.UpdatedAt,
	}

	return result, nil
}

func (FeedbackReplyMapper) FromDomainToModel(feedbackReply domain.FeedbackReply) (*model.FeedbackReplies, error) {
	var result = &model.FeedbackReplies{
		ID:         feedbackReply.ID,
		FeedbackID: feedbackReply.FeedbackID,
		UserID:     feedbackReply.UserID,
		Content:    feedbackReply.Content,
		IsEdited:   feedbackReply.IsEdited,
		CreatedAt:  feedbackReply.CreatedAt,
		UpdatedAt:  feedbackReply.UpdatedAt,
	}

	return result, nil
}
