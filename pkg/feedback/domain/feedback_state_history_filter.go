package domain

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/go-utilities/uuid"
)

const (
	feedbackStateHistoryQueryLimit int64 = 50
)

type FeedbackStateHistoryFilter struct {
	FeedbackID string
	Page       int64
	Limit      int64
}

func NewFeedbackStateHistoryFilter(feedbackID string, page int64) (*FeedbackStateHistoryFilter, error) {
	if !uuid.IsValidID(feedbackID) {
		return nil, apperrors.Feedback.InvalidFeedbackID
	}

	return &FeedbackStateHistoryFilter{
		FeedbackID: feedbackID,
		Page:       page,
		Limit:      feedbackStateHistoryQueryLimit,
	}, nil
}
