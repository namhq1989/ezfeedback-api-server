package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
)

type FeedbackVoteMapper struct{}

func (FeedbackVoteMapper) FromModelToDomain(feedbackVote model.FeedbackVotes) (*domain.FeedbackVotes, error) {
	var result = &domain.FeedbackVotes{
		ID:         feedbackVote.ID,
		FeedbackID: feedbackVote.FeedbackID,
		UserID:     feedbackVote.UserID,
		CreatedAt:  feedbackVote.CreatedAt,
	}

	return result, nil
}

func (FeedbackVoteMapper) FromDomainToModel(feedbackVote domain.FeedbackVotes) (*model.FeedbackVotes, error) {
	var result = &model.FeedbackVotes{
		ID:         feedbackVote.ID,
		FeedbackID: feedbackVote.FeedbackID,
		UserID:     feedbackVote.UserID,
		CreatedAt:  feedbackVote.CreatedAt,
	}

	return result, nil
}
