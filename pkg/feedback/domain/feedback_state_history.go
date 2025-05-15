package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/appcontext"
	"github.com/namhq1989/go-utilities/uuid"
)

type FeedbackStateHistoryRepository interface {
	Create(ctx *appcontext.AppContext, history FeedbackStateHistory) error
	FindWithFilter(ctx *appcontext.AppContext, filter FeedbackStateHistoryFilter) ([]FeedbackStateHistory, error)
}

type FeedbackStateHistory struct {
	ID         string
	FeedbackID string
	UserID     string
	State      FeedbackState
	CreatedAt  time.Time
}

func NewFeedbackStateHistory(feedbackID, userID string, state string) (*FeedbackStateHistory, error) {
	if !uuid.IsValidID(feedbackID) {
		return nil, apperrors.Feedback.InvalidFeedbackID
	}

	if !uuid.IsValidID(userID) {
		return nil, apperrors.User.InvalidUserID
	}

	dState := ToFeedbackState(state)
	if !dState.IsValid() {
		return nil, apperrors.Feedback.InvalidState
	}

	return &FeedbackStateHistory{
		ID:         uuid.New(),
		FeedbackID: feedbackID,
		UserID:     userID,
		State:      dState,
		CreatedAt:  manipulation.NowUTC(),
	}, nil
}
