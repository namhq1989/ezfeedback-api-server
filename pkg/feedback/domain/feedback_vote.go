package domain

import (
	"time"

	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/ezfeedback-api-server/internal/utils/manipulation"
	"github.com/namhq1989/go-utilities/uuid"
)

type FeedbackVotes struct {
	ID         string
	FeedbackID string
	UserID     string
	CreatedAt  time.Time
}

func NewFeedbackVote(feedbackID, userID string) (*FeedbackVotes, error) {
	var v = &FeedbackVotes{
		ID:        uuid.New(),
		CreatedAt: manipulation.NowUTC(),
	}

	if err := v.SetFeedbackID(feedbackID); err != nil {
		return nil, err
	}
	if err := v.SetUserID(userID); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *FeedbackVotes) SetFeedbackID(feedbackID string) error {
	if !uuid.IsValidID(feedbackID) {
		return apperrors.Feedback.InvalidFeedbackID
	}

	v.FeedbackID = feedbackID
	return nil
}

func (v *FeedbackVotes) SetUserID(userID string) error {
	if !uuid.IsValidID(userID) {
		return apperrors.User.InvalidUserID
	}

	v.UserID = userID
	return nil
}
