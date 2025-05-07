package domain

import (
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/namhq1989/go-utilities/uuid"
)

const (
	feedbackReplyQueryLimit int64 = 20
)

type FeedbackReplyFilter struct {
	FeedbackID string
	Page       int64
	Limit      int64
}

func NewFeedbackReplyFilter(feedbackID string, page int64) (*FeedbackReplyFilter, error) {
	if !uuid.IsValidID(feedbackID) {
		return nil, apperrors.Feedback.InvalidFeedbackID
	}

	return &FeedbackReplyFilter{
		FeedbackID: feedbackID,
		Page:       page,
		Limit:      feedbackReplyQueryLimit,
	}, nil
}
