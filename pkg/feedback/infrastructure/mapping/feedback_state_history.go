package mapping

import (
	"github.com/namhq1989/ezfeedback-api-server/internal/database/gen/ezfeedback/public/model"
	"github.com/namhq1989/ezfeedback-api-server/pkg/feedback/domain"
)

type FeedbackStateHistoryMapper struct{}

func (FeedbackStateHistoryMapper) FromModelToDomain(history model.FeedbackStateHistories) (*domain.FeedbackStateHistory, error) {
	var result = &domain.FeedbackStateHistory{
		ID:         history.ID,
		FeedbackID: history.FeedbackID,
		UserID:     history.UserID,
		State:      domain.ToFeedbackState(history.State.String()),
		CreatedAt:  history.CreatedAt,
	}

	return result, nil
}

func (FeedbackStateHistoryMapper) FromDomainToModel(history domain.FeedbackStateHistory) (*model.FeedbackStateHistories, error) {
	var result = &model.FeedbackStateHistories{
		ID:         history.ID,
		FeedbackID: history.FeedbackID,
		UserID:     history.UserID,
		State:      model.FeedbackState(history.State.String()),
		CreatedAt:  history.CreatedAt,
	}

	return result, nil
}
